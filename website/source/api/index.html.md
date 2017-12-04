---
layout: "api"
page_title: "HTTP API"
sidebar_current: "docs-http-overview"
description: |-
  Vault has an HTTP API that can be used to control every aspect of Vault.
---

# HTTP API

The Vault HTTP API gives you full access to Vault via HTTP. Every
aspect of Vault can be controlled via this API. The Vault CLI uses
the HTTP API to access Vault.

## Version Prefix

All API routes are prefixed with `/v1/`.

This documentation is only for the v1 API.

~> **Backwards compatibility:** At the current version, Vault does
not yet promise backwards compatibility even with the v1 prefix. We'll
remove this warning when this policy changes. We expect we'll reach API
stability by Vault 1.0.

## Transport

The API is expected to be accessed over a TLS connection at
all times, with a valid certificate that is verified by a well
behaved client. It is possible to disable TLS verification for
listeners, however, so API clients should expect to have to do both
depending on user settings.

## Authentication

Once the Vault is unsealed, every other operation requires a _client token_. A
user may have a client token sent to her.  The client token must be sent as the
`X-Vault-Token` HTTP header.

Otherwise, a client token can be retrieved via [authentication
backends](/docs/auth/index.html).

Each authentication backend will have one or more unauthenticated login
endpoints. These endpoints can be reached without any authentication, and are
used for authentication itself. These endpoints are specific to each
authentication backend.

Login endpoints for authentication backends that generate an identity will be
sent down via JSON. The resulting token should be saved on the client or passed
via the `X-Vault-Token` header for future requests.

## Reading, Writing, and Listing Secrets

Different backends implement different APIs according to their functionality.
The examples below are created with the `kv` backend, which acts like a
Key/Value store. Read the documentation for a particular backend for detailed
information on its API; this simply provides a general overview.

Reading a secret via the HTTP API is done by issuing a GET using the
following URL:

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

You can list secrets as well. To do this, either issue a GET with the query
parameter `list=true`, or you can use the LIST HTTP verb. For the `kv`
backend, listing is allowed on directories only, and returns the keys in the
given directory:

```shell
$ curl \
    -H "X-Vault-Token: f3b09679-3001-009d-2b80-9c306ab81aa6" \
    -X GET \
    http://127.0.0.1:8200/v1/secret/?list=true
```

To write a secret, issue a POST on the following URL:

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

Vault currently considers PUT and POST to be synonyms. Rather than trust a
client's stated intentions, Vault backends can implement an existence check to
discover whether an operation is actually a create or update operation based on
the data already stored within Vault.

For more examples, please look at the Vault API client.

## Help

To retrieve the help for any API within Vault, including mounted
backends, credential providers, etc. then append `?help=1` to any
URL. If you have valid permission to access the path, then the help text
will be returned with the following structure:

```javascript
{
  "help": "help text"
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

The following HTTP status codes are used throughout the API.

- `200` - Success with data.
- `204` - Success, no data returned.
- `400` - Invalid request, missing or invalid data.
- `403` - Forbidden, your authentication details are either
   incorrect, you don't have access to this feature, or - if CORS is
   enabled - you made a cross-origin request from an origin that is
   not allowed to make such requests.
- `404` - Invalid path. This can both mean that the path truly
   doesn't exist or that you don't have permission to view a
   specific path. We use 404 in some cases to avoid state leakage.
- `429` - Default return code for health status of standby nodes, indicating a
   warning.
- `500` - Internal server error. An internal error has occurred,
   try again later. If the error persists, report a bug.
- `503` - Vault is down for maintenance or is currently sealed.
   Try again later.

## Limits

A maximum request size of 32MB is imposed to prevent a denial
of service attack with arbitrarily large requests.
