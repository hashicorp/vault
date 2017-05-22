---
layout: "api"
page_title: "/sys/replication - HTTP API"
sidebar_current: "docs-http-system-replication"
description: |-
  The '/sys/replication' endpoint focuses on managing general operations in Vault Enterprise replication sets
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
    https://vault.rocks/v1/sys/replication/recover
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
    https://vault.rocks/v1/sys/replication/reindex
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
    https://vault.rocks/v1/sys/replication/status
```

### Sample Response

The printed status of the replication environment. As an example, for a
primary, it will look something like:

```json
{
  "mode": "primary",
  "cluster_id": "d4095d41-3aee-8791-c421-9bc7f88f7c3e",
  "known_secondaries": [],
  "last_wal": 0,
  "merkle_root": "c3260c4c682ff2d6eb3c8bfd877134b3cec022d1",
  "request_id": "009ea98c-06cd-6dc3-74f2-c4904b22e535",
  "lease_id": "",
  "renewable": false,
  "lease_duration": 0,
  "data": {
    "cluster_id": "d4095d41-3aee-8791-c421-9bc7f88f7c3e",
    "known_secondaries": [],
    "last_wal": 0,
    "merkle_root": "c3260c4c682ff2d6eb3c8bfd877134b3cec022d1",
    "mode": "primary"
  },
  "wrap_info": null,
  "warnings": null,
  "auth": null
}
```

## Enable Primary Replication

This endpoint enables replication in primary mode. This is used when replication
is currently disabled on the cluster (if the cluster is already a secondary, it
must be promoted).

!> Only one primary should be active at a given time. Multiple primaries may
result in data loss!

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `POST`   | `/sys/replication/primary/enable` | `204 (empty body)` |

### Parameters

- `primary_cluster_addr` `(string: "")` – Specifies the cluster address that the
  primary gives to secondary nodes. Useful if the primary's cluster address is
  not directly accessible and must be accessed via an alternate path/address,
  such as through a TCP-based load balancer.

### Sample Payload

```json
{}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    https://vault.rocks/v1/sys/replication/primary/enable
```

## Demote Primary

This endpoint demotes a primary cluster to a secondary. This secondary cluster
will not attempt to connect to a primary (see the update-primary call), but will
maintain knowledge of its cluster ID and can be reconnected to the same
replication set without wiping local storage.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `POST`   | `/sys/replication/primary/demote` | `204 (empty body)` |

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    https://vault.rocks/v1/sys/replication/primary/demote
```

## Disable Primary

This endptoin disables replication entirely on the cluster. Any secondaries will
no longer be able to connect. Caution: re-enabling this node as a primary or
secondary will change its cluster ID; in the secondary case this means a wipe of
the underlying storage when connected to a primary, and in the primary case,
secondaries connecting back to the cluster (even if they have connected before)
will require a wipe of the underlying storage.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `POST`   | `/sys/replication/primary/disable` | `204 (empty body)` |

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    https://vault.rocks/v1/sys/replication/primary/disable
```

## Generate Secondary Token

This endpoint generates a secondary activation token for the
cluster with the given opaque identifier, which must be unique. This
identifier can later be used to revoke a secondary's access.

**This endpoint requires 'sudo' capability.**

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `GET`    | `/sys/replication/primary/secondary-token` | `200 application/json` |

### Parameters

- `id` `(string: <required>)` – Specifies an opaque identifier, e.g. 'us-east'

- `ttl` `(string: "30m")` – Specifies the TTL for the secondary activation
  token.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    https://vault.rocks/v1/sys/replication/primary/secondary-token?id=us-east-1
```

### Sample Response

```json
{
  "request_id": "",
  "lease_id": "",
  "lease_duration": 0,
  "renewable": false,
  "data": null,
  "warnings": null,
  "wrap_info": {
    "token": "fb79b9d3-d94e-9eb6-4919-c559311133d6",
    "ttl": 300,
    "creation_time": "2016-09-28T14:41:00.56961496-04:00",
    "wrapped_accessor": ""
  }
}
```

## Revoke Secondary Token

This endpoint revokes a secondary's ability to connect to the primary cluster;
the secondary will immediately be disconnected and will not be allowed to
connect again unless given a new activation token.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `POST`   | `/sys/replication/primary/revoke-secondary` | `204 (empty body)` |

### Parameters

- `id` `(string: <required>)` – Specifies an opaque identifier, e.g. 'us-east'

### Sample Payload

```json
{
  "id": "us-east"
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    https://vault.rocks/v1/sys/replication/primary/revoke-secondary
```

## Enable Secondary

This endpoint enables replication on a secondary using a secondary activation
token.

!> This will immediately clear all data in the secondary cluster!

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `POST`   | `/sys/replication/secondary/enable` | `204 (empty body)` |

### Parameters

- `token` `(string: <required>)` – Specifies the secondary activation token fetched from the primary.

- `primary_api_addr` `(string: "")` – Set this to the API address (normal Vault
  address) to override the value embedded in the token. This can be useful if
  the primary's redirect address is not accessible directly from this cluster
  (e.g. through a load balancer).

- `ca_file` `(string: "")` – Specifies the path to a CA root file (PEM format)
  that the secondary can use when unwrapping the token from the primary. If this
  and ca_path are not given, defaults to system CA roots.

- `ca_path` `(string: "")` – Specifies  the path to a CA root directory
  containing PEM-format files that the secondary can use when unwrapping the
  token from the primary. If this and ca_file are not given, defaults to system
  CA roots.

### Sample Payload

```json
{
  "token": "..."
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    https://vault.rocks/v1/sys/replication/secondary/enable
```

## Promote Secondary

This endpoint promotes the secondary cluster to primary. For data safety and
security reasons, new secondary tokens will need to be issued to other
secondaries, and there should never be more than one primary at a time.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `POST`   | `/sys/replication/secondary/promote` | `204 (empty body)` |

### Parameters

- `primary_cluster_addr` `(string: "")` – Specifies the cluster address that the
  primary gives to secondary nodes. Useful if the primary's cluster address is
  not directly accessible and must be accessed via an alternate path/address
  (e.g. through a load balancer).

### Sample Payload

```json
{}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    https://vault.rocks/v1/sys/replication/secondary/promote
```

## Disable Secondary

This endpoint disables replication entirely on the cluster. The cluster will no
longer be able to connect to the primary.

!> Re-enabling this node as a primary or secondary will change its cluster ID;
in the secondary case this means a wipe of the underlying storage when connected
to a primary, and in the primary case, secondaries connecting back to the
cluster (even if they have connected before) will require a wipe of the
underlying storage.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `POST`   | `/sys/replication/secondary/disable` | `204 (empty body)` |

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    https://vault.rocks/v1/sys/replication/secondary/disable
```

## Update Secondary's Primary

This endpoint changes a secondary cluster's assigned primary  cluster using a
secondary activation token.  This does not wipe all data in the cluster.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `POST`   | `/sys/replication/secondary/update-primary` | `204 (empty body)` |

### Parameters

- `token` `(string: <required>)` – Specifies the secondary activation token
  fetched from the primary. If you set this to a blank string, the cluster will
  stay a secondary but clear its knowledge of any past primary (and thus not
  attempt to connect to the previous primary). This can be useful if the primary
  is down to stop the secondary from trying to reconnect to it.

- `primary_api_addr` `(string: )` – Specifies the API address (normal Vault
  address) to override the value embedded in the token. This can be useful if
  the primary's redirect address is not accessible directly from this cluster.

- `ca_file` `(string: "")` – Specifies the path to a CA root file (PEM format)
  that the secondary can use when unwrapping the token from the primary. If this
  and ca_path are not given, defaults to system CA roots.

- `ca_path` `string: ()` – Specifies the path to a CA root directory containing
  PEM-format files that the secondary can use when unwrapping the token from the
  primary. If this and ca_file are not given, defaults to system CA roots.

### Sample Payload

```json
{
  "token": "..."
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    https://vault.rocks/v1/sys/replication/secondary/update-primary
```
