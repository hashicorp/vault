---
layout: "http"
page_title: "HTTP API: /sys/wrapping/rewrap"
sidebar_current: "docs-http-wrapping-rewrap"
description: |-
  The '/sys/wrapping/rewrap' endpoint can be used to rotate a wrapping token and refresh its TTL
---

# /sys/wrapping/rewrap

## POST

<dl>
  <dt>Description</dt>
  <dd>
    Rewraps a response-wrapped token; the new token will use the same creation
    TTL as the original token and contain the same response. The old token will
    be invalidated. This can be used for long-term storage of a secret in a
    response-wrapped token when rotation is a requirement.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/sys/wrapping/rewrap`</dd>

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
        "request_id": "",
        "lease_id": "",
        "lease_duration": 0,
        "renewable": false,
        "data": null,
        "warnings": null,
        "wrap_info": {
                "token": "3b6f1193-0707-ac17-284d-e41032e74d1f",
                "ttl": 300,
                "creation_time": "2016-09-28T14:22:26.486186607-04:00",
                "wrapped_accessor": ""
        }
    }
    ```

  </dd>
</dl>
