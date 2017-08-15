---
layout: "docs"
page_title: "PostgreSQL Database Plugin - Database Secret Backend"
sidebar_current: "docs-secrets-databases-postgresql"
description: |-
  The PostgreSQL plugin for Vault's Database backend generates database credentials to access PostgreSQL.
---

# PostgreSQL Database Plugin

Name: `postgresql-database-plugin`

The PostgreSQL Database Plugin is one of the supported plugins for the Database
backend. This plugin generates database credentials dynamically based on
configured roles for the PostgreSQL database.

See the [Database Backend](/docs/secrets/databases/index.html) docs for more
information about setting up the Database Backend.

## Quick Start

After the Database Backend is mounted you can configure a PostgreSQL connection
by specifying this plugin as the `"plugin_name"` argument. Here is an example
configuration:

```
$ vault write database/config/postgresql \
    plugin_name=postgresql-database-plugin \
    allowed_roles="readonly" \
    connection_url="postgresql://root:root@localhost:5432/"

The following warnings were returned from the Vault server:
* Read access to this endpoint should be controlled via ACLs as it will return the connection details as is, including passwords, if any.
```

Once the PostgreSQL connection is configured we can add a role. The PostgreSQL
plugin replaces `{{expiration}}` in statements with a formated timestamp:

```
$ vault write database/roles/readonly \
    db_name=postgresql \
    creation_statements="CREATE ROLE \"{{name}}\" WITH LOGIN PASSWORD '{{password}}' VALID UNTIL '{{expiration}}'; \
        GRANT SELECT ON ALL TABLES IN SCHEMA public TO \"{{name}}\";" \
    default_ttl="1h" \
    max_ttl="24h"

Success! Data written to: database/roles/readonly
```

This role can be used to retrieve a new set of credentials by querying the
"database/creds/readonly" endpoint.

## API

The full list of configurable options can be seen in the [PostgreSQL database
plugin API](/api/secret/databases/postgresql.html) page.

For more information on the Database secret backend's HTTP API please see the [Database secret
backend API](/api/secret/databases/index.html) page.
