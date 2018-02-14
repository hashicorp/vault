---
layout: "docs"
page_title: "secrets move - Command"
sidebar_current: "docs-commands-secrets-move"
description: |-
  The "secrets move" command moves an existing secrets engine to a new path. Any
  leases from the old secrets engine are revoked, but all configuration
  associated with the engine is preserved.
---

# secrets move

The `secrets move` command moves an existing secrets engine to a new path. Any
leases from the old secrets engine are revoked, but all configuration associated
with the engine is preserved.

**Moving an existing secrets engine will revoke any leases from the old
engine.**

## Examples

Move the existing secrets engine at secret/ to kv/:

```text
$ vault secrets move secret/ kv/
```

## Usage

There are no flags beyond the [standard set of flags](/docs/commands/index.html)
included on all commands.
