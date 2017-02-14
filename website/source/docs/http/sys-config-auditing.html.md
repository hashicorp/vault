---
layout: "http"
page_title: "HTTP API: /sys/config/auditing"
sidebar_current: "docs-http-config-auditing"
description: |-
  The `/sys/config/auditing` endpoint is used to configure auditing settings.
---

# /sys/config/auditing/request-headers

## GET

<dl>
  <dt>Description</dt>
  <dd>
    List the request headers that are configured to be audited. _This endpoint requires `sudo`
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
        "headers":{
            "X-Forwarded-For": {
                "hmac":true
            }
        }
    }
    ```

  </dd>
</dl>

# /sys/config/auditing/request-headers/

## GET

<dl>
  <dt>Description</dt>
  <dd>
    List the information for the given request header. _This endpoint requires `sudo`
    capability._
  </dd>

  <dt>Method</dt>
  <dd>GET</dd>

  <dt>URL</dt>
  <dd>`/sys/config/auditing/request-headers/<name>`</dd>

  <dt>Parameters</dt>
  <dd>
    None
  </dd>

  <dt>Returns</dt>
  <dd>

    ```javascript
    {
        "X-Forwarded-For":{
            "hmac":true
        }
    }
    ```

  </dd>
</dl>

## PUT

<dl>
  <dt>Description</dt>
  <dd>
    Enable auditing of a header. _This endpoint requires `sudo` capability._
  </dd>

  <dt>Method</dt>
  <dd>PUT</dd>
 
  <dt>URL</dt>
  <dd>`/sys/config/auditing/request-headers/<name>`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">hmac</span>
        <span class="param-flags">optional</span>
        Bool, if this header's value should be hmac'ed in the audit logs.
        Defaults to false.
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
    Disable auditing of the given request header. _This endpoint requires `sudo`
    capability._
  </dd>

  <dt>Method</dt>
  <dd>DELETE</dd>

  <dt>URL</dt>
  <dd>`/sys/config/auditing/request-headers/<name>`</dd>

  <dt>Parameters</dt>
  <dd>None
  </dd>

  <dt>Returns</dt>
  <dd>`204` response code.
  </dd>
</dl>
