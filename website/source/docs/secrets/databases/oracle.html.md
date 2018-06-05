---
layout: "docs"
page_title: "Oracle - Database - Secrets Engines"
sidebar_current: "docs-secrets-databases-oracle"
description: |-
  Oracle is one of the supported plugins for the database secrets engine. This
  plugin generates database credentials dynamically based on configured roles
  for the Oracle database.
---

# Oracle Database Secrets Engine

Oracle is one of the supported plugins for the database secrets engine. This
plugin generates database credentials dynamically based on configured roles for
the Oracle database.

The Oracle database plugin is not bundled in the core Vault code tree and can be
found at its own git repository here:
[hashicorp/vault-plugin-database-oracle](https://github.com/hashicorp/vault-plugin-database-oracle)

See the [Database Backend](/docs/secrets/databases/index.html) docs for more
information about setting up the Database Backend.

## Setup

The Oracle Database Plugin does not live in the core Vault code tree and can be found
at its own git repository here: [hashicorp/vault-plugin-database-oracle](https://github.com/hashicorp/vault-plugin-database-oracle)

Before running the plugin you will need to have the the Oracle Instant Client
library installed. These can be downloaded from Oracle. The libraries will need to
be placed in the default library search path or somewhere defined in the
`LD_LIBRARY_PATH` environment variable.

1. Enable the database secrets engine if it is not already enabled:

    ```text
    $ vault secrets enable database
    Success! Enabled the database secrets engine at: database/
    ```

    By default, the secrets engine will enable at the name of the engine. To
    enable the secrets engine at a different path, use the `-path` argument.

1. Download and register the plugin:

    ```text
    $ vault write sys/plugins/catalog/oracle-database-plugin \
        sha_256="..." \
        command=oracle-database-plugin
    ```

1. Configure Vault with the proper plugin and connection information:

    ```text
    $ vault write database/config/my-oracle-database \
        plugin_name=oracle-database-plugin \
        connection_url="system/Oracle@localhost:1521/OraDoc.localhost" \
        allowed_roles="my-role"
    ```

1. Configure a role that maps a name in Vault to an SQL statement to execute to
create the database credential:

    ```text
    $ vault write database/roles/my-role \
        db_name=my-oracle-database \
        creation_statements="CREATE USER {{name}} IDENTIFIED BY {{password}}; GRANT CONNECT TO {{name}}; GRANT CREATE SESSION TO {{name}};" \
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

The full list of configurable options can be seen in the [Oracle database plugin
API](/api/secret/databases/oracle.html) page.

For more information on the database secrets engine's HTTP API please see the
[Database secrets engine API](/api/secret/databases/index.html) page.
