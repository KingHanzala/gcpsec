package gcp

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"cloud.google.com/go/storage"
	crmv1 "google.golang.org/api/cloudresourcemanager/v1"
	compute "google.golang.org/api/compute/v1"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

type IAMBinding struct {
	Role    string
	Members []string
}

type Bucket struct {
	Name      string
	PublicIAM bool
}

type FirewallRule struct {
	Name         string
	Direction    string
	SourceRanges []string
	Allowed      []FirewallAllowed
}

type FirewallAllowed struct {
	IPProtocol string
	Ports      []string
}

type IAMService interface {
	GetProjectIAMBindings(ctx context.Context, projectID string) ([]IAMBinding, error)
}

type StorageService interface {
	ListBuckets(ctx context.Context, projectID string) ([]Bucket, error)
}

type ComputeService interface {
	ListFirewallRules(ctx context.Context, projectID string) ([]FirewallRule, error)
}

type Services struct {
	IAM     IAMService
	Storage StorageService
	Compute ComputeService

	storageClient *storage.Client
}

func NewServices(ctx context.Context) (*Services, error) {
	storageClient, err := storage.NewClient(ctx, option.WithScopes(adcScopes...))
	if err != nil {
		return nil, wrapAPIError("storage client initialization failed", err)
	}

	crmService, err := crmv1.NewService(ctx, option.WithScopes(adcScopes...))
	if err != nil {
		storageClient.Close()
		return nil, wrapAPIError("cloud resource manager client initialization failed", err)
	}

	computeService, err := compute.NewService(ctx, option.WithScopes(adcScopes...))
	if err != nil {
		storageClient.Close()
		return nil, wrapAPIError("compute client initialization failed", err)
	}

	return &Services{
		IAM:           &iamAPI{service: crmService},
		Storage:       &storageAPI{client: storageClient},
		Compute:       &computeAPI{service: computeService},
		storageClient: storageClient,
	}, nil
}

func (s *Services) Close() error {
	if s == nil || s.storageClient == nil {
		return nil
	}
	return s.storageClient.Close()
}

type iamAPI struct {
	service *crmv1.Service
}

func (a *iamAPI) GetProjectIAMBindings(ctx context.Context, projectID string) ([]IAMBinding, error) {
	resp, err := a.service.Projects.GetIamPolicy(projectID, &crmv1.GetIamPolicyRequest{}).Context(ctx).Do()
	if err != nil {
		return nil, wrapAPIError("failed to fetch project IAM policy", err)
	}

	bindings := make([]IAMBinding, 0, len(resp.Bindings))
	for _, binding := range resp.Bindings {
		bindings = append(bindings, IAMBinding{
			Role:    binding.Role,
			Members: append([]string(nil), binding.Members...),
		})
	}
	return bindings, nil
}

type storageAPI struct {
	client *storage.Client
}

func (a *storageAPI) ListBuckets(ctx context.Context, projectID string) ([]Bucket, error) {
	it := a.client.Buckets(ctx, projectID)
	var buckets []Bucket
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, wrapAPIError("failed to list storage buckets", err)
		}

		policy, err := a.client.Bucket(attrs.Name).IAM().V3().Policy(ctx)
		if err != nil {
			return nil, wrapAPIError(fmt.Sprintf("failed to fetch IAM policy for bucket %q", attrs.Name), err)
		}

		publicIAM := false
		for _, binding := range policy.Bindings {
			for _, member := range binding.Members {
				if member == "allUsers" || member == "allAuthenticatedUsers" {
					publicIAM = true
					break
				}
			}
			if publicIAM {
				break
			}
		}

		buckets = append(buckets, Bucket{
			Name:      attrs.Name,
			PublicIAM: publicIAM,
		})
	}

	return buckets, nil
}

type computeAPI struct {
	service *compute.Service
}

func (a *computeAPI) ListFirewallRules(ctx context.Context, projectID string) ([]FirewallRule, error) {
	resp, err := a.service.Firewalls.List(projectID).Context(ctx).Do()
	if err != nil {
		return nil, wrapAPIError("failed to list firewall rules", err)
	}

	rules := make([]FirewallRule, 0, len(resp.Items))
	for _, item := range resp.Items {
		rule := FirewallRule{
			Name:         item.Name,
			Direction:    item.Direction,
			SourceRanges: append([]string(nil), item.SourceRanges...),
			Allowed:      make([]FirewallAllowed, 0, len(item.Allowed)),
		}
		for _, allowed := range item.Allowed {
			rule.Allowed = append(rule.Allowed, FirewallAllowed{
				IPProtocol: allowed.IPProtocol,
				Ports:      append([]string(nil), allowed.Ports...),
			})
		}
		rules = append(rules, rule)
	}
	return rules, nil
}

func wrapAPIError(prefix string, err error) error {
	var gerr *googleapi.Error
	if errors.As(err, &gerr) {
		switch gerr.Code {
		case 403:
			return fmt.Errorf("%s: permission denied. Ensure the account has Viewer access and required APIs are enabled", prefix)
		case 404:
			return fmt.Errorf("%s: API or resource not found. Ensure the required GCP APIs are enabled", prefix)
		}
	}

	msg := err.Error()
	if strings.Contains(msg, "SERVICE_DISABLED") {
		return fmt.Errorf("%s: API disabled. Enable cloudresourcemanager.googleapis.com, compute.googleapis.com, iam.googleapis.com, and storage.googleapis.com", prefix)
	}

	return fmt.Errorf("%s: %w", prefix, err)
}
