---
layout: "intro"
page_title: "Use Cases"
sidebar_current: "use-cases"
description: |-
  This page lists some concrete use cases for Vault, but the possible use cases are much broader than what we cover.
---

# Use Cases

Before understanding use cases, it's useful to know [what Vault is](/intro/index.html).
This page lists some concrete use cases for Vault, but the possible use cases are
much broader than what we cover.

#### General Secret Storage

At a bare minimum, Vault can be used for the storage of any secrets. For
example, Vault would be a fantastic way to store sensitive environment variables,
database credentials, API keys, etc.

Compare this with the current way to store these which might be
plaintext in files, configuration management, a database, etc. It would be
much safer to query these using `vault read` or the API. This protects
the plaintext version of these secrets as well as records access in the Vault
audit log.

#### Employee Credential Storage

While this overlaps with "General Secret Storage", Vault is a good mechanism
for storing credentials that employees share to access web services. The
audit log mechanism lets you know what secrets an employee accessed and
when an employee leaves, it is easier to roll keys and understand which keys
have and haven't been rolled.

#### API Key Generation for Scripts

The "dynamic secrets" feature of Vault is ideal for scripts: an AWS
access key can be generated for the duration of a script, then revoked.
The keypair will not exist before or after the script runs, and the
creation of the keys are completely logged.

This is an improvement over using something like Amazon IAM but still
effectively hardcoding limited-access access tokens in various places.

#### Data Encryption

In addition to being able to store secrets, Vault can be used to
encrypt/decrypt data that is stored elsewhere. The primary use of this is
to allow applications to encrypt their data while still storing it in the
primary data store.

The benefit of this is that developers do not need to worry about how to
properly encrypt data. The responsibility of encryption is on Vault
and the security team managing it, and developers just encrypt/decrypt
data as needed.
