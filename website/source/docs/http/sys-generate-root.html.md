---
layout: "http"
page_title: "HTTP API: /sys/generate-root/"
sidebar_current: "docs-http-sys-generate-root"
description: |-
  The `/sys/generate-root/` endpoints are used to create a new root key for Vault.
---

# /sys/generate-root/attempt

## GET

<dl>
  <dt>Description</dt>
  <dd>
      Reads the configuration and progress of the current root generation
      attempt.
  </dd>

  <dt>Method</dt>
  <dd>GET</dd>

  <dt>URL</dt>
  <dd>`/sys/generate-root/attempt`</dd>

  <dt>Parameters</dt>
  <dd>
    None
  </dd>

  <dt>Returns</dt>
  <dd>
    If a root generation is started, `progress` is how many unseal keys have
    been provided for this generation attempt, where `required` must be reached
    to complete. The `nonce` for the current attempt and whether the attempt is
    complete is also displayed. If a PGP key is being used to encrypt the final
    root token, its fingerprint will be returned. Note that if an OTP is being
    used to encode the final root token, it will never be returned.

    ```javascript
    {
      "started": true,
      "nonce": "2dbd10f1-8528-6246-09e7-82b25b8aba63",
      "progress": 1,
      "required": 3,
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
    Initializes a new root generation attempt. Only a single root generation
    attempt can take place at a time. One (and only one) of `otp` or `pgp_key`
    are required.
  </dd>

  <dt>Method</dt>
  <dd>PUT</dd>

  <dt>URL</dt>
  <dd>`/sys/generate-root/attempt`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">otp</span>
        <span class="param-flags">optional</span>
        A base64-encoded 16-byte value. The raw bytes of the token will be
        XOR'd with this value before being returned to the final unseal key
        provider.
      </li>
      <li>
        <span class="param">pgp_key</span>
        <span class="param-flags">optional</span>
        A base64-encoded PGP public key. The raw bytes of the token will be
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
      "nonce": "2dbd10f1-8528-6246-09e7-82b25b8aba63",
      "progress": 1,
      "required": 3,
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
    Cancels any in-progress root generation attempt. This clears any progress
    made. This must be called to change the OTP or PGP key being used.
  </dd>

  <dt>Method</dt>
  <dd>DELETE</dd>

  <dt>URL</dt>
  <dd>`/sys/generate-root/attempt`</dd>

  <dt>Parameters</dt>
  <dd>None
  </dd>

  <dt>Returns</dt>
  <dd>`204` response code.
  </dd>
</dl>

# /sys/generate-root/update

## PUT

<dl>
  <dt>Description</dt>
  <dd>
    Enter a single master key share to progress the root generation attempt.
    If the threshold number of master key shares is reached, Vault will
    complete the root generation and issue the new token.  Otherwise, this API
    must be called multiple times until that threshold is met. The attempt
    nonce must be provided with each call.
  </dd>

  <dt>Method</dt>
  <dd>PUT</dd>

  <dt>URL</dt>
  <dd>`/sys/generate-root/update`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">key</span>
        <span class="param-flags">required</span>
        A single master share key.
      </li>
      <li>
        <span class="param">nonce</span>
        <span class="param-flags">required</span>
        The nonce of the attempt.
      </li>
    </ul>
  </dd>

  <dt>Returns</dt>
  <dd>
    A JSON-encoded object indicating the attempt nonce, and completion status,
    and the encoded root token, if the attempt is complete.

    ```javascript
    {
      "started": true,
      "nonce": "2dbd10f1-8528-6246-09e7-82b25b8aba63",
      "progress": 3,
      "required": 3,
      "pgp_fingerprint": "",
      "complete": true,
      "encoded_root_token": "FPzkNBvwNDeFh4SmGA8c+w=="
    }
    ```

  </dd>
</dl>
