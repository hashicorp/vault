---
layout: "api"
page_title: "/sys/mounts - HTTP API"
sidebar_title: "<code>/sys/mounts</code>"
sidebar_current: "api-http-system-mounts"
description: |-
  The `/sys/mounts` endpoint is used manage secrets engines in Vault.
---

# `/sys/mounts`

The `/sys/mounts` endpoint is used manage secrets engines in Vault.

## List Mounted Secrets Engines

This endpoints lists all the mounted secrets engines.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `GET`    | `/sys/mounts`                | `200 application/json` |

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    http://127.0.0.1:8200/v1/sys/mounts
```

### Sample Response

```json
{
  "aws": {
    "type": "aws",
    "description": "AWS keys",
    "config": {
      "default_lease_ttl": 0,
      "max_lease_ttl": 0,
      "force_no_cache": false,
      "seal_wrap": false
    }
  },
  "sys": {
    "type": "system",
    "description": "system endpoint",
    "config": {
      "default_lease_ttl": 0,
      "max_lease_ttl": 0,
      "force_no_cache": false,
      "seal_wrap": false
    }
  }
}
```

`default_lease_ttl` or `max_lease_ttl` values of 0 mean that the system defaults
are used by this backend.

## Enable Secrets Engine

This endpoint enables a new secrets engine at the given path.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `POST`   | `/sys/mounts/:path`          | `204 (empty body)`     |

### Parameters

- `path` `(string: <required>)` – Specifies the path where the secrets engine
  will be mounted. This is specified as part of the URL.

- `type` `(string: <required>)` – Specifies the type of the backend, such as
  "aws".

- `description` `(string: "")` – Specifies the human-friendly description of the
  mount.

- `config` `(map<string|string>: nil)` – Specifies configuration options for
  this mount. This is an object with four possible values:

  - `default_lease_ttl` `(string: "")` - The default lease duration, specified
     as a string duration like "5s" or "30m".

  - `max_lease_ttl` `(string: "")` - The maximum lease duration, specified as a
     string duration like "5s" or "30m".

  - `force_no_cache` `(bool: false)` - Disable caching.

  - `audit_non_hmac_request_keys` `(array: [])` - Comma-separated list of keys
     that will not be HMAC'd by audit devices in the request data object.

  - `audit_non_hmac_response_keys` `(array: [])` - Comma-separated list of keys
     that will not be HMAC'd by audit devices in the response data object.

  - `listing_visibility` `(string: "")` - Specifies whether to show this mount
    in the UI-specific listing endpoint. Valid values are `"unauth"` or
    `"hidden"`.  If not set, behaves like `"hidden"`.

  - `passthrough_request_headers` `(array: [])` - Comma-separated list of headers
     to whitelist and pass from the request to the backend.

    These control the default and maximum lease time-to-live, and the force
    disabling backend caching. They override the global defaults if
    set on a specific mount.

    When used with supported seals (`pkcs11`, `awskms`, etc.), `seal_wrap`
    causes key material for supporting mounts to be wrapped by the seal's
    encryption capability. This is currently only supported for `transit` and
    `pki` backends. This is only available in Vault Enterprise.

- `options` `(map<string|string>: nil)` - Specifies mount type specific options
  that are passed to the backend. 
  
    *Key/Value (KV)*  
    - `version` `(string: "1")` - The version of the KV to mount. Set to "2" for mount
      KV v2.

Additionally, the following options are allowed in Vault open-source, but
relevant functionality is only supported in Vault Enterprise:

- `local` `(bool: false)` – Specifies if the secrets engine is a local mount
  only. Local mounts are not replicated nor (if a secondary) removed by
  replication.

- `seal_wrap` `(bool: false)` - Enable seal wrapping for the mount.

### Sample Payload

```json
{
  "type": "aws",
  "config": {
    "force_no_cache": true
  }
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    http://127.0.0.1:8200/v1/sys/mounts/my-mount
```

## Disable Secrets Engine

This endpoint disables the mount point specified in the URL.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `DELETE` | `/sys/mounts/:path`          | `204 (empty body)    ` |

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request DELETE \
    http://127.0.0.1:8200/v1/sys/mounts/my-mount
```

## Read Mount Configuration

This endpoint reads the given mount's configuration. Unlike the `mounts`
endpoint, this will return the current time in seconds for each TTL, which may
be the system default or a mount-specific value.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `GET`   | `/sys/mounts/:path/tune`      | `200 application/json` |

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    http://127.0.0.1:8200/v1/sys/mounts/my-mount/tune
```

### Sample Response

```json
{
  "default_lease_ttl": 3600,
  "max_lease_ttl": 7200,
  "force_no_cache": false
}
```

## Tune Mount Configuration

This endpoint tunes configuration parameters for a given mount point.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `POST`   | `/sys/mounts/:path/tune`     | `204 (empty body)`     |

### Parameters

- `default_lease_ttl` `(int: 0)` – Specifies the default time-to-live. This
  overrides the global default. A value of `0` is equivalent to the system
  default TTL.

- `max_lease_ttl` `(int: 0)` – Specifies the maximum time-to-live. This
  overrides the global default. A value of `0` are equivalent and set to the
  system max TTL.

- `description` `(string: "")` – Specifies the description of the mount. This
  overrides the current stored value, if any.

- `audit_non_hmac_request_keys` `(array: [])` - Specifies the comma-separated
  list of keys that will not be HMAC'd by audit devices in the request data
  object.

- `audit_non_hmac_response_keys` `(array: [])` - Specifies the comma-separated
  list of keys that will not be HMAC'd by audit devices in the response data
  object.

- `listing_visibility` `(string: "")` - Specifies whether to show this mount in
  the UI-specific listing endpoint. Valid values are `"unauth"` or `"hidden"`.
  If not set, behaves like `"hidden"`.

- `passthrough_request_headers` `(array: [])` - Comma-separated list of headers
    to whitelist and pass from the request to the backend.

### Sample Payload

```json
{
  "default_lease_ttl": 1800,
  "max_lease_ttl": 3600
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    http://127.0.0.1:8200/v1/sys/mounts/my-mount/tune
```
