package scanner

import (
	"context"
	"testing"

	"github.com/kinghanzala/gcpsec/internal/gcp"
)

type fakeStorageService struct {
	buckets []gcp.Bucket
	err     error
}

func (f fakeStorageService) ListBuckets(context.Context, string) ([]gcp.Bucket, error) {
	return f.buckets, f.err
}

func TestStorageCheckFindings(t *testing.T) {
	check := NewStorageCheck(fakeStorageService{
		buckets: []gcp.Bucket{
			{Name: "private-bucket", PublicIAM: false},
			{Name: "public-bucket", PublicIAM: true},
		},
	})

	findings, err := check.Run(context.Background(), "demo-project")
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].Resource != "public-bucket" {
		t.Fatalf("unexpected resource: %s", findings[0].Resource)
	}
}
