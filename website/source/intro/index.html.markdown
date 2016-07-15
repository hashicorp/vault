---
layout: "intro"
page_title: "Introduction"
sidebar_current: "what"
description: |-
  Welcome to the intro guide to Vault! This guide is the best place to start with Vault. We cover what Vault is, what problems it can solve, how it compares to existing software, and contains a quick start for using Vault.
---

# Introduction to Vault

Welcome to the intro guide to Vault! This guide is the best
place to start with Vault. We cover what Vault is, what
problems it can solve, how it compares to existing software,
and contains a quick start for using Vault.

If you are already familiar with the basics of Vault, the
[documentation](/docs/index.html) provides a better reference
guide for all available features as well as internals.

## What is Vault?

Vault is a tool for securely accessing _secrets_. A secret is anything
that you want to tightly control access to, such as API keys, passwords,
certificates, and more. Vault provides a unified interface to any
secret, while providing tight access control and recording a detailed
audit log.

A modern system requires access to a multitude of secrets: database
credentials, API keys for external services, credentials for
service-oriented architecture communication, etc. Understanding who is
accessing what secrets is already very difficult and platform-specific.
Adding on key rolling, secure storage, and detailed audit logs is almost
impossible without a custom solution. This is where Vault steps in.

Examples work best to showcase Vault. Please see the
[use cases](/intro/use-cases.html).

The key features of Vault are:

* **Secure Secret Storage**: Arbitrary key/value secrets can be stored
  in Vault. Vault encrypts these secrets prior to writing them to persistent
  storage, so gaining access to the raw storage isn't enough to access
  your secrets. Vault can write to disk, [Consul](https://www.consul.io),
  and more.

* **Dynamic Secrets**: Vault can generate secrets on-demand for some
  systems, such as AWS or SQL databases. For example, when an application
  needs to access an S3 bucket, it asks Vault for credentials, and Vault
  will generate an AWS keypair with valid permissions on demand. After
  creating these dynamic secrets, Vault will also automatically revoke them
  after the lease is up.

* **Data Encryption**: Vault can encrypt and decrypt data without storing
  it. This allows security teams to define encryption parameters and
  developers to store encrypted data in a location such as SQL without
  having to design their own encryption methods.

* **Leasing and Renewal**: All secrets in Vault have a _lease_ associated
  with them. At the end of the lease, Vault will automatically revoke that
  secret. Clients are able to renew leases via built-in renew APIs.

* **Revocation**: Vault has built-in support for secret revocation. Vault
  can revoke not only single secrets, but a tree of secrets, for example
  all secrets read by a specific user, or all secrets of a particular type.
  Revocation assists in key rolling as well as locking down systems in the
  case of an intrusion.

## Next Steps

See the page on [Vault use cases](/intro/use-cases.html) to see the
multiple ways Vault can be used. Then see
[how Vault compares to other software](/intro/vs/index.html)
to see how it fits into your existing infrastructure. Finally, continue onwards with
the [getting started guide](/intro/getting-started/install.html) to use
Vault to read, write, and create real secrets and see how it works in practice.
