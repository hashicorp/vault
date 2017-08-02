---
layout: "docs"
page_title: "Authentication"
sidebar_current: "docs-concepts-auth"
description: |-
  Before performing any operation with Vault, the connecting client must be authenticated.
---

# Authentication

Authentication in Vault is the process by which user or machine supplied
information is verified against an internal or external system. Vault supports
multiple [authentication backends](/docs/auth/index.html) including GitHub,
LDAP, AppRole, and more. Each authentication backend has a specific use case.

Before a client can interact with Vault, it must _authenticate_ against an
authentication backend. Upon authentication, a token is generated. This token is
conceptually similar to a session ID on a website. The token may have attached
policy, which is mapped at authentication time. This process is described in
detail in the [policies concepts](/docs/concepts/policies.html) documentation.

## Authentication Backends

Vault supports a number of authentication backends. Some backends are targeted
toward users while others are targeted toward machines. Most authentication
backends must be enabled before use. To enable an authentication backend:

```sh
$ vault write sys/auth/my-auth type=userpass
```

This mounts the "userpass" authentication backend at the path "my-auth". This
authentication will be accessible at the path "my-auth". Often you will see
authentications at the same path as their name, but this is not a requirement.

To learn more about this authentication, use the built-in `path-help` command:

```sh
$ vault path-help auth/my-auth
# ...
```

Vault supports multiple authentication backends simultaneously, and you can even
mount the same type of authentication backend at different paths. Only one
authentication is required to gain access to Vault, and it is not currently
possible to force a user through multiple authentication backends to gain
access, although some backends do support MFA.

## Tokens

There is an [entire page dedicated to tokens](/docs/concepts/tokens.html),
but it is important to understand that authentication works by verifying
your identity and then generating a token to associate with that identity.

For example, even though you may authenticate using something like GitHub,
Vault generates a unique access token for you to use for future requests.
The CLI automatically attaches this token to requests, but if you're using
the API you'll have to do this manually.

This token given for authentication with any backend can also be used
with the full set of token commands, such as creating new sub-tokens,
revoking tokens, and renewing tokens. This is all covered on the
[token concepts page](/docs/concepts/tokens.html).

## Authenticating

### Via the CLI

To authenticate with the CLI, `vault auth` is used. This supports many
of the built-in authentication methods. For example, with GitHub:

```
$ vault auth -method=github token=<token>
...
```

After authenticating, you will be logged in. The CLI command will also
output your raw token. This token is used for revocation and renewal.
As the user logging in, the primary use case of the token is renewal,
covered below in the "Auth Leases" section.

To determine what variables are needed for an authentication method,
supply the `-method` flag without any additional arguments and help
will be shown.

If you're using a method that isn't supported via the CLI, then the API
must be used.

### Via the API

API authentication is generally used for machine authentication. Each
auth backend implements its own login endpoint. Use the `vault path-help`
mechanism to find the proper endpoint.

For example, the GitHub login endpoint is located at `auth/github/login`.
And to determine the arguments needed, `vault path-help auth/github/login` can
be used.

## Auth Leases

Just like secrets, identities have
[leases](/docs/concepts/lease.html) associated with them. This means that
you must reauthenticate after the given lease period to continue accessing
Vault.

To set the lease associated with an identity, reference the help for
the specific authentication backend in use. It is specific to each backend
how leasing is implemented.

And just like secrets, identities can be renewed without having to
completely reauthenticate. Just use `vault token-renew <token>` with the
leased token associated with your identity to renew it.
