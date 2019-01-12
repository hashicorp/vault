---
layout: "docs"
page_title: "kv undelete - Command"
sidebar_title: "<code>undelete</code>"
sidebar_current: "docs-commands-kv-undelete"
description: |-
  The "undelete" command undeletes the data for the provided version and path in the key-value store.
  This restores the data, allowing it to be returned on get requests.
---

# kv undelete

The `kv undelete` command undeletes the data for the provided version and path in the key-value store.
  This restores the data, allowing it to be returned on get requests.

## Examples

To undelete version 3 of key "foo":

```text
$  vault kv undelete -versions=3 secret/foo
```

## Usage

```text 
$ vault kv undelete [options] KEY


Common Options:

  -versions=<string>
      Specifies the version numbers to undelete.
```
