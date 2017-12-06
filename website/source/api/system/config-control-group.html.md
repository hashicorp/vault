---
layout: "api"
page_title: "/sys/config/control-group - HTTP API"
sidebar_current: "docs-http-system-config-control-group"
description: |-
  The '/sys/config/control-group' endpoint configures control groups.
---

# `/sys/config/control-group`

The `/sys/config/control-group` endpoint is used to configure Control Group 
settings.

## Read Control Group Settings

This endpoint returns the current Control Group configuration.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `GET`    | `/sys/config/control-group` | `200 application/json` |

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    https://vault.rocks/v1/sys/config/control-group
```

### Sample Response

```json
{
  "max_ttl": "4h"
}
```

## Configure Control Group Settings

This endpoint allows configuring control groups.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `PUT`    | `/sys/config/control-group` | `204 (empty body)` |

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
    https://vault.rocks/v1/sys/config/control-group
```

## Delete Control Group Settings

This endpoint removes any control group configuration.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `DELETE` | `/sys/config/control-group` | `204 (empty body)` |

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request DELETE \
    https://vault.rocks/v1/sys/config/control-group
```
