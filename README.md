# gcpsec

<p align="center">
  <img alt="Go" src="https://img.shields.io/badge/Go-1.23.5+-00ADD8?logo=go&logoColor=white">
  <img alt="GCP" src="https://img.shields.io/badge/GCP-Security%20Scanner-4285F4?logo=googlecloud&logoColor=white">
  <img alt="CLI" src="https://img.shields.io/badge/Interface-CLI-222222?logo=gnubash&logoColor=white">
  <img alt="Status" src="https://img.shields.io/badge/Status-MVP-success">
</p>

`gcpsec` is a developer-first GCP security scanner for catching risky project configuration before deployment.

It currently scans:

- Project IAM bindings
- Cloud Storage buckets
- Compute Engine firewall rules

The CLI is built in Go and supports `scan`, `doctor`, `uninstall-info`, and `version`.

## Overview

- Developer-first GCP security scanning
- Human-readable findings with recommendations
- Policy overrides via optional `config.yaml`
- CI-friendly exit codes for blocking risky changes
- Installable with `go install github.com/kinghanzala/gcpsec@latest`

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

## Install

### Option 1: Linux one-line install

Linux users can install the latest release with:

```bash
curl -fsSL https://github.com/KingHanzala/gcpsec/releases/latest/download/install.sh | sh
```

This downloads the correct Linux archive for `amd64` or `arm64` and installs `gcpsec` into `/usr/local/bin`.

If you want a different location:

```bash
curl -fsSL https://github.com/KingHanzala/gcpsec/releases/latest/download/install.sh | INSTALL_DIR="$HOME/.local/bin" sh
```

### Option 2: Download a prebuilt release binary

No Go installation required.

1. Open the GitHub Releases page for this repository.
2. Download the archive matching your OS and CPU architecture.
3. Extract the archive.
4. Move the `gcpsec` binary somewhere on your `PATH`, such as `/usr/local/bin`.

Example for macOS or Linux after downloading a release asset:

```bash
tar -xzf gcpsec_0.1.0_darwin_arm64.tar.gz
chmod +x gcpsec
sudo mv gcpsec /usr/local/bin/gcpsec
gcpsec version
```

Release assets are built automatically by GitHub Actions when a tag like `v0.1.0` is pushed. They are uploaded to GitHub Releases together with checksum files.

### Option 3: Install with Go

Install from GitHub:

```bash
go install github.com/kinghanzala/gcpsec@latest
```

After installation, run the binary as:

```bash
gcpsec version
```

If `zsh` says `command not found: gcpsec`, your Go bin directory is probably not on `PATH`.

Add it to `~/.zshrc`:

```bash
echo 'export PATH="$(go env GOPATH)/bin:$PATH"' >> ~/.zshrc
source ~/.zshrc
```

You can also run the installed binary directly:

```bash
"$(go env GOPATH)/bin/gcpsec" version
```

Notes:

- Module path and repository path: `github.com/kinghanzala/gcpsec`
- Installed binary name: `gcpsec`
- CLI display name in help/version output: `gcpsec`

For `go install ...@latest` to work for other users, the repository must be public and pushed to GitHub.

## Uninstall

See the uninstall command for your machine:

```bash
gcpsec uninstall-info
```

Example output:

```bash
Installed binary: /Users/you/go/bin/gcpsec
Uninstall with: rm -f "/Users/you/go/bin/gcpsec"
```

If you already know the binary is on your `PATH`, the fastest option is:

```bash
rm -f "$(command -v gcpsec)"
```

## Build

Run locally:

```bash
go run . version
```

Build a binary:

```bash
go build -o bin/gcpsec .
```

Run tests:

```bash
GOCACHE=/tmp/go-build go test ./...
```

## Usage

### Scan a project

If installed from a prebuilt release:

```bash
gcpsec scan --project my-project
```

If installed with `go install`:

```bash
gcpsec scan --project my-project
```

Or, if Go bin is not yet on `PATH`:

```bash
"$(go env GOPATH)/bin/gcpsec" scan --project my-project
```

```bash
go run . scan --project my-project
```

Or with a built binary:

```bash
./bin/gcpsec scan --project my-project
```

Use a policy file:

```bash
./bin/gcpsec scan --project my-project --config ./config.yaml
```

Behavior:

- Exit code `1` when at least one non-ignored `HIGH` finding exists
- Exit code `0` when no active `HIGH` findings remain
- Nonzero error code on execution failures such as auth or API issues

### Doctor

Validate local prerequisites:

```bash
gcpsec doctor
```

Or from the built binary:

```bash
./bin/gcpsec doctor
```

### Version

```bash
gcpsec version
```

Or from the built binary:

```bash
./bin/gcpsec version
```

### Uninstall Info

```bash
gcpsec uninstall-info
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

`config.yaml` is optional. If it is missing, `gcpsec` runs with built-in default severities.

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
  run: ./bin/gcpsec scan --project=my-project
```

If your CI installs with Go instead of using a checked-in binary:

```yaml
- name: Install gcpsec
  run: go install github.com/kinghanzala/gcpsec@latest

- name: GCP Security Scan
  run: gcpsec scan --project=my-project
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

## Releases

For stable `@latest` installs:

- Publish the repository at `github.com/kinghanzala/gcpsec`
- Create semantic version tags such as `v0.1.0`, `v0.2.0`, and `v1.0.0`
- Push tags to GitHub so `go install ...@latest` resolves to the newest release version
- Pushing a version tag also triggers the GitHub Actions release workflow in [.github/workflows/release.yml](/Users/kinghanzala/gcpsec/.github/workflows/release.yml)
- The workflow runs tests, builds release archives for macOS, Linux, and Windows, and uploads them to the GitHub Release page

Without release tags, Go may install a pseudo-version derived from the latest commit instead of a clean semver release.

## Future Rename Option

If you later want the install path to be:

```bash
go install github.com/kinghanzala/gcpsec@latest
```
