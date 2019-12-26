---
layout: "docs"
page_title: "kv put - Command"
sidebar_title: "<code>put</code>"
sidebar_current: "docs-commands-kv-put"
description: |-
  The "kv put" command writes the data to the given path in the K/V secrets
  engine. The data can be of any type.
---

# kv put

The `kv put` command writes the data to the given path in the K/V secrets
engine.

If working with K/V v2, this command creates a new version of a secret at the
specified location. If working with K/V v1, this command stores the given secret
at the specified location.

Regardless of the K/V version, if the value does not yet exist at the specified
path, the calling token must have an ACL policy granting the "create"
capability. If the value already exists, the calling token must have an ACL
policy granting the "update" capability.


## Examples

Writes the data to the key "creds":

```text
$ vault kv put secret/creds passcode=my-long-passcode
```

The data can also be consumed from a file on disk by prefixing with the "@"
symbol. For example:

```text
$ vault kv put secret/foo @data.json
```

Or it can be read from stdin using the "-" symbol:

```text
$ echo "abcd1234" | vault kv put secret/foo bar=-
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

- `-cas` `(int: 0)` - Specifies to use a Check-And-Set operation. If not set the
 write will be allowed. If set to 0 a write will only be allowed if the key
 doesn’t exist. If the index is non-zero the write will only be allowed if the
 key’s current version matches the version specified in the cas parameter. The
 default is -1.
