---
layout: "docs"
page_title: "Vault Enterprise Replication"
sidebar_current: "docs-vault-enterprise-replication"
description: |-
  Vault Enterprise has support for Replication, allowing critical data to be replicated across clusters to support horizontally scaling workloads.

---

# Vault Replication

## Overview

Many organizations have infrastructure that spans multiple datacenters. Vault
provides the critical services of identity management, secrets storage, and
policy management.  This functionality is expected to be highly available and
to scale as the number of clients and their functional needs increase; at the
same time, operators would like to ensure that a common set of policies are
enforced globally, and a consistent set of secrets and keys are exposed to
applications that need to interoperate. Vault replication addresses both of
these needs.

Prior to Vault 0.7, Vault nodes could only be paired within clusters using a
common HA storage backend, where a single active node is backed up by a set of
standby nodes to ensure high availability. With replication, Vault replicates
large sets of its data to other nodes (or other HA clusters) to allow
horizontal scalability between clusters across geographically distributed data
centers.

## Architecture

Multiple Vault clusters communicate in a one-to-many near real-time flow.

The primary cluster acts as the system or record and asynchronously replicates
most Vault data to a series of remote clusters, known as secondary clusters or
secondaries.

The secondaries keep track of their own tokens and leases but share the
underlying configuration, policies, and supporting secrets (K/V values,
encryption keys for `transit`, etc). If a user action would modify underlying
shared state, the secondary forwards the request to the primary to be handled;
this is transparent to the client. In practice, most high-volume workloads
(reads in the `generic` backend, encryption/decryption operations in `transit`,
etc.) can be satisfied by the local secondary, allowing Vault to scale
relatively horizontally with the number of secondaries rather than vertically
as in the past.

All communication between primaries and secondaries is end-to-end encrypted
with mutually-authenticated TLS session, setup via replication tokens which are
exchanged during bootstrapping.

## What Is Replicated?

The data replicated in Vault 0.7 includes:

 * Secrets
 * Policies
 * Configuration details for secret backends
 * Configuration details for authentication backends
 * Configuration details for audit backends

Note that secret, authentication, and audit backends can be marked "local"
which prevents them from being affected by replication.

Access tokens for secrets are not a part of the replication process, as tokens
are local to a cluster that has generated them. Similarly, dynamic secrets
(database credentials, etc.) are issued by and their leases are tracked by each
cluster.

By not replicating leases and tokens, we avoid having each client call be
forwarded to the primary to check token validity, which enables scaling
horizontally.

## Internals

Details on the internal design of the replication feature can be found in the
[replication
internals](/docs/internals/replication.html)
document.

## Security Model

Vault is trusted all over the world to keep secrets safe. As such, we have put
extreme focus to detail to our replication model as well.

### Primary/Secondary Communication

When a cluster is marked as the primary it generates a self-signed CA
certificate. On request, and given a user-specified identifier, the primary
uses this CA certificate to generate a private key and certificate and packages
these, along with some other information, into a replication bootstrapping
bundle, a.k.a. a secondary activation token. The certificate is used to perform
TLS mutual authentication between the primary and that secondary.

This CA certificate is never shared with secondaries, and no secondary ever has
access to any other secondary’s certificate. In practice this means that
revoking a secondary’s access to the primary does not allow it continue
replication with any other machine; it also means that if a primary goes down,
there is full administrative control over which cluster becomes primary. An
attacker cannot spoof a secondary into believing that a cluster the attacker
controls is the new primary without also being able to administratively direct
the secondary to connect by giving it a new bootstrap package (which is an
ACL-protected call).

Vault makes use of Application Layer Protocol Negotiation on its cluster port.
This allows the same port to handle both request forwarding and replication,
even while keeping the certificate root of trust and feature set different.

### Secondary Activation Tokens

A secondary activation token is an extremely sensitive item and as such is
protected via response wrapping. Experienced Vault users will note that the
wrapping format for replication bootstrap packages is different from normal
response wrapping tokens: it is a signed JWT. This allows the replication token
to carry the redirect address of the primary cluster as part of the token. In
most cases this means that simply providing the token to a new secondary is
enough to activate replication, although this can also be overridden when the
token is provided to the secondary.

Secondary activation tokens should be treated like Vault root tokens. If
disclosed to a bad actor, that actor can gain access to all Vault data. It
should therefore be treated with utmost sensitivity.  Like all
response-wrapping tokens, once the token is used successfully (in this case, to
activate a secondary) it is useless, so it is only necessary to safeguard it
from one machine to the next.  Like with root tokens, HashiCorp recommends that
when a secondary activation token is live, there are multiple eyes on it from
generation until it is used.

Once a secondary is activated, its cluster information is stored safely behind
its encrypted barrier.

## Setup and Best Practices

A [setup guide](/docs/guides/replication.html) is
available to help you get started; this guide also contains best practices
around operationalizing the replication feature.

## API

The Vault replication component has a full HTTP API. Please see the
[Vault Replication API](/api/system/replication.html) for more
details.
