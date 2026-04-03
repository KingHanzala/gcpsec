package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newVersionCmd(build BuildInfo) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print build version information",
		RunE: func(cmd *cobra.Command, _ []string) error {
			_, err := fmt.Fprintf(cmd.OutOrStdout(), "gcpsec version=%s commit=%s date=%s\n", build.Version, build.Commit, build.Date)
			return err
		},
	}
}
