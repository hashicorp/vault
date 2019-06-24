---
layout: "api"
page_title: "/sys/storage/raft - HTTP API"
sidebar_title: "<code>/sys/storage/raft</code>"
sidebar_current: "api-http-system-storage-raft"
description: |-

  The `/sys/storage/raft` endpoint is used to check the high availability status and
  current leader of Vault.
---

## Join a Raft cluster

This endpoint joins a node to the Raft cluster.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `POST`   | `/sys/storage/raft/join`    |

### Parameters

- `leader_api_addr` `(string: <required>)` – Address of the leader node in the
  Raft cluster to which this node is trying to join.

- `retry` `(bool: false)` - Retry joining the Raft cluster in case of
  failures.

- `leader_ca_cert` `(string: "")` - CA certificate used to communicate with
  Raft's leader node.

- `leader_client_cert` `(string: "")` - Client certificate used to communicate
  with Raft's leader node.

- `leader_client_key` `(string: "")` - Client key used to communicate with
  Raft's leader node.

### Sample Payload
```json
{
  "leader_api_addr": "https://127.0.0.1:8200",
  "leader_ca_cert": "<pem encoded ca cert>",
  "leader_client_cert": "<pem encoded client cert>",
  "leader_client_key": "<pem encoded client key>"
}
```
### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    http://127.0.0.1:8200/v1/sys/storage/raft/join
```

## Read Raft Configuration

This endpoint returns the details of all the nodes in the raft cluster.

| Method                       | Path                           |
| :--------------------------- | :----------------------------  |
| `GET`                          | `/sys/storage/raft/configuration`  |

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    http://127.0.0.1:8200/v1/sys/storage/raft/configuration
```

## Remove a node from Raft cluster

This endpoint removes a node from the raft cluster.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `POST`   | `/sys/storage/raft/remove-peer`    |

### Sample Payload

```json
{
  "server_id": "raft_node_1"
}
```
### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    http://127.0.0.1:8200/v1/sys/storage/raft/remove-peer
```

## Take a snapshot of the Raft cluster

This endpoint returns a snapshot of the current state of the raft cluster.

| Method                       | Path                           |
| :--------------------------- | :----------------------------  |
| `GET`                          | `/sys/storage/raft/snapshot`  |

## Restore Raft using a snapshot

Installs the provided snapshot, returning the cluster to the state defined in it.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `POST`   | `/sys/storage/raft/snapshot`    |

The snapshot should be set as the request's body.

