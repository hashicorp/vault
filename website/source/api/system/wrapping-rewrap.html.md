---
layout: "api"
page_title: "/sys/wrapping/rewrap - HTTP API"
sidebar_current: "docs-http-system-wrapping-rewrap"
description: |-
  The `/sys/wrapping/rewrap` endpoint can be used to rotate a wrapping token and refresh its TTL.
---

# `/sys/wrapping/rewrap`

The `/sys/wrapping/rewrap` endpoint can be used to rotate a wrapping token and
refresh its TTL.

## Wrapping Rewrap

This endpoint rewraps a response-wrapped token. The new token will use the same
creation TTL as the original token and contain the same response. The old token
will be invalidated. This can be used for long-term storage of a secret in a
response-wrapped token when rotation is a requirement.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `POST`   | `/sys/wrapping/rewrap`       | `200 application/json` |

### Parameters

- `token` `(string: <required>)` – Specifies the wrapping token ID.

### Sample Payload

```json
{
  "token": "abcd1234...",
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
  "request_id": "",
  "lease_id": "",
  "lease_duration": 0,
  "renewable": false,
  "data": null,
  "warnings": null,
  "wrap_info": {
    "token": "3b6f1193-0707-ac17-284d-e41032e74d1f",
    "ttl": 300,
    "creation_time": "2016-09-28T14:22:26.486186607-04:00",
    "creation_path": "sys/wrapping/wrap"
  }
}
```
