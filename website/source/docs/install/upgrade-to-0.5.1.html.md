---
layout: "install"
page_title: "Upgrading to Vault 0.5.1"
sidebar_current: "docs-install-upgrade-to-0.5.1"
description: |-
  Learn how to upgrade to Vault 0.5.1
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
