---
layout: "http"
page_title: "HTTP API: /sys/remount"
sidebar_current: "docs-http-mounts-remount"
description: |-
  The '/sys/remount' endpoint is used remount a mounted backend to a new endpoint.
---

# /sys/remount

<dl>
  <dt>Description</dt>
  <dd>
    Remount an already-mounted backend to a new mount point.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">from</span>
        <span class="param-flags">required</span>
        The previous mount point.
      </li>
      <li>
        <span class="param">to</span>
        <span class="param-flags">required</span>
        The new mount point.
      </li>
    </ul>
  </dd>

  <dt>Returns</dt>
  <dd>`204` response code.
  </dd>
</dl>
