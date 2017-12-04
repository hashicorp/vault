---
layout: "docs"
page_title: "Auth Backend: RADIUS"
sidebar_current: "docs-auth-radius"
description: |-
  The "radius" auth backend allows users to authenticate with Vault using an existing RADIUS server.
---

# Auth Backend: RADIUS

Name: `radius`

The "radius" auth backend allows users to authenticate with Vault using
an existing RADIUS server that accepts the PAP authentication scheme. 

The mapping of users to Vault policies is managed by using the
`users/` path.

Optionally, a configurable set of policies can be granted to all users that 
can successfully authenticate but are not registered in the `users/` path.

## Authentication

#### Via the CLI

```
$ vault auth -method=userpass -path=radius \
    username=foo \
    password=bar
```

#### Via the API

The endpoint for the login is `auth/radius/login/<username>`.

The password should be sent in the POST body encoded as JSON.

```shell
$ curl $VAULT_ADDR/v1/auth/radius/login/mitchellh \
    -d '{ "password": "foo" }'
```

Alternatively a POST request can be made to `auth/radius/login/` 
with both `username` and `password` sent in the POST body encoded as JSON.

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

First, you must enable the RADIUS auth backend:

```
$ vault auth-enable radius
Successfully enabled 'radius' at 'radius'!
```

Now when you run `vault auth -methods`, the RADIUS backend is
available:

```
Path       Type      Description
token/     token     token based credentials
radius/    radius
```

To use the radius auth backend, it must first be configured with connection
details for your RADIUS server.
The configuration options are detailed below in the API docs.
Configuration is written to `auth/radius/config`.

To use the "radius" auth backend, an operator must configure a
mapping between users and policies. An example is shown below.
Use `vault path-help` for more details.

```
$ vault write auth/radius/users/mitchellh \
    policies=admins
...
```

The above creates a new mapping for user "mitchellh" that 
will be associated with the "admins" policy.

Alternatively, Vault can assign a configurable set of policies to 
any user that successfully authenticates with the RADIUS server but 
has no explicit mapping in the `users/` path.
This is done through the `unregistered_user_policies` configuration parameter.

## API

The RADIUS authentication backend has a full HTTP API. Please see the
[RADIUS Auth API](/api/auth/radius/index.html) for more
details.

