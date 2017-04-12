---
layout: "api"
page_title: "/sys/revoke-prefix - HTTP API"
sidebar_current: "docs-http-system-revoke-prefix"
description: |-
  The `/sys/revoke-prefix` endpoint is used to revoke secrets or tokens based on
  prefix.
---

# `/sys/revoke-prefix`

The `/sys/revoke-prefix` endpoint is used to revoke secrets or tokens based on
prefix.

## Revoke Prefix

This endpoint revokes all secrets (via a lease ID prefix) or tokens (via the
tokens' path property) generated under a given prefix immediately. This requires
`sudo` capability and access to it should be tightly controlled as it can be
used to revoke very large numbers of secrets/tokens at once.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `PUT`    | `/sys/revoke-prefix/:prefix` | `204 (empty body)`     |

### Parameters

- `prefix` `(string: <required>)` – Specifies the prefix to revoke. This is
  specified as part of the URL.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request PUT \
    https://vault.rocks/v1/sys/revoke-prefix/aws/creds
```
