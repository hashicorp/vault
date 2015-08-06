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

The token is set directly as a cookie for the HTTP API. The name
of the cookie should be "token" and the value should be the token.

## API

### /auth/token/create
#### POST

<dl class="api">
  <dt>Description</dt>
  <dd>
    Creates a new token. Certain options are only available to
    when called by a root token.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/auth/token/create`</dd>

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
        <span class="param">lease</span>
        <span class="param-flags">optional</span>
        The lease period of the token, provided as "1h", where hour is
        the largest suffix. If not provided, the token is valid for the
        [default lease duration](/docs/config/index.html), or
        indefinitely if the root policy is used.
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
        no limit to number of uses.
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

### /auth/token/lookup/
#### GET

<dl class="api">
  <dt>Description</dt>
  <dd>
    Returns information about the current client token.
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


### /auth/token/revoke/
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
  <dd>`/auth/token/revoke/<token>`</dd>

  <dt>Parameters</dt>
  <dd>
    None
  </dd>

  <dt>Returns</dt>
  <dd>`204` response code.
  </dd>
</dl>

### /auth/token/revoke-orphan/
#### POST

<dl class="api">
  <dt>Description</dt>
  <dd>
    Revokes a token but not its child tokens. When the token is revoked,
    all secrets generated with it are also revoked. All child tokens
    are orphaned, but can be revoked sub-sequently using `/auth/token/revoke/`.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/auth/token/revoke-orphan/<token>`</dd>

  <dt>Parameters</dt>
  <dd>
    None
  </dd>

  <dt>Returns</dt>
  <dd>`204` response code.
  </dd>
</dl>

### /auth/token/revoke-prefix/
#### POST

<dl class="api">
  <dt>Description</dt>
  <dd>
    Revokes all tokens generated at a given prefix, along with child tokens,
    and all secrets generated using those tokens. Uses include revoking all
    tokens generated by a credential backend during a suspected compromise.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/auth/token/revoke-prefix/<prefix>`</dd>

  <dt>Parameters</dt>
  <dd>
    None
  </dd>

  <dt>Returns</dt>
  <dd>`204` response code.
  </dd>
</dl>

### /auth/token/renew/
#### POST

<dl class="api">
  <dt>Description</dt>
  <dd>
    Renews a lease associated with a token. This is used to prevent
    the expiration of a token, and the automatic revocation of it.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/auth/token/renew/<token>`</dd>

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
</div>
