---
layout: "api"
page_title: "/sys/mounts - HTTP API"
sidebar_current: "docs-http-system-mounts"
description: |-
  The `/sys/mounts` endpoint is used manage secret backends in Vault.
---

# `/sys/mounts`

The `/sys/mounts` endpoint is used manage secret backends in Vault.

## List Mounted Secret Backends

This endpoints lists all the mounted secret backends.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `GET`    | `/sys/mounts`                | `200 application/json` |

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    https://vault.rocks/v1/sys/mounts
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
      "plugin_name": "",
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
      "plugin_name": "",
      "seal_wrap": false
    }
  }
}
```

`default_lease_ttl` or `max_lease_ttl` values of 0 mean that the system defaults
are used by this backend.

## Mount Secret Backend

This endpoint mounts a new secret backend at the given path.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `POST`   | `/sys/mounts/:path`          | `204 (empty body)`     |

### Parameters

- `path` `(string: <required>)` – Specifies the path where the secret backend
  will be mounted. This is specified as part of the URL.

- `type` `(string: <required>)` – Specifies the type of the backend, such as
  "aws".

- `description` `(string: "")` – Specifies the human-friendly description of the
  mount.

- `config` `(map<string|string>: nil)` – Specifies configuration options for
  this mount. This is an object with four possible values:

    - `default_lease_ttl`
    - `max_lease_ttl`
    - `force_no_cache`
    - `plugin_name`
    - `seal_wrap`

    These control the default and maximum lease time-to-live, force
    disabling backend caching, and option plugin name for plugin backends 
    respectively. The first three options override the global defaults if
    set on a specific mount. The plugin_name can be provided in the config
    map or as a top-level option, with the former taking precedence.
    
    When used with supported seals (`pkcs11`, `awskms`, etc.), `seal_wrap`
    causes key material for supporting mounts to be wrapped by the seal's
    encryption capability. This is currently only supported for `transit` and
    `pki` backends. This is only available in Vault Enterprise.

- `plugin_name` `(string: "")` – Specifies the name of the plugin to
  use based from the name in the plugin catalog. Applies only to plugin
  backends.

Additionally, the following options are allowed in Vault open-source, but 
relevant functionality is only supported in Vault Enterprise:

- `local` `(bool: false)` – Specifies if the secret backend is a local mount  
  only. Local mounts are not replicated nor (if a secondary) removed by
  replication.

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
    https://vault.rocks/v1/sys/mounts/my-mount
```

## Unmount Secret Backend

This endpoint un-mounts the mount point specified in the URL.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `DELETE` | `/sys/mounts/:path`          | `204 (empty body)    ` |

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request DELETE \
    https://vault.rocks/v1/sys/mounts/my-mount
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
    https://vault.rocks/v1/sys/mounts/my-mount/tune
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
    https://vault.rocks/v1/sys/mounts/my-mount/tune
```
