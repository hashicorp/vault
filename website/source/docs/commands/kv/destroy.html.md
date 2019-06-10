---
layout: "docs"
page_title: "kv destroy - Command"
sidebar_title: "<code>destroy</code>"
sidebar_current: "docs-commands-kv-destroy"
description: |-
  The "kv destroy" command permanently removes the specified version data for
  the provided key and version numbers from the key-value store.
---

# kv destroy

~> **NOTE:** This is a [K/V Version 2](/docs/secrets/kv/kv-v2.html) secrets
engine command, and not available for Version 1.

The `secrets enable` command permanently removes the specified versions' data
from the key/value secrets engine. If no key exists at the path, no action is
taken.


## Examples

Destroy version 11 of the key "creds":

```text
$ vault kv destroy -versions=11 secret/creds
Success! Data written to: secret/destroy/creds
```

## Usage

There are no flags beyond the [standard set of flags](/docs/commands/index.html)
included on all commands.

### Output Options

- `-format` `(string: "table")` - Print the output in the given format. Valid
  formats are "table", "json", or "yaml". This can also be specified via the
  `VAULT_FORMAT` environment variable.

### Command Options

- `-versions` `([]int: <required>)` - The versions to destroy. Their data will
be permanently deleted.
