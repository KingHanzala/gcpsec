package model

import "time"

type Severity string

const (
	SeverityHigh   Severity = "HIGH"
	SeverityMedium Severity = "MEDIUM"
	SeverityLow    Severity = "LOW"
	SeverityInfo   Severity = "INFO"
)

var SeverityOrder = []Severity{
	SeverityHigh,
	SeverityMedium,
	SeverityLow,
	SeverityInfo,
}

type Finding struct {
	ID             string
	Check          string
	RuleID         string
	Resource       string
	Project        string
	Severity       Severity
	Message        string
	Recommendation string
	Timestamp      time.Time
	Allowed        bool
	AllowedReason  string
}
