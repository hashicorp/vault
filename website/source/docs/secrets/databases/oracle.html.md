---
layout: "docs"
page_title: "Oracle Database Plugin"
sidebar_current: "docs-secrets-databases-oracle"
description: |-
  The Oracle Database plugin for Vault's Database backend generates database credentials to access Oracle Database severs.
---

# Oracle Database Plugin

Name: `vault-plugin-database-oracle`

The Oracle Database Plugin is an external plugin for the Database
backend. This plugin generates database credentials dynamically based on
configured roles for the Oracle database.

See the [Database Backend](/docs/secrets/databases/index.html) docs for more
information about setting up the Database Backend.

## Installation

The Oracle Database Plugin does not live in the core Vault code tree and can be found
at its own git repository here: [hashicorp/vault-plugin-database-oracle](https://github.com/hashicorp/vault-plugin-database-oracle)

Before running the plugin you will need to have the the Oracle Instant Client
library installed. These can be downloaded from Oracle. The libraries will need to
be placed in the default library search path or somewhere defined in the
`LD_LIBRARY_PATH` environment variable.

## Quick Start

After the Database Backend is mounted you can run the plugin and configure a
connection to the Oracle Database.

First the plugin must be built and registered to Vault's plugin catalog. To
build the plugin see the plugin's code repository. Once the plugin is built and
the binary is placed in Vault's plugin directory the catalog should be updated:

```
$ vault write sys/plugins/catalog/vault-plugin-database-oracle \
    sha_256=<expected SHA256 value> \
    command=vault-plugin-database-oracle
```

Once the plugin exists in the plugin catalog the Database backend can configure
a connection for the Oracle Database:

```
$ vault write database/config/oracle \
    plugin_name=vault-plugin-database-oracle \
    connection_url="system/Oracle@localhost:1521/OraDoc.localhost" \
    allowed_roles="readonly"

The following warnings were returned from the Vault server:
* Read access to this endpoint should be controlled via ACLs as it will return the connection details as is, including passwords, if any.
```

Once the Oracle connection is configured we can add a role:

```
$ vault write database/roles/readonly \
    db_name=oracle \
    creation_statements="CREATE USER {{name}} IDENTIFIED BY {{password}}; GRANT CONNECT TO {{name}}; GRANT CREATE SESSION TO {{name}};" \
    default_ttl="1h" \
    max_ttl="24h"

Success! Data written to: database/roles/readonly
```

This role can now be used to retrieve a new set of credentials by querying the
"database/creds/readonly" endpoint.

## API

The full list of configurable options can be seen in the [Oracle database
plugin API](/api/secret/databases/oracle.html) page.

For more information on the Database secret backend's HTTP API please see the [Database secret
backend API](/api/secret/databases/index.html) page.

