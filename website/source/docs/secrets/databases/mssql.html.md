---
layout: "docs"
page_title: "MSSQL Database Plugin"
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
    connection_url='sqlserver://sa:yourStrong(!)Password@localhost:1433' \
    allowed_roles="readonly"

The following warnings were returned from the Vault server:
* Read access to this endpoint should be controlled via ACLs as it will return the connection details as is, including passwords, if any.
```

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

## API

The full list of configurable options can be seen in the [MSSQL database
plugin API](/api/secret/databases/mssql.html) page.

For more information on the Database secret backend's HTTP API please see the [Database secret
backend API](/api/secret/databases/index.html) page.

