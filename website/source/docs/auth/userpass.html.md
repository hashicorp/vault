---
layout: "docs"
page_title: "Userpass - Auth Methods"
sidebar_current: "docs-auth-userpass"
description: |-
  The "userpass" auth method allows users to authenticate with Vault using a username and password.
---

# Userpass Auth Method

The `userpass` auth method allows users to authenticate with Vault using
a username and password combination.

The username/password combinations are configured directly to the auth
method using the `users/` path. This method cannot read usernames and
passwords from an external source.

The method lowercases all submitted usernames, e.g. `Mary` and `mary` are the
same entry.

## Authentication

### Via the CLI

```text
$ vault login -method=userpass \
    username=foo \
    password=bar
```

### Via the API

```shell
$ curl \
    --request POST \
    --data '{"password": "foo"}' \
    https://vault.rocks/v1/auth/userpass/login/mitchellh
```

The response will contain the token at `auth.client_token`:

```json
{
  "lease_id": "",
  "renewable": false,
  "lease_duration": 0,
  "data": null,
  "auth": {
    "client_token": "c4f280f6-fdb2-18eb-89d3-589e2e834cdb",
    "policies": [
      "admins"
    ],
    "metadata": {
      "username": "mitchellh"
    },
    "lease_duration": 0,
    "renewable": false
  }
}
```

## Configuration

Auth methods must be configured in advance before users or machines can
authenticate. These steps are usually completed by an operator or configuration
management tool.

1. Enable the userpass auth method:

    ```text
    $ vault auth enable userpass
    ```

1. Configure it with users that are allowed to authenticate:

    ```text
    $ vault write auth/userpass/users/mitchellh \
        password=foo \
        policies=admins
    ```

    This creates a new user "mitchellh" with the password "foo" that will be
    associated with the "admins" policy. This is the only configuration
    necessary.

## API

The Userpass auth method has a full HTTP API. Please see the
[Userpass auth method API](/api/auth/userpass/index.html) for more
details.
