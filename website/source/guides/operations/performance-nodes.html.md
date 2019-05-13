---
layout: "guides"
page_title: "Performance Standby Nodes - Guides"
sidebar_current: "guides-operations-performance-nodes"
description: |-
  This guide will walk you through a simple Vault Highly Available (HA) cluster
  implementation. While this is not an exhaustive or prescriptive guide that
  can be used as a drop-in production example, it covers the basics enough to
  inform your own production setup.
---

# Performance Standby Nodes

~> **Enterprise Only:** Performance Standby Nodes feature is a part of _Vault Enterprise_.

In [Vault High Availability](/guides/operations/vault-ha-consul.html) guide, it
was explained that only one Vault server will be _active_ in a cluster and
handles **all** requests (reads and writes).  The rest of the servers become the
_standby_ nodes and simply forward requests to the _active_ node.  

![HA Architecture](/img/vault-ha-consul-3.png)

If you are running **_Vault Enterprise_ 0.11** or later with the Consul storage 
backend, those standby nodes can handle most read-only requests. For example, 
performance standbys can handle encryption/decryption of data using 
[transit](/docs/secrets/transit/index.html) keys, GET requests of key/value 
secrets and other requests that do not change underlying storage. This can 
provide considerable improvements in throughput for traffic of this type, 
resulting in aggregate performance increase linearly correlated to the number 
of performance standby nodes deployed in a cluster.


## Reference Materials

- [Performance Standby Nodes](/docs/enterprise/performance-standby/index.html)
- [High Availability Mode](/docs/concepts/ha.html)
- [Consul Storage Backend](/docs/configuration/storage/consul.html)
- [Vault Reference Architecture](/guides/operations/reference-architecture.html)


## Server Configuration

Performance standbys are enabled by default when the Vault Enterprise license
includes this feature. If you wish to disable the performance standbys, you can
do so by setting the
[`disable_performance_standby`](/docs/configuration/index.html#vault-enterprise-parameters)
flag to `true`.  

Since any of the nodes in a cluster can get elected as active, it is recommended
to keep this setting consistent across all nodes in the cluster.

!> Consider a scenario where a node with performance standby _disabled_
becomes the active node. The performance standby feature is
disabled for the whole cluster although it is enabled on other nodes.


## Enterprise Cluster

A highly available Vault Enterprise cluster consists of multiple servers, and
there will be only one active node. The rest can serve as performance standby
nodes handling read-only requests locally.

![Cluster Architecture](/img/vault-perf-standby-1.png)

The number of performance standby nodes within a cluster depends on your Vault
Enterprise license.

Consider the following scenario:

- A cluster contains **five** Vault servers
- Your Vault Enterprise license allows **two** performance standby nodes

![Cluster Architecture](/img/vault-perf-standby.png)

In this scenario, the performance standby nodes running on VM 8 and VM 9 can
process read-only requests. However, the _standby_ nodes running on VM 6 and VM
10 simply forward all requests to the active node running on VM 7.


> **NOTE:** The selection of performance standby node is determined by the
active node. When a node is selected, it gets  promoted to become a performance
standby. This is a race condition that there is no configuration
parameter to specify which nodes to become performance standbys.


## Next steps

Read [Production Hardening](/guides/operations/production.html) to learn best
practices for a production hardening deployment of Vault.
