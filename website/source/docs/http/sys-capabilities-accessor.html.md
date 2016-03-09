---
layout: "http"
page_title: "HTTP API: /sys/capabilities-accessor"
sidebar_current: "docs-http-auth-capabilities-accessor"
description: |-
  The `/sys/capabilities-accessor` endpoint is used to fetch the capabilities of the token associated with an accessor, on the given path.
---

# /sys/capabilities-accessor

## POST

<dl>
  <dt>Description</dt>
  <dd>
    Returns the capabilities of the token associated with an accessor, on the given path.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">accessor</span>
        <span class="param-flags">required</span>
        Accessor of the token.
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
