---
layout: "docs"
page_title: "kv rollback - Command"
sidebar_title: "<code>rollback</code>"
sidebar_current: "docs-commands-kv-rollback"
description: |-
  The "kv rollback" command restores a given previous version to the current
  version at the given path.
---

# kv rollback

~> **NOTE:** This is a [K/V Version 2](/docs/secrets/kv/kv-v2.html) secrets
engine command, and not available for Version 1.


The `kv rollback` command restores a given previous version to the current
version at the given path. The value is written as a new version; for instance,
if the current version is 5 and the rollback version is 2, the data from version
2 will become version 6. This command makes it easy to restore unintentionally
overwritten data.

## Examples

Restores the version 2 of the data at key "creds":

```text
$ vault kv rollback -version=2 secret/creds
Key              Value
---              -----
created_time     2019-06-06T17:07:19.299831Z
deletion_time    n/a
destroyed        false
version          6
```

## Usage

There are no flags beyond the [standard set of flags](/docs/commands/index.html)
included on all commands.

### Output Options

- `-format` `(string: "table")` - Print the output in the given format. Valid
  formats are "table", "json", or "yaml". This can also be specified via the
  `VAULT_FORMAT` environment variable.

### Command Options

- `-version` `(int: 0)` - Specifies the version number that should be made
current again.
