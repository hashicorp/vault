---
layout: "docs"
page_title: "secrets disable - Command"
sidebar_current: "docs-commands-secrets-disable"
description: |-
  The "secrets disable" command disables an secrets engine at a given PATH. The
  argument corresponds to the enabled PATH of the engine, not the TYPE! All
  secrets created by this engine are revoked and its Vault data is removed.
---

# secrets disable

The `secrets disable` command disables an secrets engine at a given PATH. The
argument corresponds to the enabled PATH of the engine, not the TYPE! All
secrets created by this engine are revoked and its Vault data is removed.

Once an secrets engine is disabled, **all secrets generated via the secrets
engine are immediately revoked.**

## Examples

Disable the secrets engine enabled at aws/:

```text
$ vault secrets disable aws/
```

## Usage

There are no flags beyond the [standard set of flags](/docs/commands/index.html)
included on all commands.
