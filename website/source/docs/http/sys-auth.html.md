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
    Enable a new auth backend. The auth backend can be accessed and configured
    via the auth path specified in the URL. This auth path will be exposed
    under the `auth` prefix. For example, enabling with the `/sys/auth/foo` URL
    will make the backend available at `/auth/foo`. _This endpoint requires
    `sudo` capability on the final path._
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/sys/auth/<auth_path>`</dd>

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
    Disable the auth backend at the given auth path. _This endpoint requires
    `sudo` capability on the final path._
  </dd>

  <dt>Method</dt>
  <dd>DELETE</dd>

  <dt>URL</dt>
  <dd>`/sys/auth/<auth_path>`</dd>

  <dt>Parameters</dt>
  <dd>None
  </dd>

  <dt>Returns</dt>
  <dd>`204` response code.
  </dd>
</dl>

# /sys/auth/[auth-path]/tune

## GET

<dl>
  <dt>Description</dt>
  <dd>
    Read the given auth path's configuration. Returns the current time
    in seconds for each TTL, which may be the system default or a auth path
    specific value. _This endpoint requires `sudo` capability on the final
    path, but the same functionality can be achieved without `sudo` via
    `sys/mounts/auth/[auth-path]/tune`._
  </dd>

  <dt>Method</dt>
  <dd>GET</dd>

  <dt>URL</dt>
  <dd>`/sys/auth/[auth-path]/tune`</dd>

  <dt>Parameters</dt>
  <dd>
    None
  </dd>

  <dt>Returns</dt>
  <dd>

    ```javascript
    {
      "default_lease_ttl": 3600,
      "max_lease_ttl": 7200
    }
    ```

  </dd>
</dl>

## POST

<dl>
  <dt>Description</dt>
  <dd>
    Tune configuration parameters for a given auth path. _This endpoint
    requires `sudo` capability on the final path, but the same functionality
    can be achieved without `sudo` via `sys/mounts/auth/[auth-path]/tune`._
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/sys/auth/[auth-path]/tune`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">default_lease_ttl</span>
        <span class="param-flags">optional</span>
        The default time-to-live. If set on a specific auth path,
        overrides the global default. A value of "system" or "0"
        are equivalent and set to the system default TTL.
      </li>
      <li>
        <span class="param">max_lease_ttl</span>
        <span class="param-flags">optional</span>
        The maximum time-to-live. If set on a specific auth path,
        overrides the global default. A value of "system" or "0"
        are equivalent and set to the system max TTL.
      </li>
    </ul>
  </dd>

  <dt>Returns</dt>
  <dd>`204` response code.
  </dd>
</dl>
