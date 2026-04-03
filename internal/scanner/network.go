package scanner

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/kinghanzala/gcpsec/internal/gcp"
	"github.com/kinghanzala/gcpsec/internal/model"
)

type NetworkCheck struct {
	service gcp.ComputeService
	now     func() time.Time
}

func NewNetworkCheck(service gcp.ComputeService) NetworkCheck {
	return NetworkCheck{
		service: service,
		now:     time.Now().UTC,
	}
}

func (c NetworkCheck) Name() string {
	return "network"
}

func (c NetworkCheck) Run(ctx context.Context, projectID string) ([]model.Finding, error) {
	rules, err := c.service.ListFirewallRules(ctx, projectID)
	if err != nil {
		return nil, err
	}

	var findings []model.Finding
	for _, rule := range rules {
		if strings.ToUpper(rule.Direction) != "INGRESS" || !containsCIDR(rule.SourceRanges, "0.0.0.0/0") {
			continue
		}

		for _, allowed := range rule.Allowed {
			if strings.ToLower(allowed.IPProtocol) != "tcp" {
				continue
			}
			if exposesPort(allowed.Ports, 22) {
				findings = append(findings, c.newFinding(projectID, "open_ssh", "OPEN_SSH", rule.Name, "SSH open to the internet", "Restrict port 22 exposure to trusted source ranges"))
			}
			if exposesPort(allowed.Ports, 3389) {
				findings = append(findings, c.newFinding(projectID, "open_rdp", "OPEN_RDP", rule.Name, "RDP open to the internet", "Restrict port 3389 exposure to trusted source ranges"))
			}
		}
	}

	return findings, nil
}

func (c NetworkCheck) newFinding(projectID, ruleID, checkID, resource, message, recommendation string) model.Finding {
	return model.Finding{
		ID:             uuid.NewString(),
		Check:          checkID,
		RuleID:         ruleID,
		Resource:       resource,
		Project:        projectID,
		Severity:       model.SeverityMedium,
		Message:        fmt.Sprintf("%s: %s", message, resource),
		Recommendation: recommendation,
		Timestamp:      c.now(),
	}
}

func containsCIDR(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}

func exposesPort(ports []string, target int) bool {
	for _, port := range ports {
		if strings.Contains(port, "-") {
			bounds := strings.SplitN(port, "-", 2)
			if len(bounds) != 2 {
				continue
			}
			start, err1 := strconv.Atoi(bounds[0])
			end, err2 := strconv.Atoi(bounds[1])
			if err1 == nil && err2 == nil && target >= start && target <= end {
				return true
			}
			continue
		}

		value, err := strconv.Atoi(port)
		if err == nil && value == target {
			return true
		}
	}

	return false
}
