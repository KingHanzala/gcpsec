package scanner

import (
	"context"
	"testing"

	"github.com/kinghanzala/gcpsec/internal/gcp"
)

type fakeComputeService struct {
	rules []gcp.FirewallRule
	err   error
}

func (f fakeComputeService) ListFirewallRules(context.Context, string) ([]gcp.FirewallRule, error) {
	return f.rules, f.err
}

func TestNetworkCheckFindings(t *testing.T) {
	check := NewNetworkCheck(fakeComputeService{
		rules: []gcp.FirewallRule{
			{
				Name:         "allow-ssh-and-rdp",
				Direction:    "INGRESS",
				SourceRanges: []string{"0.0.0.0/0"},
				Allowed: []gcp.FirewallAllowed{
					{IPProtocol: "tcp", Ports: []string{"22", "3389"}},
				},
			},
		},
	})

	findings, err := check.Run(context.Background(), "demo-project")
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if len(findings) != 2 {
		t.Fatalf("expected 2 findings, got %d", len(findings))
	}
	if findings[0].RuleID != "open_ssh" {
		t.Fatalf("unexpected first rule id: %s", findings[0].RuleID)
	}
	if findings[1].RuleID != "open_rdp" {
		t.Fatalf("unexpected second rule id: %s", findings[1].RuleID)
	}
}
