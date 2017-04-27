---
layout: "api"
page_title: "/sys/lease - HTTP API"
sidebar_current: "docs-http-system-lease"
description: |-
  The `/sys/lease` endpoint is used to view and manage leases.
---

# `/sys/lease/renew`

The `/sys/lease/renew` endpoint is used to renew secrets.

## Renew Secret

This endpoint renews a secret, requesting to extend the lease.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `PUT`    | `/sys/lease/renew`           | `200 application/json` |

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
    https://vault.rocks/v1/sys/lease/renew
```

### Sample Response

```json
{
  "lease_id": "aws/creds/deploy/abcd-1234...",
  "renewable": true,
  "lease_duration": 2764790
}
```

# `/sys/lease/revoke`

The `/sys/lease/revoke` endpoint is used to revoke secrets.

## Revoke Secret

This endpoint revokes a secret immediately.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `PUT`    | `/sys/lease/revoke`          | `204 (empty body)`     |

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
    https://vault.rocks/v1/sys/lease/revoke
```

# `/sys/lease/revoke-force`

The `/sys/lease/revoke-force` endpoint is used to revoke secrets or tokens 
based on prefix while ignoring backend errors.

## Revoke Force

This endpoint revokes all secrets or tokens generated under a given prefix
immediately. Unlike `/sys/lease/revoke-prefix`, this path ignores backend errors
encountered during revocation. This is _potentially very dangerous_ and should
only be used in specific emergency situations where errors in the backend or the
connected backend service prevent normal revocation.

By ignoring these errors, Vault abdicates responsibility for ensuring that the
issued credentials or secrets are properly revoked and/or cleaned up. Access to
this endpoint should be tightly controlled.

| Method   | Path                               | Produces               |
| :------- | :--------------------------------- | :--------------------- |
| `PUT`    | `/sys/lease/revoke-force/:prefix`  | `204 (empty body)`     |

### Parameters

- `prefix` `(string: <required>)` – Specifies the prefix to revoke. This is
  specified as part of the URL.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request PUT \
    https://vault.rocks/v1/sys/lease/revoke-force/aws/creds
```

# `/sys/lease/revoke-prefix`

The `/sys/lease/revoke-prefix` endpoint is used to revoke secrets or tokens based on
prefix.

## Revoke Prefix

This endpoint revokes all secrets (via a lease ID prefix) or tokens (via the
tokens' path property) generated under a given prefix immediately. This requires
`sudo` capability and access to it should be tightly controlled as it can be
used to revoke very large numbers of secrets/tokens at once.

| Method   | Path                               | Produces               |
| :------- | :--------------------------------- | :--------------------- |
| `PUT`    | `/sys/lease/revoke-prefix/:prefix` | `204 (empty body)`     |

### Parameters

- `prefix` `(string: <required>)` – Specifies the prefix to revoke. This is
  specified as part of the URL.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request PUT \
    https://vault.rocks/v1/sys/lease/revoke-prefix/aws/creds
```
