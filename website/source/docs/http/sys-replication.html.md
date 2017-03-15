---
layout: "http"
page_title: "HTTP API: /sys/replication"
sidebar_current: "docs-http-replication"
description: |-
  The '/sys/replication' endpoint focuses on managing general operations in Vault Enterprise replication sets
---

# /sys/replication/recover

## POST

<dl>
  <dt>Description</dt>
  <dd>
    Attempts recovery if replication is in an adverse state. For example: an
    error has caused replication to stop syncing.
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
    A `200` response code and a warning.
  </dd>
</dl>


# /sys/replication/reindex

## POST

<dl>
  <dt>Description</dt>
  <dd>
    Requires ‘sudo’ capability. Reindex the local data storage. This can cause
    a very long delay depending on the number and size of objects in the data
    store.
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
    A `200` response code and a warning.
  </dd>
</dl>

# /sys/replication/status

## GET

<dl>
  <dt>Description</dt>
  <dd>
    Unauthenticated. Print information about the status of replication (mode,
    sync progress, etc).
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
    The printed status of the replication environment. As an example, for a
    primary, it will look something like:

    ```javascript
    {
      "mode": "primary",
      "cluster_id": "d4095d41-3aee-8791-c421-9bc7f88f7c3e",
      "known_secondaries": [],
      "last_wal": 0,
      "merkle_root": "c3260c4c682ff2d6eb3c8bfd877134b3cec022d1",
      "request_id": "009ea98c-06cd-6dc3-74f2-c4904b22e535",
      "lease_id": "",
      "renewable": false,
      "lease_duration": 0,
      "data": {
        "cluster_id": "d4095d41-3aee-8791-c421-9bc7f88f7c3e",
        "known_secondaries": [],
        "last_wal": 0,
        "merkle_root": "c3260c4c682ff2d6eb3c8bfd877134b3cec022d1",
        "mode": "primary"
      },
      "wrap_info": null,
      "warnings": null,
      "auth": null
    }
    ```
  </dd>
</dl>
