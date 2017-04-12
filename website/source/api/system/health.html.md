---
layout: "api"
page_title: "/sys/health - HTTP API"
sidebar_current: "docs-http-system-health"
description: |-
  The `/sys/health` endpoint is used to check the health status of Vault.
---

# `/sys/health`

The `/sys/health` endpoint is used to check the health status of Vault.

## Read Health Information

This endpoint returns the health status of Vault. This matches the semantics of
a Consul HTTP health check and provides a simple way to monitor the health of a
Vault instance.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `HEAD`   | `/sys/health`                | `000 (empty body)`     |
| `GET`    | `/sys/health`                | `000 application/json` |

The default status codes are:

- `200` if initialized, unsealed, and active
- `429` if unsealed and standby
- `501` if not initialized
- `503` if sealed

### Parameters

- `standbyok` `(bool: false)` – Specifies if being a standby should still return
  the active status code instead of the standby status code. This is useful when
  Vault is behind a non-configurable load balance that just wants a 200-level
  response.

- `activecode` `(int: 200)` – Specifies the status code that should be returned
  for an active node.

- `standbycode` `(int: 429)` – Specifies the status code that should be returned
  for a standby node.

- `sealedcode` `(int: 503)` – Specifies the status code that should be returned
  for a sealed node.

- `uninitcode` `(int: 501)` – Specifies the status code that should be returned
  for a uninitialized node.

### Sample Request

```
$ curl \
    https://vault.rocks/v1/sys/health
```

### Sample Response

This response is only returned for a `GET` request.

```json
{
  "cluster_id": "c9abceea-4f46-4dab-a688-5ce55f89e228",
  "cluster_name": "vault-cluster-5515c810",
  "version": "0.6.2",
  "server_time_utc": 1469555798,
  "standby": false,
  "sealed": false,
  "initialized": true
}
```
