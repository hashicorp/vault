---
layout: "api"
page_title: "/sys/leases - HTTP API"
sidebar_current: "docs-http-system-leases"
description: |-
  The `/sys/leases` endpoints are used to view and manage leases.
---

# `/sys/leases/lookup`

The `/sys/leases/lookup` endpoint is used to lookup metadata related to 
leases.

## Read Lease

This endpoint retrieve lease metadata.

| Method   | Path                          | Produces               |
| :------- | :---------------------------- | :--------------------- |
| `PUT`    | `/sys/leases/lookup`          | `200 application/json` |

### Parameters

- `lease_id` `(string: <required>)` – Specifies the ID of the lease to lookup.

### Sample Payload

```json
{
  "lease_id": "aws/creds/deploy/abcd-1234..."
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request PUT \
    --data @payload.json \
    https://vault.rocks/v1/sys/leases/lookup
```

### Sample Response

```json
{








}
```

## List Leases

This endpoint returns a list of lease ids.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `LIST`   | `/sys/leases/lookup/:prefix` | `200 application/json` |


### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request LIST \
    https://vault.rocks/v1/sys/leases/lookup/aws/creds/deploy/
```

### Sample Response

```json
{
  "data":{
    "keys":[
      "abcd-1234...",
      "efgh-1234...",
      "ijkl-1234..."
    ]
  }
}
```

# `/sys/leases/renew`

The `/sys/leases/renew` endpoint is used to renew secrets.

## Renew Secret

This endpoint renews a secret, requesting to extend the lease.

| Method   | Path                          | Produces               |
| :------- | :---------------------------- | :--------------------- |
| `PUT`    | `/sys/leases/renew`           | `200 application/json` |

### Parameters

- `lease_id` `(string: <required>)` – Specifies the ID of the lease to extend.
  This can be specified as part of the URL or as part of the request body.

- `increment` `(int: 0)` – Specifies the requested amount of time (in seconds)
  to extend the lease.

### Sample Payload

```json
{
  "lease_id": "aws/creds/deploy/abcd-1234...",
  "increment": 1800
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request PUT \
    --data @payload.json \
    https://vault.rocks/v1/sys/leases/renew
```

### Sample Response

```json
{
  "lease_id": "aws/creds/deploy/abcd-1234...",
  "renewable": true,
  "lease_duration": 2764790
}
```

# `/sys/leases/revoke`

The `/sys/leases/revoke` endpoint is used to revoke secrets.

## Revoke Secret

This endpoint revokes a secret immediately.

| Method   | Path                          | Produces               |
| :------- | :---------------------------- | :--------------------- |
| `PUT`    | `/sys/leases/revoke`          | `204 (empty body)`     |

### Parameters

- `lease_id` `(string: <required>)` – Specifies the ID of the lease to revoke.

### Sample Payload

```json
{
  "lease_id": "postgresql/creds/readonly/abcd-1234..."
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request PUT \
    --data @payload.json \
    https://vault.rocks/v1/sys/leases/revoke
```

# `/sys/leases/revoke-force`

The `/sys/leases/revoke-force` endpoint is used to revoke secrets or tokens 
based on prefix while ignoring backend errors.

## Revoke Force

This endpoint revokes all secrets or tokens generated under a given prefix
immediately. Unlike `/sys/leases/revoke-prefix`, this path ignores backend errors
encountered during revocation. This is _potentially very dangerous_ and should
only be used in specific emergency situations where errors in the backend or the
connected backend service prevent normal revocation.

By ignoring these errors, Vault abdicates responsibility for ensuring that the
issued credentials or secrets are properly revoked and/or cleaned up. Access to
this endpoint should be tightly controlled.

| Method   | Path                                | Produces               |
| :------- | :---------------------------------- | :--------------------- |
| `PUT`    | `/sys/leases/revoke-force/:prefix`  | `204 (empty body)`     |

### Parameters

- `prefix` `(string: <required>)` – Specifies the prefix to revoke. This is
  specified as part of the URL.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request PUT \
    https://vault.rocks/v1/sys/leases/revoke-force/aws/creds
```

# `/sys/leases/revoke-prefix`

The `/sys/leases/revoke-prefix` endpoint is used to revoke secrets or tokens based on
prefix.

## Revoke Prefix

This endpoint revokes all secrets (via a lease ID prefix) or tokens (via the
tokens' path property) generated under a given prefix immediately. This requires
`sudo` capability and access to it should be tightly controlled as it can be
used to revoke very large numbers of secrets/tokens at once.

| Method   | Path                                | Produces               |
| :------- | :---------------------------------- | :--------------------- |
| `PUT`    | `/sys/leases/revoke-prefix/:prefix` | `204 (empty body)`     |

### Parameters

- `prefix` `(string: <required>)` – Specifies the prefix to revoke. This is
  specified as part of the URL.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request PUT \
    https://vault.rocks/v1/sys/leases/revoke-prefix/aws/creds
```
