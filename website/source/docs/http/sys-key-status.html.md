---
layout: "http"
page_title: "HTTP API: /sys/key-status"
sidebar_current: "docs-http-rotate-key-status"
description: |-
  The '/sys/key-status' endpoint is used to query info about the current encryption key of Vault.
---

# /sys/key-status

<dl>
  <dt>Description</dt>
  <dd>
    Returns information about the current encryption key used by Vault.
  </dd>

  <dt>Method</dt>
  <dd>GET</dd>

  <dt>Parameters</dt>
  <dd>
    None
  </dd>

  <dt>Returns</dt>
  <dd>
    The "term" parameter is the sequential key number, and "install_time" is the time that
    encryption key was installed.

    ```javascript
    {
      "term": 3,
      "install_time": "2015-05-29T14:50:46.223692553-07:00"
    }
    ```

  </dd>
</dl>
