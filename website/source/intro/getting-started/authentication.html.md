---
layout: "intro"
page_title: "Authentication - Getting Started"
sidebar_current: "gettingstarted-auth"
description: |-
  Authentication to Vault gives a user access to use Vault. Vault can authenticate using multiple methods.
---

# Authentication

Now that we know how to use the basics of Vault, it is important to understand
how to authenticate to Vault itself. Up to this point, we haven't had to
authenticate because starting the Vault server in dev mode automatically logs
us in as the root user. In practice, you'll almost always have to manually authenticate.

On this page, we'll talk specifically about _authentication_. On the
next page, we talk about
[_authorization_](/intro/getting-started/policies.html). Authentication is the
mechanism of assigning an identity to a Vault user. The access control
and permissions associated with an identity are authorization, and will
not be covered on this page.

Vault has pluggable authentication backends, making it easy to authenticate
with Vault using whatever form works best for your organization. On this page
we'll use the token backend as well as the GitHub backend.

## Tokens

We'll first explain token authentication before going over any other
authentication backends. Token authentication is enabled by default in
Vault and cannot be disabled. It is also what we've been using up to this
point.

When you start a dev server with `vault server -dev`, it outputs your
_root token_. The root token is the initial access token to configure Vault.
It has root privileges, so it can perform any operation within Vault.
We'll cover how to limit privileges in the next section.

You can create more tokens using `vault token-create`:

```
$ vault token-create
Key             Value
token           c2c2fbd5-2893-b385-6fa5-30050439f698
token_accessor  0c1c3317-3d58-17e5-c1a9-3f54fa26610e
token_duration  0
token_renewable true
token_policies  [root]
```

By default, this will create a child token of your current token that
inherits all the same policies. The "child" concept here
is important: tokens always have a parent, and when that parent token is
revoked, children can also be revoked all in one operation. This makes it
easy when removing access for a user, to remove access for all sub-tokens
that user created as well.

After a token is created, you can revoke it with `vault token-revoke`:

```
$ vault token-revoke c2c2fbd5-2893-b385-6fa5-30050439f698
Success! Token revoked if it existed.
```

In a previous section, we use the `vault revoke` command. This command
is only used for revoking _secrets_. For revoking _tokens_, the
`vault token-revoke` command must be used.

To authenticate with a token, use the `vault auth` command:

```
$ vault auth d08e2bd5-ffb0-440d-6486-b8f650ec8c0c
Successfully authenticated! The policies that are associated
with this token are listed below:

root
```

This authenticates with Vault. It will verify your token and let you know
what access policies the token is associated with. If you want to test
`vault auth`, make sure you create a new token first.

## Auth Backends

In addition to tokens, other authentication, or _auth_, backends can be enabled.
Auth backends enable alternate methods of identifying with
Vault.  These identities are tied back to a set of access policies, just
like tokens. For example, for desktop environments, private key or
GitHub based authentication are available. For server environments, some
shared secret may be best. Auth backends give you flexibility
to choose what authentication you want to use.

As an example, let's authenticate using GitHub. First, enable the
GitHub authentication backend:

```
$ vault auth-enable github
Successfully enabled 'github' at 'github'!
```

Auth backends are mounted, just like secret backends, except auth
backends are always prefixed with `auth/`. So the GitHub backend we just
mounted can be accessed at `auth/github`. You can use `vault path-help` to
learn more about it.

With the GitHub backend enabled, we first have to configure it. For GitHub,
we tell it what organization users must be a part of, and map a team to a policy:

```
$ vault write auth/github/config organization=hashicorp
Success! Data written to: auth/github/config

$ vault write auth/github/map/teams/default value=default
Success! Data written to: auth/github/map/teams/default
```

The above configured our GitHub backend to only accept users from the
`hashicorp` organization (you should fill in your own organization)
and to map any team to the `default` policy, which is a built-in policy and is
the only policy (other than `root`) we have right now until the next section.

With GitHub enabled, we can now authenticate using `vault auth`:

```
$ vault auth -method=github token=e6919b17dd654f2b64e67b6369d61cddc0bcc7d5
Successfully authenticated! The policies that are associated
with this token are listed below:

default
```

Success! We've authenticated using GitHub. The `default` policy was associated
with my identity since we mapped that earlier. The value for `token` should be
your own [personal access
token](https://help.github.com/articles/creating-an-access-token-for-command-line-use/).

At this point, if you're following along, re-authenticate with the root token
from earlier (using `vault auth <token>`) to run the next commands.

You can revoke authentication from any authentication backend using
`vault token-revoke` as well, which can revoke any path prefix. For
example, to revoke all GitHub tokens, you could run the following.

```
$ vault token-revoke -mode=path auth/github
```

When you're done, you can disable authentication backends with
`vault auth-disable`. This will immediately invalidate all authenticated
users from this backend.

```
$ vault auth-disable github
Disabled auth provider at path 'github'!
```

## Next

In this page you learned about how Vault authenticates users. You learned
about the built-in token system as well as enabling other authentication
backends. At this point you know how Vault assigns an _identity_ to
a user.

The multiple authentication backends Vault provides let you choose the
most appropriate authentication mechanism for your organization.

In this next section, we'll learn about
[authorization and policies](/intro/getting-started/policies.html).
