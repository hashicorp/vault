---
layout: "docs"
page_title: "kv get - Command"
sidebar_title: "<code>get</code>"
sidebar_current: "docs-commands-kv-get"
description: |-
  The "get" command retrieves the value from Vault's key-value store at the given key name. If no
  key exists with that name, an error is returned. If a key exists with that
  name but has no data, nothing is returned.
---

# kv get

The `kv get` command retrieves the value from Vault's key-value store at the given key name. If no
  key exists with that name, an error is returned. If a key exists with that
  name but has no data, nothing is returned.

## Examples

To get the value of the key named "foo" in the "secret" mount:

```text
$  vault kv get secret/foo

====== Metadata ======
Key              Value
---              -----
created_time     2019-01-12T22:15:17.191030488Z
deletion_time    n/a
destroyed        false
version          1

=== Data ===
Key    Value
---    -----
bar    baz

```

To view the given key name at a specific version in time:

```text
$ vault kv get -version=1 secret/foo
```

## Usage

```text 
$ vault kv get [options] KEY


Common Options:

   -version=<int>
      If passed, the value at the version number will be returned.
```
### Output Options

- `-field` `(string: "field_name")`
  Print only the field with the given name. Specifying this option will
  take precedence over other formatting directives. The result will not
  have a trailing newline making it idea for piping to other processes.

- `-format` `(string: "table")` - Print the output in the given format. Valid
  formats are "table", "json", or "yaml". This can also be specified via the
  `VAULT_FORMAT` environment variable.
