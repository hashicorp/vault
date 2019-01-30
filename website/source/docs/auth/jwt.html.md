---
layout: "docs"
page_title: "JWT - Auth Methods"
sidebar_title: "JWT"
sidebar_current: "docs-auth-jwt"
description: |-
  The JWT auth method allows authentication using JWTs, with support for OIDC Discovery for key fetching
---

# JWT Auth Method

The `jwt` auth method can be used to authenticate with Vault using a JWT. This
JWT can be cryptographically verified using locally-provided keys, or, if
configured, an OIDC Discovery service can be used to fetch the appropriate
keys.

## Authentication

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


1. Enable the JWT auth method:

    ```text
    $ vault auth enable jwt
    ```

1. Use the `/config` endpoint to configure Vault with local keys or an OIDC Discovery URL. For the
list of available configuration options, please see the [API documentation](/api/auth/jwt/index.html).

    ```text
    $ vault write auth/jwt/config \
        oidc_discovery_url="https://myco.auth0.com/"
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
