---
layout: "http"
page_title: "HTTP API: /sys/revoke-prefix"
sidebar_current: "docs-http-lease-revoke-prefix"
description: |-
  The `/sys/revoke-prefix` endpoint is used to revoke secrets or tokens based on prefix.
---

# /sys/revoke-prefix

<dl>
  <dt>Description</dt>
  <dd>
    Revoke all secrets (via a lease ID prefix) or tokens (via the tokens' path
    property) generated under a given prefix immediately. This requires `sudo`
    capability and access to it should be tightly controlled as it can be used
    to revoke very large numbers of secrets/tokens at once.
  </dd>

  <dt>Method</dt>
  <dd>PUT</dd>

  <dt>URL</dt>
  <dd>`/sys/revoke-prefix/<path prefix>`</dd>

  <dt>Parameters</dt>
  <dd>None</dd>

  <dt>Returns</dt>
  <dd>A `204` response code.
  </dd>
</dl>
