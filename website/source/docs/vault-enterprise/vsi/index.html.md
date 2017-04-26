---
layout: "docs"
page_title: "About the Vault Secure Introduction Client"
sidebar_current: "docs-vault-enterprise-vsi"
description: |-
  Vault Secure Introduction Client provides a turnkey way to securely introduce Vault tokens and features to applications running in various environments.

---

# About the Vault Secure Introduction Client

The Vault Secure Introduction Client is a feature of [Vault
Enterprise](https://www.hashicorp.com/vault.html) that provides a turnkey way
to securely introduce Vault tokens and features to applications running in
various environments. Currently, AWS EC2 instances are supported in conjunction
with Vault's AWS authentication backend. The client stays running until
terminated and will monitor the lifetime of retrieved Vault tokens, renewing
and reauthenticating as necessary.

Configuration is simple and can generally be performed purely using CLI flags.
Please see the [Configuration](/docs/vault-enterprise/vsi/configuration.html) page for details
on client configuration.

The [Security](/docs/vault-enterprise/vsi/security.html) page contains information and
suggestions to help deploy the client in a secure fashion. It assumes
familiarity with the AWS Authentication Backend.
