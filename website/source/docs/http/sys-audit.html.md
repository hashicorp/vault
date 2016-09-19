---
layout: "http"
page_title: "HTTP API: /sys/audit"
sidebar_current: "docs-http-audits-audits"
description: |-
  The `/sys/audit` endpoint is used to enable and disable audit backends.
---

# /sys/audit

## GET

<dl>
  <dt>Description</dt>
  <dd>
    List the mounted audit backends. _This endpoint requires `sudo`
    capability._
  </dd>

  <dt>Method</dt>
  <dd>GET</dd>

  <dt>Parameters</dt>
  <dd>
    None
  </dd>

  <dt>Returns</dt>
  <dd>

    ```javascript
    {
      "file": {
        "type: "file",
        "description: "Store logs in a file",
        "options": {
          "path": "/var/log/file"
        }
      }
    }
    ```

  </dd>
</dl>

## PUT

<dl>
  <dt>Description</dt>
  <dd>
    Enable an audit backend. _This endpoint requires `sudo` capability._
  </dd>

  <dt>Method</dt>
  <dd>PUT</dd>

  <dt>URL</dt>
  <dd>`/sys/audit/<path>`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">type</span>
        <span class="param-flags">required</span>
        The type of the audit backend.
      </li>
      <li>
        <span class="param">description</span>
        <span class="param-flags">optional</span>
        A description of the audit backend for operators.
      </li>
      <li>
        <span class="param">options</span>
        <span class="param-flags">optional</span>
        An object of options to configure the backend. This is
        dependent on the backend type. Please consult the documentation
        for the backend type you intend to use.
      </li>
    </ul>
  </dd>

  <dt>Returns</dt>
  <dd>`204` response code.
  </dd>
</dl>

## DELETE

<dl>
  <dt>Description</dt>
  <dd>
    Disable the given audit backend. _This endpoint requires `sudo`
    capability._
  </dd>

  <dt>Method</dt>
  <dd>DELETE</dd>

  <dt>URL</dt>
  <dd>`/sys/audit/<path>`</dd>

  <dt>Parameters</dt>
  <dd>None
  </dd>

  <dt>Returns</dt>
  <dd>`204` response code.
  </dd>
</dl>
