---
layout: "docs"
page_title: "Auth Backend: TLS Certificates"
sidebar_current: "docs-auth-cert"
description: |-
  The "cert" auth backend allows users to authenticate with Vault using TLS client certificates.
---

# Auth Backend: TLS Certificates

Name: `cert`

The "cert" auth backend allows authentication using SSL/TLS client certificates
which are either signed by a CA or self-signed.

The trusted certificates and CAs are configured directly to the auth backend
using the `certs/` path. This backend cannot read trusted certificates from an
external source.

CA certs are associated with a role; role names and CRL names are normalized to
lower-case.

## Revocation Checking

Since Vault 0.4, the backend supports revocation checking.

An authorised user can submit PEM-formatted CRLs identified by a given name;
these can be updated or deleted at will. (Note: Vault **does not** fetch CRLs;
the CRLs themselves and any updates must be pushed into Vault when desired,
such as via a `cron` job that fetches them from the source and pushes them into
Vault.)

When there are CRLs present, at the time of client authentication:

* If the client presents any chain where no certificate in the chain matches a
  revoked serial number, authentication is allowed
* If there is no chain presented by the client without a revoked serial number,
  authentication is denied

This method provides good security while also allowing for flexibility. For
instance, if an intermediate CA is going to be retired, a client can be
configured with two certificate chains: one that contains the initial
intermediate CA in the path, and the other that contains the replacement. When
the initial intermediate CA is revoked, the chain containing the replacement
will still allow the client to successfully authenticate.

**N.B.**: Matching is performed by *serial number only*. For most CAs,
including Vault's `pki` backend, multiple CRLs can successfully be used as
serial numbers are globally unique. However, since RFCs only specify that
serial numbers must be unique per-CA, some CAs issue serial numbers in-order,
which may cause clashes if attempting to use CRLs from two such CAs in the same
mount of the backend. The workaround here is to mount multiple copies of the
`cert` backend, configure each with one CA/CRL, and have clients connect to the
appropriate mount.

In addition, since the backend does not fetch the CRLs itself, the CRL's
designated time to next update is not considered. If a CRL is no longer in use,
it is up to the administrator to remove it from the backend.

## Authentication

### Via the CLI
The below requires Vault to present a certificate signed by `ca.pem` and
presents `cert.pem` (using `key.pem`) to authenticate against the `web` cert
role. If a certificate role name is not specified, the auth backend will try to
authenticate against all trusted certificates.

```
$ vault auth -method=cert \
    -ca-cert=ca.pem -client-cert=cert.pem -client-key=key.pem \
    name=web
```

### Via the API
The endpoint for the login is `/login`. The client simply connects with their TLS
certificate and when the login endpoint is hit, the auth backend will determine
if there is a matching trusted certificate to authenticate the client. Optionally,
you may specify a single certificate role to authenticate against.

```
$ curl --cacert ca.pem --cert cert.pem --key key.pem -d name=web \
    $VAULT_ADDR/v1/auth/cert/login -XPOST
```

## Configuration

First, you must enable the certificate auth backend:

```
$ vault auth-enable cert
Successfully enabled 'cert' at 'cert'!
```

Now when you run `vault auth -methods`, the certificate backend is available:

```
Path       Type      Description
cert/      cert
token/     token     token based credentials
```

To use the "cert" auth backend, an operator must configure it with
trusted certificates that are allowed to authenticate. An example is shown below.
Use `vault path-help` for more details.

```
$ vault write auth/cert/certs/web \
    display_name=web \
    policies=web,prod \
    certificate=@web-cert.pem \
    ttl=3600
...
```

The above creates a new trusted certificate "web" with same display name
and the "web" and "prod" policies. The certificate (public key) used to verify
clients is given by the "web-cert.pem" file. Lastly, an optional `ttl` value
can be provided in seconds to limit the lease duration.

#### Via the API

The token is set directly as a header for the HTTP API. The name
of the header should be "X-Vault-Token" and the value should be the token.

## API

### /auth/cert/certs

#### DELETE

<dl class="api">
  <dt>Description</dt>
  <dd>
    Deletes the named role and CA cert from the backend mount.
  </dd>

  <dt>Method</dt>
  <dd>DELETE</dd>

  <dt>URL</dt>
  <dd>`/auth/cert/certs/<name>`</dd>

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
    Gets information associated with the named role.
  </dd>

  <dt>Method</dt>
  <dd>GET</dd>

  <dt>URL</dt>
  <dd>`/auth/cert/certs/<name>`</dd>

  <dt>Parameters</dt>
  <dd>
    None
  </dd>

  <dt>Returns</dt>
  <dd>

    ```javascript
    {
      "lease_id": "",
      "renewable": false,
      "lease_duration": 0,
      "data": {
        "certificate": "-----BEGIN CERTIFICATE-----\nMIIEtzCCA5+.......ZRtAfQ6r\nwlW975rYa1ZqEdA=\n-----END CERTIFICATE-----",
        "display_name": "test",
        "policies": "",
        "allowed_names": "",
        "ttl": 2764800
      },
      "warnings": null,
      "auth": null
    }
    ```

  </dd>
</dl>

#### LIST

<dl class="api">
  <dt>Description</dt>
  <dd>
    Lists configured certificate names.
  </dd>

  <dt>Method</dt>
  <dd>LIST/GET</dd>

  <dt>URL</dt>
  <dd>`/auth/cert/certs` (LIST) or `/auth/cert/certs?list=true` (GET)</dd>

  <dt>Parameters</dt>
  <dd>
    None
  </dd>

  <dt>Returns</dt>
  <dd>

    ```javascript
    {
      "lease_id": "",
      "renewable": false,
      "lease_duration": 0,
      "data": {
        "keys": ["cert1", "cert2"]
      },
      "warnings": null,
      "auth": null
    }
    ```

  </dd>
</dl>

#### POST

<dl class="api">
  <dt>Description</dt>
  <dd>
    Sets a CA cert and associated parameters in a role name.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/auth/cert/certs/<name>`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">certificate</span>
        <span class="param-flags">required</span>
        The PEM-format CA certificate.
      </li>
      <li>
        <span class="param">allowed_names</span>
        <span class="param-flags">optional</span>
        Constrain the Common and Alternative Names in the client certificate
        with a [globbed pattern](https://github.com/ryanuber/go-glob/blob/master/README.md#example).
        Value is a comma-separated list of patterns.
        Authentication requires at least one Name matching at least one pattern.
        If not set, defaults to allowing all names.
      </li>
      <li>
        <span class="param">policies</span>
        <span class="param-flags">optional</span>
        A comma-separated list of policies to set on tokens issued when
        authenticating against this CA certificate.
      </li>
      <li>
        <span class="param">display_name</span>
        <span class="param-flags">optional</span>
        The `display_name` to set on tokens issued when authenticating
        against this CA certificate. If not set, defaults to the name
        of the role.
      </li>
      <li>
        <span class="param">ttl</span>
        <span class="param-flags">optional</span>
        The TTL period of the token, provided as a number of seconds. If not
        provided, the token is valid for the the mount or system default TTL
        time, in that order.
      </li>
    </ul>
  </dd>

  <dt>Returns</dt>
  <dd>
    A `204` response code.
  </dd>
</dl>

### /auth/cert/crls

#### DELETE

<dl class="api">
  <dt>Description</dt>
  <dd>
    Deletes the named CRL from the backend mount.
  </dd>

  <dt>Method</dt>
  <dd>DELETE</dd>

  <dt>URL</dt>
  <dd>`/auth/cert/crls/<name>`</dd>

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
    Gets information associated with the named CRL (currently, the serial
    numbers contained within).  As the serials can be integers up to an
    arbitrary size, these are returned as strings.
  </dd>

  <dt>Method</dt>
  <dd>GET</dd>

  <dt>URL</dt>
  <dd>`/auth/cert/crls/<name>`</dd>

  <dt>Parameters</dt>
  <dd>
    None
  </dd>

  <dt>Returns</dt>
  <dd>

    ```javascript
    {
        "auth": null,
        "data": {
            "serials": {
                "13": {}
            }
        },
        "lease_duration": 0,
        "lease_id": "",
        "renewable": false,
        "warnings": null
    }

    ```

  </dd>
</dl>

#### POST

<dl class="api">
  <dt>Description</dt>
  <dd>
    Sets a named CRL.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/auth/cert/crls/<name>`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">crl</span>
        <span class="param-flags">required</span>
        The PEM-format CRL.
      </li>
    </ul>
  </dd>

  <dt>Returns</dt>
  <dd>
    A `204` response code.
  </dd>
</dl>

### /auth/cert/login

#### POST

<dl class="api">
  <dt>Description</dt>
  <dd>
    Log in and fetch a token. If there is a valid chain to a CA configured in
    the backend and all role constraints are matched, a token will be issued.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/auth/cert/login`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">name</span>
        <span class="param-flags">optional</span>
        Authenticate against only the named certificate role, returning its
        policy list if successful. If not set, defaults to trying all
        certificate roles and returning any one that matches.
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
        "lease_duration": 3600,
        "renewable": true,
      }
    }
    ```

  </dd>
</dl>

### /auth/cert/config

#### POST

<dl class="api">
  <dt>Description</dt>
  <dd>
    Configuration options for the backend.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/auth/cert/config`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">disable_binding</span>
        <span class="param-flags">optional</span>
	  If set, during renewal, skips the matching of presented client identity with the client identity used during login. Defaults to false.
      </li>
    </ul>
  </dd>

  <dt>Returns</dt>
  <dd>
    A `204` response code.
  </dd>
</dl>
