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
- `oidc_client_id` `(string: <optional>)` - The OAuth Client ID from the provider for OIDC roles.
- `oidc_client_secret` `(string: <optional>)` - The OAuth Client Secret from the provider for OIDC roles.
- `jwt_validation_pubkeys` `(comma-separated string, or array of strings: <optional>)` - A list of PEM-encoded public keys to use to authenticate signatures locally. Cannot be used with `oidc_discovery_url`.
- `bound_issuer` `(string: <optional>)` - The value against which to match the `iss` claim in a JWT.
- `jwt_supported_algs` `(comma-separated string, or array of strings: <optional>)` - A list of supported signing algorithms. Defaults to [RS256]. ([Available algorithms](https://github.com/hashicorp/vault-plugin-auth-jwt/blob/master/vendor/github.com/coreos/go-oidc/jose.go#L7))
- `default_role` `(string: <optional>)` - The default role to use if none is provided during login.

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
- `role_type` `(string: <optional>)` - Type of role, either "oidc" (default) or "jwt".
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
- `bound_claims` `(map: <optional>)` - If set, a map of claims/values to match against.
- `groups_claim` `(string: <optional>)` - The claim to use to uniquely identify
  the set of groups to which the user belongs; this will be used as the names
  for the Identity group aliases created due to a successful login. The claim
  value must be a list of strings.
- `claim_mappings` `(map: <optional>)` - If set, a map of claims (keys) to be copied to
  specified metadata fields (values).
- `oidc_scopes` `(list: <optional>)` - If set, a list of OIDC scopes to be used with an OIDC role.
  The standard scope "openid" is automatically included and need not be specified.
- `allowed_redirect_uris` `(list: <required>)` - The list of allowed values for redirect_uri
  during OIDC logins.

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
  "groups_claim": "https://vault/groups",
  "bound_claims": {
    "department": "engineering",
    "sector": "7g"
  },
  "claim_mappings": {
    "preferred_language": "language",
    "group": "group"
  }
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
| `LIST`   | `/auth/jwt/role`            | `200 application/json` |

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request LIST \
    https://127.0.0.1:8200/v1/auth/jwt/role
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

## OIDC Authorization URL Request

Obtain an authorization URL from Vault to start an OIDC login flow.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `POST`   | `/auth/jwt/oidc/auth_url`    | `200 application/json` |

### Parameters

- `role` `(string: <optional>)` - Name of the role against which the login is being
  attempted. Defaults to configured `default_role` if not provided.
- `redirect_uri` `(string: <required>)` - Path to the callback to complete the login. This will be
  of the form, "https://.../oidc/callback" where the leading portion is dependent on your Vault
  server location, port, and the mount of the JWT plugin. This must be configured with Vault and the
  provider. See [Redirect URIs](/docs/auth/jwt.html#redirect-uris) for more information.

### Sample Payload

```json
{
    "role": "dev-role",
    "redirect_uri": "https://vault.myco.com:8200/vault/ui/auth/jwt/oidc/callback"
}
```

### Sample Request

```
$ curl \
    --request POST \
    --data @payload.json \
    https://127.0.0.1:8200/v1/auth/oidc/auth_url
```

### Sample Response

```json
{
  "request_id": "c701169c-64f8-26cc-0315-078e8c3ce897",
  "data": {
    "auth_url": "https://myco.auth0.com/authorize?client_id=r3qXcK2bezU3Sbmh0K16fatW6&nonce=851b69a9bfa5a6a5668111314414e3687891a599&redirect_uri=http%3A%2F%2Flocalhost%3A8300%2Foidc%2Fcallback&response_type=code&scope=openid+email+profile&state=1011e726d24960e09cfca2e04b36b38593cb6a22"
  },
  ...
}
```

## OIDC Callback
Exchange an authorization code for an OIDC ID Token. The ID token will be further validated
against any bound claims, and if valid a Vault token will be returned.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `GET`    | `/auth/jwt/oidc/callback`    | `200 application/json` |

### Parameters

- `state` `(string: <required>)` - Opaque state ID that is part of the Authorization URL and will
  be included in the the redirect following successful authenication on the provider.
- `nonce` `(string: <required>)` - Opaque nonce that is part of the Authorization URL and will
  be included in the the redirect following successful authenication on the provider.
- `code` `(string: <required>)` - Provider-generated authorization code that Vault will exchange for
  an ID token.

### Sample Request

```
$ curl \
    https://127.0.0.1:8200/v1/auth/jwt/oidc/callback?state=n2kfh3nsl&code=mn2ldl2nv98h2jl&nonce=ni42i2idj2jj
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

## JWT Login

Fetch a token. This endpoint takes a signed JSON Web Token (JWT) and
a role name for some entity. It verifies the JWT signature to authenticate that
entity and then authorizes the entity for the given role.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `POST`   | `/auth/jwt/login`            | `200 application/json` |

### Parameters

- `role` `(string: <optional>)` - Name of the role against which the login is being
  attempted. Defaults to configured `default_role` if not provided.
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
