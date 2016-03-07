---
layout: "http"
page_title: "HTTP API: /sys/seal"
sidebar_current: "docs-http-seal-seal"
description: |-
  The '/sys/seal' endpoint seals the Vault.
---

# /sys/seal

<dl>
  <dt>Description</dt>
  <dd>
    Seals the Vault. In HA mode, only an active node can be sealed. Standby
    nodes should be restarted to get the same effect. Requires a token with
    `root` policy or `sudo` capability on the path.
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
