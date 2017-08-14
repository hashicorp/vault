---
layout: "docs"
page_title: "HANA Database Plugin - Database Secret Backend"
sidebar_current: "docs-secrets-databases-HANA"
description: |-
  The HANA plugin for Vault's Database backend generates database credentials to access SAP HANA Database.
---

# HANA Database Plugin

Name: `hana-database-plugin`

The HANA Database Plugin is one of the supported plugins for the Database
backend. This plugin generates database credentials dynamically based on
configured roles for the HANA database.

See the [Database Backend](/docs/secrets/databases/index.html) docs for more
information about setting up the Database Backend.

## Quick Start

After the Database Backend is mounted you can configure a HANA connection
by specifying this plugin as the `"plugin_name"` argument. Here is an example
configuration:

```
$ vault write database/config/hana \
    plugin_name=hana-database-plugin \
    connection_url="hdb://username:password@localhost:1433" \
    allowed_roles="readonly"

The following warnings were returned from the Vault server:
* Read access to this endpoint should be controlled via ACLs as it will
return the connection details as is, including passwords, if any.
```

Once the HANA connection is configured we can add a role:

```
$ vault write database/roles/readonly \
    db_name=hana \
    creation_statements="CREATE USER {{name}} PASSWORD {{password}} VALID UNTIL '{{expiration}}';\
        CALL GRANT_ACTIVATED_ROLE ( 'sap.hana.admin.roles::Monitoring', '{{name}}' );" \
    default_ttl="12h" \
    max_ttl="24h"

Success! Data written to: database/roles/readonly
```

This role can now be used to retrieve a new set of credentials by querying the
"database/creds/readonly" endpoint.

## API

The full list of configurable options can be seen in the [HANA database
plugin API](/api/secret/databases/hanadb.html) page.

For more information on the Database secret backend's HTTP API please see the [Database secret
backend API](/api/secret/databases/index.html) page.
