---
layout: "api"
page_title: "/sys/replication - HTTP API"
sidebar_current: "docs-http-system-replication-dr"
description: |-
  The '/sys/replication/dr' endpoint focuses on managing general operations in Vault Enterprise Disaster Recovery replication
---

# `/sys/replication/dr`

~> **Enterprise Only** – These endpoints require Vault Enterprise.

## Check DR Status

This endpoint prints information about the status of replication (mode,
sync progress, etc).

This is an authenticated endpoint.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `GET`    | `/sys/replication/dr/status`    | `200 application/json` |

### Sample Request

```
$ curl \
    https://vault.rocks/v1/sys/replication/dr/status
```

### Sample Response

The printed status of the replication environment. As an example, for a
primary, it will look something like:

```json
{
  "mode": "dr-primary",
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

## Enable DR Primary Replication

This endpoint enables DR replication in primary mode. This is used when DR replication
is currently disabled on the cluster (if the cluster is already a secondary, it
must be promoted).

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `POST`   | `/sys/replication/dr/primary/enable` | `204 (empty body)` |

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
    https://vault.rocks/v1/sys/replication/dr/primary/enable
```

## Demote DR Primary

This endpoint demotes a DR primary cluster to a secondary. This DR secondary cluster
will not attempt to connect to a primary (see the update-primary call), but will
maintain knowledge of its cluster ID and can be reconnected to the same
DR replication set without wiping local storage.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `POST`   | `/sys/replication/dr/primary/demote` | `204 (empty body)` |

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    https://vault.rocks/v1/sys/replication/dr/primary/demote
```

## Disable DR Primary

This endpoint disables DR replication entirely on the cluster. Any secondaries will
no longer be able to connect. Caution: re-enabling this node as a primary or
secondary will change its cluster ID; in the secondary case this means a wipe of
the underlying storage when connected to a primary, and in the primary case,
secondaries connecting back to the cluster (even if they have connected before)
will require a wipe of the underlying storage.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `POST`   | `/sys/replication/dr/primary/disable` | `204 (empty body)` |

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    https://vault.rocks/v1/sys/replication/dr/primary/disable
```

## Generate DR Secondary Token

This endpoint generates a DR secondary activation token for the
cluster with the given opaque identifier, which must be unique. This
identifier can later be used to revoke a DR secondary's access.

**This endpoint requires 'sudo' capability.**

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `GET`    | `/sys/replication/dr/primary/secondary-token` | `200 application/json` |

### Parameters

- `id` `(string: <required>)` – Specifies an opaque identifier, e.g. 'us-east'

- `ttl` `(string: "30m")` – Specifies the TTL for the secondary activation
  token.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    https://vault.rocks/v1/sys/replication/dr/primary/secondary-token?id=us-east-1
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

## Revoke DR Secondary Token

This endpoint revokes a DR secondary's ability to connect to the DR primary cluster;
the DR secondary will immediately be disconnected and will not be allowed to
connect again unless given a new activation token.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `POST`   | `/sys/replication/dr/primary/revoke-secondary` | `204 (empty body)` |

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
    https://vault.rocks/v1/sys/replication/dr/primary/revoke-secondary
```

## Enable DR Secondary

This endpoint enables replication on a DR secondary using a DR secondary activation
token.

!> This will immediately clear all data in the secondary cluster!

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `POST`   | `/sys/replication/dr/secondary/enable` | `204 (empty body)` |

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
    https://vault.rocks/v1/sys/replication/dr/secondary/enable
```

## Promote DR Secondary

This endpoint promotes the DR secondary cluster to DR primary. For data safety and
security reasons, new secondary tokens will need to be issued to other
secondaries, and there should never be more than one primary at a time.

If the DR secondary's primary cluster is also in a performace replication set,
the DR secondary will be promoted into that replication set. Care should be
taken when promoting to ensure multiple performance primary clusters are not
activate at the same time. 

If the DR secondary's primary cluster is a performance secondary, the promoted
cluster will attempt to connect to the performance primary cluster using the
same secondary token.

!> Only one performance primary should be active at a given time. Multiple primaries may
result in data loss!

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `POST`   | `/sys/replication/dr/secondary/promote` | `200 application/json` |

### Parameters

- `key` `(string "")` - Specifies a single master key share. This is required unless reset is true.
- `reset` `(bool false) - Specifies if previously-provided unseal keys are discarded and the promote process is reset.
- `primary_cluster_addr` `(string: "")` – Specifies the cluster address that the
  primary gives to secondary nodes. Useful if the primary's cluster address is
  not directly accessible and must be accessed via an alternate path/address
  (e.g. through a load balancer).

### Sample Payload

```json
{
  "key": "ijH8tphEHaBtgx+IvPfxDsSi2LV4j9k+Lad6eqT5cJw="
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    https://vault.rocks/v1/sys/replication/dr/secondary/promote
```

### Sample Response

```json
{
  "progress": 0,
  "required": 1,
  "complete": false,
  "request_id": "ad8f9074-0e24-d30e-83cd-595c9652ff89",
  "lease_id": "",
  "renewable": false,
  "lease_duration": 0,
  "data": {
    "complete": false,
    "progress": 0,
    "required": 1
  },
  "wrap_info": null,
  "warnings": null,
  "auth": null
}
```
