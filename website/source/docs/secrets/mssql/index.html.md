---
layout: "docs"
page_title: "Secret Backend: mssql"
sidebar_current: "docs-secrets-mssql"
description: |-
  The MSSQL secret backend for Vault generates database credentials to access Microsoft Sql Server.
---

# MSSQL Secret Backend

Name: `mssql`

The MSSQL secret backend for Vault generates database credentials
dynamically based on configured roles. This means that services that need
to access a database no longer need to hardcode credentials: they can request
them from Vault, and use Vault's leasing mechanism to more easily roll keys.

Additionally, it introduces a new ability: with every service accessing
the database with unique credentials, it makes auditing much easier when
questionable data access is discovered: you can track it down to the specific
instance of a service based on the SQL username.

Vault makes use of its own internal revocation system to ensure that users
become invalid within a reasonable time of the lease expiring.

This page will show a quick start for this backend. For detailed documentation
on every path, use `vault path-help` after mounting the backend.

## Quick Start

The first step to using the mssql backend is to mount it.
Unlike the `generic` backend, the `mssql` backend is not mounted by default.

```
$ vault mount mssql
Successfully mounted 'mssql' at 'mssql'!
```

Next, we must configure Vault to know how to connect to the MSSQL
instance. This is done by providing a DSN (Data Source Name):

```
$ vault write mssql/config/connection \
    connection_string="server=localhost;port=1433;user id=sa;password=Password!;database=AdventureWorks;app name=vault;"
Success! Data written to: mssql/config/connection
```

In this case, we've configured Vault with the user "sa" and password "Password!",
connecting to an instance at "localhost" on port 1433. It is not necessary
that Vault has the sa login, but the user must have privileges to create
logins and manage processes. The fixed server roles `securityadmin` and
`processadmin` are examples of built-in roles that grant these permissions. The
user also must have privileges to create database users and grant permissions in
the databases that Vault manages.  The fixed database roles `db_accessadmin` and
`db_securityadmin` are examples or built-in roles that grant these permissions.

Optionally, we can configure the lease settings for credentials generated
by Vault. This is done by writing to the `config/lease` key:

```
$ vault write mssql/config/lease \
    ttl=1h \
    max_ttl=24h
Success! Data written to: mssql/config/lease
```

This restricts each credential to being valid or leased for 1 hour
at a time, with a maximum use period of 24 hours. This forces an
application to renew their credentials at least hourly, and to recycle
them once per day.

The next step is to configure a role. A role is a logical name that maps
to a policy used to generate those credentials. For example, lets create
a "readonly" role:

```
$ vault write mssql/roles/readonly \
    sql="CREATE LOGIN [{{name}}] WITH PASSWORD = '{{password}}'; USE AdventureWorks; CREATE USER [{{name}}] FOR LOGIN [{{name}}]; GRANT SELECT ON SCHEMA::dbo TO [{{name}}]"
Success! Data written to: mssql/roles/readonly
```

By writing to the `roles/readonly` path we are defining the `readonly` role.
This role will be created by evaluating the given `sql` statements. By
default, the `{{name}}` and `{{password}}` fields will be populated by
Vault with dynamically generated values. This SQL statement is creating
the named login on the server, user on the AdventureWorks database, and
then granting it `SELECT` on the `dbo` schema. More complex `GRANT` queries
can be used to customize the privileges of the role.

To generate a new set of credentials, we simply read from that role:

```
$ vault read mssql/creds/readonly
Key           	Value
lease_id      	mssql/creds/readonly/cdf23ac8-6dbd-4bf9-9919-6acaaa86ba6c
lease_duration	3600
password      	ee202d0d-e4fd-4410-8d14-2a78c5c8cb76
username      	root-a147d529-e7d6-4a16-8930-4c3e72170b19
```

By reading from the `creds/readonly` path, Vault has generated a new
set of credentials using the `readonly` role configuration. Here we
see the dynamically generated username and password, along with a one
hour lease.

Using ACLs, it is possible to restrict using the mssql backend such
that trusted operators can manage the role definitions, and both
users and applications are restricted in the credentials they are
allowed to read.

## API

### /mssql/config/connection
#### POST

<dl class="api">
  <dt>Description</dt>
  <dd>
    Configures the connection DSN used to communicate with Sql Server.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/mssql/config/connection`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">connection_string</span>
        <span class="param-flags">required</span>
        The MSSQL DSN
      </li>
    </ul>
  </dd>
  <dd>
    <ul>
      <li>
        <span class="param">max_open_connections</span>
        <span class="param-flags">optional</span>
        Maximum number of open connections to the database.
	Defaults to 2.
      </li>
    </ul>
  </dd>
  <dd>
    <ul>
      <li>
        <span class="param">verify_connection</span>
        <span class="param-flags">optional</span>
	If set, connection_string is verified by actually connecting to the database.
	Defaults to true.
      </li>
    </ul>
  </dd>

  <dt>Returns</dt>
  <dd>
    A `204` response code.
  </dd>
</dl>

### /mssql/config/lease
#### POST

<dl class="api">
  <dt>Description</dt>
  <dd>
    Configures the lease settings for generated credentials.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/mssql/config/lease`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">ttl</span>
        <span class="param-flags">required</span>
        The ttl value provided as a string duration
        with time suffix. Hour is the largest suffix.
      </li>
      <li>
        <span class="param">max_ttl</span>
        <span class="param-flags">required</span>
        The maximum ttl value provided as a string duration
        with time suffix. Hour is the largest suffix.
      </li>
    </ul>
  </dd>

  <dt>Returns</dt>
  <dd>
    A `204` response code.
  </dd>
</dl>

#### GET

<dl class="api">
  <dt>Description</dt>
  <dd>
    Queries the lease configuration.
  </dd>

  <dt>Method</dt>
  <dd>GET</dd>

  <dt>URL</dt>
  <dd>`/mssql/config/lease`</dd>

  <dt>Parameters</dt>
  <dd>
     None
  </dd>

  <dt>Returns</dt>
  <dd>

    ```javascript
    {
      "lease_id": "",
      "renewable": false,
      "lease_duration": 0,
      "data": {
        "max_ttl": "5h",
        "ttl": "1h",
        "ttl_max": "5h"
      },
      "wrap_info": null,
      "warnings": ["The field ttl_max is deprecated and will be removed in a future release. Use max_ttl instead."],
      "auth": null
    }
    ```

  </dd>
</dl>

### /mssql/roles/
#### POST

<dl class="api">
  <dt>Description</dt>
  <dd>
    Creates or updates the role definition.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/mssql/roles/<name>`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">sql</span>
        <span class="param-flags">required</span>
        The SQL statements executed to create and configure the role.  The
        '{{name}}' and '{{password}}' values will be substituted. Must be a
        semicolon-separated string, a base64-encoded semicolon-separated
        string, a serialized JSON string array, or a base64-encoded serialized
        JSON string array.
      </li>
    </ul>
  </dd>

  <dt>Returns</dt>
  <dd>
    A `204` response code.
  </dd>
</dl>

#### GET

<dl class="api">
  <dt>Description</dt>
  <dd>
    Queries the role definition.
  </dd>

  <dt>Method</dt>
  <dd>GET</dd>

  <dt>URL</dt>
  <dd>`/mssql/roles/<name>`</dd>

  <dt>Parameters</dt>
  <dd>
     None
  </dd>

  <dt>Returns</dt>
  <dd>

    ```javascript
    {
      "data": {
        "sql": "CREATE LOGIN..."
      }
    }
    ```

  </dd>
</dl>

#### LIST

<dl class="api">
  <dt>Description</dt>
  <dd>
    Returns a list of available roles. Only the role names are returned, not
    any values.
  </dd>

  <dt>Method</dt>
  <dd>LIST/GET</dd>

  <dt>URL</dt>
  <dd>`/mssql/roles` (LIST) or `/mssql/roles/?list=true` (GET)</dd>

  <dt>Parameters</dt>
  <dd>
     None
  </dd>

  <dt>Returns</dt>
  <dd>

  ```javascript
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

  </dd>
</dl>

#### DELETE

<dl class="api">
  <dt>Description</dt>
  <dd>
    Deletes the role definition.
  </dd>

  <dt>Method</dt>
  <dd>DELETE</dd>

  <dt>URL</dt>
  <dd>`/mssql/roles/<name>`</dd>

  <dt>Parameters</dt>
  <dd>
     None
  </dd>

  <dt>Returns</dt>
  <dd>
    A `204` response code.
  </dd>
</dl>

### /mssql/creds/
#### GET

<dl class="api">
  <dt>Description</dt>
  <dd>
    Generates a new set of dynamic credentials based on the named role.
  </dd>

  <dt>Method</dt>
  <dd>GET</dd>

  <dt>URL</dt>
  <dd>`/mssql/creds/<name>`</dd>

  <dt>Parameters</dt>
  <dd>
     None
  </dd>

  <dt>Returns</dt>
  <dd>

    ```javascript
    {
      "data": {
        "username": "root-a147d529-e7d6-4a16-8930-4c3e72170b19",
        "password": "ee202d0d-e4fd-4410-8d14-2a78c5c8cb76"
      }
    }
    ```

  </dd>
</dl>
