---
layout: "docs"
page_title: "Cubbyhole Secret Backend"
sidebar_current: "docs-secrets-cubbyhole"
description: |-
  The cubbyhole secret backend can store arbitrary secrets scoped to a single token.
---

# Cubbyhole Secret Backend

Name: `cubbyhole`

The `cubbyhole` secret backend is used to store arbitrary secrets within
the configured physical storage for Vault. It is mounted at the `cubbyhole/`
prefix by default and cannot be mounted elsewhere or removed.

This backend differs from the `kv` backend in that the `kv` backend's
values are accessible to any token with read privileges on that path. In
`cubbyhole`, paths are scoped per token; no token can access another token's
cubbyhole, whether to read, write, list, or for any other operation. When the
token expires, its cubbyhole is destroyed.

Also unlike the `kv` backend, because the cubbyhole's lifetime is linked
to that of an authentication token, there is no concept of a TTL or refresh
interval for values contained in the token's cubbyhole.

Writing to a key in the `cubbyhole` backend will replace the old value;
the sub-fields are not merged together.

## Quick Start

The `cubbyhole` backend allows for writing keys with arbitrary values.

As an example, we can write a new key "foo" to the `cubbyhole` backend, which
is mounted at `cubbyhole/`:

```
$ vault write cubbyhole/foo \
    zip=zap
Success! Data written to: cubbyhole/foo
```

This writes the key with the "zip" field set to "zap". We can test this by doing
a read:

```
$ vault read cubbyhole/foo
Key           	Value
zip           	zap
```

As expected, the value previously set is returned to us.

## API

The Cubbyhole secret backend has a full HTTP API. Please see the
[Cubbyhole secret backend API](/api/secret/cubbyhole/index.html) for more
details.
