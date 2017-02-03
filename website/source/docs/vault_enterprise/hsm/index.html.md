---
title: "About Vault HSM Support"
---

# About Vault HSM Support

HSM support is a feature of [Vault
Enterprise](https://www.hashicorp.com/vault.html) that takes advantage of HSMs
to provide two pieces of special functionality:

 * Master Key Wrapping: Vault protects its master key by transiting it through
   the HSM for encryption rather than splitting into key shares
 * Automatic Unsealing: Vault stores its encrypted master key in storage,
   allowing for automatic unsealing

HSM support is currently limited to devices that support PKCS#11 interfaces and
provide integration libraries. This includes AWS'
[CloudHSM](https://aws.amazon.com/cloudhsm/) offering.

Some parts of Vault work differently when using an HSM. Please see the
[Behavioral Changes](/help/vault/hsm/behavior) page for important information
on these differences.

The [Configuration](/help/vault/hsm/configuration) page contains configuration
information.

Finally, the [Security](/help/vault/hsm/security) page contains information
about deploying Vault's HSM support in a secure fashion.
