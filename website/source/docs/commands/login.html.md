---
layout: "docs"
page_title: "login - Command"
sidebar_title: "<code>login</code>"
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
$ vault login s.3jnbMAKl1i4YS3QoKdbHzGXq
Success! You are now authenticated. The token information displayed below
is already stored in the token helper. You do NOT need to run "vault login"
again. Future Vault requests will automatically use this token.

Key                  Value
---                  -----
token                s.3jnbMAKl1i4YS3QoKdbHzGXq
token_accessor       7Uod1Rm0ejUAz77Oh7SxpAM0
token_duration       767h59m49s
token_renewable      true
token_policies       ["admin" "default"]
identity_policies    []
policies             ["admin" "default"]
```

To login with a different method, use `-method`:

```text
$ vault login -method=userpass username=my-username
Password (will be hidden):
Success! You are now authenticated. The token information below is already
stored in the token helper. You do NOT need to run "vault login" again. Future
requests will use this token automatically.

Key                    Value
---                    -----
token                  s.2y4SU3Sk46dK3p2Y8q2jSBwL
token_accessor         8J125x9SZyB76MI9uF2jSJZf
token_duration         768h
token_renewable        true
token_policies         ["default"]
identity_policies      []
policies               ["default"]
token_meta_username    my-username
```

~> Notice that the command option (`-method=userpass`) precedes the command
argument (`username=my-username`).

If a github auth method was enabled at the path "github-ent":

```text
$ vault login -method=github -path=github-prod
Success! You are now authenticated. The token information below is already
stored in the token helper. You do NOT need to run "vault login" again. Future
requests will use this token automatically.

Key                    Value
---                    -----
token                  s.2f3c5L1MHtnqbuNCbx90utmC
token_accessor         JLUIXJ6ltUftTt2UYRl2lTAC
token_duration         768h
token_renewable        true
token_policies         ["default"]
identity_policies      []
policies               ["default"]
token_meta_org         hashicorp
token_meta_username    my-username
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

- `-no-print` `(bool: false)` - Do not display the token. The token will be
  still be stored to the configured token helper. The default is false.

- `-no-store` `(bool: false)` - Do not persist the token to the token helper
  (usually the local filesystem) after authentication for use in future
  requests. The token will only be displayed in the command output.

- `-path` `(string: "")` - Remote path in Vault where the auth method
  is enabled. This defaults to the TYPE of method (e.g. userpass -> userpass/).

- `-token-only` `(bool: false)` - Output only the token with no verification.
  This flag is a shortcut for "-field=token -no-store". Setting those
  flags to other values will have no affect.
