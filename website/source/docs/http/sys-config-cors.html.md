---
layout: "http"
page_title: "HTTP API: /sys/config/cors"
sidebar_current: "docs-http-config-cors"
description: |-
  The '/sys/config/cors' endpoint configures how the Vault server responds to cross-origin requests.
---

# /sys/config/cors

This is a protected path, therefore all requests require a token with `root`
policy or `sudo` capability on the path.

## GET

<dl>
  <dt>Description</dt>
  <dd>
    Returns the current CORS configuration.
  </dd>

  <dt>Method</dt>
  <dd>GET</dd>

  <dt>URL</dt>
  <dd>`/sys/config/cors`</dd>

  <dt>Parameters</dt>
  <dd>
    None
  </dd>

  <dt>Returns</dt>
  <dd>

    ```javascript
    {
      "enabled": true,
      "allowed_origins": "http://www.example.com"
    }
    ```

    Sample response when CORS is disabled.

    ```javascript
    {
      "enabled": false,
      "allowed_origins": ""
    }
    ```
  </dd>
</dl>

## PUT

<dl>
  <dt>Description</dt>
  <dd>
    Configures the Vault server to return CORS headers for origins that are
    permitted to make cross-origin requests based on the `allowed_origins`
    parameter.
  </dd>

  <dt>Method</dt>
  <dd>PUT</dd>

  <dt>URL</dt>
  <dd>`/sys/config/cors`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">allowed_origins</span>
        <span class="param-flags">required</span>
        Valid values are either a wildcard (*) or a comma-separated list of
        exact origins that are permitted to make cross-origin requests.
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
    Disables the CORS functionality of the Vault server.
  </dd>

  <dt>Method</dt>
  <dd>DELETE</dd>

  <dt>URL</dt>
  <dd>`/sys/config/cors`</dd>

  <dt>Parameters</dt>
  <dd>
    None
  </dd>

  <dt>Returns</dt>
  <dd>`204` response code.
  </dd>
</dl>
