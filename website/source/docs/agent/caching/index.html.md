---
layout: "docs"
page_title: "Vault Agent Caching"
sidebar_title: "Caching"
sidebar_current: "docs-agent-caching"
description: |-
  Vault Agent's Caching functionality allows side-client caching of tokens and
  secrets.
---

# Vault Agent Caching

Vault Agent's Caching functionality allows side-client caching of tokens and
secrets.

## Functionality

Caching of tokens and its associated leased secrets are applicable under the
following scenarios:

1. Caching is performed on tokens created by authentication requests proxied
   through the agent, as well as any leased secrets that these tokens generate
   as long as the secret creation request is also proxied through the agent.

2. Similarly, caching also applies to any leased secrets created by the token
   generated from [auto-auth](/docs/agent/autoauth/index.html) if that's
   enabled. For this second case, the proxied request will use the auto-auth
   token if no token is explicitly provided (e.g. via the `X-Vault-Token`
   header).

### Eviction

Eviction of cached entries will occur automatically upon the expiration of the
token's or lease's TTL. A token's expiration will trigger any of its related
leases to be evicted to avoid having any stale entries.

Eviction also occurs when a [token revocation](/api/auth/token/index.html) or
[lease revocation](/api/system/leases.html) request is proxied through the
agent, and said token or lease was kept track of by agent. Token revocation
requests will result in eviction of the token entry as well as any of the leases
created by the token. Lease revocation will result in eviction of the said
lease. Prefix-based revocation will evict all matching leases.

### Manual Eviction

Eviction can also be done manually through the `/agent/v1/cache-clear` endpoint
that's available via the enabled listener(s).

Eviction can be based on an exact value match for token, token accessor, and
lease values, prefix match based on request path, or full cache eviction which
will reset the underlying cache.

The API endpoint accepts the following values in the request body as a  JSON
object:

- `type` `(strings: required)` - The lookup type on the entries to clear from
  the cache. This is use alongside the value parameter. Valid values are
  `request_path`, `lease`, `token`, and `token_accessor`, and `all`

- `value` `(string: required-if-not-all)` - An exact or prefix value in which to
  match the desired entries for eviction.

- `namespace` `(string: optional)` - The namespace of the . This is only
  applicable when the `type` is set to `request_path`.

## Configuration

The top level `cache` block has two configuration entries:

- `use_auto_auth_token` `(bool: false)` - Whether to use the auto-auth token, if
  present, for proxied requests. If set to true, requests made by client will
  use this token unless a token is provided explicitly via `X-Vault-Token`.

- `listener` `(array of objects: required)` - Configuration for the listeners

### Configuration (Listeners)

These configuration values are common to all Listeners:

- `type` `(string: required)` - The type of the listener to use. Valid values
  are `tcp` and `unix`. 
  *Note*: when using HCL this can be used as the key for the block, e.g.
  `listener "tcp" {...}`.

- `address` `(string: required)` - The address for the listener to listen to.
  This can either be a URL path when using `tcp` or a file path when using
  `unix`.

- `tls_disable` `(bool: false)` - Specifies if TLS will be disabled.

- `tls_key_file` `(string: optional)` - Specifies the path to the private key
  for the certificate.

- `tls_cert_file` `(string: optional)` - Specifies the path to the certificate
  for TLS.
