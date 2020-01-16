---
layout: "docs"
page_title: "Entropy Augmentation - Configuration"
sidebar_title: "<code>Entropy Augmentation</code> <sup>ENT</sup>"
sidebar_current: "docs-configuration-entropy-augmentation"
description: |-
  Entropy augmentation enables Vault to sample entropy from external cryptographic modules.
---

# `Entropy Augmentation` Seal

  Entropy augmentation enables Vault to sample entropy from external cryptographic modules.
  Sourcing external entropy is done by configuring a supported [Seal](/docs/configuration/seal/index.html) type which
  include: [PKCS11 seal](/docs/configuration/seal/pkcs11.html), [AWS KMS](/docs/configuration/seal/awskms.html), and 
  [Vault Transit](/docs/configuration/seal/transit.html).
  Vault Enterprises's external entropy support is activated by the presence of an `entropy "seal"`
  block in Vault's configuration file.

## Requirements

A valid Vault Enterprise license is required for Entropy Augmentation

Additionally, the following software packages and enterprise modules are required for sourcing entropy 
via the [PKCS11 seal](/docs/configuration/seal/pkcs11.html):
- Governance and Policy module
- PKCS#11 compatible HSM integration library. Vault targets version 2.2 or
  higher of PKCS#11. Depending on any given HSM, some functions (such as key
  generation) may have to be performed manually.
- The [GNU libltdl library](https://www.gnu.org/software/libtool/manual/html_node/Using-libltdl.html)
  — ensure that it is installed for the correct architecture of your servers


## `entropy` Example

This example shows configuring entropy augmentation through a PKCS11 HSM seal from Vault's configuration
file:

```hcl
seal "pkcs11" {
    ...
}

entropy "seal" {
    mode = "augmentation"
}
```

For a more detailed tutorial, visit the [HSM Entropy Challenge](https://learn.hashicorp.com/vault/operations/hsm-entropy)
on HashiCorp's Learn website.

## `entropy augmentation` Parameters

These parameters apply to the `entropy` stanza in the Vault configuration file:

- `mode` `(string: <required>)`: The mode determines which Vault operations requiring
entropy will sample entropy from the external source. Currently, the only mode supported
is `augmentation` which sources entropy for [Critical Security Parameters (CSPs)](/docs/enterprise/entropy-augmentation/index.html#Critical-Security-Parameters-(CSPs)).
