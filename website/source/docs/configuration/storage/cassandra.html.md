---
layout: "docs"
page_title: "Cassandra - Storage Backends - Configuration"
sidebar_current: "docs-configuration-storage-cassandra"
description: |-
  The Cassandra storage backend is used to persist Vault's data in an Apache
  Cassandra cluster.
---

# Cassandra Storage Backend

The Cassandra storage backend is used to persist Vault's data in an [Apache
Cassandra][cassandra] cluster.

- **No High Availability** – the Cassandra storage backend does not support high
  availability.

- **Community Supported** – the Cassandra storage backend is supported by the
  community. While it has undergone review by HashiCorp employees, they may not
  be as knowledgeable about the technology. If you encounter problems with it,
  you may be referred to the original author.

```hcl
storage "cassandra" {
  hosts            = "localhost"
  consistency      = "LOCAL_QUORUM"
  protocol_version = 3
}
```

The Cassandra storage backend does not automatically create the keyspace and
table. This sample configuration can be used as a guide, but you will want to
ensure the keyspace [replication options][replication-options]
are appropriate for your cluster:

```cql
CREATE KEYSPACE "vault" WITH REPLICATION = {
    'class': 'SimpleStrategy',
    'replication_factor': 1
};

CREATE TABLE "vault"."entries" (
    bucket text,
    key text,
    value blob,
    PRIMARY KEY (bucket, key)
) WITH CLUSTERING ORDER BY (key ASC);
```

## `cassandra` Parameters

* `hosts` `(string: "127.0.0.1")` – Comma-separated list of Cassandra hosts to
  connect to.

* `keyspace` `(string: "vault")` Cassandra keyspace to use.

* `table` `(string: "entries")` – Table within the `keyspace` in which to store
  data.

* `consistency` `(string: "LOCAL_QUORUM")` Consistency level to use when
  reading/writing data. If set, must be one of `"ANY"`, `"ONE"`, `"TWO"`,
  `"THREE"`, `"QUORUM"`, `"ALL"`, `"LOCAL_QUORUM"`, `"EACH_QUORUM"`, or 
  `"LOCAL_ONE"`.

* `protocol_version` `(int: 2)` Cassandra protocol version to use.

* `username` `(string: "")` – Username to use when authenticating with the
  Cassandra hosts.

* `password` `(string: "")` – Password to use when authenticating with the
  Cassandra hosts.

* `connection_timeout` `(int: 0)` - A timeout in seconds to wait until a
  connection is established with the Cassandra hosts.

* `tls` `(int: 0)` – If `1`, indicates the connection with the Cassandra hosts
  should use TLS.

* `pem_bundle_file` `(string: "")` - Specifies a file containing a
  certificate and private key; a certificate, private key, and issuing CA
  certificate; or just a CA certificate.

* `pem_json_file` `(string: "")` - Specifies a JSON file containing a certificate
  and private key; a certificate, private key, and issuing CA certificate;
  or just a CA certificate.

* `tls_skip_verify` `(int: 0)` - If `1`, then TLS host verification
  will be disabled for Cassandra. Defaults to `0`.

* `tls_min_version` `(string: "tls12")` - Minimum TLS version to use. Accepted
  values are `tls10`, `tls11` or `tls12`. Defaults to `tls12`.

[cassandra]: http://cassandra.apache.org/
[replication-options]: https://docs.datastax.com/en/cassandra/2.1/cassandra/architecture/architectureDataDistributeReplication_c.html
