---
layout: "http"
page_title: "HTTP API: /sys/rotate"
sidebar_current: "docs-http-rotate-rotate"
description: |-
  The `/sys/rotate` endpoint is used to rotate the encryption key.
---

# /sys/rotate

## PUT

<dl>
  <dt>Description</dt>
  <dd>
    Trigger a rotation of the backend encryption key. This is the key that is used
    to encrypt data written to the storage backend, and is not provided to operators.
    This operation is done online. Future values are encrypted with the new key, while
    old values are decrypted with previous encryption keys.
  </dd>

  <dt>Method</dt>
  <dd>PUT</dd>

  <dt>URL</dt>
  <dd>`/sys/rotate`</dd>

  <dt>Parameters</dt>
  <dd>
    None
  </dd>

  <dt>Returns</dt>
  <dd>`204` response code.
  </dd>
</dl>

