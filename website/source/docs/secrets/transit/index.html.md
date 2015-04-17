---
layout: "docs"
page_title: "Secret Backend: Trasit"
sidebar_current: "docs-secrets-transit"
description: |-
  The transit secret backend for Vault encrypts/decrypts data in-transit. It doesn't store any secrets.
---

# Transit Secret Backend

Name: `transit`

The transit secret backend is used to encrypt/data in-transit. Vault doesn't
store the data sent to the backend. It can also be viewed as "encryption as
a service."

The primary use case for the transit backend is to encrypt data from
applications while still storing that encrypted data in some primary data
store. This relieves the burden of proper encryption/decryption from
application developers and pushes the burden onto the operators of Vault.
Operators of Vault generally include the security team at an organization,
which means they can ensure that data is encrypted/decrypted properly.

Additionally, since encrypt/decrypt operations must enter the audit log,
any decryption event is recorded.

This page will show a quick start for this backend. For detailed documentation
on every path, use `vault help` after mounting the backend.

## Quick Start

TODO
