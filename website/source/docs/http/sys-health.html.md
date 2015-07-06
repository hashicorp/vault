---
layout: "http"
page_title: "HTTP API: /sys/health"
sidebar_current: "docs-http-debug-health"
description: |-
  The '/sys/health' endpoint is used to check the health status of Vault.
---

# /sys/health

<dl>
	<dt>Description</dt>
	<dd>
		Returns the health status of Vault. This matches the semantics of a Consul HTTP health
        check and provides a simple way to monitor the health of a Vault instance.
	</dd>

	<dt>Method</dt>
	<dd>GET</dd>

	<dt>Parameters</dt>
	<dd>
        <ul>
          <li>
            <span class="param">standbyok</span>
            <span class="param-flags">optional</span>
            A query parameter provided to indicate that being a standby should
            still return a 200 status code instead of the standard 429 status code.
          </li>
        </ul>
	</dd>

	<dt>Returns</dt>
	<dd>

```
{
    "initialized": true,
    "sealed": false,
    "standby": false
}
```

    Status Codes:

 * `200` if initialized, unsealed and active.
 * `429` if unsealed and standby.
 * `500` if not initialized or sealed.
	</dd>
</dl>
