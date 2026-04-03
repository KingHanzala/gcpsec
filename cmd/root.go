package cmd

import "github.com/spf13/cobra"

type BuildInfo struct {
	Version string
	Commit  string
	Date    string
}

func NewRootCmd(build BuildInfo) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:           "gcpsec",
		Short:         "Developer-first GCP security scanner",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	rootCmd.AddCommand(newScanCmd())
	rootCmd.AddCommand(newDoctorCmd())
	rootCmd.AddCommand(newUninstallInfoCmd())
	rootCmd.AddCommand(newVersionCmd(build))

	return rootCmd
}
