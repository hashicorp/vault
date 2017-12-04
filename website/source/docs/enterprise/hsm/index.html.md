---
layout: "docs"
page_title: "Vault Enterprise HSM Support"
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
provide integration libraries. It has successfully been tested with AWS'
[CloudHSM](https://aws.amazon.com/cloudhsm/), as well as Thales, Utimaco, and
SafeNet/Gemalto devices.

Some parts of Vault work differently when using an HSM. Please see the
[Behavioral Changes](/docs/vault-enterprise/hsm/behavior.html) page for important information
on these differences.

The [Configuration](/docs/configuration/seal/pkcs11.html) page contains configuration
information.

Finally, the [Security](/docs/vault-enterprise/hsm/security.html) page contains information
about deploying Vault's HSM support in a secure fashion.
