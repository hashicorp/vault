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
      "initialized": true
    }
    ```

  </dd>
</dl>

## PUT

<dl>
  <dt>Description</dt>
  <dd>
     Initializes a new Vault. The Vault must not have been previously
     initialized. The recovery options, as well as the stored shares option, are
     only available when using Vault HSM.
  </dd>

  <dt>Method</dt>
  <dd>PUT</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">root_token_pgp_key</span>
        <span class="param-flags">optional</span>
        A PGP public key used to encrypt the initial root token. The key
        must be base64-encoded from its original binary representation.
      </li>
      <li>
        <span class="param">secret_shares</span>
        <span class="param-flags">required</span>
        The number of shares to split the master key into.
      </li>
      <li>
        <span class="param">secret_threshold</span>
        <span class="param-flags">required</span>
        The number of shares required to reconstruct the master key. This must
        be less than or equal to <code>secret_shares</code>. If using Vault HSM
        with auto-unsealing, this value must be the same as
        <code>secret_shares</code>.
      </li>
      <li>
        <span class="param">pgp_keys</span>
        <span class="param-flags">optional</span>
        An array of PGP public keys used to encrypt the output unseal keys.
        Ordering is preserved. The keys must be base64-encoded from their
        original binary representation. The size of this array must be the
        same as <code>secret_shares</code>.
      </li>
      <li>
        <span class="param">stored_shares</span>
        <span class="param-flags">required</span>
        The number of shares that should be encrypted by the HSM and stored for
        auto-unsealing (Vault HSM only). Currently must be the same as
        <code>secret_shares</code>.
      </li>
      <li>
        <span class="param">recovery_shares</span>
        <span class="param-flags">required</span>
        The number of shares to split the recovery key into (Vault HSM only).
      </li>
      <li>
        <span class="param">recovery_threshold</span>
        <span class="param-flags">required</span>
        The number of shares required to reconstruct the recovery key (Vault
        HSM only). This must be less than or equal to
        <code>recovery_shares</code>.
      </li>
      <li>
        <span class="param">recovery_pgp_keys</span>
        <span class="param-flags">optional</span>
        An array of PGP public keys used to encrypt the output recovery keys
        (Vault HSM only). Ordering is preserved. The keys must be
        base64-encoded from their original binary representation. The size of
        this array must be the same as <code>recovery_shares</code>.
      </li>
    </ul>
  </dd>

  <dt>Returns</dt>
  <dd>
    A JSON-encoded object including the (possibly encrypted, if
    <code>pgp_keys</code> was provided) master keys, base 64 encoded master keys and initial root token:

    ```javascript
    {
      "keys": ["one", "two", "three"],
      "keys_base64": ["cR9No5cBC", "F3VLrkOo", "zIDSZNGv"],
      "root_token": "foo"
    }
    ```

  </dd>

  <dt>See Also</dt>
  <dd>
    For more information on the PGP/Keybase.io process please see the
    [Vault GPG and Keybase integration documentation](/docs/concepts/pgp-gpg-keybase.html).
  </dd>
</dl>
