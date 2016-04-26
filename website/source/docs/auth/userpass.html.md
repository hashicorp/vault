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
      "root"
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
    policies=root
...
```

The above creates a new user "mitchellh" with the password "foo" that
will be associated with the "root" policy. This is the only configuration
necessary.

## API

### /auth/userpass/users/[username]
#### POST

<dl class="api">
  <dt>Description</dt>
  <dd>
      Create a new user or update an existing user.
      This path honors the distinction between the `create` and `update` capabilities inside ACL policies.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/auth/userpass/users/<username>`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">username</span>
        <span class="param-flags">required</span>
            Username for this user.
      </li>
    </ul>
  </dd>
  <dd>
    <ul>
      <li>
        <span class="param">password</span>
        <span class="param-flags">required</span>
            Password for this user.
      </li>
    </ul>
  </dd>
  <dd>
    <ul>
      <li>
        <span class="param">policies</span>
        <span class="param-flags">optional</span>
            Comma-separated list of policies.
            If set to empty string, only the `default` policy will be applicable to the user.
      </li>
    </ul>
  </dd>
  <dd>
    <ul>
      <li>
        <span class="param">ttl</span>
        <span class="param-flags">optional</span>
            The lease duration which decides login expiration.
      </li>
    </ul>
  </dd>
  <dd>
    <ul>
      <li>
        <span class="param">max_ttl</span>
        <span class="param-flags">optional</span>
            Maximum duration after which login should expire.
      </li>
    </ul>
  </dd>

  <dt>Returns</dt>
  <dd>`204` response code.
  </dd>
</dl>

### /auth/userpass/users/[username]/password
#### POST
<dl class="api">
  <dt>Description</dt>
  <dd>
      Update the password for an existing user.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/auth/userpass/users/<username>/password`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">username</span>
        <span class="param-flags">required</span>
            Username for this user.
      </li>
    </ul>
  </dd>
  <dd>
    <ul>
      <li>
        <span class="param">password</span>
        <span class="param-flags">required</span>
            Password for this user.
      </li>
    </ul>
  </dd>

  <dt>Returns</dt>
  <dd>`204` response code.
  </dd>
</dl>

### /auth/userpass/users/[username]/policies
#### POST
<dl class="api">
  <dt>Description</dt>
  <dd>
      Update the policies associated with an existing user.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/auth/userpass/users/<username>/policies`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">username</span>
        <span class="param-flags">required</span>
            Username for this user.
      </li>
    </ul>
  </dd>
  <dd>
    <ul>
      <li>
        <span class="param">policies</span>
        <span class="param-flags">optional</span>
            Comma-separated list of policies.
            If this is field is not supplied, the policies will be unchanged.
            If set to empty string, only the `default` policy will be applicable to the user.
      </li>
    </ul>
  </dd>

  <dt>Returns</dt>
  <dd>`204` response code.
  </dd>
</dl>


### /auth/userpass/login/[username]
#### POST
<dl class="api">
  <dt>Description</dt>
  <dd>
      Login with the username and password.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/auth/userpass/login/<username>`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">password</span>
        <span class="param-flags">required</span>
            Password for this user.
      </li>
    </ul>
  </dd>

  <dt>Returns</dt>
  <dd>

   ```javascript
   {
	"lease_id": "",
	"renewable": false,
	"lease_duration": 0,
	"data": null,
	"warnings": null,
	"auth": {
		"client_token": "64d2a8f2-2a2f-5688-102b-e6088b76e344",
		"accessor": "18bb8f89-826a-56ee-c65b-1736dc5ea27d",
		"policies": ["default"],
		"metadata": {
			"username": "vishal"
		},
		"lease_duration": 7200,
		"renewable": true
	}
   }
   ```

  </dd>
</dl>
