---
layout: "docs"
page_title: "PostgreSQL - Storage Backends - Configuration"
sidebar_current: "docs-configuration-storage-postgresql"
description: |-
  The PostgreSQL storage backend is used to persist Vault's data in a PostgreSQL
  server or cluster.
---

# PostgreSQL Storage Backend

The PostgreSQL storage backend is used to persist Vault's data in a
[PostgreSQL][postgresql] server or cluster.

- **No High Availability** – the PostgreSQL storage backend does not support
  high availability.

- **Community Supported** – the PostgreSQL storage backend is supported by the
  community. While it has undergone review by HashiCorp employees, they may not
  be as knowledgeable about the technology. If you encounter problems with them,
  you may be referred to the original author.

```hcl
storage "postgresql" {
  connection_url = "postgres://user123:secret123!@localhost:5432/vault"
}
```

The PostgreSQL storage backend does not automatically create the table. Here is
some sample SQL to create the schema and indexes.

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

If you're using a version of PostgreSQL prior to 9.5, create the following function:

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

## `postgresql` Parameters

- `connection_url` `(string: <required>)` – Specifies the connection string to
  use to authenticate and connect to PostgreSQL. A full list of supported
  parameters can be found in [the pq library documentation][pglib]. For example
  connection string URLs, see the examples section below.

- `table` `(string: "vault_kv_store")` – Specifies the name of the table in
  which to write Vault data. This table must already exist (Vault will not
  attempt to create it).

- `max_parallel` `(string: "128")` – Specifies the maximum number of concurrent
  requests to PostgreSQL.

## `postgresql` Examples

### Custom SSL Verification

This example shows connecting to a PostgreSQL cluster using full SSL
verification (recommended).

```hcl
storage "postgresql" {
  connection_url = "postgres://user:pass@localhost:5432/database?sslmode=verify-full"
}
```

To disable SSL verification (not recommended), replace `verify-full` with
`disable`:

```hcl
storage "postgresql" {
  connection_url = "postgres://user:pass@localhost:5432/database?sslmode=disable"
}
```

[postgresql]: https://www.postgresql.org/
[pglib]: https://godoc.org/github.com/lib/pq#hdr-Connection_String_Parameters
