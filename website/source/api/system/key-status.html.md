---
layout: "api"
page_title: "/sys/key-status - HTTP API"
sidebar_current: "docs-http-system-key-status"
description: |-
  The `/sys/key-status` endpoint is used to query info about the current
  encryption key of Vault.
---

# `/sys/key-status`

The `/sys/key-status` endpoint is used to query info about the current
encryption key of Vault.

## Get Encryption Key Status

This endpoint returns information about the current encryption key used by
Vault.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `GET`    | `/sys/key-status`            | `200 application/json` |


### Sample Request

### Sample Response

```json
{
  "term": 3,
  "install_time": "2015-05-29T14:50:46.223692553-07:00"
}
```

The `term` parameter is the sequential key number, and `install_time` is the
time that encryption key was installed.
