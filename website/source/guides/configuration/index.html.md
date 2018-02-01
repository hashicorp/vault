---
layout: "guides"
page_title: "Vault Configuration - Guides"
sidebar_current: "guides-configuration"
description: |-
  Once a Vault instance has been installed, the next step is to configure auth
  backends, secret backends, and manage keys. Vault configuration guides addresses
  key concepts in configuring your Vault application.   
---

# Vault Configuration

This guide walks you through Vault configuration topics.

- [Policies](/guides/configuration/policies.html) are used to instrument
Role-Based Access Control (RBAC) by specifying access privileges. Authoring of
policies is probably the first step the Vault administrator performs. This guide
walks you through creating example policies for `admin` and `provisioner` users.
- [AppRole Pull Authentication](/guides/configuration/authentication.html) guide
that introduces the steps to generate tokens for machines or apps by enabling
AppRole auth backend.
- [Token and Leases](/guides/configuration/lease.html) guide helps you
understand how tokens and leases work in Vault. The understanding of the
lease hierarchy and expiration mechanism helps you plan for break glass
procedures and more.
- [Root Token Generation](/guides/configuration/generate-root.html) guide
demonstrates the workflow of regenerating root tokens. It is considered to be a
best practice not to persist the initial **root** token. If a root token needs
to be regenerated, this guide helps you walk through the task.
- [Rekeying & Rotating](/guides/configuration/rekeying-and-rotating.html) guide
provides a high-level overview of Shamir's Secret Sharing Algorithm, and how to
perform _rekey_ and _rotate_ operations in Vault.
- [Building Plugin Backends](/guides/configuration/plugin-backends.html) guide
provides steps to build, register, and mount non-database external plugin
backends.
