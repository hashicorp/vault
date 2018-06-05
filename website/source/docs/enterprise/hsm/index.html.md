---
layout: "docs"
page_title: "HSM Integration - Vault Enterprise"
sidebar_current: "docs-vault-enterprise-hsm"
description: |-
  Vault Enterprise has HSM support, allowing for external master key storage and automatic unsealing.

---

# Vault Enterprise HSM Support

HSM support is a feature of [Vault
Enterprise](https://www.hashicorp.com/vault.html) that takes advantage of HSMs
to provide three pieces of special functionality:

 * Master Key Wrapping: Vault protects its master key by transiting it through
   the HSM for encryption rather than splitting into key shares
 * Automatic Unsealing: Vault stores its HSM-wrapped master key in storage,
   allowing for automatic unsealing
 * [Seal Wrapping](/docs/enterprise/sealwrap/index.html) to provide FIPS
   KeyStorage-conforming functionality for Critical Security Parameters

HSM support is available for devices that support PKCS#11 version 2.20+
interfaces and provide integration libraries, and is currently available for
linux/amd64 platforms only. It has successfully been tested against many
different vendor HSMs; HSMs that provide only subsets of the full PKCS#11
specification can usually be supported but it depends on available
cryptographic mechanisms.

Please note however that configuration details, flags, and supported features
within PKCS#11 vary depending on HSM model and configuration. Consult your
HSM's documentation for more details.

Some parts of Vault work differently when using an HSM. Please see the
[Behavioral Changes](/docs/enterprise/hsm/behavior.html) page for
important information on these differences.

The [Configuration](/docs/configuration/seal/pkcs11.html) page contains
configuration information.

Finally, the [Security](/docs/enterprise/hsm/security.html) page contains
information about deploying Vault's HSM support in a secure fashion.
