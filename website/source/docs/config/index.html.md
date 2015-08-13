---
layout: "docs"
page_title: "Server Configuration"
sidebar_current: "docs-config"
description: |-
  Vault server configuration reference.
---

# Server Configuration

Outside of development mode, Vault servers are configured using a file.
The format of this file is [HCL](https://github.com/hashicorp/hcl) or JSON.
An example configuration is shown below:

```javascript
backend "consul" {
  address = "127.0.0.1:8500"
  path = "vault"
}

listener "tcp" {
  address = "127.0.0.1:8200"
  tls_disable = 1
}

telemetry {
  statsite_address = "127.0.0.1:8125"
  disable_hostname = true
}
```

After the configuration is written, use the `-config` flag with `vault server`
to specify where the configuration is.

## Reference

* `backend` (required) - Configures the storage backend where Vault data
  is stored. There are multiple options available for storage backends,
  and they're documented below.

* `listener` (required) - Configures how Vault is listening for API requests.
  "tcp" is currently the only option available. A full reference for the
   inner syntax is below.

* `disable_mlock` (optional) - A boolean. If true, this will disable the
  server from executing the `mlock` syscall to prevent memory from being
  swapped to disk. This is not recommended in production (see below).

* `telemetry` (optional)  - Configures the telemetry reporting system
  (see below).

* `default_lease_duration` (optional) - Configures the default lease
  duration for tokens and secrets, specified in hours. Default value
  is 30 days. This value cannot be larger than `max_lease_duration`.

* `max_lease_duration` (optional) - Configures the maximum possible
  lease duration for tokens and secrets, specified in hours. Default
  value is 30 days.

In production, you should only consider setting the `disable_mlock` option
on Linux systems that only use encrypted swap or do not use swap at all.
Vault does not currently support memory locking on Mac OS X and Windows
and so the feature is automatically disabled on those platforms.  To give
the Vault executable access to the `mlock` syscall on Linux systems:

```shell
sudo setcap cap_ipc_lock=+ep $(readlink -f $(which vault))
```

## Backend Reference

For the `backend` section, the supported backends are shown below.
Vault requires that the backend itself will be responsible for backups,
durability, etc.

  * `consul` - Store data within [Consul](http://www.consul.io). This
      backend supports HA. It is the most recommended backend for Vault
      and has been shown to work at high scale under heavy load.

  * `etcd` - Store data within [etcd](https://coreos.com/etcd/).
      This backend supports HA.

  * `zookeeper` - Store data within [Zookeeper](https://zookeeper.apache.org/).
      This backend supports HA.

  * `s3` - Store data within an S3 bucket [S3](http://aws.amazon.com/s3/).
      This backend does not support HA.

  * `mysql` - Store data within MySQL. This backend does not support HA.

  * `inmem` - Store data in-memory. This is only really useful for
      development and experimentation. Data is lost whenever Vault is
      restarted.

  * `file` - Store data on the filesystem using a directory structure.
      This backend does not support HA.

#### Common Backend Options

All backends support the following options:

  * `advertise_addr` (optional) - For backends that support HA, this
      is the address to advertise to other Vault servers in the cluster
      for request forwarding. Most HA backends will attempt to determine
      the advertise address if not provided.

#### Backend Reference: Consul

For Consul, the following options are supported:

  * `path` (optional) - The path within Consul where data will be stored.
      Defaults to "vault/".

  * `address` (optional) - The address of the Consul agent to talk to.
      Defaults to the local agent address, if available.

  * `scheme` (optional) - "http" or "https" for talking to Consul.

  * `datacenter` (optional) - The datacenter within Consul to write to.
      This defaults to the local datacenter.

  * `token` (optional) - An access token to use to write data to Consul.

  * `tls_skip_verify` (optional) - If non-empty, then TLS host verification
      will be disabled for Consul communication.
      Defaults to false.

  The following settings should be set according to your [Consul encryption settings](https://www.consul.io/docs/agent/encryption.html):

  * `tls_ca_file` (optional) - The path to the CA certificate used for Consul communication.
      Defaults to system bundle if not specified.
      Set accordingly to the [ca_file](https://www.consul.io/docs/agent/options.html#ca_file) setting in Consul.

  * `tls_cert_file` (optional) - The path to the certificate for Consul communication.
      Set accordingly to the [cert_file](https://www.consul.io/docs/agent/options.html#cert_file) setting in Consul.

  * `tls_key_file` (optional) - The path to the private key for Consul communication.
      Set accordingly to the [key_file](https://www.consul.io/docs/agent/options.html#key_file) setting in Consul.

#### Backend Reference: Zookeeper

For Zookeeper, the following options are supported:

  * `path` (optional) - The path within Zookeeper where data will be stored.
      Defaults to "vault/".

  * `address` (optional) - The address(es) of the Zookeeper instance(s) to talk to.
      Can be comma separated list (host:port) of many Zookeeper instances.
      Defaults to "localhost:2181" if not specified.

#### Backend Reference: etcd

For etcd, the following options are supported:

  * `path` (optional) - The path within etcd where data will be stored.
      Defaults to "vault/".

  * `address` (optional) - The address(es) of the etcd instance(s) to talk to.
      Can be comma separated list (protocol://host:port) of many etcd instances.
      Defaults to "http://localhost:4001" if not specified.

#### Backend Reference: S3

For S3, the following options are supported:

  * `bucket` (required) - The name of the S3 bucket to use.

  * `access_key` - (required) The AWS access key. It must be provided, but it can also be sourced from the AWS_ACCESS_KEY_ID environment variable.

  * `secret_key` - (required) The AWS secret key. It must be provided, but it can also be sourced from the AWS_SECRET_ACCESS_KEY environment variable.

  * `session_token` - (optional) The AWS session_token. It can also be sourced from the AWS_SESSION_TOKEN environment variable.

  * `region` (optional) - The AWS region. It can be sourced from the AWS_DEFAULT_REGION environment variable and will default to "us-east-1" if not specified.

#### Backend Reference: MySQL

The MySQL backend has the following options:

  * `username` (required) - The MySQL username to connect with.

  * `password` (required) - The MySQL password to connect with.

  * `address` (optional) - The address of the MySQL host. Defaults to
    "127.0.0.1:3306.

  * `database` (optional) - The name of the database to use. Defaults to "vault".

  * `table` (optional) - The name of the table to use. Defaults to "vault".

  * `tls_ca_file` (optional) - The path to the CA certificate to connect using TLS

#### Backend Reference: Inmem

The in-memory backend has no configuration options.

#### Backend Reference: File

The file backend has the following options:

  * `path` (required) - The path on disk to a directory where the
      data will be stored.

## Listener Reference

For the `listener` section, the only supported listener currently
is "tcp". Regardless of future plans, this is the recommended listener,
since it allows for HA mode.

The supported options are:

  * `address` (optional) - The address to bind to for listening. This
      defaults to "127.0.0.1:8200".

  * `tls_disable` (optional) - If non-empty, then TLS will be disabled.
      This is an opt-in; Vault assumes by default that TLS will be used.

  * `tls_cert_file` (required unless disabled) - The path to the certificate
      for TLS.

  * `tls_key_file` (required unless disabled) - The path to the private key
      for the certificate.

  * `tls_min_version` (optional) - **(Vault > 0.2)** If provided, specifies
      the minimum supported version of TLS. Accepted values are "tls10", "tls11"
      or "tls12". This defaults to "tls12". WARNING: TLS 1.1 and lower
      are generally considered less secure; avoid using these if
      possible.

## Telemetry Reference

For the `telemetry` section, there is no resource name. All configuration
is within the object itself.

* `statsite_address` (optional) - An address to a [Statsite](https://github.com/armon/statsite)
  instances for metrics. This is highly recommended for production usage.

* `statsd_address` (optional) - This is the same as `statsite_address` but
  for StatsD.

* `disable_hostname` (optional) - Whether or not to prepend runtime telemetry
  with the machines hostname. This is a global option. Defaults to false.
