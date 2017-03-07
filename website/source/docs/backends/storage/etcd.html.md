---
layout: "docs"
page_title: "Storage Backends"
sidebar_current: "docs-backends-storage-etcd"
description: |-
  TODO
---

# Etcd Storage Backend

The etcd physical backend supports both v2 and v3 APIs. To explicitly specify
the API version, use the `etcd_api` configuration parameter. The default version
is auto-detected based on the version of the etcd cluster. If the cluster
version is 3.1+ and there has been no data written using the v2 API, the
auto-detected default is v3.

The v2 API has known issues with HA support and should not be used in HA
scenarios.

- **High Availability** - the Etcd backend supports high availability.

- **Support Level** - the Etcd backend is supported by the community. While it
  has undergone review by HashiCorp employees, they may not be as knowledgeable
  about the technology. If you encounter problems with them, you may be referred
  to the original author.

```hcl
backend "etcd" {
  address  = "http://localhost:2379"
  etcd_api = "v3"
}
```

## `etcd` Parameters

- `path` (optional) - The path within etcd where data will be stored.
  Defaults to "vault/".

- `address` (optional) - The address(es) of the etcd instance(s) to talk to.
  Can be comma separated list (protocol://host:port) of many etcd instances.
  Defaults to "http://localhost:2379" if not specified. May also be specified
  via the ETCD_ADDR environment variable.

- `etcd_api` (optional) - Set to `"v2"` or `"v3"` to explicitly set the API
  version that the backend will use.

- `sync` (optional) - Should we synchronize the list of available etcd
  servers on startup?  This is a **string** value to allow for auto-sync to
  be implemented later. It can be set to "0", "no", "n", "false", "1", "yes",
  "y", or "true".  Defaults to on.  Set to false if your etcd cluster is
  behind a proxy server and syncing causes Vault to fail.

- `ha_enabled` (optional) - Setting this to `"1"`, `"t"`, or `"true"` will
  enable HA mode. _This is currently *known broken*._ This option can also be
  provided via the environment variable `ETCD_HA_ENABLED`. If you are
  upgrading from a version of Vault where HA support was enabled by default,
  it is _very important_ that you set this parameter _before_ upgrading!

- `username` (optional) - Username to use when authenticating with the etcd
  server.  May also be specified via the ETCD_USERNAME environment variable.

- `password` (optional) - Password to use when authenticating with the etcd
  server.  May also be specified via the ETCD_PASSWORD environment variable.

- `tls_ca_file` (optional) - The path to the CA certificate used for etcd
  communication.  Defaults to system bundle if not specified.

- `tls_cert_file` (optional) - The path to the certificate for etcd
  communication.

- `tls_key_file` (optional) - The path to the private key for etcd
  communication.
