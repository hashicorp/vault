---
layout: "guides"
page_title: "Upgrading to Vault 0.8.0 - Guides"
sidebar_current: "guides-upgrading-to-0.8.0"
description: |-
  This page contains the list of deprecations and important or breaking changes
  for Vault 0.8.0. Please read it carefully.
---

# Overview

This page contains the list of deprecations and important or breaking changes
for Vault 0.8.0 compared to the most recent release. Please read it carefully.

## Enterprise Upgrade Procedure

If you are upgrading from Vault Enterprise, you should take one of the
following upgrade paths. Please note that reindexing stops the ability of the
node to process requests while the indexing is happening. As a result, you may
want to plan the upgrade for a time when either no writes to the primary are
expected and secondaries can handle traffic (if using replication), or when
there is a maintenance window.

There are two reindex processes that need to happen during the upgrade, for two
different indexes. One will happen automatically when the cluster is upgraded
to 0.8.0. _This happens even if you are not currently using replication._ The
other can be done in one of a few ways, as follows.

### If Not Using Replication

If not using replication, no further action needs to be taken.

### If Using Replication

#### Option 1: Reindex the Primary, then Upgrade Secondaries

The first option is to issue a write to [`sys/replication/reindex`][reindex] on the
primary (it is not necessary on the secondaries). When the reindex on the
primary is finished, upgrade the secondaries, then upgrade the primary.

#### Option 2: Upgrade All Nodes Simultaneously

The second option is to upgrade all nodes to 0.8.0 at the same time. This
removes the need to perform an explicit reindex but may equate to more down
time since secondaries will not be able to service requests while the primary
is performing an explicit reindex.

## `sys/revoke-force` Requires `sudo` Capability

This path was meant to require `sudo` capability but was not implemented this
way. It now requires `sudo` capability to run.

[reindex]: https://www.vaultproject.io/api/system/replication.html#reindex-replication
