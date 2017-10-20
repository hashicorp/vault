---
layout: "api"
page_title: "Kubernetes Auth Plugin Backend - HTTP API"
sidebar_current: "docs-http-auth-kubernetes"
description: |-
  This is the API documentation for the Vault Kubernetes authentication
  backend plugin.
---

# Kubernetes Auth Plugin HTTP API

This is the API documentation for the Vault Kubernetes authentication backend
plugin. To learn more about the usage and operation, see the
[Vault Kubernetes backend documentation](/docs/auth/kubernetes.html).

This documentation assumes the backend is mounted at the
`/auth/kubernetes` path in Vault. Since it is possible to mount auth backends
at any location, please update your API calls accordingly.

## Configure

The Kubernetes Auth backend validates service account JWTs and verifies their
existence with the Kubernetes TokenReview API. This endpoint configures the
public key used to validate the JWT signature and the necessary information to
access the Kubernetes API.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `POST`   | `/auth/kubernetes/config`    | `204 (empty body)`     |

### Parameters
 - `kubernetes_host` `(string: <required>)` - Host must be a host string, a host:port pair, or a URL to the base of the Kubernetes API server.
 - `kubernetes_ca_cert` `(string: "")` - PEM encoded CA cert for use by the TLS client used to talk with the Kubernetes API.
 - `pem_keys` `(array: [])` - Optional list of PEM-formated public keys or certificates
    used to verify the signatures of Kubernetes service account
    JWTs. If a certificate is given, its public key will be
    extracted. Not every installation of Kubernetes exposes these
    keys. 

### Sample Payload

```json
{
  "kubernetes_host": "https://192.168.99.100:8443",
  "kubernetes_ca_cert": "-----BEGIN CERTIFICATE-----.....-----END CERTIFICATE-----",
  "pem_keys": "-----BEGIN CERTIFICATE-----.....-----END CERTIFICATE-----"
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    https://vault.rocks/v1/auth/kubernetes/config
```

## Read Config

Returns the previously configured config, including credentials.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `GET`    | `/auth/kubernetes/config`    | `200 application/json` |

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    https://vault.rocks/v1/auth/kubernetes/config
```

### Sample Response

```json
{
  "data":{
      "kubernetes_host": "https://192.168.99.100:8443",
      "kubernetes_ca_cert": "-----BEGIN CERTIFICATE-----.....-----END CERTIFICATE-----",
      "pem_keys": "-----BEGIN CERTIFICATE-----.....-----END CERTIFICATE-----"
  },
  ...
}

```

## Create Role

Registers a role in the backend. Role types have specific entities
that can perform login operations against this endpoint. Constraints specific
to the role type must be set on the role. These are applied to the authenticated
entities attempting to login.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `POST`   | `/auth/kubernetes/role/:name`| `204 (empty body)`     |

### Parameters
- `name` `(string: <required>)` - Name of the role.
- `bound_service_account_names` `(array: <required>)` - List of service account
  names able to access this role. If set to "\*" all names are allowed, both this
  and bound_service_account_namespaces can not be "\*".
- `bound_service_account_namespaces` `(array: <required>)` - List of namespaces
  allowed to access this role. If set to "\*" all namespaces are allowed, both
  this and bound_service_account_names can not be set to "\*".
- `ttl` `(string: "")` - The TTL period of tokens issued using this role in
  seconds.
- `max_ttl` `(string: "")` - The maximum allowed lifetime of tokens
  issued in seconds using this role.
- `period` `(string: "")` - If set, indicates that the token generated using
  this role should never expire. The token should be renewed within the duration
  specified by this value. At each renewal, the token's TTL will be set to the
  value of this parameter.
- `policies` `(array: [])` - Policies to be set on tokens issued using this
  role.

### Sample Payload

```json
{
  "bound_service_account_names": "vault-auth",
  "bound_service_account_namespaces": "default",
  "policies": [
    "dev",
    "prod"
  ],
  "max_ttl": 1800000,
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    https://vault.rocks/v1/auth/kubernetes/role/dev-role
```
## Read Role

Returns the previously registered role configuration.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `GET`   | `/auth/kubernetes/role/:name` | `200 application/json` |

### Parameters

- `name` `(string: <required>)` - Name of the role.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    https://vault.rocks/v1/auth/kubernetes/role/dev-role
```

### Sample Response

```json
{
    "data":{
        "bound_service_account_names": "vault-auth",
        "bound_service_account_namespaces": "default",
        "max_ttl": 1800000,,
        "ttl":0,
        "period": 0,
        "policies":[
            "dev",
            "prod"
        ],
    },
    ...
}

```

## List Roles

Lists all the roles that are registered with the backend.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `LIST`   | `/auth/kubernetes/roles`            | `200 application/json` |
| `GET`   | `/auth/kubernetes/roles?list=true`   | `200 application/json` |

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request LIST \
    https://vault.rocks/v1/auth/kubernetes/roles
```

### Sample Response

```json  
{
    "data": {
        "keys": [
            "dev-role",
            "prod-role"
        ]
    },
    ...
}
```

## Delete Role

Deletes the previously registered role.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `DELETE` | `/auth/kubernetes/role/:role`| `204 (empty body)`     |

### Parameters

- `role` `(string: <required>)` - Name of the role.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request DELETE \
    https://vault.rocks/v1/auth/kubernetes/role/dev-role
```

## Login

Fetch a token. This endpoint takes a signed JSON Web Token (JWT) and
a role name for some entity. It verifies the JWT signature to authenticate that
entity and then authorizes the entity for the given role.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `POST`   | `/auth/kubernetes/login`            | `200 application/json` |

### Sample Payload

- `role` `(string: <required>)` - Name of the role against which the login is being
  attempted.
- `jwt` `(string: <required>)` - Signed [JSON Web
  Token](https://tools.ietf.org/html/rfc7519) (JWT) for authenticating a service
  account. 

### Sample Payload

```json
{
    "role": "dev-role",
    "jwt": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### Sample Request

```
$ curl \
    --request POST \
    --data @payload.json \
    https://vault.rocks/v1/auth/kubernetes/login
```

### Sample Response

```json
{
	"auth": {
		"client_token": "62b858f9-529c-6b26-e0b8-0457b6aacdb4",
		"accessor": "afa306d0-be3d-c8d2-b0d7-2676e1c0d9b4",
		"policies": [
			"default"
		],
		"metadata": {
			"role": "test",
			"service_account_name": "vault-auth",
			"service_account_namespace": "default",
			"service_account_secret_name": "vault-auth-token-pd21c",
			"service_account_uid": "aa9aa8ff-98d0-11e7-9bb7-0800276d99bf"
		},
		"lease_duration": 2764800,
		"renewable": true
	}
    ...
}
```
