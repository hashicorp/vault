---
layout: "docs"
page_title: "secrets tune - Command"
sidebar_current: "docs-commands-secrets-tune"
description: |-
  The "secrets tune" command tunes the configuration options for the secrets
  engine at the given PATH. The argument corresponds to the PATH where the
  secrets engine is enabled, not the TYPE!
---

# secrets tune

The `secrets tune` command tunes the configuration options for the secrets
engine at the given PATH. The argument corresponds to the PATH where the secrets
engine is enabled, not the TYPE!

## Examples

Tune the default lease for the PKI secrets engine:

```text
$ vault secrets tune -default-lease-ttl=72h pki/
```

## Usage

The following flags are available in addition to the [standard set of
flags](/docs/commands/index.html) included on all commands.

- `-default-lease-ttl` `(duration: "")` - The default lease TTL for this secrets
  engine. If unspecified, this defaults to the Vault server's globally
  configured default lease TTL, or a previously configured value for the secrets
  engine.

- `-max-lease-ttl` `(duration: "")` - The maximum lease TTL for this secrets
  engine. If unspecified, this defaults to the Vault server's globally
  configured maximum lease TTL, or a previously configured value for the secrets
  engine.
