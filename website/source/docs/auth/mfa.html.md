---
layout: "docs"
page_title: "Multi-Factor Authentication"
sidebar_current: "docs-auth-mfa"
description: |-
  Multi-factor authentication is supported for several authentication backends.
---

# Multi-Factor Authentication

Several authentication backends support multi-factor authentication (MFA). Once enabled for
a backend, users are required to provide additional verification, like a one-time passcode,
before being authenticated.

Currently, the "ldap", "radius" and "userpass" backends support MFA.

## Authentication

When authenticating, users still provide the same information as before, in addition to
MFA verification. Usually this is a passcode, but in other cases, like a Duo Push
notification, no additional information is needed.

### Via the CLI

```shell
$ vault auth -method=userpass \
    username=user \
    password=test \
    passcode=111111
```
```shell
$ vault auth -method=userpass \
    username=user \
    password=test \
    method=push
```

### Via the API

The endpoint for the login is the same as for the original backend. Additional
MFA information should be sent in the POST body encoded as JSON.

```shell
$ curl $VAULT_ADDR/v1/auth/userpass/login/user \
    -d '{ "password": "test", "passcode": "111111" }'
```

The response is the same as for the original backend.

## Configuration

To enable MFA for a supported backend, the MFA type must be set in `mfa_config`. For example:

```shell
$ vault write auth/userpass/mfa_config type=duo
```

This enables the Duo MFA type, which is currently the only MFA type supported. The username
used for MFA is the same as the login username, unless the backend or MFA type provide
options to behave differently (see Duo configuration below).

### Duo

The Duo MFA type is configured through two paths: `duo/config` and `duo/access`.

`duo/access` contains connection information for the Duo Auth API. To configure:

```shell
$ vault write auth/[mount]/duo/access \
    host=[host] \
    ikey=[integration key] \
    skey=[secret key]
```

`duo/config` is an optional path that contains general configuration information
for Duo authentication. To configure:

```shell
$ vault write auth/[mount]/duo/config \
    user_agent="" \
    username_format="%s"
```

`user_agent` is the user agent to use when connecting to Duo.

`username_format` controls how the username used to login is
transformed before authenticating with Duo. This field is a format string
that is passed the original username as its first argument and outputs
the new username. For example "%s@example.com" would append "@example.com"
to the provided username before connecting to Duo.

`push_info` is a string of URL-encoded key/value pairs that provides additional
context about the authentication attempt in the Duo Mobile application.

More information can be found through the CLI `path-help` command.
