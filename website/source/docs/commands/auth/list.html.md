---
layout: "docs"
page_title: "auth list - Command"
sidebar_current: "docs-commands-auth-list"
description: |-
  The "auth list" command lists the auth methods enabled. The output lists the
  enabled auth methods and options for those methods.
---

# auth list

The `auth list` command lists the auth methods enabled. The output lists the
enabled auth methods and options for those methods.

## Examples

List all auth methods:

```text
$ vault auth list
Path         Type        Description
----         ----        -----------
token/       token       token based credentials
userpass/    userpass    n/a
```

List detailed auth method information:

```text
$ vault auth list -detailed
Path         Type        Accessor                  Plugin    Default TTL    Max TTL    Replication    Description
----         ----        --------                  ------    -----------    -------    -----------    -----------
token/       token       auth_token_b2166f9e       n/a       system         system     replicated     token based credentials
userpass/    userpass    auth_userpass_eea6507e    n/a       system         system     replicated     n/a
```

## Usage

The following flags are available in addition to the [standard set of
flags](/docs/commands/index.html) included on all commands.

### Output Options

- `-format` `(string: "table")` - Print the output in the given format. Valid
  formats are "table", "json", or "yaml". This can also be specified via the
  `VAULT_FORMAT` environment variable.

### Command Options

- `-detailed` `(bool: false)` - Print detailed information such as configuration
  and replication status about each auth method.
