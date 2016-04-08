---
layout: "docs"
page_title: "Seal/Unseal"
sidebar_current: "docs-concepts-seal"
description: |-
  A Vault must be unsealed before it can access its data. Likewise, it can be sealed to lock it down.
---

# Seal/Unseal

When a Vault server is started, it starts in a _sealed_ state. In this
state, Vault is configured to know where and how to access the physical
storage, but doesn't know how to decrypt any of it.

_Unsealing_ is the process of constructing the master key necessary to
read the decryption key to decrypt the data, allowing access to the Vault.

Prior to unsealing, almost no operations are possible with Vault. For
example authentication, managing the mount tables, etc. are all not possible.
The only possible operations are to unseal the Vault and check the status
of the unseal.

## Why?

The data stored by Vault is stored encrypted. Vault needs the
_encryption key_ in order to decrypt the data. The encryption key is
also stored with the data, but encrypted with another encryption key
known as the _master key_. The master key isn't stored anywhere.

Therefore, to decrypt the data, Vault must decrypt the encryption key
which requires the master key. Unsealing is the process of reconstructing
this master key.

Instead of distributing this master key as a single key to an operator,
Vault uses an algorithm known as
[Shamir's Secret Sharing](https://en.wikipedia.org/wiki/Shamir%27s_Secret_Sharing)
to split the key into shards. A certain threshold of shards is required to
reconstruct the master key.

This is the _unseal_ process: the shards are added one at a time (in any
order) until enough shards are present to reconstruct the key and
decrypt the data.

## Unsealing

The unseal process is done by running `vault unseal` or via the API.
This process is stateful: each key can be entered via multiple mechanisms
on multiple computers and it will work. This allows each shard of the master
key to be on a distinct machine for better security.

Once a Vault is unsealed, it remains unsealed until one of two things happens:

  1. It is resealed via the API (see below).

  2. The server is restarted.

-> **Note:** Unsealing makes the process of automating a Vault install
difficult. Automated tools can easily install, configure, and start Vault,
but unsealing it is a very manual process. We have plans in the future to
make it easier. For the time being, the best method is to manually unseal
multiple Vault servers in [HA mode](/docs/concepts/ha.html). Use a tool such
as Consul to make sure you only query Vault servers that are unsealed.

## Sealing

There is also an API to seal the Vault. This will throw away the master
key and require another unseal process to restore it. Sealing only requires
a single operator with root privileges.

This way, if there is a detected intrusion, the Vault data can be locked
quickly to try to minimize damages. It can't be accessed again without
access to the master key shards.
