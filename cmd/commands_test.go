package cmd

import (
	"bytes"
	"context"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/kinghanzala/gcpsec/internal/gcp"
	"github.com/kinghanzala/gcpsec/internal/model"
	"github.com/kinghanzala/gcpsec/internal/output"
	"github.com/kinghanzala/gcpsec/internal/policy"
	"github.com/kinghanzala/gcpsec/internal/scanner"
)

func TestScanRequiresProject(t *testing.T) {
	cmd := newScanCmdWithDeps(scanDeps{
		validateAuth: func(context.Context) error { return nil },
		newServices:  func(context.Context) (*gcp.Services, error) { return &gcp.Services{}, nil },
		loadConfig:   func(string) (policy.Config, error) { return policy.DefaultConfig(), nil },
		runChecks:    func(context.Context, string, *gcp.Services) ([]scanner.CheckResult, error) { return nil, nil },
		render:       func(_ io.Writer, _ output.ScanReport) error { return nil },
		exitWith:     func(int) {},
	})
	var stdout bytes.Buffer
	cmd.SetOut(&stdout)
	cmd.SetErr(&stdout)
	cmd.SetArgs([]string{})

	err := cmd.Execute()
	if err == nil || !strings.Contains(err.Error(), "required flag(s) \"project\" not set") {
		t.Fatalf("expected missing project flag error, got %v", err)
	}
}

func TestDoctorReportsAuthFailure(t *testing.T) {
	cmd := newDoctorCmdWithDeps(doctorDeps{
		validateAuth: func(context.Context) error { return errors.New("missing ADC") },
		newServices:  func(context.Context) (*gcp.Services, error) { return &gcp.Services{}, nil },
	})
	cmd.SetArgs([]string{})

	err := cmd.Execute()
	if err == nil || !strings.Contains(err.Error(), "missing ADC") {
		t.Fatalf("expected auth error, got %v", err)
	}
}

func TestVersionPrintsBuildInfo(t *testing.T) {
	cmd := newVersionCmd(BuildInfo{
		Version: "1.2.3",
		Commit:  "abc123",
		Date:    "2026-04-03",
	})
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	got := buf.String()
	if !strings.Contains(got, "version=1.2.3") {
		t.Fatalf("expected version in output, got %q", got)
	}
}

func TestScanExitsOnHighFindings(t *testing.T) {
	exitCode := 0
	cmd := newScanCmdWithDeps(scanDeps{
		validateAuth: func(context.Context) error { return nil },
		newServices:  func(context.Context) (*gcp.Services, error) { return &gcp.Services{}, nil },
		loadConfig:   func(string) (policy.Config, error) { return policy.DefaultConfig(), nil },
		runChecks: func(context.Context, string, *gcp.Services) ([]scanner.CheckResult, error) {
			return []scanner.CheckResult{
				{
					Name: "storage",
					Findings: []model.Finding{
						{
							ID:       "1",
							Check:    "PUBLIC_BUCKET",
							RuleID:   "public_buckets",
							Resource: "bucket-1",
							Severity: model.SeverityHigh,
							Message:  "Public bucket: bucket-1",
						},
					},
				},
			}, nil
		},
		render: func(_ io.Writer, _ output.ScanReport) error { return nil },
		exitWith: func(code int) {
			exitCode = code
		},
	})
	cmd.SetArgs([]string{"--project", "demo-project"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if exitCode != 1 {
		t.Fatalf("expected exit code 1, got %d", exitCode)
	}
}
