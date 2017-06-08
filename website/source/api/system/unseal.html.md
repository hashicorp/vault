---
layout: "api"
page_title: "/sys/unseal - HTTP API"
sidebar_current: "docs-http-system-unseal"
description: |-
  The `/sys/unseal` endpoint is used to unseal the Vault.
---

# `/sys/unseal`

The `/sys/seal-unseal` endpoint is used to unseal the Vault.

## Submit Unseal Key

This endpoint is used to enter a single master key share to progress the
unsealing of the Vault. If the threshold number of master key shares is reached,
Vault will attempt to unseal the Vault. Otherwise, this API must be called
multiple times until that threshold is met.

Either the `key` or `reset` parameter must be provided; if both are provided,
`reset` takes precedence.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `PUT`    | `/sys/unseal`                | `200 application/json` |

### Parameters

- `key` `(string: "")` – Specifies a single master key share. This is required
  unless `reset` is true.

- `reset` `(bool: false)` – Specifies if previously-provided unseal keys are
  discarded and the unseal process is reset.

### Sample Payload

```json
{
  "key": "abcd1234..."
}
```

### Sample Request

```
$ curl \
    --request PUT \
    --data @payload.json \
    https://vault.rocks/v1/sys/unseal
```

### Sample Response

The "t" parameter is the threshold, and "n" is the number of shares.

```json
{
  "sealed": true,
  "t": 3,
  "n": 5,
  "progress": 2,
  "version": "0.6.2"
}
```

Sample response when Vault is unsealed.

```json
{
  "sealed": false,
  "t": 3,
  "n": 5,
  "progress": 0,
  "version": "0.6.2",
  "cluster_name": "vault-cluster-d6ec3c7f",
  "cluster_id": "3e8b3fec-3749-e056-ba41-b62a63b997e8"
}
```
