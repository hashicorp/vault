---
layout: "docs"
page_title: "Commands (CLI)"
sidebar_title: "Commands (CLI)"
sidebar_current: "docs-commands"
description: |-
  In addition to a verbose HTTP API, Vault features a command-line interface
  that wraps common functionality and formats output. The Vault CLI is a single
  static binary. It is a thin wrapper around the HTTP API. Every CLI command
  maps directly to the HTTP API internally.
---

# Vault Commands (CLI)

~> **Note:** The Vault CLI interface was changed substantially in 0.9.2+ and may cause
confusion while using older versions of Vault with this documentation. Read our
[upgrade guide](/guides/upgrading/upgrade-to-0.9.2.html#backwards-compatible-cli-changes) for more information.

In addition to a verbose [HTTP API](/api/index.html), Vault features a
command-line interface that wraps common functionality and formats output. The
Vault CLI is a single static binary. It is a thin wrapper around the HTTP API.
Every CLI command maps directly to the HTTP API internally.

Each command is represented as a command or subcommand. Please see the sidebar
for more information about a particular command. This documentation corresponds
to the latest version of Vault. If you are running an older version, commands
may behave differently. Run `vault -h` or `vault <command> -h` to see the help
output which corresponds to your version.

To get help, run:

```text
$ vault -h
```

To get help for a subcommand, run:

```text
$ vault <subcommand> -h
```

## Exit Codes

The Vault CLI aims to be consistent and well-behaved unless documented
otherwise.

  - Local errors such as incorrect flags, failed validations, or wrong numbers
    of arguments return an exit code of 1.

  - Any remote errors such as API failures, bad TLS, or incorrect API parameters
    return an exit status of 2

Some commands override this default where it makes sense. These commands
document this anomaly.

## Autocompletion

The `vault` command features opt-in autocompletion for flags, subcommands, and
arguments (where supported).

Enable autocompletion by running:

```text
$ vault -autocomplete-install
```

~> Be sure to **restart your shell** after installing autocompletion!

When you start typing a Vault command, press the `<tab>` character to show a
list of available completions. Type `-<tab>` to show available flag completions.

If the `VAULT_*` environment variables are set, the autocompletion will
automatically query the Vault server and return helpful argument suggestions.

## Reading and Writing Data

The four most common operations in Vault are `read`, `write`, `delete`, and
`list`. These operations work on most paths in Vault. Some paths will
contain secrets, other paths might contain configuration. Whatever it is, the
primary interface for reading and writing data to Vault is similar.

To demonstrate basic read and write operations, the built-in key/value (K/V)
secrets engine will be used. This engine is automatically mounted and has no
external dependencies, making it practical for this introduction. Note that
K/V uses slightly different commands for reading and writing: `kv get`
and `kv put`, respectively.

~> The original version of K/V used the common `read` and `write` operations.
A more advanced K/V Version 2 engine was released in Vault 0.10 and introduced
the `kv get` and `kv put` commands.


### Writing Data

To write data to Vault, use the `vault kv put` command:

```text
$ vault kv put secret/password value=itsasecret
```

For some secrets engines, the key/value pairs are arbitrary. For others, they
are generally more strict. Vault's built-in help will guide you to these
restrictions where appropriate.

#### stdin

Some commands in Vault can read data from stdin using `-` as the value. If `-`
is the entire argument, Vault expects to read a JSON object from stdin:

```text
$ echo -n '{"value":"itsasecret"}' | vault kv put secret/password -
```

In addition to reading full JSON objects, Vault can read just a value from
stdin:

```text
$ echo -n "itsasecret" | vault kv put secret/password value=-
```

#### Files

Some commands can also read data from a file on disk. The usage is similar to
stdin as documented above. If an argument starts with `@`, Vault will read it as
a file:

```text
$ vault kv put secret/password @data.json
```

Or specify the contents of a file as a value:

```text
$ vault kv put secret/password value=@data.txt
```

### Reading Data

After data is persisted, read it back using `vault kv get`:

```
$ vault kv get secret/password
Key                 Value
---                 -----
refresh_interval    768h0m0s
value               itsasecret
```

## Token Helper

By default, the Vault CLI uses a "token helper" to cache the token after
authentication. This is conceptually similar to how a website securely stores
your session information as a cookie in the browser. Token helpers are
customizable, and you can even build your own.

The default token helper stores the token in `~/.vault-token`. You can delete
this file at any time to "logout" of Vault.

## Environment Variables

The CLI reads the following environment variables to set behavioral defaults.
This can alleviate the need to repetitively type a flag. Flags always take
precedence over the environment variables.

### `VAULT_TOKEN`

Vault authentication token. Conceptually similar to a session token on a
website, the `VAULT_TOKEN` environment variable holds the contents of the token.
For more information, please see the [token
concepts](/docs/concepts/tokens.html) page.

### `VAULT_ADDR`

Address of the Vault server expressed as a URL and port, for example:
`https://127.0.0.1:8200/`.

### `VAULT_CACERT`

Path to a PEM-encoded CA certificate _file_ on the local disk. This file is used
to verify the Vault server's SSL certificate. This environment variable takes
precedence over `VAULT_CAPATH`.

### `VAULT_CAPATH`

Path to a _directory_ of PEM-encoded CA certificate files on the local disk.
These certificates are used to verify the Vault server's SSL certificate.

### `VAULT_CLIENT_CERT`

Path to a PEM-encoded client certificate on the local disk. This file is used
for TLS communication with the Vault server.

### `VAULT_CLIENT_KEY`

Path to an unencrypted, PEM-encoded private key on disk which corresponds to the
matching client certificate.

### `VAULT_CLIENT_TIMEOUT`

Timeout variable. The default value is 60s.

### `VAULT_CLUSTER_ADDR`

Address that should be used for other cluster members to connect to this node
when in High Availability mode.

### `VAULT_MAX_RETRIES`

Maximum number of retries when a `5xx` error code is encountered. The default is
`2`, for three total attempts. Set this to `0` or less to disable retrying.

### `VAULT_REDIRECT_ADDR`

Address that should be used when clients are redirected to this node when in
High Availability mode.

### `VAULT_SKIP_VERIFY`

Do not verify Vault's presented certificate before communicating with it.
Setting this variable is not recommended and voids Vault's [security
model](/docs/internals/security.html).

### `VAULT_TLS_SERVER_NAME`

Name to use as the SNI host when connecting via TLS.

### `VAULT_CLI_NO_COLOR`

If provided, Vault output will not include ANSI color escape sequence characters.

### `VAULT_RATE_LIMIT`

This enviroment variable will limit the rate at which the `vault` command
sends requests to Vault.

This enviroment variable has the format `rate[:burst]` (where items in `[]` are
optional). If not specified, the burst value defaults to rate. Both rate and 
burst are specified in "operations per second". If the environment variable is
not specified, then the rate and burst will be unlimited *i.e.* rate 
limiting is off by default.

*Note:* The rate is limited for each invocation of the `vault` CLI. Since
each invocation of the `vault` CLI typically only makes a few requests,
this enviroment variable is most useful when using the Go 
[Vault client API](https://www.vaultproject.io/api/libraries.html#go).

### `VAULT_NAMESPACE`

The namespace to use for the command. Setting this is not necessary
but allows using relative paths.

### `VAULT_MFA`

**ENTERPRISE ONLY**

MFA credentials in the format `mfa_method_name[:key[=value]]` (items in `[]` are
optional). Note that when using the environment variable, only one credential
can be supplied. If a MFA method expects multiple credential values, or if there
are multiple MFA methods specified on a path, then the CLI flag `-mfa` should be
used.

## Flags

There are different CLI flags that are available depending on subcommands. Some
flags, such as those used for setting HTTP and output options, are available
globally, while others are specific to a particular subcommand. For a completely
list of available flags, run:

```text
$ vault <subcommand> -h
```
