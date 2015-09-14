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
via the `secret/` prefix that Vault mounts by default.

Writing to a key in the `secret/` backend will replace the old value,
the sub-fields are not merged together.

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
$ vault write secret/foo zip=zap ttl=1h
Success! Data written to: secret/foo
```

This writes the key with the "zip" field set to "zap" and a one hour lease. We can test
this by doing a read:

```
$ vault read secret/foo
Key           	Value
lease_id      	secret/foo/e4514713-d5d9-fb14-4177-97a7f7f64518
lease_duration	3600
ttl		1h
zip           	zap
```

As expected, we get the value previously set back as well as our custom TTL.
The lease_duration has been set to 3600 seconds (one hour) as specified.
