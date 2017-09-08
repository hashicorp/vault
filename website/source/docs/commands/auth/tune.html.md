---
layout: "docs"
page_title: "auth tune - Command"
sidebar_current: "docs-commands-auth-tune"
description: |-
  The "auth tune" command tunes the configuration options for the auth method at
  the given PATH.
---

# auth tune

The `auth tune` command tunes the configuration options for the auth method at
the given PATH. **The argument corresponds to the PATH where the auth method is
enabled, not the TYPE!**

## Examples

Tune the default lease for the auth method enabled at "github/":

```text
$ vault auth tune -default-lease-ttl=72h github/
Success! Tuned the auth method at: github/
```

## Usage

The following flags are available in addition to the [standard set of
flags](/docs/commands/index.html) included on all commands.

- `-default-lease-ttl` `(duration: "")` - The default lease TTL for this auth
  method. If unspecified, this defaults to the Vault server's globally
  configured default lease TTL, or a previously configured value for the auth
  method.

- `-max-lease-ttl` `(duration: "")` - The maximum lease TTL for this auth
  method. If unspecified, this defaults to the Vault server's globally
  configured maximum lease TTL, or a previously configured value for the auth
  method.
