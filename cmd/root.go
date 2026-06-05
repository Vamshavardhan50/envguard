// cmd/root.go
// Defines the root cobra command for envguard.

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Execute runs the envguard CLI and exits with an appropriate status code.
func Execute(version, commit, date string) {
	rootCmd := newRootCommand(version, commit, date)
	if err := rootCmd.Execute(); err != nil {
		os.Exit(2)
	}
}

func newRootCommand(version, commit, date string) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:           "envguard",
		Short:         "Audit and validate environment variables",
		Long:          "envguard is a fast CLI that audits, validates, and protects environment variable configuration.",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	rootCmd.Version = buildVersion(version, commit, date)
	rootCmd.SetVersionTemplate("envguard version {{.Version}}\n")

	rootCmd.AddCommand(
		newAuditCommand(),
		newSyncCommand(),
		newDiffCommand(),
		newDoctorCommand(),
		newInitCommand(),
		newValidateCommand(),
	)

	return rootCmd
}

func buildVersion(version, commit, date string) string {
	return fmt.Sprintf("%s (commit %s, date %s)", version, commit, date)
}
