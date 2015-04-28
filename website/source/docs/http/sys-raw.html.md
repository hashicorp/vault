---
layout: "http"
page_title: "HTTP API: /sys/raw"
sidebar_current: "docs-http-debug-raw"
description: |-
  The `/sys/raw` endpoint is access the raw underlying store in Vault.
---

# /sys/raw

## GET

<dl>
  <dt>Description</dt>
  <dd>
      Reads the value of the key at the given path. This is the raw path in the
        storage backend and not the logical path that is exposed via the mount system.
  </dd>

  <dt>Method</dt>
  <dd>GET</dd>

  <dt>URL</dt>
  <dd>`/sys/raw/<path>`</dd>

  <dt>Parameters</dt>
  <dd>
    None
  </dd>

  <dt>Returns</dt>
  <dd>

    ```javascript
    {
      "value": "{'foo':'bar'}"
    }
    ```

  </dd>
</dl>

## PUT

<dl>
  <dt>Description</dt>
  <dd>
      Update the value of the key at the given path. This is the raw path in the
        storage backend and not the logical path that is exposed via the mount system.
  </dd>

  <dt>Method</dt>
  <dd>PUT</dd>

  <dt>URL</dt>
  <dd>`/sys/raw/<path>`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">value</span>
        <span class="param-flags">required</span>
        The value of the key.
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
    Delete the key with given path. This is the raw path in the
        storage backend and not the logical path that is exposed via the mount system.
  </dd>

  <dt>Method</dt>
  <dd>DELETE</dd>

  <dt>URL</dt>
  <dd>`/sys/raw/<path>`</dd>

  <dt>Parameters</dt>
  <dd>None
  </dd>

  <dt>Returns</dt>
  <dd>`204` response code.
  </dd>
</dl>
