package policy

import (
	"testing"
	"time"

	"github.com/kinghanzala/gcpsec/internal/model"
	"github.com/kinghanzala/gcpsec/internal/scanner"
)

func TestApplyMatchesWildcardException(t *testing.T) {
	results := []scanner.CheckResult{
		{
			Name: "storage",
			Findings: []model.Finding{
				{
					ID:        "1",
					Check:     "PUBLIC_BUCKET",
					RuleID:    "public_buckets",
					Resource:  "assets-prod",
					Severity:  model.SeverityHigh,
					Timestamp: time.Now(),
				},
			},
		},
	}
	cfg := DefaultConfig()
	cfg.Rules["public_buckets"] = Rule{
		Severity:   model.SeverityHigh,
		Exceptions: []string{"assets-*"},
	}

	findings := Apply(results, "demo-project", cfg)
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if !findings[0].Allowed {
		t.Fatalf("expected finding to be allowed")
	}
	if findings[0].Severity != model.SeverityInfo {
		t.Fatalf("expected allowed finding severity to become INFO")
	}
}
