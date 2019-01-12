---
layout: "docs"
page_title: "kv metadata- Command"
sidebar_title: "<code>metadata</code>"
sidebar_current: "docs-commands-kv-metadata"
description: |-
  The "kv metdadata" command has subcommands for interacting with the metadata endpoint in
  Vault's key-value store.
---

# kv metadata

The `kv metadata` command groups subcommands interacting with the metadata endpoint in
  Vault's key-value store.

For more information, please see the [kv secrets engine
documentation](/docs/secrets/kv/kv-v2.html) page.

## Examples

Create or update a metadata entry for a key:

```text
$ vault kv metadata put -max-versions=5 secret/foo
```

Get the metadata for a key, this provides information about each existing
  version:

```text
$ vault kv metadata get secret/foo
```

Delete a key and all existing versions:

```text
$ vault kv metadata delete secret/foo
```

## Usage

```text
Usage: vault kv metadata <subcommand> [options] [args]

  # ...

Subcommands:
    delete    Deletes all versions and metadata for a key in the KV store
    get       Retrieves key metadata from the KV store
    put       Sets or updates key settings in the KV store
```

For more information, examples, and usage about a subcommand, click on the name
of the subcommand in the sidebar.
