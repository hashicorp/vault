---
layout: "guides"
page_title: "Vault Operations - Guides"
sidebar_current: "guides-operations"
description: |-
  Vault architecture guide covers Vault infrastructure discussions including
  installation.   
---

# Vault Operations

Vault Operations guides address Vault infrastructure discussions.  These
guides are designed to help the operations team to plan and install a Vault
cluster that meets your organization's needs.

- [Vault Reference Architecture](/guides/operations/reference-architecture.html)
guide provides guidance in the best practices of _Vault Enterprise_ implementations
through use of a reference architecture. This example is to convey a general
architecture, which is likely to be adapted to accommodate the specific needs of
each implementation.

- [Vault HA with Consul](/guides/operations/vault-ha-consul.html) guide
walks you through a simple Vault HA cluster implementation which is backed by
[HashiCorp Consul](https://www.consul.io/intro/index.html).

- [Production Hardening](/guides/operations/production.html) guide provides
guidance on best practices for a production hardened deployment of Vault.
The recommendations are based on the [security model](/docs/internals/security.html)
and focus on defense in depth.

- **[Enterprise Only]** [Replication Setup & Guidance](/guides/operations/replication.html)
walks you through the commands to activate the Vault servers in replication mode.
Please note that [Vault Replication](/docs/vault-enterprise/replication/index.html)
is a Vault Enterprise feature.

- **[Enterprise Only]** [Mount Filter](/guides/operations/mount-filter.html)
guide demonstrates how to selectively filter out secret engines from being
replicated across clusters. This feature can help organizations to comply with
***General Data Protection Regulation (GDPR)***.

- **[Enterprise Only]** [Vault Auto-unseal using AWS Key Management Service (KMS)](/guides/operations/autounseal-aws-kms.html) guide demonstrates an example of
how to use Terraform to provision an instance that utilizes an encryption key
from AWS Key Management Service (KMS).

- **[Enterprise Only]** [Seal Wrap / FIPS 140-2](/guides/operations/seal-wrap.html)
guide demonstrates how Vault's seal wrap feature works to encrypt your secrets
leveraging FIPS 140-2 certified HSM.

- [Root Token Generation](/guides/operations/generate-root.html) guide
demonstrates the workflow of regenerating root tokens. It is considered to be a
best practice not to persist the initial **root** token. If a root token needs
to be regenerated, this guide helps you walk through the task.

- [Rekeying & Rotating](/guides/operations/rekeying-and-rotating.html) guide
provides a high-level overview of Shamir's Secret Sharing Algorithm, and how to
perform _rekey_ and _rotate_ operations in Vault.

- [Building Plugin Backends](/guides/operations/plugin-backends.html) guide
provides steps to build, register, and mount non-database external plugin
backends.
