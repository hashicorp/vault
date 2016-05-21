---
layout: "http"
page_title: "HTTP API: /sys/step-down"
sidebar_current: "docs-http-ha-step-down"
description: |-
  The '/sys/step-down' endpoint causes the node to give up active status.
---

# /sys/step-down

<dl>
  <dt>Description</dt>
  <dd>
    Forces the node to give up active status. If the node does not have active
    status, this endpoint does nothing. Note that the node will sleep for ten
    seconds before attempting to grab the active lock again, but if no standby
    nodes grab the active lock in the interim, the same node may become the
    active node again. Requires a token with `root` policy or `sudo` capability
    on the path.
  </dd>

  <dt>Method</dt>
  <dd>PUT</dd>

  <dt>Parameters</dt>
  <dd>
    None
  </dd>

  <dt>Returns</dt>
  <dd>A `204` response code.
  </dd>
</dl>
