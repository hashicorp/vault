---
layout: "api"
page_title: "/sys/wrapping/unwrap - HTTP API"
sidebar_current: "docs-http-system-wrapping-unwrap"
description: |-
  The `/sys/wrapping/unwrap` endpoint unwraps a wrapped response.
---

# `/sys/wrapping/unwrap`

The `/sys/wrapping/unwrap` endpoint unwraps a wrapped response.

## Wrapping Unwrap

This endpoint returns the original response inside the given wrapping token.
Unlike simply reading `cubbyhole/response` (which is deprecated), this endpoint
provides additional validation checks on the token, returns the original value
on the wire rather than a JSON string representation of it, and ensures that the
response is properly audit-logged.

This endpoint can be used by using a wrapping token as the client token in the
API call, in which case the `token` parameter is not required; or, a different
token with permissions to access this endpoint can make the call and pass in the
wrapping token in the `token` parameter. Do _not_ use the wrapping token in both
locations; this will cause the wrapping token to be revoked but the value to be
unable to be looked up, as it will basically be a double-use of the token!

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `POST`   | `/sys/wrapping/unwrap`       | `200 application/json` |

### Parameters

- `token` `(string: "")` – Specifies the wrapping token ID. This is required if
  the client token is not the wrapping token. Do not use the wrapping token in
  both locations.

### Sample Payload

```json
{
  "token": "abcd1234..."
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    https://vault.rocks/v1/sys/wrapping/unwrap
```

### Sample Response

```json
{
  "request_id": "8e33c808-f86c-cff8-f30a-fbb3ac22c4a8",
  "lease_id": "",
  "lease_duration": 2592000,
  "renewable": false,
  "data": {
    "zip": "zap"
  },
  "warnings": null
}
```
