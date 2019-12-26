---
layout: "api"
page_title: "/sys/remount - HTTP API"
sidebar_title: "<code>/sys/remount</code>"
sidebar_current: "api-http-system-remount"
description: |-
  The '/sys/remount' endpoint is used remount a mounted backend to a new endpoint.
---

# `/sys/remount`

The `/sys/remount` endpoint is used remount a mounted backend to a new endpoint.

## Move Backend

This endpoint moves an already-mounted backend to a new mount point.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `POST`   | `/sys/remount`               |

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
    http://127.0.0.1:8200/v1/sys/remount
```
