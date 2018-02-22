---
layout: "api"
page_title: "/sys/capabilities - HTTP API"
sidebar_current: "docs-http-system-capabilities/"
description: |-
  The `/sys/capabilities` endpoint is used to fetch the capabilities of a token
  on a given path.
---

# `/sys/capabilities`

The `/sys/capabilities` endpoint is used to fetch the capabilities of a token
on a given path. The capabilities returned will be derived from the policies
that are on the token, and from the policies to which token is entitled to
through the entity and entity's group memberships.

## Query Token Capabilities

This endpoint returns the list of capabilities for a provided token.

| Method   | Path                 | Produces               |
| :------- | :------------------- | :--------------------- |
| `POST`   | `/sys/capabilities`  | `200 application/json` |

### Parameters

- `paths` `(list: <required>)` – Paths on which capabilities are being queried.

- `token` `(string: <required>)` – Token for which capabilities are being
  queried.

### Sample Payload

```json
{
  "token": "abcd1234",
  "paths": ["secret/foo", "secret/bar"]
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    https://vault.rocks/v1/sys/capabilities
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
```
