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

The trusted certificates and CAs are configured directly to the auth
backend using the `certs/` path. This backend cannot read trusted certificates
from an external source.

## Authentication

### Via the CLI
```
vault auth -method=cert \
  -ca-cert=ca.pem -client-cert=cert.pem -client-key=key.pem
```

### Via the API
The endpoint for the login is `/login`. The client simply connects with their TLS
certificate and when the login endpoint is hit, the auth backend will determine
if there is a matching trusted certificate to authenticate the client.

```
curl --cacert ca.pem --cert cert.pem --key key.pem \
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
$ vault write auth/cert/certs/web display_name=web policies=web,prod certificate=@web-cert.pem lease=3600
...
```

The above creates a new trusted certificate "web" with same display name
and the "web" and "prod" policies. The certificate (public key) used to verify
clients is given by the "web-cert.pem" file. Lastly, an optional lease value
can be provided in seconds to limit the lease period.

