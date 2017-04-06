---
layout: "docs"
page_title: "Vault Enterprise HSM Behavioral Changes"
sidebar_current: "docs-vault-enterprise-hsm-behavior"
description: |-
  Vault Enterprise HSM support changes the way Vault works with regard to unseal and recovery keys as well as rekey and recovery operations.
---

# Vault Enterprise HSM Behavioral Changes

This page contains information about the behavioral differences that take
effect when using Vault with an HSM.

## Key Split Between Unseal Keys and Recovery Keys

Normally, Vault uses a single set of unseal keys to perform both decryption of
the cryptographic barrier and to authorize recovery operations, such as the
[`generate-root`](/api/system/generate-root.html)
functionality.

When using an HSM, because the HSM automatically unseals the barrier but
recovery operations should still have human oversight, Vault instead uses two
sets of keys: unseal keys and recovery keys.

## Unseal (Master) Key

Vault usually generates a master key and splits it using [Shamir's Secret
Sharing](https://en.wikipedia.org/wiki/Shamir%27s_Secret_Sharing) to prevent a
single operator from being able to modify and unseal Vault (see more
information about Vault's security model
[here](/docs/internals/security.html)).

When using an HSM, Vault instead stores the master key, encrypted by the HSM,
into its internal storage. As a result, during an `init` command, the number of
key shares, threshold, and stored shares are required to be set to `1`, meaning
to not split the master key, so that the single key share is itself the master
key. (Vault does not do this automatically as it generally prefers to error
rather than change parameters set by an operator.)

Vault does not currently support rekeying the master key when protected by an
HSM; however, it _does_ continue to support rotation of the underlying data
encryption key that the master key protects via the
[`/sys/rotate`](/api/system/rotate.html) API
endpoint.

## Recovery Key

When Vault is initialized while using an HSM, rather than unseal keys being
returned to the operator, recovery keys are returned. These are generated from
an internal recovery key that is split via Shamir's Secret Sharing, similar to
Vault's treatment of unseal keys when running without an HSM.

Details about initialization and rekeying follow. When performing an operation
that uses recovery keys, such as `generate-root`, selection of the recovery
keys for this purpose, rather than the barrier unseal keys, is automatic.

### Initialization

When initializing, the split is performed according to the following CLI flags
and their API equivalents in the
[/sys/init](/api/system/init.html) endpoint:

 * `recovery-shares`: The number of shares into which to split the recovery
   key. This value is equivalent to the `recovery_shares` value in the API
   endpoint.
 * `recovery-threshold`: The threshold of shares required to reconstruct the
   recovery key. This value is equivalent to the `recovery_threshold` value in
   the API endpoint.
 * `recovery-pgp-keys`: The PGP keys to use to encrypt the returned recovery
   key shares. This value is equivalent to the `recovery_pgp_keys` value in the
   API endpoint, although as with `pgp_keys` the object in the API endpoint is
   an array, not a string.

Additionally, Vault will refuse to initialize if the option has not been set to
generate a key but no key is found. See
[Configuration](/docs/vault-enterprise/hsm/configuration.html) for more details.

### Rekeying

The recovery key can be rekeyed to change the number of shares/threshold or to
target different key holders via different PGP keys. When using the Vault CLI,
this is performed by using the `-recovery-key=true` flag to `vault rekey`.

Via the API, the rekey operation is performed with the same parameters as the
[normal `/sys/rekey`
endpoint](/api/system/rekey.html); however, the
API prefix for this operation is at `/sys/rekey-recovery-key` rather than
`/sys/rekey`.
