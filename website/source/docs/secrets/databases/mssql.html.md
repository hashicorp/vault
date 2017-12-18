---
layout: "docs"
page_title: "MSSQL Database Plugin - Database Secret Backend"
sidebar_current: "docs-secrets-databases-mssql"
description: |-
  The MSSQL plugin for Vault's Database backend generates database credentials to access Microsoft SQL Server.
---

# MSSQL Database Plugin

Name: `mssql-database-plugin`

The MSSQL Database Plugin is one of the supported plugins for the Database
backend. This plugin generates database credentials dynamically based on
configured roles for the MSSQL database.

See the [Database Backend](/docs/secrets/databases/index.html) docs for more
information about setting up the Database Backend.

## Quick Start

After the Database Backend is mounted you can configure a MSSQL connection
by specifying this plugin as the `"plugin_name"` argument. Here is an example
configuration:

```
$ vault write database/config/mssql \
    plugin_name=mssql-database-plugin \
    connection_url='server=localhost;port=1433;user id=sa;password=Password!;database=AdventureWorks;app name=vault;' \
    allowed_roles="readonly"

The following warnings were returned from the Vault server:
* Read access to this endpoint should be controlled via ACLs as it will return the connection details as is, including passwords, if any.
```

In this case, we've configured Vault with the user "sa" and password "Password!",
connecting to an instance at "localhost" on port 1433. It is not necessary
that Vault has the sa login, but the user must have privileges to create
logins and manage processes. The fixed server roles `securityadmin` and
`processadmin` are examples of built-in roles that grant these permissions. The
user also must have privileges to create database users and grant permissions in
the databases that Vault manages.  The fixed database roles `db_accessadmin` and
`db_securityadmin` are examples or built-in roles that grant these permissions.


Once the MSSQL connection is configured we can add a role:

```
$ vault write database/roles/readonly \
    db_name=mssql \
    creation_statements="CREATE LOGIN [{{name}}] WITH PASSWORD = '{{password}}';\
        CREATE USER [{{name}}] FOR LOGIN [{{name}}];\
        GRANT SELECT ON SCHEMA::dbo TO [{{name}}];" \
    default_ttl="1h" \
    max_ttl="24h"

Success! Data written to: database/roles/readonly
```

This role can now be used to retrieve a new set of credentials by querying the
"database/creds/readonly" endpoint.

## Example for Azure SQL Database

Here is a complete example using Azure SQL Database. Note that databases in Azure SQL Database are [contained databases](https://docs.microsoft.com/en-us/sql/relational-databases/databases/contained-databases) and that we do not create a login for the user; instead, we associate the password directly with the user itself. Also note that you will need a separate connection and role for each Azure SQL database for which you want to generate dynamic credentials. You can use a single database backend mount for all these databases or use a separate mount for of them. In this example, we use a custom path for the database backend.

First, we mount a database backend at the azuresql path with `vault mount -path=azuresql database`. Then we configure a connection called "testvault" to connect to a database called "test-vault", using "azuresql" at the beginning of our path:

```
$ vault write azuresql/config/testvault \
    plugin_name=mssql-database-plugin \
    connection_url='server=hashisqlserver.database.windows.net;port=1433; \
    user id=admin;password=pAssw0rd;database=test-vault;app name=vault;' \
    allowed_roles="test"
```

Now we add a role called "test" for use with the "testvault" connection:

```
$ vault write azuresql/roles/test \
    db_name=testvault \
    creation_statements="CREATE USER [{{name}}] WITH PASSWORD = '{{password}}';" \
    revocation_statements="DROP USER IF EXISTS [{{name}}]" \
    default_ttl="1h" \
    max_ttl="24h"
```
We can now use this role to dynamically generate credentials for the Azure SQL database, test-vault:

```
$ vault read azuresql/creds/test
Key            	Value
---            	-----
lease_id       	azuresql/creds/test/2e5b1e0b-a081-c7e1-5622-39f58e79a719
lease_duration 	1h0m0s
lease_renewable	true
password       	A1a-48w04t1xzw1s33z3
username       	v-token-test-tr2t4x9pxvq1z8878s9s-1513446795
```

When we no longer need the backend, we can unmount it with `vault unmount azuresql`. Now, you can use the MSSQL Database Plugin with your Azure SQL databases.

## API

The full list of configurable options can be seen in the [MSSQL database
plugin API](/api/secret/databases/mssql.html) page.

For more information on the Database secret backend's HTTP API please see the [Database secret
backend API](/api/secret/databases/index.html) page.
