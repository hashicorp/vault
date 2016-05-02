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
the number of seconds the token is valid for, _measured from the moment it is
created_. After the TTL is up, the token will no longer function, and
Vault will revoke it. You can find how much of a token's TTL is remaining
by calling `vault token-lookup`.

When a token is revoked, any leases associated with the token will be revoked
as well, even if the TTLs on the individual leases are not yet up. For example,
if a user requests AWS access keys, after the token expires the AWS access keys
will also be revoked.

To prevent your token being revoked, you should call the `vault
token-renew` command periodically, which will attempt to extend the token's TTL.
You can specify _the new TTL_ using the `increment` argument. If this is
not specified, then Vault will use the default TTL specified by the
`auth/token` mount, or the global [`default_lease_ttl`](/docs/config/index.html#default_lease_ttl).

**NOTE:** Vault enforces a cap on the maximum lifetime of non-root tokens.
This cap is called the "Max Lease TTL", and is measured in seconds
from the moment the token was created. If renewing a token would cause
its TTL to exceed this value, then Vault will cap the TTL at the Max Lease TTL.

For example, suppose the system's Max Lease TTL has been set to 12 hours,
and you have a token that was created 7 hours ago.
If you were to renew this token with an `increment` of 6 hours,
that would cause the token's TTL to exceed the max of 12 hours (7 + 6 = 13),
so Vault will cap the TTL at 12 hours.

You can configure this cap by changing the system's 
[`max_lease_ttl`](/docs/config/index.html#max_lease_ttl) (which will affect
all backends, such as AWS, PKI etc.), but the preferred method is to tune
the `auth/token` mount itself:

```console
$ vault token-create -policy=default
Key             Value
token           db7b286c-30b6-9e3f-ca2e-30d46aeedbcd
token_duration  2592000
token_renewable true

$ vault mount-tune -max-lease-ttl="2h" auth/token

$ vault token-renew db7b286c-30b6-9e3f-ca2e-30d46aeedbcd
Key             Value
token           db7b286c-30b6-9e3f-ca2e-30d46aeedbcd
token_duration  7157
```
