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
to provide two pieces of special functionality:

 * Master Key Wrapping: Vault protects its master key by transiting it through
   the HSM for encryption rather than splitting into key shares
 * Automatic Unsealing: Vault stores its encrypted master key in storage,
   allowing for automatic unsealing

HSM support is currently limited to devices that support PKCS#11 interfaces and
provide integration libraries. It has successfully been tested with the following
HSM platforms/vendors:

 * AWS [CloudHSM](https://aws.amazon.com/cloudhsm/)
 * Thales [nShield](https://www.thalesesecurity.com/products/general-purpose-hsms)
 * Utimaco [HSM](https://hsm.utimaco.com/)
 * SafeNet/Gemalto [Luna](https://safenet.gemalto.com/data-encryption/hardware-security-modules-hsms/safenet-network-hsm/)
 * Unbound [HSM](https://www.unboundtech.com/)
 * OpenDNSSEC [SoftHSM](https://www.opendnssec.org/softhsm/) 

Please note however that configuration details, flags, and supported features within PKCS#11 vary depending on HSM model and configuration. Consult your HSM's documentation for more details.

Some parts of Vault work differently when using an HSM. Please see the
[Behavioral Changes](/docs/vault-enterprise/hsm/behavior.html) page for important information
on these differences.

The [Configuration](/docs/configuration/seal/pkcs11.html) page contains configuration
information.

Finally, the [Security](/docs/vault-enterprise/hsm/security.html) page contains information
about deploying Vault's HSM support in a secure fashion.
