---
layout: "api"
page_title: "Cassandra Database Plugin - HTTP API"
sidebar_current: "docs-http-secret-databases-cassandra"
description: |-
  The Cassandra plugin for Vault's Database backend generates database credentials to access Cassandra servers.
---

# Cassandra Database Plugin HTTP API

The Cassandra Database Plugin is one of the supported plugins for the Database
backend. This plugin generates database credentials dynamically based on
configured roles for the Cassandra database.

## Configure Connection

In addition to the parameters defined by the [Database
Backend](/api/secret/databases/index.html#configure-connection), this plugin
has a number of parameters to further configure a connection.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `POST`   | `/database/config/:name`     | `204 (empty body)` |

### Parameters
- `hosts` `(string: <required>)` – Specifies a set of comma-delineated Cassandra
  hosts to connect to.

- `port` `(int: 9042)` – Specifies the default port to use if none is provided
  as part of the host URI. Defaults to Cassandra's default transport port, 9042.

- `username` `(string: <required>)` – Specifies the username to use for
  superuser access.

- `password` `(string: <required>)` – Specifies the password corresponding to
  the given username.

- `tls` `(bool: true)` – Specifies whether to use TLS when connecting to
  Cassandra.

- `insecure_tls` `(bool: false)` – Specifies whether to skip verification of the
  server certificate when using TLS.

- `pem_bundle` `(string: "")` – Specifies concatenated PEM blocks containing a
  certificate and private key; a certificate, private key, and issuing CA
  certificate; or just a CA certificate.

- `pem_json` `(string: "")` – Specifies JSON containing a certificate and
  private key; a certificate, private key, and issuing CA certificate; or just a
  CA certificate. For convenience format is the same as the output of the
  `issue` command from the `pki` backend; see
  [the pki documentation](/docs/secrets/pki/index.html).

- `protocol_version` `(int: 2)` – Specifies the CQL protocol version to use.

- `connect_timeout` `(string: "5s")` – Specifies the connection timeout to use.

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
  "plugin_name": "cassandra-database-plugin",
  "allowed_roles": "readonly",
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

## Statements

Statements are configured during role creation and are used by the plugin to
determine what is sent to the datatabse on user creation, renewing, and
revocation. For more information on configuring roles see the [Role
API](/api/secret/databases/index.html#create-role) in the Database Backend docs.

### Parameters

The following are the statements used by this plugin. If not mentioned in this
list the plugin does not support that statement type.

- `creation_statements` `(string: "")` – Specifies the database
  statements executed to create and configure a user. Must be a
  semicolon-separated string, a base64-encoded semicolon-separated string, a
  serialized JSON string array, or a base64-encoded serialized JSON string
  array. The '{{name}}' and '{{password}}' values will be substituted. If not
  provided, defaults to a generic create user statements that creates a
  non-superuser.

- `revocation_statements` `(string: "")` – Specifies the database statements to
  be executed to revoke a user. Must be a semicolon-separated string, a
  base64-encoded semicolon-separated string, a serialized JSON string array, or
  a base64-encoded serialized JSON string array. The '{{name}}' value will be
  substituted. If not provided defaults to a generic drop user statement.

- `rollback_statements` `(string: "")` – Specifies the database statements to be
  executed to rollback a create operation in the event of an error. Must be a
  semicolon-separated string, a base64-encoded semicolon-separated string, a
  serialized JSON string array, or a base64-encoded serialized JSON string
  array. The '{{name}}' value will be substituted. If not provided, defaults to
  a generic drop user statement 
