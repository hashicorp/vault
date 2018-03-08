---
layout: "api"
page_title: "/sys/capabilities - HTTP API"
sidebar_current: "docs-http-system-capabilities/"
description: |-
  The `/sys/capabilities` endpoint is used to fetch the capabilities of a token
  on the given paths.
---

# `/sys/capabilities`

The `/sys/capabilities` endpoint is used to fetch the capabilities of a token
on the given paths. The capabilities returned will be derived from the policies
that are on the token, and from the policies to which the token is entitled to
through the entity and entity's group memberships.

## Query Token Capabilities

This endpoint returns the list of capabilities of a given token on the given
paths. Multiple paths are taken in at once and the capabilities of the token
for each path is returned. For backwards compatibility, if a single path is
supplied, a `capabilities` field will also be returned.

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
  "paths": ["secret/foo"]
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
  "capabilities": [
    "delete",
    "list",
    "read",
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
