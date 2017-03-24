---
layout: "api"
page_title: "/sys/revoke-force - HTTP API"
sidebar_current: "docs-http-system-revoke-force"
description: |-
  The `/sys/revoke-force` endpoint is used to revoke secrets or tokens based on
  prefix while ignoring backend errors.
---

# `/sys/revoke-force`

The `/sys/revoke-force` endpoint is used to revoke secrets or tokens based on
prefix while ignoring backend errors.

## Revoke Force

This endpoint revokes all secrets or tokens generated under a given prefix
immediately. Unlike `/sys/revoke-prefix`, this path ignores backend errors
encountered during revocation. This is _potentially very dangerous_ and should
only be used in specific emergency situations where errors in the backend or the
connected backend service prevent normal revocation.

By ignoring these errors, Vault abdicates responsibility for ensuring that the
issued credentials or secrets are properly revoked and/or cleaned up. Access to
this endpoint should be tightly controlled.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `PUT`    | `/sys/revoke-force/:prefix`  | `204 (empty body)`     |

### Parameters

- `prefix` `(string: <required>)` – Specifies the prefix to revoke. This is
  specified as part of the URL.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request PUT \
    https://vault.rocks/v1/sys/revoke-force/aws/creds
```
