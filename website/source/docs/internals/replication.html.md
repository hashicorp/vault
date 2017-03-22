---
layout: "docs"
page_title: "Replication"
sidebar_current: "docs-internals-replication"
description: |-
  Learn about the details of multi-datacenter replication within Vault.
---

# Replication (Vault Enterprise)

Vault Enterprise 0.7 adds support for multi-datacenter replication. Before
using this feature, it is useful to understand the intended use cases, design
goals, and high level architecture.

Replication is based on a primary/secondary (1:N) model with asynchronous
replication, focusing on high availability for global deployments. The
trade-offs made in the design and implementation of replication reflect these
high level goals.

# Use Cases

Vault replication is based on a number of common use cases:

* **Multi-Datacenter Deployments**: A common challenge is providing Vault to
  applications across many datacenters in a highly-available manner. Running a
  single Vault cluster imposes high latency of access for remote clients,
  availability loss or outages during connectivity failures, and limits
  scalability.

* **Backup Sites**: Implementing a robust business continuity plan around the
  loss of a primary datacenter requires the ability to quickly and easily fail
  to a hot backup site.

* **Scaling Throughput**: Applications that use Vault for
  Encryption-as-a-Service or cryptographic offload may generate a very high
  volume of requests for Vault. Replicating keys between multiple clusters
  allows load to be distributed across additional servers to scale request
  throughput.

# Design Goals

Based on the use cases for Vault Replication, we had a number of design goals
for the implementation:

* **Availability**: Global deployments of Vault require high levels of
  availability, and can tolerate reduced consistency. During full connectivity,
  replication is nearly real-time between the primary and secondary clusters.
  Degraded connectivity between a primary and secondary does not impact the
  primary's ability to service requests, and the secondary will continue to
  service reads on last-known data.

* **Conflict Free**: Certain replication techniques allow for potential write
  conflicts to take place. Particularly, any active/active configuration where
  writes are allowed to multiple sites require a conflict resolution strategy.
  This varies from techniques that allow for data loss like last-write-wins, or
  techniques that require manual operator resolution like allowing multiple
  values per key. We avoid the possibility of conflicts to ensure there is no
  data loss or manual intervention required.

* **Transparent to Clients**: Vault replication should be transparent to
  clients of Vault, so that existing thin clients work unmodified. The Vault
  servers handle the logic of request forwarding to the primary when necessary,
  and multi-hop routing is performed internally to ensure requests are
  processed.

* **Simple to Operate**: Operating a replicated cluster should be simple to
  avoid administrative overhead and potentially introducing security gaps.
  Setup of replication is very simple, and secondaries can handle being
  arbitrarily behind the primary, avoiding the need for operator intervention
  to copy data or snapshot the primary.

# Architecture

The architecture of Vault replication is based on the design goals, focusing on
the intended use cases. When replication is enabled, a cluster is set as either
a _primary_ or _secondary_. The primary cluster is authoritative, and is the
only cluster allowed to perform actions that write to the underlying data
storage, such as modifying policies or secrets. Secondary clusters can service
all other operations, such as reading secrets or sending data through
`transit`, and forward any writes to the primary cluster. Disallowing multiple
primaries ensures the cluster is conflict free and has an authoritative state.

The primary cluster uses log shipping to replicate changes to all of the
secondaries.  This ensures writes are visible globally in near real-time when
there is full network connectivity. If a secondary is down or unable to
communicate with the primary, writes are not blocked on the primary and reads
are still serviced on the secondary. This ensures the availability of Vault.
When the secondary is initialized or recovers from degraded connectivity it
will automatically reconcile with the primary.

Lastly, clients can speak to any Vault server without a thick client. If a
client is communicating with a standby instance, the request is automatically
forwarded to a active instance. Secondary clusters will service reads locally
and forward any write requests to the primary cluster. The primary cluster is
able to service all request types.

An important optimization Vault makes is to avoid replication of tokens or
leases between clusters. Policies and secrets are the minority of data managed
by Vault and tend to be relatively stable. Tokens and leases are much more
dynamic, as they are created and expire rapidly. Keeping tokens and leases
locally reduces the amount of data that needs to be replicated, and distributes
the work of TTL management between the clusters. The caveat is that clients
will need to re-authenticate if they switch the Vault cluster they are
communicating with.

# Implementation Details

It is important to understand the high-level architecture of replication to
ensure the trade-offs are appropriate for your use case. The implementation
details may be useful for those who are curious or want to understand more
about the performance characteristics or failure scenarios.

Using replication requires a storage backend that supports transactional
updates, such as Consul.  This allows multiple key/value updates to be
performed atomically. Replication uses this to maintain a
[Write-Ahead-Log][wal] (WAL) of all updates, so that the key update happens
atomically with the WAL entry creation.  The WALs are then used to perform log
shipping between the Vault clusters. When a secondary is closely synchronized
with a primary, Vault directly streams new WALs to be applied, providing near
real-time replication. A bounded set of WALs are maintained for the
secondaries, and older WALs are garbage collected automatically.

When a secondary is initialized or is too far behind the primary there may not
be enough WALs to synchronize. To handle this scenario, Vault maintains a
[merkle index][merkle] of the encrypted keys. Any time a key is updated or
deleted, the merkle index is updated to reflect the change.  When a secondary
needs to reconcile with a primary, they compare their merkle indexes to
determine which keys are out of sync. The structure of the index allows this to
be done very efficiently, usually requiring only two round trips and a small
amount of data. The secondary uses this information to reconcile and then
switches back into WAL streaming mode.

Performance is an important concern for Vault, so WAL entries are batched and
the merkle index is not flushed to disk with every operation. Instead, the
index is updated in memory for every operation and asynchronously flushed to
disk. As a result, a crash or power loss may cause the merkle index to become
out of sync with the underlying keys. Vault uses the [ARIES][aries] recovery
algorithm to ensure the consistency of the index under those failure
conditions.

Log shipping traditionally requires the WAL stream to be synchronized, which
can introduce additional complexity when a new primary cluster is promoted.
Vault uses the merkle index as the source of truth, allowing the WAL streams to
be completely distinct and unsynchronized.  This simplifies administration of
Vault Replication for operators.

# Caveats

* **Read-After-Write Consistency**: All write requests are forwarded from
  secondaries to the primary cluster in order to avoid potential conflicts.
  While replication is near real-time, it is not instantaneous, meaning there
  is a potential for a client to write to a secondary and a subsequent read to
  return an old value. Secondaries attempt to mask this from an individual
  client making subsequent requests by stalling write requests until the write
  is replicated or a timeout is reached (2 seconds). If the timeout is reached,
  the client will receive a warning.

* **Stale Reads**: Secondary clusters service reads based on their
  locally-replicated data. During normal operation updates from a primary are
  received in near real-time by secondaries. However, during an outage or
  network service disruption, replication may stall and secondaries may have
  stale data. The cluster will automatically recover and reconcile any stale
  data once the outage has recovered, but reads in the intervening period may
  receive stale data.

[wal]: https://en.wikipedia.org/wiki/Write-ahead_logging
[merkle]: https://en.wikipedia.org/wiki/Merkle_tree
[aries]: https://en.wikipedia.org/wiki/Algorithms_for_Recovery_and_Isolation_Exploiting_Semantics
