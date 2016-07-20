---
layout: "docs"
page_title: "Secret Backend: MySQL"
sidebar_current: "docs-secrets-mysql"
description: |-
  The MySQL secret backend for Vault generates database credentials to access MySQL.
---

# MySQL Secret Backend

Name: `mysql`

The MySQL secret backend for Vault generates database credentials
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

The first step to using the mysql backend is to mount it.
Unlike the `generic` backend, the `mysql` backend is not mounted by default.

```
$ vault mount mysql
Successfully mounted 'mysql' at 'mysql'!
```

Next, we must configure Vault to know how to connect to the MySQL
instance. This is done by providing a DSN (Data Source Name):

```
$ vault write mysql/config/connection \
    connection_url="root:root@tcp(192.168.33.10:3306)/"
Success! Data written to: mysql/config/connection
```

In this case, we've configured Vault with the user "root" and password "root,
connecting to an instance at "192.168.33.10" on port 3306. It is not necessary
that Vault has the root user, but the user must have privileges to create
other users, namely the `GRANT OPTION` privilege.

Optionally, we can configure the lease settings for credentials generated
by Vault. This is done by writing to the `config/lease` key:

```
$ vault write mysql/config/lease \
    lease=1h \
    lease_max=24h
Success! Data written to: mysql/config/lease
```

This restricts each credential to being valid or leased for 1 hour
at a time, with a maximum use period of 24 hours. This forces an
application to renew their credentials at least hourly, and to recycle
them once per day.

The next step is to configure a role. A role is a logical name that maps
to a policy used to generate those credentials. For example, lets create
a "readonly" role:

```
$ vault write mysql/roles/readonly \
    sql="CREATE USER '{{name}}'@'%' IDENTIFIED BY '{{password}}';GRANT SELECT ON *.* TO '{{name}}'@'%';"
Success! Data written to: mysql/roles/readonly
```

By writing to the `roles/readonly` path we are defining the `readonly` role.
This role will be created by evaluating the given `sql` statements. By
default, the `{{name}}` and `{{password}}` fields will be populated by
Vault with dynamically generated values. This SQL statement is creating
the named user, and then granting it `SELECT` or read-only privileges
to tables in the database. More complex `GRANT` queries can be used to
customize the privileges of the role. See the [MySQL manual](https://dev.mysql.com/doc/refman/5.7/en/grant.html)
for more information.

To generate a new set of credentials, we simply read from that role:

```
$ vault read mysql/creds/readonly
Key           	Value
lease_id      	mysql/creds/readonly/bd404e98-0f35-b378-269a-b7770ef01897
lease_duration	3600
password      	132ae3ef-5a64-7499-351e-bfe59f3a2a21
username      	readonly-aefa635a-18
```

By reading from the `creds/readonly` path, Vault has generated a new
set of credentials using the `readonly` role configuration. Here we
see the dynamically generated username and password, along with a one
hour lease.

Using ACLs, it is possible to restrict using the mysql backend such
that trusted operators can manage the role definitions, and both
users and applications are restricted in the credentials they are
allowed to read.

Optionally, you may configure both the number of characters from the role name
that are truncated to form the display name portion of the mysql username
interpolated into the `{{name}}` field: the default is 10. 

You may also configure the total number of characters allowed in the entire
generated username (the sum of the display name and uuid poritions); the
default is 16. Note that versions of MySQL prior to 5.8 have a 16 character
total limit on user names, so it is probably not safe to increase this above
the default on versions prior to that.

## API

### /mysql/config/connection
#### POST

<dl class="api">
  <dt>Description</dt>
  <dd>
    Configures the connection DSN used to communicate with MySQL.
    This is a root protected endpoint.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/mysql/config/connection`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">connection_url</span>
        <span class="param-flags">required</span>
        The MySQL DSN
      </li>
    </ul>
  </dd>
  <dd>
    <ul>
      <li>
        <span class="param">value</span>
        <span class="param-flags">optional</span>
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
        <span class="param">max_idle_connections</span>
        <span class="param-flags">optional</span>
        Maximum number of idle connections to the database. A zero uses the value of `max_open_connections` and a negative value disables idle connections. If larger than `max_open_connections` it will be reduced to be equal.
      </li>
    </ul>
  </dd>
  <dd>
    <ul>
      <li>
        <span class="param">verify-connection</span>
        <span class="param-flags">optional</span>
	If set, connection_url is verified by actually connecting to the database.
	Defaults to true.
      </li>
    </ul>
  </dd>

  <dt>Returns</dt>
  <dd>
    A `204` response code.
  </dd>
</dl>

### /mysql/config/lease
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
  <dd>`/mysql/config/lease`</dd>

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

### /mysql/roles/
#### POST

<dl class="api">
  <dt>Description</dt>
  <dd>
    Creates or updates the role definition.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/mysql/roles/<name>`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">sql</span>
        <span class="param-flags">required</span>
        The SQL statements executed to create and configure the role.
        Must be semi-colon separated. The '{{name}}' and '{{password}}'
        values will be substituted.
      </li>
      <li>
        <span class="param">rolename_length</span>
        <span class="param-flags">optional</span>
        Determines how many characters from the role name will be used
        to form the mysql username interpolated into the '{{name}}' field
        of the sql parameter.  The default is 4.
      </li>
      <li>
        <span class="param">displayname_length</span>
        <span class="param-flags">optional</span>
        Determines how many characters from the token display name will be used
        to form the mysql username interpolated into the '{{name}}' field
        of the sql parameter.  The default is 4.
      </li>
      <li>
        <span class="param">username_length</span>
        <span class="param-flags">optional</span>
        Determines the maximum total length in characters of the
        mysql username interpolated into the '{{name}}' field
        of the sql parameter.  The default is 16.
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
  <dd>`/mysql/roles/<name>`</dd>

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
  <dd>GET</dd>

  <dt>URL</dt>
  <dd>`/roles/?list=true`</dd>

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
    "lease_duration": 2592000,
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
  <dd>`/mysql/roles/<name>`</dd>

  <dt>Parameters</dt>
  <dd>
     None
  </dd>

  <dt>Returns</dt>
  <dd>
    A `204` response code.
  </dd>
</dl>

### /mysql/creds/
#### GET

<dl class="api">
  <dt>Description</dt>
  <dd>
    Generates a new set of dynamic credentials based on the named role.
  </dd>

  <dt>Method</dt>
  <dd>GET</dd>

  <dt>URL</dt>
  <dd>`/mysql/creds/<name>`</dd>

  <dt>Parameters</dt>
  <dd>
     None
  </dd>

  <dt>Returns</dt>
  <dd>

    ```javascript
    {
      "data": {
        "username": "user-role-aefa63",
        "password": "132ae3ef-5a64-7499-351e-bfe59f3a2a21"
      }
    }
    ```

  </dd>
</dl>

