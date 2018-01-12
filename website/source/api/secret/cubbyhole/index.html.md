---
layout: "api"
page_title: "Cubbyhole - Secrets Engines - HTTP API"
sidebar_current: "docs-http-secret-cubbyhole"
description: |-
  This is the API documentation for the Vault Cubbyhole secrets engine.
---

# Cubbyhole Secrets Engine (API)

This is the API documentation for the Vault Cubbyhole secrets engine. For
general information about the usage and operation of the Cubbyhole secrets
engine, please see the
[Vault Cubbyhole documentation](/docs/secrets/cubbyhole/index.html).

This documentation assumes the Cubbyhole secrets engine is enabled at the
`/cubbyhole` path in Vault. Since it is possible to enable secrets engines at
any location, please update your API calls accordingly.

## Read Secret

This endpoint retrieves the secret at the specified location.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `GET`    | `/cubbyhole/:path`           | `200 application/json` |

### Parameters

- `path` `(string: <required>)` – Specifies the path of the secret to read.
  This is specified as part of the URL.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    https://vault.rocks/v1/cubbyhole/my-secret
```

### Sample Response

```json
{
  "auth": null,
  "data": {
    "foo": "bar"
  },
  "lease_duration": 0,
  "lease_id": "",
  "renewable": false
}
```

## List Secrets

This endpoint returns a list of secret entries at the specified location.
Folders are suffixed with `/`. The input must be a folder; list on a file will
not return a value. The values themselves are not accessible via this command.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `LIST`   | `/cubbyhole/:path`           | `200 application/json` |

### Parameters

- `path` `(string: <required>)` – Specifies the path of the secrets to list.
  This is specified as part of the URL.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request LIST \
    https://vault.rocks/v1/cubbyhole/my-secret
```

### Sample Response

The example below shows output for a query path of `cubbyhole/` when there are
secrets at `cubbyhole/foo` and `cubbyhole/foo/bar`; note the difference in the
two entries.

```json
{
  "auth": null,
  "data": {
    "keys": ["foo", "foo/"]
  },
  "lease_duration": 2764800,
  "lease_id": "",
  "renewable": false
}
```

## Create/Update Secret

This endpoint stores a secret at the specified location.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `POST`   | `/cubbyhole/:path`           | `204 (empty body)`     |
| `PUT`    | `/cubbyhole/:path`           | `204 (empty body)`     |

### Parameters

- `path` `(string: <required>)` – Specifies the path of the secrets to
  create/update. This is specified as part of the URL.

- `:key` `(string: "")` – Specifies a key, paired with an associated value, to
  be held at the given location. Multiple key/value pairs can be specified, and
  all will be returned on a read operation. A key called `ttl` will trigger some
  special behavior; see above for details.

### Sample Payload

```json
{
  "foo": "bar",
  "zip": "zap"
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    https://vault.rocks/v1/cubbyhole/my-secret
```

## Delete Secret

This endpoint deletes the secret at the specified location.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `DELETE` | `/cubbyhole/:path`           | `204 (empty body)`     |

### Parameters

- `path` `(string: <required>)` – Specifies the path of the secret to delete.
  This is specified as part of the URL.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request DELETE \
    https://vault.rocks/v1/cubbyhole/my-secret
```
