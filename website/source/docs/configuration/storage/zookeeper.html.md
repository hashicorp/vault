---
layout: "docs"
page_title: "Zookeeper - Storage Backends - Configuration"
sidebar_current: "docs-configuration-storage-zookeeper"
description: |-
  The Zookeeper storage backend is used to persist Vault's data in Zookeeper.
---

# Zookeeper Storage Backend

The Zookeeper storage backend is used to persist Vault's data in
[Zookeeper][zk].

- **High Availability** – the Zookeeper storage backend supports high
  availability.

- **Community Supported** – the Zookeeper storage backend is supported by the
  community. While it has undergone review by HashiCorp employees, they may not
  be as knowledgeable about the technology. If you encounter problems with them,
  you may be referred to the original author.

```hcl
storage "zookeeper" {
  address = "localhost:2181"
  path    = "vault/"
}
```

## `zookeeper` Parameters

- `address` `(string: "localhost:2181")` – Specifies the addresses of the
  Zookeeper instances as a comma-separated list.

- `path` `(string: "vault/")` – Specifies the path in Zookeeper where data will
  be stored.

The following optional settings can be used to configure zNode ACLs:

~> **Warning!** If neither `auth_info` nor `znode_owner` are set, the backend
will not authenticate with Zookeeper and will set the `OPEN_ACL_UNSAFE` ACL on
all nodes. In this scenario, anyone connected to Zookeeper could change Vault’s
znodes and, potentially, take Vault out of service.

- `auth_info` `(string: "")` – Specifies an authentication string in Zookeeper
  AddAuth format. For example, `digest:UserName:Password` could be used to
  authenticate as user `UserName` using password `Password` with the `digest`
  mechanism.

- `znode_owner` `(string: "")` – If specified, Vault will always set all
  permissions (CRWDA) to the ACL identified here via the Schema and User parts
  of the Zookeeper ACL format. The expected format is `schema:user-ACL-match`,
  for example:

    ```text
    # Access for user "UserName" with corresponding digest "HIDfRvTv623G=="
    digest:UserName:HIDfRvTv623G==
    ```

    ```text
    # Access from localhost only
    ip:127.0.0.1
    ```

    ```text
    # Access from any host on the 70.95.0.0 network (Zookeeper 3.5+)
    ip:70.95.0.0/16
    ```

## `zookeeper` Examples

### Custom Address and Path

This example shows configuring Vault to communicate with a Zookeeper
installation running on a custom port and to store data at a custom path.

```hcl
storage "zookeeper" {
  address = "localhost:3253"
  path    = "my-vault-data/"
}
```

### zNode Vault User Only

This example instructs Vault to set an ACL on all of its zNodes which permit
access only to the user "vaultUser". As per Zookeeper's ACL model, the digest
value in `znode_owner` must match the user in `znode_owner`.

```hcl
storage "zookeeper" {
  znode_owner = "digest:vaultUser:raxgVAfnDRljZDAcJFxznkZsExs="
  auth_info   = "digest:vaultUser:abc"
}
```

### zNode Localhost Only

This example instructs Vault to only allow access from localhost. As this is the
`ip` no `auth_info` is required since Zookeeper uses the address of the client
for the ACL check.

```hcl
storage "zookeeper" {
  znode_owner = "ip:127.0.0.1"
}
```

[zk]: https://zookeeper.apache.org/
