---
layout: "docs"
page_title: "kv delete - Command"
sidebar_title: "<code>delete</code>"
sidebar_current: "docs-commands-kv-delete"
description: |-
  The "delete" command deletes the data for the provided version and path in the key-value store. The
  versioned data will not be fully removed, but marked as deleted and will no
  longer be returned in normal get requests.
---

# kv delete

The "delete" command deletes the data for the provided version and path in the key-value store. The
  versioned data will not be fully removed, but marked as deleted and will no
  longer be returned in normal get requests.

## Examples

To delete the latest version of the key "foo":

```text
$  vault kv delete secret/foo
```

To delete version 3 of key foo:

```text
$ vault kv delete -versions=3 secret/foo
```

To delete all versions and metadata, see the "vault kv metadata" subcommand.

## Usage

```text 
$ vault kv delete [options] PATH


Common Options:

  -versions=<string>
      Specifies the version numbers to delete.
```
