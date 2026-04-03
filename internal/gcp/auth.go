package gcp

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/oauth2/google"
)

var adcScopes = []string{
	"https://www.googleapis.com/auth/cloud-platform.read-only",
}

func ValidateADC(ctx context.Context) error {
	_, err := google.FindDefaultCredentials(ctx, adcScopes...)
	if err == nil {
		return nil
	}

	if strings.Contains(err.Error(), "could not find default credentials") {
		return errors.New("application default credentials not found. Run: gcloud auth application-default login")
	}

	return fmt.Errorf("failed to load application default credentials: %w", err)
}
