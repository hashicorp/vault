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

The endpoint for the login is `/login`. The client simply connects with their TLS
certificate and when the login endpoint is hit, the auth backend will determine
if there is a matching trusted certificate to authenticate the client.

## Configuration

To use the "cert" auth backend, an operator must configure it with
trusted certificates that are allowed to authenticate. An example is shown below.
Use `vault help` for more details.

```
$ vault write auth/cert/certs/web display_name=web policies=web,prod certificate=@web-cert.pem
...
```

The above creates a new trusted certificate "web" with same display name
and the "web" and "prod" policies. The certificate (public key) used to verify
clients is given by the "web-cert.pem" file.

