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
        Returns the health status of Vault. This matches the semantics of a
        Consul HTTP health check and provides a simple way to monitor the
        health of a Vault instance.
    </dd>

    <dt>Method</dt>
    <dd>GET/HEAD</dd>

    <dt>Parameters</dt>
    <dd>
        <ul>
          <li>
            <span class="param">standbyok</span>
            <span class="param-flags">optional</span>
            A query parameter provided to indicate that being a standby should
            still return the active status code instead of the standby code
          </li>
          <li>
            <span class="param">activecode</span>
            <span class="param-flags">optional</span>
            A query parameter provided to indicate the status code that should
            be returned for an active node instead of the default of `200`
          </li>
          <li>
            <span class="param">standbycode</span>
            <span class="param-flags">optional</span>
            A query parameter provided to indicate the status code that should
            be returned for a standby node instead of the default of `429`
          </li>
          <li>
            <span class="param">sealedcode</span>
            <span class="param-flags">optional</span>
            A query parameter provided to indicate the status code that should
            be returned for a sealed node instead of the default of `503`
          </li>
          <li>
            <span class="param">uninitcode</span>
            <span class="param-flags">optional</span>
            A query parameter provided to indicate the status code that should
            be returned for an uninitialized Vault instead of the default of
            `501`
          </li>
        </ul>
    </dd>

    <dt>Returns (only with GET)</dt>
    <dd>

    ```javascript
{
  "cluster_id": "c9abceea-4f46-4dab-a688-5ce55f89e228",
  "cluster_name": "vault-cluster-5515c810",
  "version": "0.6.1-dev"
  "server_time_utc": 1469555798,
  "standby": false,
  "sealed": false,
  "initialized": true
}
    ```

    Default Status Codes (GET/HEAD):

 * `200` if initialized, unsealed, and active.
 * `429` if unsealed and standby.
 * `500` if sealed, or if not initialized.
	</dd>
</dl>
