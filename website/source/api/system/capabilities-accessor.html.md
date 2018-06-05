---
layout: "api"
page_title: "/sys/capabilities-accessor - HTTP API"
sidebar_current: "docs-http-system-capabilities-accessor"
description: |-
  The `/sys/capabilities-accessor` endpoint is used to fetch the capabilities of
  the token associated with an accessor, on the given paths.
---

# `/sys/capabilities-accessor`

The `/sys/capabilities-accessor` endpoint is used to fetch the capabilities of
the token associated with the given accessor. The capabilities returned will be
derived from the policies that are on the token, and from the policies to which
the token is entitled to through the entity and entity's group memberships.

## Query Token Accessor Capabilities

This endpoint returns the capabilities of the token associated with the given
accessor, for the given path. Multiple paths are taken in at once and the
capabilities of the token associated with the given accessor for each path is
returned. For backwards compatibility, if a single path is supplied, a
`capabilities` field will also be returned.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `POST`   | `/sys/capabilities-accessor` | `200 application/json` |

### Parameters

- `accessor` `(string: <required>)` – Accessor of the token for which
  capabilities are being queried.

- `paths` `(list: <required>)` – Paths on which capabilities are being
  queried.

### Sample Payload

```json
{
  "accessor": "abcd1234",
  "paths": ["secret/foo"]
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    http://127.0.0.1:8200/v1/sys/capabilities-accessor
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
