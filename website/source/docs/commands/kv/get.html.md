---
layout: "docs"
page_title: "kv get - Command"
sidebar_title: "<code>get</code>"
sidebar_current: "docs-commands-kv-get"
description: |-
  The "kv get" command retrieves the value from Vault's key-value store at the
  given key name. If no key exists with that name, an error is returned. If a
  key exists with that name but has no data, nothing is returned.
---

# kv get

The `kv get` command retrieves the value from K/V secrets engine at the given
key name. If no key exists with that name, an error is returned. If a key exists
with the name but has no data, nothing is returned.

## Examples

Retrieve the data of the key "creds":

```text
$ vault kv get secret/creds
====== Metadata ======
Key              Value
---              -----
created_time     2019-06-06T06:03:26.595978Z
deletion_time    n/a
destroyed        false
version          5

====== Data ======
Key         Value
---         -----
passcode    my-long-passcode
```

If K/V Version 1 secrets engine is enabled at "secret", the output has no
metadata since there is no versioning information associated with the data:

```text
$ vault kv get secret/creds
====== Data ======
Key         Value
---         -----
passcode    my-long-passcode
```

## Usage

There are no flags beyond the [standard set of flags](/docs/commands/index.html)
included on all commands.

### Output Options

- `-field` `(string: "")` - Print only the field with the given name. Specifying
  this option will take precedence over other formatting directives. The result
  will not have a trailing newline making it ideal for piping to other
  processes.

- `-format` `(string: "table")` - Print the output in the given format. Valid
  formats are "table", "json", or "yaml". This can also be specified via the
  `VAULT_FORMAT` environment variable.

### Command Options

- `-version` `(int: 0)` - Specifies the version to return. If not set the
 latest version is returned.
