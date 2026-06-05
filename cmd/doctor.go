// cmd/doctor.go
// Implements the doctor command.

package cmd

import (
	"fmt"
	"strings"

	"github.com/Vamshavardhan50/envguard/internal/doctor"
	"github.com/Vamshavardhan50/envguard/internal/reporter"
	"github.com/spf13/cobra"
)

func newDoctorCommand() *cobra.Command {
	var (
		path   string
		format string
	)

	cmd := &cobra.Command{
		Use:   "doctor [path]",
		Short: "Run project health checks",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				path = args[0]
			}
			if path == "" {
				path = "."
			}

			engine := &doctor.Engine{}
			result := engine.Check(path)

			return renderDoctor(cmd, format, result)
		},
	}

	cmd.Flags().StringVar(&path, "path", ".", "Path to check")
	cmd.Flags().StringVar(&format, "format", "pretty", "Output format: pretty|json")

	return cmd
}

func renderDoctor(cmd *cobra.Command, format string, result doctor.DoctorResult) error {
	switch strings.ToLower(format) {
	case "pretty":
		out := (&reporter.PrettyReporter{}).RenderDoctor(result)
		_, err := fmt.Fprintln(cmd.OutOrStdout(), out)
		return err
	case "json":
		out, err := (&reporter.JSONReporter{}).RenderDoctor(result)
		if err != nil {
			return err
		}
		_, err = fmt.Fprintln(cmd.OutOrStdout(), out)
		return err
	default:
		return fmt.Errorf("unknown format: %s", format)
	}
}
