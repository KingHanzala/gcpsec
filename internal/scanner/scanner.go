package scanner

import (
	"context"

	"github.com/kinghanzala/gcpsec/internal/gcp"
	"github.com/kinghanzala/gcpsec/internal/model"
)

type Check interface {
	Name() string
	Run(ctx context.Context, projectID string) ([]model.Finding, error)
}

type CheckResult struct {
	Name     string
	Findings []model.Finding
}

func RunAll(ctx context.Context, projectID string, services *gcp.Services) ([]CheckResult, error) {
	checks := []Check{
		NewIAMCheck(services.IAM),
		NewStorageCheck(services.Storage),
		NewNetworkCheck(services.Compute),
	}

	results := make([]CheckResult, 0, len(checks))
	for _, check := range checks {
		findings, err := check.Run(ctx, projectID)
		if err != nil {
			return nil, err
		}
		results = append(results, CheckResult{
			Name:     check.Name(),
			Findings: findings,
		})
	}

	return results, nil
}
