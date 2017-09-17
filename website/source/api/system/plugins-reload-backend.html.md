---
layout: "api"
page_title: "/sys/plugins/reload/backend - HTTP API"
sidebar_current: "docs-http-system-plugins-reload-backend"
description: |-
  The `/sys/plugins/reload/backend` endpoint is used to reload plugin backends.
---

# `/sys/plugins/reload/backend`

The `/sys/plugins/reload/backend` endpoint is used to reload mounted plugin
backends. Either the plugin name (`plugin`) or the desired plugin backend mounts
(`mounts`) must be provided, but not both. In the case that the plugin name is
provided, all mounted paths that use that plugin backend will be reloaded.

## Reload Plugins

This endpoint reloads mounted plugin backends.

| Method   | Path                      -   | Produces               |
| :------- | :---------------------------- | :--------------------- |
| `PUT`    | `/sys/plugins/reload/backend` | `204 (empty body)`     |

### Parameters

- `plugin` `(string: "")` – The name of the plugin to reload, as 
  registered in the plugin catalog.

- `mounts` `(slice: [])` – Array or comma-separated string mount paths 
  of the plugin backends to reload.

### Sample Payload

```json
{
  "plugin": "mock-plugin"
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request PUT
    https://vault.rocks/v1/sys/plugins/reload/backend
```
