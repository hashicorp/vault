---
layout: "guides"
page_title: "Vault HA with Consul - Guides"
sidebar_current: "guides-operations-vault-ha"
description: |-
  This guide will walk you through a simple Vault Highly Available (HA) cluster
  implementation. While this is not an exhaustive or prescriptive guide that
  can be used as a drop-in production example, it covers the basics enough to
  inform your own production setup.
---

# Vault High Availability (HA)

Vault can run in a high availability (HA) mode to protect against  outages by
running multiple Vault servers. Vault is typically bound by the IO limits of the
storage backend rather than the compute requirements. Certain storage backends,
such as Consul, provide additional coordination functions that enable Vault to
run in an HA configuration while others provide a more robust backup and
restoration process.

When running in HA mode, Vault servers have two additional states: ***standby***
and ***active***. Within a Vault cluster, only a single instance will be
_active_ and handles all requests (reads and writes) and all _standby_ nodes
redirect requests to the _active_ node.

![Reference Architecture](/assets/images/vault-ha-consul-3.png)


~> This guide will walk you through a simple Vault Highly Available (HA) cluster
implementation. While this is not an exhaustive or prescriptive guide that can
be used as a drop-in production example, it covers the **basics** enough to inform
your own production setup.


## Reference Materials

- [High Availability Mode](/docs/concepts/ha.html)
- [Consul Storage Backend](/docs/configuration/storage/consul.html)
- [High Availability Parameters](/docs/configuration/index.html#high-availability-parameters)
- [Consul Agent Configuration](https://www.consul.io/docs/agent/options.html)

## Estimated Time to Complete

25 minutes

## Prerequisites

This intermediate Vault operations guide assumes that you have some working
knowledge of Vault and Consul.


## Steps

Our goal in following this guide is to arrive at a Vault HA setup
consisting of the following:

- 2 Vault servers: 1 active and 1 standby
- Cluster of 3 Consul servers

### Reference Diagram

This diagram lays out the simple architecture details for reference:

![Reference Architecture](/assets/images/vault-ha-consul.png)

You perform the following:

- [Step 1: Setup a Consul Server Cluster](#step1)
- [Step 2: Start and Verify the Consul Cluster State](#step2)
- [Step 3: Setup Consul Client Agents on Vault Nodes](#step3)
- [Step 4: Configure the Vault Servers](#step4)
- [Step 5: Start Vault and Verify the State](#step5)

-> For the purpose of this guide, we will use the open source software editions of
Vault and Consul; however, the setup is the same for Enterprise editions.


### <a name="step1"></a>Step 1: Setup a Consul Server Cluster

Our Consul servers in this guide will be defined by IP address only, but also
referenced by a label:

- **`consul_s1: 10.1.42.101`**
- **`consul_s2: 10.1.42.102`**
- **`consul_s3: 10.1.42.103`**

The [Consul binary](https://www.consul.io/downloads.html) is presumed to be
located in **`/usr/local/bin/consul`**, but if your case differs, you can adjust the
path references accordingly.

With that in mind, here is a basic Consul server configuration starting point:

```plaintext
{
  "server": true,
  "node_name": "$NODE_NAME",
  "datacenter": "dc1",
  "data_dir": "$CONSUL_DATA_PATH",
  "bind_addr": "0.0.0.0",
  "client_addr": "0.0.0.0",
  "advertise_addr": "$ADVERTISE_ADDR",
  "bootstrap_expect": 3,
  "retry_join": ["$JOIN1", "$JOIN2", "$JOIN3"],
  "ui": true,
  "log_level": "DEBUG",
  "enable_syslog": true,
  "acl_enforce_version_8": false
}
```

Note that some values contain variable placeholders while the rest have
reasonable defaults. You should replace the following values in your own Consul
server configuration based on the example:

- **$NODE_NAME** this is a unique label for the node; in our case, this will be
`consul_s1`, `consul_s2`, and `consul_s3` respectively.
- **$CONSUL_DATA_PATH**: absolute path to Consul data directory; ensure that
this directory is writable by the Consul process user.
- **$ADVERTISE_ADDR**: set to address that you prefer the Consul
servers advertise to the other servers in the cluster and should not be set to
`0.0.0.0`; for this guide, it should be set to the Consul server’s IP address in
each instance of the configuration file, or `10.1.42.101`,` 10.1.42.102`, and
`10.1.42.103` respectively.
- **$JOIN1**, **$JOIN2**, **$JOIN3**: This example uses the `retry_join`
method of joining the server agents to form a cluster; as such, the values for
this guide would be `10.1.42.101`, `10.1.42.102`, and `10.1.42.103` respectively.


Note that the web user interface is enabled (`"ui": true`), and Consul will be
logging at DEBBUG level to the system log (`"log_level": "DEBUG"`). For the
purpose of this guide, the **`acl_enforce_version_8`** is set to `false` so that
we do not need to be concerned with ACLs in this guide. However, you would want
to enable ACLs in a production environment and follow the [Consul ACL Guide](https://www.consul.io/docs/guides/acl.html#acl-agent-master-token) for
details.

Create a configuration file for each Vault server and save it as
**`/usr/local/etc/consul/client_agent.json`**.

#### `consul_s1.hcl` Example

    {
      "server": true,
      "node_name": "consul_s1",
      "datacenter": "dc1",
      "data_dir": "/var/consul/data",
      "bind_addr": "0.0.0.0",
      "client_addr": "0.0.0.0",
      "advertise_addr": "10.1.42.101",
      "bootstrap_expect": 3,
      "retry_join": ["10.1.42.101", "10.1.42.102", "10.1.42.103"],
      "ui": true,
      "log_level": "DEBUG",
      "enable_syslog": true,
      "acl_enforce_version_8": false
    }

#### `consul_s2.hcl` Example

    {
      "server": true,
      "node_name": "consul_s2",
      "datacenter": "dc1",
      "data_dir": "/var/consul/data",
      "bind_addr": "0.0.0.0",
      "client_addr": "0.0.0.0",
      "advertise_addr": "10.1.42.102",
      "bootstrap_expect": 3,
      "retry_join": ["10.1.42.101", "10.1.42.102", "10.1.42.103"],
      "ui": true,
      "log_level": "DEBUG",
      "enable_syslog": true,
      "acl_enforce_version_8": false
    }

#### `consul_s3.hcl` Example

    {
      "server": true,
      "node_name": "consul_s3",
      "datacenter": "dc1",
      "data_dir": "/var/consul/data",
      "bind_addr": "0.0.0.0",
      "client_addr": "0.0.0.0",
      "advertise_addr": "10.1.42.103",
      "bootstrap_expect": 3,
      "retry_join": ["10.1.42.101", "10.1.42.102", "10.1.42.103"],
      "ui": true,
      "log_level": "DEBUG",
      "enable_syslog": true,
      "acl_enforce_version_8": false
    }

#### Consul Server `systemd` Unit file

You have Consul binaries and a reasonably basic configuration and now you just
need to start Consul on each server instance; `systemd` is popular in most
contemporary Linux distributions, so with that in mind, here is an example
`systemd` unit file:

```plaintext
### BEGIN INIT INFO
# Provides:          consul
# Required-Start:    $local_fs $remote_fs
# Required-Stop:     $local_fs $remote_fs
# Default-Start:     2 3 4 5
# Default-Stop:      0 1 6
# Short-Description: Consul agent
# Description:       Consul service discovery framework
### END INIT INFO

[Unit]
Description=Consul server agent
Requires=network-online.target
After=network-online.target

[Service]
User=consul
Group=consul
PIDFile=/var/run/consul/consul.pid
PermissionsStartOnly=true
ExecStartPre=-/bin/mkdir -p /var/run/consul
ExecStartPre=/bin/chown -R consul:consul /var/run/consul
ExecStart=/usr/local/bin/consul agent \
    -config-file=/usr/local/etc/consul/server_agent.json \
    -pid-file=/var/run/consul/consul.pid
ExecReload=/bin/kill -HUP $MAINPID
KillMode=process
KillSignal=SIGTERM
Restart=on-failure
RestartSec=42s

[Install]
WantedBy=multi-user.target
```

Note that you might be interested in changing the values of the following
depending on style, file system hierarchy standard adherence level, and so on:

- **`-config-file`**
- **`-pid-file`**

Once the unit file is defined and saved (e.g.
`/etc/systemd/system/consul.service`), be sure to perform a `systemctl daemon-reload`
and then you can start your Consul service on each server.

### <a name="step2"></a>Step 2: Start and Verify the Consul Cluster State

Be sure that the ownership and permissions are correct on the directory you
specified for the value of `data_dir`, and then start the Consul service on each
system and verify the status:

    $ sudo systemctl start consul
    $ sudo systemctl status consul
    ● consul.service - Consul server agent
       Loaded: loaded (/etc/systemd/system/consul.service; enabled; vendor preset: enabled)
       Active: active (running) since Mon 2018-03-19 17:33:14 UTC; 24h ago
     Main PID: 2068 (consul)
        Tasks: 13
       Memory: 13.6M
          CPU: 0m 52.784s
       CGroup: /system.slice/consul.service
               └─2068 /usr/local/bin/consul agent -config-file=/usr/local/etc/consul/server_agent.json -pid-file=/var/run/consul/consul.pid

After starting all Consul server agents, let’s check the Consul cluster status:

    $consul members
    Node       Address           Status  Type    Build  Protocol  DC    Segment
    consul_s1  10.1.42.101:8301  alive   server  1.0.6  2         dc1   <all>
    consul_s2  10.1.42.102:8301  alive   server  1.0.6  2         dc1   <all>
    consul_s3  10.1.42.103:8301  alive   server  1.0.6  2         dc1   <all>

The cluster looks good and all 3 servers are shown; let’s make sure we have a
leader before proceeding:

    $consul operator raft list-peers
    Node                   ID                                    Address           State     Voter  RaftProtocol
    consul_s2              536b721f-645d-544a-c10d-85c2ca24e4e4  10.1.42.102:8300  follower  true   3
    consul_s1              e10ba554-a4f9-6a8c-f662-81c8bb2a04f5  10.1.42.101:8300  follower  true   3
    consul_s3              56370ec8-da25-e7dc-dfc6-bf5f27978a7a  10.1.42.103:8300  leader    true   3

The above output shows that **`consul_s3`** is the current cluster leader in
this example.  Now, you are good to move on to the Vault server configuration.


### <a name="step3"></a>Step 3: Setup Consul Client Agents on Vault Nodes

The Vault server nodes require **both** the Consul and Vault binaries on each node. Consul will be configured as a **client** agent and Vault will be configured as a server.

![Reference Architecture](/assets/images/vault-ha-consul-2.png)


#### Consul Client Agent Configuration

Since Consul is used to provide a highly available storage backend, you need to
configure local Consul client agents on the Vault servers which will communicate
with the Consul server cluster for registering health checks, service discovery, and cluster HA failover coordination (cluster leadership).

~> Note that [it is not recommended to connect the Vault servers directly to the
Consul servers](/docs/configuration/storage/consul.html#address).

The Consul client agents will be using the same address as the Vault servers for
network communication to the Consul server cluster, but they will be binding the
**`client_address`** only to the loopback interface such that Vault can connect to
it over the loopback interface.

Here is the example configuration for the Consul client agent:

    {
      "server": false,
      "datacenter": "dc1",
      "node_name": "$NODE_NAME",
      "data_dir": "$CONSUL_DATA_PATH",
      "bind_addr": "$BIND_ADDR",
      "client_addr": "127.0.0.1",
      "retry_join": ["$JOIN1", "$JOIN2", "$JOIN3"],
      "log_level": "DEBUG",
      "enable_syslog": true,
      "acl_enforce_version_8": false
    }

Similar to what you have done in [Step 1](#step1), replace the following values
in your own Consul client agent configuration accordingly:

- **$NODE_NAME** this is a unique label for the node; in
our case, this will be `consul_c1` and `consul_c2` respectively.
- **$CONSUL_DATA_PATH**: absolute path to Consul data directory; ensure that this
directory is writable by the Consul process user.
- **$BIND_ADDR**: this should be set
to address that you prefer the Consul servers advertise to the other servers in
the cluster and should not be set to `0.0.0.0`; for this guide, it should be set
to the Vault server’s IP address in each instance of the configuration file, or
`10.1.42.201` and `10.1.42.202` respectively.
- **$JOIN1**, **$JOIN2**, **$JOIN3**: This example uses the `retry_join` method of
joining the server agents to form a cluster; as such, the values for this guide
would be `10.1.42.101`, `10.1.42.102`, and `10.1.42.103` respectively.

Create a configuration file for each Vault server and save it as
**`/usr/local/etc/consul/client_agent.json`**.

#### `consul_c1.hcl` Example

    {
      "server": false,
      "datacenter": "dc1",
      "node_name": "consul_c1",
      "data_dir": "/var/consul/data",
      "bind_addr": "10.1.42.201",
      "client_addr": "127.0.0.1",
      "retry_join": ["10.1.42.101", "10.1.42.102", "10.1.42.103"],
      "log_level": "DEBUG",
      "enable_syslog": true,
      "acl_enforce_version_8": false
    }

#### `consul_c2.hcl` Example

    {
      "server": false,
      "datacenter": "dc1",
      "node_name": "consul_c2",
      "data_dir": "/var/consul/data",
      "bind_addr": "10.1.42.202",
      "client_addr": "127.0.0.1",
      "retry_join": ["10.1.42.101", "10.1.42.102", "10.1.42.103"],
      "log_level": "DEBUG",
      "enable_syslog": true,
      "acl_enforce_version_8": false
    }

#### Consul Server `systemd` Unit file

You have Consul binaries and a reasonably basic client agent configuration and
now you just need to start Consul on each of the Vault server instances. Here
is an example `systemd` unit file:

    ### BEGIN INIT INFO
    # Provides:          consul
    # Required-Start:    $local_fs $remote_fs
    # Required-Stop:     $local_fs $remote_fs
    # Default-Start:     2 3 4 5
    # Default-Stop:      0 1 6
    # Short-Description: Consul agent
    # Description:       Consul service discovery framework
    ### END INIT INFO

    [Unit]
    Description=Consul client agent
    Requires=network-online.target
    After=network-online.target

    [Service]
    User=consul
    Group=consul
    PIDFile=/var/run/consul/consul.pid
    PermissionsStartOnly=true
    ExecStartPre=-/bin/mkdir -p /var/run/consul
    ExecStartPre=/bin/chown -R consul:consul /var/run/consul
    ExecStart=/usr/local/bin/consul agent \
        -config-file=/usr/local/etc/consul/client_agent.json \
        -pid-file=/var/run/consul/consul.pid
    ExecReload=/bin/kill -HUP $MAINPID
    KillMode=process
    KillSignal=SIGTERM
    Restart=on-failure
    RestartSec=42s

    [Install]
    WantedBy=multi-user.target

Change the following values as necessary:

- **`-config-file`**
- **`-pid-file`**

Once the unit file is defined and saved (e.g.
`/etc/systemd/system/consul.service`), be sure to perform a `systemctl
daemon-reload` and then you can start your Consul service on each Vault server.

Start the Consul and verify its cluster state to be sure that the ownership and
permissions are correct on the directory you specified for the value of
`data_dir`, and then start the Consul service on each system and verify the
status:

    $ sudo systemctl start consul
    $ sudo systemctl status consul
    ● consul.service - Consul client agent
       Loaded: loaded (/etc/systemd/system/consul.service; enabled; vendor preset: enabled)
       Active: active (running) since Tue 2018-03-20 19:36:49 UTC; 6s ago
     Main PID: 23758 (consul)
        Tasks: 11
       Memory: 9.8M
          CPU: 571ms
       CGroup: /system.slice/consul.service
               └─23758 /usr/local/bin/consul agent -config-file=/usr/local/etc/consul/client_agent.json -pid-file=/var/run/consul/consul.pid

After starting all Consul client agents, check the Consul cluster status:

    $consul members
    Node        Address           Status  Type    Build  Protocol  DC    Segment
    consul_s1   10.1.42.101:8301  alive   server  1.0.6  2         dc1   <all>
    consul_s2   10.1.42.102:8301  alive   server  1.0.6  2         dc1   <all>
    consul_s3   10.1.42.103:8301  alive   server  1.0.6  2         dc1   <all>
    consul_c1   10.1.42.201:8301  alive   client  1.0.6  2         arus  <default>
    consul_c2   10.1.42.202:8301  alive   client  1.0.6  2         arus  <default>

The above output shows 3 Consul server agents and 2 Consul client agents in the
cluster.  Now, you are ready to configure the Vault servers.

### <a name="step4"></a>Step 4: Configure the Vault Servers

Now that we have a Consul cluster consisting of 3 servers and 2 client agents
for our Vault servers, let’s get the configuration for Vault and a startup
script together so that we can bootstrap the Vault HA setup.

Our Vault servers in this guide are defined by IP address only, but referenced
by a label as well:

- **`vault_s1: 10.1.42.201`**
- **`vault_s2: 10.1.42.202`**

In our configuration file, we'll set up the following:

- [**`tcp`**](/docs/configuration/listener/tcp.html) listener
- [**`consul`**](/docs/configuration/storage/consul.html)  storage backend
- [High Availability parameters](/docs/configuration/index.html#high-availability-parameters)

This section assumes the Vault binary is located at **`/usr/local/bin/vault`**

#### Vault Configuration

    listener "tcp" {
      address          = "0.0.0.0:8200"
      cluster_address  = "0.0.0.0:8201"
      tls_disable      = "true"
    }

    storage "consul" {
      address = "127.0.0.1:8500"
      path    = "vault/"
    }

    api_addr =  "$API_ADDR"
    cluster_addr = "$CLUSTER_ADDR"


We're setting the following parameters for our `tcp` listener:

- `address` (string: "127.0.0.1:8200") – Specifies the address to bind to for listening.
- `cluster_address` (string: "127.0.0.1:8201") – Specifies the address to bind to for cluster server-to-server requests. This defaults to one port higher than the value of address. This does not usually need to be set, but can be useful in case Vault servers are isolated from each other in such a way that they need to hop through a TCP load balancer or some other scheme in order to talk.

This configuration allows for listening on all interfaces (such that a Vault
command against the loopback address would succeed, for example). Specifical

We're also explicitly setting Vault's [HA parameters](/docs/configuration/index.html#high-availability-parameters) (`api_addr` and `cluster_addr`). Often, it's not necessary to configure these two parameters when using Consul as Vault's storage backend, as Consul will attempt to automatically discover and advertise the address of the active Vault node. However, certain cluster configurations might require them to be explicitly set (accesing Vault through a load balancer, for example).

For the sake of simplicity, we will assume that clients in our scenario connect directly to the Vault nodes (rather than through a load balancer). Review the [Client Redirection](/docs/concepts/ha.html#client-redirection) documentation for more information on client access patterns and their implications.

Note that some values contain variable placeholders while the rest have
reasonable defaults. You should replace the following values in your own Vault
server configuration based on the example:

- **$API_ADDR**: Specifies the address (full URL) to advertise to other Vault servers in the cluster for client redirection. This can also be provided via the environment variable `VAULT_API_ADDR`. In general this should be set to a full URL that points to the value of the listener address. In our scenario, it will be `http://10.1.42.201:8200`
and `http://10.1.42.202:8200` respectively.

- **$CLUSTER_ADDR**: Specifies the address to advertise to other Vault servers in the cluster for request forwarding. This can also be provided via the environment variable `VAULT_CLUSTER_ADDR`. This is a full URL, like `api_addr`. In our scenario, it will be `https://10.1.42.201:8201` and
`https://10.1.42.202:8201` respectively.

> Note that the scheme here (https) is ignored; all cluster members will always
use TLS with a private key/certificate.


#### `vault_s1.hcl` Example

    listener "tcp" {
      address          = "0.0.0.0:8200"
      cluster_address  = "10.1.42.201:8201"
      tls_disable      = "true"
    }

    storage "consul" {
      address = "127.0.0.1:8500"
      path    = "vault/"
    }

    api_addr = "http://10.1.42.201:8200"
    cluster_addr = "https://10.1.42.201:8201"


#### `vault_s2.hcl` Example

    listener "tcp" {
      address          = "0.0.0.0:8200"
      cluster_address  = "10.1.42.202:8201"
      tls_disable      = "true"
    }

    storage "consul" {
      address = "127.0.0.1:8500"
      path    = "vault/"
    }

    api_addr = "http://10.1.42.202:8200"
    cluster_addr = "https://10.1.42.202:8201"


#### Vault Server `systemd` Unit file

You have Vault binaries and a reasonably basic configuration along with local
client agents configured.  Now, you just need to start Vault on each server
instance. Here is an example `systemd` unit file:

    ### BEGIN INIT INFO
    # Provides:          vault
    # Required-Start:    $local_fs $remote_fs
    # Required-Stop:     $local_fs $remote_fs
    # Default-Start:     2 3 4 5
    # Default-Stop:      0 1 6
    # Short-Description: Vault server
    # Description:       Vault secret management tool
    ### END INIT INFO

    [Unit]
    Description=Vault secret management tool
    Requires=network-online.target
    After=network-online.target

    [Service]
    User=vault
    Group=vault
    PIDFile=/var/run/vault/vault.pid
    ExecStart=/usr/local/bin/vault server -config=/etc/vault/vault_server.hcl -log-level=debug
    ExecReload=/bin/kill -HUP $MAINPID
    KillMode=process
    KillSignal=SIGTERM
    Restart=on-failure
    RestartSec=42s
    LimitMEMLOCK=infinity

    [Install]
    WantedBy=multi-user.target

Note that you might be interested in changing the values of the following
depending on style, file system hierarchy standard adherence level, and so on:

- **`-config`**
- **`-log-level`**

Once the unit file is defined and saved as e.g.
`/etc/systemd/system/vault.service`, be sure to perform a `systemctl daemon-reload`
and then you can start your Vault service on each server.

### <a name="step5"></a>Step 5: Start Vault and Verify the State

Start the Vault service on each system and verify the status:

    $ sudo systemctl start vault
    $ sudo systemctl status vault
    ● vault.service - Vault secret management tool
       Loaded: loaded (/etc/systemd/system/vault.service; enabled; vendor preset: enabled)
       Active: active (running) since Tue 2018-03-20 20:42:10 UTC; 42s ago
     Main PID: 2080 (vault)
        Tasks: 12
       Memory: 71.7M
          CPU: 50s
       CGroup: /system.slice/vault.service
               └─2080 /usr/local/bin/vault server -config=/home/ubuntu/vault_nano/config/vault_server.hcl -log-level=debu

Now you’ll need to move on to [initializing and
unsealing](/intro/getting-started/deploy.html#initializing-the-vault) each Vault
instance.

Once that is done, check Vault status on each of the servers.

The **active** Vault server:

    $ vault status
    Key             Value
    ---             -----
    Seal Type       shamir
    Sealed          false
    Total Shares    5
    Threshold       3
    Version         0.9.5
    Cluster Name    vault
    Cluster ID      0ee91bd1-55ec-c84f-3c1d-dcc7f4f644a8
    HA Enabled      true
    HA Cluster      https://10.1.42.201:8201
    HA Mode         active

The **standby** Vault server:

    vault status
    Key                     Value
    ---                     -----
    Seal Type               shamir
    Sealed                  false
    Total Shares            5
    Threshold               3
    Version                 0.9.5
    Cluster Name            vaultron
    Cluster ID              0ee91bd1-55ec-c84f-3c1d-dcc7f4f644a8
    HA Enabled              true
    HA Cluster              https://10.1.42.201:8201
    HA Mode                 standby
    Active Node Address:    http://10.1.42.201:8200


Vault servers are now operational in HA mode at this point, and you should be
able to write a secret from either the active or the standby Vault instance and
see it succeed as a test of request forwarding. Also, you can shut down the
active instance (`sudo systemctl stop vault`) to simulate a system failure and
see the standby instance assumes the leadership.


## Next steps

Read [Production Hardening](/guides/operations/production.html) to learn best
practices for a production hardening deployment of Vault.
