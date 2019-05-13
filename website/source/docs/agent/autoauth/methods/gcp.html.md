---
layout: "docs"
page_title: "Vault Agent Auto-Auth GCP Method"
sidebar_title: "GCP"
sidebar_current: "docs-agent-autoauth-methods-gcp"
description: |-
  GCP Method for Vault Agent Auto-Auth
---

# Vault Agent Auto-Auth GCP Method 

The `gcp` method performs authentication against the [GCP Auth
method](https://www.vaultproject.io/docs/auth/gcp.html). Both `gce` and `iam`
authentication types are supported.

## Credentials

Vault will use the GCP SDK's normal credential chain behavior. You can set a
static `credentials` value but it is usually not needed. If running on GCE
using Application Default Credentials, you may need to specify the service
account and project since ADC does not provide metadata used to automatically
determine these.

## Configuration

- `type` `(string: required)` - The type of authentication; must be `gce` or `iam`

- `role` `(string: required)` - The role to authenticate against on Vault

- `credentials` `(string: optional)` - When using static credentials, the
  contents of the JSON credentials file

- `service_account` `(string: optional)` - The service account to use, if it
  cannot be automatically determined

- `project` `(string: optional)` - The project to use, if it cannot be
  automatically determined

- `jwt_exp` `(string or int: optional)` - The number of minutes a generated JWT
  should be valid for when using the `iam` method; defaults to 15 minutes
