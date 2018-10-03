---
layout: "guides"
page_title: "Partnerships - Vault Integration Program"
sidebar_current: "guides-partnerships"
description: |-
  Guide to partnership integrations and creating plugins for Vault.
---

# Vault Integration Program

<p>
 The Vault Integration Program (VIP) enables vendors to build integrations with HashiCorp Vault that are officially tested and approved by HashiCorp. The program is intended to be largely self-service, with links to code samples, documentation and clearly defined integration steps.
</p>

## Types of Vault Integrations

<p>
By leveraging Vault's plugin system, vendors are able to build extensible secrets, authentication, and audit plugins to extend Vault's functionality. These integrations can be done with the OSS (open-source) version of Vault. Hardware Security Module (HSM) integrations need to be tested against Vault Enterprise since the HSM functionality is only supported in the Vault Enterprise version.
</p>

<p>
<strong>Authentication Methods</strong>: Auth methods are the components in Vault that perform authentication and are responsible for assigning identity and a set of policies to a user.
</p>

<p>
<strong>Vault Secrets Engine</strong>: Secrets engines are components which store, generate, or encrypt data. Secrets engines are incredibly flexible, so it is easiest to think about them in terms of their function. Secrets engines are provided some set of data, they take some action on that data, and they return a result.
</p>

<p>
<strong>Audit Devices</strong>: Audit devices are the components in Vault that keep a detailed log of all requests and response to Vault. Because every operation with Vault is an API request/response, the audit log contains every authenticated interaction with Vault, including errors.  (no plugin interface - built into Vault Core. Leave it there - no reqs yet but expect some soon)
</p>

<p>
<strong>
Hardware Security Module (HSM)</strong>: HSM support is a feature of Vault Enterprise that takes advantage of HSMs to provide Master Key Wrapping, Automatic Unsealing and Seal Wrapping via the PKCS#11 protocol ver. 2.2+.
</p>

<p>
<strong>
Cloud / Third Party Autounseal Integration</strong>: Non-PKCS#11 integrations with secure external data stores (e.g.: AWS KMS, Azure Key Vault) to provide Autounsealing and Seal-Wrapping.
</p>

<p>
<strong>
Storage Backend</strong>:  A storage backend is a durable storage location where Vault stores its information.
</p>

<h2>Development Process</h2>

<p>The Vault integration development process is described into the steps below. By following these steps, Vault integrations can be developed alongside HashiCorp to ensure new integrations are reviewed, certified and released as quickly as possible.</p>

<ol type="1">
<li>Engage: Initial contact between vendor and HashiCorp</li>
<li>Enable: Documentation, code samples and best practices for developing the integration</li>
<li>Develop and Test: Integration development and testing by vendor</li>
<li>Review/Certification: HashiCorp code review and certification of integration</li>
<li>Release: Vault integration released</li>
<li>Support: Ongoing maintenance and support of the integration by the vendor.</li>
</ol>

### 1. Engage</h3>
<p>
Please begin by completing Vault Integration Program webform to tell us about your company and the Vault integration you’re interested in.
</p>

### 2. Enable</h3>
<p>
Here are links to resources, documentation, examples and best practices to guide you through the Vault integration development and testing process:
</p>

<p><strong>General Vault Plugin Development:</strong></p>
<ul>
<li><a href="https://www.vaultproject.io/docs/internals/plugins.html">Plugins documentation</a></li>
<li><a href="https://www.vaultproject.io/guides/operations/plugin-backends.html">Guide to building Vault plugin backends</a></li>
<li><a href="https://github.com/hashicorp/vault">Vault's source code</a></li>
</ul>

<p><strong>Secrets Engines</p></strong>
<ul>
<li><a href="https://www.vaultproject.io/docs/secrets/index.html">Secret engine documentation</a></li>
<li><a href="https://github.com/hashicorp/vault-auth-plugin-example">Sample plugin code</a></li>
</ul>

<p><strong>Authentication Methods</strong></p>
<ul>
<li><a href="https://www.vaultproject.io/docs/auth/index.html">Auth Methods documentation</a></li>
<li><a href="https://www.hashicorp.com/blog/building-a-vault-secure-plugin">Example of how to build, install, and maintain auth method plugins plugin</a></li> 
<li><a href="https://github.com/hashicorp/vault-auth-plugin-example">Sample plugin code</a></li>
</ul>

<p><strong>Audit Devices</p></strong>
<p><a href="https://www.vaultproject.io/docs/audit/index.html">Audit devices documentation</a></p>

<p><strong>HSM Integration</strong></p>
<ul>
<li><a href="https://www.vaultproject.io/docs/enterprise/hsm/index.html">HSM documentation</a></li>
<li><a href="https://www.vaultproject.io/docs/configuration/seal/pkcs11.html">Configuration information</a></li>
</ul>

<p><strong>Storage Backends</strong></p>
<p><a href="https://www.vaultproject.io/docs/configuration/storage/index.html">Storage configuration documentation</a></p>

<p><strong>Community Forum</strong></p>
<p><a href="https://groups.google.com/forum/#!forum/vault-tool">Vault developer community forum</a></p> 

### 3. Develop and Test </h3>
<p>
The only knowledge necessary to write a plugin is basic command-line skills and knowledge of the <a href="http://www.golang.org">Go programming language</a>.  Use the plugin interface to develop your integration. All integrations should contain unit and acceptance testing.  
</p>

### 4. Review 
<p>
HashiCorp will review and certify your Vault integration. Please send the Vault logs and other relevant logs for verification at: <a href="mailto:vault-integration-dev@hashicorp.com">vault-integration-dev@hashicorp.com</a>. For Auth, Secret and Storage plugins, submit a GitHub pull request (PR) against the Vault project (https://github.com/hashicorp/vault). Where applicable, the vendor will need to provide HashiCorp with a test account.
</p>

### 5. Release 
<p>
At this stage, the Vault integration is fully developed, documented, tested and certified. Once released, HashiCorp will officially list the Vault integration.
</p>

### 6. Support
<p>
Many vendors view the release step to be the end of the journey, while at HashiCorp we view it to be the start. Getting the Vault integration built is just the first step in enabling users. Once this is done, on-going effort is required to maintain the integration and address any issues in a timely manner.
The expectation for vendors is to respond to all critical issues within 48 hours and all other issues within 5 business days. HashiCorp Vault has an extremely wide community of users and we encourage everyone to report issues however small, as well as help resolve them when possible.
</p>

## Checklist
<p>Below is a checklist of steps that should be followed during the Vault integration development process. This reiterates the steps described above.</p>
<ul>
<p><li>Complete the <a href="https://docs.google.com/forms/d/e/1FAIpQLSfQL1uj-mL59bd2EyCPI31LT9uvVT-xKyoHAb5FKIwWwwJ1qQ/viewform">Vault Integration webform</a></li></p>
<p><li>Develop and test your Vault integration following examples, documentation and best practices</li></p>
<p><li>When the integration is completed and ready for HashiCorp review, send the Vault and other relevant logs to us for review and certification at: <a href="mailto:vault-integration-dev@hashicorp.com">vault-integration-dev@hashicorp.com</a></li></p>
<p><li>Once released, plan to support the integration with additional functionality and responding to customer issues </li></p>
</ul>

## Contact Us
<p>For any questions or feedback, please contact us at: <a href="mailto:vault-integration-dev@hashicorp.com">vault-integration-dev@hashicorp.com</a></p>
