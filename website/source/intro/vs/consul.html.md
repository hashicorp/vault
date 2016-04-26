---
layout: "intro"
page_title: "Vault vs. Consul"
sidebar_current: "vs-other-consul"
description: |-
  Comparison between Vault and attempting to store secrets with Consul.
---

# Vault vs. Consul

[Consul](https://www.consul.io) is a system for service discovery, monitoring,
and configuration that is distributed and highly available. Consul also
supports an ACL system to restrict access to keys and service information.

While Consul can be used to store secret information and gate access using
ACLs, it is not designed for that purpose. As such, data is not encrypted
in transit nor at rest, it does not have pluggable authentication mechanisms,
and there is no per-request auditing mechanism.

Vault is designed from the ground up as a secret management solution. As such,
it protects secrets in transit and at rest. It provides multiple authentication
and audit logging mechanisms. Dynamic secret generation allows Vault to avoid
providing clients with root privileges to underlying systems and makes
it possible to do key rolling and revocation.

The strength of Consul is that it is fault tolerant and highly scalable.
By using Consul as a backend to Vault, you get the best of both. Consul
is used for durable storage of encrypted data at rest and provides coordination
so that Vault can be highly available and fault tolerant. Vault provides
the higher level policy management, secret leasing, audit logging, and automatic
revocation.

