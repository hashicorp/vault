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
| `POST`    | `/sys/replication/dr/primary/secondary-token` | `200 application/json` |

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

If the DR secondary's primary cluster is also in a performance replication set,
the DR secondary will be promoted into that replication set. Care should be
taken when promoting to ensure multiple performance primary clusters are not
activate at the same time. 

If the DR secondary's primary cluster is a performance secondary, the promoted
cluster will attempt to connect to the performance primary cluster using the
same secondary token.

This endpoint requires a DR Operation Token to be provided as means of
authorization. See the [DR Operation Token API
docs](/api/system/replication-dr.html#sys-generate-dr-operation-token) for more information.

!> Only one performance primary should be active at a given time. Multiple primaries may
result in data loss!

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `POST`   | `/sys/replication/dr/secondary/promote` | `200 application/json` |

### Parameters

- `dr_operation_token` `(string: <required>)` - DR operation token used to authorize this request.
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

## Update DR Secondary's Primary

This endpoint changes a DR secondary cluster's assigned primary cluster using a
secondary activation token. This does not wipe all data in the cluster.

This endpoint requires a DR Operation Token to be provided as means of
authorization. See the [DR Operation Token API
docs](/api/system/replication-dr.html#sys-generate-dr-operation-token) for more information.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `POST`   | `/sys/replication/dr/secondary/update-primary` | `204 (empty body)` |

### Parameters

- `dr_operation_token` `(string: <required>)` - DR operation token used to authorize this request.

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
  "dr_operation_token": "...",
  "token": "..."
}
``` 

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    https://vault.rocks/v1/sys/replication/dr/secondary/update-primary
```

# `/sys/replication/dr/secondary/generate-operation-token`

The `/sys/replication/dr/secondary/generate-operation-token` endpoint is used to create a new Disaster
Recovery operation token for a DR secondary. These tokens are used to authorize
certain DR Operation. They should be treated like traditional root tokens by
being generated with needed and deleted soon after.

## Read Generation Progress

This endpoint reads the configuration and process of the current generation
attempt.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `GET`    | `/sys/replication/dr/secondary/generate-operation-token/attempt` | `200 application/json` |

### Sample Request

```
$ curl \
    https://vault.rocks/v1/sys/replication/dr/secondary/generate-operation-token/attempt
```

### Sample Response

```json
{
  "started": true,
  "nonce": "2dbd10f1-8528-6246-09e7-82b25b8aba63",
  "progress": 1,
  "required": 3,
  "encoded_token": "",
  "pgp_fingerprint": "",
  "complete": false
}
```

If a generation is started, `progress` is how many unseal keys have been
provided for this generation attempt, where `required` must be reached to
complete. The `nonce` for the current attempt and whether the attempt is
complete is also displayed. If a PGP key is being used to encrypt the final
token, its fingerprint will be returned. Note that if an OTP is being used to
encode the final token, it will never be returned.

## Start Token Generation

This endpoint initializes a new generation attempt. Only a single
generation attempt can take place at a time. One (and only one) of `otp` or
`pgp_key` are required.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `PUT`    | `/sys/replication/dr/secondary/generate-operation-token/attempt` | `200 application/json` |

### Parameters

- `otp` `(string: <required-unless-pgp>)` – Specifies a base64-encoded 16-byte
  value. The raw bytes of the token will be XOR'd with this value before being
  returned to the final unseal key provider.

- `pgp_key` `(string: <required-unless-otp>)` – Specifies a base64-encoded PGP
  public key. The raw bytes of the token will be encrypted with this value
  before being returned to the final unseal key provider.

### Sample Payload

```json
{
  "otp": "CB23=="
}
```

### Sample Request

```
$ curl \
    --request PUT \
    --data @payload.json \
    https://vault.rocks/v1/sys/replication/dr/secondary/generate-operation-token/attempt
```

### Sample Response

```json
{
  "started": true,
  "nonce": "2dbd10f1-8528-6246-09e7-82b25b8aba63",
  "progress": 1,
  "required": 3,
  "encoded_token": "",
  "pgp_fingerprint": "816938b8a29146fbe245dd29e7cbaf8e011db793",
  "complete": false
}
```

## Cancel Generation

This endpoint cancels any in-progress generation attempt. This clears any
progress made. This must be called to change the OTP or PGP key being used.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `DELETE` | `/sys/replication/dr/secondary/generate-operation-token/attempt` | `204 (empty body)`     |

### Sample Request

```
$ curl \
    --request DELETE \
    https://vault.rocks/v1/sys/replication/dr/secondary/generate-operation-token/attempt
```

## Provide Key Share to Generate Token

This endpoint is used to enter a single master key share to progress the
generation attempt. If the threshold number of master key shares is reached,
Vault will complete the generation and issue the new token.  Otherwise,
this API must be called multiple times until that threshold is met. The attempt
nonce must be provided with each call.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `PUT`    | `/sys/replication/dr/secondary/generate-operation-token/update`  | `200 application/json` |

### Parameters

- `key` `(string: <required>)` – Specifies a single master key share.

- `nonce` `(string: <required>)` – Specifies the nonce of the attempt.

### Sample Payload

```json
{
  "key": "acbd1234",
  "nonce": "ad235"
}
```

### Sample Request

```
$ curl \
    --request PUT \
    --data @payload.json \
    https://vault.rocks/v1/sys/replication/dr/secondary/generate-operation-token/update
```

### Sample Response

This returns a JSON-encoded object indicating the attempt nonce, and completion
status, and the encoded token, if the attempt is complete.

```json
{
  "started": true,
  "nonce": "2dbd10f1-8528-6246-09e7-82b25b8aba63",
  "progress": 3,
  "required": 3,
  "pgp_fingerprint": "",
  "complete": true,
  "encoded_token": "FPzkNBvwNDeFh4SmGA8c+w=="
}
```


## Delete DR Operation Token
 
This endpoint revokes the DR Operation Token. This token does not have a TTL
and therefore should be deleted when it is no longer needed.


| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `POST`   | `/sys/replication/dr/secondary/operation-token/delete` | `204 (empty body)` |

### Parameters

- `dr_operation_token` `(string: <required>)` - DR operation token used to authorize this request.

### Sample Payload

```json
{
  "dr_operation_token": "..."
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    https://vault.rocks/v1/sys/replication/dr/secondary/operation-token/delete
```
