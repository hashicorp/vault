---
layout: "docs"
page_title: "Etcd - Storage Backends - Configuration"
sidebar_current: "docs-configuration-storage-etcd"
description: |-
  The Etcd storage backend is used to persist Vault's data in Etcd. It supports
  both the v2 and v3 Etcd APIs, and the version is automatically detected based
  on the version of the Etcd cluster.
---

# Etcd Storage Backend

The Etcd storage backend is used to persist Vault's data in [Etcd][etcd]. It
supports both the v2 and v3 Etcd APIs, and the version is automatically detected
based on the version of the Etcd cluster.

- **High Availability** – the Etcd storage backend supports high availability.
  The v2 API has known issues with HA support and should not be used in HA
  scenarios.

- **Community Supported** – the Etcd storage backend is supported by CoreOS.
  While it has undergone review by HashiCorp employees, they may not be as
  knowledgeable about the technology. If you encounter problems with them, you
  may be referred to the original author.

```hcl
storage "etcd" {
  address  = "http://localhost:2379"
  etcd_api = "v3"
}
```

## `etcd` Parameters

- `address` `(string: "http://localhost:2379")` – Specifies the addresses of the
  Etcd instances as a comma-separated list. This can also be provided via the
  environment variable `ETCD_ADDR`.

- `discovery_srv` `(string: "example.com")` - Specifies the domain name to
  query for SRV records describing cluster endpoints. This can also be provided
  via the environment variable `ETCD_DISCOVERY_SRV`.

- `etcd_api` `(string: "<varies>")` – Specifies the version of the API to
  communicate with. By default, this is derived automatically. If the cluster
  version is 3.1+ and there has been no data written using the v2 API, the
  auto-detected default is v3.

- `ha_enabled` `(bool: false)` – Specifies if high availability should be
  enabled. This can also be provided via the environment variable
  `ETCD_HA_ENABLED`.

- `path` `(string: "vault/")` – Specifies the path in Etcd where Vault data will
  be stored.

- `sync` `(string: "true")` – Specifies whether to sync the list of available
  Etcd services on startup. This is a string that is coerced into a boolean
  value. You may want to set this to false if your cluster is behind a proxy
  server and syncing causes Vault to fail.

- `username` `(string: "")` – Specifies the username to use when authenticating
  with the Etcd server. This can also be provided via the environment variable
  `ETCD_USERNAME`.

- `password` `(string: "")` – Specifies the password to use when authenticating
  with the Etcd server. This can also be provided via the environment variable
  `ETCD_PASSWORD`.

- `tls_ca_file` `(string: "")` – Specifies the path to the CA certificate used
  for Etcd communication. This defaults to system bundle if not specified.

- `tls_cert_file` `(string: "")` – Specifies the path to the certificate for
  Etcd communication.

- `tls_key_file` `(string: "")` – Specifies the path to the private key for Etcd
  communication.

## `etcd` Examples

### DNS Discovery of cluster members

This example configures vault to discover the Etcd cluster members via SRV
records as outlined in the
[DNS Discovery protocol documentation][dns discovery].

```hcl
storage "etcd" {
  discovery_srv = "example.com"
}
```

### Custom Authentication

This example shows connecting to the Etcd cluster using a username and password.

```hcl
storage "etcd" {
  username = "user1234"
  password = "pass5678"
}
```

### Custom Path

This example shows storing data in a custom path.

```hcl
storage "etcd" {
  path = "my-vault-data/"
}
```

### Enabling High Availability

This example show enabling high availability for the Etcd storage backend.

```hcl
api_addr = "https://vault-leader.my-company.internal"

storage "etcd" {
  ha_enabled    = true
  ...
}
```

[etcd]: https://coreos.com/etcd "Etcd by CoreOS"
[dns discovery]: https://coreos.com/etcd/docs/latest/op-guide/clustering.html#dns-discovery "Etcd cluster DNS Discovery"
