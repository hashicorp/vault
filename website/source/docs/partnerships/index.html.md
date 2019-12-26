---
layout: "docs"
page_title: "Vault Integration Program"
sidebar_current: "docs-partnerships"
description: |-
  Guide to partnership integrations and creating plugins for Vault.
---

# Vault Integration Program

The Vault Integration Program (VIP) enables vendors to build integrations with HashiCorp Vault that are officially tested and approved by HashiCorp. The program is intended to be largely self-service, with links to code samples, documentation and clearly defined integration steps.

## Types of Vault Integrations

By leveraging Vault's plugin system, vendors are able to build extensible secrets, authentication, and audit plugins to extend Vault's functionality. These integrations can be done with the OSS (open-source) version of Vault. Hardware Security Module (HSM) integrations need to be tested against Vault Enterprise since the HSM functionality is only supported in the Vault Enterprise version.

**Authentication Methods**: Auth methods are the components in Vault that perform authentication and are responsible for assigning identity and a set of policies to a user.

**Vault Secrets Engine**: Secrets engines are components which store, generate, or encrypt data. Secrets engines are incredibly flexible, so it is easiest to think about them in terms of their function. Secrets engines are provided some set of data, they take some action on that data, and they return a result.

**Audit Devices**: Audit devices are the components in Vault that keep a detailed log of all requests and response to Vault. Because every operation with Vault is an API request/response, the audit log contains every authenticated interaction with Vault, including errors.

**Hardware Security Module (HSM)**: HSM support is a feature of Vault Enterprise that takes advantage of HSMs to provide Master Key Wrapping, Automatic Unsealing and Seal Wrapping via the PKCS#11 protocol ver. 2.2+.

**Cloud / Third Party Autounseal Integration**: Non-PKCS#11 integrations with secure external data stores (e.g.: AWS KMS, Azure Key Vault) to provide Autounsealing and Seal-Wrapping.

**Storage Backend**: A storage backend is a durable storage location where Vault stores its information.

## Development Process

The Vault integration development process is described into the steps below. By following these steps, Vault integrations can be developed alongside HashiCorp to ensure new integrations are reviewed, certified and released as quickly as possible.

1.  Engage: Initial contact between vendor and HashiCorp
2.  Enable: Documentation, code samples and best practices for developing the integration
3.  Develop and Test: Integration development and testing by vendor
4.  Review/Certification: HashiCorp code review and certification of integration
5.  Release: Vault integration released
6.  Support: Ongoing maintenance and support of the integration by the vendor.

### 1. Engage

Please begin by completing Vault Integration Program webform to tell us about your company and the Vault integration you’re interested in.

### 2. Enable

Here are links to resources, documentation, examples and best practices to guide you through the Vault integration development and testing process:

**General Vault Plugin Development:**

* [Plugins documentation](https://www.vaultproject.io/docs/internals/plugins.html)
* [Guide to building Vault plugin backends](https://www.vaultproject.io/guides/operations/plugin-backends.html)
* [Vault's source code](https://github.com/hashicorp/vault)

**Secrets Engines**

* [Secret engine documentation](https://www.vaultproject.io/docs/secrets/index.html)
* There is currently no empty sample secrets plugin; however, the [AliCloud Secrets Plugin](https://github.com/hashicorp/vault-plugin-secrets-alicloud) was written recently and is fairly simple

**Authentication Methods**

* [Auth Methods documentation](https://www.vaultproject.io/docs/auth/index.html)
* [Example of how to build, install, and maintain auth method plugins plugin](https://www.hashicorp.com/blog/building-a-vault-secure-plugin)
* [Sample plugin code](https://github.com/hashicorp/vault-auth-plugin-example)

**Audit Devices**

[Audit devices documentation](https://www.vaultproject.io/docs/audit/index.html)

**HSM Integration**

* [HSM documentation](https://www.vaultproject.io/docs/enterprise/hsm/index.html)
* [Configuration information](https://www.vaultproject.io/docs/configuration/seal/pkcs11.html)

**Storage Backends**

[Storage configuration documentation](https://www.vaultproject.io/docs/configuration/storage/index.html)

**Community Forum**

[Vault developer community forum](https://groups.google.com/forum/#!forum/vault-tool)

### 3. Develop and Test

The only knowledge necessary to write a plugin is basic command-line skills and knowledge of the [Go programming language](http://www.golang.org). Use the plugin interface to develop your integration. All integrations should contain unit and acceptance testing.

### 4. Review

HashiCorp will review and certify your Vault integration. Please send the Vault logs and other relevant logs for verification at: [vault-integration-dev@hashicorp.com](mailto:vault-integration-dev@hashicorp.com). For Auth, Secret and Storage plugins, submit a GitHub pull request (PR) against the [Vault project](https://github.com/hashicorp/vault). Where applicable, the vendor will need to provide HashiCorp with a test account.

### 5. Release

At this stage, the Vault integration is fully developed, documented, tested and certified. Once released, HashiCorp will officially list the Vault integration.

### 6. Support

Many vendors view the release step to be the end of the journey, while at HashiCorp we view it to be the start. Getting the Vault integration built is just the first step in enabling users. Once this is done, on-going effort is required to maintain the integration and address any issues in a timely manner.
The expectation for vendors is to respond to all critical issues within 48 hours and all other issues within 5 business days. HashiCorp Vault has an extremely wide community of users and we encourage everyone to report issues however small, as well as help resolve them when possible.

## Checklist

Below is a checklist of steps that should be followed during the Vault integration development process. This reiterates the steps described above.

* Complete the [Vault Integration webform](https://docs.google.com/forms/d/e/1FAIpQLSfQL1uj-mL59bd2EyCPI31LT9uvVT-xKyoHAb5FKIwWwwJ1qQ/viewform)
* Develop and test your Vault integration following examples, documentation and best practices
* When the integration is completed and ready for HashiCorp review, send the Vault and other relevant logs to us for review and certification at: [vault-integration-dev@hashicorp.com](mailto:vault-integration-dev@hashicorp.com)
* Once released, plan to support the integration with additional functionality and responding to customer issues

## Contact Us

For any questions or feedback, please contact us at: [vault-integration-dev@hashicorp.com](mailto:vault-integration-dev@hashicorp.com)
