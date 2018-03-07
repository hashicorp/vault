---
layout: "guides"
page_title: "Encryption as a Service - Guides"
sidebar_current: "guides-encryption"
description: |-
  The transit secrets engine handles cryptographic functions on data in-transit.
  Vault doesn't store the data sent to the secrets engine. It can also be viewed
  as "cryptography as a service" or "encryption as a service".
---

# Encryption as a Service

Vault provides Encryption as a Service (EaaS) to enables security teams to
fortify data during transit and at rest. So even if an intrusion occurs, your
data is encrypted with AES 256-bit CBC encryption (TLS in transit). Even if an
attacker were able to access the raw data, they would only have encrypted bits.
This means attackers would need to compromise multiple systems before
exfiltrating data.

This guide walks you through Encryption as a Service topics.

- [Transit Secrets Re-wrapping](/guides/encryption/transit-rewrap.html) guide
demonstrate one possible way to re-wrap data after rotating an encryption key
in the transit engine in Vault.
