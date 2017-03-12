---
layout: "http"
page_title: "HTTP API: /sys/generate-share/"
sidebar_current: "docs-http-sys-generate-share"
description: |-
  The `/sys/generate-share/` endpoints are used to create a new share for reconstructing the master key.
---

# /sys/generate-share/attempt

## GET

<dl>
  <dt>Description</dt>
  <dd>
      Reads the configuration and progress of the current share generation
      attempt.
  </dd>

  <dt>Method</dt>
  <dd>GET</dd>

  <dt>URL</dt>
  <dd>`/sys/generate-share/attempt`</dd>

  <dt>Parameters</dt>
  <dd>
    None
  </dd>

  <dt>Returns</dt>
  <dd>
    If a share generation is started, `progress` is how many unseal keys have
    been provided for this generation attempt, where `required` must be reached
    to complete. If a PGP key is being used to encrypt the final
    key share, its fingerprint will be returned.

    ```javascript
    {
      "started": true,
      "progress": 1,
      "required": 3,
      "key": "",
      "key_base64": "",
      "pgp_fingerprint": "",
      "complete": false
    }
    ```

  </dd>
</dl>

## PUT

<dl>
  <dt>Description</dt>
  <dd>
    Initializes a new share generation attempt. Only a single share generation
    attempt can take place at a time.
  </dd>

  <dt>Method</dt>
  <dd>PUT</dd>

  <dt>URL</dt>
  <dd>`/sys/generate-share/attempt`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">pgp_key</span>
        <span class="param-flags">optional</span>
        A base64-encoded PGP public key. The raw bytes of the share will be
        encrypted with this value before being returned to the final unseal key
        provider.
      </li>
    </ul>
  </dd>

  <dt>Returns</dt>
  <dd>
    The current progress.

    ```javascript
    {
      "started": true,
      "progress": 1,
      "required": 3,
      "key": "",
      "key_bas64": "",
      "pgp_fingerprint": "816938b8a29146fbe245dd29e7cbaf8e011db793",
      "complete": false
    }
    ```

  </dd>
</dl>

## DELETE

<dl>
  <dt>Description</dt>
  <dd>
    Cancels any in-progress share generation attempt. This clears any progress
    made. This must be called to change the PGP key being used.
  </dd>

  <dt>Method</dt>
  <dd>DELETE</dd>

  <dt>URL</dt>
  <dd>`/sys/generate-share/attempt`</dd>

  <dt>Parameters</dt>
  <dd>None
  </dd>

  <dt>Returns</dt>
  <dd>`204` response code.
  </dd>
</dl>

# /sys/generate-share/update

## PUT

<dl>
  <dt>Description</dt>
  <dd>
    Enter a single master key share to progress the share generation attempt.
    If the threshold number of master key shares is reached, Vault will
    complete the share generation and issue the new share.  Otherwise, this API
    must be called multiple times until that threshold is met.
  </dd>

  <dt>Method</dt>
  <dd>PUT</dd>

  <dt>URL</dt>
  <dd>`/sys/generate-share/update`</dd>

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
    A JSON-encoded object indicating the completion status,
    and the encoded share, if the attempt is complete.

    ```javascript
    {
      "started": true,
      "nonce": "2dbd10f1-8528-6246-09e7-82b25b8aba63",
      "progress": 3,
      "required": 3,
      "pgp_fingerprint": "",
      "complete": true,
      "key": "four",
      "key_base64": "Zm91cg=="
    }
    ```

  </dd>
</dl>
