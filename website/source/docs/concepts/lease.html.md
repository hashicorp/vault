---
layout: "docs"
page_title: "Lease, Renew, and Revoke"
sidebar_current: "docs-concepts-lease"
description: |-
  Vault provides a lease with every secret. When this lease is expired, Vault will revoke that secret.
---

# Lease, Renew, and Revoke

With every secret and authentication token, Vault provides a _lease_:
an amount of time that Vault promises that the data will be valid for.
Once the lease is up, Vault can automatically revoke the data, and the
consumer of the secret can no longer be certain that it is valid.

The benefit should be clear: consumers of secrets need to check in with
Vault routinely to either renew the lease (if allowed) or request a
replacement secret. This makes the Vault audit logs more valuable and
also makes key rolling a lot easier.

All secrets in Vault are required to have a lease. Even if the data is
meant to be valid for eternity, a lease is required to force the consumer
to check in routinely.

In addition to renewals, a lease can be _revoked_. When a lease is revoked,
it invalidates that secret immediately and prevents any further renewals.
For
[dynamic secrets](#),
the secrets themselves are often immediately disabled. For example, with
the
[AWS secret backend](/docs/secrets/aws/index.html), the access keys will
be deleted from AWS the moment a secret is revoked. This renders the access
keys invalid from that point forward.

Revocation can happen manually via the API or `vault revoke`, or automatically
by Vault. When a lease is expired, Vault will automatically revoke that
lease.

## Lease IDs

When reading a secret, such as via `vault read`, Vault always returns
a `lease_id`. This is the ID used with commands such as `vault renew` and
`vault revoke` to manage the lease of the secret.

## Lease Durations and Renewal

Along with the lease ID, a _lease duration_ can be read. The lease duration
is the time in seconds that the lease is valid for. A consumer of this
secret must renew the lease within that time.

When renewing the lease, the user can request a specific amount of time
from now to extend the lease. For example: `vault renew my-lease-id 3600`
would request to extend the lease of "my-lease-id" by 1 hour (3600 seconds).

The requested increment is completely advisory. The backend in charge
of the secret can choose to completely ignore it. For most secrets, the
backend does its best to respect the increment, but often limits it to
ensure renewals every so often.

As a result, the return value of renews should be carefully inspected
to determine what the new lease is.

## Prefix-based Revocation

In addition to revoking a single secret, operators with proper access
control can revoke multiple secrets based on their lease ID prefix.

Lease IDs are structured in a way that their prefix is always the path
where the secret was requested from. This lets you revoke trees of
secrets. For example, to revoke all AWS access keys, you can do
`vault revoke -prefix aws/`.

This is very useful if there is an intrusion within a specific system:
all secrets of a specific backend or a certain configured backend can
be revoked quickly and easily.
