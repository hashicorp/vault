---
layout: "guides"
page_title: "Identity and Access Management - Guides"
sidebar_current: "guides-identity"
description: |-
  Once a Vault instance has been installed, the next step is to configure auth
  backends, secret backends, and manage keys. Vault configuration guides addresses
  key concepts in configuring your Vault application.   
---

# Identity and Access Management

This guide walks you through Identity and Access Management topics.

- [Policies](/guides/identity/policies.html) are used to instrument
Role-Based Access Control (RBAC) by specifying access privileges. Authoring of
policies is probably the first step the Vault administrator performs. This guide
walks you through creating example policies for `admin` and `provisioner` users.
- [AppRole Pull Authentication](/guides/identity/authentication.html) guide
that introduces the steps to generate tokens for machines or apps by enabling
AppRole auth method.
- [Token and Leases](/guides/identity/lease.html) guide helps you
understand how tokens and leases work in Vault. The understanding of the
lease hierarchy and expiration mechanism helps you plan for break glass
procedures and more.
