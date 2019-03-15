---
layout: "docs"
page_title: "Vault Agent Caching"
sidebar_title: "Caching <sup>BETA</sup>"
sidebar_current: "docs-agent-caching"
description: |-
  Vault Agent Caching allows client-side caching of responses containing newly
  created tokens and responses containing leased secrets generated off of these
  newly created tokens.
---

# Vault Agent Caching

Vault Agent Caching allows client-side caching of responses containing newly
created tokens and responses containing leased secrets generated off of these
newly created tokens. The renewals of the cached tokens and leases are also
managed by the agent.

-> **Note:** Vault Agent Caching works best with servers/clusters that are
running on Vault 1.1-beta2 and above due to changes that were introduced
alongside this feature, such as the exposure of the `orphan` field in token
creation responses. Agent caching functionality was tested against changes
introduced within 1.1 and thus full caching capabilities may not behave as
expected when paired with older server versions.

## Caching and Renewals

Response caching and renewals are managed by the agent only under these
specific scenarios.

1. Token creation requests are made through the agent. This means that any
   login operations performed using various auth methods and invoking the token
   creation endpoints of the token auth method via the agent will result in the
   response getting cached by the agent. Responses containing new tokens will
   be cached by the agent only if the parent token is already being managed by
   the agent or if the new token is an orphan token.

2. Leased secret creation requests are made through the agent using tokens that
   are already managed by the agent. This means that any dynamic credentials
   that are issued using the tokens managed by the agent, will be cached and
   its renewals are taken care of.

## Using Auto-Auth Token

Vault Agent allows for easy authentication to Vault in a wide variety of
environments using [Auto-Auth](/docs/agent/autoauth/index.html). By setting the
`use_auto_auth_token` (see below) configuration, clients will not be required
to provide a Vault token to the requests made to the agent. When this
configuration is set, if the request doesn't already bear a token, then the
auto-auth token will be used to forward the request to the Vault server. This
configuration will be overridden if the request already has a token attached,
in which case, the token present in the request will be used to forward the
request to the Vault server.

-> **Note:** In Vault 1.1-beta1, if the request doesn't already contain a Vault
token, then the `auto-auth` token will used to make requests. However, the
resulting secrets from these `auto-auth` token calls are not cached. This
behavior will be changed so that they get cached in the upcoming versions. To
test the caching scenarios in 1.1-beta1, please make login requests or token
creation requests via the agent. These new tokens and their respective leased
secrets will get cached.

## Cache Evictions

The eviction of cache entries pertaining to secrets will occur when the agent
can no longer renew them. This can happen when the secrets hit their maximum
TTL or if the renewals result in errors.

Agent does some best-effort cache evictions by observing specific request types
and response codes. For example, if a token revocation request is made via the
agent and if the forwarded request to the Vault server succeeds, then agent
evicts all the cache entries associated with the revoked token. Similarly, any
lease revocation operation will also be intercepted by the agent and the
respective cache entries will be evicted.

Note that while agent evicts the cache entries upon secret expirations and upon
intercepting revocation requests, it is still possible for the agent to be
completely unaware of the revocations that happen through direct client
interactions with the Vault server. This could potentially lead to stale cache
entries. For managing the stale entries in the cache, an endpoint
`/agent/v1/cache-clear`(see below) is made available to manually evict cache
entries based on some of the query criteria used for indexing the cache entries.

## Request Uniqueness

In order to detect repeat requests and return cached responses, agent will need
to have a way to uniquely identify the requests. This computation as it stands
today takes a simplistic approach (may change in future) of serializing and
hashing the HTTP request along with all the headers and the request body. This
hash value is then used as an index into the cache to check if the response is
readily available. The consequence of this approach is that the hash value for
any request will differ if any data in the request is modified. This has the
side-effect of resulting in false negatives if say, the ordering of the request
parameters are modified. As long as the requests come in without any change,
caching behavior should be consistent. Identical requests with differently
ordered request values will result in duplicated cache entries. A heuristic
assumption that the clients will use consistent mechanisms to make requests,
thereby resulting in consistent hash values per request is the idea upon which
the caching functionality is built upon.

## Renewal Management

The tokens and leases are renewed by the agent using the secret renewer that is
made available via the Vault server's [Go
API](https://godoc.org/github.com/hashicorp/vault/api#Renewer). Agent performs
all operations in memory and does not persist anything to storage. This means
that when the agent is shut down, all the renewal operations are immediately
terminated and there is no way for agent to resume renewals after the fact.
Note that shutting down the agent does not indicate revocations of the secrets,
instead it only means that renewal responsibility for all the valid unrevoked
secrets are no longer performed by the Vault agent.

### Agent CLI

Agent's listener address will be picked up by the CLI through the
`VAULT_AGENT_ADDR` environment variable. This should be a complete URL such as
"http://127.0.0.1:8200".

## API

### Cache Clear

This endpoint clears the cache based on given criteria. To be able to use this
API, some information on how the agent caches values should be known
beforehand. Each response that gets cached in the agent will be indexed on some
factors depending on the type of request. Those factors can be the `token` that
is belonging to the cached response, the `token_accessor` of the token
belonging to the cached response, the `request_path` that resulted in the
cached response, the `lease` that is attached to the cached response, the
`namespace` to which the cached response belongs to, and a few more. This API
exposes some factors through which associated cache entries are fetched and
evicted.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `POST`   | `/agent/v1/cache-clear`      | `200 application/json` |

#### Parameters

- `type` `(strings: required)` - The type of cache entries to evict. Valid
  values are `request_path`, `lease`, `token`, and `token_accessor`, and `all`.
  If the `type` is set to `all`, the entire cache is cleared.

- `value` `(string: required)` - An exact value or the prefix of the value for
  the `type` selected. This parameter is optional when the `type` is set
  to `all`.

- `namespace` `(string: optional)` - This is only applicable when the `type` is set to
  `request_path`. The namespace of which the cache entries to be evicted for
  the given request path.

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
    http://127.0.0.1:1234/agent/v1/cache-clear
```

## Configuration (`cache`)

The top level `cache` block has the following configuration entries:

- `use_auto_auth_token (bool: false)` - If set, the requests made to agent
  without a Vault token will be forwarded to the Vault server with the
  auto-auth token attached. If the requests already bear a token, this
  configuration will be overridden and the token in the request will be used to
  forward the request to the Vault server.

## Configuration (`listener`)

- `listener` `(array of objects: required)` - Configuration for the listeners.

There can be one or more `listener` blocks at the top level.
These configuration values are common to all `listener` blocks.

- `type` `(string: required)` - The type of the listener to use. Valid values
  are `tcp` and `unix`.
  *Note*: when using HCL this can be used as the key for the block, e.g.
  `listener "tcp" {...}`.

- `address` `(string: required)` - The address for the listener to listen to.
  This can either be a URL path when using `tcp` or a file path when using
  `unix`. For example, `127.0.0.1:8200` or `/path/to/socket`. Defaults to
  `127.0.0.1:8200`.

- `tls_disable` `(bool: false)` - Specifies if TLS will be disabled.

- `tls_key_file` `(string: optional)` - Specifies the path to the private key
  for the certificate.

- `tls_cert_file` `(string: optional)` - Specifies the path to the certificate
  for TLS.

### Example Configuration

An example configuration, with very contrived values, follows:

```javascript
auto_auth {
  method {
    type = "aws"
    wrap_ttl = 300
    config = {
      role = "foobar"
    }
  }

  sink {
    type = "file"
    config = {
      path = "/tmp/file-foo"
    }
  }
}

cache {
  use_auto_auth_token = true
}

listener "unix" {
  address = "/path/to/socket"
  tls_disable = true
}

listener "tcp" {
  address = "127.0.0.1:8200"
  tls_disable = true
}

vault {
  address = "http://127.0.0.1:8200"
}
```
