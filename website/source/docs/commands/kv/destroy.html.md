---
layout: "docs"
page_title: "kv destroy - Command"
sidebar_title: "<code>destroy</code>"
sidebar_current: "docs-commands-kv-destroy"
description: |-
  The "destroy" command permanently removes the specified versions' data from the key-value store. If
  no key exists at the path, no action is taken.
---

# kv destroy

The `kv destroy` command permanently removes the specified versions' data from the key-value store. If
  no key exists at the path, no action is taken.

## Examples

To destroy version 3 of key foo:

```text
$  vault kv destroy -versions=3 secret/foo
```

## Usage

```text 
$ vault kv destroy [options] KEY


Common Options:

  -versions=<string>
      Specifies the version numbers to destroy
```

### Output Options

- `-format` `(string: "table")` - Print the output in the given format. Valid
  formats are "table", "json", or "yaml". This can also be specified via the
  `VAULT_FORMAT` environment variable.
  