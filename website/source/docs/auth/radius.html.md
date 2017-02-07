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

### /auth/radius/config
#### POST

<dl class="api">
  <dt>Description</dt>
  <dd>
      Configures the connection parameters and shard secret used to communicate with RADIUS
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/auth/radius/config`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">host</span>
        <span class="param-flags">required</span>
            The RADIUS server to connect to. Examples: `radius.myorg.com`, `127.0.0.1`
      </li>
      <li>
        <span class="param">port</span>
        <span class="param-flags">optional</span>
            The UDP port where the RADIUS server is listening on. Defaults is 1812
      </li>
      <li>
        <span class="param">secret</span>
        <span class="param-flags">required</span>
            The RADIUS shared secret
      </li>
      <li>
        <span class="param">unregistered_user_policies</span>
        <span class="param-flags">optional</span>
            A Comma-Separated list of policies to be granted to unregistered users
      </li>
      <li>
        <span class="param">dial_timeout</span>
        <span class="param-flags">optional</span>
            Number of second to wait for a backend connection before timing out. Defaults is 10
      </li>
      <li>
        <span class="param">read_timeout</span>
        <span class="param-flags">optional</span>
            Number of second to wait for a backend response before timing out. Defaults is 10
      </li>
      <li>
        <span class="param">nas_port</span>
        <span class="param-flags">optional</span>
            The NAS-Port attribute of the RADIUS request. Defaults is 10
      </li>
    </ul>
  </dd>

  <dt>Returns</dt>
  <dd>`204` response code.
  </dd>
</dl>

### /auth/radius/users/[username]
#### POST

<dl class="api">
  <dt>Description</dt>
  <dd>
      Registers a new user and maps a set of policies to it.
      This path honors the distinction between the `create` and `update` capabilities inside ACL policies.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/auth/radius/users/<username>`</dd>

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
            If set to empty string, only the `default` policy will be applicable to the user.
      </li>
    </ul>
  </dd>

  <dt>Returns</dt>
  <dd>`204` response code.
  </dd>
</dl>

#### GET
<dl class="api">
  <dt>Description</dt>
  <dd>
  Reads the properties of an existing username.
  </dd>

  <dt>Method</dt>
  <dd>GET</dd>

  <dt>URL</dt>
  <dd>`/auth/radius/users/[username]`</dd>

  <dt>Parameters</dt>
  <dd>
    None.
  </dd>

  <dt>Returns</dt>
  <dd>

```javascript
{
        "request_id": "812229d7-a82e-0b20-c35b-81ce8c1b9fa6",
        "lease_id": "",
        "lease_duration": 0,
        "renewable": false,
        "data": {
                "policies": "default,dev"
        },
        "warnings": null
}
```

  </dd>
</dl>


#### DELETE
<dl class="api">
  <dt>Description</dt>
  <dd>
  Deletes an existing username from the backend.
  </dd>

  <dt>Method</dt>
  <dd>DELETE</dd>

  <dt>URL</dt>
  <dd>`/auth/radius/users/[username]`</dd>

  <dt>Parameters</dt>
  <dd>
    None.
  </dd>

  <dt>Returns</dt>
  <dd>`204` response code.
  </dd>
</dl>

### /auth/radius/login
### /auth/radius/login/[username]
#### POST
<dl class="api">
  <dt>Description</dt>
  <dd>
      Login with the username and password.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URLS</dt>
  <dd>`/auth/radius/login`</dd>
  <dd>`/auth/radius/login/[username]`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
        <li>
        <span class="param">username</span>
        <span class="param-flags">required</span>
            Username for the authenticating user.
      </li>
      <li>
        <span class="param">password</span>
        <span class="param-flags">required</span>
            Password for the authenticating user.
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

### /auth/radius/users
#### LIST
<dl class="api">
  <dt>Description</dt>
  <dd>
List the users registered with the backend.
  </dd>

  <dt>Method</dt>
  <dd>LIST/GET</dd>

  <dt>URL</dt>
  <dd>`/auth/radius/users` (LIST) `/auth/radius/users?list=true` (GET)</dd>

  <dt>Parameters</dt>
  <dd>
None
  </dd>

  <dt>Returns</dt>
  <dd>

   ```javascript
[
        "devuser",
	    "produser"
]
   ```

  </dd>
</dl>


