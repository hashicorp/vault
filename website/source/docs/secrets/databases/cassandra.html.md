---
layout: "docs"
page_title: "Cassandra Database Plugin - Database Secret Backend"
sidebar_current: "docs-secrets-databases-cassandra"
description: |-
  The Cassandra plugin for Vault's Database backend generates database credentials to access Cassandra.
---

# Cassandra Database Plugin

Name: `cassandra-database-plugin`

The Cassandra Database Plugin is one of the supported plugins for the Database
backend. This plugin generates database credentials dynamically based on
configured roles for the Cassandra database.

See the [Database Backend](/docs/secrets/databases/index.html) docs for more
information about setting up the Database Backend.

## Quick Start

After the Database Backend is mounted you can configure a cassandra connection
by specifying this plugin as the `"plugin_name"` argument. Here is an example
cassandra configuration:

```
$ vault write database/config/cassandra \
    plugin_name=cassandra-database-plugin \
    allowed_roles="readonly" \
    hosts=localhost \
    username=cassandra \
    password=cassandra

The following warnings were returned from the Vault server:
* Read access to this endpoint should be controlled via ACLs as it will return the connection details as is, including passwords, if any.
```

Once the cassandra connection is configured we can add a role:

```
$ vault write database/roles/readonly \
    db_name=cassandra \
    creation_statements="CREATE USER '{{username}}' WITH PASSWORD '{{password}}' NOSUPERUSER; \
         GRANT SELECT ON ALL KEYSPACES TO {{username}};" \
    default_ttl="1h" \
    max_ttl="24h"


Success! Data written to: database/roles/readonly
```

This role can be used to retrieve a new set of credentials by querying the
"database/creds/readonly" endpoint.

## API

The full list of configurable options can be seen in the [Cassandra database
plugin API](/api/secret/databases/cassandra.html) page.

For more information on the Database secret backend's HTTP API please see the [Database secret
backend API](/api/secret/databases/index.html) page.
