---
layout: "api"
page_title: "/sys/internal/ui/mounts - HTTP API"
sidebar_current: "docs-http-system-internal-ui-mounts"
description: |-
  The `/sys/internal/ui/mounts` endpoint is used to manage mount listing visibility.
---

# `/sys/internal/ui/mounts`

The `/sys/internal/ui/mounts` endpoint is used to manage mount listing visibility. This is currently only being used internally for the UI.

Due to the nature of its intended usage, there is no guarantee on backwards compatibility for this endpoint.

## Get Available Visible Mounts

This endpoint lists all enabled auth methods.

| Method |           Path            |        Produces        |
| :----- | :------------------------ | :--------------------- |
| `GET`  | `/sys/internal/ui/mounts` | `200 application/json` |


### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    http://127.0.0.1:8200/v1/sys/internal/ui/mounts
```

### Sample Response

```json
{
  "auth": {
    "github/": {
      "description": "GitHub auth",
      "type": "github"
    }
  },
  "secret": {
    "custom-secrets/": {
      "description": "Custom secrets",
      "type": "kv"
    }
  }
}
```