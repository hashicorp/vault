---
layout: "api"
page_title: "Consul Secret Backend - HTTP API"
sidebar_current: "docs-http-secret-consul"
description: |-
  This is the API documentation for the Vault Consul secret backend.
---

# Consul Secret Backend HTTP API

This is the API documentation for the Vault Consul secret backend. For general
information about the usage and operation of the Consul backend, please see the
[Vault Consul backend documentation](/docs/secrets/consul/index.html).

This documentation assumes the Consul backend is mounted at the `/consul` path
in Vault. Since it is possible to mount secret backends at any location, please
update your API calls accordingly.

## Configure Access

This endpoint configures the access information for Consul. This access
information is used so that Vault can communicate with Consul and generate
Consul tokens.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `POST`   | `/consul/config/access`      | `204 (empty body)`     |

### Parameters

- `address` `(string: <required>)` – Specifies the address of the Consul
  instance, provided as `"host:port"` like `"127.0.0.1:8500"`.

- `scheme` `(string: "http")` – Specifies the URL scheme to use.

- `token` `(string: <required>)` – Specifies the Consul ACL token to use. This
  must be a management type token.

### Sample Payload

```json
{
  "address": "127.0.0.1:8500",
  "scheme": "https",
  "token": "adha..."
}
```

### Sample Request

```
$ curl \
    --request POST \
    --header "X-Vault-Token: ..." \
    --data @payload.json \
    https://vault.rocks/v1/consul/config/access
```

## Create/Update Role

This endpoint creates or updates the Consul role definition. If the role does
not exist, it will be created. If the role already exists, it will receive
updated attributes.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `POST`   | `/consul/roles/:name`        | `204 (empty body)`     |

### Parameters

- `name` `(string: <required>)` – Specifies the name of an existing role against
  which to create this Consul credential. This is part of the request URL.

- `lease` `(string: "")` – Specifies the lease for this role. This is provided
  as a string duration with a time suffix like `"30s"` or `"1h"`. If not
  provided, the default Vault lease is used.

- `policy` `(string: <required>)` – Specifies the base64 encoded ACL policy. The
  ACL format can be found in the [Consul ACL
  documentation](https://www.consul.io/docs/internals/acl.html). This is
  required unless the `token_type` is `management`.

- `token_type` `(string: "client")` - Specifies the type of token to create when
  using this role. Valid values are `"client"` or `"management"`.

### Sample Payload

To create management tokens:

```json
{
  "token_type": "management"
}
```

To create a client token with a custom policy:

```json
{
  "policy": "abd2...=="
}
```

### Sample Request

```
$ curl \
    --request POST \
    --header "X-Vault-Token: ..." \
    --data @payload.json \
    https://vault.rocks/v1/consul/roles/example-role
```

## Read Role

This endpoint queries for information about a Consul role with the given name.
If no role exists with that name, a 404 is returned.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `GET`    | `/consul/roles/:name`        | `200 application/json` |

### Parameters

- `name` `(string: <required>)` – Specifies the name of the role to query. This
  is part of the request URL.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    https://vault.rocks/v1/consul/roles/example-role
```

### Sample Response

```json
{
  "data": {
    "policy": "abd2...==",
    "lease": "1h0m0s",
    "token_type": "client"
  }
}
```

## List Roles

This endpoint lists all existing roles in the backend.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `LIST`    | `/consul/roles`             | `200 application/json` |
| `GET`     | `/consul/roles?list=true`   | `200 application/json` |

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request LIST \
    https://vault.rocks/v1/consul/roles
```

### Sample Response

```json
{
  "data": {
    "keys": [
      "example-role"
    ]
  }
}
```

## Delete Role

This endpoint deletes a Consul role with the given name. Even if the role does
not exist, this endpoint will still return a successful response.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `DELETE` | `/consul/roles/:name`        | `204 (empty body)`     |

### Parameters

- `name` `(string: <required>)` – Specifies the name of the role to delete. This
  is part of the request URL.

### Sample Request

```
$ curl \
    --request DELETE \
    --header "X-Vault-Token: ..." \
    https://vault.rocks/v1/consul/roles/example-role
```

## Generate Credential

This endpoint generates a dynamic Consul token based on the given role
definition.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `GET`    | `/consul/creds/:name`        | `200 application/json` |

### Parameters

- `name` `(string: <required>)` – Specifies the name of an existing role against
  which to create this Consul credential. This is part of the request URL.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    https://vault.rocks/v1/consul/creds/example-role
```

### Sample Response

```json
{
  "data": {
    "token": "973a31ea-1ec4-c2de-0f63-623f477c2510"
  }
}
```
