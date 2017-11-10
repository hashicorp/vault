---
layout: "docs"
page_title: "High Availability"
sidebar_current: "docs-concepts-ha"
description: |-
  Vault can be highly available, allowing you to run multiple Vaults to protect against outages.
---

# High Availability Mode (HA)

Vault supports a multi-server mode for high availability. This mode protects
against outages by running multiple Vault servers. High availability mode
is automatically enabled when using a data store that supports it.

You can tell if a data store supports high availability mode ("HA") by starting
the server and seeing if "(HA available)" is output next to the data store
information. If it is, then Vault will automatically use HA mode. This
information is also available on the
[Configuration](https://www.vaultproject.io/docs/configuration/index.html) page.

To be highly available, one of the Vault server nodes grabs a lock within the
data store. The successful server node then becomes the active node; all other
nodes become standby nodes. At this point, if the standby nodes receive a
request, they will either forward the request or redirect the client depending
on the current configuration and state of the cluster -- see the sections below
for details. Due to this architecture, HA does not enable increased
scalability. In general, the bottleneck of Vault is the data store itself, not
Vault core. For example: to increase the scalability of Vault with Consul, you
would generally scale Consul instead of Vault.

Certain storage backends can support high availability mode, which enable them
to store both Vault's information in addition to the HA lock. However, Vault
also supports split data/HA mode, whereby the lock value and the rest of the
data live separately. This can be done by specifying both the `storage` and
`ha_storage` stanzas in the configuration file with different backends. For
instance, a Vault cluster can be set up to use Consul as the `ha_storage` to
manage the lock, and use Amazon S3 as the `storage` for all other persisted
data.

The sections below explain the server communication patterns and each type of
request handling in more detail. At a minimum, the requirements for redirection
mode must be met for an HA cluster to work successfully.

## Server-to-Server Communication

Both methods of request handling rely on the active node advertising
information about itself to the other nodes. Rather than over the network, this
communication takes place within Vault's encrypted storage; the active node
writes this information and unsealed standby Vault nodes can read it.

For the client redirection method, this is the extent of server-to-server
communication -- no direct communication with only encrypted entries in the
data store used to transfer state.

For the request forwarding method, the servers need direct communication with
each other. In order to perform this securely, the active node also advertises,
via the encrypted data store entry, a newly-generated private key (ECDSA-P521)
and a newly-generated self-signed certificate designated for client and server
authentication. Each standby uses the private key and certificate to open a
mutually-authenticated TLS 1.2 connection to the active node via the advertised
cluster address. When client requests come in, the requests are serialized,
sent over this TLS-protected communication channel, and acted upon by the
active node. The active node then returns a response to the standby, which
sends the response back to the requesting client.

## Request Forwarding

If request forwarding is enabled (turned on by default in 0.6.2), clients can
still force the older/fallback redirection behavior (see below) if desired by
setting the `X-Vault-No-Request-Forwarding` header to any non-empty value.

Successful cluster setup requires a few configuration parameters, although some
can be automatically determined.

## Client Redirection

If `X-Vault-No-Request-Forwarding` header in the request is set to a non-empty
value, the standby nodes will redirect the client using a `307` status code to
the *active node's* redirect address.

This is also the fallback method used when request forwarding is turned off or
there is an error performing the forwarding. As such, a redirect address is
always required for all HA setups.

Some HA data store drivers can autodetect the redirect address, but it is often
necessary to configure it manually via a top-level value in the configuration
file. The key for this value is `api_addr` and the value can also be specified
by the `VAULT_API_ADDR` environment variable, which takes precedence.

What the `api_addr` value should be set to depends on how Vault is set up.
There are two common scenarios: Vault servers accessed directly by clients, and
Vault servers accessed via a load balancer.

In both cases, the `api_addr` should be a full URL including scheme
(`http`/`https`), not simply an IP address and port.

### Direct Access

When clients are able to access Vault directly, the `api_addr` for each
node should be that node's address. For instance, if there are two Vault nodes
`A` (accessed via `https://a.vault.mycompany.com:8200`) and `B` (accessed via
`https://b.vault.mycompany.com:8200`), node `A` would set its `api_addr`
to `https://a.vault.mycompany.com:8200` and node `B` would set its
`api_addr` to `https://b.vault.mycompany.com:8200`.

This way, when `A` is the active node, any requests received by node `B` will
cause it to redirect the client to node `A`'s `api_addr` at
`https://a.vault.mycompany.com`, and vice-versa.

### Behind Load Balancers

Sometimes clients use load balancers as an initial method to access one of the
Vault servers, but actually have direct access to each Vault node. In this
case, the Vault servers should actually be set up as described in the above
section, since for redirection purposes the clients have direct access.

However, if the only access to the Vault servers is via the load balancer, the
`api_addr` on each node should be the same: the address of the load
balancer. Clients that reach a standby node will be redirected back to the load
balancer; at that point hopefully the load balancer's configuration will have
been updated to know the address of the current leader. This can cause a
redirect loop and as such is not a recommended setup when it can be avoided.

### Per-Node Cluster Listener Addresses

Each `listener` block in Vault's configuration file contains an `address` value
on which Vault listens for requests. Similarly, each `listener` block can
contain a `cluster_address` on which Vault listens for server-to-server cluster
requests. If this value is not set, its IP address will be automatically set to
same as the `address` value, and its port will be automatically set to the same
as the `address` value plus one (so by default, port `8201`).

Note that *only* active nodes have active listeners. When a node becomes active
it will start cluster listeners, and when it becomes standby it will stop them.

### Per-Node Cluster Address

Similar to the `api_addr`, `cluster_addr` is the value that each node, if
active, should advertise to the standbys to use for server-to-server
communications, and lives as a top-level value in the configuration file. On
each node, this should be set to a host name or IP address that a standby can
use to reach one of that node's `cluster_address` values set in the `listener`
blocks, including port. (Note that this will always be forced to `https` since
only TLS connections are used between servers.)

This value can also be specified by the `VAULT_CLUSTER_ADDR` environment
variable, which takes precedence.

## Storage Support

Currently there are several storage backends that support high availability
mode, including Consul, ZooKeeper and etcd. These may change over time, and the
[configuration page](/docs/configuration/index.html) should be referenced.

The Consul backend is the recommended HA backend, as it is used in production
by HashiCorp and its customers with commercial support.

If you're interested in implementing another backend or adding HA support to
another backend, we'd love your contributions. Adding HA support requires
implementing the `physical.HABackend` interface for the storage backend.
