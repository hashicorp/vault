---
layout: "docs"
page_title: "kv - Command"
sidebar_title: "<code>kv</code>"
sidebar_current: "docs-commands-kv"
description: |-
  The "kv" command groups subcommands for interacting with Vault's key-value secret store.
---

# kv

The `kv` command groups subcommands for interacting with Vault's key-value secret store.
For more information, please see the [kv secrets engine
documentation](/docs/secrets/kv/kv-v2.html) page.

## Examples

Create or update the key named "foo" in the "secret" mount with the value
  "bar=baz":

```text
$ vault kv put secret/foo bar=baz

Key              Value
---              -----
created_time     2019-01-12T22:42:43.992144933Z
deletion_time    n/a
destroyed        false
version          1
```

 Read this value back:

```text
$ vault kv get secret/foo

====== Metadata ======
Key              Value
---              -----
created_time     2019-01-12T22:42:43.992144933Z
deletion_time    n/a
destroyed        false
version          1

=== Data ===
Key    Value
---    -----
bar    baz
```

Get metadata for the key:

```text
$ vault kv metadata get secret/foo

======= Metadata =======
Key                Value
---                -----
cas_required       false
created_time       2019-01-12T22:15:17.191030488Z
current_version    2
max_versions       0
oldest_version     0
updated_time       2019-01-12T22:42:43.992144933Z

====== Version 1 ======
Key              Value
---              -----
created_time     2019-01-12T22:15:17.191030488Z
deletion_time    n/a
destroyed        false
```

Get a specific version of the key:

```text
$ vault kv get -version=2 secret/foo

====== Metadata ======
Key              Value
---              -----
created_time     2019-01-12T22:42:43.992144933Z
deletion_time    n/a
destroyed        false
version          2

=== Data ===
Key    Value
---    -----
bar    baz
```

## Usage

```text
Usage: vault kv <subcommand> [options] [args]

Subcommands:
    delete               Deletes versions in the KV store
    destroy              Permanently removes one or more versions in the KV store
    enable-versioning    Turns on versioning for a KV store
    get                  Retrieves data from the KV store
    list                 List data or secrets
    metadata             Interact with Vault's Key-Value storage
    patch                Sets or updates data in the KV store without overwriting.
    put                  Sets or updates data in the KV store
    undelete             Undeletes versions in the KV store
```

For more information, examples, and usage about a subcommand, click on the name
of the subcommand in the sidebar.
