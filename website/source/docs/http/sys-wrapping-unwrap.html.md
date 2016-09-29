---
layout: "http"
page_title: "HTTP API: /sys/wrapping/unwrap"
sidebar_current: "docs-http-wrapping-unwrap"
description: |-
  The '/sys/wrapping/unwrap' endpoint unwraps a wrapped response
---

# /sys/wrapping/unwrap

## POST

<dl>
  <dt>Description</dt>
  <dd>
    Returns the original response inside the given wrapping token. Unlike
    simply reading `cubbyhole/response` (which is deprecated), this endpoint
    provides additional validation checks on the token, returns the original
    value on the wire rather than a JSON string representation of it, and
    ensures that the response is properly audit-logged.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/sys/wrapping/unwrap`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">token</span>
        <span class="param-flags">required</span>
        The wrapping token ID.
      </li>
    </ul>
  </dd>

  <dt>Returns</dt>
  <dd>

    ```javascript
    {
        "request_id": "8e33c808-f86c-cff8-f30a-fbb3ac22c4a8",
        "lease_id": "",
        "lease_duration": 2592000,
        "renewable": false,
        "data": {
                "zip": "zap"
        },
        "warnings": null
    }
    ```

  </dd>
</dl>
