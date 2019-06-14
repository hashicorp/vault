---
layout: "docs"
page_title: "kv patch - Command"
sidebar_title: "<code>patch</code>"
sidebar_current: "docs-commands-kv-patch"
description: |-
  The "kv patch" command writes the data to the given path in the key-value
  store. The data can be of any type.
---

# kv patch

~> **NOTE:** This is a [K/V Version 2](/docs/secrets/kv/kv-v2.html) secrets
engine command, and not available for Version 1.


The `kv patch` command writes the data to the given path in the K/V v2 secrets
engine. The data can be of any type. Unlike the `kv put` command, the `patch`
command combines the change with existing data instead of replacing them.
Therefore, this command makes it easy to make a partial updates to an existing
data.

## Examples

If you wish to add an additional key-value (`ttl=48h`) to the existing data at
the key "creds":

```text
$ vault kv patch secret/creds ttl=48h
Key              Value
---              -----
created_time     2019-06-06T16:46:22.090654Z
deletion_time    n/a
destroyed        false
version          6
```

**NOTE:** The `kv put` command requires both the existing data and
the data you wish to add in order to accomplish the same result.

```text
$ vault kv put secret/creds ttl=48h passcode=my-long-passcode
```

The data can also be consumed from a file on disk by prefixing with the "@"
symbol. For example:

```text
$ vault kv patch secret/creds @data.json
```

Or it can be read from stdin using the "-" symbol:

```text
$ echo "abcd1234" | vault kv patch secret/foo bar=-
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
