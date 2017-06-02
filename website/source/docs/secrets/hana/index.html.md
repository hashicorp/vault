---
layout: "docs"
page_title: "Secret Backend: hana"
sidebar_current: "docs-secrets-hana"
description: |-
  The HANA secret backend for Vault generates database credentials to access HANA Database Server.
---

# HANA Secret Backend

Name: `hana`

The HANA secret backend for Vault generates database credentials
dynamically based on configured roles. This means that services that need
to access a database no longer need to hardcode credentials: they can request
them from Vault, and use Vault's leasing mechanism to more easily roll keys.

Additionally, it introduces a new ability: with every service accessing
the database with unique credentials, it makes auditing much easier when
questionable data access is discovered: you can track it down to the specific
instance of a service based on the SQL username.

Credentials generated from the HANA secret backend will be activated for the
configured lease TTL, unless renewed. After lease expiry, the generated
account will be deactivated HANA-side without action from Vault.

In addition, when the credentials are revoked due to lease expiry or explicit
revocation, a soft attempt will be made to drop the user from the database.

This page will show a quick start for this backend. For detailed documentation
on every path, use `vault path-help` after mounting the backend.

## Quick Start

The first step to using the Hana backend is to mount it.
Unlike the `generic` backend, the `hana` backend is not mounted by default.

```
$ vault mount hana
Successfully mounted 'hana' at 'hana'!
```

Next, we must configure Vault to know how to connect to the HANA
instance. This is done by providing a DSN (Data Source Name):

```
$ vault write hana/config/connection \
    connection_string="hdb://username:password@127.0.0.1:30015" verify_connection=true
Success! Data written to: hana/config/connection
```

In this case, we've configured Vault with the user "username" and password "password",
connecting to an instance at "127.0.0.1" on port 30015. The configured user only
requires the priviledges to create users, delete users, and grant permissions as required.

Optionally, we can configure the lease settings for credentials generated
by Vault. This is done by writing to the `config/lease` key:

```
$ vault write hana/config/lease \
    ttl=1h \
    ttl_max=24h
Success! Data written to: hana/config/lease
```

This restricts each credential to being valid or leased for 1 hour
at a time, with a maximum use period of 24 hours. This forces an
application to renew their credentials at least hourly, and to recycle
them once per day. When the credentials expire, HANA will deactivate
the user, and Vault will attempt a restrict drop on the user.


The next step is to configure a role. A role is a logical name that maps
to a policy used to generate those credentials. For example, lets create
a "monitoring" role:

```
$ vault write hana/roles/monitoring \
    sql="CREATE USER {{name}} PASSWORD {{password}} VALID UNTIL '{{valid_until}}'; CALL GRANT_ACTIVATED_ROLE ( 'sap.hana.admin.roles::Monitoring', '{{name}}' );"
Success! Data written to: hana/roles/monitoring
```

By writing to the `roles/monitoring` path we are defining the `monitoring` role.
This role will be created by evaluating the given `sql` statements. By
default, the `{{name}}` and `{{password}}` fields will be populated by
Vault with dynamically generated values, and the '{{valid_until}}' field will
be substituted with the lease TTL. This SQL statement is creating a new user
on the server, which will be deactivated and possibly dropped when the lease
expires, and has been granted the monitoring role. Other GRANT queries may be
added in conjunction to fit needs. The given SQL statements are not tested
due to HANA rejecting unknown usernames when preparing GRANT statements.


To generate a new set of credentials, we simply read from that role:

```
$ vault read hana/creds/monitoring
Key             Value
lease_id        hana/creds/monitoring/cdf23ac8-6dbd-4bf9-9919-6acaaa86ba6c
lease_duration  3600
password        A1acd07ff89_9832_316c_bb53_1a7ef59ac76c
username        ROOT_5589A574_6E66_9F6E_C8E7_F553B35ED592
valid_until     2016-10-27 15:12:47
```

By reading from the `creds/monitoring` path, Vault has generated a new
set of credentials using the `monitoring` role configuration. Here we
see the dynamically generated username and password, along with a one
hour lease. The valid_until value is the time at which the HANA instance
will deactivate the given credentials (in HANA-side time).

Using ACLs, it is possible to restrict using the HANA backend such
that trusted operators can manage the role definitions, and both
users and applications are restricted in the credentials they are
allowed to read.

## API

### /hana/config/connection
#### POST

<dl class="api">
  <dt>Description</dt>
  <dd>
    Configures the connection DSN used to communicate with HANA Server.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/hana/config/connection`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">connection_string</span>
        <span class="param-flags">required</span>
        The HANA DSN
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

### /hana/config/lease
#### POST

<dl class="api">
  <dt>Description</dt>
  <dd>
    Configures the lease settings for generated credentials.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/hana/config/lease`</dd>

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
        <span class="param">ttl_max</span>
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

### /hana/roles/
#### POST

<dl class="api">
  <dt>Description</dt>
  <dd>
    Creates or updates the role definition.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/hana/roles/<name>`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">sql</span>
        <span class="param-flags">required</span>
        The SQL statements executed to create and configure the role.  The
        '{{name}}', '{{password}}', and '{{valid_until}}' values will be substituted.
        Must be a semicolon-separated string, a base64-encoded semicolon-separated
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
  <dd>`/hana/roles/<name>`</dd>

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
  <dd>`/hana/roles/<name>`</dd>

  <dt>Parameters</dt>
  <dd>
     None
  </dd>

  <dt>Returns</dt>
  <dd>
    A `204` response code.
  </dd>
</dl>

### /hana/creds/
#### GET

<dl class="api">
  <dt>Description</dt>
  <dd>
    Generates a new set of dynamic credentials based on the named role.
    The password must be changed upon login
  </dd>

  <dt>Method</dt>
  <dd>GET</dd>

  <dt>URL</dt>
  <dd>`/hana/creds/<name>`</dd>

  <dt>Parameters</dt>
  <dd>
     None
  </dd>

  <dt>Returns</dt>
  <dd>

    ```javascript
    {
      "data": {
        "username": "Vault_5589A574_6E66_9F6E_C8E7_F553B35ED592",
        "password": "A1acd07ff89_9832_316c_bb53_1a7ef59ac76c"
      }
    }
    ```

  </dd>
</dl>