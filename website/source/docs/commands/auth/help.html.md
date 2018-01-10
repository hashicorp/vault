---
layout: "docs"
page_title: "auth help - Command"
sidebar_current: "docs-commands-auth-help"
description: |-
  The "auth help" command prints usage and help for an auth method.
---

# auth help

The `auth help` command prints usage and help for an auth method.

  - If given a TYPE, this command prints the default help for the auth method of
    that type.

  - If given a PATH, this command prints the help output for the auth method
    enabled at that path. This path must already exist.

Each auth method produces its own help output.

## Examples

Get usage instructions for the userpass auth method:

```text
$ vault auth help userpass
Usage: vault login -method=userpass [CONFIG K=V...]

  The userpass auth method allows users to authenticate using Vault's
  internal user database.

# ...
```

Print usage for the auth method enabled at my-method/

```text
$ vault auth help my-method/
# ...
```

## Usage

There are no flags beyond the [standard set of flags](/docs/commands/index.html)
included on all commands.
