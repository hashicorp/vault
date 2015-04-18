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

## Token Leases

Every token has a lease associated with it. These leases behave in much
the same way as [leases for secrets](/docs/concepts/lease.html). After
the lease period is up, the token will no longer function. In addition
to no longer functioning, Vault will revoke it.

In order to avoid your token being revoked, the `vault token-renew`
command should be used to renew the lease on the token periodically.

After a token is revoked, all of the secrets in use by that token will
also be revoked. Therefore, if a user requests AWS access keys, for example,
then after the token expires the AWS access keys will also be expired even
if they had remaining lease time.
