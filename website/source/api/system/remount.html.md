---
layout: "api"
page_title: "/sys/remount - HTTP API"
sidebar_current: "docs-http-system-remount"
description: |-
  The '/sys/remount' endpoint is used remount a mounted backend to a new endpoint.
---

# `/sys/remount`

The `/sys/remount` endpoint is used remount a mounted backend to a new endpoint.

## Remount Backend

This endpoint remounts an already-mounted backend to a new mount point.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `POST`   | `/sys/remount`               | `204 (empty body)`     |

### Parameters

- `from` `(string: <required>)` – Specifies the previous mount point.

- `to` `(string: <required>)` – Specifies the new destination mount point.

### Sample Payload

```json
{
  "from": "secret",
  "to": "new-secret"
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    https://vault.rocks/v1/sys/remount
```
