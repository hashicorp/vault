---
layout: "http"
page_title: "HTTP API: /sys/mounts"
sidebar_current: "docs-http-mounts-mounts"
description: |-
  The '/sys/mounts' endpoint is used manage secret backends in Vault.
---

# /sys/mounts

## GET

<dl>
  <dt>Description</dt>
  <dd>
    Lists all the mounted secret backends.
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
      "aws": {
        "type": "aws",
        "description": "AWS keys"
      },

      "sys": {
        "type": "system",
        "description": "system endpoint"
      }
    }
    ```

  </dd>
</dl>

## POST

<dl>
  <dt>Description</dt>
  <dd>
    Mount a new secret backend to the mount point in the URL.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/sys/mounts/<mount point>`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">type</span>
        <span class="param-flags">required</span>
        The name of the backend type, such as "aws"
      </li>
      <li>
        <span class="param">description</span>
        <span class="param-flags">optional</span>
        A human-friendly description of the mount.
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
    Unmount the mount point specified in the URL.
  </dd>

  <dt>Method</dt>
  <dd>DELETE</dd>

  <dt>URL</dt>
  <dd>`/sys/mounts/<mount point>`</dd>

  <dt>Parameters</dt>
  <dd>None
  </dd>

  <dt>Returns</dt>
  <dd>`204` response code.
  </dd>
</dl>
