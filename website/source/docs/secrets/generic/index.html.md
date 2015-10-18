---
layout: "docs"
page_title: "Secret Backend: Generic"
sidebar_current: "docs-secrets-generic"
description: |-
  The generic secret backend can store arbitrary secrets.
---

# Generic Secret Backend

Name: `generic`

The generic secret backend is used to store arbitrary secrets within
the configured physical storage for Vault. If you followed along with
the getting started guide, you interacted with a generic secret backend
via the `secret/` prefix that Vault mounts by default. You can mount as many
of these backends at different mount points as you like.

Writing to a key in the `secret/` backend will replace the old value;
sub-fields are not merged together.

## Quick Start

The generic backend allows for writing keys with arbitrary values. A `ttl` value
can be provided, which affects the duration of generated leases. Specifically,
this can be used as a hint from the writer of a secret to consumers of a secret
that the consumer should wait no more than the `ttl` duration before checking
for a new value. If you expect a secret to change frequently, or if you need
clients to react quickly to a change in the secret's value, specify a low value
of `ttl`. Keep in mind that a low `ttl` value may add significant additional load
to the Vault server if it results in clients accessing the value very frequently.
Also note that setting `ttl` does not actually expire the data; it is
informational only.

N.B.: Prior to version 0.3, the `ttl` parameter was called `lease`. Both will
work for 0.3, but in 0.4 `lease` will be removed.

As an example, we can write a new key "foo" to the generic backend
mounted at "secret/" by default:

```
$ vault write secret/foo \
    zip=zap \
    ttl=1h
Success! Data written to: secret/foo
```

This writes the key with the "zip" field set to "zap" and a one hour TTL.
We can test this by doing a read:

```
$ vault read secret/foo
Key             Value
ttl_seconds     3600
ttl             1h
zip             zap
```

As expected, we get the value previously set back as well as our custom TTL
both as specified and translated to seconds. The TTL has been set to 3600
seconds (one hour) as specified.

## API

### /secret
#### GET

<dl class="api">
  <dt>Description</dt>
  <dd>
    Retrieves the secret at the specified location.
  </dd>

  <dt>Method</dt>
  <dd>GET</dd>

  <dt>URL</dt>
  <dd>`/secret/<path>`</dd>

  <dt>Parameters</dt>
  <dd>
     None
  </dd>

  <dt>Returns</dt>
  <dd>

  ```javascript
  {
    "auth": null,
    "data": {
      "foo": "bar"
    },
    "lease_duration": 2592000,
    "lease_id": "",
    "renewable": false
  }
  ```

  </dd>
</dl>

#### POST/PUT

<dl class="api">
  <dt>Description</dt>
  <dd>
    Stores a secret at the specified location.
  </dd>

  <dt>Method</dt>
  <dd>POST/PUT</dd>

  <dt>URL</dt>
  <dd>`/secret/<path>`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">(key)</span>
        <span class="param-flags">optional</span>
        A key, paired with an associated value, to be held at the
        given location. Multiple key/value pairs can be specified,
        and all will be returned on a read operation.
      </li>
      <li>
        <span class="param">ttl</span>
        <span class="param-flags">optional</span>
        The Time To Live for the entry. This value, converted to
        seconds, is round-tripped on read operations as the
        `lease_duration` parameter. Vault takes no action when this
        value expires; it is only meant as a way for a writer of
        a value to indicate to readers how often they should check
        for new entries.
      </li>
    </ul>
  </dd>

  <dt>Returns</dt>
  <dd>
  A `204` response code.
  </dd>
</dl>
