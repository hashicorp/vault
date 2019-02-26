---
layout: "docs"
page_title: "Vault Agent Caching"
sidebar_title: "Caching"
sidebar_current: "docs-agent-caching"
description: |-
  Vault Agent Caching allows client-side caching of responses containing newly
  created tokens and responses containing leased secrets generated off of these
  newly created tokens.
---

# Vault Agent Caching

Vault Agent Caching allows client-side caching of responses containing newly
created tokens and responses containing leased secrets generated off of these
newly created tokens. The renewals of the secrets that are cached are also
managed by the agent.

## Caching and Renewals

Caching and renewals are managed by the agent only under these scenarios.

1. Token creation requests are made through the agent and not directly to the
   Vault server. Login endpoints from various auth methods also fall under this
   category along with token creation endpoints of the token auth method.

2. Leased secret creation requests are made through the agent using tokens that
   are already managed by the agent.

## Using Auto-Auth Token

Clients do not need to provide a Vault token along with the proxied request if
the [auto-auth](/docs/autoauth/index.html) feature is enabled. This feature is
enabled by setting the `use_auto_auth_token` (see below) configuration field.
However, even when enabled, if requests that reach the agent already have a
token attached on them, the attached token will be put to use instead of the
auto-auth token.

-> **Note:** In Vault 1.1-beta, if the request doesn't already contain a Vault
token, then the `auto-auth` token will used to make requests. The resulting
secrets from these `auto-auth` token calls are not cached. They will be in the
non-beta version. To test out the caching scenarios, please make a login
request or a token creation request via the agent. The secrets generated from
these new tokens will get cached.

## Cache Evictions

The eviction of cache entries will occur when the agent fails to renew secrets.
This can happen when the secret that is cached hits it's maximum TTL or if the
renewal results in an error.

Agent also does some best-effort cache evictions by observing specific request
types and response codes. For example, if a token revocation request is made
via the agent and if the request succeeds, then agent evicts all the cache
entries associated with the revoked token. Similar behavior is exercised for
lease revocations as well.

While agent tries to observe some requests and evicts cache entries
automatically, agent is completely unaware of revocations that happen outside of
Agent's context. This is when stale entries are created in the agent.

For managing the stale entries in the cache, an endpoint
`/v1/agent/cache-clear`(see below) is made available to manually evict cache
entries based on some of the criteria.

## API

### Cache Clear

This endpoint clears the cache based on given parameters. To be able to use
this API, some information on how the agent caches values should be known
beforehand. Each response that gets cached in the agent is indexed on some
factors depending on the type of request. Those factors can be the `token` that
is being returned by the response, the `token_accessor` of the token being
returned by the response, the `request_path` that resulted in the response, the
`lease` that is attached to the response, the `namespace` to which the request
belongs to, and a few more. This API exposes some factors through which
associated cache entries are fetched and evicted.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `POST`   | `/v1/agent/cache-clear`      | `200 application/json` |

#### Parameters

- `type` `(strings: required)` - The type of cache entries to evict. Valid
  values are `request_path`, `lease`, `token`, and `token_accessor`, and `all`.
  If the `type` is set to `all`, the entire cache is cleared.

- `value` `(string: required)` - An exact value or the prefix of the value for
  the `type` selected. This parameter is optional when the `type` is set
  to `all`.

- `namespace` `(string: optional)` - The namespace in which to match along with
  the provided request path. This is only applicable when the `type` is set to
  `request_path`.

### Sample Payload

```json
{
  "type": "token",
  "value": "s.rlNjegSKykWcplOkwsjd8bP9"
}
```

### Sample Request

```
$ curl \
    --request POST \
    --data @payload.json \
    http://127.0.0.1:1234/v1/agent/cache-clear
```

## Configuration

The top level `cache` block has two configuration entries:

- `use_auto_auth_token (bool: false)` - If set, the requests made to agent
  without a Vault token will be forwarded to Vault with the auto-auth token
  attached. If the requests already bear a token, this configuration will be
  ignored.

- `listener` `(array of objects: required)` - Configuration for the listeners

### Configuration (Listeners)

These configuration values are common to all Listeners.

- `type` `(string: required)` - The type of the listener to use. Valid values
  are `tcp` and `unix`.
  *Note*: when using HCL this can be used as the key for the block, e.g.
  `listener "tcp" {...}`.

- `address` `(string: required)` - The address for the listener to listen to.
  This can either be a URL path when using `tcp` or a file path when using
  `unix`. For example, `127.0.0.1:1234` or `/path/to/socket`.

- `tls_disable` `(bool: false)` - Specifies if TLS will be disabled.

- `tls_key_file` `(string: optional)` - Specifies the path to the private key
  for the certificate.

- `tls_cert_file` `(string: optional)` - Specifies the path to the certificate
  for TLS.

### Example Configuration

An example configuration, with very contrived values, follows:

```javascript
cache {
  use_auto_auth_token = true

  listener "unix" {
    address = "/path/to/socket"
    tls_disable = true
  }

  listener "tcp" {
    address = "127.0.0.1:8100"
    tls_disable = true
  }
}
```

### Environment Variable

Agent's listener address will be picked up by the CLI through the
`VAULT_AGENT_ADDR` environment variable. This should be a complete URL such as
"http://127.0.0.1:8100".
