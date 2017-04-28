---
layout: "api"
page_title: "/sys/revoke - HTTP API"
sidebar_current: "docs-http-system-revoke/"
description: |-
  The `/sys/revoke` endpoint is used to revoke secrets.
---

# `/sys/revoke`

The `/sys/revoke` endpoint is used to revoke secrets.

## Revoke Secret

This endpoint revokes a secret immediately.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `PUT`    | `/sys/revoke`                | `204 (empty body)`     |

### Parameters

- `lease_id` `(string: <required>)` – Specifies the ID of the lease to revoke.

### Sample Payload

```json
{
  "lease_id": "postgresql/creds/readonly/abcd-1234..."
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request PUT \
    --data @payload.json \
    https://vault.rocks/v1/sys/revoke
```
