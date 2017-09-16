---
layout: "docs"
page_title: "MySQL Secret Backend"
sidebar_current: "docs-secrets-mysql"
description: |-
  The MySQL secret backend for Vault generates database credentials to access MySQL.
---

# MySQL Secret Backend

Name: `mysql`

~> **Deprecation Note:** This backend is deprecated in favor of the
combined databases backend added in v0.7.1. See the documentation for
the new implementation of this backend at
[MySQL/MariaDB Database Plugin](/docs/secrets/databases/mysql-maria.html).

The MySQL secret backend for Vault generates database credentials
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

The first step to using the mysql backend is to mount it.
Unlike the `kv` backend, the `mysql` backend is not mounted by default.

```
$ vault mount mysql
Successfully mounted 'mysql' at 'mysql'!
```

Next, we must configure Vault to know how to connect to the MySQL
instance. This is done by providing a [DSN (Data Source Name)](https://github.com/go-sql-driver/mysql#dsn-data-source-name):

```
$ vault write mysql/config/connection \
    connection_url="root:root@tcp(192.168.33.10:3306)/"
Success! Data written to: mysql/config/connection
```

In this case, we've configured Vault with the user "root" and password "root,
connecting to an instance at "192.168.33.10" on port 3306. It is not necessary
that Vault has the root user, but the user must have privileges to create
other users, namely the `GRANT OPTION` privilege.

For using UNIX socket use: `root:root@unix(/path/to/socket)/`.

Optionally, we can configure the lease settings for credentials generated
by Vault. This is done by writing to the `config/lease` key:

```
$ vault write mysql/config/lease \
    lease=1h \
    lease_max=24h
Success! Data written to: mysql/config/lease
```

This restricts each credential to being valid or leased for 1 hour
at a time, with a maximum use period of 24 hours. This forces an
application to renew their credentials at least hourly, and to recycle
them once per day.

The next step is to configure a role. A role is a logical name that maps
to a policy used to generate those credentials. For example, lets create
a "readonly" role:

```
$ vault write mysql/roles/readonly \
    sql="CREATE USER '{{name}}'@'%' IDENTIFIED BY '{{password}}';GRANT SELECT ON *.* TO '{{name}}'@'%';"
Success! Data written to: mysql/roles/readonly
```

By writing to the `roles/readonly` path we are defining the `readonly` role.
This role will be created by evaluating the given `sql` statements. By
default, the `{{name}}` and `{{password}}` fields will be populated by
Vault with dynamically generated values. This SQL statement is creating
the named user, and then granting it `SELECT` or read-only privileges
to tables in the database. More complex `GRANT` queries can be used to
customize the privileges of the role. See the [MySQL manual](https://dev.mysql.com/doc/refman/5.7/en/grant.html)
for more information.

To generate a new set of credentials, we simply read from that role:

```
$ vault read mysql/creds/readonly
Key           	Value
lease_id      	mysql/creds/readonly/bd404e98-0f35-b378-269a-b7770ef01897
lease_duration	3600
password      	132ae3ef-5a64-7499-351e-bfe59f3a2a21
username      	readonly-aefa635a-18
```

By reading from the `creds/readonly` path, Vault has generated a new
set of credentials using the `readonly` role configuration. Here we
see the dynamically generated username and password, along with a one
hour lease.

Using ACLs, it is possible to restrict using the mysql backend such
that trusted operators can manage the role definitions, and both
users and applications are restricted in the credentials they are
allowed to read.

Optionally, you may configure both the number of characters from the role name
that are truncated to form the display name portion of the mysql username
interpolated into the `{{name}}` field: the default is 10.

You may also configure the total number of characters allowed in the entire
generated username (the sum of the display name and uuid portions); the
default is 16. Note that versions of MySQL prior to 5.8 have a 16 character
total limit on user names, so it is probably not safe to increase this above
the default on versions prior to that.

## API

The MySQL secret backend has a full HTTP API. Please see the
[MySQL secret backend API](/api/secret/mysql/index.html) for more
details.
