package cmd

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/kinghanzala/gcpsec/internal/gcp"
	"github.com/kinghanzala/gcpsec/internal/output"
	"github.com/kinghanzala/gcpsec/internal/policy"
	"github.com/kinghanzala/gcpsec/internal/scanner"
)

type scanDeps struct {
	validateAuth func(context.Context) error
	newServices  func(context.Context) (*gcp.Services, error)
	loadConfig   func(string) (policy.Config, error)
	runChecks    func(context.Context, string, *gcp.Services) ([]scanner.CheckResult, error)
	render       func(io.Writer, output.ScanReport) error
	exitWith     func(int)
}

func newScanCmd() *cobra.Command {
	return newScanCmdWithDeps(scanDeps{
		validateAuth: gcp.ValidateADC,
		newServices:  gcp.NewServices,
		loadConfig:   policy.LoadConfig,
		runChecks:    scanner.RunAll,
		render:       output.RenderScan,
		exitWith:     os.Exit,
	})
}

func newScanCmdWithDeps(deps scanDeps) *cobra.Command {
	var projectID string
	var configPath string

	cmd := &cobra.Command{
		Use:   "scan --project <id>",
		Short: "Scan a GCP project for insecure configuration",
		RunE: func(cmd *cobra.Command, _ []string) error {
			if projectID == "" {
				return errors.New("required flag(s) \"project\" not set")
			}

			ctx := cmd.Context()
			if err := deps.validateAuth(ctx); err != nil {
				return err
			}

			services, err := deps.newServices(ctx)
			if err != nil {
				return err
			}
			defer services.Close()

			cfg, err := deps.loadConfig(configPath)
			if err != nil {
				return err
			}

			results, err := deps.runChecks(ctx, projectID, services)
			if err != nil {
				return err
			}

			evaluated := policy.Apply(results, projectID, cfg)
			report := output.NewScanReport(projectID, evaluated)
			if err := deps.render(cmd.OutOrStdout(), report); err != nil {
				return err
			}

			if report.HasHigh {
				deps.exitWith(1)
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&projectID, "project", "", "GCP project ID to scan")
	defaultPath := filepath.Join(".", "config.yaml")
	cmd.Flags().StringVar(&configPath, "config", defaultPath, "Optional path to a policy config file")

	return cmd
}
