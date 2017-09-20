---
layout: "docs"
page_title: "Database - Secrets Engines"
sidebar_current: "docs-secrets-databases"
description: |-
  The database secrets engine generates database credentials dynamically based
  on configured roles. It works with a number of different databases through a
  plugin interface. There are a number of builtin database types and an exposed
  framework for running custom database types for extendability.
---

# Databases

The database secrets engine generates database credentials dynamically based on
configured roles. It works with a number of different databases through a plugin
interface. There are a number of builtin database types and an exposed framework
for running custom database types for extendability. This means that services
that need to access a database no longer need to hardcode credentials: they can
request them from Vault, and use Vault's leasing mechanism to more easily roll
keys.

Since every service accessing the database with unique credentials, it makes
auditing much easier when questionable data access is discovered. You can track
it down to the specific instance of a service based on the SQL username.

Vault makes use of its own internal revocation system to ensure that users
become invalid within a reasonable time of the lease expiring.

## Setup

Most secrets engines must be configured in advance before they can perform their
functions. These steps are usually completed by an operator or configuration
management tool.

1. Enable the database secrets engine:

    ```text
    $ vault secrets enable database
    Success! Enabled the database secrets engine at: database/
    ```

    By default, the secrets engine will enable at the name of the engine. To
    enable the secrets engine at a different path, use the `-path` argument.

1. Configure Vault with the proper plugin and connection information:

    ```text
    $ vault write database/config/my-database \
        plugin_name="..." \
        connection_url="..." \
        allowed_roles="..."
    ```

    This secrets engine can configure multiple database connections. For details
    on the specific configuration options, please see the database-specific
    documentation.

1. Configure a role that maps a name in Vault to an SQL statement to execute to create the database credential:

    ```text
    $ vault write database/roles/my-role \
        db_name=my-database \
        creation_statements="..." \
        default_ttl="1h" \
        max_ttl="24h"
    Success! Data written to: database/roles/my-role
    ```

    The `{{name}}` and `{{password}}` fields will be populated by the plugin
    with dynamically generated values. In some plugins the `{{expiration}}`
    field is also be supported.

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

## Custom Plugins

This secrets engine allows custom database types to be run through the exposed
plugin interface. Please see the [custom database
plugin](/docs/secrets/databases/custom.html) for more information.

## API

The database secrets engine has a full HTTP API. Please see the [Database secret
secrets engine API](/api/secret/databases/index.html) for more details.
