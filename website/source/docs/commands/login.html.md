---
layout: "docs"
page_title: "login - Command"
sidebar_current: "docs-commands-login"
description: |-
  The "login" command authenticates users or machines to Vault using the
  provided arguments. A successful authentication results in a Vault token -
  conceptually similar to a session token on a website.
---

# login

The `login` command authenticates users or machines to Vault using the provided
arguments. A successful authentication results in a Vault token - conceptually
similar to a session token on a website. By default, this token is cached on the
local machine for future requests.

The `-method` flag allows using other auth methods, such as userpass,
github, or cert. For these, additional "K=V" pairs may be required.  For more
information about the list of configuration parameters available for a given
auth method, use the "vault auth help TYPE". You can also use "vault
auth list" to see the list of enabled auth methods.

If an auth method is enabled at a non-standard path, the `-method`
flag still refers to the canonical type, but the `-path` flag refers to the
enabled path.

If the authentication is requested with response wrapping (via `-wrap-ttl`),
the returned token is automatically unwrapped unless:

  - The `-token-only` flag is used, in which case this command will output
    the wrapping token.

  - The `-no-store` flag is used, in which case this command will output the
    details of the wrapping token.

## Examples

By default, login uses a "token" method:

```text
$ vault login 10862232-fd55-701c-9013-d764b5bc3953
Success! You are now authenticated. The token information below is already
stored in the token helper. You do NOT need to run "vault login" again. Future
requests will use this token automatically.

token: 10862232-fd55-701c-9013-d764b5bc3953
accessor: 121533e1-20e7-0b4e-04d6-a8c18b8566d5
renewable: true
policies: [my-policy]
```

To login with a different method, use `-method`:

```text
$ vault login -method=userpass username=my-username
Password (will be hidden):
Success! You are now authenticated. The token information below is already
stored in the token helper. You do NOT need to run "vault login" again. Future
requests will use this token automatically.

token: a700ded8-28ed-907d-abf4-23514b783d52
accessor: e0857619-3912-9981-4e03-8d6c4b2f6c56
duration: 768h
renewable: true
policies: [default]
```

If a github auth method was enabled at the path "github-ent":

```text
$ vault login -method=github -path=github-prod
Success! You are now authenticated. The token information below is already
stored in the token helper. You do NOT need to run "vault login" again. Future
requests will use this token automatically.

token: 7eab2aba-b476-af57-e0af-dfcab7c541f6
accessor: 2ae9b1cd-6d17-3428-bd44-986e97f6d2f3
renewable: 22bc4d76-aa3b-1c53-4349-b230b459b56b
policies: [root]
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

### Command Options

- `-method` `(string "token")` - Type of authentication to use such as
  "userpass" or "ldap". Note this corresponds to the TYPE, not the enabled path.
  Use -path to specify the path where the authentication is enabled.

- `-no-store` `(bool: false)` - Do not persist the token to the token helper
  (usually the local filesystem) after authentication for use in future
  requests. The token will only be displayed in the command output.

- `-path` `(string: "")` - Remote path in Vault where the auth method
  is enabled. This defaults to the TYPE of method (e.g. userpass -> userpass/).

- `-token-only` `(bool: false)` - Output only the token with no verification.
  This flag is a shortcut for "-field=token -no-store". Setting those
  flags to other values will have no affect.
