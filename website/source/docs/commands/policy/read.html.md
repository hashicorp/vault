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

The following flags are available in addition to the [standard set of
flags](/docs/commands/index.html) included on all commands.

### Output Options

- `-format` `(string: "table")` - Print the output in the given format. Valid
  formats are "table", "json", or "yaml". This can also be specified via the
  `VAULT_FORMAT` environment variable.
