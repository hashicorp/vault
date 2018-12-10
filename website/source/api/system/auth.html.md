---
layout: "api"
page_title: "/sys/auth - HTTP API"
sidebar_title: "<code>/sys/auth</code>"
sidebar_current: "api-http-system-auth"
description: |-
  The `/sys/auth` endpoint is used to manage auth methods in Vault.
---

# `/sys/auth`

The `/sys/auth` endpoint is used to list, create, update, and delete auth
methods. Auth methods convert user or machine-supplied information into a
token which can be used for all future requests.

## List Auth Methods

This endpoint lists all enabled auth methods.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `GET`    | `/sys/auth`                  | `200 application/json` |

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    http://127.0.0.1:8200/v1/sys/auth
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

## Enable Auth Method

This endpoint enables a new auth method. After enabling, the auth method can
be accessed and configured via the auth path specified as part of the URL. This
auth path will be nested under the `auth` prefix.

For example, enable the "foo" auth method will make it accessible at
`/auth/foo`.

- **`sudo` required** – This endpoint requires `sudo` capability in addition to
  any path-specific capabilities.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `POST`   | `/sys/auth/:path`            | `204 (empty body)`     |

### Parameters

- `path` `(string: <required>)` – Specifies the path in which to enable the auth
  method. This is part of the request URL.

- `description` `(string: "")` – Specifies a human-friendly description of the
  auth method.

- `type` `(string: <required>)` – Specifies the name of the authentication
  method type, such as "github" or "token".

- `config` `(map<string|string>: nil)` – Specifies configuration options for
  this auth method. These are the possible values:

  - `default_lease_ttl` `(string: "")` - The default lease duration, specified
     as a string duration like "5s" or "30m".

  - `max_lease_ttl` `(string: "")` - The maximum lease duration, specified as a
     string duration like "5s" or "30m".

  - `audit_non_hmac_request_keys` `(array: [])` - Comma-separated list of keys
     that will not be HMAC'd by audit devices in the request data object.

  - `audit_non_hmac_response_keys` `(array: [])` - Comma-separated list of keys
     that will not be HMAC'd by audit devices in the response data object.

  - `listing_visibility` `(string: "")` - Specifies whether to show this mount
     in the UI-specific listing endpoint.

  - `passthrough_request_headers` `(array: [])` - Comma-separated list of headers
     to whitelist and pass from the request to the backend.

Additionally, the following options are allowed in Vault open-source, but
relevant functionality is only supported in Vault Enterprise:

- `local` `(bool: false)` – Specifies if the auth method is local only. Local
  auth methods are not replicated nor (if a secondary) removed by replication.

  ~> ** Warning:** Remember, policies when using replication secondaries are
  validated by the local cluster. An administrator that can set up a local auth
  method mount can assign policies to tokens that are valid on the replication
  primary if a request is forwarded. Never give untrusted administrators the
  ability to assign policies or configure authentication methods.

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
    http://127.0.0.1:8200/v1/sys/auth/my-auth
```

## Disable Auth Method

This endpoint disables the auth method at the given auth path.

- **`sudo` required** – This endpoint requires `sudo` capability in addition to
  any path-specific capabilities.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `DELETE` | `/sys/auth/:path`            | `204 (empty body)`     |

### Parameters

- `path` `(string: <required>)` – Specifies the path to disable. This is part of
  the request URL.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request DELETE \
    http://127.0.0.1:8200/v1/sys/auth/my-auth
```

## Read Auth Method Tuning

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
    http://127.0.0.1:8200/v1/sys/auth/my-auth/tune
```

### Sample Response

```json
{
  "default_lease_ttl": 3600,
  "max_lease_ttl": 7200
}
```

## Tune Auth Method

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

- `description` `(string: "")` – Specifies the description of the mount. This
  overrides the current stored value, if any.

- `audit_non_hmac_request_keys` `(array: [])` - Specifies the comma-separated
  list of keys that will not be HMAC'd by audit devices in the request data
  object.

- `audit_non_hmac_response_keys` `(array: [])` - Specifies the comma-separated
  list of keys that will not be HMAC'd by audit devices in the response data
  object.

- `listing_visibility` `(string: "")` - Specifies whether to show this mount
    in the UI-specific listing endpoint. Valid values are `"unauth"` or `""`.

- `passthrough_request_headers` `(array: [])` - Comma-separated list of headers
    to whitelist and pass from the request to the backend.

- `token_type` `(string: "")` – Specifies the type of tokens that should be
  returned by the mount. The following values are available:

  - `default-service`: Unless the auth method requests a different type, issue
    service tokens
  - `default-batch`: Unless the auth method requests a different type, issue
    batch tokens
  - `service`: Override any auth method preference and always issue service
    tokens from this mount
  - `batch`: Override any auth method preference and always issue batch tokens
    from this mount

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
    http://127.0.0.1:8200/v1/sys/auth/my-auth/tune
```
