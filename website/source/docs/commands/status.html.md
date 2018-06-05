---
layout: "docs"
page_title: "status - Command"
sidebar_current: "docs-commands-status"
description: |-
  The "status" command prints the current state of Vault including whether it is
  sealed and if HA mode is enabled. This command prints regardless of whether
  the Vault is sealed.
---

# status

The `status` command prints the current state of Vault including whether it is
sealed and if HA mode is enabled. This command prints regardless of whether the
Vault is sealed.

The exit code reflects the seal status:

- 0 - unsealed
- 1 - error
- 2 - sealed

## Examples

Check the status:

```text
$ vault status
Sealed: false
Key Shares: 5
Key Threshold: 3
Unseal Progress: 0
Unseal Nonce:
Version: x.y.z
Cluster Name: vault-cluster-49ffd45f
Cluster ID: d2dad792-fb99-1c8d-452e-528d073ba205

High-Availability Enabled: false
```

## Usage

The following flags are available in addition to the [standard set of
flags](/docs/commands/index.html) included on all commands.

### Output Options

- `-format` `(string: "table")` - Print the output in the given format. Valid
  formats are "table", "json", or "yaml". This can also be specified via the
  `VAULT_FORMAT` environment variable.