---
layout: "guides"
page_title: "Vault Reference Architecture - Guides"
sidebar_title: "Reference Architecture"
sidebar_current: "guides-operations-reference-architecture"
description: |-
  This guide provides guidance in the best practices of Vault
  implementations through use of a reference architecture.
ea_version: 1.1
---

# Vault Reference Architecture

The goal of this document is to recommend _HashiCorp Vault_ deployment
practices. This reference architecture conveys a general architecture
that should be adapted to accommodate the specific needs of each implementation.

The following topics are addressed in this guide:

- [Design Summary](#design)
    - [Network Connectivity](#network-connectivity-details)
    - [Failure Tolerance](#failure-tolerance)
- [Recommended Architecture](#recommended-architecture)
    - [Single Region Deployment](#single-region)
    - [Multiple Region Deployment](#multiple-region)
- [Best Case Architecture](#best-case-architecture)
    - [Single Availability Zone](#single-zone)
    - [Two Availability Zones - OSS](#two-zone-oss)
    - [Two Availability Zones - Enterprise](#two-zone-enterprise)
    - [Three Availability Zones - OSS](#three-zone-oss)
- [Vault Replication](#vault-replication)  
- [Deployment System Requirements](#deployment-system-requirements)
    - [Hardware Considerations](#hardware-considerations)
    - [Load Balancing](#load-balancing)
- [Additional References](#additional-references)

## Glossary
##### Vault Cluster
A Vault cluster is a set of Vault processes that together run a Vault service. These Vault processes could be running on physical or virtual servers, or in containers.
##### Consul storage backend cluster
HashiCorp recommends and supports Consul being used as the storage backend for Vault. A Consul cluster is a set of Consul server processes that together run a Consul service. These Consul processes could be running on physical or virtual servers, or in containers.
##### Availability Zone
A single failure domain on a location level that hosts part of, or all of a Vault cluster. The latency between availability zones should be < 8ms for a round trip. A single Vault cluster may be spread across multiple availability zones.
Examples of an availability zone in this context are:
- An isolated datacenter
- An isolated cage in a datacenter if it is isolated from other cages by all other means (power, network, etc)
- An availability zone in AWS, Azure or GCP

##### Region
A geographically separate collection of one or more availability zones. A region would host one or more Vault clusters. There is no defined maximum latency requirement between regions in Vault architecture. A single Vault cluster would not be spread across multiple regions.

## Design Summary

This design is the recommended architecture for production environments, as it provides flexibility and resilience.

It is a major architecture recommendation that the Consul servers are separate from the Vault servers and that the Consul cluster is only used as a storage backend for Vault and not for other Consul-focused functionality (eg service segmentation and service discovery) which can introduce unpredictable resource utilisation. Separating Vault and Consul allows each to have a system that can be sized appropriately in terms of CPU, memory and disk. Consul is a memory intensive application and so separating it to its own resources is advantageous to prevent resource contention or starvation. Dedicating a Consul cluster as a Vault storage backend is also advantageous as this allows the Consul cluster to be upgraded only as required to improve Vault storage backend functionality. This is likely to be much less frequently than a Consul cluster that is also used for service discovery and service segmentation.

Vault to Consul backend connectivity is over HTTP and should be secured with TLS as well as a Consul token to provide encryption of all traffic. See the Vault [Deployment Guide](/guides/operations/deployment-guide.html) for more information. As the Consul cluster for Vault storage may be used in addition and separate to a Consul cluster for service discovery, it is recommended that the storage Consul process be run on non-default ports so that it does not conflict with other Consul functionality. Setting the Consul storage cluster to run on 7xxx ports and using this as the storage port in the Vault configuration will achieve this.
It is also recommended that Consul be run using TLS.

-> Refer to the online documentation to learn more about running [Consul in encrypted mode](https://www.consul.io/docs/agent/options.html#encrypt).

### Network Connectivity Details
![Network Connectivity Details](/img/vault-ra-network.png)

The following table outlines the network traffic requirements for Vault cluster nodes.

| Source | Destination | port | protocol | Direction | Purpose |
|-------|----------|-----------------|-----------|-----------------|-----------|
| Consul clients and servers | Consul Server | 7300 | tcp | incoming | Server RPC |
| Consul clients | Consul clients | 7301 | tcp and udp | bidirectional | Lan gossip communications |
| Vault clients | Vault servers | 8200 | tcp | incomming | Vault API |
| Vault servers | Vault servers | 8201 | tcp | bidirectional | Vault replication traffic, request forwarding |

##### Alternative Network Configurations
Vault can be configured in several separate ways for communications between the Vault and Consul Clusters.
Using host IP addresses or hostnames that are resolvable via standard named subsystem.
Using loadbalancer IP addresses or hostnames that are resolvable via standard named subsystem.
Using the attached Consul cluster DNS as service discovery to resolve Vault endpoints.
Using a separate Consul service discovery cluster DNS as service discovery to resolve Vault endpoints.

All of these options are explored more in the Vault [Deployment Guide](/guides/operations/deployment-guide.html).


### Failure Tolerance

Vault is designed to handle different failure scenarios that have different probabilities. When deploying a Vault cluster, the failure tolerance that you require should be considered and designed for.
In OSS Vault the recommended number of instances is 3 in a cluster as any more would have limited value. In Vault Enterprise the recommended number is also 3 in a cluster, but more can be used if they were performance replicas to help with workload.
The Consul cluster is from one to seven instances. It is recommended that the Consul cluster is at least five instances that are dedicated to performing backend storage functions for the Vault cluster only.

-> Refer to the online documentation to learn more about the [Consul leader election process](https://www.consul.io/docs/guides/leader-election.html).

#### Node
The Vault and Consul cluster software allows for a failure domain at the node level by having replication within the cluster. Vault achieves this by replicating all data to all nodes within the cluster and standby servers in Vault and one of the Vault servers obtaining a lock within the data store to become the leader. If at any time the leader is lost then another Vault node will seamlessly take its place as the leader. To achieve n-2 redundancy (where the loss of 2 objects within the failure domain can be tolerated), an ideal size for a Vault cluster would be 3.
Consul achieves replication and leadership through the use of its consensus and gossip protocols. In these protocols, a leader is elected by consensus and so a quorum of active servers must always exist. To achieve n-2 redundancy, an ideal size of a Consul cluster is 5.
See Consul Internals for more details.
#### Availability Zone
Typical distribution in a cloud environment is to spread Consul/Vault nodes into separate Availability Zones (AZs) within a high bandwidth, low latency network, such as an AWS Region, however this may not be possible in a datacenter installation where there is only one DC within the level of latency required.
It is important to understand a shift in requirement or best practices that has come about as a result of the shift towards greater utilization of highly distributed systems such as Consul. When operating environments comprised of distributed systems, a shift is required in the redundancy coefficient of underlying components. Consul relies upon consensus negotiation to organize and replicate information and so the environment must provide 3 unique resilient paths in order to provide meaningful reliability. Essentially, a consensus system requires a simple majority of nodes to be available at any time. In the example of 3 nodes, you must have 2 available. If those 3 nodes are places in two failure domains, there is a 50% chance that losing a single failure domain would result in a complete outage.
#### Region
To protect against a failure at the region level, Vault enterprise offers two types of replication that can address this.
Disaster Recovery Replication
Performance Replication
Please see the Recommended Patterns on Vault Replication for a full description of these options.

Because of the constraints listed above, the recommended architecture is with Vault and Consul Enterprise distributed across three availability zones within a cluster and for clusters to be replicated across regions using DR and Performance replication. There are also several “Best Case” architecture solutions for one and two Availability Zones and also for Consul OSS. These are not the recommended architecture, but are the best solutions if your deployment is restricted by Consul version or number of availability zones.

## Recommended Architecture
The architecture below is the recommended best approach to Vault deployment and should be the target architecture for any installation. This is split into two parts:
- Vault cluster - This is the recommended architecture for a vault cluster as a single entity, but should also utilise replication as per the second diagram
- Vault replication - This is the recommended architecture for multiple vault clusters to allow for regional, performance and disaster recovery.

### Single Region Deployment (Enterprise)
##### Reference Diagram
![Reference Diagram](/img/vault-ra-3_az_ent.png)

In this scenario, the nodes in the Vault and associated Consul cluster are hosted between three Availability Zones. This solution has an n-2 at the node level for Vault and an n-3 at the node level for Consul. At the Availability Zone level, Vault is at n-2 and Consul at n-1. This differs from the OSS design in that the Consul cluster has six nodes with three of them as a [non-voting member](https://www.consul.io/docs/agent/options.html#_non_voting_server)s. If any Zone were to fail a non-voting member would be promoted by Autopilot to become a full member and so maintain Quorum.

### Multiple Region Deployment (Enterprise)
##### Reference Diagram
![Reference Diagram](/img/vault-ra-full-replication.png)

In this scenario, there is one Primary Vault cluster with a DR cluster and three Performance Replicas each with a DR cluster of its own. Each cluster has its associated Consul cluster for storage backend.  
In this setup it is recommended that the Primary Vault cluster be only used for handling replication and cluster leadership and should not be used for client connections of secrets. All client connections should go through the associated regional Performance Replica. This architecture allows for n-2 at the region level provided all secrets and secret engines are replicated across all clusters.
Failure of Region 1 would require one of the DR clusters to be promoted to primary.  
Failure of any region would result in all identities that have authenticated with Vault to have to re-authenticate with a PR in another region and this would involve additional work on the application logic side to handle.  
The advantage of this architecture over the previous one is that failure at the cluster level.  
Failure of the Primary cluster would result in the promotion of one of the DR clusters.  
Failure of a Performance Replica would result in its DR cluster being promoted to a Performance Replica and there would not need to be further logic in the client applications as all tokens and leases are also replicated in a DR replica.  
Another advantage of this architecture is that namespaces, secrets and authentication methods can be limited to different clusters so that secrets can be maintained within a region if this is required for governance purposes.
The pattern in Region 2 and 3 could be repeated multiple times, though there would unlikely be the need for further DR replicas if there were more regions.

## Best Case Architecture
In some deployments there may be insurmountable restrictions that mean the recommended architecture is not possible. This could be due to lack of availability zones or because of using Vault OSS. In these cases, the architectures below detail the best case options available.  
Note that in these following architectures the Consul leader could be any of the five Consul server nodes and the Vault active node could be any of the three Vault nodes

### Deployment of Vault in one Availability Zone (all)
##### Reference Diagram
![Reference Diagram](/img/vault-ra-1-az.png)
In this scenario, all nodes in the Vault and associated Consul cluster are hosted within one Availability Zone. This solution has a single point of failure at the availability zone level, but an n-2 at the node level for both Consul and Vault.  
This is not Hashicorp recommended architecture for production systems are there is no redundancy at the Availability Zone level. Also there is no DR capability and so as a minimum this should at least have a DR replica in a separate Region.

### Deployment of Vault in two Availability Zones (OSS)
##### Reference Diagram
![Reference Diagram](/img/vault-ra-2-az.png)

In this scenario, the nodes in the Vault and associated Consul cluster are hosted between two Availability Zones. This solution has an n-2 at the node level for Vault and Consul and n-1 for Vault at the Availability Zone level, but the addition of an Availability Zone does not significantly increase the availability of the Consul cluster. This is because the Raft protocol requires a quorum of (n/2)+1 and if Zone B were to fail in the above diagram then the cluster would not be quorate and so would also fail.  
This is not Hashicorp recommended architecture for production systems are there is only partial redundancy at the Availability Zone level and an Availability Zone failure may or may not result in an outage.

### Deployment of Vault in two Availability Zones (Enterprise)
##### Reference Diagram
![Reference Diagram](/img/vault-ra-2-az-ent.png)

In this scenario, the nodes in the Vault and associated Consul cluster are hosted between two Availability Zones. This solution has an n-2 at the node level for Vault and Consul and n-1 for Vault and Consul at the Availability Zone level. This differs from the OSS design in that the Consul cluster has six nodes with one of them as a non-voting member. If Zone B were to fail the non-voting member would be promoted by Autopilot to become a full member and so maintain Quorum. This configuration option is only available in the Enterprise version of Consul.

### Deployment of Vault in three Availability Zones (OSS)
##### Reference Diagram
![Reference Diagram](/img/vault-ra-3-az.png)

In this scenario, the nodes in the Vault and associated Consul cluster are hosted between three Availability Zones. This solution has an n-2 at the node level for Vault and Consul and n-2 for Vault at the Availability Zone level. This also has an n-1 at the Availability Zone level for Consul and as such is considered the most resilient of all architectures for a single Vault cluster with a Consul storage backend for the OSS product.

## Vault Replication (Enterprise Only)
In these architectures the “Vault Cluster” (Primary, Secondary (Performance, Disaster Recovery)) is illustrated as a single entity, and would be one of the single clusters detailed above based on your number of Availability Zones. Multiple Vault clusters acting as a single Vault solution and replicating between them is available in Enterprise Vault only. OSS Vault can be set up in multiple clusters, but they would each be individual Vault solutions and would not support replication between clusters.  
The [Vault documentation](https://www.vaultproject.io/docs/enterprise/replication/index.html) provides more detailed information on the replication capabilities within Vault Enterprise.

#### Performance Replication
Vault performance replication allows for secrets management across many sites. Secrets, authentication methods, authorization policies and other details are replicated to be active and available in multiple locations.

NOTE: Refer to the [Vault Mount Filter guide](https://learn.hashicorp.com/vault/operations/mount-filter) about filtering out secret engines from being replicated across regions.

#### Disaster Recovery Replication
Vault disaster recovery replication ensures that a standby Vault cluster is kept synchronised with an active Vault cluster. This mode of replication includes data such as ephemeral authentication tokens, time-based token information as well as token usage data. This provides for aggressive recovery point objective in environments where preventing loss of ephemeral operational data is of the utmost concern.  
NOTE: Refer to the [Vault Disaster Recovery Setup guide](https://learn.hashicorp.com/vault/operations/ops-disaster-recovery.html) for additional information.
##### Corruption or Sabotage Disaster Recovery
Another common scenario to protect against, more prevalent in cloud environments that provide very high levels of intrinsic resiliency, might be the purposeful or accidental corruption of data and configuration, and or a loss of cloud account control. Vault's DR Replication is designed to replicate live data, which would propagate intentional or accidental data corruption or deletion. To protect against these possibilities, you should backup Vault's storage backend. This is supported through the Consul Snapshot feature, which can be automated for regular archival backups. A cold site or new infrastructure could be re-hydrated from a Consul snapshot.

NOTE: Refer to the online documentation to learn more about [Consul snapshots](https://www.consul.io/docs/commands/snapshot.html)

#### Replication Notes
There is no set limit on number of clusters within a replication set. Largest deployments today are in the 30+ cluster range.  
Any cluster within a Performance replication set can act as a Disaster Recovery primary cluster.  
A cluster within a Performance replication set can also replicate to multiple Disaster Recovery secondary clusters.  
While a Vault cluster can possess a replication role (or roles), there are no special considerations required in terms of infrastructure, and clusters can assume (or be promoted) to another role. Special circumstances related to mount filters and HSM usage may limit swapping of roles, but those are based on specific organisation configurations.
#### Considerations Related to Unseal proxy_protocol_behavior
Using replication with Vault clusters integrated with HSM devices for automated unseal operations has some details that should be understood during the planning phase.

- If a performance primary cluster utilises an HSM, all other clusters within that replication set must use an HSM as well.
- If a performance primary cluster does NOT utilize an HSM (uses Shamir secret sharing method), the clusters within that replication set can be mixed, such that some may use an HSM, others may use Shamir.
- For the sake of this discussion, the cloud auto-unseal feature is treated as an HSM.

### Deployment System Requirements

The following table provides guidelines for server sizing. Of particular note is
the strong recommendation to avoid non-fixed performance CPUs, or "Burstable
CPU" in AWS terms, such as T-series instances.

#### Sizing for Vault Servers

| Size  | CPU      | Memory          | Disk      | Typical Cloud Instance Types               |
|-------|----------|-----------------|-----------|--------------------------------------------|
| Small | 2 core   | 4-8 GB RAM      | 25 GB     | **AWS:** m5.large                          |
|       |          |                 |           | **Azure:** Standard_D2_v3                  |
|       |          |                 |           | **GCE:** n1-standard-2, n1-standard-4      |
| Large | 4-8 core | 16-32 GB RAM    | 50 GB     | **AWS:** m5.xlarge, m5.2xlarge             |
|       |          |                 |           | **Azure:** Standard_D4_v3, Standard_D8_v3  |
|       |          |                 |           | **GCE:** n1-standard-8, n1-standard-16     |

#### Sizing for Consul Servers

| Size  | CPU      | Memory          | Disk      | Typical Cloud Instance Types               |
|-------|----------|-----------------|-----------|--------------------------------------------|
| Small | 2 core   | 8-16 GB RAM     | 50 GB     | **AWS:** m5.large, m5.xlarge               |
|       |          |                 |           | **Azure:** Standard_D2_v3, Standard_D4_v3  |
|       |          |                 |           | **GCE:** n1-standard-4, n1-standard-8      |
| Large | 4-8 core | 32-64+ GB RAM   | 100 GB    | **AWS:** m5.2xlarge, m5.4xlarge            |
|       |          |                 |           | **Azure:** Standard_D4_v3, Standard_D8_v3  |
|       |          |                 |           | **GCE:** n1-standard-16, n1-standard-32    |

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
written to the key-value store section of Consul's Service Catalog, which is
required to be stored in its entirety in-memory on each Consul server. This
means that memory can be a constraint in scaling as more clients authenticate to
Vault, more secrets are persistently stored in Vault, and more temporary secrets
are leased from Vault. This also has the effect of requiring vertical scaling on
Consul server's memory if additional space is required, as the entire Service
Catalog is stored in memory on each Consul server.

Furthermore, network throughput is a common consideration for Vault and Consul
servers. As both systems are HTTPS API driven, all incoming requests,
communications between Vault and Consul, underlying gossip communication between
Consul cluster members, communications with external systems (per auth or secret
engine configuration, and some audit logging configurations) and responses
consume network bandwidth.

Due to network performance considerations in Consul cluster operations,
replication of Vault datasets across network boundaries should be achieved
through Performance or DR Replication, rather than spreading the Consul cluster
across network and physical boundaries.  If a single consul cluster is spread
across network segments that are distant or inter-regional, this can cause
synchronization issues within the cluster or additional data transfer charges
in some cloud providers.

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

![Vault Behind a Load Balancer](/img/vault-ref-arch-9.png)

External load balancers are supported as well, and would be placed in front of the
Vault cluster, and would poll specific Vault URL's to detect the active node and
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

There are two supported methods for handling client IP addressing behind a proxy
or load balancer;
[X-Forwarded-For Headers](https://www.vaultproject.io/docs/configuration/listener/tcp.html#x_forwarded_for_authorized_addrs)
and [PROXY v1](https://www.vaultproject.io/docs/configuration/listener/tcp.html#proxy_protocol_authorized_addrs).  Both require a trusted load balancer and require IP address whitelisting to
adhere to security best practices.

## Additional References

- Vault [architecture](/docs/internals/architecture.html) documentation explains
each Vault component
- To integrate Vault with existing LDAP server, refer to
[LDAP Auth Method](/docs/auth/ldap.html) documentation
- Refer to the [AppRole Pull
Authentication](/guides/identity/authentication.html) guide to programmatically
generate a token for a machine or app
- Consul is an integral part of running a resilient Vault cluster, regardless of
location.  Refer to the online [Consul documentation](https://www.consul.io/intro/getting-started/install.html) to
learn more.

## Next steps

- Read [Production Hardening](/guides/operations/production.html) to learn best
  practices for a production hardening deployment of Vault.

- Read [Deployment Guide](/guides/operations/deployment-guide.html) to learn
  the steps required to install and configure a single HashiCorp Vault cluster.
