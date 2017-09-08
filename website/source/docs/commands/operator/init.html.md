---
layout: "docs"
page_title: "operator init - Command"
sidebar_current: "docs-commands-operator-init"
description: |-
  The "operator init" command initializes a Vault server. Initialization is the
  process by which Vault's storage backend is prepared to receive data. Since
  Vault server's share the same storage backend in HA mode, you only need to
  initialize one Vault to initialize the storage backend.
---

# operator init

The `operator init` command initializes a Vault server. Initialization is the
process by which Vault's storage backend is prepared to receive data. Since
Vault server's share the same storage backend in HA mode, you only need to
initialize one Vault to initialize the storage backend.

During initialization, Vault generates an in-memory master key and applies
Shamir's secret sharing algorithm to disassemble that master key into a
configuration number of key shares such that a configurable subset of those key
shares must come together to regenerate the master key. These keys are often
called "unseal keys" in Vault's documentation.

This command cannot be run against already-initialized Vault cluster.

For more information on sealing and unsealing, please the [seal concepts page](/docs/concepts/seal.html).

## Examples

Start initialization with the default options:

```text
$ vault operator init
```

Initialize, but encrypt the unseal keys with pgp keys:

```text
$ vault operator init \
    -key-shares=3 \
    -key-threshold=2 \
    -pgp-keys="keybase:hashicorp,keybase:jefferai,keybase:sethvargo"
```

Encrypt the initial root token using a pgp key:

```text
$ vault operator init -root-token-pgp-key="keybase:hashicorp"
```

## Usage

The following flags are available in addition to the [standard set of
flags](/docs/commands/index.html) included on all commands.

### Output Options

- `-format` `(string: "")` - Print the output in the given format. Valid formats
  are "table", "json", or "yaml". The default is table. This can also be
  specified via the `VAULT_FORMAT` environment variable.

### Common Options

- `-key-shares` `(int: 5)` - Number of key shares to split the generated master
  key into. This is the number of "unseal keys" to generate. This is aliased as
  `-n`.

- `-key-threshold` `(int: 3)` - Number of key shares required to reconstruct the
  master key. This must be less than or equal to -key-shares. This is aliased as
  `-t`.

- `-pgp-keys` `(string: "...")` - Comma-separated list of paths to files on disk
  containing public GPG keys OR a comma-separated list of Keybase usernames
  using the format "keybase:<username>". When supplied, the generated unseal
  keys will be encrypted and base64-encoded in the order specified in this list.
  The number of entires must match -key-shares, unless -store-shares are used.

- `-root-token-pgp-key` `(string: "")` - Path to a file on disk containing a
  binary or base64-encoded public GPG key. This can also be specified as a
  Keybase username using the format "keybase:<username>". When supplied, the
  generated root token will be encrypted and base64-encoded with the given
  public key.

- `-status` `(bool": false)` - Print the current initialization status. An exit
  code of 0 means the Vault is already initialized. An exit code of 1 means an
  error occurred. An exit code of 2 means the mean is not initialized.

### Consul Options

- `-consul-auto` `(bool: false)` - Perform automatic service discovery using
  Consul in HA mode. When all nodes in a Vault HA cluster are registered with
  Consul, enabling this option will trigger automatic service discovery based on
  the provided -consul-service value. When Consul is Vault's HA backend, this
  functionality is automatically enabled. Ensure the proper Consul environment
  variables are set (CONSUL_HTTP_ADDR, etc). When only one Vault server is
  discovered, it will be initialized automatically. When more than one Vault
  server is discovered, they will each be output for selection. The default is
  false.

- `-consul-service` `(string: "vault")` - Name of the service in Consul under
  which the Vault servers are registered.

### HSM Options

- `-recovery-pgp-keys` `(string: "...")` - Behaves like `-pgp-keys`, but for the
  recovery key shares. This is only used in HSM mode.

- `-recovery-shares` `(int: 5)` - Number of key shares to split the recovery key
  into. This is only used in HSM mode.

- `-recovery-threshold` `(int: 3)` - Number of key shares required to
  reconstruct the recovery key. This is only used in HSM mode.

- `-stored-shares` `(int: 0)` - Number of unseal keys to store on an HSM. This
  must be equal to `-key-shares`.
