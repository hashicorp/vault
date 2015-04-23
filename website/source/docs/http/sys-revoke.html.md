---
layout: "http"
page_title: "HTTP API: /sys/revoke"
sidebar_current: "docs-http-lease-revoke-single"
description: |-
  The `/sys/revoke` endpoint is used to revoke secrets.
---

# /sys/revoke

<dl>
  <dt>Description</dt>
  <dd>
    Revoke a secret immediately.
  </dd>

  <dt>Method</dt>
  <dd>PUT</dd>

  <dt>URL</dt>
  <dd>`/sys/revoke/<lease id>`</dd>

  <dt>Parameters</dt>
  <dd>None</dd>

  <dt>Returns</dt>
  <dd>A `204` response code.
  </dd>
</dl>
