---
layout: api
page_title: Consul - Secrets Engines - HTTP API
sidebar_title: Consul
description: This is the API documentation for the Vault Consul secrets engine.
---

# Consul Secrets Engine (API)

This is the API documentation for the Vault Consul secrets engine. For general
information about the usage and operation of the Consul secrets engine, please
see the [Vault Consul documentation](/docs/secrets/consul).

This documentation assumes the Consul secrets engine is enabled at the `/consul`
path in Vault. Since it is possible to enable secrets engines at any location,
please update your API calls accordingly.

## Configure Access

This endpoint configures the access information for Consul. This access
information is used so that Vault can communicate with Consul and generate
Consul tokens.

| Method | Path                    |
| :----- | :---------------------- |
| `POST` | `/consul/config/access` |

### Parameters

- `address` `(string: <required>)` – Specifies the address of the Consul
  instance, provided as `"host:port"` like `"127.0.0.1:8500"`.

- `scheme` `(string: "http")` – Specifies the URL scheme to use.

- `token` `(string: <required>)` – Specifies the Consul ACL token to use. This
  must be a management type token.

- `ca_cert` `(string: "")` - CA certificate to use when verifying Consul server certificate,
  must be x509 PEM encoded.

- `client_cert` `(string: "")` - Client certificate used for Consul's TLS communication,
  must be x509 PEM encoded and if this is set you need to also set client_key.

- `client_key` `(string: "")` - Client key used for Consul's TLS communication,
  must be x509 PEM encoded and if this is set you need to also set client_cert.

### Sample Payload

```json
{
  "address": "127.0.0.1:8500",
  "scheme": "https",
  "token": "adha..."
}
```

### Sample Request

```shell-session
$ curl \
    --request POST \
    --header "X-Vault-Token: ..." \
    --data @payload.json \
    http://127.0.0.1:8200/v1/consul/config/access
```

## Create/Update Role

This endpoint creates or updates the Consul role definition. If the role does
not exist, it will be created. If the role already exists, it will receive
updated attributes.

| Method | Path                  |
| :----- | :-------------------- |
| `POST` | `/consul/roles/:name` |

### Parameters for Consul versions 1.4 and above

- `name` `(string: <required>)` – Specifies the name of an existing role against
  which to create this Consul credential. This is part of the request URL.

- `token_type` `(string: "client")` - Specifies the type of token to create when
  using this role. Valid values are `"client"` or `"management"`.

- `policy` `(string: <policy or policies>)` – Specifies the base64 encoded ACL policy. The
  ACL format can be found in the [Consul ACL
  documentation](https://www.consul.io/docs/internals/acl). This is
  required unless the `token_type` is `management`.

- `policies` `(list: <policy or policies>)` – The list of policies to assign to the generated
  token. This is only available in Consul 1.4 and greater.

- `local` `(bool: false)` - Indicates that the token should not be replicated
  globally and instead be local to the current datacenter. Only available in Consul
  1.4 and greater.

- `ttl` `(duration: "")` – Specifies the TTL for this role. This is provided
  as a string duration with a time suffix like `"30s"` or `"1h"` or as seconds. If not
  provided, the default Vault TTL is used.

- `max_ttl` `(duration: "")` – Specifies the max TTL for this role. This is provided
  as a string duration with a time suffix like `"30s"` or `"1h"` or as seconds. If not
  provided, the default Vault Max TTL is used.

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

```shell-session
$ curl \
    --request POST \
    --header "X-Vault-Token: ..." \
    --data @payload.json \
    http://127.0.0.1:8200/v1/consul/roles/example-role
```

### Parameters for Consul version below 1.4

- `lease` `(string: "")` – Specifies the lease for this role. This is provided
  as a string duration with a time suffix like `"30s"` or `"1h"`. If not
  provided, the default Vault lease is used.

- `policies` `(string: <required>)` – Comma separated list of policies to be applied
  to the tokens.

### Sample payload

```json
{
  "policies": "global-management"
}
```

### Sample request

```sh
curl \
→     --request POST \
→     --header "X-Vault-Token: ..."\
→     --data @payload.json \
→     http://127.0.0.1:8200/v1/consul/roles/example-role
```

## Read Role

This endpoint queries for information about a Consul role with the given name.
If no role exists with that name, a 404 is returned.

| Method | Path                  |
| :----- | :-------------------- |
| `GET`  | `/consul/roles/:name` |

### Parameters

- `name` `(string: <required>)` – Specifies the name of the role to query. This
  is part of the request URL.

### Sample Request

```shell-session
$ curl \
    --header "X-Vault-Token: ..." \
    http://127.0.0.1:8200/v1/consul/roles/example-role
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

This endpoint lists all existing roles in the secrets engine.

| Method | Path            |
| :----- | :-------------- |
| `LIST` | `/consul/roles` |

### Sample Request

```shell-session
$ curl \
    --header "X-Vault-Token: ..." \
    --request LIST \
    http://127.0.0.1:8200/v1/consul/roles
```

### Sample Response

```json
{
  "data": {
    "keys": ["example-role"]
  }
}
```

## Delete Role

This endpoint deletes a Consul role with the given name. Even if the role does
not exist, this endpoint will still return a successful response.

| Method   | Path                  |
| :------- | :-------------------- |
| `DELETE` | `/consul/roles/:name` |

### Parameters

- `name` `(string: <required>)` – Specifies the name of the role to delete. This
  is part of the request URL.

### Sample Request

```shell-session
$ curl \
    --request DELETE \
    --header "X-Vault-Token: ..." \
    http://127.0.0.1:8200/v1/consul/roles/example-role
```

## Generate Credential

This endpoint generates a dynamic Consul token based on the given role
definition.

| Method | Path                  |
| :----- | :-------------------- |
| `GET`  | `/consul/creds/:name` |

### Parameters

- `name` `(string: <required>)` – Specifies the name of an existing role against
  which to create this Consul credential. This is part of the request URL.

### Sample Request

```shell-session
$ curl \
    --header "X-Vault-Token: ..." \
    http://127.0.0.1:8200/v1/consul/creds/example-role
```

### Sample Response

```json
{
  "data": {
    "token": "973a31ea-1ec4-c2de-0f63-623f477c2510"
  }
}
```
