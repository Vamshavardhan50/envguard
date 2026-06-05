// cmd/sync.go
// Implements the sync command.

package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/Vamshavardhan50/envguard/internal/parser"
	"github.com/Vamshavardhan50/envguard/pkg/envfile"
	"github.com/spf13/cobra"
)

func newSyncCommand() *cobra.Command {
	var (
		envPath    string
		outputPath string
		force      bool
	)

	cmd := &cobra.Command{
		Use:   "sync",
		Short: "Generate .env.example from .env",
		RunE: func(cmd *cobra.Command, args []string) error {
			vars, err := parser.ParseDotEnv(envPath)
			if err != nil {
				return fmt.Errorf("sync: %w", err)
			}

			cfg, err := parser.ParseConfig(defaultConfigPath())
			if err != nil {
				return fmt.Errorf("sync: %w", err)
			}

			content := buildEnvExample(vars, cfg)

			if !force {
				exists, err := fileExists(outputPath)
				if err != nil {
					return fmt.Errorf("sync: %w", err)
				}
				if exists {
					ok, err := confirmOverwrite(cmd.InOrStdin(), cmd.OutOrStdout(), outputPath)
					if err != nil {
						return fmt.Errorf("sync: %w", err)
					}
					if !ok {
						return nil
					}
				}
			}

			cleanPath := filepath.Clean(outputPath)
			if err := os.WriteFile(cleanPath, []byte(content), 0o600); err != nil {
				return fmt.Errorf("sync: write output: %w", err)
			}

			_, err = fmt.Fprintln(cmd.OutOrStdout(), "Wrote "+cleanPath)
			return err
		},
	}

	cmd.Flags().StringVar(&envPath, "env", ".env", "Source .env file")
	cmd.Flags().StringVar(&outputPath, "output", ".env.example", "Output file")
	cmd.Flags().BoolVar(&force, "force", false, "Overwrite without confirmation")

	return cmd
}

func buildEnvExample(vars []envfile.EnvVar, cfg parser.Config) string {
	keys := envfile.Keys(vars)
	if len(keys) == 0 {
		return ""
	}

	sort.Strings(keys)
	descriptions := make(map[string]string, len(cfg.Rules))
	for key, rule := range cfg.Rules {
		if rule.Description == "" {
			continue
		}
		descriptions[key] = rule.Description
	}

	var b strings.Builder
	for _, key := range keys {
		line := key + "="
		if desc, ok := descriptions[key]; ok {
			line += " # " + desc
		}
		b.WriteString(line)
		b.WriteString("\n")
	}
	return b.String()
}

func confirmOverwrite(in io.Reader, out io.Writer, path string) (bool, error) {
	_, err := fmt.Fprintf(out, "Overwrite %s? (y/N): ", path)
	if err != nil {
		return false, err
	}

	reader := bufio.NewReader(in)
	answer, err := reader.ReadString('\n')
	if err != nil {
		return false, err
	}
	answer = strings.TrimSpace(strings.ToLower(answer))
	return answer == "y" || answer == "yes", nil
}

func fileExists(path string) (bool, error) {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return !info.IsDir(), nil
}

func defaultConfigPath() string {
	return ".envguard.yaml"
}
