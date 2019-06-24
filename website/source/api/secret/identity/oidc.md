---
layout: "api"
page_title: "Identity Secret Backend: OIDC - HTTP API"
sidebar_title: "OIDC"
sidebar_current: "api-http-secret-identity-oidc"
description: |-
  This is the API documentation for configuring and acquiring vault issued OIDC tokens.
---

## Configure OIDC Backend

This endpoint updates vault OIDC configurations.

| Method   | Path                |
| :------------------ | :----------------------|
| `POST`   | `identity/oidc/config`  |

### Parameters

- `issuer` `(string: "")` – Issuer URL to be used in the iss claim of the token. If not set, Vault's app_addr will be used. The issuer is a case sensitive URL using the https scheme that contains scheme, host, and optionally, port number and path components and no query or fragment components.

### Sample Payload

```json
{
  "issuer": "https://example.com:1234"
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    http://127.0.0.1:8200/v1/identity/oidc/config
```

### Sample Response

```json
{
  "data": null,
  "warnings": [
    "If \"issuer\" is set explicitly, all tokens must be validated against that address, including those issued by secondary clusters. Setting issuer to \"\" will restore the default behavior of using the cluster's api_addr as the issuer."
  ],
}
```

## Read OIDC Backend Configurations

This endpoint queries vault OIDC configurations.

| Method   | Path                |
| :------------------ | :----------------------|
| `GET`   | `identity/oidc/config`  |

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request GET \
    http://127.0.0.1:8200/v1/identity/oidc/config
```

### Sample Response

```json
{
  "data": {
    "issuer": "https://example.com:1234"
  },
}
```

## Create a Signing Key

This endpoint creates or updates a Key which is used to sign tokens.

| Method   | Path                |
| :------------------ | :----------------------|
| `POST`   | `/oidc/key/:name`  |

### Parameters

- `name` `(string)` – Name of the key.

- `rotation_period` `(int or time string: "24h")` - How often to generate a new keypair. Can be specified as a number of seconds or a time string like "30m" or "6h".
  corresponding existing entity.

- `verification_ttl` `(int or time string: "24h")` - Controls how long the public portion of a key will be available for verification after being rotated.

- "algorithm" `(string: "RS256")` - Signing algorithm to use. This will default to RS256, and is currently the only allowed value.

### Sample Payload

```json
{
  "rotation_period":"12h",
  "verification_ttl":43200,
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    http://127.0.0.1:8200/v1/identity/entity/key/signing-key-001
```

## Read a Signing Key

This endpoint queries a signing Key.

| Method   | Path                |
| :------------------ | :----------------------|
| `GET`   | `/oidc/key/:name`  |

### Parameters

- `name` `(string)` – Name of the key.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request GET \
    http://127.0.0.1:8200/v1/identity/entity/key/signing-key-001
```

### Sample Response

```json
{
  "data": {
    "algorithm": "RS256",
    "rotation_period": 43200,
    "verification_ttl": 43200
  },
}
```

## Delete a Signing Key

This endpoint deletes a signing key.

| Method   | Path                |
| :------------------ | :----------------------|
| `DELETE`   | `/oidc/role/:name`  |

### Parameters

- `name` `(string)` – Name of the key.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request DELETE \
    http://127.0.0.1:8200/v1/identity/oidc/key/signing-key-001
```

## Rotate a Signing Key

This endpoint rotates a Signing Key.

| Method   | Path                |
| :------------------ | :----------------------|
| `POST`   | `identity/oidc/key/:name/rotate`  |

### Parameters

- `name` `(string)` – Name of the key to be rotated.

- `verification_ttl` `(string: <optional>)` - Controls how long the public portion of the key will be available for verification after being rotated. Setting verification_ttl here will override the verification_ttl set on the key."

### Sample Payload

```json
{
  "verification_ttl": 0
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    http://127.0.0.1:8200/v1/identity/oidc/key/signing-key-001/rotate
```





## List Signing Keys

This endpoint will List all Signing Keys.

| Method   | Path                |
| :------------------ | :----------------------|
| `LIST`   | `identity/oidc/key`  |

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request LIST \
    http://127.0.0.1:8200/v1/identity/oidc/key
```

### Sample Response

```json
{
  "data": {
    "keys": [
      "signing-key-001",
      "signing-key-002"
    ]
  },
}
```

## Read .well-known Configurations

Query this path to retrieve the configured OIDC Issuer, Keys endpoint, Subjects, and signing algorithms used by the OIDC backend.

| Method   | Path                |
| :------------------ | :----------------------|
| `GET`   | `identity/oidc/.well-known/openid-configuration`  |

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request LIST \
    http://127.0.0.1:8200/v1/identity/oidc/.well-known/openid-configuration
```

### Sample Response

```json
{
  "issuer": "https://example.com:1234",
  "authorization_endpoint": "",
  "token_endpoint": "",
  "jwks_uri": "https://example.com:1234/.well-known/keys",
  "response_types_supported": null,
  "subject_types_supported": [
    "public"
  ],
  "id_token_signing_alg_values_supported": [
    "RS256"
  ],
  "scopes_supported": null,
  "token_endpoint_auth_methods_supported": null,
  "claims_supported": null
}
```

## Read Active Public Keys

Query this path to retrieve the public portion of signing keys which are able to verify signed tokens. Clients can use this to validate the authenticity of an OIDC token.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request LIST \
    http://127.0.0.1:8200/v1/identity/oidc/.well-known/keys
```

### Sample Response

```json
{
  "keys": [
    {
      "use": "sig",
      "kty": "RSA",
      "kid": "94178020-55b5-e18d-b32b-1010ba5a35b4",
      "alg": "RS256",
      "n": "1bt-V8T7g0zr7koNbdppFrUM5YrnybPDOt-cK3MKmL1FcN3aOltCw9tCYStHgm8mIz_DJ1HgIjA-DcK_O9gacEGFCidUuudV8O4TixToHEVyRe1yXu-Q98hwkm9JtFF9PvMzDXhn4s3bLanOZzO15JAdVCo0JnwSIT9Ay3LxPLbWHYbPj7ROScuvic99OyvWz87qBK-AoXmxo9lRNY39LtieMr1D2iq0HvtjHkfiarr34CSTcuksknOsY49BU5ktrs_YJSEVpeRQ8RywY1sWrq8w_UmGsNFfPr--crXQw0ekJCXzmotsRHE5jwMuhjuucVlnyQFBYEdfDB_iPbC7Hw",
      "e": "AQAB"
    }
  ]
}
```

## Create or Update a Role

Create or update a role. OIDC tokens are generated against a role.

| Method   | Path                |
| :------------------ | :----------------------|
| `POST`   | `identity/oidc/role/:name`  |

### Parameters

- `name` `(string)` – Name of the role.

- `key` `(string)` – A configured Signing Key

-	`template` `(string: <optional>)` - The template string to use for generating tokens. This may be in string-ified JSON or base64 format.

- `ttl` `(int or time string: "24h") - TTL of the tokens generated against the role. Can be specified as a number of seconds or a time string like "30m" or "6h".

### Sample Payload

```json
{
  "key": "signing-key-001",
  "ttl":"12h"
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    http://127.0.0.1:8200/v1/identity/oidc/role/role-001
```

## Read a Role

This endpoint queries a Role.

| Method   | Path                |
| :------------------ | :----------------------|
| `GET`   | `/oidc/role/:name`  |

### Parameters

- `name` `(string)` – Name of the key.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request GET \
    http://127.0.0.1:8200/v1/identity/role/role-001
```

### Sample Response

```json
{
  "data": {
    "client_id": "PGE8tf4RmJkDwvjI1FgARkXEmH",
    "key": "signing-key-001",
    "template": "",
    "ttl": 43200
  },
}
```

## List Roles

This endpoint will List all Signing Keys.

| Method   | Path                |
| :------------------ | :----------------------|
| `LIST`   | `identity/oidc/role`  |

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request LIST \
    http://127.0.0.1:8200/v1/identity/oidc/role
```

### Sample Response

```json
{
  "data": {
    "keys": [
      "role-001",
      "role-002",
      "testrole"
    ]
  },
}
```

## Delete a Role

This endpoint deletes a role.

| Method   | Path                |
| :------------------ | :----------------------|
| `DELETE`   | `/oidc/role/:name`  |

### Parameters

- `name` `(string)` – Name of the key.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request DELETE \
    http://127.0.0.1:8200/v1/identity/oidc/role/role-001
```

## Generate a Signed ID Token

Use this endpoint to generate a signed ID (OIDC) token.

| Method   | Path                |
| :------------------ | :----------------------|
| `POST`   | `identity/oidc/token/:name`  |

### Parameters

- `name` `(string: "")` – The name of the role against which to generate a signed ID token

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request GET \
    --data @payload.json \
    http://127.0.0.1:8200/v1/identity/oidc/token/role-001
```

### Sample Response

```json
{
  "data": {
    "client_id": "P6CfCzyHsQY4pMcA6kWAOCItA7",
    "token": "eyJhbGciOiJSUzI1NiIsImtpZCI6IjJkMGI4YjlkLWYwNGQtNzFlYy1iNjc0LWM3MzU4NDMyYmM1YiJ9.eyJhdWQiOiJQNkNmQ3p5SHNRWTRwTWNBNmtXQU9DSXRBNyIsImV4cCI6MTU2MTQ4ODQxMiwiaWF0IjoxNTYxNDAyMDEyLCJpc3MiOiJodHRwczovL2V4YW1wbGUuY29tOjEyMzQiLCJzdWIiOiI2YzY1ZWFmNy1kNGY0LTEzMzMtMDJiYy0xYzc1MjE5YzMxMDIifQ.IcbWTmks7P5eVtwmIBl5rL1B88MI55a9JJuYVLIlwE9aP_ilXpX5fE38CDm5PixDDVJb8TI2Q_FO4GMMH0ymHDO25ZvA917WcyHCSBGaQlgcS-WUL2fYTqFjSh-pezszaYBgPuGvH7hJjlTZO6g0LPCyUWat3zcRIjIQdXZum-OyhWAelQlveEL8sOG_ldyZ8v7fy7GXDxfJOK1kpw5AX9DXJKylbwZTBS8tLb-7edq8uZ0lNQyWy9VPEW_EEIZvGWy0AHua-Loa2l59GRRP8mPxuMYxH_c88x1lsSw0vH9E3rU8AXLyF3n4d40PASXEjZ-7dnIf4w4hf2P4L0xs_g",
    "ttl": 86400
  },
}
```


















		{
			Pattern: "oidc/role/?$",
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.ListOperation: i.pathOIDCListRole,
			},
			HelpSynopsis:    "List configured OIDC roles",
			HelpDescription: "List all configured OIDC roles in the identity backend.",
		},
		{
			Pattern: "oidc/introspect/?$",
			Fields: map[string]*framework.FieldSchema{
				"token": {
					Type:        framework.TypeString,
					Description: "Token to verify",
				},
				"client_id": {
					Type:        framework.TypeString,
					Description: "Optional client_id to verify",
				},
			},
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.UpdateOperation: i.pathOIDCIntrospect,
			},
			HelpSynopsis:    "Verify the authenticity of an OIDC token",
			HelpDescription: "Use this path to verify the authenticity of an OIDC token and whether the associated entity is active and enabled.",
		},
	}
}