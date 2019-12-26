---
layout: "guides"
page_title: "Identity and Access Management - Guides"
sidebar_title: "Identity and Access Management"
sidebar_current: "guides-identity"
description: |-
  Once a Vault instance has been installed, the next step is to configure auth
  backends, secret backends, and manage keys. Vault configuration guides addresses
  key concepts in configuring your Vault application.   
---

# Identity and Access Management

This guide walks you through Identity and Access Management topics.

- [Secure Introduction of Vault Clients](/guides/identity/secure-intro.html)
explains the mechanism of the client authentication against a Vault server.

- [Policies](/guides/identity/policies.html) are used to instrument
Role-Based Access Control (RBAC) by specifying access privileges. Authoring of
policies is probably the first step the Vault administrator performs. This guide
walks you through creating example policies for `admin` and `provisioner` users.

- [ACL Policy Path Templating](/guides/identity/policy-templating.html) guide
demonstrates templated policy paths to allow non-static paths.

- [AppRole Pull Authentication](/guides/identity/authentication.html) guide
that introduces the steps to generate tokens for machines or apps by enabling
AppRole auth method.

- [AppRole with Terraform and Chef](/guides/identity/approle-trusted-entities.html)
guide explains how AppRole auth method integrates with Terraform and Chef.
This guide is for anyone who wishes to reproduce the demo introduced during
the [Delivering Secret Zero: Vault AppRole with Terraform and
Chef](https://www.hashicorp.com/resources/delivering-secret-zero-vault-approle-terraform-chef)
webinar.

- [Token and Leases](/guides/identity/lease.html) guide helps you
understand how tokens and leases work in Vault. The understanding of the
lease hierarchy and expiration mechanism helps you plan for break glass
procedures and more.

- [Identity - Entities & Groups](/guides/identity/identity.html) guide
demonstrates the usage of _Entities_ and _Groups_ to manage Vault clients'
identity.

## Vault Enterprise

- [Sentinel Policies](/guides/identity/sentinel.html) guide
walks through the creation and usage of _Role Governing Policies_ (RGPs) and
_Endpoint Governing Policies_ (EGPs) in Vault.

- [Control Groups](/guides/identity/control-groups.html) can be used to enforce
additional authorization factors before the request can be completed. This
guide walks through the implementation of a Control Group.
