---
layout: "http"
page_title: "HTTP API: /sys/leader"
sidebar_current: "docs-http-ha-leader"
description: |-
  The '/sys/leader' endpoint is used to check the high availability status and current leader of Vault.
---

# /sys/leader

<dl>
  <dt>Description</dt>
  <dd>
    Returns the high availability status and current leader instance of Vault.
  </dd>

  <dt>Method</dt>
  <dd>GET</dd>

  <dt>Parameters</dt>
  <dd>
    None
  </dd>

  <dt>Returns</dt>
  <dd>

    ```javascript
    {
      "ha_enabled": true,
      "is_self": false,
      "leader_address": "https://127.0.0.1:8200/"
    }
    ```

  </dd>
</dl>
