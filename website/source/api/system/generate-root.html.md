---
layout: "api"
page_title: "/sys/generate-root - HTTP API"
sidebar_current: "docs-http-system-generate-root"
description: |-
  The `/sys/generate-root/` endpoints are used to create a new root key for
  Vault.
---

# `/sys/generate-root`

The `/sys/generate-root` endpoint is used to create a new root key for Vault.

## Read Root Generation Progress

This endpoint reads the configuration and process of the current root generation
attempt.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `GET`    | `/sys/generate-root/attempt` | `200 application/json` |

### Sample Request

```
$ curl \
    https://vault.rocks/v1/sys/generate-root/attempt
```

### Sample Response

```json
{
  "started": true,
  "nonce": "2dbd10f1-8528-6246-09e7-82b25b8aba63",
  "progress": 1,
  "required": 3,
  "encoded_token": "",
  "pgp_fingerprint": "",
  "complete": false
}
```

If a root generation is started, `progress` is how many unseal keys have been
provided for this generation attempt, where `required` must be reached to
complete. The `nonce` for the current attempt and whether the attempt is
complete is also displayed. If a PGP key is being used to encrypt the final root
token, its fingerprint will be returned. Note that if an OTP is being used to
encode the final root token, it will never be returned.

## Start Root Token Generation

This endpoint initializes a new root generation attempt. Only a single root
generation attempt can take place at a time. One (and only one) of `otp` or
`pgp_key` are required.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `PUT`    | `/sys/generate-root/attempt` | `200 application/json` |

### Parameters

- `otp` `(string: <required-unless-pgp>)` – Specifies a base64-encoded 16-byte
  value. The raw bytes of the token will be XOR'd with this value before being
  returned to the final unseal key provider.

- `pgp_key` `(string: <required-unless-otp>)` – Specifies a base64-encoded PGP
  public key. The raw bytes of the token will be encrypted with this value
  before being returned to the final unseal key provider.

### Sample Payload

```json
{
  "otp": "CB23=="
}
```

### Sample Request

```
$ curl \
    --request PUT \
    --data @payload.json \
    https://vault.rocks/v1/sys/generate-root/attempt    
```

### Sample Response

```json
{
  "started": true,
  "nonce": "2dbd10f1-8528-6246-09e7-82b25b8aba63",
  "progress": 1,
  "required": 3,
  "encoded_token": "",
  "pgp_fingerprint": "816938b8a29146fbe245dd29e7cbaf8e011db793",
  "complete": false
}
```

## Cancel Root Generation

This endpoint cancels any in-progress root generation attempt. This clears any
progress made. This must be called to change the OTP or PGP key being used.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `DELETE` | `/sys/generate-root/attempt` | `204 (empty body)`     |

### Sample Request

```
$ curl \
    --request DELETE \
    https://vault.rocks/v1/sys/generate-root/attempt
```

## Provide Key Share to Generate Root

This endpoint is used to enter a single master key share to progress the root
generation attempt. If the threshold number of master key shares is reached,
Vault will complete the root generation and issue the new token.  Otherwise,
this API must be called multiple times until that threshold is met. The attempt
nonce must be provided with each call.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `PUT`    | `/sys/generate-root/update`  | `200 application/json` |

### Parameters

- `key` `(string: <required>)` – Specifies a single master key share.

- `nonce` `(string: <required>)` – Specifies the nonce of the attempt.

### Sample Payload

```json
{
  "key": "acbd1234",
  "nonce": "ad235"
}
```

### Sample Request

```
$ curl \
    --request PUT \
    --data @payload.json \
    https://vault.rocks/v1/sys/generate-root/update
```

### Sample Response

This returns a JSON-encoded object indicating the attempt nonce, and completion
status, and the encoded root token, if the attempt is complete.

```json
{
  "started": true,
  "nonce": "2dbd10f1-8528-6246-09e7-82b25b8aba63",
  "progress": 3,
  "required": 3,
  "pgp_fingerprint": "",
  "complete": true,
  "encoded_token": "FPzkNBvwNDeFh4SmGA8c+w=="
}
```
