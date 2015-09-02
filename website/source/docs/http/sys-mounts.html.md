---
layout: "http"
page_title: "HTTP API: /sys/mounts"
sidebar_current: "docs-http-mounts-mounts"
description: |-
  The '/sys/mounts' endpoint is used manage secret backends in Vault.
---

# /sys/mounts

## GET

<dl>
  <dt>Description</dt>
  <dd>
    Lists all the mounted secret backends. `default_lease_ttl`
    or `max_lease_ttl` values of `0` mean that the system
    defaults are used by this backend.
  </dd>

  <dt>Method</dt>
  <dd>GET</dd>

  <dt>URL</dt>
  <dd>`/sys/mounts`</dd>

  <dt>Parameters</dt>
  <dd>
    None
  </dd>

  <dt>Returns</dt>
  <dd>

    ```javascript
    {
      "aws": {
        "type": "aws",
        "description": "AWS keys",
        "config": {
          "default_lease_ttl": 0,
          "max_lease_ttl": 0
        }
      },

      "sys": {
        "type": "system",
        "description": "system endpoint",
        "config": {
          "default_lease_ttl": 0,
          "max_lease_ttl": 0
        }
      }
    }
    ```

  </dd>
</dl>

<dl>
  <dt>Description</dt>
  <dd>
    List the given secret backends configuration. `default_lease_ttl`
    or `max_lease_ttl` values of `0` mean that the system
    defaults are used by this backend.
  </dd>

  <dt>Method</dt>
  <dd>GET</dd>

  <dt>URL</dt>
  <dd>`/sys/mounts/<mount point>/tune`</dd>

  <dt>Parameters</dt>
  <dd>
    None
  </dd>

  <dt>Returns</dt>
  <dd>

    ```javascript
    {
      "config": {
          "default_lease_ttl": 0,
          "max_lease_ttl": 0
      }
    }
    ```

  </dd>
</dl>

## POST

<dl>
  <dt>Description</dt>
  <dd>
    Mount a new secret backend to the mount point in the URL.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/sys/mounts/<mount point>`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">type</span>
        <span class="param-flags">required</span>
        The name of the backend type, such as "aws"
      </li>
      <li>
        <span class="param">description</span>
        <span class="param-flags">optional</span>
        A human-friendly description of the mount.
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

<dl>
  <dt>Description</dt>
  <dd>
    Tune configuration parameters for a given mount point.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/sys/mounts/<mount point>/tune`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">config</span>
        <span class="param-flags">required</span>
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

## DELETE

<dl>
  <dt>Description</dt>
  <dd>
    Unmount the mount point specified in the URL.
  </dd>

  <dt>Method</dt>
  <dd>DELETE</dd>

  <dt>URL</dt>
  <dd>`/sys/mounts/<mount point>`</dd>

  <dt>Parameters</dt>
  <dd>None
  </dd>

  <dt>Returns</dt>
  <dd>`204` response code.
  </dd>
</dl>
