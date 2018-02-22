---
layout: "api"
page_title: "/sys/capabilities-self - HTTP API"
sidebar_current: "docs-http-system-capabilities-self"
description: |-
  The `/sys/capabilities-self` endpoint is used to fetch the capabilities of
  client token on a given path.
---

# `/sys/capabilities-self`

The `/sys/capabilities-self` endpoint is used to fetch the capabilities of a
the supplied token.  The capabilities returned will be derived from the
policies that are on the token, and from the policies to which token is
entitled to through the entity and entity's group memberships.

## Query Self Capabilities

This endpoint returns the capabilities of client token on the given path. The
client token is the Vault token with which this API call is made.

| Method   | Path                     | Produces               |
| :------- | :----------------------- | :--------------------- |
| `POST`   | `/sys/capabilities-self` | `200 application/json` |


### Parameters

- `paths` `(list: <required>)` – Paths on which capabilities are being queried.

### Sample Payload

```json
{
  "paths": ["secret/foo", "secret/bar"]
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    https://vault.rocks/v1/sys/capabilities-self
```

### Sample Response

```json
{
  "secret/bar": [
    "sudo",
    "update"
  ],
  "secret/foo": [
    "delete",
    "list",
    "read",
    "update"
  ]
}
