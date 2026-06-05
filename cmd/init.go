// cmd/init.go
// Implements the init command.

package cmd

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/Vamshavardhan50/envguard/internal/scanner"
	"github.com/Vamshavardhan50/envguard/pkg/envfile"
	"github.com/spf13/cobra"
)

func newInitCommand() *cobra.Command {
	var (
		yes  bool
		path string
	)

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize envguard configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			basePath := strings.TrimSpace(path)
			if basePath == "" {
				basePath = "."
			}

			configPath := filepath.Join(basePath, ".envguard.yaml")
			exists, err := fileExists(configPath)
			if err != nil {
				return fmt.Errorf("init: %w", err)
			}
			if exists && !yes {
				ok, err := confirmOverwrite(cmd.InOrStdin(), cmd.OutOrStdout(), configPath)
				if err != nil {
					return fmt.Errorf("init: %w", err)
				}
				if !ok {
					return nil
				}
			}

			language := detectPrimaryLanguage(basePath)
			rules, err := loadEnvRules(basePath)
			if err != nil {
				return fmt.Errorf("init: %w", err)
			}

			configBody := buildConfigYAML(language, rules)
			if err := os.WriteFile(configPath, []byte(configBody), 0o600); err != nil {
				return fmt.Errorf("init: write config: %w", err)
			}

			if err := ensureGitignoreEntry(basePath, ".env"); err != nil {
				return fmt.Errorf("init: %w", err)
			}

			_, err = fmt.Fprintln(cmd.OutOrStdout(), "Created "+configPath)
			return err
		},
	}

	cmd.Flags().BoolVar(&yes, "yes", false, "Skip prompts and use defaults")
	cmd.Flags().StringVar(&path, "path", ".", "Path to initialize")

	return cmd
}

func detectPrimaryLanguage(root string) string {
	counts := make(map[string]int)
	_ = filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			name := strings.ToLower(entry.Name())
			if name == ".git" || name == "node_modules" || name == "vendor" || name == "dist" || name == "build" {
				return filepath.SkipDir
			}
			return nil
		}
		lang := scanner.DetectLanguage(path)
		if lang != "" {
			counts[lang]++
		}
		return nil
	})

	var (
		bestLang string
		bestCnt  int
	)
	for lang, cnt := range counts {
		if cnt > bestCnt {
			bestLang = lang
			bestCnt = cnt
		}
	}
	if bestLang == "" {
		return "auto"
	}
	return bestLang
}

func loadEnvRules(root string) (map[string]envfile.EnvVar, error) {
	envPath := filepath.Join(root, ".env")
	exists, err := fileExists(envPath)
	if err != nil {
		return nil, err
	}
	if !exists {
		return map[string]envfile.EnvVar{}, nil
	}

	vars, err := envfile.ParseFile(envPath)
	if err != nil {
		return nil, err
	}

	indexed := make(map[string]envfile.EnvVar, len(vars))
	for _, entry := range vars {
		if entry.Key == "" {
			continue
		}
		indexed[entry.Key] = entry
	}
	return indexed, nil
}

func buildConfigYAML(language string, rules map[string]envfile.EnvVar) string {
	if language == "" {
		language = "auto"
	}

	keys := make([]string, 0, len(rules))
	for key := range rules {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	var b strings.Builder
	b.WriteString("version: 1\n\n")
	b.WriteString("scan:\n")
	b.WriteString("  paths:\n")
	b.WriteString("    - \".\"\n")
	b.WriteString("  ignore:\n")
	b.WriteString("    - \"node_modules\"\n")
	b.WriteString("    - \".git\"\n")
	b.WriteString("    - \"dist\"\n")
	b.WriteString("    - \"build\"\n")
	b.WriteString("    - \"vendor\"\n")
	b.WriteString("    - \".goreleaser.yaml\"\n")
	b.WriteString("  languages:\n")
	b.WriteString("    - \"" + yamlEscape(language) + "\"\n\n")

	b.WriteString("rules:\n")
	if len(keys) == 0 {
		b.WriteString("  {}\n")
	} else {
		for _, key := range keys {
			rule := defaultRuleForKey(key)
			b.WriteString("  " + key + ":\n")
			b.WriteString("    required: " + formatBool(rule.Required) + "\n")
			b.WriteString("    type: " + rule.Type + "\n")
			if len(rule.Values) > 0 {
				b.WriteString("    values:\n")
				for _, v := range rule.Values {
					b.WriteString("      - " + yamlEscape(v) + "\n")
				}
			}
			if rule.Default != "" {
				b.WriteString("    default: \"" + yamlEscape(rule.Default) + "\"\n")
			}
			if rule.Description != "" {
				b.WriteString("    description: \"" + yamlEscape(rule.Description) + "\"\n")
			}
		}
	}

	b.WriteString("\noutput:\n")
	b.WriteString("  format: pretty\n")
	b.WriteString("  exit_code: true\n")
	return b.String()
}

type ruleTemplate struct {
	Required    bool
	Type        string
	Values      []string
	Default     string
	Description string
}

func defaultRuleForKey(key string) ruleTemplate {
	upper := strings.ToUpper(key)

	switch upper {
	case "DATABASE_URL":
		return ruleTemplate{Required: true, Type: "url", Description: "Primary database connection string"}
	case "PORT":
		return ruleTemplate{Required: false, Type: "number", Default: "3000", Description: "Server port"}
	case "NODE_ENV":
		return ruleTemplate{Required: true, Type: "enum", Values: []string{"development", "production", "test"}, Description: "Application environment"}
	}

	if strings.Contains(upper, "URL") {
		return ruleTemplate{Required: true, Type: "url"}
	}
	if strings.HasSuffix(upper, "_PORT") {
		return ruleTemplate{Required: true, Type: "number"}
	}
	if strings.Contains(upper, "DEBUG") || strings.HasSuffix(upper, "_ENABLED") {
		return ruleTemplate{Required: true, Type: "boolean"}
	}

	return ruleTemplate{Required: true, Type: "string"}
}

func yamlEscape(value string) string {
	replacer := strings.NewReplacer("\\", "\\\\", "\"", "\\\"")
	return replacer.Replace(value)
}

func formatBool(value bool) string {
	if value {
		return "true"
	}
	return "false"
}

func ensureGitignoreEntry(root, entry string) error {
	path := filepath.Join(root, ".gitignore")
	exists, err := fileExists(path)
	if err != nil {
		return err
	}

	if !exists {
		return os.WriteFile(path, []byte(entry+"\n"), 0o600)
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) == entry {
			return nil
		}
	}

	lines = append(lines, entry)
	output := strings.Join(lines, "\n")
	if !strings.HasSuffix(output, "\n") {
		output += "\n"
	}
	return os.WriteFile(path, []byte(output), 0o600)
}
