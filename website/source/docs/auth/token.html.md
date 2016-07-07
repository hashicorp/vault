---
layout: "docs"
page_title: "Auth Backend: Token"
sidebar_current: "docs-auth-token"
description: |-
  The token store auth backend is used to authenticate using tokens.
---

# Auth Backend: Token

The token backend is the only auth backend that is built-in and
automatically available at `/auth/token` as well as with first-class
built-in CLI methods such as `vault token-create`. It allows users to
authenticate using a token, as well to create new tokens, revoke
secrets by token, and more.

When any other auth backend returns an identity, Vault core invokes the
token backend to create a new unique token for that identity.

The token store can also be used to bypass any other auth backend:
you can create tokens directly, as well as perform a variety of other
operations on tokens such as renewal and revocation.

Please see the [token concepts](/docs/concepts/tokens.html) page dedicated
to tokens.

## Authentication

#### Via the CLI

```
$ vault auth <token>
...
```

#### Via the API

The token is set directly as a header for the HTTP API. The name
of the header should be "X-Vault-Token" and the value should be the token.

## API

### /auth/token/create
### /auth/token/create-orphan
### /auth/token/create/[role_name]
#### POST

<dl class="api">
  <dt>Description</dt>
  <dd>
    Creates a new token. Certain options are only available when called by a
    root token. If used via the `/auth/token/create-orphan` endpoint, a root
    token is not required to create an orphan token (otherwise set with the
    `no_parent` option). If used with a role name in the path, the token will
    be created against the specified role name; this may override options set
    during this call.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URLs</dt>
  <dd>`/auth/token/create`</dd>
  <dd>`/auth/token/create-orphan`</dd>
  <dd>`/auth/token/create/<role_name>`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">id</span>
        <span class="param-flags">optional</span>
        The ID of the client token. Can only be specified by a root token.
        Otherwise, the token ID is a randomly generated UUID.
      </li>
      <li>
        <span class="param">policies</span>
        <span class="param-flags">optional</span>
        A list of policies for the token. This must be a subset of the
        policies belonging to the token making the request, unless root.
        If not specified, defaults to all the policies of the calling token.
      </li>
      <li>
        <span class="param">meta</span>
        <span class="param-flags">optional</span>
        A map of string to string valued metadata. This is passed through
        to the audit backends.
      </li>
      <li>
        <span class="param">no_parent</span>
        <span class="param-flags">optional</span>
        If true and set by a root caller, the token will not have the
        parent token of the caller. This creates a token with no parent.
      </li>
      <li>
        <span class="param">no_default_policy</span>
        <span class="param-flags">optional</span>
        If true the `default` policy will not be a part of this token's
        policy set.
      </li>
      <li>
        <span class="param">renewable</span>
        <span class="param-flags">optional</span>
        Set to `false` to disable the ability of the token to be renewed past
        its initial TTL. Specifying `true`, or omitting this option, will allow
        the token to be renewable up to the system/mount maximum TTL.
      </li>
      <li>
        <span class="param">lease</span>
        <span class="param-flags">optional</span>
        DEPRECATED; use "ttl" instead.
      </li>
      <li>
        <span class="param">ttl</span>
        <span class="param-flags">optional</span>
        The TTL period of the token, provided as "1h", where hour is
        the largest suffix. If not provided, the token is valid for the
        [default lease TTL](/docs/config/index.html), or
        indefinitely if the root policy is used.
      </li>
      <li>
        <span class="param">explicit_max_ttl</span>
        <span class="param-flags">optional</span>
        If set, the token will have an explicit max TTL set upon it. This
        maximum token TTL *cannot* be changed later, and unlike with normal
        tokens, updates to the system/mount max TTL value will have no effect
        at renewal time -- the token will never be able to be renewed or used
        past the value set at issue time. 
      </li>
      <li>
        <span class="param">display_name</span>
        <span class="param-flags">optional</span>
        The display name of the token. Defaults to "token".
      </li>
      <li>
        <span class="param">num_uses</span>
        <span class="param-flags">optional</span>
        The maximum uses for the given token. This can be used to create
        a one-time-token or limited use token. Defaults to 0, which has
        no limit to the number of uses.
      </li>
    </ul>
  </dd>

  <dt>Returns</dt>
  <dd>

    ```javascript
    {
      "auth": {
        "client_token": "ABCD",
        "policies": ["web", "stage"],
        "metadata": {"user": "armon"},
        "lease_duration": 3600,
        "renewable": true,
      }
    }
    ```

  </dd>
</dl>

### /auth/token/lookup-self
#### GET

<dl class="api">
  <dt>Description</dt>
  <dd>
    Returns information about the current client token.
  </dd>

  <dt>Method</dt>
  <dd>GET</dd>

  <dt>Parameters</dt>
  <dd>
    None
  </dd>

  <dt>Returns</dt>
  <dd>

    ```javascript
    {
      "data": {
        "id": "ClientToken",
        "policies": ["web", "stage"],
        "path": "auth/github/login",
        "meta": {"user": "armon", "organization": "hashicorp"},
        "display_name": "github-armon",
        "num_uses": 0,
      }
    }
    ```
  </dd>
</dl>

### /auth/token/lookup[/token]
#### GET

<dl class="api">
  <dt>Description</dt>
  <dd>
    Returns information about the client token provided in the request path.
  </dd>

  <dt>Method</dt>
  <dd>GET</dd>

  <dt>URL</dt>
  <dd>`/auth/token/lookup/<token>`</dd>

  <dt>Parameters</dt>
  <dd>
    None
  </dd>

  <dt>Returns</dt>
  <dd>

    ```javascript
    {
      "data": {
        "id": "ClientToken",
        "policies": ["web", "stage"],
        "path": "auth/github/login",
        "meta": {"user": "armon", "organization": "hashicorp"},
        "display_name": "github-armon",
        "num_uses": 0,
      }
    }
    ```

  </dd>
</dl>


#### POST

<dl class="api">
  <dt>Description</dt>
  <dd>
    Returns information about the client token provided in the request body.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/auth/token/lookup`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">token</span>
        <span class="param-flags">required</span>
            Token to lookup.
      </li>
    </ul>
  </dd>

  <dt>Returns</dt>
  <dd>

    ```javascript
    {
      "data": {
        "id": "ClientToken",
        "policies": ["web", "stage"],
        "path": "auth/github/login",
        "meta": {"user": "armon", "organization": "hashicorp"},
        "display_name": "github-armon",
        "num_uses": 0,
      }
    }
    ```

  </dd>
</dl>

### /auth/token/renew-self
#### POST

<dl class="api">
  <dt>Description</dt>
  <dd>
	Renews a lease associated with the calling token. This is used to prevent
	the expiration of a token, and the automatic revocation of it. Token
	renewal is possible only if there is a lease associated with it.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/auth/token/renew-self`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">increment</span>
        <span class="param-flags">optional</span>
            An optional requested lease increment can be provided. This
            increment may be ignored.
      </li>
    </ul>
  </dd>

  <dt>Returns</dt>
  <dd>

    ```javascript
    {
      "auth": {
        "client_token": "ABCD",
        "policies": ["web", "stage"],
        "metadata": {"user": "armon"},
        "lease_duration": 3600,
        "renewable": true,
      }
    }
    ```

  </dd>
</dl>

### /auth/token/renew[/token]
#### POST

<dl class="api">
  <dt>Description</dt>
  <dd>
    Renews a lease associated with a token. This is used to prevent the
    expiration of a token, and the automatic revocation of it. Token
    renewal is possible only if there is a lease associated with it.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/auth/token/renew</token>`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">token</span>
        <span class="param-flags">required</span>
            Token to renew. This can be part of the URL or the body.
      </li>
    </ul>
  </dd>
  <dd>
    <ul>
      <li>
        <span class="param">increment</span>
        <span class="param-flags">optional</span>
            An optional requested lease increment can be provided. This
            increment may be ignored.
      </li>
    </ul>
  </dd>

  <dt>Returns</dt>
  <dd>

    ```javascript
    {
      "auth": {
        "client_token": "ABCD",
        "policies": ["web", "stage"],
        "metadata": {"user": "armon"},
        "lease_duration": 3600,
        "renewable": true,
      }
    }
    ```

  </dd>
</dl>

### /auth/token/revoke[/token]
#### POST

<dl class="api">
  <dt>Description</dt>
  <dd>
    Revokes a token and all child tokens. When the token is revoked,
    all secrets generated with it are also revoked.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/auth/token/revoke</token>`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">token</span>
        <span class="param-flags">required</span>
            Token to revoke. This can be part of the URL or the body.
      </li>
    </ul>
  </dd>

  <dt>Returns</dt>
  <dd>`204` response code.
  </dd>
</dl>

### /auth/token/revoke-self/
#### POST

<dl class="api">
  <dt>Description</dt>
  <dd>
    Revokes the token used to call it and all child tokens.
    When the token is revoked, all secrets generated with
    it are also revoked.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/auth/token/revoke-self`</dd>

  <dt>Parameters</dt>
  <dd>
    None
  </dd>

  <dt>Returns</dt>
  <dd>`204` response code.
  </dd>
</dl>

### /auth/token/revoke-orphan[/token]
#### POST

<dl class="api">
  <dt>Description</dt>
  <dd>
    Revokes a token but not its child tokens. When the token is revoked, all
    secrets generated with it are also revoked. All child tokens are orphaned,
    but can be revoked sub-sequently using `/auth/token/revoke/`. This is a
    root-protected endpoint.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/auth/token/revoke-orphan</token>`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">token</span>
        <span class="param-flags">required</span>
            Token to revoke. This can be part of the URL or the body.
      </li>
    </ul>
  </dd>

  <dt>Returns</dt>
  <dd>`204` response code.
  </dd>
</dl>

### /auth/token/roles/[role_name]

#### DELETE 

<dl class="api">
  <dt>Description</dt>
  <dd>
    Deletes the named role.
  </dd>

  <dt>Method</dt>
  <dd>DELETE</dd>

  <dt>URL</dt>
  <dd>`/auth/token/roles/<role_name>`</dd>

  <dt>Parameters</dt>
  <dd>
    None
  </dd>

  <dt>Returns</dt>
  <dd>
    A `204` response code.
  </dd>
</dl>

#### GET 

<dl class="api">
  <dt>Description</dt>
  <dd>
    Fetches the named role configuration.
  </dd>

  <dt>Method</dt>
  <dd>GET</dd>

  <dt>URL</dt>
  <dd>`/auth/token/roles/<role_name>`</dd>

  <dt>Parameters</dt>
  <dd>
    None
  </dd>

  <dt>Returns</dt>
  <dd>

    ```javascript
    {
      "data": {
        "period": 3600,
        "allowed_policies": ["web", "stage"],
        "orphan": true,
        "path_suffix": ""
      }
    }
    ```

  </dd>
</dl>

#### LIST 

<dl class="api">
  <dt>Description</dt>
  <dd>
    Lists available roles.
  </dd>

  <dt>Method</dt>
  <dd>GET</dd>

  <dt>URL</dt>
  <dd>`/auth/token/roles?list=true`<dd>

  <dt>Parameters</dt>
  <dd>
    None
  </dd>

  <dt>Returns</dt>
  <dd>

    ```javascript
    {
      "data": {
        "keys": ["role1", "role2"]
      }
    }
    ```

  </dd>
</dl>

#### POST

<dl class="api">
  <dt>Description</dt>
  <dd>
    Creates (or replaces) the named role. Roles enforce specific behavior when
    creating tokens that allow token functionality that is otherwise not
    available or would require `sudo`/root privileges to access. Role
    parameters, when set, override any provided options to the `create`
    endpoints. The role name is also included in the token path, allowing all
    tokens created against a role to be revoked using the `sys/revoke-prefix`
    endpoint.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/auth/token/roles/<role_name>`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">allowed_policies</span>
        <span class="param-flags">optional</span>
        If set, tokens can be created with any subset of the policies in this
        list, rather than the normal semantics of tokens being a subset of the
        calling token's policies. The parameter is a comma-delimited string of
        policy names.
      </li>
      <li>
        <span class="param">orphan</span>
        <span class="param-flags">optional</span>
        If `true`, tokens created against this policy will be orphan tokens
        (they will have no parent). As such, they will not be automatically
        revoked by the revocation of any other token.
      </li>
      <li>
        <span class="param">period</span>
        <span class="param-flags">optional</span>
        If set, tokens created against this role will <i>not</i> have a maximum
        lifetime. Instead, they will have a fixed TTL that is refreshed with
        each renewal. So long as they continue to be renewed, they will never
        expire. The parameter is an integer duration of seconds. Tokens issued
        track updates to the role value; the new period takes effect upon next
        renew. This cannot be used in conjunction with `explicit_max_ttl`.
      </li>
      <li>
        <span class="param">renewable</span>
        <span class="param-flags">optional</span>
        Set to `false` to disable the ability of token created against this
        role to be renewed past their initial TTL. Defaults to `true`, which
        allows tokens to be renewed up to the system/mount maximum TTL.
      </li>
      <li>
        <span class="param">path_suffix</span>
        <span class="param-flags">optional</span>
        If set, tokens created against this role will have the given suffix as
        part of their path in addition to the role name. This can be useful in
        certain scenarios, such as keeping the same role name in the future but
        revoking all tokens created against it before some point in time. The
        suffix can be changed, allowing new callers to have the new suffix as
        part of their path, and then tokens with the old suffix can be revoked
        via `sys/revoke-prefix`.
      </li>
      <li>
        <span class="param">explicit_max_ttl</span>
        <span class="param-flags">optional</span>
        If set, tokens created with this role have an explicit max TTL set upon
        them. This maximum token TTL *cannot* be changed later, and unlike with
        normal tokens, updates to the role or the system/mount max TTL value
        will have no effect at renewal time -- the token will never be able to
        be renewed or used past the value set at issue time. This cannot be
        used in conjunction with `period`.
      </li>
    </ul>
  </dd>

  <dt>Returns</dt>
  <dd>
    A `204` return code.
  </dd>
</dl>

### /auth/token/lookup-accessor[/accessor]
#### POST

<dl class="api">
  <dt>Description</dt>
  <dd>
      Fetch the properties of the token associated with the accessor, except the token ID.
      This is meant for purposes where there is no access to token ID but there is need
      to fetch the properties of a token.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/auth/token/lookup-accessor</accessor>`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">accessor</span>
        <span class="param-flags">required</span>
            Accessor of the token to lookup. This can be part of the URL or the body.
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
	"data": {
		"creation_time": 1457533232,
		"creation_ttl": 2592000,
		"display_name": "token",
		"id": "",
		"meta": null,
		"num_uses": 0,
		"orphan": false,
		"path": "auth/token/create",
		"policies": ["default", "web"],
		"ttl": 2591976
	},
	"warnings": null,
	"auth": null
   }
    ```
  </dd>
</dl>

### /auth/token/revoke-accessor[/accessor]
#### POST

<dl class="api">
  <dt>Description</dt>
  <dd>
      Revoke the token associated with the accessor and all the child tokens.
      This is meant for purposes where there is no access to token ID but
      there is need to revoke a token and its children.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/auth/token/revoke-accessor</accessor>`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">accessor</span>
        <span class="param-flags">required</span>
            Accessor of the token. This can be part of the URL or the body.
      </li>
    </ul>
  </dd>

  <dt>Returns</dt>
  <dd>`204` response code.
  </dd>
</dl>

