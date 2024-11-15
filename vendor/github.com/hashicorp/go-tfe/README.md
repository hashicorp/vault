HCP Terraform and Terraform Enterprise Go Client
==============================

[![Tests](https://github.com/hashicorp/go-tfe/actions/workflows/ci.yml/badge.svg)](https://github.com/hashicorp/go-tfe/actions/workflows/ci.yml)
[![GitHub license](https://img.shields.io/github/license/hashicorp/go-tfe.svg)](https://github.com/hashicorp/go-tfe/blob/main/LICENSE)
[![GoDoc](https://godoc.org/github.com/hashicorp/go-tfe?status.svg)](https://godoc.org/github.com/hashicorp/go-tfe)
[![Go Report Card](https://goreportcard.com/badge/github.com/hashicorp/go-tfe)](https://goreportcard.com/report/github.com/hashicorp/go-tfe)
[![GitHub issues](https://img.shields.io/github/issues/hashicorp/go-tfe.svg)](https://github.com/hashicorp/go-tfe/issues)

The official Go API client for [HCP Terraform and Terraform Enterprise](https://www.hashicorp.com/products/terraform).

This client supports the [HCP Terraform V2 API](https://developer.hashicorp.com/terraform/cloud-docs/api-docs).
As Terraform Enterprise is a self-hosted distribution of HCP Terraform, this
client supports both HCP Terraform and Terraform Enterprise use cases. In all package
documentation and API, the platform will always be stated as 'Terraform
Enterprise' - but a feature will be explicitly noted as only supported in one or
the other, if applicable (rare).

## Version Information

Almost always, minor version changes will indicate backwards-compatible features and enhancements. Occasionally, function signature changes that reflect a bug fix may appear as a minor version change. Patch version changes will be used for bug fixes, performance improvements, and otherwise unimpactful changes.

## Example Usage

Construct a new TFE client, then use the various endpoints on the client to
access different parts of the Terraform Enterprise API. The following example lists
all organizations.

### (Recommended Approach) Using custom config to provide configuration details to the API client

```go
import (
  "context"
  "log"

  "github.com/hashicorp/go-tfe"
)

config := &tfe.Config{
	Address: "https://tfe.local",
	Token: "insert-your-token-here",
  RetryServerErrors: true,
}

client, err := tfe.NewClient(config)
if err != nil {
	log.Fatal(err)
}

orgs, err := client.Organizations.List(context.Background(), nil)
if err != nil {
	log.Fatal(err)
}
```

### Using the default config with env vars
The default configuration makes use of the `TFE_ADDRESS` and `TFE_TOKEN` environment variables.

1. `TFE_ADDRESS` - URL of a HCP Terraform or Terraform Enterprise instance. Example: `https://tfe.local`
1. `TFE_TOKEN` - An [API token](https://developer.hashicorp.com/terraform/cloud-docs/users-teams-organizations/api-tokens) for the HCP Terraform or Terraform Enterprise instance.

**Note:** Alternatively, you can set `TFE_HOSTNAME` which serves as a fallback for `TFE_ADDRESS`. It will only be used if `TFE_ADDRESS` is not set and will resolve the host to an `https` scheme. Example: `tfe.local` => resolves to `https://tfe.local`

The environment variables are used as a fallback to configure TFE client if the Address or Token values are not provided as in the cases below:

#### Using the default configuration
```go
import (
  "context"
  "log"

  "github.com/hashicorp/go-tfe"
)

// Passing nil to tfe.NewClient method will also use the default configuration
client, err := tfe.NewClient(tfe.DefaultConfig())
if err != nil {
	log.Fatal(err)
}

orgs, err := client.Organizations.List(context.Background(), nil)
if err != nil {
	log.Fatal(err)
}
```

#### When Address or Token has no value
```go
import (
  "context"
  "log"

  "github.com/hashicorp/go-tfe"
)

config := &tfe.Config{
	Address: "",
	Token: "",
}

client, err := tfe.NewClient(config)
if err != nil {
	log.Fatal(err)
}

orgs, err := client.Organizations.List(context.Background(), nil)
if err != nil {
	log.Fatal(err)
}
```

## Documentation

For complete usage of the API client, see the [full package docs](https://pkg.go.dev/github.com/hashicorp/go-tfe).

## API Coverage

This API client covers most of the existing HCP Terraform API calls and is updated regularly to add new or missing endpoints.

- [x] Account
- [x] Agents
- [x] Agent Pools
- [x] Agent Tokens
- [x] Applies
- [x] Audit Trails
- [x] Changelog
- [x] Comments
- [x] Configuration Versions
- [x] Cost Estimation
- [ ] Feature Sets
- [ ] Invoices
- [x] IP Ranges
- [x] Notification Configurations
- [x] OAuth Clients
- [x] OAuth Tokens
- [x] Organizations
- [x] Organization Memberships
- [x] Organization Tags
- [x] Organization Tokens
- [x] Plan Exports
- [x] Plans
- [x] Policies
- [x] Policy Checks
- [x] Policy Sets
- [x] Policy Set Parameters
- [x] Private Registry
	- [x] Modules
	  - [x] No-Code Modules
	- [x] Providers
	- [x] Provider Versions and Platforms
	- [x] GPG Keys
- [x] Projects
- [x] Runs
- [x] Run Events
- [x] Run Tasks
- [x] Run Tasks Integration
- [x] Run Triggers
- [x] SSH Keys
- [x] Stability Policy
- [x] State Versions
- [x] State Version Outputs
- [ ] Subscriptions
- [x] Team Access
- [x] Team Membership
- [x] Team Tokens
- [x] Teams
- [x] Test Runs
- [x] User Tokens
- [x] Users
- [x] Variable Sets
- [x] Variables
- [ ] VCS Events
- [x] Workspaces
- [x] Workspace-Specific Variables
- [x] Workspace Resources
- [x] Admin
  - [x] Module Sharing
  - [x] Organizations
  - [x] Runs
  - [x] Settings
  - [x] Terraform Versions
  - [x] OPA Versions
  - [x] Sentinel Versions
  - [x] Users
  - [x] Workspaces


## Examples

See the [examples directory](https://github.com/hashicorp/go-tfe/tree/main/examples).

## Running tests

See [TESTS.md](docs/TESTS.md).

## Issues and Contributing

See [CONTRIBUTING.md](docs/CONTRIBUTING.md)

## Releases

See [RELEASES.md](docs/RELEASES.md)
