---
layout: "docs"
page_title: "Key Rotation"
sidebar_current: "docs-internals-rotation"
description: |-
  Learn about the details of key rotation within Vault.
---

# Key Rotation

Vault has multiple encryption keys that are used for various purposes. These keys support
rotation so that they can be periodically changed or in response to a potential leak or
compromise. It is useful to first understand the
[high-level architecture](/docs/internals/architecture.html) before learning about key rotation.

As a review, Vault starts in a _sealed_ state. Vault is unsealed by providing the unseal keys.
By default, Vault uses a technique known as [Shamir's secret sharing algorithm](https://en.wikipedia.org/wiki/Shamir's_Secret_Sharing)
to split the master key into 5 shares, any 3 of which are required to reconstruct the master
key. The master key is used to protect the encryption key, which is ultimately used to protect
data written to the storage backend.

[![Vault Shamir Secret Sharing Algorithm](/assets/images/vault-shamir-secret-sharing.svg)](/assets/images/vault-shamir-secret-sharing.svg)

To support key rotation, we need to support changing the unseal keys, master key, and the
backend encryption key. We split this into two separate operations, `rekey` and `rotate`.

The `rekey` operation is used to generate a new master key. When this is being done,
it is possible to change the parameters of the key splitting, so that the number of shares
and the threshold required to unseal can be changed. To perform a rekey a threshold of the
current unseal keys must be provided. This is to prevent a single malicious operator from
performing a rekey and invalidating the existing master key.

Performing a rekey is fairly straightforward. The rekey operation must be initialized with
the new parameters for the split and threshold. Once initialized, the current unseal keys
must be provided until the threshold is met. Once met, Vault will generate the new master
key, perform the splitting, and re-encrypt the encryption key with the new master key.
The new unseal keys are then provided to the operator, and the old unseal keys are no
longer usable.

The `rotate` operation is used to change the encryption key used to protect data written
to the storage backend. This key is never provided or visible to operators, who only
have unseal keys. This simplifies the rotation, as it does not require the current key
holders unlike the `rekey` operation. When `rotate` is triggered, a new encryption key
is generated and added to a keyring. All new values written to the storage backend are
encrypted with the new key. Old values written with previous encryption keys can still
be decrypted since older keys are saved in the keyring. This allows key rotation to be
done online, without an expensive re-encryption process.

Both the `rekey` and `rotate` operations can be done online and in a highly available
configuration. Only the active Vault instance can perform either of the operations
but standby instances can still assume an active role after either operation. This is
done by providing an online upgrade path for standby instances. If the current encryption
key is `N` and a rotation installs `N+1`, Vault creates a special "upgrade" key, which
provides the `N+1` encryption key protected by the `N` key. This upgrade key is only available
for a few minutes enabling standby instances to do a periodic check for upgrades.
This allows standby instances to update their keys and stay in-sync with the active Vault
without requiring operators to perform another unseal.
