---
layout: "api"
page_title: "MongoDB Database Plugin - HTTP API"
sidebar_current: "docs-http-secret-databases-mongodb"
description: |-
  The MongoDB plugin for Vault's Database backend generates database credentials to access MongoDB servers.
---

# MongoDB Database Plugin HTTP API

The MongoDB Database Plugin is one of the supported plugins for the Database
backend. This plugin generates database credentials dynamically based on
configured roles for the MongoDB database.

## Configure Connection

In addition to the parameters defined by the [Database
Backend](/api/secret/databases/index.html#configure-connection), this plugin
has a number of parameters to further configure a connection.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `POST`   | `/database/config/:name`     | `204 (empty body)` |

### Parameters
- `connection_url` `(string: <required>)` – Specifies the MongoDB standard connection string (URI).

### Sample Payload

```json
{
  "plugin_name": "mongodb-database-plugin",
  "allowed_roles": "readonly",
  "connection_url": "mongodb://admin:Password!@mongodb.acme.com:27017/admin?ssl=true"
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    https://vault.rocks/v1/database/config/mongodb
```
