package cmd

import (
	"fmt"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

func newUninstallInfoCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "uninstall-info",
		Short: "Show how to uninstall gcpsec from the current machine",
		RunE: func(cmd *cobra.Command, _ []string) error {
			installedPath, err := resolveInstalledPath()
			if err != nil {
				_, writeErr := fmt.Fprintln(cmd.OutOrStdout(), "gcpsec was not found on PATH.")
				if writeErr != nil {
					return writeErr
				}
				_, writeErr = fmt.Fprintln(cmd.OutOrStdout(), "If you installed it with Go, try:")
				if writeErr != nil {
					return writeErr
				}
				_, writeErr = fmt.Fprintln(cmd.OutOrStdout(), `rm -f "$(go env GOPATH)/bin/gcpsec"`)
				return writeErr
			}

			_, err = fmt.Fprintf(cmd.OutOrStdout(), "Installed binary: %s\n", installedPath)
			if err != nil {
				return err
			}
			_, err = fmt.Fprintf(cmd.OutOrStdout(), "Uninstall with: rm -f %q\n", installedPath)
			return err
		},
	}
}

func resolveInstalledPath() (string, error) {
	path, err := exec.LookPath("gcpsec")
	if err != nil {
		return "", err
	}
	return filepath.Clean(path), nil
}
