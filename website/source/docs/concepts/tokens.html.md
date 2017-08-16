---
layout: "docs"
page_title: "Tokens"
sidebar_current: "docs-concepts-tokens"
description: |-
  Tokens are a core authentication method in Vault. Concepts and important features.
---

# Tokens

Tokens are the core method for _authentication_ within Vault. Tokens
can be used directly or [authentication backends](/docs/concepts/auth.html)
can be used to dynamically generate tokens based on external identities.

If you've gone through the getting started guide, you probably noticed that
`vault server -dev` (or `vault init` for a non-dev server) outputs an
initial "root token." This is the first method of authentication for Vault.
It is also the only authentication backend that cannot be disabled.

As stated in the [authentication concepts](/docs/concepts/auth.html),
all external authentication mechanisms, such as GitHub, map down to dynamically
created tokens. These tokens have all the same properties as a normal manually
created token.

Within Vault, tokens map to information. The most important information mapped
to a token is a set of one or more attached
[policies](/docs/concepts/policies.html). These policies control what the token
holder is allowed to do within Vault. Other mapped information includes
metadata that can be viewed and is added to the audit log, creation time, last
renewal time, and more.

Read on for a deeper dive into token concepts.

## Token Concepts

### The Token Store

Often in documentation or in help channels, the "token store" is referenced.
This is the same as the [`token` authentication
backend](/docs/auth/token.html). This is a special
backend in that it is responsible for creating and storing tokens, and cannot
be disabled. It is also the only authentication backend that has no login
capability -- all actions require existing authenticated tokens.

### Root Tokens

Root tokens are tokens that have the `root` policy attached to them. Root
tokens can do anything in Vault. _Anything_. In addition, they are the only
type of token within Vault that can be set to never expire without any renewal
needed. As a result, it is purposefully hard to create root tokens; in fact, as
of version 0.6.1, there are only three ways to create root tokens:

1. The initial root token generated at `vault init` time -- this token has no
   expiration
2. By using another root token; a root token with an expiration cannot create a
   root token that never expires
3. By using `vault generate-root` ([example](/guides/generate-root.html))
   with the permission of a quorum of unseal key holders

Root tokens are useful in development but should be extremely carefully guarded
in production. In fact, the Vault team recommends that root tokens are only
used for just enough initial setup (usually, setting up authentication backends
and policies necessary to allow administrators to acquire more limited tokens)
or in emergencies, and are revoked immediately after they are no longer needed.
If a new root token is needed, the `generate-root` command and associated [API
endpoint](/api/system/generate-root.html) can be
used to generate one on-the-fly.

It is also good security practice for there to be multiple eyes on a terminal
whenever a root token is live. This way multiple people can verify as to the
tasks performed with the root token, and that the token was revoked immediately
after these tasks were completed.

### Token Hierarchies and Orphan Tokens

Normally, when a token holder creates new tokens, these tokens will be created
as children of the original token; tokens they create will be children of them;
and so on. When a parent token is revoked, all of its child tokens -- and all
of their leases -- are revoked as well. This ensures that a user cannot escape
revocation by simply generating a never-ending tree of child tokens.

Often this behavior is not desired, so users with appropriate access can create
`orphan` tokens. These tokens have no parent -- they are the root of their own
token tree. These orphan tokens can be created:

1. Via the `auth/token/create-orphan` endpoint
2. By having `sudo` capability or `root` policy when accessing
   `auth/token/create` and setting the `orphan` parameter to `true`
3. Via token store roles
4. By logging in with any other (non-`token`) authentication backend

Users with appropriate permissions can also use the `auth/token/revoke-orphan`
endpoint, which revokes the given token but rather than revoke the rest of the
tree, it instead sets the tokens' immediate children to be orphans. Use with
caution!

### Token Accessors

When tokens are created, a token accessor is also created and returned. This
accessor is a value that acts as a reference to a token and can only be used to
perform limited actions:

1. Look up a token's properties (not including the actual token ID)
2. Look up a token's capabilities on a path
3. Revoke the token

There are many useful workflows around token accessors. As an example, a
service that creates tokens on behalf of another service (such as the
[Nomad](https://www.nomadproject.io/) scheduler) can store the accessor
correlated with a particular job ID. When the job is complete, the accessor can
be used to instantly revoke the token given to the job and all of its leased
credentials, limiting the chance that a bad actor will discover and use them.

Audit backends can optionally be set to not obfuscate token accessors in audit
logs. This provides a way to quickly revoke tokens in case of an emergency.
However, it also means that the audit logs can be used to perform a
larger-scale denial of service attack.

Finally, the only way to "list tokens" is via the `auth/token/accessors`
command, which actually gives a list of token accessors. While this is still a
dangerous endpoint (since listing all of the accessors means that they can then
be used to revoke all tokens), it also provides a way to audit and revoke the
currently-active set of tokens.

### Token Time-To-Live, Periodic Tokens, and Explicit Max TTLs

Every non-root token has a time-to-live (TTL) associated with it, which is a
current period of validity since either the token's creation time or last
renewal time, whichever is more recent. (Root tokens may have a TTL associated,
but the TTL may also be 0, indicating a token that never expires). After the
current TTL is up, the token will no longer function -- it, and its associated
leases, are revoked.

If the token is renewable, Vault can be asked to extend the token validity
period using `vault token-renew` or the appropriate renewal endpoint. At this
time, various factors come into play. What happens depends upon whether the
token is a periodic token (available for creation by `root`/`sudo` users, token
store roles, or some authentication backends), has an explicit maximum TTL
attached, or neither.

#### The General Case

In the general case, where there is neither a period nor explicit maximum TTL
value set on the token, the token's lifetime since it was created will be
compared to the maximum TTL. This maximum TTL value is dynamically generated
and can change from renewal to renewal, so the value cannot be displayed when a
token's information is looked up. It is based on a combination of factors:

1. The system max TTL, which is 32 days but can be changed in Vault's
   configuration file
2. The max TTL set on a mount using [mount
   tuning](/api/system/mounts.html). This value
   is allowed to override the system max TTL -- it can be longer or shorter,
   and if set this value will be respected.
3. A value suggested by the authentication backend that issued the token. This
   might be configured on a per-role, per-group, or per-user basis. This value
   is allowed to be less than the mount max TTL (or, if not set, the system max
   TTL), but it is not allowed to be longer.

Note that the values in (2) and (3) may change at any given time, which is why
a final determination about the current allowed max TTL is made at renewal time
using the current values. It is also why it is important to always ensure that
the TTL returned from a renewal operation is within an allowed range; if this
value is not extending, likely the TTL of the token cannot be extended past its
current value and the client may want to reauthenticate and acquire a new
token. However, outside of direct operator interaction, Vault will never revoke
a token before the returned TTL has expired.

#### Explicit Max TTLs

Tokens can have an explicit max TTL set on them. This value becomes a hard
limit on the token's lifetime -- no matter what the values in (1), (2), and (3)
from the general case are, the token cannot live past this explicitly-set
value. This has an effect even when using periodic tokens to escape the normal
TTL mechanism.

#### Periodic Tokens

In some cases, having a token be revoked would be problematic -- for instance,
if a long-running service needs to maintain its SQL connection pool over a long
period of time. In this scenario, a periodic token can be used. Periodic tokens
can be created in a few ways:

1. By having `sudo` capability or a `root` token with the `auth/token/create`
   endpoint
2. By using token store roles
3. By using an authentication backend that supports issuing these, such as
   AppRole

At issue time, the TTL of a periodic token will be equal to the configured
period. At every renewal time, the TTL will be reset back to this configured
period, and as long as the token is successfully renewed within each of these
periods of time, it will never expire. Outside of `root` tokens, it is
currently the only way for a token in Vault to have an unlimited lifetime.

The idea behind periodic tokens is that it is easy for systems and services to
perform an action relatively frequently -- for instance, every two hours, or
even every five minutes. Therefore, as long as a system is actively renewing
this token -- in other words, as long as the system is alive -- the system is
allowed to keep using the token and any associated leases. However, if the
system stops renewing within this period (for instance, if it was shut down),
the token will expire relatively quickly. It is good practice to keep this
period as short as possible, and generally speaking it is not useful for humans
to be given periodic tokens.

There are a few important things to know when using periodic tokens:

* When a periodic token is created via a token store role, the _current_ value of the role's period setting will be used at renewal time
* A token with both a period and an explicit max TTL will act like a periodic token but will be revoked when the explicit max TTL is reached
