---
layout: "http"
page_title: "HTTP API: /sys/repliation/primary"
sidebar_current: "docs-http-replication-primary"
description: |-
  The '/sys/replication/primary' endpoint focuses on managing replication behavior for a primary cluster, including management of secondaries.
---

# /sys/replication/primary/enable

## POST

<dl>
  <dt>Description</dt>
  <dd>
    Enables replication in primary mode. This is used when replication is
    currently disabled on the cluster (if the cluster is already a secondary,
    it must be promoted).
    
    Caution: only one primary should be active at a given time. Multiple
    primaries may result in data loss!

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
        Can be used to override the cluster address that the primary gives to
        secondary nodes. Useful if the primary’s cluster address is not
        directly accessible and must be accessed via an alternate path/address,
        such as through a TCP-based load balancer.
      </li>
    </ul>
  </dd>

  <dt>Returns</dt>
 <dd>
   `200` response code with a warning.
 </dd>
</dl>

# /sys/replication/primary/demote

## POST

<dl>
  <dt>Description</dt>
  <dd>
    Demotes a primary cluster to a secondary. This secondary cluster will not
    attempt to connect to a primary (see the update-primary call), but will
    maintain knowledge of its cluster ID and can be reconnected to the same
    replication set without wiping local storage.
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
 <dd>
   `200` response code with a warning.
 </dd>
</dl>


# /sys/replication/primary/disable

## POST

<dl>
  <dt>Description</dt>
  <dd>
   Disable replication entirely on the cluster. Any secondaries will no longer
   be able to connect. Caution: re-enabling this node as a primary or secondary
   will change its cluster ID; in the secondary case this means a wipe of the
   underlying storage when connected to a primary, and in the primary case,
   secondaries connecting back to the cluster (even if they have connected
   before) will require a wipe of the underlying storage.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/sys/repliation/primary/disable`</dd>

  <dt>Parameters</dt>
  <dd>
      None
  </dd>

  <dt>Returns</dt>
 <dd>
   `200` response code with a warning..
 </dd>
</dl>

# /sys/replication/primary/secondary-token

## GET

<dl>
  <dt>Description</dt>
  <dd>
    Requires ‘sudo’ capability. Generate a secondary activation token for the
    cluster with the given opaque identifier, which must be unique. This
    identifier can later be used to revoke a secondary's access.
  </dd>

  <dt>Method</dt>
  <dd>GET</dd>

  <dt>URL</dt>
  <dd>`/sys/replication/primary/secondary-token`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">id</span>
        <span class="param-flags">required</span>
        An opaque identifier, e.g. ‘us-east’
      </li>
      <li>
        <span class="param">ttl</span>
        <span class="param-flags">optional</span>
        The TTL for the secondary activation token. Defaults to ‘"30m"’.
      </li>
    </ul>
  </dd>

  <dt>Returns</dt>
  <dd>

    ```javascript
    {
      "request_id": "",
      "lease_id": "",
      "lease_duration": 0,
      "renewable": false,
      "data": null,
      "warnings": null,
      "wrap_info": {
        "token": "fb79b9d3-d94e-9eb6-4919-c559311133d6",
        "ttl": 300,
        "creation_time": "2016-09-28T14:41:00.56961496-04:00",
        "wrapped_accessor": ""
      }
    }
    ```

  </dd>
</dl>

# /sys/replication/primary/revoke-secondary

## POST

<dl>
  <dt>Description</dt>
  <dd>
    Revoke a secondary’s ability to connect to the primary cluster; the
    secondary will immediately be disconnected and will not be allowed to
    connect again unless given a new activation token.
  </dd>

  <dt>Method</dt>
  <dd></dd>

  <dt>URL</dt>
  <dd>`/sys/replication/secondary/revoke-secondary`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">id</span>
        <span class="param-flags">required</span>
        The identifier used when fetching the secondary token.
      </li>
    </ul>
  </dd>

  <dt>Returns</dt>
  <dd>
    `200` response code with a warning.
  </dd>
</dl>


