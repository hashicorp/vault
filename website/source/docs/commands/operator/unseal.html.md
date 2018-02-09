---
layout: "docs"
page_title: "operator unseal - Command"
sidebar_current: "docs-commands-operator-unseal"
description: |-
  The "operator unseal" allows the user to provide a portion of the master key
  to unseal a Vault server.
---

# operator unseal

The `operator unseal` allows the user to provide a portion of the master key to
unseal a Vault server. Vault starts in a sealed state. It cannot perform
operations until it is unsealed. This command accepts a portion of the master
key (an "unseal key").

The unseal key can be supplied as an argument to the command, but this is
not recommended as the unseal key will be available in your history:

```text
$ vault operator unseal IXyR0OJnSFobekZMMCKCoVEpT7wI6l+USMzE3IcyDyo=
```

Instead, run the command with no arguments and it will prompt for the key:

```text
$ vault operator unseal
Key (will be hidden): IXyR0OJnSFobekZMMCKCoVEpT7wI6l+USMzE3IcyDyo=
```

For more information on sealing and unsealing, please the [seal concepts
page](/docs/concepts/seal.html).


## Examples

Provide an unseal key:

```text
$ vault operator unseal
Key (will be hidden):
Sealed: false
Key Shares: 1
Key Threshold: 1
Unseal Progress: 0
```

## Usage

The following flags are available in addition to the [standard set of
flags](/docs/commands/index.html) included on all commands.

### Output Options

- `-format` `(string: "table")` - Print the output in the given format. Valid
  formats are "table", "json", or "yaml". This can also be specified via the
  `VAULT_FORMAT` environment variable.

### Command Options

- `-reset` `(bool: false)` - Discard any previously entered keys to the unseal
  process.
