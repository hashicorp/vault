---
layout: "docs"
page_title: "Secret Backend: PostgreSQL"
sidebar_current: "docs-secrets-postgresql"
description: |-
  The PostgreSQL secret backend for Vault generates database credentials to access PostgreSQL.
---

# PostgreSQL Secret Backend

Name: `postgresql`

The PostgreSQL secret backend for Vault generates database credentials
dynamically based on configured roles. This means that services that need
to access a database no longer need to hardcode credentials: they can request
them from Vault, and use Vault's leasing mechanism to more easily roll keys.

Additionally, it introduces a new ability: with every service accessing
the database with unique credentials, it makes auditing much easier when
questionable data access is discovered: you can track it down to the specific
instance of a service based on the SQL username.

Vault makes use both of its own internal revocation system as well as the
`VALID UNTIL` setting when creating PostgreSQL users to ensure that users
become invalid within a reasonable time of the lease expiring.

This page will show a quick start for this backend. For detailed documentation
on every path, use `vault help` after mounting the backend.

## Quick Start

The first step to using the PostgreSQL backend is to mount it.
Unlike the `generic` backend, the `postgresql` backend is not mounted by default.

```text
$ vault mount postgresql
Successfully mounted 'postgresql' at 'postgresql'!
```

Vault must be configured to connect to the PostgreSQL:

```text
$ vault write postgresql/config/connection \
    value="host=localhost user=username password=password"
```

This value must be a PG-style connection string, and the specified user must
have permission to manage the database.

Note, if you get an error like:

```text
pq: SSL is not enabled on the server
```

this means your PostgreSQL server has not enabled SSL. It is highly recommended
that you configure your PostgreSQL server to communicate via SSL.

The SSL check can be disabled by specifying the `sslmode=disable` attribute in
the PostgreSQL connection string:

```text
$ vault write postgresql/config/connection \
    value="host=localhost user=username password=password sslmode=disable"
```

Vault's PostgreSQL integration is role-based, so you must create a role for
which to request credentials:

```text
$ vault write postgresql/roles/production \
    name=production
```

Vault is now configured to create and manage credentials for Postgres!

```text
$ vault read postgresql/creds/production
Key             Value
lease_id        postgresql/creds/production/8ade2cde-5081-e3b7-af1a-3b9fb070df66
lease_duration  3600
password        56b43bc3-b285-4803-abdf-662d6a105bd0
username        vault-root-1430141210-1847
```

If you get stuck at any time, simply run `vault help postgresql` or with a
subpath for interactive help output.
