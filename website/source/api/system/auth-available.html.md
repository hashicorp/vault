---
layout: "api"
page_title: "/sys/auth-available - HTTP API"
sidebar_current: "docs-http-system-auth-available"
description: |-
  The `/sys/auth-available` endpoint is used to anonymously view auth backends in Vault.
---

# `/sys/auth-available`

The `/sys/auth-available` endpoint is an unauthenticated endpoint used to
list mounted auth backends. It is similar to a `GET` request on the `/sys/auth`
endpoint, except that less data is returned.

## List Auth Backends

This endpoint lists all enabled auth backends.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `GET`    | `/sys/auth-available`        | `200 application/json` |

### Sample Request

```
$ curl https://vault.rocks/v1/sys/auth-available
```

### Sample Response

```json
{
  "github/": {
    "type": "github",
  },
  "token/": {
    "type": "token"
  }
}
```


