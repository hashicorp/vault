---
layout: "api"
page_title: "/sys/wrapping/wrap - HTTP API"
sidebar_current: "docs-http-system-wrapping-wrap"
description: |-
  The `/sys/wrapping/wrap` endpoint wraps the given values in a
  response-wrapped token.
---

# `/sys/wrapping/wrap`

The `/sys/wrapping/wrap` endpoint wraps the given values in a response-wrapped
token.

## Wrapping Wrap

This endpoint wraps the given user-supplied data inside a response-wrapped
token.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `POST`   | `/sys/wrapping/wrap`         | `200 application/json` |

### Parameters

- `:any` `(map<string|string>: nil)` – Parameters should be supplied as
  keys/values in a JSON object. The exact set of given parameters will be
  contained in the wrapped response.

### Sample Payload

```json
{
  "foo": "bar",
  "zip": "zap"
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --header "X-Vault-Wrap-TTL: 60" \
    --request POST \
    --data @payload.json \
    https://vault.rocks/v1/sys/wrapping/wrap
```

### Sample Response

```json
{
  "request_id": "",
  "lease_id": "",
  "lease_duration": 0,
  "renewable": false,
  "data": null,
  "warnings": null,
  "wrap_info": {
    "token": "fb79b9d3-d94e-9eb6-4919-c559311133d6",
    "ttl": 300,
    "creation_time": "2016-09-28T14:41:00.56961496-04:00",
    "creation_path": "sys/wrapping/wrap",
  }
}
```
