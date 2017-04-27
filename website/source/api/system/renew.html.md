---
layout: "api"
page_title: "/sys/renew - HTTP API"
sidebar_current: "docs-http-system-renew"
description: |-
  The `/sys/renew` endpoint is used to renew secrets.
---

# `/sys/renew`

The `/sys/renew` endpoint is used to renew secrets.

## Renew Secret

This endpoint renews a secret, requesting to extend the lease.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `PUT`    | `/sys/renew`                 | `200 application/json` |

### Parameters

- `lease_id` `(string: <required>)` – Specifies the ID of the lease to extend.
  This can be specified as part of the URL or as part of the request body.

- `increment` `(int: 0)` – Specifies the requested amount of time (in seconds)
  to extend the lease.

### Sample Payload

```json
{
  "lease_id": "aws/creds/deploy/abcd-1234...",
  "increment": 1800
}
```

### Sample Request

With the `lease_id` in the request body:

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request PUT \
    --data @payload.json \
    https://vault.rocks/v1/sys/renew
```

### Sample Response

```json
{
  "lease_id": "aws/creds/deploy/abcd-1234...",
  "renewable": true,
  "lease_duration": 2764790
}
```
