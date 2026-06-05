// cmd/diff.go
// Implements the diff command.

package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/Vamshavardhan50/envguard/internal/differ"
	"github.com/Vamshavardhan50/envguard/internal/parser"
	"github.com/Vamshavardhan50/envguard/internal/reporter"
	"github.com/spf13/cobra"
)

func newDiffCommand() *cobra.Command {
	var (
		format string
		ciMode bool
	)

	cmd := &cobra.Command{
		Use:   "diff <file1> <file2>",
		Short: "Compare two env files",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			firstVars, err := parser.ParseDotEnv(args[0])
			if err != nil {
				return fmt.Errorf("diff: %w", err)
			}
			secondVars, err := parser.ParseDotEnv(args[1])
			if err != nil {
				return fmt.Errorf("diff: %w", err)
			}

			engine := &differ.Engine{}
			result := engine.Diff(firstVars, secondVars)

			if err := renderDiff(cmd, format, result); err != nil {
				return err
			}

			if ciMode && (len(result.OnlyInFirst) > 0 || len(result.OnlyInSecond) > 0) {
				os.Exit(1)
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&format, "format", "pretty", "Output format: pretty|json")
	cmd.Flags().BoolVar(&ciMode, "ci", false, "Exit 1 if any differences")

	return cmd
}

func renderDiff(cmd *cobra.Command, format string, result differ.DiffResult) error {
	switch strings.ToLower(format) {
	case "pretty":
		out := (&reporter.PrettyReporter{}).RenderDiff(result)
		_, err := fmt.Fprintln(cmd.OutOrStdout(), out)
		return err
	case "json":
		out, err := (&reporter.JSONReporter{}).RenderDiff(result)
		if err != nil {
			return err
		}
		_, err = fmt.Fprintln(cmd.OutOrStdout(), out)
		return err
	default:
		return fmt.Errorf("unknown format: %s", format)
	}
}
