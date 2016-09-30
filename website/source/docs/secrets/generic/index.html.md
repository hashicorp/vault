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

Writing to a key in the `generic` backend will replace the old value;
sub-fields are not merged together.

This backend honors the distinction between the `create` and `update`
capabilities inside ACL policies.

**Note**: Path and key names are _not_ obfuscated or encrypted; only the values
set on keys are. You should not store sensitive information as part of a
secret's path.

## Quick Start

The generic backend allows for writing keys with arbitrary values. When data is
returned, the `lease_duration` field (in the API JSON) or `refresh_interval`
field (on the CLI) gives a hint as to how often a reader should look for a new
value. This comes from the value of the `default_lease_ttl` set on the mount,
or the system value.

There is one piece of special data handling: if a `ttl` key is provided, it
will be treated as normal data, but on read the backend will attempt to parse
it as a duration (either as a string like `1h` or an integer number of seconds
like `3600`). If successful, the backend will use this value in place of the
normal `lease_duration`. However, the given value will also still be returned
exactly as specified, so you are free to use that key in any way that you like
if it fits your input data.

As an example, we can write a new key "foo" to the generic backend mounted at
"secret/" by default:

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
Key               Value
---               -----
refresh_interval  3600
ttl               1h
zip               zap
```

As expected, we get the values previously set back as well as our custom TTL
both as specified and translated to seconds. The duration has been set to 3600
seconds (one hour) as specified.

## API

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
    "lease_duration": 2764800,
    "lease_id": "",
    "renewable": false
  }
  ```

  </dd>
</dl>

#### LIST

<dl class="api">
  <dt>Description</dt>
  <dd>
    Returns a list of key names at the specified location. Folders are
    suffixed with `/`. The input must be a folder; list on a file will not
    return a value. Note that no policy-based filtering is performed on keys;
    do not encode sensitive information in key names. The values themselves
    are not accessible via this command.
  </dd>

  <dt>Method</dt>
  <dd>LIST/GET</dd>

  <dt>URL</dt>
  <dd>`/secret/<path>` (LIST) or `/secret/<path>?list=true` (GET)</dd>

  <dt>Parameters</dt>
  <dd>
     None
  </dd>

  <dt>Returns</dt>
  <dd>
  The example below shows output for a query path of `secret/` when there are
  secrets at `secret/foo` and `secret/foo/bar`; note the difference in the two
  entries.

  ```javascript
  {
    "auth": null,
    "data": {
      "keys": ["foo", "foo/"]
    },
    "lease_duration": 2764800,
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
    Stores a secret at the specified location. If the value does not yet exist,
    the calling token must have an ACL policy granting the `create` capability.
    If the value already exists, the calling token must have an ACL policy
    granting the `update` capability.
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
        A key, paired with an associated value, to be held at the given
        location. Multiple key/value pairs can be specified, and all will be
        returned on a read operation. A key called `ttl` will trigger some
        special behavior; see above for details.
      </li>
    </ul>
  </dd>

  <dt>Returns</dt>
  <dd>
  A `204` response code.
  </dd>
</dl>

#### DELETE

<dl class="api">
  <dt>Description</dt>
  <dd>
    Deletes the secret at the specified location.
  </dd>

  <dt>Method</dt>
  <dd>DELETE</dd>

  <dt>URL</dt>
  <dd>`/secret/<path>`</dd>

  <dt>Parameters</dt>
  <dd>
     None
  </dd>

  <dt>Returns</dt>
  <dd>
  A `204` response code.
  </dd>
</dl>
