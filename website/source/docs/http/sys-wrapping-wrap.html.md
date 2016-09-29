---
layout: "http"
page_title: "HTTP API: /sys/wrapping/wrap"
sidebar_current: "docs-http-wrapping-wrap"
description: |-
  The '/sys/wrapping/wrap' endpoint wraps the given values in a response-wrapped token
---

# /sys/wrapping/wrap

## POST

<dl>
  <dt>Description</dt>
  <dd>
    Wraps the given user-supplied data inside a response-wrapped token.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/sys/wrapping/wrap`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">[any]</span>
        <span class="param-flags">optional</span>
        Parameters should be supplied as keys/values in a JSON object. The
        exact set of given parameters will be contained in the wrapped
        response.
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
                "token": "fb79b9d3-d94e-9eb6-4919-c559311133d6",
                "ttl": 300,
                "creation_time": "2016-09-28T14:41:00.56961496-04:00",
                "wrapped_accessor": ""
        }
    }
    ```

  </dd>
</dl>
