---
layout: "docs"
page_title: "Lease, Renew, and Revoke"
sidebar_current: "docs-concepts-lease"
description: |-
  Vault provides a lease with every secret. When this lease is expired, Vault will revoke that secret.
---

# Lease, Renew, and Revoke

With every dynamic secret and authentication token, Vault creates a _lease_:
metadata containing information such as a time duration, renewability, and
more. Vault promises that the data will be valid for the given duration, or
Time To Live (TTL). Once the lease is expired, Vault can automatically revoke
the data, and the consumer of the secret can no longer be certain that it is
valid.

The benefit should be clear: consumers of secrets need to check in with
Vault routinely to either renew the lease (if allowed) or request a
replacement secret. This makes the Vault audit logs more valuable and
also makes key rolling a lot easier.

All dynamic secrets in Vault are required to have a lease. Even if the data is
meant to be valid for eternity, a lease is required to force the consumer
to check in routinely.

In addition to renewals, a lease can be _revoked_. When a lease is revoked, it
invalidates that secret immediately and prevents any further renewals. For
example, with the [AWS secret backend](/docs/secrets/aws/index.html), the
access keys will be deleted from AWS the moment a secret is revoked. This
renders the access keys invalid from that point forward.

Revocation can happen manually via the API, via the `vault revoke` cli command,
or automatically by Vault. When a lease is expired, Vault will automatically
revoke that lease.

**Note**: The [Key/Value Backend](/docs/secrets/kv/index.html) which stores
arbitrary secrets does not issue leases.

## Lease IDs

When reading a dynamic secret, such as via `vault read`, Vault always returns a
`lease_id`. This is the ID used with commands such as `vault renew` and `vault
revoke` to manage the lease of the secret.

## Lease Durations and Renewal

Along with the lease ID, a _lease duration_ can be read. The lease duration is
a Time To Live value: the time in seconds for which the lease is valid.  A
consumer of this secret must renew the lease within that time.

When renewing the lease, the user can request a specific amount of time from
now to extend the lease. For example: `vault renew my-lease-id 3600` would
request to extend the lease of "my-lease-id" by 1 hour (3600 seconds).

The requested increment is completely advisory. The backend in charge of the
secret can choose to completely ignore it. For most secrets, the backend does
its best to respect the increment, but often limits it to ensure renewals every
so often.

As a result, the return value of renewals should be carefully inspected to
determine what the new lease is.

**Note**: Prior to version 0.3, Vault documentation and help text did not
distinguish sufficiently between a _lease_ and a _lease duration_.  Starting
with version 0.3, Vault will start migrating to the term _ttl_ to describe
lease durations, at least for user-facing text. As _lease duration_ is still a
legitimate (but more verbose) description, there are currently no plans to
change the JSON key used in responses, in order to retain
backwards-compatibility.

## Prefix-based Revocation

In addition to revoking a single secret, operators with proper access control
can revoke multiple secrets based on their lease ID prefix.

Lease IDs are structured in a way that their prefix is always the path where
the secret was requested from. This lets you revoke trees of secrets. For
example, to revoke all AWS access keys, you can do `vault revoke -prefix aws/`.

This is very useful if there is an intrusion within a specific system: all
secrets of a specific backend or a certain configured backend can be revoked
quickly and easily.
