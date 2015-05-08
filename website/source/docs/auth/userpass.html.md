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

## Authentication

#### Via the CLI

```
$ vault auth -method=userpass \
  -var="username=foo" \
  -var="password=bar"
...
```

#### Via the API

The endpoint for the login is `auth/userpass/login/USERNAME`.

The password should be sent, in the POST body encoded as JSON.

cURL example:
```
$ curl https://192.168.1.100:8200/v1/auth/userpass/login/USERNAME -XPOST
-d'{"password":"very_secure_password"}'
```

Response:
```
"lease_id":"","renewable":false,"lease_duration":0,"data":null,"auth":{"client_token":"really_long_token","policies":["root"],"metadata":{"username":"man"},"lease_duration":2592000,"renewable":false}}
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
Use `vault help` for more details.

```
$ vault write auth/userpass/users/mitchellh password=foo policies=root
...
```

The above creates a new user "mitchellh" with the password "foo" that
will be associated with the "root" policy. This is the only configuration
necessary.
