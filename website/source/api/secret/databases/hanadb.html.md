---
layout: "api"
page_title: "HANA Database Plugin - HTTP API"
sidebar_current: "docs-http-secret-databases-hana"
description: |-
  The HANA plugin for Vault's Database backend generates database credentials to access HANA servers.
---

# HANA Database Plugin HTTP API

The HANA Database Plugin is one of the supported plugins for the Database
backend. This plugin generates database credentials dynamically based on
configured roles for the HANA database.

## Configure Connection

In addition to the parameters defined by the [Database
Backend](/api/secret/databases/index.html#configure-connection), this plugin
has a number of parameters to further configure a connection.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `POST`   | `/database/config/:name`     | `204 (empty body)` |

### Parameters
- `connection_url` `(string: <required>)` - Specifies the HANA DSN.

- `max_open_connections` `(int: 2)` - Specifies the maximum number of open
  connections to the database.

- `max_idle_connections` `(int: 0)` - Specifies the maximum number of idle
  connections to the database. A zero uses the value of `max_open_connections`
  and a negative value disables idle connections. If larger than
  `max_open_connections` it will be reduced to be equal.

- `max_connection_lifetime` `(string: "0s")` - Specifies the maximum amount of
  time a connection may be reused. If <= 0s connections are reused forever.

### Sample Payload

```json
{
  "plugin_name": "hana-database-plugin",
  "allowed_roles": "readonly",
  "connection_url": "hdb://username:password@localhost:1433",
  "max_open_connections": 5,
  "max_connection_lifetime": "5s",
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    https://vault.rocks/v1/database/config/hana
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
  semicolon-separated string, a base64-encoded semicolon-separated string, a
  serialized JSON string array, or a base64-encoded serialized JSON string
  array. The '{{name}}', '{{password}}', and '{{expiration}}' values will be
  substituted.
  - The expiration time will be HANA server time plus the role's `default_ttl`.
    If `default_ttl` is 0 or not set, a SQL HdbError 438 will be returned.

- `revocation_statements` `(string: "")` – Specifies the database statements to
  be executed to revoke a user. Must be a semicolon-separated string, a
  base64-encoded semicolon-separated string, a serialized JSON string array, or
  a base64-encoded serialized JSON string array. The '{{name}}' value will be
  substituted. If not provided, defaults to dropping the user only if they have
  no dependent objects.
