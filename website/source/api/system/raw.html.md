---
layout: "api"
page_title: "/sys/raw - HTTP API"
sidebar_current: "docs-http-system-raw"
description: |-
  The `/sys/raw` endpoint is access the raw underlying store in Vault.
---

# `/sys/raw`

The `/sys/raw` endpoint is access the raw underlying store in Vault.

## Read Raw

This endpoint reads the value of the key at the given path. This is the raw path
in the storage backend and not the logical path that is exposed via the mount
system.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `GET`    | `/sys/raw/:path`             | `200 application/json` |

### Parameters

- `path` `(string: <required>)` – Specifies the raw path in the storage backend.
  This is specified as part of the URL.

### Sample Request

```
$ curl \
    ---header "X-Vault-Token: ..." \
    https://vault.rocks/v1/sys/raw/secret/foo
```

### Sample Response

```json
{
  "value": "{'foo':'bar'}"
}
```

## Create/Update Raw

This endpoint updates the value of the key at the given path. This is the raw
path in the storage backend and not the logical path that is exposed via the
mount system.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `PUT`    | `/sys/raw/:path`             | `204 (empty body)`     |

### Parameters

- `path` `(string: <required>)` – Specifies the raw path in the storage backend.
  This is specified as part of the URL.

- `value` `(string: <required>)` – Specifies the value of the key.

### Sample Payload

```json
{
  "value": "{\"foo\": \"bar\"}"
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request PUT \
    --data @payload.json \
    https://vault.rocks/v1/sys/raw/secret/foo
```

## Delete Raw

This endpoint deletes the key with given path. This is the raw path in the
storage backend and not the logical path that is exposed via the mount system.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `DELETE` | `/sys/raw/:path`             | `204 (empty body)`     |

### Parameters

- `path` `(string: <required>)` – Specifies the raw path in the storage backend.
  This is specified as part of the URL.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request DELETE \
    https://vault.rocks/v1/sys/raw/secret/foo
```
