---
layout: "docs"
page_title: "policy list - Command"
sidebar_current: "docs-commands-policy-list"
description: |-
  The "policy list" command Lists the names of the policies that are installed
  on the Vault server.
---

# policy list

The `policy list` command Lists the names of the policies that are installed on
the Vault server.

## Examples

List the available policies:

```text
$ vault policy list
default
root
```

## Usage

There are no flags beyond the [standard set of flags](/docs/commands/index.html)
included on all commands.
