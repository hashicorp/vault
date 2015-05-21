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
  address = "demo.consul.io:80"
  path = "vault"
}

listener "tcp" {
  address = "127.0.0.1:8200"
  tls_disable = 1
}
```

After the configuration is written, use the `-config` flag with `vault server`
to specify where the configuration is.

## Reference

* `backend` (required) - Configures the physical backend where Vault data
  is stored. There are multiple options available for physical backends,
  and they're documented below.

* `listener` (required) - Configures how Vault is listening for API requests.
  "tcp" is currently the only option available. A full reference for the
   inner syntax is below.

* `disable_mlock` (optional) - A boolean. If true, this will disable the
  server from executing the `mlock` syscall to prevent memory from being
  swapped to disk. This is not recommended.

* `statsite_addr` (optional) - An address to a [Statsite](https://github.com/armon/statsite)
  instances for metrics. This is highly recommended for production usage.

* `statsd_addr` (optional) - This is the same as `statsite_addr` but
  for StatsD.

## Backend Reference

For the `backend` section, the supported backends are shown below.
Vault requires that the backend itself will be responsible for backups,
durability, etc.

  * `consul` - Store data within [Consul](http://www.consul.io). This
      backend supports HA. It is the most recommended backend for Vault
      and has been shown to work at high scale under heavy load.

  * `zookeeper` - Store data within [Zookeeper](https://zookeeper.apache.org/).
      This backend does not support HA.

  * `s3` - Store data within an S3 bucket [S3](http://aws.amazon.com/s3/).
      This backend does not support HA.

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

#### Backend Reference: Zookeeper

For Zookeeper, the following options are supported:

  * `path` (optional) - The path within Zookeeper where data will be stored.
      Defaults to "vault/".

  * `address` (optional) - The address(es) of the Zookeeper instance(s) to talk to.
      Can be comma separated list (host:port) of many Zookeeper instances.
      Defaults to "localhost:2181" if not specified.

#### Backend Reference: S3

For S3, the following options are supported:

  * `bucket` (required) - The name of the S3 bucket to use.

  * `access_key` - (Required) This is the AWS access key. It must be provided, but it can also be sourced from the AWS_ACCESS_KEY_ID environment variable.

  * `secret_key` - (Required) This is the AWS secret key. It must be provided, but it can also be sourced from the AWS_SECRET_ACCESS_KEY environment variable.

  * `region` (optional) - This is the AWS region. It can be sourced from the AWS_DEFAULT_REGION environment variable and will default to "us-east-1" if not specified.

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
