---
layout: "docs"
page_title: "High Availability"
sidebar_current: "docs-internals-ha"
description: |-
  Learn about the high availability design of Vault.
---

# High Availability

Vault is primarily used in production environments to manage secrets.
As a result, any downtime of the Vault service can affect downstream clients.
Vault is designed to support a highly available deploy to ensure a machine
or process failure is minimally disruptive.

~> **Advanced Topic!** This page covers technical details
of Vault. You don't need to understand these details to
effectively use Vault. The details are documented here for
those who wish to learn about them without having to go
spelunking through the source code. However, if you're an
operator of Vault, we recommend learning about the architecture
due to the importance of Vault in an environment.

# Design Overview

The primary design goal in making Vault highly available (HA) was to
minimize downtime and not horizontal scalability. Vault is typically
bound by the IO limits of the storage backend rather than the compute
requirements. This simplifies the HA approach and allows more complex
coordination to be avoided.

Certain storage backends, such as Consul, provide additional coordination
functions that enable Vault to run in an HA configuration. When supported
by the backend, Vault will automatically run in HA mode without additional
configuration.

When running in HA mode, Vault servers have two additional states they
can be in: standby and active. For multiple Vault servers sharing a storage
backend, only a single instance will be active at any time while all other
instances are hot standbys.

The active server operates in a standard fashion and processes all requests.
The standby servers do not process requests, and instead redirect to the active
Vault. Meanwhile, if the active server is sealed, fails, or loses network connectivity
then one of the standbys will take over and become the active instance.

It is important to note that only _unsealed_ servers act as a standby.
If a server is still in the sealed state, then it cannot act as a standby
as it would be unable to serve any requests should the active server fail.
