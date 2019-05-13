---
layout: "api"
page_title: "/sys/policy - HTTP API"
sidebar_title: "<code>/sys/policy</code>"
sidebar_current: "api-http-system-policy"
description: |-
  The `/sys/policy` endpoint is used to manage ACL policies in Vault.
---

# `/sys/policy`

The `/sys/policy` endpoint is used to manage ACL policies in Vault.

## List Policies

This endpoint lists all configured policies.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `GET`    | `/sys/policy`                |

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    http://127.0.0.1:8200/v1/sys/policy
```

### Sample Response

```json
{
  "policies": ["root", "deploy"]
}
```

## Read Policy

This endpoint retrieve the policy body for the named policy.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `GET`    | `/sys/policy/:name`          |

### Parameters

- `name` `(string: <required>)` – Specifies the name of the policy to retrieve.
  This is specified as part of the request URL.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    http://127.0.0.1:8200/v1/sys/policy/my-policy
```

### Sample Response

```json
{
  "name": "my-policy",
  "rules": "path \"secret/*\"...
}
```

## Create/Update Policy

This endpoint adds a new or updates an existing policy. Once a policy is
updated, it takes effect immediately to all associated users.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `PUT`    | `/sys/policy/:name`          |

### Parameters

- `name` `(string: <required>)` – Specifies the name of the policy to create.
  This is specified as part of the request URL.

- `policy` `(string: <required>)` - Specifies the policy document.

### Sample Payload

```json
{
  "policy": "path \"secret/foo\" {..."
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request PUT \
    --data @payload.json \
    http://127.0.0.1:8200/v1/sys/policy/my-policy
```

## Delete Policy

This endpoint deletes the policy with the given name. This will immediately
affect all users associated with this policy.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `DELETE` | `/sys/policy/:name`          |

### Parameters

- `name` `(string: <required>)` – Specifies the name of the policy to delete.
  This is specified as part of the request URL.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request DELETE \
    http://127.0.0.1:8200/v1/sys/policy/my-policy
```
