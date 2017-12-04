---
layout: "docs"
page_title: "MongoDB Database Plugin - Database Secret Backend"
sidebar_current: "docs-secrets-databases-mongodb"
description: |-
  The MongoDB plugin for Vault's Database backend generates database credentials to access MongoDB.
---

# MongoDB Database Plugin

Name: `mongodb-database-plugin`

The MongoDB Database Plugin is one of the supported plugins for the Database
backend. This plugin generates database credentials dynamically based on
configured roles for the MongoDB database.

See the [Database Backend](/docs/secrets/databases/index.html) docs for more
information about setting up the Database Backend.

## Quick Start

After the Database Backend is mounted you can configure a MongoDB connection
by specifying this plugin as the `"plugin_name"` argument. Here is an example
MongoDB configuration:

```
$ vault write database/config/mongodb \
    plugin_name=mongodb-database-plugin \
    allowed_roles="readonly" \
    connection_url="mongodb://admin:Password!@mongodb.acme.com:27017/admin?ssl=true"

The following warnings were returned from the Vault server:
* Read access to this endpoint should be controlled via ACLs as it will return the connection details as is, including passwords, if any.
```

Once the MongoDB connection is configured we can add a role:

```
$ vault write database/roles/readonly \
    db_name=mongodb \
    creation_statements='{ "db": "admin", "roles": [{ "role": "readWrite" }, {"role": "read", "db": "foo"}] }' \
    default_ttl="1h" \
    max_ttl="24h"

Success! Data written to: database/roles/readonly
```

This role can be used to retrieve a new set of credentials by querying the
"database/creds/readonly" endpoint.

## API

The full list of configurable options can be seen in the [MongoDB database
plugin API](/api/secret/databases/mongodb.html) page.

For more information on the Database secret backend's HTTP API please see the [Database secret
backend API](/api/secret/databases/index.html) page.
