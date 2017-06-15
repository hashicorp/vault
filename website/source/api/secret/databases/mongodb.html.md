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

## Statements

Statements are configured during role creation and are used by the plugin to
determine what is sent to the datatabse on user creation, renewing, and
revocation. For more information on configuring roles see the [Role
API](/api/secret/databases/index.html#create-role) in the Database Backend docs.

### Parameters

The following are the statements used by this plugin. If not mentioned in this
list the plugin does not support that statement type.

- `creation_statements` `(string: <required>)` – Specifies the database
  statements executed to create and configure a user. Must be a
  serialized JSON object, or a base64-encoded serialized JSON object. 
  The object can optionally contain a "db" string for session connection,
  and must contain a "roles" array. This array contains objects that holds
  a "role", and an optional "db" value, and is similar to the BSON document that
  is accepted by MongoDB's `roles` field. Vault will transform this array into
  such format. For more information regarding the `roles` field, refer to
  [MongoDB's documentation](https://docs.mongodb.com/manual/reference/method/db.createUser/).

- `revocation_statements` `(string: "")` – Specifies the database statements to
  be executed to revoke a user. Must be a serialized JSON object, or a base64-encoded 
  serialized JSON object. The object can optionally contain a "db" string. If no
  "db" value is provided, it defaults to the "admin" database.

### Sample Creation Statement

```json
{
  "db": "admin",
  "roles": [
    {
      "role": "read",
      "db": "foo",
    }
  ]
}
```