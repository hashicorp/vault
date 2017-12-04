---
layout: "api"
page_title: "RabbitMQ Secret Backend - HTTP API"
sidebar_current: "docs-http-secret-rabbitmq"
description: |-
  This is the API documentation for the Vault RabbitMQ secret backend.
---

# RabbitMQ Secret Backend HTTP API

This is the API documentation for the Vault RabbitMQ secret backend. For general
information about the usage and operation of the RabbitMQ backend, please see
the [Vault RabbitMQ backend documentation](/docs/secrets/rabbitmq/index.html).

This documentation assumes the RabbitMQ backend is mounted at the `/rabbitmq`
path in Vault. Since it is possible to mount secret backends at any location,
please update your API calls accordingly.

## Configure Connection

This endpoint configures the connection string used to communicate with
RabbitMQ.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `POST`   | `/rabbitmq/config/connection` | `204 (empty body)` |

### Parameters

- `connection_uri` `(string: <required>)` – Specifies the RabbitMQ connection
  URI.

- `username` `(string: <required>)` – Specifies the RabbitMQ management
  administrator username.

- `password` `(string: <required>)` – Specifies the RabbitMQ management
  administrator password.

- `verify_connection` `(bool: true)` – Specifies whether to verify connection
  URI, username, and password.

### Sample Payload

```json
{
  "connection_uri": "https://...",
  "username": "user",
  "password": "password"
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    https://vault.rocks/v1/rabbitmq/config/connection
```

## Configure Lease

This endpoint configures the lease settings for generated credentials.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `POST`   | `/rabbitmq/config/lease`     | `204 (empty body)` |

### Parameters

- `ttl` `(int: 0)` – Specifies the lease ttl provided in seconds.

- `max_ttl` `(int: 0)` – Specifies the maximum ttl provided in seconds.

### Sample Payload

```json
{
  "ttl": 1800,
  "max_ttl": 3600
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    https://vault.rocks/v1/rabbitmq/config/lease
```

## Create Role

This endpoint creates or updates the role definition.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `POST`   | `/rabbitmq/roles/:name`      | `204 (empty body)`     |

### Parameters

- `name` `(string: <required>)` – Specifies the name of the role to create. This
  is specified as part of the URL.

- `tags` `(string: "")` – Specifies a comma-separated RabbitMQ management tags.

- `vhost` `(string: "")` – Specifies a map of virtual hosts to
  permissions.

### Sample Payload

```json
{
  "tags": "tag1,tag2",
  "vhost": "{\"/\": {\"configure\":\".*\", \"write\":\".*\", \"read\": \".*\"}}"
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    https://vault.rocks/v1/rabbitmq/roles/my-role
```

## Read Role

This endpoint queries the role definition.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `GET`    | `/rabbitmq/roles/:name`      | `200 application/json` |

### Parameters

- `name` `(string: <required>)` – Specifies the name of the role to read. This
  is specified as part of the URL.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    https://vault.rocks/v1/rabbitmq/roles/my-role
```

### Sample Response

```json
{
  "data": {
    "tags": "",
    "vhost": "{\"/\": {\"configure\":\".*\", \"write\":\".*\", \"read\": \".*\"}}"
  }
}
```

## Delete Role

This endpoint deletes the role definition.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `DELETE` | `/rabbitmq/roles/:name`     | `204 (empty body)`     |

### Parameters

- `name` `(string: <required>)` – Specifies the name of the role to delete. This
  is specified as part of the URL.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request DELETE \
    https://vault.rocks/v1/rabbitmq/roles/my-role
```

## Generate Credentials

This endpoint generates a new set of dynamic credentials based on the named
role.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `GET`    | `/rabbitmq/creds/:name`      | `200 application/json` |

### Parameters

- `name` `(string: <required>)` – Specifies the name of the role to create
  credentials against. This is specified as part of the URL.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    https://vault.rocks/v1/rabbitmq/creds/my-role
```

### Sample Response

```json
{
  "data": {
    "username": "root-4b95bf47-281d-dcb5-8a60-9594f8056092",
    "password": "e1b6c159-ca63-4c6a-3886-6639eae06c30"
  }
}
```
