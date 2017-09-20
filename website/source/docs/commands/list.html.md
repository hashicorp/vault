---
layout: "docs"
page_title: "list - Command"
sidebar_current: "docs-commands-list"
description: |-
  The "list" command lists data from Vault at the given path. This can be used
  to list keys in a, given secrets engine.
---

# list

The `list` command lists data from Vault at the given path. This can be used to
list keys in a, given secrets engine.

## Examples

List values under the "my-app" folder of the KV secrets engine:

```text
$ vault list secret/my-app/
```

## Usage

There are no flags beyond the [standard set of flags](/docs/commands/index.html)
included on all commands.
