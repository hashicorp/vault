---
layout: "api"
page_title: "/sys/leader - HTTP API"
sidebar_current: "docs-http-system-leader"
description: |-
  The `/sys/leader` endpoint is used to check the high availability status and
  current leader of Vault.
---

# `/sys/leader`

The `/sys/leader` endpoint is used to check the high availability status and
current leader of Vault.

## Read Leader Status

This endpoint returns the high availability status and current leader instance
of Vault.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `GET`    | `/sys/leader`                | `200 application/json` |

### Sample Request

```
$ curl \
    https://vault.rocks/v1/sys/leader
```

### Sample Response

```json
{
  "ha_enabled": true,
  "is_self": false,
  "leader_address": "https://127.0.0.1:8200/",
  "leader_cluster_address": "https://127.0.0.1:8201/"
}
```
