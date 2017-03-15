---
layout: "http"
page_title: "HTTP API: /sys/replication/secondary"
sidebar_current: "docs-http-replication-secondary"
description: |-
  The '/sys/replication/secondary' endpoint focuses on replication management operations on secondary clusters.
---

# /sys/replication/secondary/enable

## POST

<dl>
  <dt>Description</dt>
  <dd>
    Enables replication on a secondary using a secondary activation token.

    Caution: this will immediately clear all data in the cluster!
  </dd>

  <dt>Method</dt>
  <dd></dd>

  <dt>URL</dt>
  <dd>`/sys/replication/secondary/enable`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">token</span>
        <span class="param-flags">required</span>
        The secondary activation token fetched from the primary.
      </li>
      <li>
        <span class="param">primary_api_addr</span>
        <span class="param-flags">optional</span>
        Set this to the API address (normal Vault address) to override the
        value embedded in the token. This can be useful if the primary’s
        redirect address is not accessible directly from this cluster (e.g.
        through a load balancer).
      </li>
      <li>
        <span class="param">ca_file</span>
        <span class="param-flags">optional</span>
        The path to a CA root file (PEM format) that the secondary can use when
        unwrapping the token from the primary. If this and ca_path are not
        given, defaults to system CA roots.
      </li>
      <li>
        <span class="param">ca_path</span>
        <span class="param-flags">optional</span>
        The path to a CA root directory containing PEM-format files that the
        secondary can use when unwrapping the token from the primary. If this
        and ca_file are not given, defaults to system CA roots.
      </li>
    </ul>
  </dd>

  <dt>Returns</dt>
 <dd>
   `200` response code and a warning.
 </dd>
</dl>


# /sys/replication/secondary/promote

## POST

<dl>
  <dt>Description</dt>
  <dd>
    Promotes the secondary cluster to primary. For data safety and security
    reasons, new secondary tokens will need to be issued to other secondaries,
    and there should never be more than one primary at a time.
  </dd>

  <dt>Method</dt>
  <dd></dd>

  <dt>URL</dt>
  <dd>`/sys/replication/secondary/promote`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">primary_cluster_addr</span>
        <span class="param-flags">optional</span>
        Can be used to override the cluster address that the primary gives to
        secondary nodes. Useful if the primary’s cluster address is not
        directly accessible and must be accessed via an alternate path/address
        (e.g. through a load balancer).
      </li>
    </ul>
  </dd>

  <dt>Returns</dt>
 <dd>
   `200` response code and a warning.
 </dd>
</dl>

# /sys/replication/secondary/disable

## POST

<dl>
  <dt>Description</dt>
  <dd>
    Disable replication entirely on the cluster. The cluster will no longer be
    able to connect to the primary.

    Caution: re-enabling this node as a primary or secondary will change its
    cluster ID; in the secondary case this means a wipe of the underlying
    storage when connected to a primary, and in the primary case, secondaries
    connecting back to the cluster (even if they have connected before) will
    require a wipe of the underlying storage.
  </dd>

  <dt>Method</dt>
  <dd></dd>

  <dt>URL</dt>
  <dd>`/sys/replication/secondary/disable`</dd>

  <dt>Parameters</dt>
  <dd>
    None
  </dd>

  <dt>Returns</dt>
 <dd>
   `200` response code and a warning.
 </dd>
</dl>

# /sys/replication/secondary/update-primary

## POST

<dl>
  <dt>Description</dt>
  <dd>
    Change a secondary cluster’s assigned primary 
    cluster using a secondary activation token. 
    This does not wipe all data in the cluster.
  </dd>

  <dt>Method</dt>
  <dd></dd>

  <dt>URL</dt>
  <dd>`/sys/replication/secondary/update-primary`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">token</span>
        <span class="param-flags">required</span>
        The secondary activation token fetched from the primary. If you set
        this to a blank string, the cluster will stay a secondary but clear its
        knowledge of any past primary (and thus not attempt to connect to the
        previous primary). This can be useful if the primary is down to stop
        the secondary from trying to reconnect to it.
      </li>
      <li>
       <span class="param">primary_api_addr</span>
        <span class="param-flags">optional</span>
         Set this to the API address (normal Vault address) to override the
         value embedded in the token. This can be useful if the primary’s
         redirect address is not accessible directly from this cluster.
      </li>
      <li>
       <span class="param">ca_file</span>
        <span class="param-flags">optional</span>
        The path to a CA root file (PEM format) that the secondary can use when
        unwrapping the token from the primary. If this and ca_path are not
        given, defaults to system CA roots.
      </li>
      <li>
       <span class="param">ca_path</span>
        <span class="param-flags">optional</span>
        The path to a CA root directory containing PEM-format files that the
        secondary can use when unwrapping the token from the primary. If this
        and ca_file are not given, defaults to system CA roots.
      </li>
    </ul>
  </dd>

  <dt>Returns</dt>
 <dd>
   `200` response code and a warning.
 </dd>
</dl>
