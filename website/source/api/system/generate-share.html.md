---
layout: "api"
page_title: "/sys/generate-share - HTTP API"
sidebar_current: "docs-http-system-generate-share"
description: |-
  The `/sys/generate-share/` endpoints are used to create new master key 
  shares for Vault.
---

# `/sys/generate-share`

The `/sys/generate-share` endpoint is used to create new master key shares
for Vault.

## Read Master Key Share Generation Progress

This endpoint reads the configuration and progress of the current share generation
attempt.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `GET`    | `/sys/generate-share/attempt` | `200 application/json` |

### Sample Request

```
$ curl \
    https://vault.rocks/v1/sys/generate-share/attempt
```

### Sample Response

```json
{
  "started": true,
  "progress": 1,
  "required": 3,
  "key": "",
  "pgp_fingerprint": "",
  "complete": false
}
```

If a share generation is started, `progress` is how many unseal keys have been
provided for this generation attempt, where `required` must be reached to
complete. Whether the attempt is complete is also displayed. The fingerprint 
of the PGP key being used to encrypt the new master key share is also returned.

## Start Master Key Share Generation

This endpoint initializes a new share generation attempt. Only a single share
generation attempt can take place at a time. `pgp_key` is required.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `PUT`    | `/sys/generate-share/attempt` | `200 application/json` |

### Parameters

- `pgp_key` `(string: <required>)` – Specifies a base64-encoded PGP
  public key. The raw bytes of the token will be encrypted with this value
  before being returned to the final unseal key provider.

### Sample Payload

```json
{
  "pgp-key": "-----BEGIN PGP..."
}
```

### Sample Request

```
$ curl \
    --request PUT \
    --data payload.json \
    https://vault.rocks/v1/sys/generate-share/attempt    
```

### Sample Response

```json
{
  "started": true,
  "progress": 1,
  "required": 3,
  "key": "",
  "pgp_fingerprint": "816938b8a29146fbe245dd29e7cbaf8e011db793",
  "complete": false
}
```

## Cancel Master Key Share Generation

This endpoint cancels any in-progress share generation attempt. This clears any
progress made. This must be called to change the PGP key being used.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `DELETE` | `/sys/generate-share/attempt` | `204 (empty body)`     |

### Sample Request

```
$ curl \
    --request DELETE \
    https://vault.rocks/v1/sys/generate-share/attempt
```

## Provide Key Share to Generate New Master Key Share

This endpoint is used to enter a single master key share to progress the share
generation attempt. If the threshold number of master key shares is reached,
Vault will complete the share generation and issue the new master key share. 
Otherwise, this API must be called multiple times until that threshold is met. 

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `PUT`    | `/sys/generate-share/update`  | `200 application/json` |

### Parameters

- `key` `(string: <required>)` – Specifies a single master key share.

### Sample Payload

```json
{
  "key": "acbd1234"
}
```

### Sample Request

```
$ curl \
    --request PUT \
    --data payload.json \
    https://vault.rocks/v1/sys/generate-share/update
```

### Sample Response

This returns a JSON-encoded object indicating the completion status 
and the encrypted master key share, if the attempt is complete.

```json
{
  "started": true,
  "progress": 3,
  "required": 3,
  "pgp_fingerprint": "",
  "complete": true,
  "key": "wcFMA72ihosZG..."
}
```
