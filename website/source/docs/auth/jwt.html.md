---
layout: "docs"
page_title: "JWT - Auth Methods"
sidebar_title: "JWT"
sidebar_current: "docs-auth-jwt"
description: |-
  The JWT auth method allows authentication using OIDC and user-provided JWTs
---

# JWT Auth Method

The `jwt` auth method can be used to authenticate with Vault using
[OIDC](https://en.wikipedia.org/wiki/OpenID_Connect) or by providing a [JWT](https://en.wikipedia.org/wiki/JSON_Web_Token).
The OIDC method will allow authentication via a configured OIDC provider using the user's web browser.
This method may be initiated from the Vault UI or the command line. Alternatively, a JWT can be provided
directly. It will be cryptographically verified using locally-provided keys, or, if configured, an
OIDC Discovery service can be used to fetch the appropriate keys. The choice of method is configured per role.

Both methods allow additional processing of the claims data in the JWT. Some of the concepts common
to both methods will be covered first, followed by specific examples of OIDC and JWT usage.

### Bound Claims

Once a JWT has been validated as being properly signed and not expired, the authorization flow will validate that any
configured "bound" parameters match. In some cases there are dedicated parameters, for example `bound_subject`,  which must
match the JWT's `sub` parameter. A role may also be configured to check arbitrary claims through the `bound_claims`
map. The map contains a set of claims and their required values. For example, assume `bound_claims` is set to:

```json
{
  "division": "Europe",
  "department": "Engineering"
}
```

Only JWTs containing both the "division" and "department" claims, and respective matching values of "Europe" and "Engineering",
would be authorized.

### Claims as Metadata

Data from claims can be copied into the resulting auth token and alias metadata by configuring `claim_mappings`. This role
parameter is a map of items to copy. The map elements are of the form: `"<JWT claim>":"<metadata key>"`. Assume
`claim_mappings` is set to:

```json
{
  "division": "organization",
  "department": "department"
}
```

This specifies that the value in the JWT claim "division" should be copied to the metadata key "organization". The JWT
"department" claim value will also be copied into metadata but will retain the key name. If a claim is configured in `claim_mappings`,
it must existing in the JWT or else the authentication will fail.

Note: the metadata key name "role" is reserved and may not be used for claim mappings.


### Claim specifications and JSON Pointer

Some parameters (e.g. `bound_claims` and `groups_claim`) are used to point to data within the JWT. If
the desired key is at the top of level of the JWT, the name can be provided directly. If it is nested at a
lower level, a JSON Pointer may be used.

Assume the following JSON data to be referenced:

```json
{
  "division": "North America",
  "groups": {
    "primary": "Engineering",
    "secondary": "Software"
  }
}
```

A parameter of `"division"` will reference "North America", as this is a top level key. A parameter
`"/groups/primary"` uses JSON Pointer syntax to reference "Engineering" at a lower level. Any valid
JSON Pointer can be used as a selector. Refer to the
[JSON Pointer RFC](https://tools.ietf.org/html/rfc6901) for a full description of the syntax


## OIDC Authentication

This section covers the setup and use of OIDC roles. If a JWT is to be provided directly,
refer to the [JWT Authentication](/docs/auth/jwt.html#jwt-authentication) section below. Basic
familiarity with [OIDC concepts] (https://developer.okta.com/blog/2017/07/25/oidc-primer-part-1)
is assumed.

Vault includes two built-in OIDC login flows: the Vault UI, and the CLI
using a `vault login`.

### Redirect URIs

An important part of OIDC role configuration is properly setting redirect URIs. This must be
done both in Vault and with the OIDC provider, and these configurations must align. The
redirect URIs are specified for a role with the `allowed_redirect_uris` parameter. There are
different redirect URIs to configure the Vault UI and CLI flows, so one or both will need to
be set up depending on the installation.

**CLI**

If you plan to support authentication via `vault login -method=oidc`, a localhost redirect URI
must be set. This can usually be: `http://localhost:8250/oidc/callback`. Logins via the CLI may
specify a different listening port if needed, and a URI with this port must match one of the
configured redirected URIs. These same "localhost" URIs must be added to the provider as well.

**Vault UI**

Logging in via the Vault UI requires a redirect URI of the form:
`https://{host:port}/ui/vault/auth/{path}/oidc/callback`

The "host:port" must be correct for the Vault server, and "path" must match the path the JWT
backend is mounted at (e.g. "oidc" or "jwt").
If [namespaces](https://www.vaultproject.io/docs/enterprise/namespaces/index.html) are being used,
they must be added as query parameters, for example:

`https://vault.example.com:8200/ui/vault/auth/oidc/oidc/callback?namespace=my_ns`

### OIDC Login (Vault UI)

1. Select the "OIDC" login method.
1. Enter a role name if necessary.
1. Press "Sign In" and complete the authentication with the configured provider.

### OIDC Login (CLI)

The CLI login defaults to path of `/oidc`. If this auth method was enabled at a
different path, specify `-path=/my-path` in the CLI. The default local listening port is 8250. This
can be changed with the `port` option.

```text
$ vault login -method=oidc port=8400 role=test

Complete the login via your OIDC provider. Launching browser to:

    https://myco.auth0.com/authorize?redirect_uri=http%3A%2F%2Flocalhost%3A8400%2Foidc%2Fcallback&client_id=r3qXc2bix9eF...
```

The browser will open to the generated URL to complete the provider's login. The
URL may be entered manually if the browser cannot be automatically opened.

### OIDC Provider Configuration

The OIDC authentication flow has been successfully tested with a number of providers. A full
guide to configuring OAuth/OIDC applications is beyond the scope of Vault documentation, but a
collection of provider configuration steps has been collected to help get started:
[OIDC Provider Setup](/docs/auth/jwt_oidc_providers.html)

## JWT Authentication

The authentication flow for roles of type "jwt" is simpler than OIDC since Vault
only needs to validate the provided JWT.

### Via the CLI

The default path is `/jwt`. If this auth method was enabled at a
different path, specify `-path=/my-path` in the CLI.

```text
$ vault write auth/jwt/login role=demo jwt=...
```

### Via the API

The default endpoint is `auth/jwt/login`. If this auth method was enabled
at a different path, use that value instead of `jwt`.

```shell
$ curl \
    --request POST \
    --data '{"jwt": "your_jwt", "role": "demo"}' \
    http://127.0.0.1:8200/v1/auth/jwt/login
```

The response will contain a token at `auth.client_token`:

```json
{
  "auth": {
    "client_token": "38fe9691-e623-7238-f618-c94d4e7bc674",
    "accessor": "78e87a38-84ed-2692-538f-ca8b9f400ab3",
    "policies": [
      "default"
    ],
    "metadata": {
      "role": "demo"
    },
    "lease_duration": 2764800,
    "renewable": true
  }
}
```

## Configuration

Auth methods must be configured in advance before users or machines can
authenticate. These steps are usually completed by an operator or configuration
management tool.


1. Enable the JWT auth method. Either the "jwt" or "oidc" name may be used. The
backend will be mounted at the chosen name.

    ```text
    $ vault auth enable jwt
     or
    $ vault auth enable oidc
    ```

1. Use the `/config` endpoint to configure Vault. To support JWT roles, either local keys or an OIDC
Discovery URL must be present. For OIDC roles, OIDC Discovery URL, OIDC Client ID and OIDC Client Secret are required. For the
list of available configuration options, please see the [API documentation](/api/auth/jwt/index.html).

    ```text
    $ vault write auth/jwt/config \
        oidc_discovery_url="https://myco.auth0.com/"
        oidc_client_id="m5i8bj3iofytj",
        oidc_client_secret="f4ubv72nfiu23hnsj",
        default_role="demo"
    ```

1. Create a named role:

    ```text
    vault write auth/jwt/role/demo \
        bound_subject="r3qX9DljwFIWhsiqwFiu38209F10atW6@clients" \
        bound_audiences="https://vault.plugin.auth.jwt.test" \
        user_claim="https://vault/user" \
        groups_claim="https://vault/groups" \
        policies=webapps \
        ttl=1h
    ```

    This role authorizes JWTs with the given subject and audience claims, gives
    it the `webapps` policy, and uses the given user/groups claims to set up
    Identity aliases.

    For the complete list of configuration options, please see the API
    documentation.

## API

The JWT Auth Plugin has a full HTTP API. Please see the
[API docs](/api/auth/jwt/index.html) for more details.
