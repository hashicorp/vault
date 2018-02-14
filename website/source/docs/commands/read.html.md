---
layout: "docs"
page_title: "read - Command"
sidebar_current: "docs-commands-read"
description: |-
  The "read" command reads data from Vault at the given path. This can be used
  to read secrets, generate dynamic credentials, get configuration details, and
  more.
---

# read

The `read` command reads data from Vault at the given path. This can be used to
read secrets, generate dynamic credentials, get configuration details, and more.

For a full list of examples and paths, please see the documentation that
corresponds to the secrets engine in use.

## Examples

Read a secret from the static secrets engine:

```text
$ vault read secret/my-secret
Key                 Value
---                 -----
refresh_interval    768h
foo                 bar
```

## Usage

The following flags are available in addition to the [standard set of
flags](/docs/commands/index.html) included on all commands.

### Output Options

- `-field` `(string: "")` - Print only the field with the given name. Specifying
  this option will take precedence over other formatting directives. The result
  will not have a trailing newline making it ideal for piping to other processes.

- `-format` `(string: "table")` - Print the output in the given format. Valid
  formats are "table", "json", or "yaml". This can also be specified via the
  `VAULT_FORMAT` environment variable.
