---
layout: "api"
page_title: "HTTP API"
sidebar_title: "Overview"
sidebar_current: "api-http-overview"
description: |-
  Vault has an HTTP API that can be used to control every aspect of Vault.
---

# HTTP API

The Vault HTTP API gives you full access to Vault via HTTP. Every aspect of
Vault can be controlled via this API. The Vault CLI uses the HTTP API to access
Vault.

All API routes are prefixed with `/v1/`.

This documentation is only for the v1 API, which is currently the only version.

  ~> **Backwards compatibility:** At the current version, Vault does not yet
  promise backwards compatibility even with the v1 prefix. We'll remove this
  warning when this policy changes. At this point in time the core API (that
  is, `sys/` routes) change very infrequently, but various secrets engines/auth
  methods/etc. sometimes have minor changes to accommodate new features as
  they're developed.

## Transport

The API is expected to be accessed over a TLS connection at all times, with a
valid certificate that is verified by a well-behaved client. It is possible to
disable TLS verification for listeners, however, so API clients should expect
to have to do both depending on user settings.

## Authentication

Once Vault is unsealed, almost every other operation requires a _client token_.
A user may have a client token sent to them.  The client token must be sent as
either the `X-Vault-Token` HTTP Header or as `Authorization` HTTP Header using
the `Bearer <token>` scheme.

Otherwise, a client token can be retrieved via [authentication
backends](/docs/auth/index.html).

Each auth method has one or more unauthenticated login endpoints. These
endpoints can be reached without any authentication, and are used for
authentication to Vault itself. These endpoints are specific to each auth
method.

Responses from auth login methods that generate an authentication token are
sent back to the client via JSON. The resulting token should be saved on the
client or passed via the `X-Vault-Token` or `Authorization` header for future requests.

## Namespaces

If using the [Namespaces](/docs/enterprise/namespaces/index.html) feature, API
operations are relative to the namespace value passed in via the
`X-Vault-Namespace` header. For instance, if the request path is to
`secret/foo`, and the header is set to `ns1/ns2/`, the final request path Vault
uses will be `ns1/ns2/secret/foo`. Note that it is semantically equivalent to
use a full path rather than the `X-Vault-Namespace` header, as the operation in
Vault will always look up the correct namespace based on the final given path.
Thus, it would be equivalent to the above example to set `X-Vault-Namespace` to
`ns1/` and a request path of `ns2/secret/foo`, or to not set
`X-Vault-Namespace` at all and use a request path of `ns1/ns2/secret/foo`.

For example, the following two commands result in equivalent requests:

```shell
$ curl \
    -H "X-Vault-Token: f3b09679-3001-009d-2b80-9c306ab81aa6" \
    -H "X-Vault-Namespace: ns1/ns2/" \
    -X GET \
    http://127.0.0.1:8200/v1/secret/foo
```

```shell
$ curl \
    -H "X-Vault-Token: f3b09679-3001-009d-2b80-9c306ab81aa6" \
    -X GET \
    http://127.0.0.1:8200/v1/ns1/ns2/secret/foo
```

## API Operations

With few documented exceptions, all request body data and response data from
Vault is via JSON. Vault will set the `Content-Type` header appropriately but
does not require that clients set it.

Different plugins implement different APIs according to their functionality.
The examples below are created with the `KVv1` backend, which acts like a very
simple Key/Value store. Read the documentation for a particular backend for
detailed information on its API; this simply provides a general overview.

For `KVv1`, reading a secret via the HTTP API is done by issuing a GET:

```text
/v1/secret/foo
```

This maps to `secret/foo` where `foo` is the key in the `secret/` mount, which
is mounted by default on a fresh Vault install and is of type `kv`.

Here is an example of reading a secret using cURL:

```shell
$ curl \
    -H "X-Vault-Token: f3b09679-3001-009d-2b80-9c306ab81aa6" \
    -X GET \
    http://127.0.0.1:8200/v1/secret/foo
```

A few endpoints consume query parameters via `GET` calls, but only if those
parameters are not sensitive, as some load balancers will log these. Most
endpoints that consume parameters use `POST` instead and put the parameters in
the request body.

You can list secrets as well. To do this, either issue a GET with the query
parameter `list=true`, or you can use the `LIST` HTTP verb. For the `kv`
backend, listing is allowed on directories only, and returns the keys in the
given directory:

```shell
$ curl \
    -H "X-Vault-Token: f3b09679-3001-009d-2b80-9c306ab81aa6" \
    -X LIST \
    http://127.0.0.1:8200/v1/secret/
```

The API documentation uses `LIST` as the HTTP verb, but you can still use `GET`
with the `?list=true` query string.

To use an API that consumes data via request body, issue a `POST` or `PUT`:

```text
/v1/secret/foo
```

with a JSON body like:

```javascript
{
  "value": "bar"
}
```

Here is an example of writing a secret using cURL:

```shell
$ curl \
    -H "X-Vault-Token: f3b09679-3001-009d-2b80-9c306ab81aa6" \
    -H "Content-Type: application/json" \
    -X POST \
    -d '{"value":"bar"}' \
    http://127.0.0.1:8200/v1/secret/baz
```

Vault currently considers `PUT` and `POST` to be synonyms. Rather than trust a
client's stated intentions, Vault backends can implement an existence check to
discover whether an operation is actually a create or update operation based on
the data already stored within Vault. This makes permission management via ACLs
more flexible.

For more examples, please look at the Vault API client.

## Help

To retrieve the help for any API within Vault, including mounted backends, auth
methods, etc. then append `?help=1` to any URL. If you have valid permission to
access the path, then the help text will be return a markdown-formatted block in the `help` attribute of the response.

Additionally, with the [OpenAPI generation](/api/system/internal-specs-openapi.html) in Vault, you will get back a small
OpenAPI document in the `openapi` attribute. This document is relevant for the path you're looking up and any paths under it - also note paths in the OpenAPI document are relative to the initial path queried.

Example request:

```shell
$ curl \
    -H "X-Vault-Token: f3b09679-3001-009d-2b80-9c306ab81aa6" \
    http://127.0.0.1:8200/v1/secret?help=1
```

Example response: 

```javascript

{
  "help": "## DESCRIPTION\n\nThis backend provides a versioned key-value store. The kv backend reads and\nwrites arbitrary secrets to the storage backend. The secrets are\nencrypted/decrypted by Vault: they are never stored unencrypted in the backend\nand the backend never has an opportunity to see the unencrypted value. Each key\ncan have a configured number of versions, and versions can be retrieved based on\ntheir version numbers.\n\n## PATHS\n\nThe following paths are supported by this backend. To view help for\nany of the paths below, use the help command with any route matching\nthe path pattern. Note that depending on the policy of your auth token,\nyou may or may not be able to access certain paths.\n\n    ^.*$\n\n\n    ^config$\n        Configures settings for the KV store\n\n    ^data/(?P<path>.*)$\n        Write, Read, and Delete data in the Key-Value Store.\n\n    ^delete/(?P<path>.*)$\n        Marks one or more versions as deleted in the KV store.\n\n    ^destroy/(?P<path>.*)$\n        Permanently removes one or more versions in the KV store\n\n    ^metadata/(?P<path>.*)$\n        Configures settings for the KV store\n\n    ^undelete/(?P<path>.*)$\n        Undeletes one or more versions from the KV store.",
  "openapi": {
    "openapi": "3.0.2",
    "info": {
      "title": "HashiCorp Vault API",
      "description": "HTTP API that gives you full access to Vault. All API routes are prefixed with `/v1/`.",
      "version": "1.0.0",
      "license": {
        "name": "Mozilla Public License 2.0",
        "url": "https://www.mozilla.org/en-US/MPL/2.0"
      }
    },
    "paths": {
      "/.*": {},
      "/config": {
        "description": "Configures settings for the KV store",
        "x-vault-create-supported": true,
        "get": {
          "summary": "Read the backend level settings.",
          "tags": [
            "secrets"
          ],
          "responses": {
            "200": {
              "description": "OK"
            }
          }
        },
     ...[output truncated]...
     }
  }
}
```


## Error Response

A common JSON structure is always returned to return errors:

```javascript
{
  "errors": [
    "message",
    "another message"
  ]
}
```

This structure will be sent down for any HTTP status greater than
or equal to 400.

## HTTP Status Codes

The following HTTP status codes are used throughout the API. Vault tries to
adhere to these whenever possible, but in some cases may not -- feel free to
file a bug in that case to point our attention to it!

~> *Note*: Applications should be prepared to accept both `200` and `204` as
success. `204` is simply an indication that there is no response body to parse,
but API endpoints that indicate that they return a `204` may return a `200` if
warnings are generated during the operation.

- `200` - Success with data.
- `204` - Success, no data returned.
- `400` - Invalid request, missing or invalid data.
- `403` - Forbidden, your authentication details are either incorrect, you
  don't have access to this feature, or - if CORS is enabled - you made a
  cross-origin request from an origin that is not allowed to make such
  requests.
- `404` - Invalid path. This can both mean that the path truly doesn't exist or
  that you don't have permission to view a specific path. We use 404 in some
  cases to avoid state leakage.
- `429` - Default return code for health status of standby nodes. This will
  likely change in the future.
- `473` - Default return code for health status of performance standby nodes.
- `500` - Internal server error. An internal error has occurred, try again
  later. If the error persists, report a bug.
- `502` - A request to Vault required Vault making a request to a third party;
  the third party responded with an error of some kind.
- `503` - Vault is down for maintenance or is currently sealed.  Try again
  later.

## Limits

A maximum request size of 32MB is imposed to prevent a denial of service attack
with arbitrarily large requests; this can be tuned per `listener` block in
Vault's server configuration file.
