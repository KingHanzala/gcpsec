# gcp-sec

<p align="center">
  <img alt="Go" src="https://img.shields.io/badge/Go-1.23.5+-00ADD8?logo=go&logoColor=white">
  <img alt="GCP" src="https://img.shields.io/badge/GCP-Security%20Scanner-4285F4?logo=googlecloud&logoColor=white">
  <img alt="CLI" src="https://img.shields.io/badge/Interface-CLI-222222?logo=gnubash&logoColor=white">
  <img alt="Status" src="https://img.shields.io/badge/Status-MVP-success">
</p>

`gcp-sec` is a developer-first GCP security scanner for catching risky project configuration before deployment.

It currently scans:

- Project IAM bindings
- Cloud Storage buckets
- Compute Engine firewall rules

The CLI is built in Go and supports `scan`, `doctor`, and `version`.

## Overview

- Developer-first GCP security scanning
- Human-readable findings with recommendations
- Policy overrides via optional `config.yaml`
- CI-friendly exit codes for blocking risky changes

## What It Checks

### IAM

- Identity and Access Management checks

- `roles/owner` assignments
- Public members such as `allUsers`
- Public members such as `allAuthenticatedUsers`

### Cloud Storage

- Bucket exposure checks

- Publicly accessible buckets

### Network

- Internet exposure checks

- Firewall rules exposing SSH on port `22` to `0.0.0.0/0`
- Firewall rules exposing RDP on port `3389` to `0.0.0.0/0`

## Prerequisites

- Go `1.23.5+`
- A GCP project to scan
- Application Default Credentials configured locally

Authenticate with ADC:

```bash
gcloud auth application-default login
```

Enable the required APIs:

```bash
gcloud services enable \
  cloudresourcemanager.googleapis.com \
  iam.googleapis.com \
  compute.googleapis.com \
  storage.googleapis.com
```

Recommended GCP access:

- `Viewer`
- `Security Reviewer` (optional, depending on your environment)

## Build

Run locally:

```bash
go run . version
```

Build a binary:

```bash
go build -o bin/gcp-sec .
```

Run tests:

```bash
GOCACHE=/tmp/go-build go test ./...
```

## Usage

### Scan a project

```bash
go run . scan --project my-project
```

Or with a built binary:

```bash
./bin/gcp-sec scan --project my-project
```

Use a policy file:

```bash
./bin/gcp-sec scan --project my-project --config ./config.yaml
```

Behavior:

- Exit code `1` when at least one non-ignored `HIGH` finding exists
- Exit code `0` when no active `HIGH` findings remain
- Nonzero error code on execution failures such as auth or API issues

### Doctor

Validate local prerequisites:

```bash
./bin/gcp-sec doctor
```

### Version

```bash
./bin/gcp-sec version
```

## Example Output

```text
Project: my-project

[HIGH] Public bucket: user-data
[HIGH] Owner role assigned: user@example.com
[MEDIUM] SSH open to the internet: allow-ssh
[INFO] Public bucket: assets-prod (allowed by assets-*)

Summary:
HIGH: 2
MEDIUM: 1
LOW: 0
INFO: 1
```

## Policy Configuration

`config.yaml` is optional. If it is missing, `gcp-sec` runs with built-in default severities.

Example:

```yaml
rules:
  public_buckets:
    severity: HIGH
    exceptions:
      - assets-*

  owner_roles:
    severity: HIGH

  public_iam_members:
    severity: HIGH

  open_ssh:
    severity: MEDIUM

  open_rdp:
    severity: MEDIUM
```

Supported behavior:

- Override severity per rule
- Ignore matching resources with wildcard exceptions
- Render ignored findings as `INFO`

## CI Example

```yaml
- name: GCP Security Scan
  run: ./bin/gcp-sec scan --project=my-project
```

## Project Layout

```text
cmd/              Cobra CLI commands
internal/gcp/     ADC validation and GCP clients
internal/model/   Shared finding types
internal/output/  Terminal rendering and summary logic
internal/policy/  YAML policy loading and exception handling
internal/scanner/ Scanner runner and resource checks
main.go           CLI entrypoint
```

## Current Scope

This project is intentionally narrow. It does not try to replace the full GCP security ecosystem.

Current MVP focus:

- Fast local scans
- Clear findings
- Simple policy overrides
- CI-friendly exit codes

Planned future areas include service account checks, multi-project scans, and machine-readable output formats.
