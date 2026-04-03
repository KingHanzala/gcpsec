package scanner

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/kinghanzala/gcpsec/internal/gcp"
	"github.com/kinghanzala/gcpsec/internal/model"
)

type IAMCheck struct {
	service gcp.IAMService
	now     func() time.Time
}

func NewIAMCheck(service gcp.IAMService) IAMCheck {
	return IAMCheck{
		service: service,
		now:     time.Now().UTC,
	}
}

func (c IAMCheck) Name() string {
	return "iam"
}

func (c IAMCheck) Run(ctx context.Context, projectID string) ([]model.Finding, error) {
	bindings, err := c.service.GetProjectIAMBindings(ctx, projectID)
	if err != nil {
		return nil, err
	}

	var findings []model.Finding
	for _, binding := range bindings {
		for _, member := range binding.Members {
			switch {
			case binding.Role == "roles/owner":
				findings = append(findings, model.Finding{
					ID:             uuid.NewString(),
					Check:          "OWNER_ROLE_ASSIGNED",
					RuleID:         "owner_roles",
					Resource:       member,
					Project:        projectID,
					Severity:       model.SeverityHigh,
					Message:        fmt.Sprintf("Owner role assigned: %s", member),
					Recommendation: "Remove broad owner access and replace it with least-privilege roles",
					Timestamp:      c.now(),
				})
			case isPublicMember(member):
				findings = append(findings, model.Finding{
					ID:             uuid.NewString(),
					Check:          "PUBLIC_IAM_MEMBER",
					RuleID:         "public_iam_members",
					Resource:       fmt.Sprintf("%s -> %s", binding.Role, member),
					Project:        projectID,
					Severity:       model.SeverityHigh,
					Message:        fmt.Sprintf("Public IAM member bound to %s", binding.Role),
					Recommendation: "Remove allUsers and allAuthenticatedUsers members from project IAM bindings",
					Timestamp:      c.now(),
				})
			}
		}
	}

	return findings, nil
}

func isPublicMember(member string) bool {
	member = strings.TrimSpace(member)
	return member == "allUsers" || member == "allAuthenticatedUsers"
}
