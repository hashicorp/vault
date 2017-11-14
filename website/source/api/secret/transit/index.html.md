---
layout: "api"
page_title: "Transit Secret Backend - HTTP API"
sidebar_current: "docs-http-secret-transit"
description: |-
  This is the API documentation for the Vault Transit secret backend.
---

# Transit Secret Backend HTTP API

This is the API documentation for the Vault Transit secret backend. For general
information about the usage and operation of the Transit backend, please see the
[Vault Transit backend documentation](/docs/secrets/transit/index.html).

This documentation assumes the Transit backend is mounted at the `/transit`
path in Vault. Since it is possible to mount secret backends at any location,
please update your API calls accordingly.

## Create Key

This endpoint creates a new named encryption key of the specified type. The
values set here cannot be changed after key creation.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `POST`   | `/transit/keys/:name`        | `204 (empty body)`     |

### Parameters

- `name` `(string: <required>)` – Specifies the name of the encryption key to
  create. This is specified as part of the URL.

- `convergent_encryption` `(bool: false)` – If enabled, the key will support
  convergent encryption, where the same plaintext creates the same ciphertext.
  This requires _derived_ to be set to `true`. When enabled, each
  encryption(/decryption/rewrap/datakey) operation will derive a `nonce` value
  rather than randomly generate it. Note that while this is useful for
  particular situations, all nonce values used with a given context value **must
  be unique** or it will compromise the security of your key, and the key space
  for nonces is 96 bit -- not as large as the AES key itself.

- `derived` `(bool: false)` – Specifies if key derivation is to be used. If
  enabled, all encrypt/decrypt requests to this named key must provide a context
  which is used for key derivation.

- `exportable` `(bool: false)` – Specifies if the raw key is exportable.

- `type` `(string: "aes256-gcm96")` – Specifies the type of key to create. The
  currently-supported types are:

    - `aes256-gcm96` – AES-256 wrapped with GCM using a 12-byte nonce size
      (symmetric, supports derivation)
    - `ecdsa-p256` – ECDSA using the P-256 elliptic curve (asymmetric)
    - `ed25519` – ED25519 (asymmetric, supports derivation)
    - `rsa-2048` - RSA with bit size of 2048 (asymmetric)
    - `rsa-4096` - RSA with bit size of 4096 (asymmetric)

### Sample Payload

```json
{
  "type": "ecdsa-p256",
  "derived": true
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    https://vault.rocks/v1/transit/keys/my-key
```

## Read Key

This endpoint returns information about a named encryption key. The `keys`
object shows the creation time of each key version; the values are not the keys
themselves. Depending on the type of key, different information may be returned,
e.g. an asymmetric key will return its public key in a standard format for the
type.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `GET`    | `/transit/keys/:name`        | `200 application/json` |

### Parameters

- `name` `(string: <required>)` – Specifies the name of the encryption key to
  read. This is specified as part of the URL.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    https://vault.rocks/v1/transit/keys/my-key
```

### Sample Response

```json
{
  "data": {
    "type": "aes256-gcm96",
    "deletion_allowed": false,
    "derived": false,
    "exportable": false,
    "keys": {
      "1": 1442851412
    },
    "min_decryption_version": 1,
    "min_encryption_version": 0,
    "name": "foo",
    "supports_encryption": true,
    "supports_decryption": true,
    "supports_derivation": true,
    "supports_signing": false
  }
}
```

## List Keys

This endpoint returns a list of keys. Only the key names are returned (not the
actual keys themselves).

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `LIST`   | `/transit/keys`              | `200 application/json` |
| `GET`    | `/transit/keys?list=true`    | `200 application/json` |

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request LIST \
    https://vault.rocks/v1/transit/keys
```

### Sample Response

```json
{
  "data": {
    "keys": ["foo", "bar"]
  },
  "lease_duration": 0,
  "lease_id": "",
  "renewable": false
}
```

## Delete Key

This endpoint deletes a named encryption key. It will no longer be possible to
decrypt any data encrypted with the named key. Because this is a potentially
catastrophic operation, the `deletion_allowed` tunable must be set in the key's
`/config` endpoint.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `DELETE` | `/transit/keys/:name`        | `204 (empty body)`     |

### Parameters

- `name` `(string: <required>)` – Specifies the name of the encryption key to
  delete. This is specified as part of the URL.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request DELETE \
    https://vault.rocks/v1/transit/keys/my-key
```

## Update Key Configuration

This endpoint allows tuning configuration values for a given key. (These values
are returned during a read operation on the named key.)

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `POST`   | `/transit/keys/:name/config` | `204 (empty body)`     |

### Parameters

- `min_decryption_version` `(int: 0)` – Specifies the minimum version of
  ciphertext allowed to be decrypted. Adjusting this as part of a key rotation
  policy can prevent old copies of ciphertext from being decrypted, should they
  fall into the wrong hands. For signatures, this value controls the minimum
  version of signature that can be verified against. For HMACs, this controls
  the minimum version of a key allowed to be used as the key for verification.

- `min_encryption_version` `(int: 0)` – Specifies the minimum version of the
  key that can be used to encrypt plaintext, sign payloads, or generate HMACs.
  Must be `0` (which will use the latest version) or a value greater or equal
  to `min_decryption_version`.

- `deletion_allowed` `(bool: false)`- Specifies if the key is allowed to be
  deleted.

### Sample Payload

```json
{
  "deletion_allowed": true
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    https://vault.rocks/v1/transit/keys/my-key/config
```

## Rotate Key

This endpoint rotates the version of the named key. After rotation, new
plaintext requests will be encrypted with the new version of the key. To upgrade
ciphertext to be encrypted with the latest version of the key, use the `rewrap`
endpoint. This is only supported with keys that support encryption and
decryption operations.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `POST`   | `/transit/keys/:name/rotate` | `204 (empty body)`     |

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    https://vault.rocks/v1/transit/keys/my-key/rotate
```

## Export Key

This endpoint returns the named key. The `keys` object shows the value of the
key for each version. If `version` is specified, the specific version will be
returned. If `latest` is provided as the version, the current key will be
provided. Depending on the type of key, different information may be returned.
The key must be exportable to support this operation and the version must still
be valid.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `GET`    | `/transit/export/:key_type/:name(/:version)` | `200 application/json` |

### Parameters

- `key_type` `(string: <required>)` – Specifies the type of the key to export.
  This is specified as part of the URL. Valid values are:

    - `encryption-key`
    - `signing-key`
    - `hmac-key`

- `name` `(string: <required>)` – Specifies the name of the key to read
  information about. This is specified as part of the URL.

- `version` `(string: "")` – Specifies the version of the key to read. If omitted,
  all versions of the key will be returned. This is specified as part of the
  URL. If the version is set to `latest`, the current key will be returned.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    https://vault.rocks/v1/transit/export/encryption-key/my-key/1
```

### Sample Response

```json
{
  "data": {
    "name": "foo",
    "keys": {
      "1": "eyXYGHbTmugUJn6EtYD/yVEoF6pCxm4R/cMEutUm3MY=",
      "2": "Euzymqx6iXjS3/NuGKDCiM2Ev6wdhnU+rBiKnJ7YpHE="
    }
  }
}
```

## Encrypt Data

This endpoint encrypts the provided plaintext using the named key. Currently,
this only supports symmetric keys. This path supports the `create` and `update`
policy capabilities as follows: if the user has the `create` capability for this
endpoint in their policies, and the key does not exist, it will be upserted with
default values (whether the key requires derivation depends on whether the
context parameter is empty or not). If the user only has `update` capability and
the key does not exist, an error will be returned.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `POST`   | `/transit/encrypt/:name`     | `200 application/json` |

### Parameters

- `name` `(string: <required>)` – Specifies the name of the encryption key to
  encrypt against. This is specified as part of the URL.

- `plaintext` `(string: <required>)` – Specifies **base64 encoded** plaintext to
  be encoded.

- `context` `(string: "")` – Specifies the **base64 encoded** context for key
  derivation. This is required if key derivation is enabled for this key.

- `key_version` `(int: 0)` – Specifies the version of the key to use for
  encryption. If not set, uses the latest version. Must be greater than or
  equal to the key's `min_encryption_version`, if set.

- `nonce` `(string: "")` – Specifies the **base64 encoded** nonce value. This
  must be provided if convergent encryption is enabled for this key and the key
  was generated with Vault 0.6.1. Not required for keys created in 0.6.2+. The
  value must be exactly 96 bits (12 bytes) long and the user must ensure that
  for any given context (and thus, any given encryption key) this nonce value is
  **never reused**.

- `batch_input` `(array<object>: nil)` – Specifies a list of items to be
  encrypted in a single batch. When this parameter is set, if the parameters
  'plaintext', 'context' and 'nonce' are also set, they will be ignored. The
  format for the input is:

    ```json
    [
      {
        "context": "c2FtcGxlY29udGV4dA==",
        "plaintext": "dGhlIHF1aWNrIGJyb3duIGZveA=="
      },
      {
        "context": "YW5vdGhlcnNhbXBsZWNvbnRleHQ=",
        "plaintext": "dGhlIHF1aWNrIGJyb3duIGZveA=="
      },
    ]
    ```

- `type` `(string: "aes256-gcm96")` –This parameter is required when encryption
  key is expected to be created. When performing an upsert operation, the type
  of key to create. Currently, "aes256-gcm96" (symmetric) is the only type
  supported.

- `convergent_encryption` `(string: "")` – This parameter will only be used when
  a key is expected to be created.  Whether to support convergent encryption.
  This is only supported when using a key with key derivation enabled and will
  require all requests to carry both a context and 96-bit (12-byte) nonce. The
  given nonce will be used in place of a randomly generated nonce. As a result,
  when the same context and nonce are supplied, the same ciphertext is
  generated. It is _very important_ when using this mode that you ensure that
  all nonces are unique for a given context.  Failing to do so will severely
  impact the ciphertext's security.

### Sample Payload

```json
{
  "plaintext": "dGhlIHF1aWNrIGJyb3duIGZveA=="
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    https://vault.rocks/v1/transit/encrypt/my-key
```

### Sample Response

```json
{
  "data": {
    "ciphertext": "vault:v1:abcdefgh"
  }
}
```

## Decrypt Data

This endpoint decrypts the provided ciphertext using the named key. Currently,
this only supports symmetric keys.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `POST`   | `/transit/decrypt/:name`     | `200 application/json` |

### Parameters

- `name` `(string: <required>)` – Specifies the name of the encryption key to
  decrypt against. This is specified as part of the URL.

- `ciphertext` `(string: <required>)` – Specifies the ciphertext to decrypt.

- `context` `(string: "")` – Specifies the **base64 encoded** context for key
  derivation. This is required if key derivation is enabled.

- `nonce` `(string: "")` – Specifies a base64 encoded nonce value used during
  encryption. Must be provided if convergent encryption is enabled for this key
  and the key was generated with Vault 0.6.1. Not required for keys created in
  0.6.2+.

- `batch_input` `(array<object>: nil)` – Specifies a list of items to be
  decrypted in a single batch. When this parameter is set, if the parameters
  'ciphertext', 'context' and 'nonce' are also set, they will be ignored. Format
  for the input goes like this:

    ```json
    [
      {
        "context": "c2FtcGxlY29udGV4dA==",
        "ciphertext": "vault:v1:/DupSiSbX/ATkGmKAmhqD0tvukByrx6gmps7dVI="
      },
      {
        "context": "YW5vdGhlcnNhbXBsZWNvbnRleHQ=",
        "ciphertext": "vault:v1:XjsPWPjqPrBi1N2Ms2s1QM798YyFWnO4TR4lsFA="
      },
    ]
    ```

### Sample Payload

```json
{
  "ciphertext": "vault:v1:XjsPWPjqPrBi1N2Ms2s1QM798YyFWnO4TR4lsFA="
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    https://vault.rocks/v1/transit/decrypt/my-key
```

### Sample Response

```json
{
  "data": {
    "plaintext": "dGhlIHF1aWNrIGJyb3duIGZveAo="
  }
}
```

## Rewrap Data

This endpoint rewraps the provided ciphertext using the latest version of the
named key. Because this never returns plaintext, it is possible to delegate this
functionality to untrusted users or scripts.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `POST`   | `/transit/rewrap/:name`      | `200 application/json` |

### Parameters

- `name` `(string: <required>)` – Specifies the name of the encryption key to
  re-encrypt against. This is specified as part of the URL.

- `ciphertext` `(string: <required>)` – Specifies the ciphertext to re-encrypt.

- `context` `(string: "")` – Specifies the **base64 encoded** context for key
  derivation. This is required if key derivation is enabled.

- `key_version` `(int: 0)` – Specifies the version of the key to use for the
  operation. If not set, uses the latest version. Must be greater than or equal
  to the key's `min_encryption_version`, if set.

- `nonce` `(string: "")` – Specifies a base64 encoded nonce value used during
  encryption. Must be provided if convergent encryption is enabled for this key
  and the key was generated with Vault 0.6.1. Not required for keys created in
  0.6.2+.

- `batch_input` `(array<object>: nil)` – Specifies a list of items to be
  decrypted in a single batch. When this parameter is set, if the parameters
  'ciphertext', 'context' and 'nonce' are also set, they will be ignored. Format
  for the input goes like this:

    ```json
    [
      {
        "context": "c2FtcGxlY29udGV4dA==",
        "ciphertext": "vault:v1:/DupSiSbX/ATkGmKAmhqD0tvukByrx6gmps7dVI="
      },
      {
        "context": "YW5vdGhlcnNhbXBsZWNvbnRleHQ=",
        "ciphertext": "vault:v1:XjsPWPjqPrBi1N2Ms2s1QM798YyFWnO4TR4lsFA="
      },
    ]
    ```

### Sample Payload

```json
{
  "ciphertext": "vault:v1:XjsPWPjqPrBi1N2Ms2s1QM798YyFWnO4TR4lsFA="
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    https://vault.rocks/v1/transit/rewrap/my-key
```

### Sample Response

```json
{
  "data": {
    "ciphertext": "vault:v2:abcdefgh"
  }
}
```

## Generate Data Key

This endpoint generates a new high-entropy key and the value encrypted with the
named key. Optionally return the plaintext of the key as well. Whether plaintext
is returned depends on the path; as a result, you can use Vault ACL policies to
control whether a user is allowed to retrieve the plaintext value of a key. This
is useful if you want an untrusted user or operation to generate keys that are
then made available to trusted users.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `POST`   | `/transit/datakey/:type/:name` | `200 application/json` |

### Parameters

- `type` `(string: <required>)` – Specifies the type of key to generate. If
  `plaintext`, the plaintext key will be returned along with the ciphertext. If
  `wrapped`, only the ciphertext value will be returned. This is specified as
  part of the URL.

- `name` `(string: <required>)` – Specifies the name of the encryption key to
  use to encrypt the datakey. This is specified as part of the URL.

- `context` `(string: "")` – Specifies the key derivation context, provided as a
  base64-encoded string. This must be provided if derivation is enabled.

- `nonce` `(string: "")` – Specifies a nonce value, provided as base64 encoded.
  Must be provided if convergent encryption is enabled for this key and the key
  was generated with Vault 0.6.1. Not required for keys created in 0.6.2+. The
  value must be exactly 96 bits (12 bytes) long and the user must ensure that
  for any given context (and thus, any given encryption key) this nonce value is
  **never reused**.

- `bits` `(int: 256)` – Specifies the number of bits in the desired key. Can be
  128, 256, or 512.

### Sample Payload

```json
{
  "context": "Ab3=="
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    https://vault.rocks/v1/transit/datakey/plaintext/my-key
```

### Sample Response

```json
{
  "data": {
    "plaintext": "dGhlIHF1aWNrIGJyb3duIGZveAo=",
    "ciphertext": "vault:v1:abcdefgh"
  }
}
```

## Generate Random Bytes

This endpoint returns high-quality random bytes of the specified length.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `POST`   | `/transit/random(/:bytes)`   | `200 application/json` |

### Parameters

- `bytes` `(int: 32)` – Specifies the number of bytes to return. This value can
  be specified either in the request body, or as a part of the URL.

- `format` `(string: "base64")` – Specifies the output encoding. Valid options
  are `hex` or `base64`.

### Sample Payload

```json
{
  "format": "hex"
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    https://vault.rocks/v1/transit/random/164
```

### Sample Response

```json
{
  "data": {
    "random_bytes": "dGhlIHF1aWNrIGJyb3duIGZveAo="
  }
}
```

## Hash Data

This endpoint returns the cryptographic hash of given data using the specified
algorithm.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `POST`   | `/transit/hash(/:algorithm)` | `200 application/json` |

### Parameters

- `algorithm` `(string: "sha2-256")` – Specifies the hash algorithm to use. This
  can also be specified as part of the URL. Currently-supported algorithms are:

    - `sha2-224`
    - `sha2-256`
    - `sha2-384`
    - `sha2-512`

- `input` `(string: <required>)` – Specifies the **base64 encoded** input data.

- `format` `(string: "hex")` – Specifies the output encoding. This can be either
  `hex` or `base64`.

### Sample Payload

```json
{
  "input": "adba32=="
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    https://vault.rocks/v1/transit/hash/sha2-512
```

### Sample Response

```json
{
  "data": {
    "sum": "dGhlIHF1aWNrIGJyb3duIGZveAo="
  }
}
```

## Generate HMAC

This endpoint returns the digest of given data using the specified hash
algorithm and the named key. The key can be of any type supported by `transit`;
the raw key will be marshaled into bytes to be used for the HMAC function. If
the key is of a type that supports rotation, the latest (current) version will
be used.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `POST`   | `/transit/hmac/:name(/:algorithm)` | `200 application/json` |

### Parameters

- `name` `(string: <required>)` – Specifies the name of the encryption key to
  generate hmac against. This is specified as part of the URL.

- `key_version` `(int: 0)` – Specifies the version of the key to use for the
  operation. If not set, uses the latest version. Must be greater than or equal
  to the key's `min_encryption_version`, if set.

- `algorithm` `(string: "sha2-256")` – Specifies the hash algorithm to use. This
  can also be specified as part of the URL. Currently-supported algorithms are:

    - `sha2-224`
    - `sha2-256`
    - `sha2-384`
    - `sha2-512`

- `input` `(string: <required>)` – Specifies the **base64 encoded** input data.

### Sample Payload

```json
{
  "input": "adba32=="
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    https://vault.rocks/v1/transit/hmac/my-key/sha2-512
```

### Sample Response

```json
{
  "data": {
    "hmac": "dGhlIHF1aWNrIGJyb3duIGZveAo="
  }
}
```

## Sign Data

This endpoint returns the cryptographic signature of the given data using the
named key and the specified hash algorithm. The key must be of a type that
supports signing.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `POST`   | `/transit/sign/:name(/:algorithm)` | `200 application/json` |

### Parameters

- `name` `(string: <required>)` – Specifies the name of the encryption key to
  use for signing. This is specified as part of the URL.

- `key_version` `(int: 0)` – Specifies the version of the key to use for
  signing. If not set, uses the latest version. Must be greater than or equal
  to the key's `min_encryption_version`, if set.

- `algorithm` `(string: "sha2-256")` – Specifies the hash algorithm to use for
  supporting key types (notably, not including `ed25519` which specifies its
  own hash algorithm). This can also be specified as part of the URL.
  Currently-supported algorithms are:

    - `none`
    - `sha2-224`
    - `sha2-256`
    - `sha2-384`
    - `sha2-512`

- `input` `(string: <required>)` – Specifies the **base64 encoded** input data.

- `context` `(string: "")` - Base64 encoded context for key derivation.
   Required if key derivation is enabled; currently only available with ed25519
   keys.

 - `prehashed` `(bool: false)` - Set to `true` when the input is already
   hashed. If the key type is `rsa-2048` or `rsa-4096`, then the algorithm used
   to hash the input should be indicated by the `algorithm` parameter.


### Sample Payload

```json
{
  "input": "adba32=="
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    https://vault.rocks/v1/transit/sign/my-key/sha2-512
```

### Sample Response

```json
{
  "data": {
    "signature": "vault:v1:MEUCIQCyb869d7KWuA0hBM9b5NJrmWzMW3/pT+0XYCM9VmGR+QIgWWF6ufi4OS2xo1eS2V5IeJQfsi59qeMWtgX0LipxEHI="
  }
}
```

## Verify Signed Data

This endpoint returns whether the provided signature is valid for the given
data.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `POST`   | `/transit/verify/:name(/:algorithm)` | `200 application/json` |

### Parameters

- `name` `(string: <required>)` – Specifies the name of the encryption key that
  was used to generate the signature or HMAC.

- `algorithm` `(string: "sha2-256")` – Specifies the hash algorithm to use. This
  can also be specified as part of the URL. Currently-supported algorithms are:

    - `none`
    - `sha2-224`
    - `sha2-256`
    - `sha2-384`
    - `sha2-512`

- `input` `(string: <required>)` – Specifies the **base64 encoded** input data.

- `signature` `(string: "")` – Specifies the signature output from the
  `/transit/sign` function. Either this must be supplied or `hmac` must be
  supplied.

- `hmac` `(string: "")` – Specifies the signature output from the
  `/transit/hmac` function. Either this must be supplied or `signature` must be
  supplied.

 - `context` `(string: "")` - Base64 encoded context for key derivation.
   Required if key derivation is enabled; currently only available with ed25519
   keys.

 - `prehashed` `(bool: false)` - Set to `true` when the input is already
   hashed. If the key type is `rsa-2048` or `rsa-4096`, then the algorithm used
   to hash the input should be indicated by the `algorithm` parameter.

### Sample Payload

```json
{
  "input": "abcd13==",
  "signature": "vault:v1:MEUCIQCyb869d7KWuA..."
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    https://vault.rocks/v1/transit/verify/my-key/sha2-512
```

### Sample Response

```json
{
  "data": {
    "valid": true
  }
}
```
