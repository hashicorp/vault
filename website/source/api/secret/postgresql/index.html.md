---
layout: "api"
page_title: "PostgreSQL Secret Backend - HTTP API"
sidebar_current: "docs-http-secret-postgresql"
description: |-
  This is the API documentation for the Vault PostgreSQL secret backend.
---

# PostgreSQL Secret Backend HTTP API

~> **Deprecation Note:** This backend is deprecated in favor of the
combined databases backend added in v0.7.1. See the API documentation for
the new implementation of this backend at
[PostgreSQL Database Plugin HTTP API](/api/secret/databases/postgresql.html).

This is the API documentation for the Vault PostgreSQL secret backend. For
general information about the usage and operation of the PostgreSQL backend,
please see the
[Vault PostgreSQL backend documentation](/docs/secrets/postgresql/index.html).

This documentation assumes the PostgreSQL backend is mounted at the
`/postgresql` path in Vault. Since it is possible to mount secret backends at
any location, please update your API calls accordingly.

## Configure Connection

This endpoint configures the connection string used to communicate with
PostgreSQL.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `POST`   | `/postgresql/config/connection` | `204 (empty body)` |

### Parameters

- `connection_url` `(string: <required>)` – Specifies the PostgreSQL connection
  URL or PG-style string, for example `"user=foo host=bar"`.

- `max_open_connections` `(int: 2)` – Specifies the maximum number of open
  connections to the database. A negative value means unlimited.

- `max_idle_connections` `(int: 0)` – Specifies the maximum number of idle
  connections to the database. A zero uses the value of `max_open_connections`
  and a negative value disables idle connections. If this is larger than
  `max_open_connections` it will be reduced to be equal.

- `verify_connection` `(bool: true)` – Specifies if the connection is verified
  during initial configuration.

### Sample Payload

```json
{
  "connection_url": "postgresql://user:pass@localhost/my-db"
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    https://vault.rocks/v1/postgresql/config/connection
```

## Configure Lease

This configures the lease settings for generated credentials. If not configured,
leases default to 1 hour. This is a root protected endpoint.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `POST`   | `/postgresql/config/lease`   | `204 (empty body)`     |

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
    https://vault.rocks/v1/postgresql/config/lease
```

## Create Role

This endpoint creates or updates a role definition.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `POST`   | `/postgresql/roles/:name`    | `204 (empty body)`     |

### Parameters

- `name` `(string: <required>)` – Specifies the name of the role to create. This
  is specified as part of the URL.

- `sql` `(string: <required>)` – Specifies the SQL statements executed to create
  and configure the role. Must be a semicolon-separated string, a base64-encoded
  semicolon-separated string, a serialized JSON string array, or a
  base64-encoded serialized JSON string array. The '{{name}}', '{{password}}'
  and '{{expiration}}' values will be substituted.

- `revocation_sql` `(string: "")` – Specifies the SQL statements to be executed
  to revoke a user. Must be a semicolon-separated string, a base64-encoded
  semicolon-separated string, a serialized JSON string array, or a
  base64-encoded serialized JSON string array. The '{{name}}' value will be
  substituted.

### Sample Payload

```json
{
  "sql": "CREATE USER WITH ROLE {{name}}"
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    https://vault.rocks/v1/postgresql/roles/my-role
```

## Read Role

This endpoint queries the role definition.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `GET`    | `/postgresql/roles/:name`    | `200 application/json` |

### Parameters

- `name` `(string: <required>)` – Specifies the name of the role to read. This
  is specified as part of the URL.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    https://vault.rocks/v1/postgresql/roles/my-role
```

### Sample Response

```json
{
  "data": {
    "sql": "CREATE USER..."
  }
}
```

## List Roles

This endpoint returns a list of available roles. Only the role names are
returned, not any values.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `LIST`   | `/postgresql/roles`          | `200 application/json` |
| `GET`    | `/postgresql/roles?list=true` | `200 application/json` |

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request LIST \
    https://vault.rocks/v1/postgresql/roles
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
| `DELETE` | `/postgresql/roles/:name`    | `204 (empty body)`     |

### Parameters

- `name` `(string: <required>)` – Specifies the name of the role to delete. This
  is specified as part of the URL.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request DELETE \
    https://vault.rocks/v1/postgresql/roles/my-role
```

## Generate Credentials

This endpoint generates a new set of dynamic credentials based on the named
role.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `GET`    | `/postgresql/creds/:name`    | `200 application/json` |

### Parameters

- `name` `(string: <required>)` – Specifies the name of the role to create
  credentials against. This is specified as part of the URL.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    https://vault.rocks/v1/postgresql/creds/my-role
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
