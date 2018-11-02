---
layout: "api"
page_title: "JWT - Auth Methods - HTTP API"
sidebar_title: "JWT"
sidebar_current: "api-http-auth-jwt"
description: |-
  This is the API documentation for the Vault JWT authentication
  method plugin.
---

# JWT Auth Method (API)

This is the API documentation for the Vault JWT auth method
plugin. To learn more about the usage and operation, see the
[Vault JWT method documentation](/docs/auth/jwt.html).

This documentation assumes the plugin method is mounted at the
`/auth/jwt` path in Vault. Since it is possible to enable auth methods
at any location, please update your API calls accordingly.

## Configure

Configures the validation information to be used globally across all roles. One
(and only one) of `oidc_discovery_url` and `jwt_validation_pubkeys` must be
set.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `POST`   | `/auth/jwt/config`           | `204 (empty body)`     |

### Parameters

- `oidc_discovery_url` `(string: <optional>)` - The OIDC Discovery URL, without any .well-known component (base path). Cannot be used with `jwt_validation_pubkeys`.
- `oidc_discovery_ca_pem` `(string: <optional>)` - The CA certificate or chain of certificates, in PEM format, to use to validate connections to the OIDC Discovery URL. If not set, system certificates are used.
- `jwt_validation_pubkeys` `(comma-separated string, or array of strings: <optional>)` - A list of PEM-encoded public keys to use to authenticate signatures locally. Cannot be used with `oidc_discovery_url`.
- `bound_issuer` `(string: <optional>)` - The value against which to match the `iss` claim in a JWT.

### Sample Payload

```json
{
  "oidc_discovery_url": "https://myco.auth0.com/",
  "bound_issuer": "https://myco.auth0.com/"
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    https://127.0.0.1:8200/v1/auth/jwt/config
```

# Read Config

Returns the previously configured config.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `GET`    | `/auth/jwt/config`           | `200 application/json` |

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    https://127.0.0.1:8200/v1/auth/jwt/config
```

### Sample Response

```json
{
  "data":{
    "oidc_discovery_url": "https://myco.auth0.com/",
    "oidc_discovery_ca_pem": [],
    "bound_issuer": "https://myco.auth0.com/",
    "jwt_validation_pubkeys": []
  },
  ...
}
```

## Create Role

Registers a role in the method. Role types have specific entities
that can perform login operations against this endpoint. Constraints specific
to the role type must be set on the role. These are applied to the authenticated
entities attempting to login. At least one of the bound values must be set.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `POST`   | `/auth/jwt/role/:name`       | `204 (empty body)`     |

### Parameters
- `name` `(string: <required>)` - Name of the role.
- `bound_audiences` `(array: <required>)` - List of `aud` claims to match
  against. Any match is sufficient.
- `user_claim` `(string: <required>)` - The claim to use to uniquely identify
  the user; this will be used as the name for the Identity entity alias created
  due to a successful login. The claim value must be a string.
- `policies` `(array: <optional>)` - Policies to be set on tokens issued using
  this role.
- `ttl` `(int: <optional>)` - The initial/renewal TTL of tokens issued using
  this role, in seconds.
- `max_ttl` `(int: <optional>)` - The maximum allowed lifetime of tokens issued
  using this role, in seconds.
- `period` `(int: <optional>)` - If set, indicates that the token generated
  using this role should never expire, but instead always use the value set
  here as the TTL for every renewal.
- `num_uses` `(int: <optional>)` - If set, puts a use-count limitation on the
  issued token.
- `bound_subject` `(string: <optional>)` - If set, requires that the `sub`
  claim matches this value.
- `bound_cidrs` `(array: <optional>)` - If set, a list of CIDRs valid as the
  source address for login requests. This value is also encoded into any
  resulting token.
- `groups_claim` `(string: <optional>)` - The claim to use to uniquely identify
  the set of groups to which the user belongs; this will be used as the names
  for the Identity group aliases created due to a successful login. The claim
  value must be a list of strings.
- `groups_claim_delimiter_pattern` `(string: optional)` - A pattern of
  delimiters used to allow the `groups_claim` to live outside of the top-level
  JWT structure. For instance, a `groups_claim` of `meta/user.name/groups` with
  this field set to `//` will expect nested structures named `meta`,
  `user.name`, and `groups`. If this field was set to `/./` the groups
  information would expect to be via nested structures of `meta`, `user`,
  `name`, and `groups`.

### Sample Payload

```json
{
  "policies": [
    "dev",
    "prod"
  ],
  "bound_subject": "sl29dlldsfj3uECzsU3Sbmh0F29Fios1@clients",
  "bound_audiences": "https://myco.test",
  "user_claim": "https://vault/user",
  "groups_claim": "https://vault/groups"
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    https://127.0.0.1:8200/v1/auth/jwt/role/dev-role
```

## Read Role

Returns the previously registered role configuration.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `GET`   | `/auth/jwt/role/:name`        | `200 application/json` |

### Parameters

- `name` `(string: <required>)` - Name of the role.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    https://127.0.0.1:8200/v1/auth/jwt/role/dev-role
```

### Sample Response

```json
{
  "data":{
    "bound_subject": "sl29dlldsfj3uECzsU3Sbmh0F29Fios1@clients",
    "bound_audiences": [
      "https://myco.test"
    ],
    "bound_cidrs": [],
    "user_claim": "https://vault/user",
    "groups_claim": "https://vault/groups",
    "policies": [
      "dev",
      "prod"
    ],
    "period": 0,
    "ttl": 0,
    "num_uses": 0,
    "max_ttl": 0
  },
  ...
}

```

## List Roles

Lists all the roles that are registered with the plugin.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `LIST`   | `/auth/jwt/roles`            | `200 application/json` |

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request LIST \
    https://127.0.0.1:8200/v1/auth/jwt/roles
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
| `DELETE` | `/auth/jwt/role/:name`       | `204 (empty body)`     |

### Parameters

- `name` `(string: <required>)` - Name of the role.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request DELETE \
    https://127.0.0.1:8200/v1/auth/jwt/role/dev-role
```

## Login

Fetch a token. This endpoint takes a signed JSON Web Token (JWT) and
a role name for some entity. It verifies the JWT signature to authenticate that
entity and then authorizes the entity for the given role.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `POST`   | `/auth/jwt/login`            | `200 application/json` |

### Sample Payload

- `role` `(string: <required>)` - Name of the role against which the login is being
  attempted.
- `jwt` `(string: <required>)` - Signed [JSON Web Token](https://tools.ietf.org/html/rfc7519) (JWT).

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
    https://127.0.0.1:8200/v1/auth/jwt/login
```

### Sample Response

```json
{
    "auth":{
        "client_token":"f33f8c72-924e-11f8-cb43-ac59d697597c",
        "accessor":"0e9e354a-520f-df04-6867-ee81cae3d42d",
        "policies":[
            "default",
            "dev",
            "prod"
        ],
        "lease_duration":2764800,
        "renewable":true
    },
    ...
}
```
