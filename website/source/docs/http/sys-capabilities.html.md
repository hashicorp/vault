---
layout: "http"
page_title: "HTTP API: /sys/capabilities"
sidebar_current: "docs-http-auth-capabilities"
description: |-
  The `/sys/capabilities` endpoint is used to fetch the capabilities of a token on a given path.
---

# /sys/capabilities

## POST

<dl>
  <dt>Description</dt>
  <dd>
    Returns the capabilities of the token on the given path.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">token</span>
        <span class="param-flags">required</span>
        Token for which capabilities are being queried.
      </li>
      <li>
        <span class="param">path</span>
        <span class="param-flags">required</span>
        Path on which the token's capabilities will be checked.
      </li>
    </ul>
  </dd>

  <dt>Returns</dt>
  <dd>

    ```javascript
    {
        "capabilities": ["read", "list"]
    }
    ```

  </dd>
</dl>
