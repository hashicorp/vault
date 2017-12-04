---
layout: "api"
page_title: "/sys/rekey - HTTP API"
sidebar_current: "docs-http-system-rekey"
description: |-
  The `/sys/rekey` endpoints are used to rekey the unseal keys for Vault.
---

# `/sys/rekey`

The `/sys/rekey` endpoints are used to rekey the unseal keys for Vault.

On seals that support stored keys (e.g. HSM PKCS11), the recovery key share(s)
can be provided to rekey the master key since no unseal keys are available. The
secret shares, secret threshold, and stored shares parameteres must be set to 1.
Upon successful rekey, no split unseal key shares are returned.

## Read Rekey Progress

This endpoint reads the configuration and progress of the current rekey attempt.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `GET`    | `/sys/rekey/init`            | `200 application/json` |

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    https://vault.rocks/v1/sys/rekey/init
```

### Sample Response

```json
{
  "started": true,
  "nonce": "2dbd10f1-8528-6246-09e7-82b25b8aba63",
  "t": 3,
  "n": 5,
  "progress": 1,
  "required": 3,
  "pgp_fingerprints": ["abcd1234"],
  "backup": true
}
```

If a rekey is started, then `n` is the new shares to generate and `t` is the
threshold required for the new shares. `progress` is how many unseal keys have
been provided for this rekey, where `required` must be reached to complete. The
`nonce` for the current rekey operation is also displayed. If PGP keys are being
used to encrypt the final shares, the key fingerprints and whether the final
keys will be backed up to physical storage will also be displayed.


## Start Rekey

This endpoint initializes a new rekey attempt. Only a single rekey attempt can
take place at a time, and changing the parameters of a rekey requires canceling
and starting a new rekey, which will also provide a new nonce.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `PUT`    | `/sys/rekey/init`            | `204 (empty body)`     |

### Parameters

- `secret_shares` `(int: <required>)` – Specifies the number of shares to split
  the master key into.

- `secret_threshold` `(int: <required>)` – Specifies the number of shares
  required to reconstruct the master key. This must be less than or equal to
  `secret_shares`.

- `pgp_keys` `(array<string>: nil)` – Specifies an array of PGP public keys used
  to encrypt the output unseal keys. Ordering is preserved. The keys must be
  base64-encoded from their original binary representation. The size of this
  array must be the same as `secret_shares`.

- `backup` `(bool: false)` – Specifies if using PGP-encrypted keys, whether
  Vault should also store a plaintext backup of the PGP-encrypted keys at
  `core/unseal-keys-backup` in the physical storage backend. These can then
  be retrieved and removed via the `sys/rekey/backup` endpoint.

### Sample Payload

```json
{
  "secret_shares": 10,
  "secret_threshold": 5
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request PUT \
    --data @payload.json \
    https://vault.rocks/v1/sys/rekey/init
```

## Cancel Rekey

This endpoint cancels any in-progress rekey. This clears the rekey settings as
well as any progress made. This must be called to change the parameters of the
rekey.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `DELETE` | `/sys/rekey/init`            | `204 (empty body)`     |

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request DELETE \
    https://vault.rocks/v1/sys/rekey/init
```

## Read Backup Key

This endpoint returns the backup copy of PGP-encrypted unseal keys. The returned
value is the nonce of the rekey operation and a map of PGP key fingerprint to
hex-encoded PGP-encrypted key.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `GET`    | `/sys/rekey/backup`          | `200 application/json` |

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    https://vault.rocks/v1/sys/rekey/backup
```

### Sample Response

```json
{
  "nonce": "2dbd10f1-8528-6246-09e7-82b25b8aba63",
  "keys": {
    "abcd1234": "..."
  }
}
```

## Delete Backup Key

This endpoint deletes the backup copy of PGP-encrypted unseal keys.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `DELETE` | `/sys/rekey/backup`          | `204 (empty body)`     |

### Sample Request

```
$ curl \
    --header "X-Vault-Token" \
    --request DELETE \
    https://vault.rocks/v1/sys/rekey/backup
```

## Submit Key

This endpoint is used to enter a single master key share to progress the rekey
of the Vault. If the threshold number of master key shares is reached, Vault
will complete the rekey. Otherwise, this API must be called multiple times until
that threshold is met. The rekey nonce operation must be provided with each
call.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `PUT`    | `/sys/rekey/update`          | `200 application/json` |

### Parameters

- `key` `(string: <required>)` – Specifies a single master share key.

- `nonce` `(string: <required>)` – Specifies the nonce of the rekey operation.

### Sample Payload

```json
{
  "key": "abcd1234...",
  "nonce": "AB32..."
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token" \
    --request PUT \
    --data @payload.json \
    https://vault.rocks/v1/sys/rekey/update
```

### Sample Response

```json
{
  "complete": true,
  "keys": ["one", "two", "three"],
  "nonce": "2dbd10f1-8528-6246-09e7-82b25b8aba63",
  "pgp_fingerprints": ["abcd1234"],
  "keys_base64": ["base64keyvalue"],
  "backup": true
}
```

If the keys are PGP-encrypted, an array of key fingerprints will also be
provided (with the order in which the keys were used for encryption) along with
whether or not the keys were backed up to physical storage.
