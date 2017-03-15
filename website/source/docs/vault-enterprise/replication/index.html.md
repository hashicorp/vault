---
layout: "docs"
page_title: "Vault Enterprise HSM Support"
sidebar_current: "docs-vault-enterprise"
description: |-
  Vault Enterprise has support for Replication, allowing critical data to be replicated across clusters. 

---

# Vault Replication

## Overview

Many organizations have infrastructure that spans multiple datacenters. Vault provides 
the critical services of identity management, secrets storage, and policy management. 
This functionality is expected to be highly available and operators would like to ensure
a common set of policies are enforced globally, and a consistent set of secrets and keys
are exposed to applications that need to interoperate. Vault replication enables
redundancy across datacenters.

Prior to Vault 0.7, Vault nodes could only be paired within clusters using a common HA
storage backend. A single node is locked with a set of other nodes in a cluster to ensure
high availability. With replication, we intend to allow Vault to allow horizontal 
scalability between clusters a cross geographically distributed data centers.

## Architecture

When replicating clusters, a single cluster is designed the primary cluster. The
primary cluster is the active cluster in use within an infrastructure, which 
serves as the system  of record for replication.

The primary cluster asynchronously replicates data to a series of remote standby
clusters. These clusters are known as secondary clusters or secondaries. 

Roles for primaries and secondaries, as well as authentication for each and the process
of replicating data between them, is setup using replication tokens which are exchanged
during bootstrapping. Once setup, replication uses end-to-end TLS, and Vault manages the
lifecycle of certificates.

## What Is Replicated?

The data replicated in Vault 0.7 will include:

 * Secrets
 * Policies
 * Configuration details for secret backends
 * Configuration details for authentication backends
 * Configuration details for audit backends

Access tokens for secrets are not replicated during the replication process, as tokens are 
local to a cluster that has generated them.  

## Activating Replication 

### Activating the Primary

To activate the primary, run vault write -f sys/replication/primary/enable. 

There is currently one optional argument: primary_cluster_addr. This can be used to override
the cluster address that the primary advertises to the secondary, in case the internal network
address/pathing is different between members of a single cluster and primary/secondary clusters.

### Fetching a Secondary Token

To fetch a secondary bootstrap token, run vault write sys/replication/primary/secondary-token id=<id>. 

The value for ID is opaque to Vault and can be any identifying value you want; this can be used
later to revoke the secondary and will be listed when you read replication status on the 
primary. You will get back a normal wrapped response, except that the token will be a JWT instead
of UUID-formatted random bytes. 

### Activating a Secondary

To activate a secondary, run vault write sys/replication/secondary/enable token=<token>. 

You must provide the full token value. Be very careful when running this command, as it will 
destroy all data currently stored in the secondary.

There is an optional argument, primary_api_addr, which can be used to override the API address
of the primary cluster; otherwise the secondary will use the value embedded in the bootstrap 
token, which is the primary’s redirect address.

Once the secondary is activated and has bootstrapped, it will be ready for service and will 
maintain state with the primary. It is safe to seal/shutdown the primary and/or secondary; when
both are available again, they will synchronize back into a replicated state.

Note: if the secondary is in an HA cluster, you will need to ensure that each standby is 
sealed/unsealed with the new (primary’s) unseal keys. If one of the standbys takes over on active
duty before this happens it will seal itself to remove it from rotation (e.g. if using Consul for
service discovery), but if a standby does not attempt taking over it will throw errors. We plan
to make this workflow better in a future update.

### Dev-Mode Root Tokens 

To ease development and testing, when both the primary and secondary are running in development 
mode, the initial root token created by the primary (including those with custom IDs specified with
 -dev-root-token-id) will be populated into the secondary upon activation. This allows a developer
to keep a consistent ~/.vault-token file or VAULT_TOKEN environment variable when working with both
clusters.

On a production system, after a secondary is activated, the enabled authentication backends should
be used to get tokens with appropriate policies as policies and auth backend configuration are replicated.

The generate-root command can be also be used to generate a root token local to the secondary cluster.

### Managing Vault Replication 

Note: this section describes the replication model in the initial release of the replication feature.
It is possible that there may be more replication modes in the future.

Vault’s current replication model is intended to allow horizontally scaling Vault’s functions rather
than to act in a strict Disaster Recovery (DR) capacity. As a result, Vault replication acts on static
items within Vault, meaning information that is not part of Vault’s lease-tracking system. In a practical
sense, this means that all Vault information is replicated from the primary to secondaries except for 
tokens and secret leases.

Because token information must be checked and possibly rewritten with each use (e.g. to decrement its
use count), replicated tokens would require every call to be forwarded to the primary, decreasing 
rather than increasing total Vault throughput.

Secret leases are tracked independently for two reasons: one, because every such lease is tied to a 
token and tokens are local to each cluster; and two, because tracking large numbers of leases is 
memory-intensive and tracking all leases in a replicated fashion could dramatically increase the memory
 requirements across all Vault nodes.

We believe that this replication model provides significant utility and the benefits of horizontally
scaling Vault’s functionality dramatically outweigh the drawbacks of not providing a full DR-ready system.
However, it does mean that certain principles must be kept in mind.

### Always Use the Local Cluster

First and foremost, when designing systems to take advantage of replicated Vault, you must ensure
that they always use the same Vault cluster for all operations, as only that cluster will know about
the client’s Vault token.

### Enabling a Secondary Wipes Storage

Replication relies on having a shared keyring between primary and secondaries and also relies on
having a shared understanding of the data store state. As a result, when replication is enabled,
all of the secondary’s existing storage will be wiped. This is irrevocable. Make a backup first
if there is a remote chance you’ll need some of this data at some future point.

Generally, activating as a secondary will be the first thing that is done upon setting up a new 
cluster for replication.

### Replicated vs. Local Backend Mounts

All backend mounts (of all types) that can be enabled within Vault default to being mounted as
a replicated mount. This means that mounts cannot be enabled on a secondary, and mounts enabled
on the primary will replicate to secondaries.

Mounts can also be marked local (via the -local flag on the Vault CLI or setting the local parameter
true in the API). This can only be performed at mount time; if a mount is local but should have
been replicated, or vice versa, you must unmount the backend and mount a new instance at that
path with the local flag enabled.

Local mounts do not propagate data from the primary to secondaries, and local mounts on secondaries
do not have their data removed during the syncing process.

### Audit Backends

If Vault has at least one audit backend configured and is unable to successfully log to at least
one backend, it will block further requests. 

Replicated audit mounts must be able to successfully log on all replicated clusters. For example,
if using the file backend, the configured path must be able to be written to by all secondaries.
It may be useful to use at least one local audit mount on each cluster to prevent such a scenario.

### Never Have Two Primaries

The replication model is not designed for active-active usage and enabling two primaries should
never be done, as it can lead to data loss if they or their secondaries are ever reconnected.

### Disaster Recovery
At the moment, because leases and tokens are not replicated, if you need true DR, you will need
a DR solution per cluster (similar to non-replicated Vault). 

Local backend mounts are not replicated and their use will require existing DR mechanisms if DR
is necessary in your implementation. 

We may pursue a dedicated Disaster Recovery-focused Replication Mode at a future time. 

## Security Model

Vault is trusted all over the world to keep secrets safe. As such, we have put extreme focus
to detail to our replication model as well.

### Primary/Secondary Communication

When a cluster is marked as the primary it generates a self-signed CA certificate. On request,
and given a user-specified identifier, the primary uses this CA certificate to generate a private
key and certificate and packages these, along with some other information, into a replication
bootstrapping bundle, a.k.a. a secondary activation token. The certificate is used to perform TLS
mutual authentication between the primary and that secondary.

This CA certificate is never shared with secondaries, and no secondary ever has access to any other
secondary’s certificate. In practice this means that revoking a secondary’s access to the primary 
does not allow it continue replication with any other machine; it also means that if a primary goes
down, there is full administrative control over which cluster becomes primary. An attacker cannot 
spoof a secondary into believing that a cluster the attacker controls is the new primary without 
also being able to administratively direct the secondary to connect by giving it a new bootstrap 
package (which is an ACL-protected call).

Vault makes use of Application Layer Protocol Negotiation on its cluster port. This allows the same
port to handle both request forwarding and replication, even while keeping the certificate root of
trust and feature set different.

### Secondary Activation Tokens

A secondary activation token is an extremely sensitive item and as such is protected via response
wrapping. Experienced Vault users will note that the wrapping format for replication bootstrap
packages is different from normal response wrapping tokens: it is a signed JWT. This allows the 
replication token to carry the redirect address of the primary cluster as part of the token. In 
most cases this means that simply providing the token to a new secondary is enough to activate 
replication, although this can also be overridden when the token is provided to the secondary.

Secondary activation tokens should be treated like Vault root tokens. If disclosed to a bad actor, 
that actor can gain access to all Vault data. It should therefore be treated with utmost sensitivity.
Like all response-wrapping tokens, once the token is used successfully (in this case, to activate 
a secondary) it is useless, so it is only necessary to safeguard it from one machine to the next. 
Like with root tokens, HashiCorp recommends that when a secondary activation token is live, there 
are multiple eyes on it from generation until it is used.

Once a secondary is activated, its cluster information is stored safely behind its encrypted barrier.





