---
layout: "api"
page_title: "/sys/audit - HTTP API"
sidebar_current: "docs-http-system-audit/"
description: |-
  The `/sys/audit` endpoint is used to enable and disable audit backends.
---

# `/sys/audit`

The `/sys/audit` endpoint is used to list, mount, and unmount audit backends.
Audit backends must be enabled before use, and more than one backend may be
enabled at a time.

## List Mounted Audit Backends

This endpoint lists only the mounted audit backends (it does not list all
available audit backends).

- **`sudo` required** – This endpoint requires `sudo` capability in addition to
  any path-specific capabilities.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `GET`    | `/sys/audit`                 | `200 application/json` |

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    https://vault.rocks/v1/sys/audit
```

### Sample Response

```javascript
{
  "file": {
    "type": "file",
    "description": "Store logs in a file",
    "options": {
      "path": "/var/log/vault.log"
    }
  }
}
```

## Mount Audit Backend

This endpoint mounts a new audit backend at the supplied path. The path can be a
single word name or a more complex, nested path.

- **`sudo` required** – This endpoint requires `sudo` capability in addition to
  any path-specific capabilities.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `PUT`    | `/sys/audit/:path`           | `204 (empty body)`     |

### Parameters

- `path` `(string: <required>)` – Specifies the path in which to mount the audit
  backend. This is part of the request URL.

- `description` `(string: "")` – Specifies a human-friendly description of the
  audit backend.

- `options` `(map<string|string>: nil)` – Specifies configuration options to
  pass to the audit backend itself. This is dependent on the audit backend type.

- `type` `(string: <required>)` – Specifies the type of the audit backend.

Additionally, the following options are allowed in Vault open-source, but
relevant functionality is only supported in Vault Enterprise:

- `local` `(bool: false)` – Specifies if the audit backend is a local mount  
  only. Local mounts are not replicated nor (if a secondary) removed by
  replication.

### Sample Payload

```json
{
  "type": "file",
  "options": {
    "path": "/var/log/vault/log"
  }
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request PUT \
    --data @payload.json \
    https://vault.rocks/v1/sys/audit/example-audit
```

## Unmount Audit Backend

This endpoint un-mounts the audit backend at the given path.

- **`sudo` required** – This endpoint requires `sudo` capability in addition to
  any path-specific capabilities.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `DELETE` | `/sys/audit/:path`           | `204 (empty body)`     |

### Parameters

- `path` `(string: <required>)` – Specifies the path of the audit backend to
  delete. This is part of the request URL.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request DELETE \
    https://vault.rocks/v1/sys/audit/example-audit
```
