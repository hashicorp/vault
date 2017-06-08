---
layout: "api"
page_title: "/sys/policy - HTTP API"
sidebar_current: "docs-http-system-policy"
description: |-
  The `/sys/policy` endpoint is used to manage ACL policies in Vault.
---

# `/sys/policy`

The `/sys/policy` endpoint is used to manage ACL policies in Vault.

## List Policies

This endpoint lists all configured policies.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `GET`    | `/sys/policy`                | `200 application/json` |

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    https://vault.rocks/v1/sys/policy
```

### Sample Response

```json
{
  "policies": ["root", "deploy"]
}
```

## Read Policy

This endpoint retrieve the rules for the named policy.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `GET`    | `/sys/policy/:name`          | `200 application/json` |

### Parameters

- `name` `(string: <required>)` – Specifies the name of the policy to retrieve.
  This is specified as part of the request URL.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    https://vault.rocks/v1/sys/policy/my-policy
```

### Sample Response

```json
{
  "rules": "path \"secret/foo\" {..."
}
```

## Create/Update Policy

This endpoint adds a new or updates an existing policy. Once a policy is
updated, it takes effect immediately to all associated users.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `PUT`    | `/sys/policy/:name`          | `204 (empty body)`     |

### Parameters

- `name` `(string: <required>)` – Specifies the name of the policy to create.
  This is specified as part of the request URL.

- `rules` `(string: <required>)` - Specifies the policy document.

### Sample Payload

```json
{
  "rules": "path \"secret/foo\" {..."
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request PUT \
    --data @payload.json \
    https://vault.rocks/v1/sys/policy/my-policy
```

## Delete Policy

This endpoint deletes the policy with the given name. This will immediately
affect all users associated with this policy.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `DELETE` | `/sys/policy/:name`          | `204 (empty body)`     |

### Parameters

- `name` `(string: <required>)` – Specifies the name of the policy to delete.
  This is specified as part of the request URL.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request DELETE \
    https://vault.rocks/v1/sys/policy/my-policy
```
