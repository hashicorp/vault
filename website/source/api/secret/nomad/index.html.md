---
layout: "api"
page_title: "Nomad Secret Backend - HTTP API"
sidebar_current: "docs-http-secret-nomad"
description: |-
  This is the API documentation for the Vault Nomad secret backend.
---

# Nomad Secret Backend HTTP API

This is the API documentation for the Vault Nomad secret backend. For general
information about the usage and operation of the Nomad backend, please see the
[Vault Nomad backend documentation](/docs/secrets/nomad/index.html).

This documentation assumes the Nomad backend is mounted at the `/nomad` path
in Vault. Since it is possible to mount secret backends at any location, please
update your API calls accordingly.

## Configure Access

This endpoint configures the access information for Nomad. This access
information is used so that Vault can communicate with Nomad and generate
Nomad tokens.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `POST`   | `/nomad/config/access`       | `204 (empty body)`     |

### Parameters

- `address` `(string: <required>)` – Specifies the address of the Nomad
  instance, provided as `"protocol://host:port"` like `"http://127.0.0.1:4646"`.

- `token` `(string: <required>)` – Specifies the Nomad Management token to use.

### Sample Payload

```json
{
  "address": "http://127.0.0.1:4646",
  "token": "adha..."
}
```

### Sample Request

```
$ curl \
    --request POST \
    --header "X-Vault-Token: ..." \
    --data @payload.json \
    https://vault.rocks/v1/nomad/config/access
```

## Configure Lease

This endpoint configures the lease settings for generated tokens.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `POST`   | `/nomad/config/lease`     | `204 (empty body)` |

### Parameters

- `ttl` `(int: 0)` – Specifies the lease ttl provided in seconds.

- `max_ttl` `(int: 0)` – Specifies the maximum ttl provided in seconds.

### Sample Payload

```json
{
  "ttl": 1800,
  "max_ttl": 3600
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    https://vault.rocks/v1/nomad/config/lease
```

## Create/Update Role

This endpoint creates or updates the Nomad role definition in Vault. If the role does not exist, it will be created. If the role already exists, it will receive
updated attributes.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `POST`   | `/nomad/role/:name`         | `204 (empty body)`     |

### Parameters

- `name` `(string: <required>)` – Specifies the name of an existing role against
  which to create this Nomad tokens. This is part of the request URL.

- `lease` `(string: "")` – Specifies the lease for this role. This is provided
  as a string duration with a time suffix like `"30s"` or `"1h"`. If not
  provided, the default Vault lease is used.

- `policies` `(string: "")` – Comma separated list of Nomad policies the token is going to be created against. These need to be created beforehand in Nomad.

- `global` `(bool: "false")` – Specifies if the token should be global, as defined in the [Nomad Documentation](https://www.nomadproject.io/guides/acl.html#acl-tokens).
ma

- `type` `(string: "client")` - Specifies the type of token to create when
  using this role. Valid values are `"client"` or `"management"`.

### Sample Payload

To create a client token with a custom policy:

```json
{
  "policies": "readonly"
}
```

### Sample Request

```
$ curl \
    --request POST \
    --header "X-Vault-Token: ..." \
    --data @payload.json \
    https://vault.rocks/v1/nomad/role/monitoring
```

## Read Role

This endpoint queries for information about a Nomad role with the given name.
If no role exists with that name, a 404 is returned.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `GET`    | `/nomad/role/:name`         | `200 application/json` |

### Parameters

- `name` `(string: <required>)` – Specifies the name of the role to query. This
  is part of the request URL.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    https://vault.rocks/v1/nomad/role/monitoring
```

### Sample Response

```json
{
    "auth": null,
    "data": {
        "lease": "0s",
        "policies": [
            "example"
        ],
        "token_type": "client"
    },
    "lease_duration": 0,
    "lease_id": "",
    "renewable": false,
    "request_id": "f4c7ee18-72aa-3b20-a910-93b6274a9dc0",
    "warnings": null,
    "wrap_info": null
}
```

## List Roles

This endpoint lists all existing roles in the backend.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `LIST`    | `/nomad/role`              | `200 application/json` |
| `GET`     | `/nomad/role?list=true`    | `200 application/json` |

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request LIST \
    https://vault.rocks/v1/nomad/role
```

### Sample Response

```json
{
    "auth": null,
    "data": {
        "keys": [
            "example"
        ]
    },
    "lease_duration": 0,
    "lease_id": "",
    "renewable": false,
    "request_id": "d7bb167b-81c5-9606-c214-b34fcda45634",
    "warnings": null,
    "wrap_info": null
}
```

## Delete Role

This endpoint deletes a Nomad role with the given name. Even if the role does
not exist, this endpoint will still return a successful response.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `DELETE` | `/nomad/role/:name`         | `204 (empty body)`     |

### Parameters

- `name` `(string: <required>)` – Specifies the name of the role to delete. This
  is part of the request URL.

### Sample Request

```
$ curl \
    --request DELETE \
    --header "X-Vault-Token: ..." \
    https://vault.rocks/v1/nomad/role/example-role
```

## Generate Credential

This endpoint generates a dynamic Nomad token based on the given role
definition.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `GET`    | `/nomad/creds/:name`         | `200 application/json` |

### Parameters

- `name` `(string: <required>)` – Specifies the name of an existing role against
  which to create this Nomad token. This is part of the request URL.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    https://vault.rocks/v1/nomad/creds/example
```

### Sample Response

```json
{
    "auth": null,
    "data": {
        "accessor_id": "c834ba40-8d84-b0c1-c084-3a31d3383c03",
        "secret_id": "65af6f07-7f57-bb24-cdae-a27f86a894ce"
    },
    "lease_duration": 2764800,
    "lease_id": "nomad/creds/example/c2686da3-2431-b6d6-7bbf-c5b9496dd6d7",
    "renewable": true,
    "request_id": "37a06ca1-8a1d-7f17-bda8-4661289c392b",
    "warnings": null,
    "wrap_info": null
}
```
