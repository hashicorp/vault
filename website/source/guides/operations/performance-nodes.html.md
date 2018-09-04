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

![HA Architecture](/assets/images/vault-ha-consul-3.png)

If you are running **_Vault Enterprise_ 0.11** or later, those standby
nodes can handle most read-only requests and behave as read-replica.  

~> This Performance Standby Nodes feature is included in _Vault Enterprise
Premium_, and also available for _Vault Enterprise Pro_ with additional fee.

This is particularly useful for processing high volume [_Encryption as a
Service_](/docs/secrets/transit/index.html) requests.

## Reference Materials

- [Performance Standby Nodes](/docs/enterprise/performance-standby/index.html)
- [High Availability Mode](/docs/concepts/ha.html)
- [Consul Storage Backend](/docs/configuration/storage/consul.html)
- [Vault Reference Architecture](/guides/operations/reference-architecture.html)


## Server Configuration

Performance standby is enabled by default. If you wish to disable the
performance standbys, you can do so by setting the
[`disable_performance_standby`](/docs/configuration/index.html#vault-enterprise-parameters)
flag to `true`.  

Since any of the nodes in a cluster can get elected as active, it is recommended
to keep this setting consistent across all nodes in the cluster.

!> Consider the situation where a node with performance standby _disabled_
becomes the active node. In such a case, the performance standby feature is
disabled for the whole cluster although it is enabled on other nodes.


## Enterprise Cluster

A highly available Vault Enterprise cluster consists of multiple servers, and
there will be only one active node, and the rest can serve as performance
standby nodes handling read-only requests locally.

![Cluster Architecture](/assets/images/vault-perf-standby-1.png)

The number of performance standby nodes within a cluster depends on your Vault
Enterprise license.

Let's assume the following:

- A cluster contains 5 Vault servers
- Vault Enterprise license allows 2 performance standby nodes

![Cluster Architecture](/assets/images/vault-perf-standby.png)

In this scenario, the performance standby nodes running on VM 8 and VM 9 can
process read-only requests. However, the _standby_ nodes running on VM 6 and VM
10 simply forward all requests to the active node running on VM 7.

> **NOTE:** The selection of performance standby nodes is similar to the
selection of the active node.


Your global deployment might implement [performance
replication](/guides/operations/mount-filter.html) and [disaster recovery
replication](/guides/operations/disaster-recovery.html) to scale with multiple
datacenters.

![Global Architecture](/assets/images/vault-perf-standby-2.png)



## Next steps

Read [Production Hardening](/guides/operations/production.html) to learn best
practices for a production hardening deployment of Vault.
