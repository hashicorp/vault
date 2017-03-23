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

This backend differs from the `generic` backend in that the `generic` backend's
values are accessible to any token with read privileges on that path. In
`cubbyhole`, paths are scoped per token; no token can access another token's
cubbyhole, whether to read, write, list, or for any other operation. When the
token expires, its cubbyhole is destroyed.

Also unlike the `generic` backend, because the cubbyhole's lifetime is linked
to that of an authentication token, there is no concept of a TTL or refresh
interval for values contained in the token's cubbyhole.

Writing to a key in the `cubbyhole` backend will replace the old value;
the sub-fields are not merged together.

## Response Wrapping

Starting in Vault 0.6, almost any response (except those from `sys/` endpoints)
from Vault can be wrapped (see the [Response
Wrapping](/docs/concepts/response-wrapping.html)
concept page for details).

The TTL for the token is set by the client using the `X-Vault-Wrap-TTL` header
and can be either an integer number of seconds or a string duration of seconds
(`15s`), minutes (`20m`), or hours (`25h`). When using the Vault CLI, you can
set this via the `-wrap-ttl` parameter. Response wrapping is per-request; it is
the presence of a value in this header that activates wrapping of the response.

If a client requests wrapping:

1. The original response is serialized to JSON
2. A new single-use token is generated with a TTL as supplied by the client
3. Internally, the original response JSON is stored in the single-use token's
   cubbyhole.
4. A new response is generated, with the token ID and the token TTL stored in
   the new response's `wrap_info` dict
5. The new response is returned to the caller

To get the original value, if using the API, perform a write on
`sys/wrapping/unwrap`, passing in the wrapping token ID. The original value
will be returned.

If using the CLI, passing the wrapping token's ID to the `vault unwrap` command
will return the original value; `-format` and `-field` can be set like with
`vault read`.

If the original response is an authentication response containing a token, the
token's accessor will be made available to the caller. This allows a privileged
caller to generate tokens for clients and be able to manage the tokens'
lifecycle while not being exposed to the actual client token IDs.

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
