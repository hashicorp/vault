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

Vault provides Encryption as a Service (EaaS) to enable security teams to
fortify data during transit and at rest. So even if an intrusion occurs, your
data is encrypted and the attacker would never get a hold of the raw data.

This guide walks you through Encryption as a Service topics.

- [Java Application Demo](/guides/encryption/spring-demo.html) guide walks
through a sample application which relies on Vault to generate database
credentials as well as encrypting sensitive data.  This guide is for anyone who
wishes to reproduce the demo introduced in
the [Manage secrets, access, and encryption in the public cloud with Vault](https://www.hashicorp.com/resources/solutions-engineering-webinar-series-episode-2-vault)
webinar.

- [Transit Secrets Re-wrapping](/guides/encryption/transit-rewrap.html) guide
demonstrates one possible way to re-wrap data after rotating an encryption key
in the transit engine in Vault.
