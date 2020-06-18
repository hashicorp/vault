# MongoDB Atlas Database Secrets Engine

This plugin provides unique, short-lived credentials for [MongoDB Atlas](https://www.mongodb.com/cloud/atlas).
It is to be used with [Hashicorp Vault](https://www.github.com/hashicorp/vault).

**Please note**: We take Vault's security and our users' trust very seriously. If you believe you have found a security issue in Vault, _please responsibly disclose_ by contacting us at [security@hashicorp.com](mailto:security@hashicorp.com).

## Quick Links

- [Database Secrets Engine for MongoDB Atlas](https://www.vaultproject.io/docs/secrets/databases/mongodbatlas.html)
- [MongoDB Atlas Website](https://www.mongodb.com/cloud/atlas)
- [Vault Website](https://www.vaultproject.io)
- [Vault Github](https://www.github.com/hashicorp/vault)

## Getting Started

This is a Vault plugin and is meant to work with Vault. This guide assumes you have already installed Vault
and have a basic understanding of how Vault works.

Otherwise, first read this guide on how to [get started with Vault](https://www.vaultproject.io/intro/getting-started/install.html).

To learn specifically about how plugins work, see documentation on [Vault plugins](https://www.vaultproject.io/docs/internals/plugins.html).

## Installation

This plugin is bundled in Vault version 1.4.0 or later. It may also be built and mounted externally
with earlier versions of Vault. For details on this process please see the documentation for Vault's
[plugin system](https://www.vaultproject.io/docs/internals/plugins.html).

## Setup

1. Enable the database secrets engine if it is not already enabled:

    ```text
    $ vault secrets enable database
    Success! Enabled the database secrets engine at: database/
    ```

    The secrets engine will be enabled at the default path which is name of the engine. To
    enable the secrets engine at a different path use the `-path` argument.

1. Configure Vault with the proper plugin and connection information:

    ```text
    $ vault write database/config/my-mongodbatlas-database \
        plugin_name=mongodbatlas-database-plugin \
        allowed_roles="my-role" \
        public_key="a-public-key" \
        private_key="a-private-key!" \
        project_id="a-project-id"
    ```

2. Configure a role that maps a name in Vault to a MongoDB Atlas command that executes and
   creates the Database User credential:

    ```text
    $ vault write database/roles/my-role \
        db_name=my-mongodbatlas-database \
        creation_statements='{ "database_name": "admin", "roles": [{"databaseName":"admin","roleName":"atlasAdmin"}]}' \
        default_ttl="1h" \
        max_ttl="24h"
    Success! Data written to: database/roles/my-role
    ```

## Usage

After the secrets engine is configured and a user/machine has a Vault token with
the proper permissions, it can generate credentials.

1. Generate a new credential by reading from the `/creds` endpoint with the name
of the role:

    ```text
    $ vault read database/creds/my-role
    Key                Value
    ---                -----
    lease_id           database/creds/my-role/2f6a614c-4aa2-7b19-24b9-ad944a8d4de6
    lease_duration     1h
    lease_renewable    true
    password           A1a-QwxApKgnfCp1AJYN
    username           v-5WFTBKdwOTLOqWLgsjvH-1565815206
    ```


For more details on configuring and using the plugin, refer to the [Database Secrets Engine for MongoDB Atlas](https://www.vaultproject.io/docs/secrets/databases/mongodbatlas.html)
documentation.
