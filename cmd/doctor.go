package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/kinghanzala/gcpsec/internal/gcp"
)

type doctorDeps struct {
	validateAuth func(context.Context) error
	newServices  func(context.Context) (*gcp.Services, error)
}

func newDoctorCmd() *cobra.Command {
	return newDoctorCmdWithDeps(doctorDeps{
		validateAuth: gcp.ValidateADC,
		newServices:  gcp.NewServices,
	})
}

func newDoctorCmdWithDeps(deps doctorDeps) *cobra.Command {
	return &cobra.Command{
		Use:   "doctor",
		Short: "Validate local GCP scanner prerequisites",
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := cmd.Context()
			if err := deps.validateAuth(ctx); err != nil {
				return err
			}

			services, err := deps.newServices(ctx)
			if err != nil {
				return err
			}
			defer services.Close()

			fmt.Fprintln(cmd.OutOrStdout(), "OK Authenticated with Application Default Credentials")
			fmt.Fprintln(cmd.OutOrStdout(), "OK GCP client initialization succeeded")
			fmt.Fprintln(cmd.OutOrStdout(), "OK Ready to scan")
			return nil
		},
	}
}
