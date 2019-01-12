---
layout: "docs"
page_title: "kv enable-versioning - Command"
sidebar_title: "<code>enable-versioning</code>"
sidebar_current: "docs-commands-kv-enable-versioning"
description: |-
  The "enable-versioning" command turns on versioning for the backend at the provided path.
---

# kv enable-versioning

The `kv enable-versioning` command turns on versioning for the backend at the provided path.

## Examples

To enable versioning of the 'secret' path:

```text
$  vault kv enable-versioning secret

Success! Tuned the secrets engine at: secret/
```


To delete all versions and metadata, see the "vault kv metadata" subcommand.

## Usage

```text
$ vault kv enable-versioning [options] KEY
```

### Output Options

- `-format` `(string: "table")` - Print the output in the given format. Valid
  formats are "table", "json", or "yaml". This can also be specified via the
  `VAULT_FORMAT` environment variable.

