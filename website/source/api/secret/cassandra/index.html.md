---
layout: "api"
page_title: "Cassandra Secret Backend - HTTP API"
sidebar_current: "docs-http-secret-cassandra"
description: |-
  This is the API documentation for the Vault Cassandra secret backend.
---

# Cassandra Secret Backend HTTP API

~> **Deprecation Note:** This backend is deprecated in favor of the
combined databases backend added in v0.7.1. See the API documentation for
the new implementation of this backend at
[Cassandra Database Plugin HTTP API](/api/secret/databases/cassandra.html).

This is the API documentation for the Vault Cassandra secret backend. For
general information about the usage and operation of the Cassandra backend,
please see the
[Vault Cassandra backend documentation](/docs/secrets/cassandra/index.html).

This documentation assumes the Cassandra backend is mounted at the `/cassandra`
path in Vault. Since it is possible to mount secret backends at any location,
please update your API calls accordingly.

## Configure Connection

This endpoint configures the connection information used to communicate with
Cassandra.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `POST`   | `/cassandra/config/connection` | `204 (empty body)`   |

### Parameters

- `hosts` `(string: <required>)` – Specifies a set of comma-delineated Cassandra
  hosts to connect to.

- `username` `(string: <required>)` – Specifies the username to use for
  superuser access.

- `password` `(string: <required>)` – Specifies the password corresponding to
  the given username.

- `tls` `(bool: true)` – Specifies whether to use TLS when connecting to
  Cassandra.

- `insecure_tls` `(bool: false)` – Specifies whether to skip verification of the
  server certificate when using TLS.

- `pem_bundle` `(string: "")` – Specifies concatenated PEM blocks containing a
  certificate and private key; a certificate, private key, and issuing CA
  certificate; or just a CA certificate.

- `pem_json` `(string: "")` – Specifies JSON containing a certificate and
  private key; a certificate, private key, and issuing CA certificate; or just a
  CA certificate. For convenience format is the same as the output of the
  `issue` command from the `pki` backend; see
  [the pki documentation](/docs/secrets/pki/index.html).

- `protocol_version` `(int: 2)` – Specifies the CQL protocol version to use.

- `connect_timeout` `(string: "5s")` – Specifies the connection timeout to use.

TLS works as follows:

- If `tls` is set to true, the connection will use TLS; this happens
  automatically if `pem_bundle`, `pem_json`, or `insecure_tls` is set

- If `insecure_tls` is set to true, the connection will not perform verification
  of the server certificate; this also sets `tls` to true

- If only `issuing_ca` is set in `pem_json`, or the only certificate in
  `pem_bundle` is a CA certificate, the given CA certificate will be used for
  server certificate verification; otherwise the system CA certificates will be
  used

- If `certificate` and `private_key` are set in `pem_bundle` or `pem_json`,
  client auth will be turned on for the connection

`pem_bundle` should be a PEM-concatenated bundle of a private key + client
certificate, an issuing CA certificate, or both. `pem_json` should contain the
same information; for convenience, the JSON format is the same as that output by
the issue command from the PKI backend.

### Sample Payload

```json
{
  "hosts": "cassandra1.local",
  "username": "user",
  "password": "pass"
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    https://vault.rocks/v1/cassandra/config/connection
```

## Create Role

This endpoint creates or updates the role definition.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `POST`   | `/cassandra/roles/:name`     | `204 (empty body)`     |

### Parameters

- `creation_cql` `(string: "")` – Specifies the CQL statements executed to
  create and configure the new user. Must be a semicolon-separated string, a
  base64-encoded semicolon-separated string, a serialized JSON string array, or
  a base64-encoded serialized JSON string array. The '{{username}}' and
  '{{password}}' values will be substituted; it is required that these
  parameters are in single quotes. The default creates a non-superuser user with
  no authorization grants.

- `rollback_cql` `(string: "")` – Specifies the CQL statements executed to
  attempt a rollback if an error is encountered during user creation. The
  default is to delete the user. Must be a semicolon-separated string, a
  base64-encoded semicolon-separated string, a serialized JSON string array, or
  a base64-encoded serialized JSON string array. The '{{username}}' and
  '{{password}}' values will be substituted; it is required that these
  parameters are in single quotes.

- `lease` `(string: "")` – Specifies the lease value provided as a string
  duration with time suffix. "h" hour is the largest suffix.

- `consistency` `(string: "Quorum")` – Specifies the consistency level value
  provided as a string. Determines the consistency level used for operations
  performed on the Cassandra database.

### Sample Payload

```json
{
  "creation_cql": "CREATE USER ..."
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    https://vault.rocks/v1/cassandra/roles/my-role
```

## Read Role

This endpoint queries the role definition.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `GET`    | `/cassandra/roles/:name`     | `200 application/json` |

### Parameters

- `name` `(string: <required>)` – Specifies the name of the role to read. This
  is part of the request URL.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    https://vault.rocks/v1/cassandra/roles/my-role
```

### Sample Response

```json
{
  "data": {
    "creation_cql": "CREATE USER...",
    "rollback_cql": "DROP USER...",
    "lease": "12h",
    "consistency": "Quorum"
  }
}
```

## Delete Role

This endpoint deletes the role definition.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `DELETE` | `/cassandra/roles/:name`     | `204 (no body)`        |

### Parameters

- `name` `(string: <required>)` – Specifies the name of the role to delete. This
  is part of the request URL.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request DELETE \
    https://vault.rocks/v1/cassandra/roles/my-role
```

## Generate Credentials

This endpoint generates a new set of dynamic credentials based on the named
role.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `GET`    | `/cassandra/creds/:name`     | `200 application/json` |

### Parameters

- `name` `(string: <required>)` – Specifies the name of the role to create
  credentials against. This is part of the request URL.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    https://vault.rocks/v1/cassandra/creds/my-role
```

### Sample Response

```json
{
  "data": {
    "username": "vault-root-1430158508-126",
    "password": "132ae3ef-5a64-7499-351e-bfe59f3a2a21"
  }
}
```
