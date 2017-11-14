---
layout: "api"
page_title: "/sys/capabilities-accessor - HTTP API"
sidebar_current: "docs-http-system-capabilities-accessor"
description: |-
  The `/sys/capabilities-accessor` endpoint is used to fetch the capabilities of
  the token associated with an accessor, on the given path.
---

# `/sys/capabilities-accessor`

The `/sys/capabilities-accessor` endpoint is used to fetch the capabilities of
a token associated with an accessor. The capabilities returned will be derived
from the policies that are on the token, and from the policies to which token
is entitled to through the entity and entity's group memberships.


## Query Token Accessor Capabilities

This endpoint returns the capabilities of the token associated with an accessor,
for the given path.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `POST`   | `/sys/capabilities-accessor` | `200 application/json` |

### Parameters

- `accessor` `(string: <required>)` – Specifies the accessor of the token to
  check.

- `path` `(string: <required>)` – Specifies the path on which the token's
  capabilities will be checked.

### Sample Payload

```json
{
  "accessor": "abcd1234",
  "path": "secret/foo"
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    https://vault.rocks/v1/sys/capabilities-accessor
```

### Sample Response

```json
{
  "capabilities": ["read", "list"]
}
```
