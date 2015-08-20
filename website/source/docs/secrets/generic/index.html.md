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

The generic backend allows for writing keys with arbitrary values. The
only value that is special is the `ttl` key, which can be provided with
any key to restrict the lease duration of the secret. This is useful to ensure
clients periodically renew so that key rolling can be time bounded. Note
that this does not actually expire the data, it is simply a hint to clients
that they should not go longer than the `ttl` value before refreshing the
value locally.

N.B.: Prior to version 0.3, the `ttl` parameter was called `lease`. Both will
work for 0.3, but in 0.4 `lease` will be removed. When providing a `lease` value
in 0.3, both `lease` and `ttl` will be returned with the same data.

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
