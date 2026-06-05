// cmd/validate.go
// Implements the validate command.

package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/Vamshavardhan50/envguard/internal/parser"
	"github.com/Vamshavardhan50/envguard/internal/reporter"
	"github.com/Vamshavardhan50/envguard/internal/validator"
	"github.com/spf13/cobra"
)

func newValidateCommand() *cobra.Command {
	var (
		envPath    string
		configPath string
		format     string
		ciMode     bool
	)

	cmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate env vars against schema",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := parser.ParseConfig(configPath)
			if err != nil {
				return err
			}

			vars, err := parser.ParseDotEnv(envPath)
			if err != nil {
				return fmt.Errorf("validate: %w", err)
			}

			engine := &validator.Engine{}
			result := engine.Validate(vars, cfg.Rules)

			if err := renderValidate(cmd, format, result); err != nil {
				return err
			}

			if ciMode && len(result.Invalid) > 0 {
				os.Exit(1)
			}
			if !ciMode && cfg.Output.ExitCode && len(result.Invalid) > 0 {
				os.Exit(1)
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&envPath, "env", ".env", "Path to .env file")
	cmd.Flags().StringVar(&configPath, "config", ".envguard.yaml", "Path to config file")
	cmd.Flags().StringVar(&format, "format", "pretty", "Output format: pretty|json")
	cmd.Flags().BoolVar(&ciMode, "ci", false, "Exit 1 on validation errors")

	return cmd
}

func renderValidate(cmd *cobra.Command, format string, result validator.ValidationResult) error {
	switch strings.ToLower(format) {
	case "pretty":
		out := (&reporter.PrettyReporter{}).RenderValidate(result)
		_, err := fmt.Fprintln(cmd.OutOrStdout(), out)
		return err
	case "json":
		out, err := (&reporter.JSONReporter{}).RenderValidate(result)
		if err != nil {
			return err
		}
		_, err = fmt.Fprintln(cmd.OutOrStdout(), out)
		return err
	default:
		return fmt.Errorf("unknown format: %s", format)
	}
}
