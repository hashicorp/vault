---
layout: "http"
page_title: "HTTP API: /sys/init"
sidebar_current: "docs-http-sys-init"
description: |-
  The '/sys/init' endpoint is used to initialize a new Vault.
---

# /sys/init

## GET

<dl>
  <dt>Description</dt>
  <dd>
    Return the initialization status of a Vault.
  </dd>

  <dt>Method</dt>
  <dd>GET</dd>

  <dt>Parameters</dt>
  <dd>None</dd>

  <dt>Returns</dt>
  <dd>

    ```javascript
    {
      "initialize": true
    }
    ```

  </dd>
</dl>

## PUT

<dl>
  <dt>Description</dt>
  <dd>
    Initializes a new Vault. The Vault must've not been previously
    initialized.
  </dd>

  <dt>Method</dt>
  <dd>PUT</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">secret_shares</span>
        <span class="param-flags">required</span>
        The number of shares to split the master key into.
      </li>
      <li>
        <span class="param">secret_threshold</span>
        <span class="param-flags">required</span>
        The number of shares required to reconstruct the master key.
        This must be less than or equal to <code>secret_shares</code>.
      </li>
    </ul>
  </dd>

  <dt>Returns</dt>
  <dd>
    A JSON-encoded object including the master keys and initial root token:

    ```javascript
    {
      "keys": ["one", "two", "three"],
      "root_token": "foo"
    }
    ```

  </dd>
</dl>
