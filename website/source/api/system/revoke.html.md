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
| `PUT`    | `/sys/revoke/:lease_id`      | `204 (empty body)`     |

### Parameters

- `lease_id` `(string: <required>)` – Specifies the ID of the lease to renew.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request PUT \
    https://vault.rocks/v1/sys/revoke/aws/creds/readonly-acbd1234
```
