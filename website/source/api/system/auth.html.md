---
layout: "api"
page_title: "/sys/auth - HTTP API"
sidebar_current: "docs-http-system-auth"
description: |-
  The `/sys/auth` endpoint is used to manage auth backends in Vault.
---

# `/sys/auth`

The `/sys/auth` endpoint is used to list, create, update, and delete auth
backends. Auth backends convert user or machine-supplied information into a
token which can be used for all future requests.

## List Auth Backends

This endpoint lists all enabled auth backends.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `GET`    | `/sys/auth`                  | `200 application/json` |

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    https://vault.rocks/v1/sys/auth
```

### Sample Response

```json
{
  "github/": {
    "type": "github",
    "description": "GitHub auth"
  },
  "token/": {
    "config": {
      "default_lease_ttl": 0,
      "max_lease_ttl": 0
    },
    "description": "token based credentials",
    "type": "token"
  }
}
```

## Mount Auth Backend

This endpoint enables a new auth backend. After mounting, the auth backend can
be accessed and configured via the auth path specified as part of the URL. This
auth path will be nested under the `auth` prefix.

For example, mounting the "foo" auth backend will make it accessible at
`/auth/foo`.

- **`sudo` required** – This endpoint requires `sudo` capability in addition to
  any path-specific capabilities.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `POST`   | `/sys/auth/:path`            | `204 (empty body)`     |

### Parameters

- `path` `(string: <required>)` – Specifies the path in which to mount the auth
  backend. This is part of the request URL.

- `description` `(string: "")` – Specifies a human-friendly description of the
  auth backend.

- `type` `(string: <required>)` – Specifies the name of the authentication
  backend type, such as "github" or "token".

- `config` `(map<string|string>: nil)` – Specifies configuration options for
  this mount. These are the possible values:

    - `plugin_name`

    The plugin_name can be provided in the config map or as a top-level option, 
    with the former taking precedence.

- `plugin_name` `(string: "")` – Specifies the name of the auth plugin to
  use based from the name in the plugin catalog. Applies only to plugin
  backends.

Additionally, the following options are allowed in Vault open-source, but
relevant functionality is only supported in Vault Enterprise:

- `local` `(bool: false)` – Specifies if the auth backend is a local mount
  only. Local mounts are not replicated nor (if a secondary) removed by
  replication.

### Sample Payload

```json
{
  "type": "github",
  "description": "Login with GitHub"
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    https://vault.rocks/v1/sys/auth/my-auth
```

## Unmount Auth Backend

This endpoint un-mounts the auth backend at the given auth path.

- **`sudo` required** – This endpoint requires `sudo` capability in addition to
  any path-specific capabilities.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `DELETE` | `/sys/auth/:path`            | `204 (empty body)`     |

### Parameters

- `path` `(string: <required>)` – Specifies the path to unmount. This is part of
  the request URL.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request DELETE \
    https://vault.rocks/v1/sys/auth/my-auth
```

## Read Auth Backend Tuning

This endpoint reads the given auth path's configuration. _This endpoint requires
`sudo` capability on the final path, but the same functionality can be achieved
without `sudo` via `sys/mounts/auth/[auth-path]/tune`._

- **`sudo` required** – This endpoint requires `sudo` capability in addition to
  any path-specific capabilities.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `GET`    | `/sys/auth/:path/tune`       | `200 application/json` |

### Parameters

- `path` `(string: <required>)` – Specifies the path in which to tune.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    https://vault.rocks/v1/sys/auth/my-auth/tune
```

### Sample Response

```json
{
  "default_lease_ttl": 3600,
  "max_lease_ttl": 7200
}
```

## Tune Auth Backend

Tune configuration parameters for a given auth path. _This endpoint
requires `sudo` capability on the final path, but the same functionality
can be achieved without `sudo` via `sys/mounts/auth/[auth-path]/tune`._

- **`sudo` required** – This endpoint requires `sudo` capability in addition to
  any path-specific capabilities.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `POST`   | `/sys/auth/:path/tune`       | `204 (empty body)`     |

### Parameters

- `default_lease_ttl` `(int: 0)` – Specifies the default time-to-live. If set on
  a specific auth path, this overrides the global default.

- `max_lease_ttl` `(int: 0)` – Specifies the maximum time-to-live. If set on a
  specific auth path, this overrides the global default.

### Sample Payload

```json
{
  "default_lease_ttl": 1800,
  "max_lease_ttl": 86400
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    https://vault.rocks/v1/sys/auth/my-auth/tune
```
