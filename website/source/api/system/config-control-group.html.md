---
layout: "api"
page_title: "/sys/config/control-group - HTTP API"
sidebar_title: "<code>/sys/config/control-group</code>"
sidebar_current: "api-http-system-config-control-group"
description: |-
  The '/sys/config/control-group' endpoint configures control groups.
---

# `/sys/config/control-group`

~> **Enterprise Only** – These endpoints require Vault Enterprise.

The `/sys/config/control-group` endpoint is used to configure Control Group 
settings.

## Read Control Group Settings

This endpoint returns the current Control Group configuration.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `GET`    | `/sys/config/control-group` |

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    http://127.0.0.1:8200/v1/sys/config/control-group
```

### Sample Response

```json
{
  "max_ttl": "4h"
}
```

## Configure Control Group Settings

This endpoint allows configuring control groups.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `PUT`    | `/sys/config/control-group` |

### Parameters

- `max_ttl` `int` – The maximum ttl for a control group wrapping token.  This can be provided in seconds or duration (2h).

### Sample Payload

```json
{
  "max_ttl": "4h"
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request PUT \
    --data @payload.json \
    http://127.0.0.1:8200/v1/sys/config/control-group
```

## Delete Control Group Settings

This endpoint removes any control group configuration.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `DELETE` | `/sys/config/control-group` |

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request DELETE \
    http://127.0.0.1:8200/v1/sys/config/control-group
```
