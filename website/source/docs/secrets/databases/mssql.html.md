---
layout: "docs"
page_title: "MSSQL - Database - Secrets Engines"
sidebar_current: "docs-secrets-databases-mssql"
description: |-

  MSSQL is one of the supported plugins for the database secrets engine. This
  plugin generates database credentials dynamically based on configured roles
  for the MSSQL database.
---

# MSSQL Database Secrets Engine

MSSQL is one of the supported plugins for the database secrets engine. This
plugin generates database credentials dynamically based on configured roles for
the MSSQL database.

See the [database secrets engine](/docs/secrets/databases/index.html) docs for
more information about setting up the database secrets engine.

## Setup

1. Enable the database secrets engine if it is not already enabled:

    ```text
    $ vault secrets enable database
    Success! Enabled the database secrets engine at: database/
    ```

    By default, the secrets engine will enable at the name of the engine. To
    enable the secrets engine at a different path, use the `-path` argument.

1. Configure Vault with the proper plugin and connection information:

    ```text
    $ vault write database/config/my-mssql-database \
        plugin_name=mssql-database-plugin \
        connection_url='sqlserver://sa:yourStrong(!)Password@localhost:1433' \
        allowed_roles="my-role"
    ```

    In this case, we've configured Vault with the user "sa" and password
    "Password!", connecting to an instance at "localhost" on port 1433. It is
    not necessary that Vault has the sa login, but the user must have privileges
    to create logins and manage processes. The fixed server roles
    `securityadmin` and `processadmin` are examples of built-in roles that grant
    these permissions. The user also must have privileges to create database
    users and grant permissions in the databases that Vault manages.  The fixed
    database roles `db_accessadmin` and `db_securityadmin` are examples or
    built-in roles that grant these permissions.

1. Configure a role that maps a name in Vault to an SQL statement to execute to
create the database credential:

    ```text
    $ vault write database/roles/my-role \
        db_name=my-mssql-database \
        creation_statements="CREATE LOGIN [{{name}}] WITH PASSWORD = '{{password}}';\
            CREATE USER [{{name}}] FOR LOGIN [{{name}}];\
            GRANT SELECT ON SCHEMA::dbo TO [{{name}}];" \
        default_ttl="1h" \
        max_ttl="24h"
    Success! Data written to: database/roles/my-role
    ```

## Usage

After the secrets engine is configured and a user/machine has a Vault token with
the proper permission, it can generate credentials.

1. Generate a new credential by reading from the `/creds` endpoint with the name
of the role:

    ```text
    $ vault read database/creds/my-role
    Key                Value
    ---                -----
    lease_id           database/creds/my-role/2f6a614c-4aa2-7b19-24b9-ad944a8d4de6
    lease_duration     1h
    lease_renewable    true
    password           8cab931c-d62e-a73d-60d3-5ee85139cd66
    username           v-root-e2978cd0-
    ```

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

For more information on the database secrets engine's HTTP API please see the
[Database secrets engine API](/api/secret/databases/index.html) page.
