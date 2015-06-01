---
layout: "http"
page_title: "HTTP API: /sys/rekey/"
sidebar_current: "docs-http-rotate-rekey"
description: |-
  The `/sys/rekey/` endpoints are used to rekey the unseal keys for Vault.
---

# /sys/rekey/init

## GET

<dl>
  <dt>Description</dt>
  <dd>
      Reads the configuration and progress of the current rekey attempt.
  </dd>

  <dt>Method</dt>
  <dd>GET</dd>

  <dt>URL</dt>
  <dd>`/sys/rekey/init`</dd>

  <dt>Parameters</dt>
  <dd>
    None
  </dd>

  <dt>Returns</dt>
  <dd>
    If a rekey is started, then "n" is the new shares to generate and "t" is
    the threshold required for the new shares. The "progress" is how many unseal
    keys have been provided for this rekey, where "required" must be reached to
    complete.

    ```javascript
    {
      "started": true,
      "t": 3,
      "n": 5,
      "progress": 1,
      "required": 3
    }
    ```

  </dd>
</dl>

## PUT

<dl>
  <dt>Description</dt>
  <dd>
    Initializes a new rekey attempt. Only a single rekey attempt can take place
    at a time, and changing the parameters of a rekey requires canceling and starting
    a new rekey.
  </dd>

  <dt>Method</dt>
  <dd>PUT</dd>

  <dt>URL</dt>
  <dd>`/sys/rekey/init`</dd>

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
  <dd>`204` response code.
  </dd>
</dl>

## DELETE

<dl>
  <dt>Description</dt>
  <dd>
    Cancels any in-progress rekey. This clears the rekey settings as well as any
    progress made. This must be called to change the parameters of the rekey.
  </dd>

  <dt>Method</dt>
  <dd>DELETE</dd>

  <dt>URL</dt>
  <dd>`/sys/rekey/init`</dd>

  <dt>Parameters</dt>
  <dd>None
  </dd>

  <dt>Returns</dt>
  <dd>`204` response code.
  </dd>
</dl>

# /sys/rekey/update

## PUT

<dl>
  <dt>Description</dt>
  <dd>
    Enter a single master key share to progress the rekey of the Vault.
    If the threshold number of master key shares is reached, Vault
    will complete the rekey. Otherwise, this API must be called multiple
    times until that threshold is met.
  </dd>

  <dt>Method</dt>
  <dd>PUT</dd>

  <dt>URL</dt>
  <dd>`/sys/rekey/update`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">key</span>
        <span class="param-flags">required</span>
        A single master share key.
      </li>
    </ul>
  </dd>

  <dt>Returns</dt>
  <dd>
    A JSON-encoded object indicating completion and if so with the new master keys:

    ```javascript
    {
      "complete": true,
      "keys": ["one", "two", "three"]
    }
    ```

  </dd>
</dl>

