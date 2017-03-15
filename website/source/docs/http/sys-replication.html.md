---
layout: "http"
page_title: "HTTP API: /sys/replication"
sidebar_current: "docs-http-replication"
description: |-
  The '/sys/replication' endpoint focuses on managing general operations in Vault Enterprise Replication
---

# /sys/replication/recover

## POST

<dl>
  <dt>Description</dt>
  <dd>
    Attempts recovery if replication is in an adverse state. For example: an error has caused replication to stop syncing.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/sys/replication/recover`</dd>

  <dt>Parameters</dt>
  <dd>
    None
  </dd>

  <dt>Returns</dt>
  <dd>
    A '204' response code 
  </dd>
</dl>


# /sys/replication/reindex

## POST

<dl>
  <dt>Description</dt>
  <dd>
    Requires ‘sudo’ capability. Reindex the local data storage. This can cause a very long delay depending on the number and size of objects in the data store.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/sys/replication/reindex`</dd>

  <dt>Parameters</dt>
  <dd>
    None
  </dd>

  <dt>Returns</dt>
  <dd>
    A '204' response code 
  </dd>
</dl>

# /sys/replication/status

## GET

<dl>
  <dt>Description</dt>
  <dd>
    Unauthenticated. Print information about the status of replication (mode, sync progress, etc).
  </dd>

  <dt>Method</dt>
  <dd>GET</dd>

  <dt>URL</dt>
  <dd>`/sys/replication/status`</dd>

  <dt>Parameters</dt>
  <dd>
    None
  </dd>

  <dt>Returns</dt>
  <dd>
    The printed status of the replication environment. 
  </dd>
</dl>