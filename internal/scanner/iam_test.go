package scanner

import (
	"context"
	"testing"

	"github.com/kinghanzala/gcpsec/internal/gcp"
)

type fakeIAMService struct {
	bindings []gcp.IAMBinding
	err      error
}

func (f fakeIAMService) GetProjectIAMBindings(context.Context, string) ([]gcp.IAMBinding, error) {
	return f.bindings, f.err
}

func TestIAMCheckFindings(t *testing.T) {
	check := NewIAMCheck(fakeIAMService{
		bindings: []gcp.IAMBinding{
			{Role: "roles/owner", Members: []string{"user:admin@example.com"}},
			{Role: "roles/viewer", Members: []string{"allUsers"}},
		},
	})

	findings, err := check.Run(context.Background(), "demo-project")
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if len(findings) != 2 {
		t.Fatalf("expected 2 findings, got %d", len(findings))
	}
	if findings[0].RuleID != "owner_roles" {
		t.Fatalf("unexpected first rule id: %s", findings[0].RuleID)
	}
	if findings[1].RuleID != "public_iam_members" {
		t.Fatalf("unexpected second rule id: %s", findings[1].RuleID)
	}
}
