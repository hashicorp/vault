---
layout: "guides"
page_title: "Upgrading to Vault 0.6.0 - Guides"
sidebar_current: "guides-upgrading-to-0.6.0"
description: |-
  This page contains the list of breaking changes for Vault 0.6. Please read it
  carefully.
---

# Overview

This page contains the list of breaking changes for Vault 0.6 compared to the
previous release. Please read it carefully.

## PKI Backend Does Not Issue Leases for CA Certificates

When a token expires, it revokes all leases associated with it. This means that
long-lived CA certs need correspondingly long-lived tokens, something that is
easy to forget, resulting in an unintended revocation of the CA certificate
when the token expires. To prevent this, root and intermediate CA certs no
longer have associated leases. To revoke these certificates, use the
`pki/revoke` endpoint.

CA certificates that have already been issued and acquired leases will report
to the lease manager that revocation was successful, but will not actually be
revoked and placed onto the CRL.

## The `auth/token/revoke-prefix` Endpoint Has Been Removed

As part of addressing a minor security issue, this endpoint has been removed in
favor of using `sys/revoke-prefix` for prefix-based revocation of both tokens
and secrets leases.

## Go API Uses `json.Number` For Decoding

When using the Go API, it now calls `UseNumber()` on the decoder object. As a
result, rather than always decode as a `float64`, numbers are returned as a
`json.Number`, where they can be converted, with proper error checking, to
`int64`, `float64`, or simply used as a `string` value. This fixes some display
errors where numbers were being decoded as `float64` and printed in scientific
notation.

## List Operations Return `404` On No Keys Found

Previously, list operations on an endpoint with no keys found would return an
empty response object. Now, a `404` will be returned instead.

## Consul TTL Checks Automatically Registered

If using the Consul HA storage backend, Vault will now automatically register
itself as the `vault` service and perform its own health checks/lifecycle
status management. This behavior can be adjusted or turned off in Vault's
configuration; see the
[documentation](/docs/configuration/index.html#check_timeout)
for details.
