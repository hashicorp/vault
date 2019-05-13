---
layout: "docs"
page_title: "Upgrading to Vault 0.11.0 - Guides"
sidebar_title: "Upgrade to 0.11.0"
sidebar_current: "docs-upgrading-to-0.11.0"
description: |-
  This page contains the list of deprecations and important or breaking changes
  for Vault 0.11.0. Please read it carefully.
---

# Overview

This page contains the list of deprecations and important or breaking changes
for Vault 0.11.0 compared to 0.10.0. Please read it carefully.

## Known Issues

### Nomad Integration

Users that integrate Vault with Nomad should hold off on upgrading.  A modification to
Vault's API is causing a runtime issue with the Nomad to Vault integration.

### Minified JSON Policies

Users that generate policies in minfied JSON may cause a parsing errors due to
a regression in the policy parser when it encounters repeating brackets. Although
HCL is the official language for policies in Vault, HCL is JSON compatible and JSON
should work in place of HCL. To work around this error, pretty print the JSON policies
or add spaces between repeating brackets.  This regression will be addressed in
a future release.

### Common Mount Prefixes

Before running the upgrade, users should run `vault secrets list` and `vault auth list` 
to check their mount table to ensure that mounts do not have common prefix "folders".  
For example, if there is a mount with path `team1/` and a mount with path `team1/secrets`, 
Vault will fail to unseal. Before upgrade, these mounts must be remounted at a path that 
does not share a common prefix.

## Changes Since 0.10.4

### Request Timeouts

A default request timeout of 90s is now enforced. This setting can be
overwritten in the config file. If you anticipate requests taking longer than
90s this setting should be configured before upgrading.

### `sys/` Top Level Injection

For the last two years for backwards compatibility data for various `sys/`
routes has been injected into both the Secret's Data map and into the top level
of the JSON response object. However, this has some subtle issues that pop up
from time to time and is becoming increasingly complicated to maintain, so it's
finally being removed.

### Path Fallback for List Operations

For a very long time Vault has automatically adjusted `list` operations to
always end in a `/`, as list operations operates on prefixes, so all list
operations by definition end with `/`. This was done server-side so affects all
clients. However, this has also led to a lot of confusion for users writing
policies that assume that the path that they use in the CLI is the path used
internally. Starting in 0.11, ACL policies gain a new fallback rule for
listing: they will use a matching path ending in `/` if available, but if not
found, they will look for the same path without a trailing `/`. This allows
putting `list` capabilities in the same path block as most other capabilities
for that path, while not providing any extra access if `list` wasn't actually
provided there.

### Performance Standbys On By Default

If your flavor/license of Vault Enterprise supports Performance Standbys, they
are on by default. You can disable this behavior per-node with the
`disable_performance_standby` configuration flag.

### AWS Secret Engine Roles
Roles in the AWS Secret Engine were previously ambiguous. For example, if the
`arn` parameter had been specified, that could have been interpreted as the ARN
of an AWS IAM policy to attach to an IAM user or it could have been the ARN of
an AWS role to assume. Now, types are explicit, both in terms of what
credential type is being requested (e.g., an IAM User or an Assumed Role?) as
well as the parameters being sent to vault (e.g., the IAM policy document
attached to an IAM user or used during a GetFederationToken call). All
credential retrieval remains backwards compatible as does updating role data.
However, the data returned when reading role data is now different and
breaking, so anything which reads role data out of Vault will need to be
updated to handle the new role data format.

While creating/updating roles remains backwards compatible, the old parameters
are now considered deprecated. You should use the new parameters as documented
in the API docs.

As part of this, the `/aws/creds/` and `/aws/sts/` endpoints have been merged,
with the behavior only differing as specified below. The `/aws/sts/` endpoint
is considered deprecated and should only be used when needing backwards
compatibility.

All roles will be automatically updated to the new role format when accessed.
However, due to the way role data was previously being stored in Vault, it's
possible that invalid data was stored that both make the upgrade impossible as
well as would have made the role unable to retrieve credentials. In this
situation, the previous role data is returned in an `invalid_data` key so you
can inspect what used to be in the role and correct the role data if desired.
One consequence of the prior AWS role storage format is that a single Vault
role could have led to two different AWS credential types being retrieved when
a `policy` parameter was stored. In this case, these legacy roles will be
allowed to retrieve both IAM User and Federation Token credentials, with the
credential type depending on the path used to access it (IAM User if accessed
via the `/aws/creds/<role_name>` endpoint and Federation Token if accessed via
the `/aws/sts/<role_name>` endpoint).

## Full List Since 0.10.0

### Revocations of dynamic secrets leases now asynchronous

Dynamic secret lease revocation are now queued/asynchronous rather
than synchronous. This allows Vault to take responsibility for revocation
even if the initial attempt fails. The previous synchronous behavior can be
attained via the `-sync` CLI flag or `sync` API parameter. When in
synchronous mode, if the operation results in failure it is up to the user
to retry.

### CLI Retries

The CLI will no longer retry commands on 5xx errors. This was a
source of confusion to users as to why Vault would "hang" before returning a
5xx error. The Go API client still defaults to two retries.

### Identity Entity Alias metadata

You can no longer manually set metadata on
entity aliases. All alias data (except the canonical entity ID it refers to)
is intended to be managed by the plugin providing the alias information, so
allowing it to be set manually didn't make sense.

### Convergent Encryption version 3

If you are using `transit`'s convergent encryption feature, which prior to this
release was at version 2, we recommend
[rotating](https://www.vaultproject.io/api/secret/transit/index.html#rotate-key)
your encryption key (the new key will use version 3) and
[rewrapping](https://www.vaultproject.io/api/secret/transit/index.html#rewrap-data)
your data to mitigate the chance of offline plaintext-confirmation attacks.

### PKI duration return types

The PKI backend now returns durations (e.g. when reading a role) as an integer
number of seconds instead of a Go-style string.
