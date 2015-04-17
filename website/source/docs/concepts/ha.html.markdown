---
layout: "docs"
page_title: "High Availibility"
sidebar_current: "docs-concepts-ha"
description: |-
  Vault can be highly available, allowing you to run multiple Vaults to protect against outages.
---

# High Availability Mode (HA)

Vault supports multi-server mode for high availability. This mode protects
against outages by running multiple Vault servers. High availability mode
is automatically enabled when using a physical backend that supports it.

You can tell if a backend supports high availability mode ("HA") by
starting the server and seeing if "(HA available)" is outputted next to
the backend information. If it is, then HA will begin happening automatically.

To be highly available, Vault elects a leader and does request forwarding to
the leader. Due to this architecture, HA does not enable increased scalability.
In general, the bottleneck of Vault is the physical backend itself, not
Vault core. For example: to increase scalability of Vault with Consul, you
would scale Consul instead of Vault.

In addition to using a backend that supports HA, you have to configure
Vault with an _advertise address_. This is the address that Vault advertises
to other Vault servers in the cluster for request forwarding. By default,
Vault will use the first private IP address it finds, but you can override
this to any address you want.

## Backend Support

Currently, the only backend that supports HA is Consul.

If you're interested in implementing another backend or adding HA support
to another backend, we'd love your contributions. Adding HA support
requires implementing the `physical.HABackend` interface for the physical
backend.
