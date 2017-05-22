---
layout: "docs"
page_title: "Generic Secret Backend"
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

The backend _never_ removes data on its own; the `ttl` key is merely advisory.

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

The Generic secret backend has a full HTTP API. Please see the
[Generic secret backend API](/api/secret/generic/index.html) for more
details.
