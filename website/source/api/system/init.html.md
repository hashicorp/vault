---
layout: "api"
page_title: "/sys/init - HTTP API"
sidebar_current: "docs-http-system-init"
description: |-
  The `/sys/init` endpoint is used to initialize a new Vault.
---

# `/sys/init`

The `/sys/init` endpoint is used to initialize a new Vault.

## Read Initialization Status

This endpoint returns the initialization status of Vault.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `GET`    | `/sys/init`                  | `200 application/json` |

### Sample Request

```
$ curl \
    https://vault.rocks/v1/sys/init
```

### Sample Response

```json
{
  "initialized": true
}
```

## Start Initialization

This endpoint initializes a new Vault. The Vault must not have been previously
initialized. The recovery options, as well as the stored shares option, are only
available when using Vault HSM.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `PUT`    | `/sys/init`                  | `200 application/json` |

### Parameters

- `pgp_keys` `(array<string>: nil)` – Specifies an array of PGP public keys used
  to encrypt the output unseal keys. Ordering is preserved. The keys must be
  base64-encoded from their original binary representation. The size of this
  array must be the same as `secret_shares`.

- `root_token_pgp_key` `(string: "")` – Specifies a PGP public key used to
  encrypt the initial root token. The key must be base64-encoded from its
  original binary representation.

- `secret_shares` `(int: <required>)` – Specifies the number of shares to
  split the master key into.

- `secret_threshold` `(int: <required>)` – Specifies the number of shares
  required to reconstruct the master key. This must be less than or equal
  `secret_shares`. If using Vault HSM with auto-unsealing, this value must be
  the same as `secret_shares`.

Additionally, the following options are only supported on Vault Pro/Enterprise:

- `stored_shares` `(int: <required>)` – Specifies the number of shares that
  should be encrypted by the HSM and stored for auto-unsealing. Currently must
  be the same as `secret_shares`.

- `recovery_shares` `(int: <required>)` – Specifies rhe number of shares to
  split the recovery key into.

- `recovery_threshold` `(int: <required>)` – Specifies rhe number of shares
  required to reconstruct the recovery key. This must be less than or equal to
  `recovery_shares`.

- `recovery_pgp_keys` `(array<string>: nil)` – Specifies an array of PGP public
  keys used to encrypt the output recovery keys. Ordering is preserved. The keys
  must be base64-encoded from their original binary representation. The size of
  this array must be the same as `recovery_shares`.

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
    --request PUT \
    --data @payload.json \
    https://vault.rocks/v1/sys/init
```

### Sample Response

A JSON-encoded object including the (possibly encrypted, if `pgp_keys` was
provided) master keys, base 64 encoded master keys and initial root token:

```json
{
  "keys": ["one", "two", "three"],
  "keys_base64": ["cR9No5cBC", "F3VLrkOo", "zIDSZNGv"],
  "root_token": "foo"
}
```
