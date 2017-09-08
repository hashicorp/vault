---
layout: "docs"
page_title: "lease - Command"
sidebar_current: "docs-commands-lease"
description: |-
  The "lease" command groups subcommands for interacting with leases attached to
  secrets.
---

# lease

The `lease` command groups subcommands for interacting with leases attached to
secrets. For leases attached to tokens, use the [`vault
token`](/docs/commands/token.html) subcommand.

## Examples

Renew a lease:

```text
$ vault lease renew database/creds/readonly/27e1b9a1-27b8-83d9-9fe0-d99d786bdc83
Key                Value
---                -----
lease_id           database/creds/readonly/27e1b9a1-27b8-83d9-9fe0-d99d786bdc83
lease_duration     5m
lease_renewable    true
```

Revoke a lease:

```text
$ vault lease revoke database/creds/readonly/27e1b9a1-27b8-83d9-9fe0-d99d786bdc83
Success! Revoked lease: database/creds/readonly/27e1b9a1-27b8-83d9-9fe0-d99d786bdc83
```

## Usage

```text
Usage: vault lease <subcommand> [options] [args]

  # ...

Subcommands:
    renew     Renews the lease of a secret
    revoke    Revokes leases and secrets
```

For more information, examples, and usage about a subcommand, click on the name
of the subcommand in the sidebar.
