---
layout: "api"
page_title: "/sys/replication - HTTP API"
sidebar_current: "docs-http-system-replication"
description: |-
  The '/sys/replication' endpoint focuses on managing general operations in Vault Enterprise replication
---

# `/sys/replication`

~> **Enterprise Only** – These endpoints require Vault Enterprise.

## Attempt Recovery

This endpoint attempts recovery if replication is in an adverse state. For
example: an error has caused replication to stop syncing.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `POST`   | `/sys/replication/recover`   | `200 application/json` |

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    http://127.0.0.1:8200/v1/sys/replication/recover
```

### Sample Response

```json
{
  "warnings": ["..."]
}
```

## Reindex Replication

This endpoint reindexes the local data storage. This can cause a very long delay
depending on the number and size of objects in the data store.

**This endpoint requires 'sudo' capability.**

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `POST`   | `/sys/replication/reindex`   | `200 application/json` |

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    http://127.0.0.1:8200/v1/sys/replication/reindex
```

### Sample Response

```json
{
  "warnings": ["..."]
}
```

## Check Status

This endpoint print information about the status of replication (mode,
sync progress, etc).

This is an authenticated endpoint.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `GET`    | `/sys/replication/status`    | `200 application/json` |

### Sample Request

```
$ curl \
    http://127.0.0.1:8200/v1/sys/replication/status
```

### Sample Response

The printed status of the replication environment. As an example, for a
performance primary and DR primary node, it will look something like:

```json
{
  "data": {
    "dr": {
      "cluster_id": "f2c21cb5-523f-617b-20ac-c913d9154ba6",
      "known_secondaries": [
        "3"
      ],
      "last_wal": 291,
      "merkle_root": "38543b95d44132138003939addbaf94125ec184e",
      "mode": "primary",
      "primary_cluster_addr": ""
    },
    "performance": {
      "cluster_id": "1598d434-dfec-1f48-f019-3d22a8075bf9",
      "known_secondaries": [
        "2"
      ],
      "last_wal": 291,
      "merkle_root": "43f40fc775b40cc76cd5d7e289b2e6eaf4ba138c",
      "mode": "primary",
      "primary_cluster_addr": ""
    }
  },
}
```

### Sample Response from Performance Secondary & DR Primary

The printed status of the replication environment. As an example, for a
performnace secondary and DR primary node, it will look something like:

```json
{
  "data": {
    "dr": {
      "cluster_id": "e4bfa800-002e-7b6d-14c2-617855ece02f",
      "known_secondaries": [
        "4"
      ],
      "last_wal": 455,
      "merkle_root": "cdcf796619240ce19dd8af30fa700f64c8006e3d",
      "mode": "primary",
      "primary_cluster_addr": ""
    },
    "performance": {
      "cluster_id": "1598d434-dfec-1f48-f019-3d22a8075bf9",
      "known_primary_cluster_addrs": [
        "https://127.0.0.1:8201"
      ],
      "last_remote_wal": 291,
      "merkle_root": "43f40fc775b40cc76cd5d7e289b2e6eaf4ba138c",
      "mode": "secondary",
      "primary_cluster_addr": "https://127.0.0.1:8201",
      "secondary_id": "2",
      "state": "stream-wals"
    }
  },
}
```
