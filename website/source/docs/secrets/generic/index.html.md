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
only value that is special is the `lease` key, which can be provided with
any key to restrict the lease time of the secret. This is useful to ensure
clients periodically renew so that key rolling can be time bounded.

As an example, we can write a new key "foo" to the generic backend
mounted at "secret/" by default:

```
$ vault write secret/foo zip=zap lease=1h
Success! Data written to: secret/foo
```

This writes the key with the "zip" field set to "zap" and a one hour lease. We can test
this by doing a read:

```
$ vault read secret/foo
Key           	Value
lease_id      	secret/foo/e4514713-d5d9-fb14-4177-97a7f7f64518
lease_duration	3600
lease         	1h
zip           	zap
```

As expected, we get the value previously set back as well as our custom lease.
The lease_duration has been set to 3600 seconds, or one hour as specified.

