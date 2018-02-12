---
layout: "docs"
page_title: "policy read - Command"
sidebar_current: "docs-commands-policy-read"
description: |-
  The "policy read" command prints the contents and metadata of the Vault policy
  named NAME. If the policy does not exist, an error is returned.
---

# policy read

The `policy read` command prints the contents and metadata of the Vault policy
named NAME. If the policy does not exist, an error is returned.

## Examples

Read the policy named "my-policy":

```text
$ vault policy read my-policy
```

## Usage

There are no flags beyond the [standard set of flags](/docs/commands/index.html)
included on all commands.
