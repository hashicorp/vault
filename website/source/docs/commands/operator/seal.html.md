---
layout: "docs"
page_title: "operator seal - Command"
sidebar_current: "docs-commands-operator-seal"
description: |-
  The "operator seal" seals the Vault server. Sealing tells the Vault server to
  stop responding to any operations until it is unsealed. When sealed, the Vault
  server discards its in-memory master key to unlock the data, so it is
  physically blocked from responding to operations unsealed.
---

# operator seal

The `operator seal` seals the Vault server. Sealing tells the Vault server to
stop responding to any operations until it is unsealed. When sealed, the Vault
server discards its in-memory master key to unlock the data, so it is physically
blocked from responding to operations unsealed.

If an unseal is in progress, sealing the Vault will reset the unsealing process.
Users will have to re-enter their portions of the master key again.

This command does nothing if the Vault server is already sealed.

For more information on sealing and unsealing, please the [seal concepts
page](/docs/concepts/seal.html).

## Examples

Seal a Vault server:

```text
$ vault operator seal
Success! Vault is sealed.
```

## Usage

There are no flags beyond the [standard set of flags](/docs/commands/index.html)
included on all commands.
