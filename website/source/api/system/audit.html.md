---
layout: "api"
page_title: "/sys/audit - HTTP API"
sidebar_current: "docs-http-system-audit/"
description: |-
  The `/sys/audit` endpoint is used to enable and disable audit devices.
---

# `/sys/audit`

The `/sys/audit` endpoint is used to list, enable, and disable audit devices.
Audit devices must be enabled before use, and more than one device may be
enabled at a time.

## List Enabled Audit Devices

This endpoint lists only the enabled audit devices (it does not list all
available audit devices).

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

## Enable Audit Device

This endpoint enables a new audit device at the supplied path. The path can be a
single word name or a more complex, nested path.

- **`sudo` required** – This endpoint requires `sudo` capability in addition to
  any path-specific capabilities.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `PUT`    | `/sys/audit/:path`           | `204 (empty body)`     |

### Parameters

- `path` `(string: <required>)` – Specifies the path in which to enable the audit
  device. This is part of the request URL.

- `description` `(string: "")` – Specifies a human-friendly description of the
  audit device.

- `options` `(map<string|string>: nil)` – Specifies configuration options to
  pass to the audit device itself. This is dependent on the audit device type.

- `type` `(string: <required>)` – Specifies the type of the audit device.

Additionally, the following options are allowed in Vault open-source, but
relevant functionality is only supported in Vault Enterprise:

- `local` `(bool: false)` – Specifies if the audit device is a local only. Local
  audit devices are not replicated nor (if a secondary) removed by replication.

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

## Disable Audit Device

This endpoint disables the audit device at the given path.

- **`sudo` required** – This endpoint requires `sudo` capability in addition to
  any path-specific capabilities.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `DELETE` | `/sys/audit/:path`           | `204 (empty body)`     |

### Parameters

- `path` `(string: <required>)` – Specifies the path of the audit device to
  delete. This is part of the request URL.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request DELETE \
    https://vault.rocks/v1/sys/audit/example-audit
```
