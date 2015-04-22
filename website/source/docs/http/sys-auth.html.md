---
layout: "http"
page_title: "HTTP API: /sys/auth"
sidebar_current: "docs-http-auth-auth"
description: |-
  The `/sys/auth` endpoint is used to manage auth backends in Vault.
---

# /sys/auth

## GET

<dl>
  <dt>Description</dt>
  <dd>
    Lists all the enabled auth backends.
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
      "github": {
        "type": "github",
        "description": "GitHub auth"
      }
    }
    ```

  </dd>
</dl>

## POST

<dl>
  <dt>Description</dt>
  <dd>
    Enable a new auth backend. The auth backend can be accessed
    and configured via the mount point specified in the URL. This
    mount point will be exposed under the `auth` prefix. For example,
    enabling with the `/sys/auth/foo` URL will make the backend
    available at `/auth/foo`.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/sys/auth/<mount point>`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">type</span>
        <span class="param-flags">required</span>
        The name of the auth backend type, such as "github"
      </li>
      <li>
        <span class="param">description</span>
        <span class="param-flags">optional</span>
        A human-friendly description of the auth backend.
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
    Disable the auth backend at the given mount point.
  </dd>

  <dt>Method</dt>
  <dd>DELETE</dd>

  <dt>URL</dt>
  <dd>`/sys/auth/<mount point>`</dd>

  <dt>Parameters</dt>
  <dd>None
  </dd>

  <dt>Returns</dt>
  <dd>`204` response code.
  </dd>
</dl>
