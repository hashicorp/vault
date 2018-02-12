---
layout: "guides"
page_title: "Secrets Management - Guides"
sidebar_current: "guides-secret-mgmt"
description: |-
   A very common use case of Vault is to manage your organization's secrets from
   storing credentials and API keys to encrypting passwords for user signups.
   Vault is meant to be a solution for all secret management needs.
---

# Secrets Management

Vault is a tool for securely accessing secrets. A secret is anything that you
want to tightly control access to, such as API keys, passwords, certificates,
and more. Vault provides a unified interface to any secret while providing
tight access control and recording a detailed audit log.

Secrets Management guides demonstrate features in Vault to securely store your
secrets.

- [Static Secrets](/guides/secret-mgmt/static-secrets.html) guide walks you
through the steps to write secrets in Vault, and control who can access them.

- [Secret as a Service: Dynamic Secrets](/guides/secret-mgmt/dynamic-secrets.html)
 guide demonstrates the Vault feature to generate database credentials
 on-demand so that each application or system can obtain its own credentials,
 and its permissions can be tightly controlled.

- [Cubbyhole Response Wrapping](/guides/secret-mgmt/cubbyhole.html) guide
demonstrates a secure method to distribute secrets by wrapping them where only
the expecting client can unwrap.
