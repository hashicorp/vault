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
    If a rekey is started, then `n` is the new shares to generate and `t` is
    the threshold required for the new shares. `progress` is how many unseal
    keys have been provided for this rekey, where `required` must be reached to
    complete. The `nonce` for the current rekey operation is also displayed. If
    PGP keys are being used to encrypt the final shares, the key fingerprints
    and whether the final keys will be backed up to physical storage will also
    be displayed.

    ```javascript
    {
      "started": true,
      "nonce": "2dbd10f1-8528-6246-09e7-82b25b8aba63",
      "t": 3,
      "n": 5,
      "progress": 1,
      "required": 3,
      "pgp_fingerprints": ["abcd1234"],
      "backup": true
    }
    ```

  </dd>
</dl>

## PUT

<dl>
  <dt>Description</dt>
  <dd>
    Initializes a new rekey attempt. Only a single rekey attempt can take place
    at a time, and changing the parameters of a rekey requires canceling and
    starting a new rekey, which will also provide a new nonce.
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
      <li>
        <spam class="param">pgp_keys</span>
        <span class="param-flags">optional</spam>
        An array of PGP public keys used to encrypt the output unseal keys.
        Ordering is preserved. The keys must be base64-encoded from their
        original binary representation. The size of this array must be the
        same as <code>secret_shares</code>.
      </li>
      <li>
        <spam class="param">backup</span>
        <span class="param-flags">optional</spam>
        If using PGP-encrypted keys, whether Vault should also back them up to
        a well-known location in physical storage (`core/unseal-keys-backup`).
        These can then be retrieved and removed via the `sys/rekey/backup`
        endpoint.
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

# /sys/rekey/backup

## GET

<dl>
  <dt>Description</dt>
  <dd>
    Return the backup copy of PGP-encrypted unseal keys. The returned value is
    the nonce of the rekey operation and a map of PGP key fingerprint to
    hex-encoded PGP-encrypted key.
  </dd>

  <dt>Method</dt>
  <dd>GET</dd>

  <dt>URL</dt>
  <dd>`/sys/rekey/backup`</dd>

  <dt>Parameters</dt>
  <dd>
    None
  </dd>

  <dt>Returns</dt>
  <dd>

    ```javascript
    {
      "nonce": "2dbd10f1-8528-6246-09e7-82b25b8aba63",
      "keys": {
        "abcd1234": "..."
      }
    }
    ```

  </dd>
</dl>

## DELETE

<dl>
  <dt>Description</dt>
  <dd>
    Delete the backup copy of PGP-encrypted unseal keys.
  </dd>

  <dt>Method</dt>
  <dd>DELETE</dd>

  <dt>URL</dt>
  <dd>`/sys/rekey/backup`</dd>

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
    times until that threshold is met. The rekey nonce operation must be
    provided with each call.
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
      <li>
        <span class="param">nonce</span>
        <span class="param-flags">required</span>
        The nonce of the rekey operation.
      </li>
    </ul>
  </dd>

  <dt>Returns</dt>
  <dd>
    A JSON-encoded object indicating the rekey operation nonce and completion
    status; if completed, the new master keys are returned. If the keys are
    PGP-encrypted, an array of key fingerprints will also be provided (with the
    order in which the keys were used for encryption) along with whether or not
    the keys were backed up to physical storage:

    ```javascript
    {
      "complete": true,
      "keys": ["one", "two", "three"],
      "nonce": "2dbd10f1-8528-6246-09e7-82b25b8aba63",
      "pgp_fingerprints": ["abcd1234"],
      "keys_base64": ["base64keyvalue"],
      "backup": true
    }
    ```

  </dd>
</dl>
