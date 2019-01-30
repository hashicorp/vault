---
layout: "api"
page_title: "Google Cloud KMS - Secrets Engines - HTTP API"
sidebar_title: "Google Cloud KMS"
sidebar_current: "api-http-secret-gcpkms"
description: |-
  This is the API documentation for the Vault Google Cloud KMS secrets engine.
---

# Google Cloud KMS Secrets Engine (API)

This is the API documentation for the Vault Google Cloud KMS secrets engine. For
general information about the usage and operation of the Google Cloud KMS
secrets engine, please see the
[Google Cloud KMS documentation](/docs/secrets/gcpkms/index.html).

This documentation assumes the Google Cloud KMS secrets engine is enabled at the
`/gcpkms` path in Vault. Since it is possible to enable secrets engines at any
location, please update your API calls accordingly.

## Configure Credentials

This endpoint configures the Google Cloud KMS secrets engine with credentials
and manages the requested scope(s) for authentication.

| Method   | Path                     | Produces                  |
| :------- | :------------------------| :------------------------ |
| `POST`   | `gcpkms/config`          | `204 (empty body)`        |

### Example Policy

```hcl
path "gcpkms/config" {
  capabilities = ["create", "update"]
}
```

### Parameters

- `credentials` (`string: ""`) -
The credentials to use for authenticating to Google Cloud. Leave this blank to
use the Default Application Credentials or instance metadata authentication.

- `scopes` (`array<string>: []`) -
The list of full-URL scopes to request when authenticating. By default, this
requests https://www.googleapis.com/auth/cloudkms.

### Sample Payload

```json
{
  "credentials": "< JSON credentials... >"
}
```

### Sample Request

```text
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    https://127.0.0.1:8200/v1/gcpkms/config
```

## Read Configuration

This endpoint returns the configuration endpoint for the Google Cloud KMS
secrets engine. The credentials are not returned.

| Method   | Path                     | Produces                  |
| :------- | :------------------------| :------------------------ |
| `GET`    | `gcpkms/config`          | `200 application/json` |

### Example Policy

```hcl
path "gcpkms/config" {
  capabilities = ["read"]
}
```

### Sample Request

```text
$ curl \
    --header "X-Vault-Token: ..." \
    --request GET \
    https://127.0.0.1:8200/v1/gcpkms/config
```

### Sample Response

```json
{
  "data": {
    "scopes": [
      "https://www.googleapis.com/auth/cloudkms"
    ]
  }
}
```

## Delete Configuration

This endpoint deletes any configuration for the Google Cloud KMS secrets engine.
If there is no configuration, the endpoint still returns successfully.

| Method   | Path                     | Produces                  |
| :------- | :------------------------| :------------------------ |
| `DELETE` | `gcpkms/config`          | `204 (empty body)`        |

### Example Policy

```hcl
path "gcpkms/config" {
  capabilities = ["delete"]
}
```

### Sample Request

```text
$ curl \
    --header "X-Vault-Token: ..." \
    --request DELETE \
    https://127.0.0.1:8200/v1/gcpkms/config
```


## Decrypt Ciphertext

This endpoint uses the named encryption key to decrypt the ciphertext string. For symmetric key types, the provided ciphertext must come from a previous invocation of the `/encrypt` endpoint. For asymmetric key types, the provided ciphertext must be from the encrypt operation against the corresponding key version's public key.

| Method   | Path                       | Produces                  |
| :------- | :--------------------------| :------------------------ |
| `POST`   | `gcpkms/decrypt/:key`      | `200 application/json`    |

### Example Policy

```hcl
path "gcpkms/decrypt/my-key" {
  capabilities = ["create", "update"]
}
```

### Parameters

- `key` (`string: ""`) -
Name of the key in Vault to use for decryption. This key must already exist in
Vault and must map back to a Google Cloud KMS key.
This is specified as part of the URL.

- `additional_authenticated_data` (`string: ""`) -
Optional data that was specified during encryption of this payload.

- `ciphertext` (`string: ""`) -
Ciphertext to decrypt as previously returned from an encrypt operation. This
must be base64-encoded ciphertext as previously returned from an encrypt
operation.

- `key_version` (`int: 0`) -
Integer version of the crypto key version to use for decryption. This is
required for asymmetric keys. For symmetric keys, Cloud KMS will choose the
correct version automatically.

### Sample Payload

```json
{
  "ciphertext": "CiQAuMv0..."
}
```

### Sample Request

```text
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    https://127.0.0.1:8200/v1/gcpkms/decrypt/my-key
```

### Sample Response

```json
{
  "data": {
    "plaintext": "hello world"
  }
}
```

## Encrypt Plaintext

This endpoint uses the named encryption key to encrypt arbitrary plaintext
string data. The response will be base64-encoded encrypted ciphertext.

| Method   | Path                       | Produces                  |
| :------- | :--------------------------| :------------------------ |
| `POST`   | `gcpkms/encrypt/:key`      | `200 application/json`    |

### Example Policy

```hcl
path "gcpkms/encrypt/my-key" {
  capabilities = ["create", "update"]
}
```

### Parameters

- `key` (`string: ""`) -
Name of the key in Vault to use for encryption. This key must already exist in
Vault and must map back to a Google Cloud KMS key.
This is specified as part of the URL.

- `additional_authenticated_data` (`string: ""`) -
Optional base64-encoded data that, if specified, must also be provided to
decrypt this payload.

- `key_version` (`int: 0`) -
Integer version of the crypto key version to use for encryption. If unspecified,
this defaults to the latest active crypto key version.

- `plaintext` (`string: ""`) -
Plaintext value to be encrypted. This can be a string or binary, but the size
is limited. See the Google Cloud KMS documentation for information on size
limitations by key types.

### Sample Payload

```json
{
  "plaintext": "hello world"
}
```

### Sample Request

```text
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    https://127.0.0.1:8200/v1/gcpkms/encrypt/my-key
```

### Sample Response

```json
{
  "data": {
    "ciphertext": "CiQAuMv0...",
    "key_version": "1"
  }
}
```

## Re-Encrypt Existing Ciphertext

This endpoint uses the named encryption key to re-encrypt the underlying
cryptokey to the latest version for this ciphertext without disclosing the
original plaintext value to the requestor. This is similar to "rewrapping" in
Vault's transit secrets engine.

| Method   | Path                       | Produces                  |
| :------- | :--------------------------| :------------------------ |
| `POST`   | `gcpkms/reencrypt/:key`    | `200 application/json`    |

### Example Policy

```hcl
path "gcpkms/reencrypt/my-key" {
  capabilities = ["create", "update"]
}
```

### Parameters

- `key` (`string: ""`) -
Name of the key to use for encryption. This key must already exist in Vault and
Google Cloud KMS.
This is specified as part of the URL.

- `additional_authenticated_data` (`string: ""`) -
Optional data that, if specified, must also be provided during decryption.

- `ciphertext` (`string: ""`) -
Ciphertext to be re-encrypted to the latest key version. This must be ciphertext
that Vault previously generated for this named key.

- `key_version` (`int: 0`) -
Integer version of the crypto key version to use for re-encryption. If unspecified,
this defaults to the latest active crypto key version.

### Sample Payload

```json
{
  "ciphertext": "CiQAuMv0..."
}
```

### Sample Request

```text
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    https://127.0.0.1:8200/v1/gcpkms/reencrypt/my-key
```

### Sample Response
```json
{
  "data": {
    "ciphertext": "0lX848IG...",
    "key_version": "3"
  },
}
```

## Sign Digest

This endpoint uses the named encryption key to sign digest string data. The
response will include the base64-encoded signature.

| Method   | Path                       | Produces                  |
| :------- | :--------------------------| :------------------------ |
| `POST`   | `gcpkms/sign/:key`         | `200 application/json`    |

### Example Policy

```hcl
path "gcpkms/sign/my-key" {
  capabilities = ["create", "update"]
}
```

### Parameters

- `key` (`string: ""`) -
Name of the key in Vault to use for signing. This key must already exist in
Vault and must map back to a Google Cloud KMS key.
This is specified as part of the URL.

- `key_version` (`int: <required>`) -
Integer version of the crypto key version to use for signing.

- `digest` (`string: <required>`) -
Digest to sign. This digest is the base64 encoded binary value, and must match
the signing algorithm digest of the Cloud KMS key, for example:

    ```text
    $ openssl dgst -sha256 -binary /my/file | base64
    ```

### Sample Payload

```json
{
  "key_version": "1",
  "digest": "LoM6lxd8YS+hUynZwrlCG20ViUUqqbNNNmh7HCtOkSc="
}
```

### Sample Request

```text
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    https://127.0.0.1:8200/v1/gcpkms/sign/my-key
```

### Sample Response

```json
{
  "data": {
    "signature": "MGYCMQCOfWMc21jBevoRRo4zGjYsCXer8s..."
  }
}
```

## Verify Digest

This endpoint uses the named encryption key to verify a signature and digest
string data.

| Method   | Path                       | Produces                  |
| :------- | :--------------------------| :------------------------ |
| `POST`   | `gcpkms/verify/:key`       | `200 application/json`    |

### Example Policy

```hcl
path "gcpkms/verify/my-key" {
  capabilities = ["create", "update"]
}
```

### Parameters

- `key` (`string: ""`) -
Name of the key in Vault to use for verifying. This key must already exist in
Vault and must map back to a Google Cloud KMS key.
This is specified as part of the URL.

- `key_version` (`int: <required>`) -
Integer version of the crypto key version to use for verifying.

- `digest` (`string: <required>`) -
Digest that was signed. This digest is the base64 encoded binary value, and must match
the signing algorithm digest of the Cloud KMS key. For example:

    ```text
    $ openssl dgst -sha256 -binary /my/file | base64
    ```

- `signature` (`string: <required>`) -
Signature of the digest as returned from a signing operation.

### Sample Payload

```json
{
  "key_version": "1",
  "digest": "LoM6lxd8YS+hUynZwrlCG20ViUUqqbNNNmh7HCtOkSc=",
  "signature": "MGQCMEN2rgg6sj2vUEC3IcKDD+UprtMnxDoB3..."
}
```

### Sample Request

```text
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    https://127.0.0.1:8200/v1/gcpkms/verify/my-key
```

### Sample Response

```json
{
  "data": {
    "valid": true
  }
}
```

## List Keys

This endpoint lists the named keys available for use in Vault. It does not list
all Google Cloud KMS keys.

| Method   | Path                     | Produces                  |
| :------- | :------------------------| :------------------------ |
| `LIST`   | `gcpkms/keys`            | `200 application/json`    |

### Example Policy

```hcl
path "gcpkms/keys" {
  capabilities = ["list"]
}
```

### Sample Request

```text
$ curl \
    --header "X-Vault-Token: ..." \
    --request LIST \
    https://127.0.0.1:8200/v1/gcpkms/keys
```

### Sample Response
```json
{
  "data": {
    "keys": [
      "my-key"
    ]
  }
}
```

## Create/Update Google Cloud KMS key

This endpoint is used to create or update a Google Cloud KMS key. In addition to
registering the key in Vault, this endpoint will also create the corresponding
Google Cloud KMS key with the given configuration options.


| Method   | Path                     | Produces                  |
| :------- | :------------------------| :------------------------ |
| `POST`   | `gcpkms/keys/:key`       | `204 (empty body)`        |

### Example Policy

```hcl
path "gcpkms/keys/my-key" {
  capabilities = ["create", "update"]
}
```

### Parameters

- `key` (`string: ""`) -
Name of the key in Vault.
This is specified as part of the URL.

- `crypto_key` (`string: ""`) -
Name of the crypto key to use. If the given crypto key does not exist, Vault
will try to create it. This defaults to the name of the key given to Vault as
the parameter if unspecified.

- `key_ring` (`string: ""`) -
Full Google Cloud resource ID of the key ring with the project and location
(e.g. projects/my-project/locations/global/keyRings/my-keyring). If the given
key ring does not exist, Vault will try to create it during a create operation.

- `label` (`map<string>string: nil`) -
Arbitrary key=value label to apply to the crypto key. To specify multiple
labels, specify this argument multiple times (e.g. label="a=b" label="c=d").

- `rotation_period` (`string: ""`) -
Amount of time between crypto key version rotations. This is specified as a
time duration value like 72h (72 hours). The smallest possible value is 24h.

### Sample Payload

```json
{
  "key_ring": "projects/my-project/locations/my-location/keyRings/my-keyring",
  "labels": {
    "foo": "bar"
  },
  "rotation_period": "72h",
}
```

### Sample Request

```text
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    https://127.0.0.1:8200/v1/gcpkms/keys/my-key
```

## Delete Google Cloud KMS Key

This endpoint deletes a key from both Vault and Google Cloud KMS. This will
disable all crypto key versions for this crypto key in Google Cloud KMS and
delete Vault's reference to the crypto key.

| Method   | Path                     | Produces                  |
| :------- | :------------------------| :------------------------ |
| `DELETE` | `gcpkms/keys/:key`       | `200 application/json`    |

### Example Policy

```hcl
path "gcpkms/keys/my-key" {
  capabilities = ["delete"]
}
```

### Sample Request

```text
$ curl \
    --header "X-Vault-Token: ..." \
    --request DELETE \
    https://127.0.0.1:8200/v1/gcpkms/keys/my-key
```

## Read Google Cloud KMS Key

This endpoint reads data about a Google Cloud KMS crypto key, including the key
status and current primary key version.

| Method   | Path                     | Produces                  |
| :------- | :------------------------| :------------------------ |
| `GET`    | `gcpkms/keys/:key`       | `200 application/json`    |

### Example Policy

```hcl
path "gcpkms/keys/my-key" {
  capabilities = ["read"]
}
```

### Sample Request

```text
$ curl \
    --header "X-Vault-Token: ..." \
    --request GET \
    https://127.0.0.1:8200/v1/gcpkms/keys/my-key
```

### Sample Response

```json
{
  "data": {
    "id": "projects/my-project/locations/my-location/keyRings/my-keyring/cryptoKeys/my-crypto-key",
    "labels": {
      "foo": "bar"
    },
    "next_rotation_time_seconds": 1536613424,
    "primary_version": "3",
    "purpose": "encrypt_decrypt",
    "rotation_schedule_seconds": 259200,
    "state": "enabled"
  }
}
```

## Read Vault Key Configuration

This endpoint reads data about a Vault's configuration of the key.

| Method   | Path                      | Produces                  |
| :------- | :-------------------------| :------------------------ |
| `GET`    | `gcpkms/keys/config/:key` | `200 application/json`    |

### Example Policy

```hcl
path "gcpkms/keys/config/my-key" {
  capabilities = ["read"]
}
```

### Sample Request

```text
$ curl \
    --header "X-Vault-Token: ..." \
    --request GET \
    https://127.0.0.1:8200/v1/gcpkms/keys/config/my-key
```

### Sample Response

```json
{
  "data": {
    "name": "my-key",
    "crypto_key": "projects/my-project/locations/my-location/keyRings/my-keyring/cryptoKeys/my-crypto-key",
    "min_version": 10
  }
}
```

## Update Vault Key Configuration

This endpoint is used to update Vault's information about an existing key.


| Method   | Path                      | Produces                  |
| :------- | :-------------------------| :------------------------ |
| `POST`   | `gcpkms/keys/config/:key` | `204 (empty body)`        |

### Example Policy

```hcl
path "gcpkms/keys/my-key" {
  capabilities = ["create", "update"]
}
```

### Parameters

- `key` (`string: ""`) -
Name of the key in Vault.
This is specified as part of the URL.

- `min_version` (`int: 0`) -
Minimum allowed crypto key version. If set to a positive value, key versions
less than the given value are not permitted to be used. If set to 0 or a
negative value, there is no minimum key version.

- `max_version` (`int: 0`) -
Maximum allowed crypto key version. If set to a positive value, key versions
greater than the given value are not permitted to be used. If set to 0 or a
negative value, there is no maximum key version.

### Sample Payload

```json
{
  "min_version": 10
}
```

### Sample Request

```text
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    https://127.0.0.1:8200/v1/gcpkms/keys/config/my-key
```

## Deregister Crypto Key

This endpoint deregisters an existing reference Vault has to a crypto key in
Google Cloud KMS. The underlying Google Cloud KMS key remains unchanged.

| Method   | Path                          | Produces                  |
| :------- | :-----------------------------| :------------------------ |
| `POST`   | `gcpkms/keys/deregister/:key` | `204 (empty body)`        |

### Example Policy

```hcl
path "gcpkms/keys/deregister/my-key" {
  capabilities = ["create", "update"]
}
```

### Sample Request

```text
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    https://127.0.0.1:8200/v1/gcpkms/keys/deregister/my-key
```

## Register Crypto Key

This endpoint registers an existing crypto key in Google Cloud KMS and makes it
available for encryption and decryption in Vault.

| Method   | Path                        | Produces                  |
| :------- | :---------------------------| :------------------------ |
| `POST`   | `gcpkms/keys/register/:key` | `204 (empty body)`        |

### Example Policy

```hcl
path "gcpkms/keys/register/my-key" {
  capabilities = ["create", "update"]
}
```

### Parameters

- `key` (`string: ""`) -
Name of the key to register in Vault. This will be the named used to refer to
the underlying crypto key when encrypting or decrypting data.
This is specified as part of the URL.

- `crypto_key` (`string: ""`) -
Full resource ID of the crypto key including the project, location, key ring,
and crypto key like "projects/%s/locations/%s/keyRings/%s/cryptoKeys/%s". This
crypto key must already exist in Google Cloud KMS unless verify is set to
"false".

- `verify` (`bool: true`) -
Verify that the given Google Cloud KMS crypto key exists and is accessible
before creating the storage entry in Vault. Set this to "false" if the key will
not exist at creation time.

### Sample Payload

```json
{
  "crypto_key": "projects/my-project/locations/my-location/keyRings/my-keyring/cryptoKeys/my-crypto-key",
  "verify": true,
}
```

### Sample Request

```text
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    https://127.0.0.1:8200/v1/gcpkms/keys/register/my-key
```

## Rotate Crypto Key

This endpoint rotates a crypto key by creating a new crypto key version for the
corresponding Google Cloud KMS key and updates the new crypto key to be the
primary key for future encryptions.

**It can take up to 2 hours for a new crypto key version to become the primary,
so be sure to issue a read operation if you require new data to be encrypted
with this key.**

| Method   | Path                      | Produces                  |
| :------- | :-------------------------| :------------------------ |
| `POST`   | `gcpkms/keys/rotate/:key` | `200 application/json`    |

### Example Policy

```hcl
path "gcpkms/keys/rotate/my-key" {
  capabilities = ["create", "update"]
}
```

### Sample Request

```text
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    https://127.0.0.1:8200/v1/gcpkms/keys/rotate/my-key
```

### Sample Response

```json
{
  "data": {
    "key_version": "3"
  }
}
```

## Trim KMS Key Versions

This endpoint deletes old crypto key versions that are older than the key's specified `min_version`.

**Data encrypted with older key versions will be irrecoverable!**

| Method   | Path                      | Produces            |
| :------- | :-------------------------| :------------------ |
| `POST`   | `gcpkms/keys/trim/:key`   | `204 (empty body)`  |

### Example Policy

```hcl
path "gcpkms/keys/trim/my-key" {
  capabilities = ["create", "update"]
}
```

### Sample Request

```text
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    https://127.0.0.1:8200/v1/gcpkms/keys/trim/my-key
```
