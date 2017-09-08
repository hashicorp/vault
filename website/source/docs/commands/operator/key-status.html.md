---
layout: "docs"
page_title: "operator key-status - Command"
sidebar_current: "docs-commands-operator-key-status"
description: |-
  The "operator key-status" provides information about the active encryption
  key.
---

# operator key-status

The `operator key-status` provides information about the active encryption key.
Specifically, the current key term and the key installation time.

## Examples

Get the key status:

```text
$ vault operator key-status
Key Term        2
Install Time    01 Jan 17 12:30 UTC
```

## Usage

There are no flags beyond the [standard set of flags](/docs/commands/index.html)
included on all commands.
