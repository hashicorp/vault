---
layout: "docs"
page_title: "lease renew - Command"
sidebar_current: "docs-commands-lease-renew"
description: |-
  The "lease renew" command renews the lease on a secret, extending the time
  that it can be used before it is revoked by Vault.
---

# lease renew

The `lease renew` command renews the lease on a secret, extending the time that
it can be used before it is revoked by Vault.

Every secret in Vault has a lease associated with it. If the owner of the secret
wants to use it longer than the lease, then it must be renewed. Renewing the
lease does not change the contents of the secret.

## Examples

Renew a lease:

```text
$ vault lease renew database/creds/readonly/27e1b9a1-27b8-83d9-9fe0-d99d786bdc83
Success! Revoked lease: database/creds/readonly/27e1b9a1-27b8-83d9-9fe0-d99d786bdc83
```

## Usage

The following flags are available in addition to the [standard set of
flags](/docs/commands/index.html) included on all commands.

- `-increment` `(duration: "")` - Request a specific increment in seconds. Vault
  is not required to honor this request.
