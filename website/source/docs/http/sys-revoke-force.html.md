---
layout: "http"
page_title: "HTTP API: /sys/revoke-force"
sidebar_current: "docs-http-lease-revoke-force"
description: |-
  The `/sys/revoke-force` endpoint is used to revoke secrets or tokens based on prefix while ignoring backend errors.
---

# /sys/revoke-force

<dl>
  <dt>Description</dt>
  <dd>
    Revoke all secrets or tokens generated under a given prefix immediately.
    Unlike `/sys/revoke-prefix`, this path ignores backend errors encountered
    during revocation. This is <i>potentially very dangerous</i> and should
    only be used in specific emergency situations where errors in the backend
    or the connected backend service prevent normal revocation. <i>By ignoring
    these errors, Vault abdicates responsibility for ensuring that the issued
    credentials or secrets are properly revoked and/or cleaned up. Access to
    this endpoint should be tightly controlled.</i>
  </dd>

  <dt>Method</dt>
  <dd>PUT</dd>

  <dt>URL</dt>
  <dd>`/sys/revoke-force/<path prefix>`</dd>

  <dt>Parameters</dt>
  <dd>None</dd>

  <dt>Returns</dt>
  <dd>A `204` response code.
  </dd>
</dl>
