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

## API

The full list of configurable options can be seen in the [MSSQL database
plugin API](/api/secret/databases/mssql.html) page.

For more information on the database secrets engine's HTTP API please see the
[Database secrets engine API](/api/secret/databases/index.html) page.
