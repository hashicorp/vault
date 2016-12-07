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

Starting with 0.5.2, limited configuration options can be changed on-the-fly by
sending a SIGHUP to the server process. These are denoted below.

## Reference

* `backend` (required) - Configures the storage backend where Vault data
  is stored. There are multiple options available for storage backends,
  and they're documented below.

* `ha_backend` (optional) - Configures the storage backend where Vault HA
  coordination will take place. Must be an HA-supporting backend using the
  configuration options as documented below. If not set, HA will be attempted
  on the backend given in the `backend` parameter.

* `cluster_name` (optional) - An identifier for your Vault cluster. If omitted,
  Vault will generate a value for `cluster_name`. If connecting to Vault
  Enterprise, this value will be used in the interface.

* `listener` (required) - Configures how Vault is listening for API requests.
  "tcp" and "atlas" are valid values. A full reference for the
   inner syntax is below.

* `cache_size` (optional) - If set, the size of the read cache used
  by the physical storage subsystem will be set to this value. The
  value is in number of entries so the total cache size is dependent
  on the entries being stored. Defaults to 32k entries.

* `disable_cache` (optional) - A boolean. If true, this will disable all caches
  within Vault, including the read cache used by the physical storage
  subsystem. This will very significantly impact performance.

* `disable_mlock` (optional) - A boolean. If true, this will disable the
  server from executing the `mlock` syscall to prevent memory from being
  swapped to disk. This is not recommended in production (see below).

* `telemetry` (optional)  - Configures the telemetry reporting system
  (see below).

* `default_lease_ttl` (optional) - Configures the default lease duration
  for tokens and secrets. This is a string value using a suffix, e.g. "768h".
  Default value is 32 days. This value cannot be larger than `max_lease_ttl`.

* `max_lease_ttl` (optional) - Configures the maximum possible
  lease duration for tokens and secrets. This is a string value using a suffix,
  e.g. "768h". Default value is 32 days.

* `ui` (optional, Vault Enterprise only) - If set `true`, enables the built-in
  web-based UI. Once enabled, the UI will be available to browsers at the
  standard Vault address.

In production it is a risk to run Vault on systems where `mlock` is
unavailable or the setting has been disabled via the `disable_mlock`.
Disabling `mlock` is not recommended unless the systems running Vault only
use encrypted swap or do not use swap at all.  Vault only supports memory
locking on UNIX-like systems (Linux, FreeBSD, Darwin, etc).  Non-UNIX like
systems (e.g. Windows, NaCL, Android) lack the primitives to keep a process's
entire memory address space from spilling to disk and is therefore automatically
disabled on unsupported platforms.

On Linux, to give the Vault executable the ability to use the `mlock` syscall
without running the process as root, run:

```shell
sudo setcap cap_ipc_lock=+ep $(readlink -f $(which vault))
```

## Listener Reference

For the `listener` section, the only required listener is "tcp".
Regardless of future plans, this is the recommended listener,
as it allows for HA mode. If you wish to use the Vault
Enterprise interface in HashiCorp Atlas, you may add an ["atlas" listener block](#connecting-to-vault-enterprise-in-hashicorp-atlas)
in addition to the "tcp" one.

The supported options are:

  * `address` (optional) - The address to bind to for listening. This
      defaults to "127.0.0.1:8200".

  * `cluster_address` (optional) - The address to bind to for cluster
      server-to-server requests. This defaults to one port higher than the
      value of `address`, so with the default value of `address`, this would be
      "127.0.0.1:8201".

  * `tls_disable` (optional) - If true, then TLS will be disabled.
      This will parse as boolean value, and can be set to "0", "no",
      "false", "1", "yes", or "true". This is an opt-in; Vault assumes
      by default that TLS will be used.

  * `tls_cert_file` (required unless disabled) - The path to the certificate
      for TLS. To configure the listener to use a CA certificate, concatenate
      the primary certificate and the CA certificate together. The primary
      certificate should appear first in the combined file. This is reloaded
      via SIGHUP.

  * `tls_key_file` (required unless disabled) - The path to the private key
      for the certificate. This is reloaded via SIGHUP.

  * `tls_min_version` (optional) - **(Vault > 0.2)** If provided, specifies
      the minimum supported version of TLS. Accepted values are "tls10", "tls11"
      or "tls12". This defaults to "tls12". WARNING: TLS 1.1 and lower
      are generally considered less secure; avoid using these if
      possible.

### Connecting to Vault Enterprise in HashiCorp Atlas

Adding an "atlas" block will initiate a long-running connection to the
[SCADA](https://scada.hashicorp.com) service. The SCADA connection allows the
Vault Enterprise interface to securely communicate with and operate on your
Vault cluster.

The "atlas" `listener` supports these options:

  * `endpoint` (optional) - The endpoint address used for Vault Enterprise interface
      integration. Defaults to the public Vault Enterprise endpoints on Atlas.

  * `infrastructure` (required) - Used to provide the Atlas infrastructure name and
      the SCADA connection. The format of this is `username/environment`.

  * `node_id` (required) - The identifier for an individual node—used in
      the Vault Enterprise dashboard.

  * `token` (required) - A token from Atlas used to authenticate SCADA session. Generate
      one in the [Atlas](https://atlas.hashicorp.com/settings/tokens).

Additionally, the [`cluster_name`](#cluster_name) config option will be used to
identify your cluster members inside the infrastructure in the Vault Enterprise
interface. It is important for operators to use the same value for
`cluster_name` across cluster members because Vault overwrites this value
internally on instance instantiation.

This allows the connection of multiple clusters to a single `infrastructure`.

For more on Vault Enterprise, see the [help documentation](https://atlas.hashicorptest.com/help/vault/features).


## Telemetry Reference

For the `telemetry` section, there is no resource name. All configuration
is within the object itself.

* `statsite_address` (optional) - An address to a [Statsite](https://github.com/armon/statsite)
  instance for metrics. This is highly recommended for production usage.

* `statsd_address` (optional) - This is the same as `statsite_address` but
  for StatsD.

* `disable_hostname` (optional) - Whether or not to prepend runtime telemetry
  with the machines hostname. This is a global option. Defaults to false.

* `circonus_api_token`
  A valid [Circonus](http://circonus.com/) API Token used to create/manage check. If provided, metric management is enabled.

* `circonus_api_app`
  A valid app name associated with the API token. By default, this is set to "consul".

* `circonus_api_url`
  The base URL to use for contacting the Circonus API. By default, this is set to "https://api.circonus.com/v2".

* `circonus_submission_interval`
  The interval at which metrics are submitted to Circonus. By default, this is set to "10s" (ten seconds).

* `circonus_submission_url`
  The `check.config.submission_url` field, of a Check API object, from a previously created HTTPTRAP check.

* `circonus_check_id`
  The Check ID (not **check bundle**) from a previously created HTTPTRAP check. The numeric portion of the `check._cid` field in the Check API object.

* `circonus_check_force_metric_activation`
  Force activation of metrics which already exist and are not currently active. If check management is enabled, the default behavior is to add new metrics as they are encountered. If the metric already exists in the check, it will **not** be activated. This setting overrides that behavior. By default, this is set to "false".

* `circonus_check_instance_id`
  Serves to uniquely identify the metrics coming from this *instance*.  It can be used to maintain metric continuity with transient or ephemeral instances as they move around within an infrastructure. By default, this is set to hostname:application name (e.g. "host123:vault").

* `circonus_check_search_tag`
  A special tag which, when coupled with the instance id, helps to narrow down the search results when neither a Submission URL or Check ID is provided. By default, this is set to service:app (e.g. "service:vault").

* `circonus_check_display_name`
  Specifies a name to give a check when it is created. This name is displayed in the Circonus UI Checks list.

* `circonus_check_tags`
  Comma separated list of additional tags to add to a check when it is created.

* `circonus_broker_id`
  The ID of a specific Circonus Broker to use when creating a new check. The numeric portion of `broker._cid` field in a Broker API object. If metric management is enabled and neither a Submission URL nor Check ID is provided, an attempt will be made to search for an existing check using Instance ID and Search Tag. If one is not found, a new HTTPTRAP check will be created. By default, this is not used and a random Enterprise Broker is selected, or, the default Circonus Public Broker.

* `circonus_broker_select_tag`
  A special tag which will be used to select a Circonus Broker when a Broker ID is not provided. The best use of this is to as a hint for which broker should be used based on *where* this particular instance is running (e.g. a specific geo location or datacenter, dc:sfo). By default, this is not used.

## Backend Reference

For the `backend` section, the supported physical backends are shown below.
Vault requires that the backend itself will be responsible for backups,
durability, etc.

__*Please note*__: The only physical backends actively maintained by HashiCorp
are `consul`, `inmem`, and `file`. The other backends are community-derived and
community-supported. We include them in the hope that they will be useful to
those users that wish to utilize them, but they receive minimal validation and
testing from HashiCorp, and HashiCorp staff may not be knowledgeable about the
data store being utilized. If you encounter problems with them, we will attempt
to help you, but may refer you to the backend author.

  * `consul` - Store data within [Consul](https://www.consul.io). This
    backend supports HA. It is the most recommended backend for Vault and has
    been shown to work at high scale under heavy load.

  * `etcd` - Store data within [etcd](https://coreos.com/etcd/).
    This backend supports HA. This is a community-supported backend.

  * `zookeeper` - Store data within [Zookeeper](https://zookeeper.apache.org/).
    This backend supports HA. This is a community-supported backend.

  * `dynamodb` - Store data in a [DynamoDB](https://aws.amazon.com/dynamodb/) table.
    This backend optionally supports HA. This is a community-supported backend.

  * `s3` - Store data within an S3 bucket [S3](https://aws.amazon.com/s3/).
    This backend does not support HA. This is a community-supported backend.

  * `gcs` - Store data within a [Google Cloud Storage](https://cloud.google.com/storage/) bucket.
    This backend does not support HA. This is a community-supported backend.

  * `azure` - Store data in an Azure Storage container [Azure](https://azure.microsoft.com/en-us/services/storage/).
    This backend does not support HA. This is a community-supported backend.

  * `swift` - Store data within an OpenStack Swift container [Swift](http://docs.openstack.org/developer/swift/).
    This backend does not support HA. This is a community-supported backend.

  * `mysql` - Store data within MySQL. This backend does not support HA. This
    is a community-supported backend.

  * `postgresql` - Store data within PostgreSQL. This backend does not support HA. This
    is a community-supported backend.

  * `inmem` - Store data in-memory. This is only really useful for
    development and experimentation. Data is lost whenever Vault is
    restarted.

  * `file` - Store data on the filesystem using a directory structure.
    This backend does not support HA.


#### High Availability Options

All HA backends support the following options. These are discussed in much more
detail in the [High Availability concepts
page](https://www.vaultproject.io/docs/concepts/ha.html).

  * `redirect_addr` (optional) - This is the address to advertise to other
    Vault servers in the cluster for client redirection. This can also be
    set via the `VAULT_REDIRECT_ADDR` environment variable, which takes
    precedence.

  * `cluster_addr` (optional) - This is the address to advertise to other Vault
    servers in the cluster for request forwarding. This can also be set via the
    `VAULT_CLUSTER_ADDR` environment variable, which takes precedence.

  * `disable_clustering` (optional) - This controls whether clustering features
    (currently, request forwarding) are enabled. Setting this on a node will
    disable these features _when that node is the active node_.

#### Backend Reference: Consul

For Consul, the following options are supported:

  * `path` (optional) - The path within Consul where data will be stored.
    Defaults to "vault/".

  * `address` (optional) - The address of the Consul agent to talk to.
    Defaults to the local agent address, if available.

  * `scheme` (optional) - "http" or "https" for talking to Consul.

  * `check_timeout` (optional) - The check interval used to send health check
    information to Consul.  Defaults to "5s".

  * `disable_registration` (optional) - If true, then Vault will not register
    itself with Consul.  Defaults to "false".

  * `service` (optional) - The name of the service to register with Consul.
    Defaults to "vault".

  * `service_tags` (optional) - Comma separated list of tags that are to be
    applied to the service that gets registered with Consul.

  * `token` (optional) - An access token to use to write data to Consul.

  * `max_parallel` (optional) - The maximum number of concurrent requests to Consul.
    Defaults to `"128"`.

  * `tls_skip_verify` (optional) - If non-empty, then TLS host verification
    will be disabled for Consul communication.  Defaults to false.

  * `tls_min_version` (optional) - Minimum TLS version to use. Accepted values
    are 'tls10', 'tls11' or 'tls12'. Defaults to 'tls12'.

The following settings should be set according to your [Consul encryption
settings](https://www.consul.io/docs/agent/encryption.html):

  * `tls_ca_file` (optional) - The path to the CA certificate used for Consul
    communication.  Defaults to system bundle if not specified.  Set
    accordingly to the
    [ca_file](https://www.consul.io/docs/agent/options.html#ca_file) setting in
    Consul.

  * `tls_cert_file` (optional) - The path to the certificate for Consul
    communication.  Set accordingly to the
    [cert_file](https://www.consul.io/docs/agent/options.html#cert_file)
    setting in Consul.

  * `tls_key_file` (optional) - The path to the private key for Consul
    communication.  Set accordingly to the
    [key_file](https://www.consul.io/docs/agent/options.html#key_file) setting
    in Consul.

```
// Sample Consul Backend configuration with local Consul Agent
backend "consul" {
  // address MUST match Consul's `addresses.http` config value (or
  // `addresses.https` depending on the scheme provided below).
  address = "127.0.0.1:8500"
  #address = "unix:///tmp/.consul.http.sock"

  // scheme defaults to "http" (suitable for loopback and UNIX sockets), but
  // should be "https" when Consul exists on a remote node (a non-standard
  // deployment).  All decryption happen within Vault so this value does not
  // change Vault's Threat Model.
  scheme = "http"

  // token is a Consul ACL Token that has write privileges to the path
  // specified below.  Use of a Consul ACL Token is a best pracitce.
  token = "[redacted]" // Vault's Consul ACL Token

  // path must be writable by the Consul ACL Token
  path = "vault/"
}
```

Once properly configured, an unsealed Vault installation should be available
on the network at `active.vault.service.consul`. Unsealed Vault instances in
the standby state are available at `standby.vault.service.consul`.  All
unsealed Vault instances are available as healthy in the
`vault.service.consul` pool.  Sealed Vault instances will mark themselves as
critical to avoid showing up by default in Consul's service discovery.

```
% dig active.vault.service.consul srv
; <<>> DiG 9.8.3-P1 <<>> active.vault.service.consul srv
; (1 server found)
;; global options: +cmd
;; Got answer:
;; ->>HEADER<<- opcode: QUERY, status: NOERROR, id: 11331
;; flags: qr aa rd; QUERY: 1, ANSWER: 1, AUTHORITY: 0, ADDITIONAL: 1
;; WARNING: recursion requested but not available

;; QUESTION SECTION:
;active.vault.service.consul.   IN      SRV

;; ANSWER SECTION:
active.vault.service.consul. 0  IN      SRV     1 1 8200 vault1.node.dc1.consul.

;; ADDITIONAL SECTION:
vault1.node.dc1.consul.      0  IN      A       172.17.33.46

;; Query time: 0 msec
;; SERVER: 127.0.0.1#53(127.0.0.1)
;; WHEN: Sat Apr 23 17:33:14 2016
;; MSG SIZE  rcvd: 172
% dig +short standby.vault.service.consul srv
1 1 8200 vault3.node.dc1.consul.
1 1 8200 vault2.node.dc1.consul.
% dig +short vault.service.consul srv
1 1 8200 vault3.node.dc1.consul.
1 1 8200 vault1.node.dc1.consul.
1 1 8200 vault2.node.dc1.consul.
% dig +short vault.service.consul a
172.17.33.46
172.17.34.32
172.17.35.29
vault1% vault seal
% dig +short vault.service.consul srv
1 1 8200 vault3.node.dc1.consul.
1 1 8200 vault2.node.dc1.consul.
vault1% vault unseal
Key (will be hidden):
Sealed: false
Key Shares: 5
Key Threshold: 3
Unseal Progress: 0
% dig +short vault.service.consul srv
1 1 8200 vault1.node.dc1.consul.
1 1 8200 vault3.node.dc1.consul.
1 1 8200 vault2.node.dc1.consul.
```

#### Backend Reference: etcd (Community-Supported)

For etcd, the following options are supported:

  * `path` (optional) - The path within etcd where data will be stored.
    Defaults to "vault/".

  * `address` (optional) - The address(es) of the etcd instance(s) to talk to.
    Can be comma separated list (protocol://host:port) of many etcd instances.
    Defaults to "http://localhost:2379" if not specified. May also be specified
    via the ETCD_ADDR environment variable.

  * `sync` (optional) - Should we synchronize the list of available etcd
    servers on startup?  This is a **string** value to allow for auto-sync to
    be implemented later. It can be set to "0", "no", "n", "false", "1", "yes",
    "y", or "true".  Defaults to on.  Set to false if your etcd cluster is
    behind a proxy server and syncing causes Vault to fail.

  * `ha_enabled` (optional) - Setting this to `"1"`, `"t"`, or `"true"` will
    enable HA mode. _This is currently *known broken*._ This option can also be
    provided via the environment variable `ETCD_HA_ENABLED`. If you are
    upgrading from a version of Vault where HA support was enabled by default,
    it is _very important_ that you set this parameter _before_ upgrading!

  * `username` (optional) - Username to use when authenticating with the etcd
    server.  May also be specified via the ETCD_USERNAME environment variable.

  * `password` (optional) - Password to use when authenticating with the etcd
    server.  May also be specified via the ETCD_PASSWORD environment variable.

  * `tls_ca_file` (optional) - The path to the CA certificate used for etcd
    communication.  Defaults to system bundle if not specified.

  * `tls_cert_file` (optional) - The path to the certificate for etcd
    communication.

  * `tls_key_file` (optional) - The path to the private key for etcd
    communication.

#### Backend Reference: Zookeeper (Community-Supported)

For Zookeeper, the following options are supported:

  * `path` (optional) - The path within Zookeeper where data will be stored.
    Defaults to "vault/".

  * `address` (optional) - The address(es) of the Zookeeper instance(s) to talk
    to. Can be comma separated list (host:port) of many Zookeeper instances.
    Defaults to "localhost:2181" if not specified.

The following optional settings can be used to configure zNode ACLs:

  * `auth_info` (optional) - Authentication string in Zookeeper AddAuth format
    (`schema:auth`). As an example, `digest:UserName:Password` could be used to
    authenticate as user `UserName` using password `Password` with the `digest`
    mechanism.

  * `znode_owner` (optional) - If specified, Vault will always set all
    permissions (CRWDA) to the ACL identified here via the Schema and User
    parts of the Zookeeper ACL format. The expected format is
    `schema:user-ACL-match`. Some examples:
    * `digest:UserName:HIDfRvTv623G==` - Access for the user `UserName` with
      the corresponding digest `HIDfRvTv623G==`
    * `ip:127.0.0.1` - Access from localhost only
    * `ip:70.95.0.0/16` - Any host on the 70.95.0.0 network (CIDRs are
      supported starting from Zookeeper 3.5.0)

If neither of these is set, the backend will not authenticate with Zookeeper
and will set the OPEN_ACL_UNSAFE ACL on all nodes. In this scenario, anyone
connected to Zookeeper could change Vault’s znodes and, potentially, take Vault
out of service.

Some sample configurations:

```
backend "zookeeper" {
  znode_owner = "digest:vaultUser:raxgVAfnDRljZDAcJFxznkZsExs="
  auth_info = "digest:vaultUser:abc"
}
```

The above configuration causes Vault to set an ACL on all of its zNodes
permitting access to vaultUser only. If the `digest` schema is used, please
protect this file as it contains the cleartext password. As per Zookeeper's ACL
model, the digest value (in znode_owner) must match the user (in znode_owner).

```
backend "zookeeper" {
  znode_owner = "ip:127.0.0.1"
}
```

The above example allows access from localhost only - as this is the `ip` no
auth_info is required since Zookeeper uses the address of the client for the
ACL check.

#### Backend Reference: DynamoDB (Community-Supported)

The DynamoDB optionally supports HA. Because Dynamo does not support session
lifetimes on its locks, a Vault node that has failed, rather than shut down in
an orderly fashion, will require manual cleanup rather than failing over
automatically. See the documentation of `recovery_mode` to better understand
this process. To enable HA, set the `ha_enabled` option.

The DynamoDB backend has the following options:

  * `table` (optional) - The name of the DynamoDB table to store data in. The
    default table name is `vault-dynamodb-backend`. This option can also be
    provided via the environment variable `AWS_DYNAMODB_TABLE`. If the
    specified table does not yet exist, it will be created during
    initialization.

  * `read_capacity` (optional) - The read capacity to provision when creating
    the DynamoDB table. This is the maximum number of reads consumed per second
    on the table. The default value is 5. This option can also be provided via
    the environment variable `AWS_DYNAMODB_READ_CAPACITY`.

  * `write_capacity` (optional) - The write capacity to provision when creating
    the DynamoDB table. This is the maximum number of writes performed per
    second on the table. The default value is 5. This option can also be
    provided via the environment variable `AWS_DYNAMODB_WRITE_CAPACITY`.

  * `access_key` - (required) The AWS access key. It must be provided, but it
    can also be sourced from the `AWS_ACCESS_KEY_ID` environment variable.

  * `secret_key` - (required) The AWS secret key. It must be provided, but it
    can also be sourced from the `AWS_SECRET_ACCESS_KEY` environment variable.

  * `session_token` - (optional) The AWS session token. It can also be sourced
    from the `AWS_SESSION_TOKEN` environment variable.

  * `endpoint` - (optional) An alternative (AWS compatible) DynamoDB endpoint
    to use. It can also be sourced from the `AWS_DYNAMODB_ENDPOINT` environment
    variable.

  * `region` (optional) - The AWS region. It can be sourced from the
    `AWS_DEFAULT_REGION` environment variable and will default to `us-east-1`
    if not specified.

  * `max_parallel` (optional) - The maximum number of concurrent requests to
    DynamoDB. Defaults to `"128"`.

  * `ha_enabled` (optional) - Setting this to `"1"`, `"t"`, or `"true"` will
    enable HA mode. Please ensure you have read the documentation for the
    `recovery_mode` option before enabling this. This option can also be
    provided via the environment variable `DYNAMODB_HA_ENABLED`. If you are
    upgrading from a version of Vault where HA support was enabled by default,
    it is _very important_ that you set this parameter _before_ upgrading!

  * `recovery_mode` (optional) - When the Vault leader crashes or is killed
    without being able to shut down properly, no other node can become the new
    leader because the DynamoDB table still holds the old leader's lock record.
    To recover from this situation, one can start a single Vault node with this
    option set to `"1"`, `"t"`, or `"true"` and the node will remove the old
    lock from DynamoDB. It is important that only one node is running in
    recovery mode! After this node has become the leader, other nodes can be
    started with regular configuration. This option can also be provided via
    the environment variable `RECOVERY_MODE`.

For more information about the read/write capacity of DynamoDB tables, see the
[official AWS DynamoDB
docs](https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/WorkingWithTables.html#ProvisionedThroughput).
If you are running your Vault server on an EC2 instance, you can also make use
of the EC2 instance profile service to provide the credentials Vault will use
to make DynamoDB API calls. Leaving the `access_key` and `secret_key` fields
empty will cause Vault to attempt to retrieve credentials from the metadata
service.

#### Backend Reference: S3 (Community-Supported)

For S3, the following options are supported:

  * `bucket` (required) - The name of the S3 bucket to use. It must be provided, but it can also be sourced from the `AWS_S3_BUCKET` environment variable.

  * `access_key` - (required) The AWS access key. It must be provided, but it can also be sourced from the `AWS_ACCESS_KEY_ID` environment variable.

  * `secret_key` - (required) The AWS secret key. It must be provided, but it can also be sourced from the `AWS_SECRET_ACCESS_KEY` environment variable.

  * `session_token` - (optional) The AWS session token. It can also be sourced from the `AWS_SESSION_TOKEN` environment variable.

  * `endpoint` - (optional) An alternative (AWS compatible) S3 endpoint to use. It can also be sourced from the `AWS_S3_ENDPOINT` environment variable.

  * `region` (optional) - The AWS region. It can be sourced from the `AWS_DEFAULT_REGION` environment variable and will default to `us-east-1` if not specified.

If you are running your Vault server on an EC2 instance, you can also make use
of the EC2 instance profile service to provide the credentials Vault will use to
make S3 API calls.  Leaving the `access_key` and `secret_key` fields empty
will cause Vault to attempt to retrieve credentials from the metadata service.
You are responsible for ensuring your instance is launched with the appropriate
profile enabled. Vault will handle renewing profile credentials as they rotate.

#### Backend Reference: Google Cloud Storage (Community-Supported)

For Google Cloud Storage, the following options are supported:

  * `bucket` (required) - The name of the Google Cloud Storage bucket to use. It must be provided, but it can also be sourced from the `GOOGLE_STORAGE_BUCKET` environment variable.

  * `credentials_file` - (required) The path to a GCP [service account](https://cloud.google.com/compute/docs/access/service-accounts) private key file in [JSON format](https://cloud.google.com/storage/docs/authentication#generating-a-private-key). It must be provided, but it can also be sourced from the `GOOGLE_APPLICATION_CREDENTIALS` environment variable.

  * `max_parallel` (optional) - The maximum number of concurrent requests to Google Cloud Storage.
    Defaults to `"128"`.

#### Backend Reference: Azure (Community-Supported)

  * `accountName` (required) - The Azure Storage account name

  * `accountKey`  (required) - The Azure Storage account key

  * `container`   (required) - The Azure Storage Blob container name

  * `max_parallel` (optional) - The maximum number of concurrent requests to Azure. Defaults to `"128"`.

The current implementation is limited to a maximum of 4 MBytes per blob/file.

#### Backend Reference: Swift (Community-Supported)

For Swift, the following options are valid; only v1.0 auth endpoints are supported:

  * `container` (required) - The name of the Swift container to use. It must be provided, but it can also be sourced from the `OS_CONTAINER` environment variable.

  * `username` - (required) The OpenStack account/username. It must be provided, but it can also be sourced from the `OS_USERNAME` environment variable.

  * `password` - (required) The OpenStack password. It must be provided, but it can also be sourced from the `OS_PASSWORD` environment variable.

  * `auth_url` - (required) Then OpenStack auth endpoint to use. It can also be sourced from the `OS_AUTH_URL` environment variable.

  * `tenant` (optional) - The name of Tenant to use. It can be sourced from the `OS_TENANT_NAME` environment variable and will default to default tenant of for the username if not specified.

  * `max_parallel` (optional) - The maximum number of concurrent requests to Swift. Defaults to `"128"`.

#### Backend Reference: MySQL (Community-Supported)

The MySQL backend has the following options:

  * `username` (required) - The MySQL username to connect with.

  * `password` (required) - The MySQL password to connect with.

  * `address` (optional) - The address of the MySQL host. Defaults to
    "127.0.0.1:3306.

  * `database` (optional) - The name of the database to use. Defaults to "vault".

  * `table` (optional) - The name of the table to use. Defaults to "vault".

  * `tls_ca_file` (optional) - The path to the CA certificate to connect using TLS

#### Backend Reference: PostgreSQL (Community-Supported)

The PostgreSQL backend has the following options:

  * `connection_url` (required) - The connection string used to connect to PostgreSQL.

    Examples:

    * postgres://username:password@localhost:5432/database?sslmode=disable

    * postgres://username:password@localhost:5432/database?sslmode=verify-full

    A list of all supported parameters can be found in [the pq library documentation](https://godoc.org/github.com/lib/pq#hdr-Connection_String_Parameters).

  * `table` (optional) - The name of the table to write vault data to. Defaults
    to "vault_kv_store".

Add the following table and index to a new or existing PostgreSQL database:

```sql
CREATE TABLE vault_kv_store (
  parent_path TEXT COLLATE "C" NOT NULL,
  path        TEXT COLLATE "C",
  key         TEXT COLLATE "C",
  value       BYTEA,
  CONSTRAINT pkey PRIMARY KEY (path, key)
);

CREATE INDEX parent_path_idx ON vault_kv_store (parent_path);
```

If you're using a version of PostgreSQL prior to 9.5, create the following
function:

```sql
CREATE FUNCTION vault_kv_put(_parent_path TEXT, _path TEXT, _key TEXT, _value BYTEA) RETURNS VOID AS
$$
BEGIN
    LOOP
        -- first try to update the key
        UPDATE vault_kv_store
          SET (parent_path, path, key, value) = (_parent_path, _path, _key, _value)
          WHERE _path = path AND key = _key;
        IF found THEN
            RETURN;
        END IF;
        -- not there, so try to insert the key
        -- if someone else inserts the same key concurrently,
        -- we could get a unique-key failure
        BEGIN
            INSERT INTO vault_kv_store (parent_path, path, key, value)
              VALUES (_parent_path, _path, _key, _value);
            RETURN;
        EXCEPTION WHEN unique_violation THEN
            -- Do nothing, and loop to try the UPDATE again.
        END;
    END LOOP;
END;
$$
LANGUAGE plpgsql;
```

More info can be found in the [PostgreSQL documentation](http://www.postgresql.org/docs/9.4/static/plpgsql-control-structures.html#PLPGSQL-UPSERT-EXAMPLE):

#### Backend Reference: Inmem

The in-memory backend has no configuration options.

#### Backend Reference: File

The file backend has the following options:

  * `path` (required) - The path on disk to a directory where the
      data will be stored.
