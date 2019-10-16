---
layout: "api"
page_title: "/sys/sealwrap/rewrap - HTTP API"
sidebar_title: "<code>/sys/sealwrap/rewrap</code>"
sidebar_current: "api-http-system-sealwrap-rewrap"
description: |-
  The `/sys/sealwrap/rewrap` endpoint is used to rewrap all seal wrapped entries.
---

# `/sys/sealwrap/rewrap`

~> **Enterprise Only** – These endpoints require Vault Enterprise.

The `/sys/sealwrap/rewrap` endpoint is used to rewrap all seal wrapped entries.

## Read Rewrap Status

This endpoint retrieves if a seal rewrap process is currently running.

| Method   | Path                   |
| :------- | :--------------------- |
| `GET`    | `/sys/sealwrap/rewrap` |

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    http://127.0.0.1:8200/v1/sys/sealwrap/rewrap
```

### Sample Response

```json
{
  "request_id": "a6a51003-2576-be0b-9a43-d3bdfeafc2f7",
  "lease_id": "",
  "lease_duration": 0,
  "renewable": false,
  "data": {
    "entries": {
      "failed": 0,
      "processed": 30,
      "succeeded": 30
    },
    "is_running": false
  },
  "warnings": null
}
```

## Start a Seal Rewrap Process

This endpoint starts a seal rewrap process if one is not currently running.
The process will run in the background. Check the vault server logs for status
and progress updates.

| Method   | Path                   |
| :------- | :--------------------- |
| `POST`   | `/sys/sealwrap/rewrap` |

The default status codes are:

- `200` if a seal rewrap process is already running
- `204` if a seal rewrap process was started

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    http://127.0.0.1:8200/v1/sys/sealwrap/rewrap
```
