---
layout: "intro"
page_title: "Vault vs. Amazon Key Management Service"
sidebar_current: "vs-other-kms"
description: |-
  Comparison between Vault and Amazon Key Management Service.
---

# Vault vs. Amazon KMS

Amazon Key Management Service (KMS) is a service provided in the AWS ecosystem
for encryption key management. It is backed by Hardware Security Modules (HSM)
for physical security.

Vault and KMS differ in the scope of problems they are trying to solve.  KMS is
focused on securely storing encryption keys and supporting cryptographic
operations (encrypt and decrypt) using those keys. It supports access controls
and auditing as well.

In contrast, Vault provides a comprehensive secret management solution.  The
[`transit` backend](/docs/secrets/transit/index.html) provides similar
capabilities as KMS, allowing for encryption keys to be stored and
cryptographic operations to be performed. However, Vault goes much further than
just key management.

The flexible secret backends allow Vault to handle any type of secret data,
including database credentials, API keys, PKI keys, and encryption keys.  Vault
also supports dynamic secrets, generating credentials on-demand for
fine-grained security controls, auditing, and non-repudiation.

Lastly Vault forces a mandatory lease contract with clients. All secrets read
from Vault have an associated lease which enables operations to audit key
usage, perform key rolling, and ensure automatic revocation. Vault provides
multiple revocation mechanisms to give operators a clear "break glass"
procedure after a potential compromise.

Vault is an open source tool that can be deployed to any environment, and does
not require any special hardware. This makes it well suited for cloud
environments where HSMs are not available or are cost prohibitive.

