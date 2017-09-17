---
layout: "docs"
page_title: "MSSQL Secret Backend"
sidebar_current: "docs-secrets-mssql"
description: |-
  The MSSQL secret backend for Vault generates database credentials to access Microsoft Sql Server.
---

# MSSQL Secret Backend

Name: `mssql`

~> **Deprecation Note:** This backend is deprecated in favor of the
combined databases backend added in v0.7.1. See the documentation for
the new implementation of this backend at
[MSSQL Database Plugin](/docs/secrets/databases/mssql.html).

The MSSQL secret backend for Vault generates database credentials
dynamically based on configured roles. This means that services that need
to access a database no longer need to hardcode credentials: they can request
them from Vault, and use Vault's leasing mechanism to more easily roll keys.

Additionally, it introduces a new ability: with every service accessing
the database with unique credentials, it makes auditing much easier when
questionable data access is discovered: you can track it down to the specific
instance of a service based on the SQL username.

Vault makes use of its own internal revocation system to ensure that users
become invalid within a reasonable time of the lease expiring.

This page will show a quick start for this backend. For detailed documentation
on every path, use `vault path-help` after mounting the backend.

## Quick Start

The first step to using the mssql backend is to mount it.
Unlike the `kv` backend, the `mssql` backend is not mounted by default.

```
$ vault mount mssql
Successfully mounted 'mssql' at 'mssql'!
```

Next, we must configure Vault to know how to connect to the MSSQL
instance. This is done by providing a DSN (Data Source Name):

```
$ vault write mssql/config/connection \
    connection_string="server=localhost;port=1433;user id=sa;password=Password!;database=AdventureWorks;app name=vault;"
Success! Data written to: mssql/config/connection
```

In this case, we've configured Vault with the user "sa" and password "Password!",
connecting to an instance at "localhost" on port 1433. It is not necessary
that Vault has the sa login, but the user must have privileges to create
logins and manage processes. The fixed server roles `securityadmin` and
`processadmin` are examples of built-in roles that grant these permissions. The
user also must have privileges to create database users and grant permissions in
the databases that Vault manages.  The fixed database roles `db_accessadmin` and
`db_securityadmin` are examples or built-in roles that grant these permissions.

Optionally, we can configure the lease settings for credentials generated
by Vault. This is done by writing to the `config/lease` key:

```
$ vault write mssql/config/lease \
    ttl=1h \
    max_ttl=24h
Success! Data written to: mssql/config/lease
```

This restricts each credential to being valid or leased for 1 hour
at a time, with a maximum use period of 24 hours. This forces an
application to renew their credentials at least hourly, and to recycle
them once per day.

The next step is to configure a role. A role is a logical name that maps
to a policy used to generate those credentials. For example, lets create
a "readonly" role:

```
$ vault write mssql/roles/readonly \
    sql="CREATE LOGIN [{{name}}] WITH PASSWORD = '{{password}}'; USE AdventureWorks; CREATE USER [{{name}}] FOR LOGIN [{{name}}]; GRANT SELECT ON SCHEMA::dbo TO [{{name}}]"
Success! Data written to: mssql/roles/readonly
```

By writing to the `roles/readonly` path we are defining the `readonly` role.
This role will be created by evaluating the given `sql` statements. By
default, the `{{name}}` and `{{password}}` fields will be populated by
Vault with dynamically generated values. This SQL statement is creating
the named login on the server, user on the AdventureWorks database, and
then granting it `SELECT` on the `dbo` schema. More complex `GRANT` queries
can be used to customize the privileges of the role.

To generate a new set of credentials, we simply read from that role:

```
$ vault read mssql/creds/readonly
Key           	Value
lease_id      	mssql/creds/readonly/cdf23ac8-6dbd-4bf9-9919-6acaaa86ba6c
lease_duration	3600
password      	ee202d0d-e4fd-4410-8d14-2a78c5c8cb76
username      	root-a147d529-e7d6-4a16-8930-4c3e72170b19
```

By reading from the `creds/readonly` path, Vault has generated a new
set of credentials using the `readonly` role configuration. Here we
see the dynamically generated username and password, along with a one
hour lease.

Using ACLs, it is possible to restrict using the mssql backend such
that trusted operators can manage the role definitions, and both
users and applications are restricted in the credentials they are
allowed to read.

## API

The MSSQL secret backend has a full HTTP API. Please see the
[MSSQL secret backend API](/api/secret/mssql/index.html) for more
details.
