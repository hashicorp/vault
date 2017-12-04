---
layout: "api"
page_title: "/sys/wrapping/lookup - HTTP API"
sidebar_current: "docs-http-system-wrapping-lookup"
description: |-
  The `/sys/wrapping/lookup` endpoint returns wrapping token properties.
---

# `/sys/wrapping/lookup`

The `/sys/wrapping/lookup` endpoint returns wrapping token properties.

## Wrapping Lookup

This endpoint looks up wrapping properties for the given token.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `POST`   | `/sys/wrapping/lookup`       | `200 application/json` |

### Parameters

- `token` `(string: <required>)` – Specifies the wrapping token ID.

### Sample Payload

```json
{
  "token": "abcd1234"
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    https://vault.rocks/v1/sys/wrapping/lookup
```

### Sample Response

```json
{
  "request_id": "481320f5-fdf8-885d-8050-65fa767fd19b",
  "lease_id": "",
  "lease_duration": 0,
  "renewable": false,
  "data": {
    "creation_path": "sys/wrapping/wrap",
    "creation_time": "2016-09-28T14:16:13.07103516-04:00",
    "creation_ttl": 300
  },
  "wrap_info": null,
  "warnings": null,
  "auth": null
}
```
