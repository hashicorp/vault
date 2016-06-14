---
layout: "docs"
page_title: "Response Wrapping"
sidebar_current: "docs-concepts-response-wrapping"
description: |-
  Wrapping responses in cubbyholes for secure distribution.
---

# Response Wrapping

In many Vault deployments, clients can access Vault directly and consume
returned secrets. In other situations, it may make sense to or be desired to
separate privileges such that one trusted entity is responsible for interacting
with most of the Vault API and passing secrets to the end consumer.

However, the more relays a secret travels through, the more possibilities for
accidental disclosure, especially if the secret is being transmitted in
plaintext.

In Vault 0.3 the
[`cubbyhole`](https://www.vaultproject.io/docs/secrets/cubbyhole/index.html)
backend was introduced, providing storage scoped to a single token. The
[Cubbyhole Principles blog
post](https://www.hashicorp.com/blog/vault-cubbyhole-principles.html) described
how this, along with the limited-use and time-to-live features of Vault tokens,
could be used to securely authenticate a Vault client in such a way that the
final Vault token was only readable by the end consumer, and malfeasance could
be detected. The major downside to this operation was the need to write
programs to perform this wrapping (and by extension, those programs need to be
trusted).

Starting in 0.6, this concept is taken to its logical conclusion: almost every
response that Vault generates can be automatically wrapped inside a single-use,
limited-time-to-live token's cubbyhole. Details can be found in the
[`cubbyhole` backend
documentation](https://www.vaultproject.io/docs/secrets/cubbyhole/index.html).

This capability should be carefully considered when planning your security
architecture. For instance, many Vault deployments use the
[`pki`](https://www.vaultproject.io/docs/secrets/pki/index.html) backend to
generate TLS certificates and private keys for services. If you do not wish
these services to have access to the generation API, a trusted third party
could generate the certificates and private keys and pass the resulting
wrapping tokens directly to the services in need. A simple API call will return
the original PKI information; if the call fails, a security alert can be
raised.

To look at the above example another way, response wrapping also frees end
services from needing to generate a CSR and pass it to Vault through the
trusted third party simply to ensure that the private key corresponding to the
eventual certificate remains private. The end service can be assured that only
it will see the generated private key and that any malfeasance is detected.
This can significantly reduce the complexity of any relaying third party.

One final note: if the wrapped response is an authentication response
containing a Vault token, the token's accessor will be made available in the
returned wrap information. This allows privileged callers to generate tokens
for clients and revoke these tokens (and their created leases) at an
appropriate time, while never being exposed to the actual generated token IDs.
