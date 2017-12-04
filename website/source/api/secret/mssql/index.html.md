---
layout: "api"
page_title: "MSSQL Secret Backend - HTTP API"
sidebar_current: "docs-http-secret-mssql"
description: |-
  This is the API documentation for the Vault MSSQL secret backend.
---

# MSSQL Secret Backend HTTP API

~> **Deprecation Note:** This backend is deprecated in favor of the
combined databases backend added in v0.7.1. See the API documentation for
the new implementation of this backend at
[MSSQL Database Plugin HTTP API](/api/secret/databases/mssql.html).

This is the API documentation for the Vault MSSQL secret backend. For general
information about the usage and operation of the MSSQL backend, please see
the [Vault MSSQL backend documentation](/docs/secrets/mssql/index.html).

This documentation assumes the MSSQL backend is mounted at the `/mssql`
path in Vault. Since it is possible to mount secret backends at any location,
please update your API calls accordingly.

## Configure Connection

This endpoint configures the connection DSN used to communicate with Microsoft
SQL Server.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `POST`   | `/mssql/config/connection`   | `204 (empty body)`     |

### Parameters

- `connection_string` `(string: <required>)` – Specifies the MSSQL DSN.

- `max_open_connections` `(int: 2)` – Specifies the maximum number of open
  connections to the database.

- `max_idle_connections` `(int: 0)` – Specifies the maximum number of idle
  connections to the database. A zero uses the value of `max_open_connections`
  and a negative value disables idle connections. If larger than
  `max_open_connections` it will be reduced to be equal.

### Sample Payload

```json
{
  "connection_string": "Server=myServerAddress;Database=myDataBase;User Id=myUsername; Password=myPassword;"
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    https://vault.rocks/v1/mssql/config/connection
```

## Configure Lease

This endpoint configures the lease settings for generated credentials.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `POST`   | `/mysql/config/lease`        | `204 (empty body)`     |

### Parameters

- `lease` `(string: <required>)` – Specifies the lease value provided as a
  string duration with time suffix. "h" (hour) is the largest suffix.

- `lease_max` `(string: <required>)` – Specifies the maximum lease value
  provided as a string duration with time suffix. "h" (hour) is the largest
  suffix.

### Sample Payload

```json
{
  "lease": "12h",
  "lease_max": "24h"
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    https://vault.rocks/v1/mssql/config/lease
```

## Create Role

This endpoint creates or updates the role definition.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `POST`   | `/mssql/roles/:name`         | `204 (empty body)`     |

### Parameters

- `sql` `(string: <required>)` – Specifies the SQL statements executed to create
  and configure the role.  The '{{name}}' and '{{password}}' values will be
  substituted. Must be a semicolon-separated string, a base64-encoded
  semicolon-separated string, a serialized JSON string array, or a
  base64-encoded serialized JSON string array.

### Sample Payload

```json
{
  "sql": "CREATE LOGIN ..."
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    https://vault.rocks/v1/mssql/roles/my-role
```

## Read Role

This endpoint queries the role definition.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `GET`    | `/mssql/roles/:name`         | `200 application/json` |

### Parameters

- `name` `(string: <required>)` – Specifies the name of the role to read. This
  is specified as part of the URL.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    https://vault.rocks/v1/mssql/roles/my-role
```

### Sample Response

```json
{
  "data": {
    "sql": "CREATE LOGIN..."
  }
}
```

## List Roles

This endpoint returns a list of available roles. Only the role names are
returned, not any values.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `LIST`   | `/mssql/roles`               | `200 application/json` |
| `GET`   | `/mssql/roles?list=true`      | `200 application/json` |

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request LIST \
    https://vault.rocks/v1/mssql/roles
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
| `DELETE` | `/mssql/roles/:name`         | `204 (empty body)`     |

### Parameters

- `name` `(string: <required>)` – Specifies the name of the role to delete. This
  is specified as part of the URL.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request DELETE \
    https://vault.rocks/v1/mssql/roles/my-role
```

## Generate Credentials

This endpoint generates a new set of dynamic credentials based on the named
role.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `GET`    | `/mssql/creds/:name`         | `200 application/json` |

### Parameters

- `name` `(string: <required>)` – Specifies the name of the role to create
  credentials against. This is specified as part of the URL.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    https://vault.rocks/v1/mssql/creds/my-role
```

### Sample Response

```json
{
  "data": {
    "username": "root-a147d529-e7d6-4a16-8930-4c3e72170b19",
    "password": "ee202d0d-e4fd-4410-8d14-2a78c5c8cb76"
  }
}
```
