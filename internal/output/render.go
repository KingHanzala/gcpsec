package output

import (
	"fmt"
	"io"

	"github.com/kinghanzala/gcpsec/internal/model"
)

type ScanReport struct {
	Project  string
	Findings []model.Finding
	Counts   map[model.Severity]int
	HasHigh  bool
}

func NewScanReport(project string, findings []model.Finding) ScanReport {
	counts := map[model.Severity]int{
		model.SeverityHigh:   0,
		model.SeverityMedium: 0,
		model.SeverityLow:    0,
		model.SeverityInfo:   0,
	}

	hasHigh := false
	for _, finding := range findings {
		counts[finding.Severity]++
		if finding.Severity == model.SeverityHigh && !finding.Allowed {
			hasHigh = true
		}
	}

	return ScanReport{
		Project:  project,
		Findings: findings,
		Counts:   counts,
		HasHigh:  hasHigh,
	}
}

func RenderScan(w io.Writer, report ScanReport) error {
	if _, err := fmt.Fprintf(w, "Project: %s\n\n", report.Project); err != nil {
		return err
	}

	if len(report.Findings) == 0 {
		if _, err := fmt.Fprintln(w, "No findings detected."); err != nil {
			return err
		}
	} else {
		for _, finding := range report.Findings {
			if finding.Allowed {
				if _, err := fmt.Fprintf(w, "[INFO] %s (allowed by %s)\n", finding.Message, finding.AllowedReason); err != nil {
					return err
				}
				continue
			}
			if _, err := fmt.Fprintf(w, "[%s] %s\n", finding.Severity, finding.Message); err != nil {
				return err
			}
		}
	}

	if _, err := fmt.Fprintln(w, "\nSummary:"); err != nil {
		return err
	}
	for _, severity := range model.SeverityOrder {
		if _, err := fmt.Fprintf(w, "%s: %d\n", severity, report.Counts[severity]); err != nil {
			return err
		}
	}
	return nil
}
