---
layout: "docs"
page_title: "MSSQL - Storage Backends - Configuration"
sidebar_current: "docs-configuration-storage-mssql"
description: |-
  The MSSQL storage backend is used to persist Vault's data in a Microsoft SQL Server.
---

# MSSQL Storage Backend

The MSSQL storage backend is used to persist Vault's data in a Microsoft SQL Server.

- **No High Availability** – the MSSQL storage backend does not support high
  availability.

- **Community Supported** – the MSSQL storage backend is supported by the
  community. While it has undergone review by HashiCorp employees, they may not
  be as knowledgeable about the technology. If you encounter problems with them,
  you may be referred to the original author.

```hcl
storage "mssql" {
  server = "localhost"
  username = "user1234"
  password = "secret123!"
  database = "vault"
  table = "vault"
  appname = "vault"
  schema = "dbo"
  connectionTimeout = 30
  logLevel = 0
}
```

## `mssql` Parameters

- `server` `(string: <required>)` – host or host\instance.

- `username` `(string: "")` - enter the SQL Server Authentication user id or 
  the Windows Authentication user id in the DOMAIN\User format. 
  On Windows, if user id is empty or missing Single-Sign-On is used.

- `password` `(string: "")` – specifies the MSSQL password to connect to
  the database.

- `database` `(string: "Vault")` – Specifies the name of the database. If the
  database does not exist, Vault will attempt to create it.

- `table` `(string: "Vault")` – Specifies the name of the table. If the table
  does not exist, Vault will attempt to create it.

- `schema` `(string: "dbo")` – Specifies the name of the schema. If the schema
  does not exist, Vault will attempt to create it.

- `appname` `(string: "Vault")` – the application name.

- `connectionTimeout` `(int: 30)` – in seconds (default is 30).

- `logLevel` `(int: 0)` – logging flags (default 0/no logging, 63 for full logging) .

- `max_parallel` `(string: "128")` – Specifies the maximum number of concurrent
  requests to MSSQL.

## `mssql` Examples

### Custom Database, Table and Schema

This example shows configuring the MSSQL backend to use a custom database and
table name.

```hcl
storage "mssql" {
  database = "my-vault"
  table    = "vault-data"
  schema   = "vlt"
  username = "user1234"
  password = "pass5678"
}
```


