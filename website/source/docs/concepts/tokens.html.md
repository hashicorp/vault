---
layout: "docs"
page_title: "Tokens"
sidebar_current: "docs-concepts-tokens"
description: |-
  Tokens are a core authentication method in Vault. Child tokens, token-based revocation, and more.
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
all external authentication mechanisms such as GitHub map down to dynamically
created tokens. These tokens have all the same properties as a normal
manually created token.

On this page, we'll show you how to create and manage tokens.

## Token Creation

Tokens are created via the API or using `vault token-create` from the CLI.
This will create a new token that is a child of the currently authenticated
token. As a child, the new token will automatically be revoked if the parent
is revoked.

If you're logged in as root, you can create an _orphan_ token by
specifying the `-orphan` flag. An orphan token has no parent, and therefore
when your token is revoked, it will not revoke the orphan.

Metadata associated with the token with `-metadata` is used to annotate
the token with information that is added to the audit log.

Finally, the `-policy` flag can be used to set the policies associated
with the token. Learn more about policies on the
[policies concepts](/docs/concepts/policies.html) page.

<a name="ttls-and-leases" />
## Token Time-To-Live and Leases

Every non-root token has a time-to-live (TTL) associated with it, which is
the number of seconds the token is valid for from this moment in time.
The TTL is constantly decreasing, and Vault will revoke the token when the TTL hits 0.
You can find out how much of a token's TTL remains using the `vault token-lookup` command.

When a token is revoked, any leases associated with the token will also be revoked,
even if the TTLs on the individual leases are not yet up. For example,
if a developer uses a token to generate AWS access keys, the AWS access keys
will be revoked when the token is revoked. The developer will then need to
create a new token, and use it to request a new set of AWS access keys.

To stop your token being revoked, you should call the `vault
token-renew` command periodically, which will attempt to extend the token's TTL.
You can specify _the new TTL_ using the `increment` argument. If this is
not specified, then Vault will use the default TTL specified by the
`auth/token` mount, or the global [`default_lease_ttl`](/docs/config/index.html#default_lease_ttl).

**NOTE:** Vault enforces an upper bound on how long non-root tokens can exist.
This bound is called the "max lease TTL", and is measured from
the moment the token was created. Vault will not allow a token to be renewed
beyond this point. If you attempt to do so, Vault will cap the TTL at a value
that will not cause the token to exist beyond the max lease TTL.

You can detect if the token has hit the max lease TTL by comparing the token's
TTL after renewal to the `increment` you provided when renewing it. If the TTL
is less than half of the requested `increment` then you should generate a new token.

As an example, suppose Vault's "max lease TTL" has been set to 10 hours,
and you have a token that was created 6 hours ago.
If you were to renew this token with an `increment` of 9 hours,
that would cause the token to exist for longer than the "max lease TTL" of
10 hours (`6 + 9 = 15`), so instead Vault will ignore `increment` and
cap the TTL at 4 hours (`10 - 6 = 4`). 4 hours is less than half of the requested
`increment`, so you should generate a new token and prepare for the fact that
Vault will soon be revoking the old one.

You can configure this cap by changing the system's
[`max_lease_ttl`](/docs/config/index.html#max_lease_ttl) (which will affect
all backends, such as AWS, PKI etc.), but the preferred method is to tune
the `auth/token` mount itself. Configuration options set on the `auth/token`
backend will take precedence over values set in the server config.
Here's an example:

```console
# We don't specify a TTL, so Vault uses the configured `default_lease_ttl`
$ vault token-create -policy=default
Key             Value
token           db7b286c-30b6-9e3f-ca2e-30d46aeedbcd
token_duration  2592000
token_renewable true

$ vault mount-tune -max-lease-ttl="2h" auth/token

# Vault tries to renew the token with an increment equal to the default
# lease TTL, but the `auth/token` backend is now restricting the overall lifetime of
# tokens to 2 hours (7200 seconds), so the TTL is restricted
# accordingly.
$ vault token-renew db7b286c-30b6-9e3f-ca2e-30d46aeedbcd
Key             Value
token           db7b286c-30b6-9e3f-ca2e-30d46aeedbcd
token_duration  7157
```
