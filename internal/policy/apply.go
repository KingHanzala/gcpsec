package policy

import (
	"path/filepath"

	"github.com/kinghanzala/gcpsec/internal/model"
	"github.com/kinghanzala/gcpsec/internal/scanner"
)

func Apply(results []scanner.CheckResult, projectID string, cfg Config) []model.Finding {
	findings := make([]model.Finding, 0)
	for _, result := range results {
		for _, finding := range result.Findings {
			finding.Project = projectID
			if rule, ok := cfg.Rules[finding.RuleID]; ok {
				if rule.Severity != "" {
					finding.Severity = rule.Severity
				}
				for _, pattern := range rule.Exceptions {
					if matches(pattern, finding.Resource) {
						finding.Allowed = true
						finding.AllowedReason = pattern
						finding.Severity = model.SeverityInfo
						break
					}
				}
			}
			findings = append(findings, finding)
		}
	}
	return findings
}

func matches(pattern, resource string) bool {
	ok, err := filepath.Match(pattern, resource)
	if err != nil {
		return false
	}
	return ok
}
