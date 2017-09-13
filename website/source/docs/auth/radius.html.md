---
layout: "docs"
page_title: "RADIUS - Auth Methods"
sidebar_current: "docs-auth-radius"
description: |-
  The "radius" auth method allows users to authenticate with Vault using an
  existing RADIUS server.
---

# RADIUS Auth Method

The `radius` auth method allows users to authenticate with Vault using an
existing RADIUS server that accepts the PAP authentication scheme.

## Authentication

The default path is `/radius`. If this auth method was enabled at a different
path, specify `-path=/my-path` in the CLI.

### Via the CLI

```text
$ vault login -path=radius username=sethvargo
```

### Via the API

The default endpoint is `auth/radius/login`. If this auth method was enabled
at a different path, use that value instead of `radius`.

```shell
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data '{"password": "..."}' \
    https://vault.rocks/v1/auth/radius/login/sethvargo
```

The response will contain a token at `auth.client_token`:

```json
{
  "auth": {
    "client_token": "c4f280f6-fdb2-18eb-89d3-589e2e834cdb",
    "policies": [
      "admins"
    ],
    "metadata": {
      "username": "mitchellh"
    }
  }
}
```

## Configuration

### Via the CLI

1. Enable the radius auth method:

    ```text
    $ vault auth enable radius
    ```

1. Configure connection details for your RADIUS server.

    ```text
    $ vault write auth/radius/users/mitchellh policies=admins
    ```

    For the complete list of configuration options, please see the API
    documentation.

    The above creates a new mapping for user "mitchellh" that will be associated
    with the "admins" policy.

    Alternatively, Vault can assign a configurable set of policies to any user
    that successfully authenticates with the RADIUS server but has no explicit
    mapping in the `users/` path. This is done through the
    `unregistered_user_policies` configuration parameter.

## API

The RADIUS auth method has a full HTTP API. Please see the
[RADIUS Auth API](/api/auth/radius/index.html) for more
details.
