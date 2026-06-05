// cmd/audit.go
// Implements the audit command.

package cmd

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/Vamshavardhan50/envguard/internal/parser"
	"github.com/Vamshavardhan50/envguard/internal/reporter"
	"github.com/Vamshavardhan50/envguard/internal/scanner"
	"github.com/Vamshavardhan50/envguard/pkg/envfile"
	"github.com/spf13/cobra"
)

func newAuditCommand() *cobra.Command {
	var (
		flagPath   string
		envPath    string
		configPath string
		format     string
		ciMode     bool
		ignore     string
	)

	cmd := &cobra.Command{
		Use:   "audit [path]",
		Short: "Audit environment variable usage",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 && cmd.Flags().Changed("path") {
				return fmt.Errorf("audit: specify either [path] or --path, not both")
			}

			cfg, err := parser.ParseConfig(configPath)
			if err != nil {
				return err
			}

			paths := resolveAuditPaths(cmd, args, flagPath, cfg)
			ignoreList := mergeStringSlices(cfg.Scan.Ignore, parseCommaList(ignore))
			languages := cfg.Scan.Languages
			if len(languages) == 0 {
				languages = []string{"auto"}
			}

			vars, err := parser.ParseDotEnv(envPath)
			if err != nil {
				return fmt.Errorf("audit: %w", err)
			}

			customPatterns := make([]scanner.Pattern, 0, len(cfg.Scan.CustomPatterns))
			for _, cp := range cfg.Scan.CustomPatterns {
				rx, err := regexp.Compile(cp.Regex)
				if err != nil {
					return fmt.Errorf("audit: compile custom regex %q: %w", cp.Regex, err)
				}
				customPatterns = append(customPatterns, scanner.Pattern{
					Language: cp.Language,
					Regex:    rx,
				})
			}

			scannerPatterns := append(scanner.Patterns(), customPatterns...)
			scannerEngine := &scanner.FileScanner{Patterns: scannerPatterns}
			allUsages := make([]scanner.VarUsage, 0, 128)
			totalFiles := 0
			start := time.Now()
			for _, path := range paths {
				result, err := scannerEngine.Scan(path, ignoreList, languages)
				if err != nil {
					return fmt.Errorf("audit: %w", err)
				}
				totalFiles += result.ScannedFiles
				allUsages = append(allUsages, result.Usages...)
			}
			duration := time.Since(start).Milliseconds()

			report := buildAuditResult(allUsages, vars, totalFiles, duration)

			if !cmd.Flags().Changed("format") {
				format = cfg.Output.Format
			}
			if format == "" {
				format = "pretty"
			}

			if err := renderAudit(cmd, format, report); err != nil {
				return err
			}

			if ciMode && (report.ErrorCount > 0 || report.WarnCount > 0) {
				os.Exit(1)
			}
			if !ciMode && cfg.Output.ExitCode && report.ErrorCount > 0 {
				os.Exit(1)
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&flagPath, "path", ".", "Path to scan")
	cmd.Flags().StringVar(&envPath, "env", ".env", "Path to .env file")
	cmd.Flags().StringVar(&configPath, "config", ".envguard.yaml", "Path to config file")
	cmd.Flags().StringVar(&format, "format", "pretty", "Output format: pretty|json|github")
	cmd.Flags().BoolVar(&ciMode, "ci", false, "Exit 1 on any findings")
	cmd.Flags().StringVar(&ignore, "ignore", "", "Comma-separated paths to ignore")

	return cmd
}

func resolveAuditPaths(cmd *cobra.Command, args []string, flagPath string, cfg parser.Config) []string {
	if len(args) > 0 {
		return []string{args[0]}
	}
	if cmd.Flags().Changed("path") {
		return []string{flagPath}
	}
	if len(cfg.Scan.Paths) > 0 {
		return cfg.Scan.Paths
	}
	return []string{"."}
}

func renderAudit(cmd *cobra.Command, format string, report reporter.AuditResult) error {
	switch strings.ToLower(format) {
	case "pretty":
		out := (&reporter.PrettyReporter{}).RenderAudit(report)
		_, err := fmt.Fprintln(cmd.OutOrStdout(), out)
		return err
	case "json":
		out, err := (&reporter.JSONReporter{}).RenderAudit(report)
		if err != nil {
			return err
		}
		_, err = fmt.Fprintln(cmd.OutOrStdout(), out)
		return err
	case "github":
		out := (&reporter.GitHubReporter{}).RenderAudit(report)
		_, err := fmt.Fprintln(cmd.OutOrStdout(), out)
		return err
	default:
		return fmt.Errorf("unknown format: %s", format)
	}
}

func buildAuditResult(usages []scanner.VarUsage, vars []envfile.EnvVar, scannedFiles int, durationMs int64) reporter.AuditResult {
	usageByKey := make(map[string][]scanner.VarUsage)
	for _, usage := range usages {
		usageByKey[usage.Key] = append(usageByKey[usage.Key], usage)
	}

	usedKeys := make(map[string]struct{}, len(usageByKey))
	for key := range usageByKey {
		usedKeys[key] = struct{}{}
	}

	envKeys := envfile.Keys(vars)
	envSet := make(map[string]struct{}, len(envKeys))
	for _, key := range envKeys {
		envSet[key] = struct{}{}
	}

	missing := make([]reporter.MissingVar, 0)
	for key, usedIn := range usageByKey {
		if _, ok := envSet[key]; !ok {
			sort.Slice(usedIn, func(i, j int) bool {
				if usedIn[i].File == usedIn[j].File {
					return usedIn[i].Line < usedIn[j].Line
				}
				return usedIn[i].File < usedIn[j].File
			})
			missing = append(missing, reporter.MissingVar{Key: key, UsedIn: usedIn})
		}
	}

	unused := make([]reporter.UnusedVar, 0)
	for _, key := range envKeys {
		if _, ok := usedKeys[key]; !ok {
			unused = append(unused, reporter.UnusedVar{Key: key})
		}
	}

	valid := make([]string, 0)
	for key := range usedKeys {
		if _, ok := envSet[key]; ok {
			valid = append(valid, key)
		}
	}

	sort.Slice(missing, func(i, j int) bool { return missing[i].Key < missing[j].Key })
	sort.Slice(unused, func(i, j int) bool { return unused[i].Key < unused[j].Key })
	sort.Strings(valid)

	return reporter.AuditResult{
		ScannedFiles: scannedFiles,
		DurationMs:   durationMs,
		Missing:      missing,
		Unused:       unused,
		Valid:        valid,
		ErrorCount:   len(missing),
		WarnCount:    len(unused),
	}
}

func parseCommaList(input string) []string {
	input = strings.TrimSpace(input)
	if input == "" {
		return nil
	}
	parts := strings.Split(input, ",")
	items := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		items = append(items, part)
	}
	return items
}

func mergeStringSlices(first, second []string) []string {
	if len(first) == 0 && len(second) == 0 {
		return nil
	}
	merged := make([]string, 0, len(first)+len(second))
	merged = append(merged, first...)
	merged = append(merged, second...)
	return merged
}
