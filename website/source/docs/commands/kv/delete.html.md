---
layout: "docs"
page_title: "kv delete - Command"
sidebar_title: "<code>delete</code>"
sidebar_current: "docs-commands-kv-delete"
description: |-
  The "kv delete" command disables a secrets engine at a given PATH. The
  argument corresponds to the enabled PATH of the engine, not the TYPE! All
  secrets created by this engine are revoked and its Vault data is removed.
---

# kv delete

The `kv delete` command deletes the data for the provided path in
the key/value secrets engine. If using K/V Version 2, its versioned data will
not be fully removed, but marked as deleted and will no longer be returned in
normal get requests.

## Examples

Delete the latest version of the key "creds":

```text
$ vault kv delete secret/creds
Success! Data deleted (if it existed) at: secret/creds
```

**[K/V Version 2]** Delete version 11 of key "creds":

```text
$ vault kv delete -versions=11 secret/creds
Success! Data deleted (if it existed) at: secret/creds
```

## Usage

There are no flags beyond the [standard set of flags](/docs/commands/index.html)
included on all commands.


### Command Options

- `-versions` `([]int: <required>)` - The versions to be deleted. The versioned
data will not be deleted, but it will no longer be returned in normal get
requests.

~> **NOTE:** This command option is only for K/V v2.
