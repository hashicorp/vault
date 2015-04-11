---
layout: "docs"
page_title: "Security Model"
sidebar_current: "docs-internals-security"
description: |-
  Learn about the security model of Vault.
---

# Security Model

Due to the nature of Vault and the confidentiality of data it is managing,
the Vault security model is very critical. The overall goal of Vault's security
model is to provide [confidentiality, integrity, availability, accountability,
authentication](http://en.wikipedia.org/wiki/Information_security).

This means that data at rest and in transit must be secure from eavesdropping
or tampering. Clients must be appropriately authenticated and authorized
to access data or modify policy. All interactions must be auditable and traced
uniquely back to the origin entity. The system must be robust against intentional
attempts to bypass any of its access controls.

~> **Advanced Topic!** This page covers the technical details
of the security model of Vault. You don't need to understand these details to
effectively use Vault. The details are documented here for
those who wish to learn about them without having to go
spelunking through the source code.

# Threat Model

The following are the various parts of the Vault threat model:

* Eavesdropping on any Vault communication. Client communication with Vault
  should be secure from eavesdropping as well as communication from Vault to
  its storage backend.

* Tampering with data at rest or in transit. Any tampering should be detectable
  and cause Vault to abort processing of the transaction.

* Access to data or controls without authentication or authorization. All requests
  must be proceeded by the applicable security policies.

* Access to data or controls without accountability. If audit logging
  is enabled, requests and responses must be logged before the client receives
  any secret material.

* Confidentiality of stored secrets. Any data that leaves Vault to rest in the
  storage backend must be safe from eavesdropping. In practice, this means all
  data at rest must be encrypted.

* Availability of secret material in the face of failure. Vault supports
  running in a highly available configuration to avoid loss of availability.

The following are not parts of the Vault threat model:

* Protecting against arbitrary control of the storage backend. An attacker
  that can perform arbitrary operations against the storage backend can
  undermine in any number of ways that are difficult or impossible to protect
  against. As an example, an attacker could delete or corrupt all the contents
  of the storage backend causing  total data loss for Vault. The ability to controls
  reads would allow an attacker to snapshot in a well-known state and rollback state
  changes if that would be beneficial to them.

* Protecting against the leakage of the existance of secret material. An attacker
  that can read from the storage backend may observe that secret material exists
  and is stored, even if it is kept confidential.

* Protecting against memory analysis of a running Vault. If an attacker is able
  to inspect the memory state of a running Vault instance then the confidentiality
  of data may be compromised.

