---
layout: "docs"
page_title: "CockroachDB - Storage Backends - Configuration"
sidebar_current: "docs-configuration-storage-cockroachdb"
description: |-
  The CockroachDB storage backend is used to persist Vault's data in a CockroachDB
  server or cluster.
---

# CockroachDB Storage Backend

The CockroachDB storage backend is used to persist Vault's data in a
[CockroachDB][cockroachdb] server or cluster.

- **No High Availability** – the CockroachDB storage backend does not support
  high availability.

- **Community Supported** – the CockroachDB storage backend is supported by the
  community. While it has undergone development and review by HashiCorp 
  employees, they may not be as knowledgeable about the technology.

```hcl
storage "cockroachdb" {
  connection_url = "postgres://user123:secret123!@localhost:5432/vault"
}
```

**Note** - CockroachDB is compatible with the PostgreSQL database driver and 
uses that driver to interact with the database.

## `cockroachdb` Parameters

- `connection_url` `(string: <required>)` – Specifies the connection string to
  use to authenticate and connect to CockroachDB. A full list of supported
  parameters can be found in [the pq library documentation][pglib]. For example
  connection string URLs, see the examples section below.

- `table` `(string: "vault_kv_store")` – Specifies the name of the table in
  which to write Vault data. This table must already exist (Vault will not
  attempt to create it).

- `max_parallel` `(string: "128")` – Specifies the maximum number of concurrent
  requests to CockroachDB.

## `cockroachdb` Examples

This example shows connecting to a PostgresSQL cluster using full SSL
verification (recommended).

```hcl
storage "cockroachdb" {
  connection_url = "postgres://user:pass@localhost:5432/database?sslmode=verify-full"
}
```

To disable SSL verification (not recommended), replace `verify-full` with
`disable`:

```hcl
storage "cockroachdb" {
  connection_url = "postgres://user:pass@localhost:5432/database?sslmode=disable"
}
```

[cockroachdb]: https://www.cockroachlabs.com/
[pglib]: https://godoc.org/github.com/lib/pq#hdr-Connection_String_Parameters
