---
layout: "docs"
page_title: "auth - Command"
sidebar_current: "docs-commands-auth"
description: |-
  The "auth" command groups subcommands for interacting with Vault's auth
  methods. Users can list, enable, disable, and get help for different auth
  methods.
---

# auth

The `auth` command groups subcommands for interacting with Vault's auth methods.
Users can list, enable, disable, and get help for different auth methods.

For more information, please see the [auth method
documentation](/docs/auth/index.html) or the [authentication
concepts](/docs/concepts/auth.html) page.

To authenticate to Vault as a user or machine, use the [`vault
login`](/docs/commands/login.html) command instead. This command is for
interacting with the auth methods themselves, not authenticating to Vault.

## Examples

Enable an auth method:

```text
$ vault auth enable userpass
Success! Enabled userpass auth method at: userpass/
```

List all auth methods:

```text
$ vault auth list
Path         Type        Description
----         ----        -----------
token/       token       token based credentials
userpass/    userpass    n/a
```

Get help about how to authenticate to a particular auth method:

```text
$ vault auth help userpass/
Usage: vault login -method=userpass [CONFIG K=V...]
# ...
```

Disable an auth method:

```text
$ vault auth disable userpass/
Success! Disabled the auth method (if it existed) at: userpass/
```

Tune an auth method:

```text
$ vault auth tune -max-lease-ttl=30m userpass/
Success! Tuned the auth method at: userpass/
```

## Usage

```text
Usage: vault auth <subcommand> [options] [args]

  # ...

Subcommands:
    disable    Disables an auth method
    enable     Enables a new auth method
    help       Prints usage for an auth method
    list       Lists enabled auth methods
    tune       Tunes an auth method configuration
```

For more information, examples, and usage about a subcommand, click on the name
of the subcommand in the sidebar.
