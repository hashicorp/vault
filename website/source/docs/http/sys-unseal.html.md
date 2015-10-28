---
layout: "http"
page_title: "HTTP API: /sys/seal-unseal"
sidebar_current: "docs-http-seal-unseal"
description: |-
  The '/sys/seal-unseal' endpoint is used to unseal the Vault.
---

# /sys/unseal

<dl>
  <dt>Description</dt>
  <dd>
    Enter a single master key share to progress the unsealing of the Vault.
    If the threshold number of master key shares is reached, Vault
    will attempt to unseal the Vault. Otherwise, this API must be
    called multiple times until that threshold is met.<br/><br/>Either
    the `key` or `reset` parameter must be provided; if both are provided,
    `reset` takes precedence.
  </dd>

  <dt>Method</dt>
  <dd>PUT</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">key</span>
        <span class="param-flags">optional</span>
        A single master share key.
      </li>
      <li>
        <span class="param">reset</span>
        <span class="param-flags">optional</span>
        A boolean; if true, the previously-provided unseal keys are discarded
        from memory and the unseal process is reset.
      </li>
    </ul>
  </dd>
  <dt>Returns</dt>
  <dd>The same result as `/sys/seal-status`.
  </dd>
</dl>
