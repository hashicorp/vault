---
layout: "docs"
page_title: "Secret Backend: PostgreSQL"
sidebar_current: "docs-secrets-postgresql"
description: |-
  The PostgreSQL secret backend for Vault generates database credentials to access PostgreSQL.
---

# PostgreSQL Secret Backend

Name: `postgresql`

The PostgreSQL secret backend for Vault generates database credentials
dynamically based on configured roles. This means that services that need
to access a database no longer need to hardcode credentials: they can request
them from Vault, and use Vault's leasing mechanism to more easily roll keys.

Additionally, it introduces a new ability: with every service accessing
the database with unique credentials, it makes auditing much easier when
questionable data access is discovered: you can track it down to the specific
instance of a service based on the SQL username.

Vault makes use both of its own internal revocation system as well as the
`VALID UNTIL` setting when creating PostgreSQL users to ensure that users
become invalid within a reasonable time of the lease expiring.

This page will show a quick start for this backend. For detailed documentation
on every path, use `vault path-help` after mounting the backend.

## Quick Start

The first step to using the PostgreSQL backend is to mount it.
Unlike the `generic` backend, the `postgresql` backend is not mounted by default.

```text
$ vault mount postgresql
Successfully mounted 'postgresql' at 'postgresql'!
```

Next, Vault must be configured to connect to the PostgreSQL. This is done by
writing either a PostgreSQL URL or PG connection string:

```text
$ vault write postgresql/config/connection \
    connection_url="postgresql://root:vaulttest@vaulttest.ciuvljjni7uo.us-west-1.rds.amazonaws.com:5432/postgres"
```

In this case, we've configured Vault with the user "root" and password "vaulttest",
connecting to a PostgreSQL instance in AWS RDS. The "postgres" database name is being used.
It is important that the Vault user have the `GRANT OPTION` privilege to manage users.

Optionally, we can configure the lease settings for credentials generated
by Vault. This is done by writing to the `config/lease` key:

```
$ vault write postgresql/config/lease lease=1h lease_max=24h
Success! Data written to: postgresql/config/lease
```

This restricts each credential to being valid or leased for 1 hour
at a time, with a maximum use period of 24 hours. This forces an
application to renew their credentials at least hourly, and to recycle
them once per day.

The next step is to configure a role. A role is a logical name that maps
to a policy used to generated those credentials. For example, lets create
a "readonly" role:

```text
$ vault write postgresql/roles/readonly \
    sql="CREATE ROLE \"{{name}}\" WITH LOGIN PASSWORD '{{password}}' VALID UNTIL '{{expiration}}';
    GRANT SELECT ON ALL TABLES IN SCHEMA public TO \"{{name}}\";"
Success! Data written to: postgresql/roles/readonly
```

By writing to the `roles/readonly` path we are defining the `readonly` role.
This role will be created by evaluating the given `sql` statements. By
default, the `{{name}}`, `{{password}}` and `{{expiration}}` fields will be populated by
Vault with dynamically generated values. This SQL statement is creating
the named user, and then granting it `SELECT` or read-only privileges
to tables in the database. More complex `GRANT` queries can be used to
customize the privileges of the role. See the [PostgreSQL manual](http://www.postgresql.org/docs/9.4/static/sql-grant.html)
for more information.

To generate a new set of credentials, we simply read from that role:
Vault is now configured to create and manage credentials for Postgres!

```text
$ vault read postgresql/creds/readonly
Key           	Value
lease_id      	postgresql/creds/readonly/c888a097-b0e2-26a8-b306-fc7c84b98f07
lease_duration	3600
password      	34205e88-0de1-68b7-6267-72d8e32c5d3d
username      	root-1430162075-7887
```

By reading from the `creds/readonly` path, Vault has generated a new
set of credentials using the `readonly` role configuration. Here we
see the dynamically generated username and password, along with a one
hour lease.

Using ACLs, it is possible to restrict using the postgresql backend such
that trusted operators can manage the role definitions, and both
users and applications are restricted in the credentials they are
allowed to read.

If you get stuck at any time, simply run `vault path-help postgresql` or with a
subpath for interactive help output.

## API

### /postgresql/config/connection
#### POST

<dl class="api">
  <dt>Description</dt>
  <dd>
    Configures the connection string used to communicate with PostgreSQL.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/postgresql/config/connection`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">connection_url</span>
        <span class="param-flags">required</span>
        The PostgreSQL connection URL or PG style string. e.g. "user=foo host=bar"
      </li>
    </ul>
  </dd>
  <dd>
    <ul>
      <li>
        <span class="param">value</span>
        <span class="param-flags">optional</span>
        The PostgreSQL connection URL or PG style string. e.g. "user=foo host=bar". Use `connection_url` instead.
      </li>
    </ul>
  </dd>
  <dd>
    <ul>
      <li>
        <span class="param">max_open_connections</span>
        <span class="param-flags">optional</span>
        Maximum number of open connections to the database. A zero uses the
        default value of 2 and a negative value means unlimited.
      </li>
    </ul>
  </dd>
  <dd>
    <ul>
        <span class="param">max_idle_connections</span>
        <span class="param-flags">optional</span>
        Maximum number of idle connections to the database. A zero uses the
        value of `max_open_connections` and a negative value disables idle
        connections. If larger than `max_open_connections` it will be reduced
        to be equal.
    </ul>
  </dd>
  <dd>
    <ul>
        <span class="param">verify_connection</span>
        <span class="param-flags">optional</span>
	If set, connection_url is verified by actually connecting to the database.
	Defaults to true.
    </ul>
  </dd>

  <dt>Returns</dt>
  <dd>
    A `204` response code.
  </dd>
</dl>

### /postgresql/config/lease
#### POST

<dl class="api">
  <dt>Description</dt>
  <dd>
    Configures the lease settings for generated credentials.
    If not configured, leases default to 1 hour. This is a root
    protected endpoint.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/postgresql/config/lease`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">lease</span>
        <span class="param-flags">required</span>
        The lease value provided as a string duration
        with time suffix. Hour is the largest suffix.
      </li>
      <li>
        <span class="param">lease_max</span>
        <span class="param-flags">required</span>
        The maximum lease value provided as a string duration
        with time suffix. Hour is the largest suffix.
      </li>
    </ul>
  </dd>

  <dt>Returns</dt>
  <dd>
    A `204` response code.
  </dd>
</dl>

### /postgresql/roles/
#### POST

<dl class="api">
  <dt>Description</dt>
  <dd>
    Creates or updates the role definition.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/postgresql/roles/<name>`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">sql</span>
        <span class="param-flags">required</span>
        The SQL statements executed to create and configure the role. Must be
        a semicolon-separated string, a base64-encoded semicolon-separated
        string, a serialized JSON string array, or a base64-encoded serialized
        JSON string array. The '{{name}}', '{{password}}' and '{{expiration}}'
        values will be substituted.
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
  <dd>`/postgresql/roles/<name>`</dd>

  <dt>Parameters</dt>
  <dd>
     None
  </dd>

  <dt>Returns</dt>
  <dd>

    ```javascript
    {
      "data": {
        "sql": "CREATE USER..."
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
  <dd>`/postgresql/roles` (LIST) or `/postgresql/roles/?list=true` (GET)</dd>

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
  <dd>`/postgresql/roles/<name>`</dd>

  <dt>Parameters</dt>
  <dd>
     None
  </dd>

  <dt>Returns</dt>
  <dd>
    A `204` response code.
  </dd>
</dl>

### /postgresql/creds/
#### GET

<dl class="api">
  <dt>Description</dt>
  <dd>
    Generates a new set of dynamic credentials based on the named role.
  </dd>

  <dt>Method</dt>
  <dd>GET</dd>

  <dt>URL</dt>
  <dd>`/postgresql/creds/<name>`</dd>

  <dt>Parameters</dt>
  <dd>
     None
  </dd>

  <dt>Returns</dt>
  <dd>

    ```javascript
    {
      "data": {
        "username": "root-1430158508-126",
        "password": "132ae3ef-5a64-7499-351e-bfe59f3a2a21"
      }
    }
    ```

  </dd>
</dl>

