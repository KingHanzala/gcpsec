package scanner

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/kinghanzala/gcpsec/internal/gcp"
	"github.com/kinghanzala/gcpsec/internal/model"
)

type StorageCheck struct {
	service gcp.StorageService
	now     func() time.Time
}

func NewStorageCheck(service gcp.StorageService) StorageCheck {
	return StorageCheck{
		service: service,
		now:     time.Now().UTC,
	}
}

func (c StorageCheck) Name() string {
	return "storage"
}

func (c StorageCheck) Run(ctx context.Context, projectID string) ([]model.Finding, error) {
	buckets, err := c.service.ListBuckets(ctx, projectID)
	if err != nil {
		return nil, err
	}

	var findings []model.Finding
	for _, bucket := range buckets {
		if !bucket.PublicIAM {
			continue
		}
		findings = append(findings, model.Finding{
			ID:             uuid.NewString(),
			Check:          "PUBLIC_BUCKET",
			RuleID:         "public_buckets",
			Resource:       bucket.Name,
			Project:        projectID,
			Severity:       model.SeverityHigh,
			Message:        fmt.Sprintf("Public bucket: %s", bucket.Name),
			Recommendation: "Restrict bucket IAM access to trusted principals only",
			Timestamp:      c.now(),
		})
	}
	return findings, nil
}
