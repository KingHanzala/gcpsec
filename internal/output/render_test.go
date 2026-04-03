package output

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/kinghanzala/gcpsec/internal/model"
)

func TestNewScanReportAndRender(t *testing.T) {
	findings := []model.Finding{
		{
			ID:        "1",
			Message:   "Public bucket: user-data",
			Severity:  model.SeverityHigh,
			Timestamp: time.Now(),
		},
		{
			ID:            "2",
			Message:       "Public bucket: assets-prod",
			Severity:      model.SeverityInfo,
			Allowed:       true,
			AllowedReason: "assets-*",
			Timestamp:     time.Now(),
		},
	}

	report := NewScanReport("demo-project", findings)
	if !report.HasHigh {
		t.Fatalf("expected active HIGH findings")
	}
	if report.Counts[model.SeverityInfo] != 1 {
		t.Fatalf("expected INFO count of 1")
	}

	var buf bytes.Buffer
	if err := RenderScan(&buf, report); err != nil {
		t.Fatalf("RenderScan() error = %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "[HIGH] Public bucket: user-data") {
		t.Fatalf("missing high finding in output")
	}
	if !strings.Contains(output, "(allowed by assets-*)") {
		t.Fatalf("missing allowed finding context")
	}
	if !strings.Contains(output, "HIGH: 1") {
		t.Fatalf("missing summary count")
	}
}
