---
layout: "docs"
page_title: "server - Command"
sidebar_current: "docs-commands-server"
description: |-
  The "server" command starts a Vault server that responds to API requests. By
  default, Vault will start in a "sealed" state. The Vault cluster must be
  initialized before use.
---

# server

The `server` command starts a Vault server that responds to API requests. By
default, Vault will start in a "sealed" state. The Vault cluster must be
initialized before use, usually by the `vault operator init` command. Each Vault
server must also be unsealed using the `vault operator unseal` command or the
API before the server can respond to requests.

For more information, please see:

- [`operator init` command](/docs/commands/operator/init.html) for information
  on initializing a Vault server.

- [`operator unseal` command](/docs/commands/operator/unseal.html) for
  information on providing unseal keys.

- [Vault configuration](/docs/configuration/index.html) for the syntax and
  various configuration options for a Vault server.

## Examples

Start a server with a configuration file:

```text
$ vault server -config=/etc/vault/config.hcl
```

Run in "dev" mode with a custom initial root token:

```text
$ vault server -dev -dev-root-token-id="root"
```

## Usage

The following flags are available in addition to the [standard set of
flags](/docs/commands/index.html) included on all commands.

### Command Options

- `-config` `(string: "")` - Path to a configuration file or directory of
  configuration files. This flag can be specified multiple times to load
  multiple configurations. If the path is a directory, all files which end in
  .hcl or .json are loaded.

- `-log-level` `(string: "info")` - Log verbosity level. Supported values (in
  order of detail) are "trace", "debug", "info", "warn", and "err". This can
  also be specified via the VAULT_LOG environment variable.

### Dev Options

- `-dev` `(bool: false)` - Enable development mode. In this mode, Vault runs
  in-memory and starts unsealed. As the name implies, do not run "dev" mode in
  production.

- `-dev-listen-address` `(string: "127.0.0.1:8200")` - Address to bind to in
  "dev" mode. This can also be specified via the `VAULT_DEV_LISTEN_ADDRESS`
  environment variable.

- `-dev-root-token-id` `(string: "")` - Initial root token. This only applies
  when running in "dev" mode. This can also be specified via the
  `VAULT_DEV_ROOT_TOKEN_ID` environment variable.
