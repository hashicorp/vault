---
layout: "http"
page_title: "HTTP API: /sys/policy"
sidebar_current: "docs-http-auth-policy"
description: |-
  The `/sys/policy` endpoint is used to manage ACL policies in Vault.
---

# /sys/policy

## GET

<dl>
  <dt>Description</dt>
  <dd>
    Lists all the available policies.
  </dd>

  <dt>Method</dt>
  <dd>GET</dd>

  <dt>Parameters</dt>
  <dd>
    None
  </dd>

  <dt>Returns</dt>
  <dd>

    ```javascript
    {
      "policies": ["root", "deploy"]
    }
    ```

  </dd>
</dl>

# /sys/policy/

## GET

<dl>
  <dt>Description</dt>
  <dd>
    Retrieve the rules for the named policy.
  </dd>

  <dt>Method</dt>
  <dd>GET</dd>

  <dt>URL</dt>
  <dd>`/sys/policy/<name>`</dd>

  <dt>Parameters</dt>
  <dd>
    None
  </dd>

  <dt>Returns</dt>
  <dd>

    ```javascript
    {
      "rules": "path..."
    }
    ```

  </dd>
</dl>


## PUT

<dl>
  <dt>Description</dt>
  <dd>
    Add or update a policy. Once a policy is updated, it takes effect
    immediately to all associated users.
  </dd>

  <dt>Method</dt>
  <dd>PUT</dd>

  <dt>URL</dt>
  <dd>`/sys/policy/<name>`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">rules</span>
        <span class="param-flags">required</span>
        The policy document.
      </li>
    </ul>
  </dd>

  <dt>Returns</dt>
  <dd>`204` response code.
  </dd>
</dl>

## DELETE

<dl>
  <dt>Description</dt>
  <dd>
    Delete the policy with the given name. This will immediately
    affect all associated users.
  </dd>

  <dt>Method</dt>
  <dd>DELETE</dd>

  <dt>URL</dt>
  <dd>`/sys/policy/<name>`</dd>

  <dt>Parameters</dt>
  <dd>None
  </dd>

  <dt>Returns</dt>
  <dd>`204` response code.
  </dd>
</dl>
