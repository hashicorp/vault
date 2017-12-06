---
layout: "guides"
page_title: "Setting up Vault Enterprise Performance Replication - Guides"
sidebar_current: "guides-replication"
description: |-
  Learn how to set up and manage Vault Enterprise Performance Replication.
---

# Replication Setup &amp; Guidance

If you're unfamiliar with Vault Replication concepts, please first look at the
[general information page](/docs/vault-enterprise/replication/index.html). More
details can be found in the
[replication internals](/docs/internals/replication.html) document.

Vault replication also includes a complete API. For more information, please see
the [Vault Replication API documentation](/api/system/replication.html)


## Activating Performance Replication

### Activating the Primary

To activate the primary, run:


    $ vault write -f sys/replication/performance/primary/enable


There is currently one optional argument: `primary_cluster_addr`. This can be
used to override the cluster address that the primary advertises to the
secondary, in case the internal network address/pathing is different between
members of a single cluster and primary/secondary clusters.

### Fetching a Secondary Token

To fetch a secondary bootstrap token, run:


    $ vault write sys/replication/performance/primary/secondary-token id=<id>


The value for `id` is opaque to Vault and can be any identifying value you want;
this can be used later to revoke the secondary and will be listed when you read
replication status on the primary. You will get back a normal wrapped response,
except that the token will be a JWT instead of UUID-formatted random bytes.

### Activating a Secondary

To activate a secondary using the fetched token, run:


    $ vault write sys/replication/performance/secondary/enable token=<token>


You must provide the full token value. Be very careful when running this
command, as it will destroy all data currently stored in the secondary.

There are a few optional arguments, with the one you'll most likely need being
`primary_api_addr`, which can be used to override the API address of the
primary cluster; otherwise the secondary will use the value embedded in the
bootstrap token, which is the primary’s redirect address. If the primary has no
redirect address (for instance, if it's not in an HA cluster), you'll need to
set this value at secondary enable time.

Once the secondary is activated and has bootstrapped, it will be ready for
service and will maintain state with the primary. It is safe to seal/shutdown
the primary and/or secondary; when both are available again, they will
synchronize back into a replicated state.

Note: if the secondary is in an HA cluster, you will need to ensure that each
standby is sealed/unsealed with the new (primary’s) unseal keys. If one of the
standbys takes over on active duty before this happens it will seal itself to
remove it from rotation (e.g. if using Consul for service discovery), but if a
standby does not attempt taking over it will throw errors. We hope to make this
workflow better in a future update.

### Dev-Mode Root Tokens

To ease development and testing, when both the primary and secondary are
running in development mode, the initial root token created by the primary
(including those with custom IDs specified with `-dev-root-token-id`) will be
populated into the secondary upon activation. This allows a developer to keep a
consistent `~/.vault-token` file or `VAULT_TOKEN` environment variable when
working with both clusters.

On a production system, after a secondary is activated, the enabled
authentication backends should be used to get tokens with appropriate policies,
as policies and auth backend configuration are replicated.

The generate-root command can also be used to generate a root token local to
the secondary cluster.

## Managing Vault Performance Replication

Vault’s performance replication model is intended to allow horizontally scaling Vault’s
functions rather than to act in a strict Disaster Recovery (DR) capacity. For more information on Vault's disaster recovery replication, look at the
[general information page](/docs/vault-enterprise/replication/index.html).

As a result, Vault performance replication acts on static items within Vault, meaning
information that is not part of Vault’s lease-tracking system. In a practical
sense, this means that all Vault information is replicated from the primary to
secondaries except for tokens and secret leases.

Because token information must be checked and possibly rewritten with each use
(e.g. to decrement its use count), replicated tokens would require every call
to be forwarded to the primary, decreasing rather than increasing total Vault
throughput.

Secret leases are tracked independently for two reasons: one, because every
such lease is tied to a token and tokens are local to each cluster; and two,
because tracking large numbers of leases is memory-intensive and tracking all
leases in a replicated fashion could dramatically increase the memory
requirements across all Vault nodes.

We believe that this performance replication model provides significant utility for horizontally scaling Vault’s functionality.  However, it does mean
that certain principles must be kept in mind.

### Always Use the Local Cluster

First and foremost, when designing systems to take advantage of replicated
Vault, you must ensure that they always use the same Vault cluster for all
operations, as only that cluster will know about the client’s Vault token.

### Enabling a Secondary Wipes Storage

Replication relies on having a shared keyring between primary and secondaries
and also relies on having a shared understanding of the data store state. As a
result, when replication is enabled, all of the secondary’s existing storage
will be wiped. This is irrevocable. Make a backup first if there is a remote
chance you’ll need some of this data at some future point.

Generally, activating as a secondary will be the first thing that is done upon
setting up a new cluster for replication.

### Replicated vs. Local Backend Mounts

All backend mounts (of all types) that can be enabled within Vault default to
being mounted as a replicated mount. This means that mounts cannot be enabled
on a secondary, and mounts enabled on the primary will replicate to
secondaries.

Mounts can also be marked local (via the `-local` flag on the Vault CLI or
setting the `local` parameter to `true` in the API). This can only be performed
at mount time; if a mount is local but should have been replicated, or vice
versa, you must unmount the backend and mount a new instance at that path with
the local flag enabled.

Local mounts do not propagate data from the primary to secondaries, and local
mounts on secondaries do not have their data removed during the syncing
process. The exception is during initial bootstrapping of a secondary from a
state where replication is disabled; all data, including local mounts, is
deleted at this time (as the encryption keys will have changed so data in local
mounts would be unable to be read).

### Audit Backends

In normal Vault usage, if Vault has at least one audit backend configured and
is unable to successfully log to at least one backend, it will block further
requests.

Replicated audit mounts must be able to successfully log on all replicated
clusters. For example, if using the file backend, the configured path must be
able to be written to by all secondaries. It may be useful to use at least one
local audit mount on each cluster to prevent such a scenario.

### Never Have Two Primaries

The replication model is not designed for active-active usage and enabling two
primaries should never be done, as it can lead to data loss if they or their
secondaries are ever reconnected.

### Disaster Recovery

Local backend mounts are not replicated and their use will require existing DR
mechanisms if DR is necessary in your implementation.

If you need true DR, look at the
[general information page](/docs/vault-enterprise/replication/index.html) for information on Vault's disaster recovery replication.



