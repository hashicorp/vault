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
    ensures that the response is properly audit-logged.<br/><br/>This endpoint
    can be used by using a wrapping token as the client token in the API call,
    in which case the `token` parameter is not required; or, a different token
    with permissions to access this endpoint can make the call and pass in the
    wrapping token in the `token` parameter. Do _not_ use the wrapping token in
    both locations; this will cause the wrapping token to be revoked but the
    value to be unable to be looked up, as it will basically be a double-use of
    the token!
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
        <span class="param-flags">optional</span>
        The wrapping token ID; required if the client token is not the wrapping
        token. Do not use the wrapping token in both locations.
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
