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
on every path, use `vault help` after mounting the backend.

## Quick Start

The first step to using the PostgreSQL backend is to mount it.
Unlike the `generic` backend, the `postgresql` backend is not mounted by default.

```text
$ vault mount postgresql
Successfully mounted 'postgresql' at 'postgresql'!
```

Vault must be configured to connect to the PostgreSQL:

```text
$ vault write postgresql/config/connection \
    value="host=localhost user=username password=password"
```

This value must be a PG-style connection string, and the specified user must
have permission to manage the database.

Note, if you get an error like:

```text
pq: SSL is not enabled on the server
```

this means your PostgreSQL server has not enabled SSL. It is highly recommended
that you configure your PostgreSQL server to communicate via SSL.

The SSL check can be disabled by specifying the `sslmode=disable` attribute in
the PostgreSQL connection string:

```text
$ vault write postgresql/config/connection \
    value="host=localhost user=username password=password sslmode=disable"
```

Vault's PostgreSQL integration is role-based, so you must create a role for
which to request credentials:

```text
$ vault write postgresql/roles/production \
    name=production
```

Vault is now configured to create and manage credentials for Postgres!

```text
$ vault read postgresql/creds/production
Key             Value
lease_id        postgresql/creds/production/8ade2cde-5081-e3b7-af1a-3b9fb070df66
lease_duration  3600
password        56b43bc3-b285-4803-abdf-662d6a105bd0
username        vault-root-1430141210-1847
```

If you get stuck at any time, simply run `vault help postgresql` or with a
subpath for interactive help output.

## API

### /postgresql/config/connection
#### POST

<dl class="api">
  <dt>Description</dt>
  <dd>
    Configures the connection string used to communicate with PostgreSQL.
    This is a root protected endpoint.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/postgresql/config/connection`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">value</span>
        <span class="param-flags">required</span>
        The PostgreSQL connection URL or PG style string. e.g. "user=foo host=bar"
      </li>
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
        The SQL statements executed to create and configure the role.
        Must be semi-colon seperated. The '{{name}}', '{{password}}' and
        '{{expiration}}' values will be substituted.
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
            "username": "vault-root-1430158508-126",
            "password": "132ae3ef-5a64-7499-351e-bfe59f3a2a21"
        }
    }
    ```

  </dd>
</dl>

