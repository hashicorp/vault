---
layout: "http"
page_title: "HTTP API: /sys/seal-status"
sidebar_current: "docs-http-seal-status"
description: |-
  The '/sys/seal-status' endpoint is used to check the seal status of a Vault.
---

# /sys/seal-status

<dl>
  <dt>Description</dt>
  <dd>
    Returns the seal status of the Vault.<br/><br/>This is an unauthenticated endpoint.
  </dd>

  <dt>Method</dt>
  <dd>GET</dd>

  <dt>Parameters</dt>
  <dd>
    None
  </dd>

  <dt>Returns</dt>
  <dd>
    The "t" parameter is the threshold, and "n" is the number of shares.

    ```javascript
    {
      "sealed": true,
      "t": 3,
      "n": 5,
      "progress": 2,
      "version": "0.6.1-dev"
    }
    ```
    
    Sample response when Vault is unsealed.
    
    ```javascript
    {
      "sealed": false,
      "t": 3,
      "n": 5,
      "progress": 0,
      "version": "0.6.1-dev",
      "cluster_name": "vault-cluster-d6ec3c7f",
      "cluster_id": "3e8b3fec-3749-e056-ba41-b62a63b997e8"
    }
    ```

  </dd>
</dl>
