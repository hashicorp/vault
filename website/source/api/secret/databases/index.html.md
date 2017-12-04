---
layout: "api"
page_title: "Databases - HTTP API"
sidebar_current: "docs-http-secret-databases"
description: |-
  Top page for database secret backend information
---

# Database Secret Backend HTTP API

This is the API documentation for the Vault Database secret backend. For
general information about the usage and operation of the Database backend,
please see the
[Vault Database backend documentation](/docs/secrets/databases/index.html).

This documentation assumes the Database backend is mounted at the
`/database` path in Vault. Since it is possible to mount secret backends at
any location, please update your API calls accordingly.

## Configure Connection

This endpoint configures the connection string used to communicate with the
desired database. In addition to the parameters listed here, each Database
plugin has additional, database plugin specifig,  parameters for this endpoint.
Please read the HTTP API for the plugin you'd wish to configure to see the full
list of additional parameters.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `POST`   | `/database/config/:name`     | `204 (empty body)` |

### Parameters
- `name` `(string: <required>)` – Specifies the name for this database
  connection. This is specified as part of the URL.

- `plugin_name` `(string: <required>)` - Specifies the name of the plugin to use
  for this connection.

- `verify_connection` `(bool: true)` – Specifies if the connection is verified
  during initial configuration. Defaults to true.

- `allowed_roles` `(slice: [])` - Array or comma separated string of the roles
  allowed to use this connection. Defaults to empty (no roles), if contains a
  "*" any role can use this connection. 

### Sample Payload

```json
{
  "plugin_name": "mysql-database-plugin",
  "allowed_roles": "readonly",
  "connection_url": "root:mysql@tcp(127.0.0.1:3306)/"
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    https://vault.rocks/v1/database/config/mysql
```

## Read Connection

This endpoint returns the configuration settings for a connection.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `GET`    | `/database/config/:name`     | `200 application/json` |

### Parameters

- `name` `(string: <required>)` – Specifies the name of the connection to read.
  This is specified as part of the URL.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request GET \
    https://vault.rocks/v1/database/config/mysql
```

### Sample Response

```json
{
	"data": {
		"allowed_roles": [
			"readonly"
		],
		"connection_details": {
			"connection_url": "root:mysql@tcp(127.0.0.1:3306)/",
		},
		"plugin_name": "mysql-database-plugin"
	},
}
```

## Delete Connection

This endpoint deletes a connection.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `DELETE` | `/database/config/:name`     | `204 (empty body)`     |

### Parameters

- `name` `(string: <required>)` – Specifies the name of the connection to delete.
  This is specified as part of the URL.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request DELETE \
    https://vault.rocks/v1/database/config/mysql
```

## Reset Connection

This endpoint closes a connection and it's underlying plugin and restarts it
with the configuration stored in the barrier.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `POST`   | `/database/reset/:name`      | `204 (empty body)`     |

### Parameters

- `name` `(string: <required>)` – Specifies the name of the connection to delete.
  This is specified as part of the URL.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    https://vault.rocks/v1/database/reset/mysql
```

## Create Role

This endpoint creates or updates a role definition.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `POST`   | `/database/roles/:name`      | `204 (empty body)`     |

### Parameters

- `name` `(string: <required>)` – Specifies the name of the role to create. This
  is specified as part of the URL.

- `db_name` `(string: <required>)` - The name of the database connection to use
  for this role. 

- `default_ttl` `(string/int: 0)` - Specifies the TTL for the leases
  associated with this role. Accepts time suffixed strings ("1h") or an integer
  number of seconds. Defaults to system/backend default TTL time.

- `max_ttl` `(string/int: 0)` - Specifies the maximum TTL for the leases
  associated with this role. Accepts time suffixed strings ("1h") or an integer
  number of seconds. Defaults to system/backend default TTL time.

- `creation_statements` `(string: <required>)` – Specifies the database
  statements executed to create and configure a user. See the plugin's API page
  for more information on support and formatting for this parameter. 

- `revocation_statements` `(string: "")` – Specifies the database statements to
  be executed to revoke a user. See the plugin's API page for more information
  on support and formatting for this parameter. 

- `rollback_statements` `(string: "")` – Specifies the database statements to be
  executed rollback a create operation in the event of an error. Not every
  plugin type will support this functionality. See the plugin's API page for
  more information on support and formatting for this parameter. 

- `renew_statements` `(string: "")` – Specifies the database statements to be
  executed to renew a user. Not every plugin type will support this
  functionality. See the plugin's API page for more information on support and
  formatting for this parameter. 



### Sample Payload

```json
{
    "db_name": "mysql",
    "creation_statements": "CREATE USER '{{name}}'@'%' IDENTIFIED BY '{{password}}';GRANT SELECT ON *.* TO '{{name}}'@'%';",
    "default_ttl": "1h",
    "max_ttl": "24h"
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    https://vault.rocks/v1/database/roles/my-role
```

## Read Role

This endpoint queries the role definition.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `GET`    | `/database/roles/:name`    | `200 application/json` |

### Parameters

- `name` `(string: <required>)` – Specifies the name of the role to read. This
  is specified as part of the URL.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    https://vault.rocks/v1/database/roles/my-role
```

### Sample Response

```json
{
    "data": {
		"creation_statements": "CREATE ROLE \"{{name}}\" WITH LOGIN PASSWORD '{{password}}' VALID UNTIL '{{expiration}}';         GRANT SELECT ON ALL TABLES IN SCHEMA public TO \"{{name}}\";",
		"db_name": "mysql",
		"default_ttl": 3600,
		"max_ttl": 86400,
		"renew_statements": "",
		"revocation_statements": "",
		"rollback_statements": ""
	},
}
```

## List Roles

This endpoint returns a list of available roles. Only the role names are
returned, not any values.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `LIST`   | `/database/roles`          | `200 application/json` |
| `GET`   | `/database/roles?list=true` | `200 application/json` |

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request LIST \
    https://vault.rocks/v1/database/roles
```

### Sample Response

```json
{
  "auth": null,
  "data": {
    "keys": ["dev", "prod"]
  },
  "lease_duration": 2764800,
  "lease_id": "",
  "renewable": false
}
```

## Delete Role

This endpoint deletes the role definition.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `DELETE` | `/database/roles/:name`    | `204 (empty body)`     |

### Parameters

- `name` `(string: <required>)` – Specifies the name of the role to delete. This
  is specified as part of the URL.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request DELETE \
    https://vault.rocks/v1/database/roles/my-role
```

## Generate Credentials

This endpoint generates a new set of dynamic credentials based on the named
role.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `GET`    | `/database/creds/:name`    | `200 application/json` |

### Parameters

- `name` `(string: <required>)` – Specifies the name of the role to create
  credentials against. This is specified as part of the URL.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    https://vault.rocks/v1/database/creds/my-role
```

### Sample Response

```json
{
  "data": {
    "username": "root-1430158508-126",
    "password": "132ae3ef-5a64-7499-351e-bfe59f3a2a21"
  }
}
```
