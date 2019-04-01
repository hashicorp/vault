---
layout: "api"
page_title: "/sys/config/ui - HTTP API"
sidebar_title: "<code>/sys/config/ui</code>"
sidebar_current: "api-http-system-config-ui"
description: |-
  The '/sys/config/ui' endpoint configures the UI.
---

# `/sys/config/ui`

The `/sys/config/ui` endpoint is used to configure UI settings.

- **`sudo` required** – All UI endpoints require `sudo` capability in
  addition to any path-specific capabilities.

## Read UI Settings

This endpoint returns the given UI header configuration.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `GET`    | `/sys/config/ui/headers/:name` |

### Parameters

- `name` `(string: <required>)` – The name of the custom header.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    http://127.0.0.1:8200/v1/sys/config/ui/headers/X-Custom-Header
```

### Sample Response

```json
{
  "value": "custom-value"
}
```

## Configure UI Headers

This endpoint allows configuring the values to be returned for the UI header.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `PUT`    | `/sys/config/ui/headers/:name` |

### Parameters

- `name` `(string: <required>)` – The name of the custom header.

- `values` `(list: <required>)` - The values to be returned from the header.

### Sample Payload

```json
{
  "values": ["custom value 1", "custom value 2"]
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request PUT \
    --data @payload.json \
    http://127.0.0.1:8200/v1/sys/config/ui/headers/X-Custom-Header
```

## Delete a UI Header

This endpoint removes a UI header.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `DELETE` | `/sys/config/ui/headers/:name`|

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request DELETE \
    http://127.0.0.1:8200/v1/sys/config/ui/headers/X-Custom-Header
```

## List UI Headers

This endpoint returns a list of configured UI headers.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `LIST`   | `/sys/config/ui/headers`   |


### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request LIST \
    http://127.0.0.1:8200/v1/sys/config/ui/headers
```

### Sample Response

```json
{
  "data":{
    "keys":[
      "X-Custom...",
      "X-Header...",
    ]
  }
}
```
