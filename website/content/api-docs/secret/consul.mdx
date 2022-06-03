---
layout: api
page_title: Consul - Secrets Engines - HTTP API
description: This is the API documentation for the Vault Consul secrets engine.
---

# Consul Secrets Engine (API)

@include 'x509-sha1-deprecation.mdx'

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

- `token` `(string: "")` – Specifies the Consul ACL token to use. This
  must be a management type token. If this is not provided, Vault will try to
  bootstrap the ACL system of the Consul cluster.

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

### Parameters for Consul versions 1.11 and above

- `partition` `(string: "")` <EnterpriseAlert inline /> - Specifies the Consul admin partition the token
  will be generated within. The partition must exist, and the policies or roles assigned to the
  Vault role must also exist inside the given partition. If not provided, the "default"
  partition is used.

To create a client token within a particular Consul admin partition:

```json
{
  "partition": "admin1"
}
```

### Parameters for Consul versions 1.8 and above

- `consul_namespace` `(string: "")` <EnterpriseAlert inline /> - Specifies the Consul namespace the token
  will be generated within. The namespace must exist, and the policies or roles assigned to the
  Vault role must also exist inside the given Consul namespace. If not provided, the "default"
  namespace is used.

- `node_identities` `(list: <node identity or identities>)` - The list of node identities to
assign to the generated token. This may be a comma-separated list to attach multiple node identities
to a token.

To create a client token within a particular Consul namespace:

```json
{
  "consul_namespace": "ns1"
}
```

To create a client token with node identities attached:

```json
{
  "node_identities": "client-1:dc1,client-2:dc1"
}
```

### Parameters for Consul versions 1.5 and above

- `name` `(string: <required>)` – Specifies the name of an existing role against
  which to create this Consul credential. This is part of the request URL.

- `token_type` `(string: "client")` - Specifies the type of token to create when
  using this role. Valid values are `"client"` or `"management"`. If a `"management"`
  token, the `policy`, `policies`, and `consul_roles` parameters are not required.
  Defaults to `"client`".

- `policy` `(string: <policy>)` – Specifies the base64-encoded ACL policy. This is
  required unless the `token_type` is `"management"`. [Deprecated as of Consul 1.4 and
  removed as of Consul 1.11.](https://www.consul.io/api/acl/legacy)

- `policies` `(list: <policy or policies>)` – The list of policies to assign to the
  generated token. Either `policies` or `consul_roles` are required for Consul 1.5 and
  above, or just `policies` if using Consul 1.4.

- `consul_roles` `(list: <role or roles>)` – The list of Consul roles to assign to the
  generated token. Either `policies` or `consul_roles` are required for Consul 1.5 and above.

- `service_identities` `(list: <service identity or identities>)` - The list of
  service identities to assign to the generated token. This may be a semicolon-separated list to
  attach multiple service identities to a token.

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

To create a client token with defined policies:

```json
{
  "policies": "global-management,policy-2"
}
```

To create a client token with defined roles:

```json
{
  "consul_roles": "role-a,role-b"
}
```

To create a client token with service identities attached:

```json
{
  "service_identities": "myservice-1:dc1,dc2;myservice-2:dc1"
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

- `policy` `(string: <policy>)` – Specifies the base64-encoded ACL policy. The
  ACL format can be found in the [Consul ACL
  documentation](https://www.consul.io/docs/security/acl/acl-legacy). This is
  required unless the `token_type` is `"management"`.

### Sample payload

To create a client token with a custom base64-encoded policy:

```json
{
  "policy": "a2V5ICIi...=="
}
```

### Sample request

```shell-session
$ curl \
    --request POST \
    --header "X-Vault-Token: ..." \
    --data @payload.json \
    http://127.0.0.1:8200/v1/consul/roles/example-role
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
