---
layout: "api"
page_title: "/sys/rotate - HTTP API"
sidebar_title: "<code>/sys/rotate</code>"
sidebar_current: "api-http-system-rotate"
description: |-
  The `/sys/rotate` endpoint is used to rotate the encryption key.
---

# `/sys/rotate`

The `/sys/rotate` endpoint is used to rotate the encryption key.

## Rotate Encryption Key

This endpoint triggers a rotation of the backend encryption key. This is the key
that is used to encrypt data written to the storage backend, and is not provided
to operators. This operation is done online. Future values are encrypted with
the new key, while old values are decrypted with previous encryption keys.

This path requires `sudo` capability in addition to `update`.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `PUT`    | `/sys/rotate`                |

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request PUT \
    http://127.0.0.1:8200/v1/sys/rotate
```
