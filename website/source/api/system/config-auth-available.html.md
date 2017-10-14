---
layout: "api"
page_title: "/sys/config/auth-available - HTTP API"
sidebar_current: "docs-http-system-config-auth-available"
description: |-
  The `/sys/config/auth-available` endpoint is used to configure which auth
  backends are available via the auth-available endpoint
---

# `/sys/config/auth-available`

The `/sys/config/auth-available` endpoint is used to configure auth backends
visible in the unauthenticated auth-available endpoint.

## Read All Enabled Backend Mounts

This endpoint lists all backends enabled for access via the auth-available
endpoint.


| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `GET`    | `/sys/config/auth-available` | `200 application/json` |

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    https://vault.rocks/v1/sys/config/auth-available
```

### Sample Response

```json
{
  "paths": {
    "token/": {}
  }
}
```

## Read Single Enabled Backend Mount

This endpoint lists information about a single entry in whitelist

- **`sudo` required** – This endpoint requires `sudo` capability in addition to
  any path-specific capabilities.

| Method   | Path                               | Produces               |
| :------- | :--------------------------------- | :--------------------- |
| `GET`    | `/sys/config/auth-available/:path` | `200 application/json` |

### Parameters

- `path` `(string: <required>)` – Specifies the mount path of the auth backend
  to query. Vault will automatically append a "/" if necessary (and not a
  wildcard entry).

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    https://vault.rocks/v1/sys/config/auth-available/token
```

### Sample Response

```json
{
  "token": {}
}
```

## Enable Auth Backend Mount In `auth-available` Endpoint

This endpoint enables an auth backend mount to be returned by the
`sys/auth-available` endpoint. If the auth backend is subsequently unmounted,
this entry will also be removed.

- **`sudo` required** – This endpoint requires `sudo` capability in addition to
  any path-specific capabilities.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `PUT`    | `/sys/config/auth-available/:path` | `204 (empty body)` |

### Parameters

- `path` – Specifies the auth backend path to whitelist in the
  `sys/auth-available` endpoint, e.g., `token`. Accepts wildcard suffixes, e.g.,
  `tok*` would match `token` while `*` would match all backends. Errors out if
   `path` is not currently mounted as an auth backend.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request PUT \
    https://vault.rocks/v1/sys/config/auth-available/token
```

## Disable Auth Backend Mount From `auth-available` Endpoint

This endpoint removes an auth backend from being returned by the
`sys/auth-available` endpoint.

- **`sudo` required** – This endpoint requires `sudo` capability in addition to
  any path-specific capabilities.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `DELETE` | `/sys/config/auth-available/:path` | `204 (empty body)` |

### Parameters

- `path` – Specifies the auth backend path to remove from the
  `sys/auth-available` endpoint whitelist.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request DELETE \
    https://vault.rocks/v1/sys/config/auth-available/token
```
