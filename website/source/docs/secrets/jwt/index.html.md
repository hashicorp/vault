# JWT Secret Backend

Name: `jwt`

The JWT secret backend for Vault generates JSON Web Tokens dynamically based on configured roles. This means services can get tokens needed for authentication without going through the usual manual process of generating a private key and signing the token and maintaining the private key's security. Vault's built-in authentication and authorization mechanisms provide the verification functionality.

This page will show a quick start for this backend. For detailed documentation on every path, use `vault path-help` after mounting the backend.

The JWT secret backend acts like the `transit` backend, it does not store any information.

## Algorithms

### RSA 
* RS256
* RS384
* RS512

These require a RSA private/public keypair for signing and verification.

### ECDSA
* EC256
* EC384
* EC512

These require an ECDSA private/public keypair for signing and verification.

### HMAC
* HS256
* HS384
* HS512

These require a shared secret for signing and verification.

## Roles

Roles are defined with the signing algorithm, the secret key or private key to be used, as well as allowing for default but optional JWT Token claims. Once you write a private key or a secret to the role, it CANNOT be read back out.

## Quick Start

The first step to using the jwt backend is to mount it.
Unlike the `generic` backend, the `jwt` backend is not mounted by default.

```text
$ vault mount jwt
Successfully mounted 'jwt' at 'jwt'!
```

The next step is to configure a role. A role is a logical name that maps
to a few settings used to generated the tokens. For example, lets create
a "webauth" role:

```text
$ vault write jwt/roles/webauth \
    algorithm=RS256 \
    key=@/path/to/private.key
```

Each role requires a secret or a private key to be associated against it.

Generating a token requires passing of additional information so we use the
"jwt/issue/ROLE" path.

```text
$ vault write jwt/issue/webauth \
    issuer="Vault" \
    audience="Vault Client" \
    expiration="1538096292" \
    claims=@extra.json
```


## API

### /jwt/roles/
#### POST

<dl class="api">
  <dt>Description</dt>
  <dd>
    Creates or updates a named role.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/jwt/roles/<name>`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">algorithm</span>
        <span class="param-flags">required</span>
        The algorithm used by JWT to sign the token.
      </li>
      <li>
        <span class="param">key</span>
        <span class="param-flags">required</span>
        The private key or string used to sign the token.
      </li>
      <li>
        <span class="param">default_issuer</span>
        <span class="param-flags">required</span>
        The default issuer claim for the role, can be overridden at issue time.
      </li>
      <li>
        <span class="param">default_subject</span>
        <span class="param-flags">required</span>
        The default subject claim for the role, can be overridden at issue time.
      </li>
      <li>
        <span class="param">default_audience</span>
        <span class="param-flags">required</span>
        The default audience claim for the role, can be overridden at issue time.
      </li>
    </ul>
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
    Queries a named role.
  </dd>

  <dt>Method</dt>
  <dd>GET</dd>

  <dt>URL</dt>
  <dd>`/jwt/roles/<name>`</dd>

  <dt>Parameters</dt>
  <dd>
    None
  </dd>

  <dt>Returns</dt>
  <dd>

    ```javascript
    {
        "data": {
            "algorithm": "..."
        }
    }
    ```

  </dd>
</dl>

#### DELETE

<dl class="api">
  <dt>Description</dt>
  <dd>
    Deletes a named role.
  </dd>

  <dt>Method</dt>
  <dd>DELETE</dd>

  <dt>URL</dt>
  <dd>`/jwt/roles/<name>`</dd>

  <dt>Parameters</dt>
  <dd>
    None
  </dd>

  <dt>Returns</dt>
  <dd>
    A `204` response code.
  </dd>
</dl>


### /jwt/issue/
#### POST

<dl class="api">
  <dt>Description</dt>
  <dd>
    Generates a JWT token based on the named role.
  </dd>

  <dt>Method</dt>
  <dd>GET</dd>

  <dt>URL</dt>
  <dd>`/jwt/issue/<role>`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">issuer</span>
        <span class="param-flags">optional</span>
        The Issuer of the token.
      </li>
      <li>
        <span class="param">audience</span>
        <span class="param-flags">optional</span>
        The Audience of the token.
      </li>
      <li>
        <span class="param">subject</span>
        <span class="param-flags">optional</span>
        The Subject of the token.
      </li>
      <li>
        <span class="param">expiration</span>
        <span class="param-flags">optional</span>
        The expiration of the token, expressed in seconds (unix time).
      </li>
      <li>
        <span class="param">issued_at</span>
        <span class="param-flags">optional</span>
        The issued at time of the token, expressed in seconds (unix time). (Default: current time)
      </li>
      <li>
        <span class="param">not_before</span>
        <span class="param-flags">optional</span>
        Not Before: the time at which the token is not useful before. Expressed as seconds, unix time. (Default: current time)
      </li>
      <li>
        <span class="param">jti</span>
        <span class="param-flags">optional</span>
        JSONWebToken Identifier. Unique ID useful for preventing replay attacks. (Default: Random UUID)
      </li>
      <li>
        <span class="param">claims</span>
        <span class="param-flags">optional</span>
        Should be a JSON Object of additional key/values you want in the token.
      </li>
    </ul>
  </dd>

  <dt>Returns</dt>
  <dd>

    ```javascript
    {
        "data": {
            "jti": "...",
            "token": "..."
        }
    }
    ```

  </dd>
</dl>
