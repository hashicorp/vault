---
layout: "guides"
page_title: "Upgrading to Vault 0.5.1 - Guides"
sidebar_current: "guides-upgrading-to-0.5.1"
description: |-
  This page contains the list of breaking changes for Vault 0.5.1. Please read
  it carefully.
---

# Overview

This page contains the list of breaking changes for Vault 0.5.1. Please read it
carefully.

## PKI Backend Disallows RSA Keys < 2048 Bits

The PKI backend now refuses to issue a certificate, sign a CSR, or save a role
that specifies an RSA key of less than 2048 bits. Smaller keys are considered
unsafe and are disallowed in the Internet PKI. Although this may require
updating your roles, we do not expect any other breaks from this; since its
inception the PKI backend has mandated SHA256 hashes for its signatures, and
software able to handle these certificates should be able to handle
certificates with >= 2048-bit RSA keys as well.

## PKI Backend Does Not Automatically Delete Expired Certificates

The PKI backend now does not automatically delete expired certificates,
including from the CRL. Doing so could lead to a situation where a time
mismatch between the Vault server and clients could result in a certificate
that would not be considered expired by a client being removed from the CRL.

Vault strives for determinism and putting the operator in control, so expunging
expired certificates has been moved to a new function at `pki/tidy`. You can
flexibly determine whether to tidy up from the revocation list, the general
certificate storage, or both. In addition, you can specify a safety buffer
(defaulting to 72 hours) to ensure that any time discrepancies between your
hosts is accounted for.

## Cert Authentication Backend Performs Client Checking During Renewals

The `cert` backend now performs a variant of channel binding at renewal time
for increased security. In order to not overly burden clients, a notion of
identity is used, as follows:

- At both login and renewal time, the validity of the presented client
  certificate is checked
- At login time, the key ID of both the client certificate and its issuing
  certificate are stored
- At renewal time, the key ID of both the client certificate and its issuing
  certificate must match those stored at login time

Matching on the key ID rather than the serial number allows tokens to be
renewed even if the CA or the client certificate used are rotated; so long as
the same key was used to generate the certificate (via a CSR) and sign the
certificate, renewal is allowed. As Vault encourages short-lived secrets,
including client certificates (for instance, those issued by the `pki`
backend), this is a useful approach compared to strict issuer/serial number
checking.

You can use the new `cert/config` endpoint to disable this behavior.
