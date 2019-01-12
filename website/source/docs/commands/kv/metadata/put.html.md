---
layout: "docs"
page_title: "kv metadata put - Command"
sidebar_title: "<code>put</code>"
sidebar_current: "docs-commands-kv-metadata-put"
description: |-
  The "metadata put" command can be used to create a blank key in the key-value store or to
  update key configuration for a specified key.
---

# kv metadata put

The `kv metadata put` command can be used to create a blank key in the key-value store or to
  update key configuration for a specified key.

## Examples

Create a key in the key-value store with no data:

```text
$   vault kv metadata put secret/foo
```

Set a max versions setting on the key:

```text
$ vault kv metadata put -max-versions=5 secret/foo
```

Require Check-and-Set for this key:

```text
$  vault kv metadata put -require-cas secret/foo
```

## Usage

```text 
$ vault metadata kv put [options] KEY [DATA]


Common Options:

 -cas-required
      If true the key will require the cas parameter to be set on all write
      requests. If false, the backend’s configuration will be used. The
      default is false.

  -max-versions=<int>
      The number of versions to keep. If not set, the backend’s configured
      max version is used.
```

### Output Options
  
- `-format` `(string: "table")` - Print the output in the given format. Valid
  formats are "table", "json", or "yaml". This can also be specified via the
  `VAULT_FORMAT` environment variable.
