---
layout: "guides"
page_title: "HSM Autounseal and Seal Wrap - Guides"
sidebar_current: "guides-operations-hsm-seal-wrap"
description: |-
  In this guide,
---


# HSM Auto-unseal and Seal Wrap (FIPS 140-2)

~> **Enterprise Only:** Seal Wrap feature is developed for FIPS 140-2 compliance;
therefore, it requires FIPS 140-2 compliant HSM.  Vault's HSM auto-unseal feature
is a part of _Vault Enterprise Premium_.

[Vault Enterprise](https://www.hashicorp.com/vault.html) integrates with [a
number of HSM platforms](/docs/enterprise/hsm/index.html) to opt-in automatic
unsealing. HSM integration provides two pieces of special functionality:

- **Master Key Wrapping**: Vault protects its master key by transiting it through
the HSM for encryption rather than splitting into key shares
- **Automatic Unsealing**: Vault stores its encrypted master key in storage,
allowing for automatic unsealing

![Unseal with HSM](/assets/images/vault-hsm-autounseal.png)

Vault pulls its encrypted master key from its storage (e.g. Consul) and transit
it through the HSM for decryption via  **PKCS \#11 API**. Once the master key is
decrypted, Vault uses the master key to decrypt the encryption key so that Vault
operations can be performed.


## Reference Material

-

## Estimated Time to Complete

10 minutes

## Personas

The steps described in this guide are typically performed by **operations**
persona.

## Challenge



## Solution

Vault encrypts secrets using 256-bit AES in GCM mode with a randomly generated
nonce prior to writing them to its persistent storage (e.g. Consul). By
integrating Vault with HSM, Vault wraps your secrets with an extra layer of
encryption. This additional data protection is useful in some compliance and
regulatory environments including FIPS 140-2 environment.




## Prerequisites

This guide assumes the following:   

- Access to **Vault Enterprise 0.9.0 or later** which supports AWS KMS as an unseal mechanism
- A URL to download Vault Enterprise from (an Amazon S3 bucket will suffice)
- AWS account for provisioning cloud resources
- [Terraform installed](https://www.terraform.io/intro/getting-started/install.html)
and basic understanding of its usage




## Steps

To configure Vault Enterprise integration with AWS CloudHSM

1. Follow the steps in [Set Up Prerequisites](#vault-autounseal-prerequisites.html) to prepare your environment

1. Follow the steps in [Configure the Vault Server](#vault-autounseal-config.html) to integrate Vault with your AWS CloudHSM cluster.

1. [Seal Wrap Enabled by AWS CloudHSM](#vault-autounseal-sealwrap.html) describes the feature to protect your data with additional layer of encryption.


You are going to perform the following steps:

1. [Provision the Cloud Resources](#step-1-provision-the-cloud-resources)
1. [Test the Auto-unseal Feature](#step-2-test-the-auto-unseal-feature)
1. [Clean Up](#step-3-clean-up)


-> There is just a single master key that's created upon init, is encrypted by the
HSM using PKCS \#11, then placed in the storage backend to be used later on when
unsealing is needed. When it goes to unseal, Vault grabs the HSM encrypted
master key from the storage backend, round trips it through the HSM to decrypt,
and unseals itself. The pin and slot are used for Vault to authenticate with the
HSM I believe.



### Step 1: Provision the Cloud Resources



## Next steps

Once you have a Vault environment setup, the next step is to write policies.
Read [Policies](/guides/identity/policies.html) to learn how to write policies
to govern the behavior of clients.
