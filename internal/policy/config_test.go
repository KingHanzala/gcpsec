package policy

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/kinghanzala/gcpsec/internal/model"
)

func TestLoadConfigMissingUsesDefaults(t *testing.T) {
	cfg, err := LoadConfig(filepath.Join(t.TempDir(), "missing.yaml"))
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}
	if cfg.Rules["public_buckets"].Severity != model.SeverityHigh {
		t.Fatalf("expected default HIGH severity")
	}
}

func TestLoadConfigOverridesSeverityAndExceptions(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	content := []byte("rules:\n  public_buckets:\n    severity: LOW\n    exceptions:\n      - assets-*\n")
	if err := os.WriteFile(path, content, 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	cfg, err := LoadConfig(path)
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}
	if cfg.Rules["public_buckets"].Severity != model.SeverityLow {
		t.Fatalf("expected LOW severity, got %s", cfg.Rules["public_buckets"].Severity)
	}
	if len(cfg.Rules["public_buckets"].Exceptions) != 1 {
		t.Fatalf("expected 1 exception")
	}
}
