---
layout: "docs"
page_title: "Database Secret Backend"
sidebar_current: "docs-secrets-databases"
description: |-
  Top page for database secret backend information
---

# Databases

Name: `Database`

The Database secret backend for Vault generates database credentials dynamically
based on configured roles. It works with a number of different databases through
a plugin interface. There are a number of builtin database types and an exposed
framework for running custom database types for extendability. This means that
services that need to access a database no longer need to hardcode credentials:
they can request them from Vault, and use Vault's leasing mechanism to more
easily roll keys.

Additionally, it introduces a new ability: with every service accessing the
database with unique credentials, it makes auditing much easier when
questionable data access is discovered: you can track it down to the specific
instance of a service based on the SQL username.

Vault makes use of its own internal revocation system to ensure that users
become invalid within a reasonable time of the lease expiring.

This page will show a quick start for this backend. For detailed documentation
on every path, use vault path-help after mounting the backend.

## Quick Start

The first step in using the Database backend is mounting it.

```text
$ vault mount database
Successfully mounted 'database' at 'database'!
```

Next, we must configure this backend to connect to a database. In this example
we will connect to a MySQL database, but the configuration details needed for
other plugin types can be found in their docs pages. This backend can configure
multiple database connections, therefore a name for the connection must be
provided; we'll call this one simply "mysql".

```
$ vault write database/config/mysql \
    plugin_name=mysql-database-plugin \
    connection_url="root:mysql@tcp(127.0.0.1:3306)/" \
    allowed_roles="readonly"

The following warnings were returned from the Vault server:
* Read access to this endpoint should be controlled via ACLs as it will return the connection details as is, including passwords, if any.
```

The next step is to configure a role. A role is a logical name that maps to a
policy used to generate those credentials. A role needs to be configured with
the database name we created above, and the default/max TTLs. For example, lets
create a "readonly" role:

```
$ vault write database/roles/readonly \
    db_name=mysql \
    creation_statements="CREATE USER '{{name}}'@'%' IDENTIFIED BY '{{password}}';GRANT SELECT ON *.* TO '{{name}}'@'%';" \
    default_ttl="1h" \
    max_ttl="24h"
Success! Data written to: database/roles/readonly
```
By writing to the roles/readonly path we are defining the readonly role. This
role will be created by evaluating the given creation statements. By default,
the {{name}} and {{password}} fields will be populated by the plugin with
dynamically generated values. In other plugins the {{expiration}} field could
also be supported. This SQL statement is creating the named user, and then
granting it SELECT or read-only privileges to tables in the database. More
complex GRANT queries can be used to customize the privileges of the role.
Custom revocation statements could be passed too, but this plugin has a default
statement we can use.

To generate a new set of credentials, we simply read from that role:

```
$ vault read database/creds/readonly
Key            	Value
---            	-----
lease_id       	database/creds/readonly/2f6a614c-4aa2-7b19-24b9-ad944a8d4de6
lease_duration 	1h0m0s
lease_renewable	true
password       	8cab931c-d62e-a73d-60d3-5ee85139cd66
username       	v-root-e2978cd0-
```

## Custom Plugins

This backend allows custom database types to be run through the exposed plugin
interface. Please see the [Custom database
plugin](/docs/secrets/databases/custom.html) for more information.

## API

The Database secret backend has a full HTTP API. Please see the [Database secret
backend API](/api/secret/databases/index.html) for more details.
