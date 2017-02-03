---
title: "About the Vault Secure Introduction Client"
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
Please see the [Configuration](/help/vault/vsi/configuration) page for details
on client configuration.

The [Security](/help/vault/vsi/security) page contains information and
suggestions to help deploy the client in a secure fashion. It assumes
familiarity with the AWS Authentication Backend.

The client provides version information by specifying the `-v` flag to the CLI.
The [Changelog](/help/vault/vsi/changelog) details the changes and bug fixes
between versions.
