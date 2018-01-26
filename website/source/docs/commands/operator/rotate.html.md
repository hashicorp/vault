---
layout: "docs"
page_title: "operator rotate - Command"
sidebar_current: "docs-commands-operator-rotate"
description: |-
  The "operator rotate" rotates the underlying encryption key which is used to
  secure data written to the storage backend. This installs a new key in the key
  ring. This new key is used to encrypted new data, while older keys in the ring
  are used to decrypt older data.
---

# operator rotate

The `operator rotate` rotates the underlying encryption key which is used to
secure data written to the storage backend. This installs a new key in the key
ring. This new key is used to encrypted new data, while older keys in the ring
are used to decrypt older data.

This is an online operation and does not cause downtime. This command is run
per-cluster (not per-server), since Vault servers in HA mode share the same
storage backend.

## Examples

Rotate Vault's encryption key:

```text
$ vault operator rotate
Key Term        3
Install Time    01 May 17 10:30 UTC
```

## Usage

There are no flags beyond the [standard set of flags](/docs/commands/index.html)
included on all commands.
