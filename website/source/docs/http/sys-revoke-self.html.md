---
layout: "http"
page_title: "HTTP API: /sys/revoke-self"
sidebar_current: "docs-http-lease-revoke-self"
description: |-
  The `/sys/revoke-self` endpoint is used for a token to revoke itself.
---

# /sys/revoke-self

<dl>
  <dt>Description</dt>
  <dd>
    Revoke the calling token immediately.
  </dd>

  <dt>Method</dt>
  <dd>PUT</dd>

  <dt>URL</dt>
  <dd>`/sys/revoke-self`</dd>

  <dt>Parameters</dt>
  <dd>None</dd>

  <dt>Returns</dt>
  <dd>A `204` response code.
  </dd>
</dl>
