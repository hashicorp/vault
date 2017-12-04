---
layout: "docs"
page_title: "MySQL - Storage Backends - Configuration"
sidebar_current: "docs-configuration-storage-mysql"
description: |-
  The MySQL storage backend is used to persist Vault's data in a MySQL server or
  cluster.
---

# MySQL Storage Backend

The MySQL storage backend is used to persist Vault's data in a [MySQL][mysql]
server or cluster.

- **No High Availability** – the MySQL storage backend does not support high
  availability.

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

- `password` `(string: <required)` – Specifies the MySQL password to connect to
  the database.

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
