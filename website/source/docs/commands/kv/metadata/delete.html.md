---
layout: "docs"
page_title: "kv metadata delete - Command"
sidebar_title: "<code>delete</code>"
sidebar_current: "docs-commands-kv-metadata-delete"
description: |-
  The "metadata delete" command deletes all versions and metadata for the provided key.
---

# kv metadata delete

The `metadata delete` command deletes all versions and metadata for the provided key.

## Examples

Delete a key and all existing versions:

```text
$ vault kv metadata delete secret/foo
```

## Usage

There are no flags beyond the [standard set of flags](/docs/commands/index.html)
included on all commands.