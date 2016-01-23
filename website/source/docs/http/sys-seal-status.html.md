---
layout: "http"
page_title: "HTTP API: /sys/seal-status"
sidebar_current: "docs-http-seal-status"
description: |-
  The '/sys/seal-status' endpoint is used to check the seal status of a Vault.
---

# /sys/seal-status

<dl>
  <dt>Description</dt>
  <dd>
    Returns the seal status of the Vault.<br/><br/>This is an unauthenticated endpoint.
  </dd>

  <dt>Method</dt>
  <dd>GET</dd>

  <dt>Parameters</dt>
  <dd>
    None
  </dd>

  <dt>Returns</dt>
  <dd>
    The "t" parameter is the threshold, and "n" is the number of shares.

    ```javascript
    {
      "sealed": true,
      "t": 3,
      "n": 5,
      "progress": 2
    }
    ```

  </dd>
</dl>
