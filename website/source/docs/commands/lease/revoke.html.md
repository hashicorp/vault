---
layout: "docs"
page_title: "lease revoke - Command"
sidebar_current: "docs-commands-lease-revoke"
description: |-
  The "lease revoke" command revokes the lease on a secret, invalidating the
  underlying secret.
---

# lease revoke

The `lease revoke` command revokes the lease on a secret, invalidating the
underlying secret.

## Examples

Revoke a lease:

```text
$ vault lease revoke database/creds/readonly/27e1b9a1-27b8-83d9-9fe0-d99d786bdc83
Success! Revoked lease: database/creds/readonly/27e1b9a1-27b8-83d9-9fe0-d99d786bdc83
```

Revoke a lease which starts with a prefix:

```text
$ vault lease revoke -prefix database/creds
Success! Revoked any leases with prefix: database/creds
```

## Usage

The following flags are available in addition to the [standard set of
flags](/docs/commands/index.html) included on all commands.

- `-force` `(bool: false)` - Delete the lease from Vault even if the secret
  engine revocation fails. This is meant for recovery situations where the
  secret in the target secrets engine was manually removed. If this flag is
  specified, -prefix is also required. This is aliased as "-f". The default is
  false.

- `-prefix` `(bool: false)` - Treat the ID as a prefix instead of an exact lease
  ID. This can revoke multiple leases simultaneously. The default is false.
