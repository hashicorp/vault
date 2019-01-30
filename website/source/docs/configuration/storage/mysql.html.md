---
layout: "docs"
page_title: "MySQL - Storage Backends - Configuration"
sidebar_title: "MySQL"
sidebar_current: "docs-configuration-storage-mysql"
description: |-
  The MySQL storage backend is used to persist Vault's data in a MySQL server or
  cluster.
---

# MySQL Storage Backend

The MySQL storage backend is used to persist Vault's data in a [MySQL][mysql]
server or cluster.

- **High Availability** – the MySQL storage backend supports high availability.
  Note that due to the way mysql locking functions work they are lost if a connection
  dies. If you would like to not have frequent changes in your elected leader you
  can increase interactive_timeout and wait_timeout MySQL config to much higher than
  default which is set at 8 hours.

- **Community Supported** – the MySQL storage backend is supported by the
  community. While it has undergone review by HashiCorp employees, they may not
  be as knowledgeable about the technology. If you encounter problems with them,
  you may be referred to the original author.

```hcl
storage "mysql" {
  username = "user1234"
  password = "secret123!"
  database = "vault"
}
```

## `mysql` Parameters

- `address` `(string: "127.0.0.1:3306")` – Specifies the address of the MySQL
  host.

- `database` `(string: "vault")` – Specifies the name of the database. If the
  database does not exist, Vault will attempt to create it.

- `table` `(string: "vault")` – Specifies the name of the table. If the table
  does not exist, Vault will attempt to create it.

- `tls_ca_file` `(string: "")` – Specifies the path to the CA certificate to
  connect using TLS.

- `max_parallel` `(string: "128")` – Specifies the maximum number of concurrent
  requests to MySQL.

Additionally, Vault requires the following authentication information.

- `username` `(string: <required>)` – Specifies the MySQL username to connect to
  the database.

- `password` `(string: <required>)` – Specifies the MySQL password to connect to
  the database.

### High Availability Parameters

- `ha_enabled` `(string: "true")` -  Specifies if high availability mode is
  enabled. This is a boolean value, but it is specified as a string like "true"
  or "false".

- `lock_table` `(string: "vault_lock")` – Specifies the name of the table to
  use for storing high availability information. By default, this is the name
  of the `table` suffixed with `_lock`. If the table does not exist, Vault will
  attempt to create it.

## `mysql` Examples

### Custom Database and Table

This example shows configuring the MySQL backend to use a custom database and
table name.

```hcl
storage "mysql" {
  database = "my-vault"
  table    = "vault-data"
  username = "user1234"
  password = "pass5678"
}
```

[mysql]: https://dev.mysql.com
