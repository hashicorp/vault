---
layout: "docs"
page_title: "Auth Backend: Username & Password"
sidebar_current: "docs-auth-userpass"
description: |-
  The "userpass" auth backend allows users to authenticate with Vault using a username and password.
---

# Auth Backend: Username & Password

Name: `userpass`

The "userpass" auth backend allows users to authenticate with Vault using
a username and password combination.

The username/password combinations are configured directly to the auth
backend using the `users/` path. This backend cannot read usernames and
passwords from an external source.

The backend lowercases all submitted usernames, e.g. `Mary` and `mary` are the
same entry.

## Authentication

#### Via the CLI

```
$ vault auth -method=userpass \
    username=foo \
    password=bar
```

#### Via the API

The endpoint for the login is `auth/userpass/login/<username>`.

The password should be sent in the POST body encoded as JSON.

```shell
$ curl $VAULT_ADDR/v1/auth/userpass/login/mitchellh \
    -d '{ "password": "foo" }'
```

The response will be in JSON. For example:

```javascript
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

First, you must enable the username/password auth backend:

```
$ vault auth-enable userpass
Successfully enabled 'userpass' at 'userpass'!
```

Now when you run `vault auth -methods`, the username/password backend is
available:

```
Path       Type      Description
token/     token     token based credentials
userpass/  userpass
```

To use the "userpass" auth backend, an operator must configure it with
users that are allowed to authenticate. An example is shown below.
Use `vault path-help` for more details.

```
$ vault write auth/userpass/users/mitchellh \
    password=foo \
    policies=admins
...
```

The above creates a new user "mitchellh" with the password "foo" that
will be associated with the "admins" policy. This is the only configuration
necessary.

## API

The Username & Password authentication backend has a full HTTP API. Please see the
[Userpass auth backend API](/api/auth/userpass/index.html) for more
details.

