---
layout: "guides"
page_title: "Vault Deployment Reference Architecture - Guides"
sidebar_current: "guides-operations-reference-architecture"
description: |-
  This guide provides guidance in the best practices of Vault
  implementations through use of a reference architecture.
---

# Vault Deployment Reference Architecture

The goal of this document is to recommend _HashiCorp Vault_ deployment
practices. This reference architecture conveys a general architecture
that should be adapted to accommodate the specific needs of each implementation.

The following topics are addressed in this guide:

- [Deployment Topology within One Datacenter](#one-dc)
    - [Network Connectivity](#network-connectivity-details)
    - [Deployment System Requirements](#deployment-system-requirements)
    - [Hardware Considerations](#hardware-considerations)
    - [Load Balancing](#load-balancing)
    - [High Availability](#high-availability)
- [Deployment Topology for Multiple Datacenters](#multi-dc)
    - [Vault Replication](#vault-replication)
- [Additional References](#additional-references)

-> This document assumes Vault uses Consul as the [storage
backend](/docs/internals/architecture.html) since that is the recommended
storage backend for production deployments.

## <a name="one-dc"></a>Deployment Topology within One Datacenter

This section explains how to deploy a Vault open source cluster in one datacenter.
Once a Vault cluster is running succesfully within one datacenter,
move on to deploying Vault Enterprise across [multiple datacenters](#multi-dc).

### Reference Diagram

Eight Nodes with [Consul Storage Backend](/docs/configuration/storage/consul.html)
![Reference diagram](/assets/images/vault-ref-arch-2.png)

#### Design Summary

This design is the recommended architecture for production environments, as it
provides flexibility and resilience. Consul servers are separate
from the Vault servers so that software upgrades are easier to perform. Additionally,
separate Consul and Vault servers allows for separate sizing for each.
Vault to Consul backend connectivity is over HTTP and should be
secured with TLS as well as a Consul token to provide encryption of all traffic.

#### Failure Tolerance

Typical distribution in a cloud environment is to spread Consul/Vault nodes into
separate Availability Zones (AZs) within a high bandwidth, low latency network,
such as an AWS Region. The diagram below shows Vault and Consul spread between
AZs, with Consul servers in Redundancy Zone configurations, promoting a single
voting member per AZ, providing both Zone and Node level failure protection.

![Failure tolerance|40%](/assets/images/vault-ref-arch-3.png)

### Network Connectivity Details

![Network Connectivity Details](/assets/images/vault-ref-arch.png)

### Deployment System Requirements

The following table provides guidelines for server sizing. Of particular note is
the strong recommendation to avoid non-fixed performance CPUs, or “Burstable
CPU” in AWS terms, such as T-series instances.

#### Sizing for Vault Servers

| Size  | CPU      | Memory          | Disk      | Typical Cloud Instance Types               |
|-------|----------|-----------------|-----------|--------------------------------------------|
| Small | 2 core   | 4-8 GB RAM      | 25 GB     | **AWS:** m5.large                          |
|       |          |                 |           | **Azure:** Standard_A2_v2, Standard_A4_v2  |
|       |          |                 |           | **GCE:** n1-standard-4, n1-standard-8      |
| Large | 4-8 core | 16-32 GB RAM    | 50 GB     | **AWS:** m5.xlarge, m5.2xlarge             |
|       |          |                 |           | **Azure:** Standard_D4_v3, Standard_D8_v3  |
|       |          |                 |           | **GCE:** n1-standard-16, n1-standard-32    |

#### Sizing for Consul Servers

| Size  | CPU      | Memory          | Disk      | Typical Cloud Instance Types               |
|-------|----------|-----------------|-----------|--------------------------------------------|
| Small | 2 core   | 8-16 GB RAM     | 50 GB     | **AWS:** m5.large, m5.xlarge               |
|       |          |                 |           | **Azure:** Standard_A4_v2, Standard_A8_v2  |
|       |          |                 |           | **GCE:** n1-standard-8, n1-standard-16     |
| Large | 4-8 core | 32-64+ GB RAM   | 100 GB    | **AWS:** m5.2xlarge, m5.4xlarge            |
|       |          |                 |           | **Azure:** Standard_D4_v3, Standard_D5_v3  |
|       |          |                 |           | **GCE:** n1-standard-32, n1-standard-64    |

### Hardware Considerations

The small size category would be appropriate for most initial production
deployments, or for development/testing environments.  

The large size is for production environments where there is a consistent high
workload. That might be a large number of transactions, a large number of
secrets, or a combination of the two.

In general, processing requirements will be dependent on encryption workload and
messaging workload (operations per second, and types of operations).  Memory
requirements will be dependent on the total size of secrets/keys stored in
memory and should be sized according to that data (as should the hard drive
storage).  Vault itself has minimal storage requirements, but the underlying
storage backend should have a relatively high-performance hard disk subsystem.
If many secrets are being generated/rotated frequently, this information will
need to flush to disk often and can impact performance if slower hard drives are
used.

Consul servers function in this deployment is to serve as the storage backend
for Vault. This means that all content stored for persistence in Vault is
encrypted by Vault, and written to the storage backend at rest. This data is
written to the key-value store section of Consul’s Service Catalog, which is
required to be stored in its entirety in-memory on each Consul server. This
means that memory can be a constraint in scaling as more clients authenticate to
Vault, more secrets are persistently stored in Vault, and more temporary secrets
are leased from Vault. This also has the effect of requiring vertical scaling on
Consul server’s memory if additional space is required, as the entire Service
Catalog is stored in memory on each Consul server.

Furthermore, network throughput is a common consideration for Vault and Consul
servers. As both systems are HTTPS API driven, all incoming requests,
communications between Vault and Consul, underlying gossip communication between
Consul cluster members, communications with external systems (per auth or secret
engine configuration, and some audit logging configurations) and responses
consume network bandwidth.

### Other Considerations

[Vault Production Hardening Recommendations](/guides/operations/production.html)
provides guidance on best practices for a production hardened deployment of
Vault.

## Load Balancing

### <a name="consul-lb"></a>Load Balancing Using Consul Interface

Consul can provide load balancing capabilities, but it requires that any Vault
clients are Consul aware. This means that a client can either utilize Consul DNS
or API interfaces to resolve the active Vault node. A client might access Vault
via a URL like the following: `http://active.vault.service.consul:8200`

This relies upon the operating system DNS resolution system, and
the request could be forwarded to Consul for the actual IP address response.
The operation can be completely transparent to legacy applications and would
operate just as a typical DNS resolution operation.

### <a name="external-lb"></a>Load Balancing Using External Load Balancer

External load balancers are supported as well, and would be placed in front of the
Vault cluster, and would poll specific Vault URL’s to detect the active node and
route traffic accordingly. An HTTP request to the active node with the following
URL will respond with a 200 status: `http://<Vault Node URL>:8200/v1/sys/health`

The following is a sample configuration block from HAProxy to illustrate:

```plaintext
listen vault
    bind 0.0.0.0:80
    balance roundrobin
    option httpchk GET /v1/sys/health
    server vault1 192.168.33.10:8200 check
    server vault2 192.168.33.11:8200 check
    server vault3 192.168.33.12:8200 check
```

Note that the above block could be generated by Consul (with consul-template)
when a software load balancer is used. This could be the case when the load
balancer is software like Nginx, HAProxy, or Apache.

**Example Consul Template for the above HAProxy block:**

```plaintext
listen vault
   bind 0.0.0.0:8200
   balance roundrobin
   option httpchk GET /v1/sys/health{{range service "vault"}}
   server {{.Node}} {{.Address}}:{{.Port}} check{{end}}
```

#### Client IP Address Handling

Vault does not support X-Forwarded-For at this time due to security concerns, as
an attacker can perform a header injection to impersonate another client.  Vault
does have PROXY v1 protocol support. This is well supported by
[HAProxy](https://www.haproxy.org/download/1.8/doc/proxy-protocol.txt) and [AWS
ELB’s](http://docs.aws.amazon.com/elasticloadbalancing/latest/classic/enable-proxy-protocol.html).
Another [HAProxy blog post](http://www.haproxy.com/blog/haproxy/proxy-protocol/)
with some good info on PROXY protocol. It is worth pointing out that [F5 load
balancers do support it as
well](https://devcentral.f5.com/codeshare/proxy-protocol-initiator).

Note that the Vault listener must be [properly
configured](/docs/configuration/listener/tcp.html#proxy_protocol_behavior)
to support this functionality.


### High Availability

A Vault cluster is the high-available unit of deployment within one datacenter.
A recommended approach is three Vault servers with a Consul storage backend.
With this configuration, during a Vault server outage, failover is handled
immediately without human intervention. For more information on this topic,
please see the [Vault site documentation about HA
concepts](/docs/concepts/ha.html).

High-availability and data locality across datacenters requires
Vault Enterprise.


## <a name="multi-dc"></a>Deployment Topology for Multiple Datacenters

<img src="/assets/images/vault-ref-arch-6.png">

### Vault Replication

~> **Enterprise Only:** Vault replication feature is a part of _Vault Enterprise_.

HashiCorp Vault Enterprise provides two modes of replication, performance
and disaster recovery. The [Vault documentation](/docs/enterprise/replication/index.html)
provides more detailed information on the replication capabilities within Vault Enterprise.

#### Performance Replication

Vault performance replication allows for secrets management across many sites.
Secrets, authentication methods, authorization policies and other details are
replicated to be active and available in multiple locations.

#### Disaster Recovery Replication

Vault disaster recovery replication ensures that a standby Vault cluster is kept
synchronized with an active Vault cluster.  This mode of replication includes
data such as ephemeral authentication tokens, time-based token information as
well as token usage data. This provides for aggressive recovery point objective
in environments where high availability is of the utmost concern.

#### Diagrams

The following diagrams illustrate some possible replication scenarios.
![Replication Pattern|40%](/assets/images/vault-ref-arch-4.png)

#### Replication Notes

- There is no set limit on number of clusters within a replication set. Largest
deployments today are in the 30+ cluster range.
- Any cluster within a Performance replication set can act as a Disaster
Recovery primary cluster.
- A cluster within a Performance replication set can also replicate to multiple
Disaster Recovery secondary clusters.
- While a Vault cluster can possess a replication role (or roles), there are no
special considerations required in terms of infrastructure, and clusters can
assume (or be promoted) to another role. Special circumstances related to mount
filters and HSM usage may limit swapping of roles, but those are based on
specific organization configurations.

#### Considerations Related to Unseal proxy_protocol_behavior

Using replication with Vault clusters integrated with HSM devices for automated
unseal operations has some details that should be understood during the planning
phase.

- If a **performance** primary cluster utilizes an HSM, all other clusters
within that replication set must use an HSM as well.
- If a **performance** primary cluster does NOT utilize an HSM (uses Shamir
  secret sharing method), the clusters within that replication set can be mixed,
  such that some may use an HSM, others may use Shamir.

For sake of this discussion, the [cloud
auto-unseal](/docs/enterprise/auto-unseal/index.html) feature is treated as an
HSM.

## Additional References

- Vault [architecture](/docs/internals/architecture.html) documentation explains
each Vault component
- To integrate Vault with existing LDAP server, refer to
[LDAP Auth Method](/docs/auth/ldap.html) documentation
- Refer to the [AppRole Pull
Authentication](/guides/identity/authentication.html) guide to programmatically
generate a token for a machine or app

## Next steps

Read [Production Hardening](/guides/operations/production.html) to learn best
practices for a production hardening deployment of Vault.
