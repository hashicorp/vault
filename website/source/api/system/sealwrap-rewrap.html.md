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
  "request_id": "4435a8e9-bc19-480e-c512-ad2031d4fe27",
  "lease_id": "",
  "renewable": false,
  "lease_duration": 0,
  "data": {
    "is_running": false
  },
  "wrap_info": null,
  "warnings": null,
  "auth": null
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
