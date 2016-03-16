---
layout: "http"
page_title: "HTTP API: /sys/audit-hash"
sidebar_current: "docs-http-audits-hash"
description: |-
  The `/sys/audit-hash` endpoint is used to hash data using an audit backend's hash function and salt.
---

# /sys/audit-hash

## POST

<dl>
  <dt>Description</dt>
  <dd>
    Hash the given input data with the specified audit backend's hash function
    and salt. This endpoint can be used to discover whether a given plaintext
    string (the `input` parameter) appears in the audit log in obfuscated form.
    Note that the audit log records requests and responses; since the Vault API
    is JSON-based, any binary data returned from an API call (such as a
    DER-format certificate) is base64-encoded by the Vault server in the
    response, and as a result such information should also be base64-encoded to
    supply into the `input` parameter.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/sys/audit-hash/<path>`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">input</span>
        <span class="param-flags">required</span>
        The input string to hash.
      </li>
    </ul>
  </dd>

  <dt>Returns</dt>
  <dd>

    ```javascript
    {
      "hash": "hmac-sha256:08ba357e274f528065766c770a639abf6809b39ccfd37c2a3157c7f51954da0a"
    }
    ```

  </dd>
</dl>
