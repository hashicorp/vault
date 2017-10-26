---
layout: "api"
page_title: "/sys/tools - HTTP API"
sidebar_current: "docs-http-system-tools"
description: |-
  This is the API documentation for a general set of crypto  tools.
---

# `/sys/tools`

The `/sys/tools` endpoints are a general set of tools.

## Generate Random Bytes

This endpoint returns high-quality random bytes of the specified length.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `POST`   | `/sys/tools/random(/:bytes)`   | `200 application/json` |

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
    https://vault.rocks/v1/sys/tools/random/164
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
| `POST`   | `/sys/tools/hash(/:algorithm)` | `200 application/json` |

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
    https://vault.rocks/v1/sys/tools/hash/sha2-512
```

### Sample Response

```json
{
  "data": {
    "sum": "dGhlIHF1aWNrIGJyb3duIGZveAo="
  }
}
```

