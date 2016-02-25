---
layout: "install"
page_title: "Upgrading to Vault 0.6"
sidebar_current: "docs-install-upgrade-to-0.6"
description: |-
  Learn how to upgrade to Vault 0.6
---

# Overview

This page contains the list of breaking changes for Vault 0.6. Please read it
carefully.

Please note that this includes the full list of breaking changes _since Vault 0.5_. Some of these changes were introduced in later releases in the Vault 0.5.x series.

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
