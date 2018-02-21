---
layout: "docs"
page_title: "operator rekey - Command"
sidebar_current: "docs-commands-operator-rekey"
description: |-
  The "operator rekey" command generates a new set of unseal keys. This can
  optionally change the total number of key shares or the required threshold of
  those key shares to reconstruct the master key. This operation is zero
  downtime, but it requires the Vault is unsealed and a quorum of existing
  unseal keys are provided.
---

# operator rekey

The `operator rekey` command generates a new set of unseal keys. This can
optionally change the total number of key shares or the required threshold of
those key shares to reconstruct the master key. This operation is zero downtime,
but it requires the Vault is unsealed and a quorum of existing unseal keys are
provided.

An unseal key may be provided directly on the command line as an argument to the
command. If key is specified as "-", the command will read from stdin. If a TTY
is available, the command will prompt for text.

Please see the [rotating and rekeying](/guides/configuration/rekeying-and-rotating.html) for
step-by-step instructions.

## Examples

Initialize a rekey:

```text
$ vault operator rekey \
    -init \
    -key-shares=15 \
    -key-threshold=9
```

Rekey and encrypt the resulting unseal keys with PGP:

```text
$ vault operator rekey \
    -init \
    -key-shares=3 \
    -key-threshold=2 \
    -pgp-keys="keybase:hashicorp,keybase:jefferai,keybase:sethvargo"
```

Store encrypted PGP keys in Vault's core:

```text
$ vault operator rekey \
    -init \
    -pgp-keys="..." \
    -backup
```

Retrieve backed-up unseal keys:

```text
$ vault operator rekey -backup-retrieve
```

Delete backed-up unseal keys:

```text
$ vault operator rekey -backup-delete
```

## Usage

The following flags are available in addition to the [standard set of
flags](/docs/commands/index.html) included on all commands.

### Output Options

- `-format` `(string: "table")` - Print the output in the given format. Valid
  formats are "table", "json", or "yaml". This can also be specified via the
  `VAULT_FORMAT` environment variable.

### Command Options

- `-cancel` `(bool: false)` - Reset the rekeying progress. This will discard any submitted unseal keys
      or configuration. The default is false.

- `-init` `(bool: false)` - Initialize the rekeying operation. This can only be
  done if no rekeying operation is in progress. Customize the new number of key
  shares and key threshold using the `-key-shares` and `-key-threshold flags`.

- `-key-shares` `(int: 5)` - Number of key shares to split the generated master
  key into. This is the number of "unseal keys" to generate. This is aliased as
  `-n`

- `-key-threshold` `(int: 3)` - Number of key shares required to reconstruct the
  master key. This must be less than or equal to -key-shares. This is aliased as
  `-t`.

- `-nonce` `(string: "")` - Nonce value provided at initialization. The same
  nonce value must be provided with each unseal key.

- `-pgp-keys` `(string: "...")` - Comma-separated list of paths to files on disk
  containing public GPG keys OR a comma-separated list of Keybase usernames
  using the format "keybase:<username>". When supplied, the generated unseal
  keys will be encrypted and base64-encoded in the order specified in this list.

- `-status` `(bool: false)` - Print the status of the current attempt without
  providing an unseal key. The default is false.

- `-target` `(string: "barrier")` - Target for rekeying. "recovery" only applies
  when HSM support is enabled.

### Backup Options

- `-backup` `(bool: false)` - Store a backup of the current PGP encrypted unseal
  keys in Vault's core. The encrypted values can be recovered in the event of
  failure or discarded after success. See the -backup-delete and
  -backup-retrieve options for more information. This option only applies when
  the existing unseal keys were PGP encrypted.

- `-backup-delete` `(bool: false)` - Delete any stored backup unseal keys.

- `-backup-retrieve` `(bool: false)` - Retrieve the backed-up unseal keys. This
  option is only available if the PGP keys were provided and the backup has not
  been deleted.
