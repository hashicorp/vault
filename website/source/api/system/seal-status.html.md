---
layout: "api"
page_title: "/sys/seal-status - HTTP API"
sidebar_current: "docs-http-system-seal-status"
description: |-
  The `/sys/seal-status` endpoint is used to check the seal status of a Vault.
---

# `/sys/seal-status`

The `/sys/seal-status` endpoint is used to check the seal status of a Vault.

## Seal Status

This endpoint returns the seal status of the Vault. This is an unauthenticated
endpoint.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `GET`    | `/sys/seal-status`           | `200 application/json` |

### Sample Request

```
$ curl \
    https://vault.rocks/v1/sys/seal-status
```

### Sample Response

The "t" parameter is the threshold, and "n" is the number of shares.

```json
{
  "type": "shamir",
  "sealed": true,
  "t": 3,
  "n": 5,
  "progress": 2,
  "nonce": "",
  "version": "0.9.0"
}
```

Sample response when Vault is unsealed.

```json
{
  "type": "shamir",
  "sealed": false,
  "t": 3,
  "n": 5,
  "progress": 0,
  "version": "0.9.0",
  "cluster_name": "vault-cluster-d6ec3c7f",
  "cluster_id": "3e8b3fec-3749-e056-ba41-b62a63b997e8",
  "nonce": "ef05d55d-4d2c-c594-a5e8-55bc88604c24"
}
```
