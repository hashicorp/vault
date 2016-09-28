---
layout: "http"
page_title: "HTTP API: /sys/wrapping/lookup"
sidebar_current: "docs-http-wrapping-lookup"
description: |-
  The '/sys/wrapping/lookup' endpoint returns wrapping token properties
---

# /sys/wrapping/lookup

## POST

<dl>
  <dt>Description</dt>
  <dd>
    Looks up wrapping properties for the given token.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/sys/wrapping/lookup`</dd>

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
        "request_id": "481320f5-fdf8-885d-8050-65fa767fd19b",
        "lease_id": "",
        "lease_duration": 0,
        "renewable": false,
        "data": {
                "creation_time": "2016-09-28T14:16:13.07103516-04:00",
                "creation_ttl": 300
        },
        "warnings": null
    }
    ```

  </dd>
</dl>
