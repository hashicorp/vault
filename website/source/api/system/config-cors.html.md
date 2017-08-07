---
layout: "api"
page_title: "/sys/config/cors - HTTP API"
sidebar_current: "docs-http-system-config-cors"
description: |-
  The '/sys/config/cors' endpoint configures how the Vault server responds to cross-origin requests.
---

# `/sys/config/cors`

The `/sys/config/cors` endpoint is used to configure CORS settings.

- **`sudo` required** – All CORS endpoints require `sudo` capability in
  addition to any path-specific capabilities.

## Read CORS Settings

This endpoint returns the current CORS configuration.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `GET`    | `/sys/config/cors` | `200 application/json` |

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    https://vault.rocks/v1/sys/config/cors
```

### Sample Response

```json
{
  "enabled": true,
  "allowed_origins": ["http://www.example.com"],
  "allowed_headers": [
    "Content-Type",
    "X-Requested-With",
    "X-Vault-AWS-IAM-Server-ID",
    "X-Vault-No-Request-Forwarding",
    "X-Vault-Token",
    "X-Vault-Wrap-Format",
    "X-Vault-Wrap-TTL",
  ]
}
```

## Configure CORS Settings

This endpoint allows configuring the origins that are permitted to make
cross-origin requests, as well as headers that are allowed on cross-origin requests.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `PUT`    | `/sys/config/cors` | `204 (empty body)` |

### Parameters

- `allowed_origins` `(string or string array: <required>)` – A wildcard (`*`), comma-delimited string, or array of strings specifying the origins that are permitted to make cross-origin requests.

- `allowed_headers` `(string or string array: "" or [])` – A comma-delimited string or array of strings specifying headers that are permitted to be on cross-origin requests. Headers set via this parameter will be appended to the list of headers that Vault allows by default.

### Sample Payload

```json
{
  "allowed_origins": "*",
  "allowed_headers": "X-Custom-Header"
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request PUT \
    --data @payload.json \
    https://vault.rocks/v1/sys/config/cors
```

## Delete CORS Settings

This endpoint removes any CORS configuration.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `DELETE` | `/sys/config/cors` | `204 (empty body)` |

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request DELETE \
    https://vault.rocks/v1/sys/config/cors
```
