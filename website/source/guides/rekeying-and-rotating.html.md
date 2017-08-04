---
layout: "guides"
page_title: "Rekeying & Rotating Vault - Guides"
sidebar_current: "guides-rekeying-and-rotating"
description: |-
  Vault supports generating new unseal keys as well as rotating the underlying
  encryption keys. This guide covers rekeying and rotating Vault's encryption
  keys.
---

# Rekeying &amp; Rotating Vault

~> **Advanced Topic** This guide presents an advanced topic that is not required
for a basic understanding of Vault. Knowledge of this topic is not required for
daily Vault use.

## Background

In order to prevent no one person from having complete access to the system,
Vault employs [Shamir's Secret Sharing Algorithm][shamir]. Under this process,
a secret is divided into a subset of parts such that a subset of those parts are
needed to reconstruct the original secret. Vault makes heavy use of this
algorithm as part of the [unsealing process](/docs/concepts/seal.html).

When a Vault server is first initialized, Vault generates a master key and
immediately splits this master key into a series of key shares following
Shamir's Secret Sharing Algorithm. Vault never stores the master key, therefore,
the only way to retrieve the master key is to have a quorum of unseal keys
re-generate it.

The master key is used to decrypt the underlying encryption key. Vault uses the
encryption key to encrypt data at rest in a storage backend like the filesystem
or Consul.

Typically each of these key shares is distributed to trusted parties in the
organization. These parties must come together to "unseal" the Vault by entering
their key share.

[![Vault Shamir Secret Sharing Algorithm](/assets/images/vault-shamir-secret-sharing.svg)](/assets/images/vault-shamir-secret-sharing.svg)

[shamir]: https://en.wikipedia.org/wiki/Shamir%27s_Secret_Sharing

In some cases, you may want to re-generate the master key and key shares. Here
are a few examples:

- Someone joins or leaves the organization
- Security wants to change the number of shares or threshold of shares
- Compliance mandates the master key be rotated at a regular interval

In addition to rekeying the master key, there may be an independent desire to
rotate the underlying encryption key Vault uses to encrypt data at rest.

[![Vault Rekey vs Rotate](/assets/images/vault-rekey-vs-rotate.svg)](/assets/images/vault-rekey-vs-rotate.svg)

In Vault, _rekeying_ and _rotating_ are two separate operations. The process for
generating a new master key and applying Shamir's algorithm is called
"rekeying". The process for generating a new encryption key for Vault to encrypt
data at rest is called "rotating".

Both rekeying the Vault and rotating Vault's underlying encryption key are fully
online operations. Vault will continue to service requests uninterrupted during
either of these processes.

## Rekeying Vault

Rekeying the Vault requires a quorum of unseal keys. Before continuing, you
should ensure all unseal key holders are available to assist with the rekeying.

First, initialize a rekeying operation. The flags represent the **newly
desired** number of keys and threshold:

```text
$ vault rekey -init -key-shares=3 -key-threshold=2
```

This will generate a nonce value and start the rekeying process. All other
unseal keys must also provide this nonce value. This nonce value is not a
secret, so it is safe to distribute over insecure channels like chat, email, or
carrier pigeon.

```text
Nonce: 22657753-9cca-189a-65b8-cb743d104ffc
Started: true
Key Shares: 3
Key Threshold: 2
Rekey Progress: 0
Required Keys: 1
```

Each unseal key holder runs the following command and enters their unseal key:

```text
$ vault rekey -nonce=<nonce>
Rekey operation nonce: 22657753-9cca-189a-65b8-cb743d104ffc
Key (will be hidden):
```

When the final unseal key holder enters their key, Vault will output the new
unseal keys:

```text
Key 1: EDj4NZK6z5Y9rpr+TtihTulfdHvFzXtBYQk36dmBczuQ
Key 2: sCkM1i5BGGNDFk5GsqtVolWRPyd5mWn2eZG0gUySiCF7
Key 3: e5DUvDIH0cPU8Q+hh1KNVkkMc9lliliPVe9u3Fzbzv38

Operation nonce: 22657753-9cca-189a-65b8-cb743d104ffc

Vault rekeyed with 3 keys and a key threshold of 2. Please
securely distribute the above keys. When the vault is re-sealed,
restarted, or stopped, you must provide at least 2 of these keys
to unseal it again.

Vault does not store the master key. Without at least 2 keys,
your vault will remain permanently sealed.
```

Like the initialization process, Vault supports PGP encrypting the resulting
unseal keys and creating backup encryption keys for disaster recovery.

## Rotating the Encryption Key

Unlike rekeying the Vault, rotating Vault's encryption key does not require a
quorum of unseal keys. Anyone with the proper permissions in Vault can perform
the encryption key rotation.

To trigger a key rotation, execute the command:

```text
$ vault rotate
```

This will output the key version and installation time:

```text
Key Term: 2
Installation Time: ...
```

This will add a new key to the keyring. All new values written to the storage
backend will be encrypted with this new key.
