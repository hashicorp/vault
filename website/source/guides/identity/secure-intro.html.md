---
layout: "guides"
page_title: "Secure Introduction of Vault Clients - Guides"
sidebar_current: "guides-identity-secure-intro"
description: |-
  This introductory guide walk through the mechanism of Vault clients to
  authenticate with Vault. There are two approaches at a high-level: platform
  integration, and trusted orchestrator.
---

# Secure Introduction of Vault Clients

A _secret_ is something that will elevate the risk if  exposed to unauthorized
entities and results in undesired consequences. For example:

- Unauthorized data access
- Identity spoofing
- Private data egress
- Regulatory fines

This means that only the ***trusted entities*** should have an access to your secrets. 
































## Next steps

Read the [_AppRole with Terraform and
Chef_](/guides/identity/approle-trusted-entities.html) guide to better
understand the role of trusted entities using Terraform and Chef as an example.

To learn more about response wrapping, go to the [Cubbyhole Response
Wrapping](/guides/secret-mgmt/cubbyhole.html) guide.
