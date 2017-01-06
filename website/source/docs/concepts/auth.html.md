---
layout: "docs"
page_title: "Authentication"
sidebar_current: "docs-concepts-auth"
description: |-
  Before performing any operation with Vault, the connecting client must be authenticated.
---

# Authentication

Before performing any operation with Vault, the connecting client must be
_authenticated_. Authentication is the process of verifying a person or
machine is who they say they are and assigning an identity to them. This
identity is then used when making requests with Vault.

Authentication in Vault is pluggable via authentication backends. This
allows you to authenticate with Vault using a method that works best for your
organization. For example, you can authenticate using GitHub, certs, etc.

## Authentication Backends

There are many authentication backends available for Vault. They
are enabled using `vault auth-enable`. After they're enabled, you can
learn more about them using `vault path-help auth/<name>`. For example,
if you enable GitHub, you can use `vault path-help auth/github` to learn more
about how to configure it and login.

Multiple authentication backends can be enabled, but only one is required
to gain authentication. It is not currently possible to force a user through
multiple authentication backends to gain access.

This allows you to enable human-friendly as well as machine-friendly
backends at the same time. For example, for humans you might use the
"github" auth backend, and for machines you might use the "approle" backend.

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

#### Via the CLI

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

#### Via the API

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
