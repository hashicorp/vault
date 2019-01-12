---
layout: "docs"
page_title: "kv list - Command"
sidebar_title: "<code>list</code>"
sidebar_current: "docs-commands-kv-list"
description: |-
  The "list" command lists data from Vault's key-value store at the given path.
---

# kv list

The `kv list` command lists data from Vault's key-value store at the given path.

## Examples

To list values under the "my-app" folder of the key-value store:

```text
$  vault kv list secret/my-app
```

## Usage

```text 
$ vault kv list [options] PATH
```
The following flags are available in addition to the [standard set of
flags](/docs/commands/index.html) included on all commands.

### Output Options

- `-format` `(string: "table")` - Print the output in the given format. Valid
  formats are "table", "json", or "yaml". This can also be specified via the
  `VAULT_FORMAT` environment variable.
