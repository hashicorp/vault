---
layout: "docs"
page_title: "kv metadata get - Command"
sidebar_title: "<code>get</code>"
sidebar_current: "docs-commands-kv-metadata-get"
description: |-
  The "metadata get" command retrieves the metadata from Vault's key-value store at the given key name. If no
  key exists with that name, an error is returned.
---

# kv metadata get

The `kv metadata get` command retrieves the metadata from Vault's key-value store at the given key name. If no
  key exists with that name, an error is returned.

## Examples

Get the metadata for a key, this provides information about each existing
  version:

```text
$ vault kv metadata get secret/foo

======= Metadata =======
Key                Value
---                -----
cas_required       false
created_time       2019-01-12T22:15:17.191030488Z
current_version    5
max_versions       0
oldest_version     0
updated_time       2019-01-12T22:50:59.754692971Z

====== Version 1 ======
Key              Value
---              -----
created_time     2019-01-12T22:15:17.191030488Z
deletion_time    n/a
destroyed        false

====== Version 2 ======
Key              Value
---              -----
created_time     2019-01-12T22:42:43.992144933Z
deletion_time    n/a
destroyed        false
```


## Usage

The following flags are available in addition to the [standard set of
flags](/docs/commands/index.html) included on all commands.

### Output Options

- `-format` `(string: "table")` - Print the output in the given format. Valid
  formats are "table", "json", or "yaml". This can also be specified via the
  `VAULT_FORMAT` environment variable.
