---
layout: "docs"
page_title: "Response Wrapping"
sidebar_current: "docs-concepts-response-wrapping"
description: |-
  Wrapping responses in cubbyholes for secure distribution.
---

# Response Wrapping

_Note_: Some of this information relies on features of response-wrapping tokens
introduced in Vault 0.8 and may not be available in earlier releases.

## Overview

In many Vault deployments, clients can access Vault directly and consume
returned secrets. In other situations, it may make sense to or be desired to
separate privileges such that one trusted entity is responsible for interacting
with most of the Vault API and passing secrets to the end consumer.

However, the more relays a secret travels through, the more possibilities for
accidental disclosure, especially if the secret is being transmitted in
plaintext. For instance, you may wish to get a TLS private key to a machine
that has been cold-booted, but since you do not want to store a decryption key
in persistent storage, you cannot encrypt this key in transit.

To help address this problem, Vault includes a feature called _response
wrapping_. When requested, Vault can take the response it would have sent to an
HTTP client and instead insert it into the
[`cubbyhole`](/docs/secrets/cubbyhole/index.html) of a single-use token,
returning that single-use token instead. Logically speaking, the response is
wrapped by the token, and retrieving it requires an unwrap operation against
this token.

This provides a powerful mechanism for information sharing in many
environments. In the types of scenarios, described above, often the best
practical option is to provide _cover_ for the secret information, be able to
_detect malfeasance_ (interception, tampering), and limit _lifetime_ of the
secret's exposure. Response wrapping performs all three of these duties:

 * It provides _cover_ by ensuring that the value being transmitted across the
   wire is not the actual secret but a reference to such a secret, namely the
   response-wrapping token. Information stored in logs or captured along the
   way do not directly see the sensitive information.
 * It provides _malfeasance detection_ by ensuring that only a single party can
   ever unwrap the token and see what's inside. A client receiving a token that
   cannot be unwrapped can trigger an immediate security incident. In addition,
   a client can inspect a given token before unwrapping to ensure that its
   origin is from the expected location in Vault.
 * It _limits the lifetime_ of secret exposure because the response-wrapping
   token has a lifetime that is separate from the wrapped secret (and often can
   be much shorter), so if a client fails to come up and unwrap the token, the
   token can expire very quickly.

## Response-Wrapping Tokens

When a response is wrapped, the normal API response from Vault does not contain
the original secret, but rather contains a set of information related to the
response-wrapping token:

 * TTL: The TTL of the response-wrapping token itself
 * Token: The actual token value
 * Creation Time: The time that the response-wrapping token was created
 * Creation Path: The API path that was called in the original request
 * Wrapped Accessor: If the wrapped response is an authentication response
   containing a Vault token, this is the value of the wrapped token's accessor.
   This is useful for orchestration systems (such as Nomad) to able to control
   the lifetime of secrets based on their knowledge of the lifetime of jobs,
   without having to actually unwrap the response-wrapping token or gain
   knowledge of the token ID inside.

Vault currently does not provide signed response-wrapping tokens, as it
provides little extra protection. If you are being pointed to the correct Vault
server, token validation is performed by interacting with the server itself; a
signed token does not remove the need to validate the token with the server,
since the token is not carrying data but merely an access mechanism and the
server will not release data without validating it. If you are being attacked
and pointed to the wrong Vault server, the same attacker could trivially give
you the wrong signing public key that corresponds to the wrong Vault server.
You could cache a previously valid key, but could also cache a previously valid
address (and in most cases the Vault address will not change or will be set via
a service discovery mechanism). As such, we rely on the fact that the token
itself is not carrying authoritative data and do not sign it.

## Response-Wrapping Token Operations

Via the `sys/wrapping` path, several operations can be run against wrapping
tokens:

 * Lookup (`sys/wrapping/lookup`): This allows fetching the response-wrapping
   token's creation time, creation path, and TTL. This path is unauthenticated
   and available to response-wrapping tokens themselves. In other words, a
   response-wrapping token holder wishing to perform validation is always
   allowed to look up the properties of the token.
 * Unwrap (`sys/wrapping/unwrap`): Unwrap the token, returning the response
   inside. The response that is returned will be the original wire-format
   response; it can be used directly with API clients.
 * Rewrap (`sys/wrapping/rewrap`): Allows migrating the wrapped data to a new
   response-wrapping token. This can be useful for long-lived secrets. For
   example, an organization may wish (or be required in a compliance scenario)
   to have the `pki` backend's root CA key be returned in a long-lived
   response-wrapping token to ensure that nobody has seen the key (easily
   verified by performing lookups on the response-wrapping token) but available
   for signing CRLs in case they ever accidentally change or lose the `pki`
   mount.  Often, compliance schemes require periodic rotation of secrets, so
   this helps achieve that compliance goal without actually exposing what's
   inside.
 * Wrap (`sys/wrapping/wrap`): A helper endpoint that echoes back the data sent
   to it in a response-wrapping token. Note that blocking access to this
   endpoint does not remove the ability for arbitrary data to be wrapped, as it
   can be done elsewhere in Vault.

## Response-Wrapping Token Creation

Response wrapping is per-request and is triggered by providing to Vault the
desired TTL for a response-wrapping token for that request. This is set by the
client using the `X-Vault-Wrap-TTL` header and can be either an integer number
of seconds or a string duration of seconds (`15s`), minutes (`20m`), or hours
(`25h`). When using the Vault CLI, you can set this via the `-wrap-ttl`
parameter. When using the Go API, wrapping is triggered by [setting a helper
function](https://godoc.org/github.com/hashicorp/vault/api#Client.SetWrappingLookupFunc)
that tells the API the conditions under which to request wrapping, by mapping
an operation and path to a desired TTL.

If a client requests wrapping:

1. The original HTTP response is serialized
2. A new single-use token is generated with the TTL supplied by the client
3. Internally, the original serialized response is stored in the single-use
   token's cubbyhole
4. A new response is generated, with the token ID, TTL, and path stored in the
   new response's wrap information object
5. The new response is returned to the caller

Note that policies can control minimum/maximum wrapping TTLs; see the [policies
concepts page](https://www.vaultproject.io/docs/concepts/policies.html) for
more information.

## Response-Wrapping Token Validation

Proper validation of response-wrapping tokens is essential to ensure that any
malfeasance is detected. It's also pretty straightforward.

Validation is best performed by the following steps:

1. If a client has been expecting delivery of a response-wrapping token and
   none arrives, this may be due to an attacker intercepting the token and then
   preventing it from traveling further. This should cause an alert to trigger
   an immediate investigation.
2. Perform a lookup on the response-wrapping token. This immediately tells you
   if the token has already been unwrapped or is expired (or otherwise
   revoked). If the lookup indicates that a token is invalid, it does not
   necessarily mean that the data was intercepted (for instance, perhaps the
   client took a long time to start up and the TTL expired) but should trigger
   an alert for immediate investigation, likely with the assistance of Vault's
   audit logs to see if the token really was unwrapped.
3. With the token information in hand, validate that the creation path matches
   expectations. If you expect to find a TLS key/certificate inside, chances
   are the path should be something like `pki/issue/...`. If the path is not
   what you expect, it is possible that the data contained inside was read and
   then put into a new response-wrapping token. (This is especially likely if
   the path starts with `cubbyhole` or `sys/wrapping/wrap`.) Particular care
   should be taken with `kv` mounts: exact matches on the path are best
   there.  For example, if you expect a secret to come from `secret/foo` and
   the interceptor provides a token with `secret/bar` as the path, simply
   checking for a prefix of `secret/` is not enough.
4. After prefix validation, unwrap the token. If the unwrap fails, the response
   is similar to if the initial lookup fails: trigger an alert for immediate
   investigation.

Following those steps provides very strong assurance that the data contained
within the response-wrapping token has never been seen by anyone other than the
intended client and that any interception or tampering has resulted in a
security alert.
