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
        The new mount point. This can be the same
        as `from` if you simply want to change
        backend configuration with the `config`
        parameter.
      </li>
      <li>
        <span class="param">config</span>
        <span class="param-flags">optional</span>
        Config options for this mount. This is an object with
        two possible values: `default_lease_ttl` and
        `max_lease_ttl`. These control the default and
        maximum lease time-to-live, respectively. If set
        on a specific mount, this overrides the global
        defaults.
    </li>
    </ul>
  </dd>

  <dt>Returns</dt>
  <dd>`204` response code.
  </dd>
</dl>
