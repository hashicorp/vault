---
layout: "api"
page_title: "/sys/health - HTTP API"
sidebar_title: "<code>/sys/health</code>"
sidebar_current: "api-http-system-health"
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
- `472` if data recovery mode replication secondary and active
- `473` if performance standby 
- `501` if not initialized
- `503` if sealed

### Parameters

- `standbyok` `(bool: false)` – Specifies if being a standby should still return
  the active status code instead of the standby status code. This is useful when
  Vault is behind a non-configurable load balance that just wants a 200-level
  response. This will not apply if the node is a performance standby.
  
- `perfstandbyok` `(bool: false)` – Specifies if being a performance standby should
  still return the active status code instead of the performance standby status code.
  This is useful when Vault is behind a non-configurable load balance that just wants
  a 200-level response.

- `activecode` `(int: 200)` – Specifies the status code that should be returned
  for an active node.

- `standbycode` `(int: 429)` – Specifies the status code that should be returned
  for a standby node.

- `drsecondarycode` `(int: 472)` – Specifies the status code that should be
  returned for a DR secondary node.

- `performancestandbycode` `(int: 473)` – Specifies the status code that should be
  returned for a performance standby node.

- `sealedcode` `(int: 503)` – Specifies the status code that should be returned
  for a sealed node.

- `uninitcode` `(int: 501)` – Specifies the status code that should be returned
  for a uninitialized node.

### Sample Request

```
$ curl \
    http://127.0.0.1:8200/v1/sys/health
```

### Sample Response

This response is only returned for a `GET` request.

Note: `replication_perf_mode` and `replication_dr_mode` reflect the state of
the active node in the cluster; if you are querying it for a standby that has
just come up, it can take a small time for the active node to inform the
standby of its status.

```json
{
  "initialized": true,
  "sealed": false,
  "standby": false,
  "performance_standby": false,
  "replication_perf_mode": "disabled",
  "replication_dr_mode": "disabled",
  "server_time_utc": 1516639589,
  "version": "0.9.1",
  "cluster_name": "vault-cluster-3bd69ca2",
  "cluster_id": "00af5aa8-c87d-b5fc-e82e-97cd8dfaf731"
}
```
