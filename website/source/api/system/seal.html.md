---
layout: "api"
page_title: "/sys/seal - HTTP API"
sidebar_current: "docs-http-system-seal/"
description: |-
  The `/sys/seal` endpoint seals the Vault.
---

# `/sys/seal`

The `/sys/seal` endpoint seals the Vault.

## Seal

This endpoint seals the Vault. In HA mode, only an active node can be sealed.
Standby nodes should be restarted to get the same effect. Requires a token with
`root` policy or `sudo` capability on the path.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `PUT`    | `/sys/seal`                  | `204 (empty body)`     |

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request PUT \
    https://vault.rocks/v1/sys/seal
```
