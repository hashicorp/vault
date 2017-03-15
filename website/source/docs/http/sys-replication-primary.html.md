---
layout: "http"
page_title: "HTTP API: /sys/repliation/primary"
sidebar_current: "docs-http-replication-primary"
description: |-
  The '/sys/replication/primary' endpoint focuses on managing replication operations for primary clusters. 
---

# /sys/replication/primary/enable

## POST

<dl>
  <dt>Description</dt>
  <dd>
    Enables replication in primary mode. This is used when replication is currently disabled on the cluster
	(if the cluster is already a secondary, it must be promoted).
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/sys/repliation/primary/enable`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">primary_cluster_addr /span>
        <span class="param-flags">optional</span>
       Can be used to override the cluster address that
	   the primary gives to secondary nodes. Useful if the 
	   primary’s cluster address is not directly accessible
	   and must be accessed via an alternate path/address.
      </li>
    </ul>
  </dd>

  <dt>Returns</dt>
 <dd>`204` response code.
 </dd>
</dl>

# /sys/replication/primary/demote

## POST

<dl>
  <dt>Description</dt>
  <dd>
    Demotes a primary cluster to a secondary. This secondary cluster 
	will not attempt to connect to a primary (see the update-primary call),
	but will maintain knowledge of its cluster ID and can be reconnected
	to the same replication set without wiping local storage.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/sys/repliation/primary/demote`</dd>

  <dt>Parameters</dt>
  <dd>
      None
  </dd>

  <dt>Returns</dt>
 <dd>`204` response code.
 </dd>
</dl>


# /sys/replication/primary/disable

## POST

<dl>
  <dt>Description</dt>
  <dd>
   Disable replication entirely on the cluster. Any secondaries will no longer be able to connect.
   Caution: re-enabling this node as a primary or secondary will change its cluster ID; in the secondary
   case this means a wipe of the underlying storage when connected to a primary, and in the primary case,
   secondaries connecting back to the cluster (even if they have connected before) will require a wipe of
   the underlying storage.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/sys/repliation/primary/demote`</dd>

  <dt>Parameters</dt>
  <dd>
      None
  </dd>

  <dt>Returns</dt>
 <dd>`204` response code.
 </dd>
</dl>