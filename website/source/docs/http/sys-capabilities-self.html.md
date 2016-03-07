---
layout: "http"
page_title: "HTTP API: /sys/capabilities-self"
sidebar_current: "docs-http-auth-capabilities-self"
description: |-
  The `/sys/capabilities-self` endpoint is used to fetch the capabilities of client token on a given path.
---

# /sys/capabilities-self

## POST

<dl>
  <dt>Description</dt>
  <dd>
    Returns the capabilities of client token on the given path.
    Client token is the Vault token with which this API call is made.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">path</span>
        <span class="param-flags">required</span>
        Path on which the client token's capabilities will be checked.
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
