---
layout: "http"
page_title: "HTTP API: /sys/revoke-prefix"
sidebar_current: "docs-http-lease-revoke-prefix"
description: |-
  The `/sys/revoke-prefix` endpoint is used to revoke secrets based on prefix.
---

# /sys/revoke-prefix

<dl>
  <dt>Description</dt>
  <dd>
    Revoke all secrets generated under a given prefix immediately.
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
