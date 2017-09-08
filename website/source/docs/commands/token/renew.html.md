---
layout: "docs"
page_title: "token renew - Command"
sidebar_current: "docs-commands-token-renew"
description: |-
  The "token renew" renews a token's lease, extending the amount of time it can
  be used. If a TOKEN is not provided, the locally authenticated token is used.
  Lease renewal will fail if the token is not renewable, the token has already
  been revoked, or if the token has already reached its maximum TTL.
---

# token renew

The `token renew` renews a token's lease, extending the amount of time it can be
used. If a TOKEN is not provided, the locally authenticated token is used. Lease
renewal will fail if the token is not renewable, the token has already been
revoked, or if the token has already reached its maximum TTL.

## Examples

Renew a token (this uses the `/auth/token/renew` endpoint and permission):

```text
$ vault token renew 96ddf4bc-d217-f3ba-f9bd-017055595017
```

Renew the currently authenticated token (this uses the `/auth/token/renew-self`
endpoint and permission):

```text
$ vault token renew
```

Renew a token requesting a specific increment value:

```text
$ vault token renew -increment=30m 96ddf4bc-d217-f3ba-f9bd-017055595017
```

## Usage

The following flags are available in addition to the [standard set of
flags](/docs/commands/index.html) included on all commands.

### Output Options

- `-format` `(default: "table")` - Print the output in the given format. Valid
  formats are "table", "json", or "yaml". This can also be specified via the
  `VAULT_FORMAT` environment variable.

### Command Options

- `-increment` `(duration: "")` - Request a specific increment for renewal.
  Vault is not required to honor this request. If not supplied, Vault will use
  the default TTL. This is specified as a numeric string with suffix like "30s"
  or "5m". This is aliased as "-i".
