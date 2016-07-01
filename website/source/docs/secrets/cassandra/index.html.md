---
layout: "docs"
page_title: "Secret Backend: Cassandra"
sidebar_current: "docs-secrets-cassandra"
description: |-
  The Cassandra secret backend for Vault generates database credentials to access Cassandra.
---

# Cassandra Secret Backend

Name: `cassandra`

The Cassandra secret backend for Vault generates database credentials
dynamically based on configured roles. This means that services that need
to access a database no longer need to hardcode credentials: they can request
them from Vault, and use Vault's leasing mechanism to more easily roll keys.

Additionally, it introduces a new ability: with every service accessing
the database with unique credentials, it makes auditing much easier when
questionable data access is discovered: you can track it down to the specific
instance of a service based on the Cassandra username.

This page will show a quick start for this backend. For detailed documentation
on every path, use `vault path-help` after mounting the backend.

## Quick Start

The first step to using the Cassandra backend is to mount it.
Unlike the `generic` backend, the `cassandra` backend is not mounted by default.

```text
$ vault mount cassandra
Successfully mounted 'cassandra' at 'cassandra'!
```

Next, Vault must be configured to connect to Cassandra. This is done by
writing one or more hosts, a username, and a password:

```text
$ vault write cassandra/config/connection \
    hosts=localhost \
    username=cassandra \
    password=cassandra
```

In this case, we've configured Vault with the user "cassandra" and password "cassandra",
It is important that the Vault user is a superuser, in order to manage other user accounts.

The next step is to configure a role. A role is a logical name that maps
to a policy used to generated those credentials. For example, lets create
a "readonly" role:

```text
$ vault write cassandra/roles/readonly \
    creation_cql="CREATE USER '{{username}}' WITH PASSWORD '{{password}}' NOSUPERUSER; \
    GRANT SELECT ON ALL KEYSPACES TO {{username}};"
Success! Data written to: cassandra/roles/readonly
```

By writing to the `roles/readonly` path we are defining the `readonly` role.
This role will be created by evaluating the given `creation_cql` statements. By
default, the `{{username}}` and `{{password}}` fields will be populated by
Vault with dynamically generated values. This CQL statement is creating
the named user, and then granting it `SELECT` or read-only privileges
to keyspaces. More complex `GRANT` queries can be used to
customize the privileges of the role. See the [CQL Reference Manual](https://docs.datastax.com/en/cql/3.1/cql/cql_reference/grant_r.html)
for more information.

To generate a new set of credentials, we simply read from that role:
Vault is now configured to create and manage credentials for Cassandra!

```text
$ vault read cassandra/creds/readonly
Key           	Value
lease_id       	cassandra/creds/test/7a23e890-3a26-531d-529b-92d18d1fa63f
lease_duration 	3600
lease_renewable	true
password       	dfa80eea-ccbe-b228-ebf7-e2f62b245e71
username       	vault-root-1434647667-9313
```

By reading from the `creds/readonly` path, Vault has generated a new
set of credentials using the `readonly` role configuration. Here we
see the dynamically generated username and password, along with a one
hour lease.

Using ACLs, it is possible to restrict using the `cassandra` backend such
that trusted operators can manage the role definitions, and both
users and applications are restricted in the credentials they are
allowed to read.

If you get stuck at any time, simply run `vault path-help cassandra` or with a
subpath for interactive help output.

## API

### /cassandra/config/connection
#### POST

<dl class="api">
  <dt>Description</dt>
  <dd>
    Configures the connection information used to communicate with Cassandra.
    TLS works as follows:<br /><br />
    <ul>
      <li>
        • If `tls` is set to true, the connection will use TLS; this happens
        automatically if `pem_bundle`, `pem_json`, or `insecure_tls` is set
      </li>
      <li>
        • If `insecure_tls` is set to true, the connection will not perform
        verification of the server certificate; this also sets `tls` to true
      </li>
      <li>
        • If only `issuing_ca` is set in `pem_json`, or the only certificate in
        `pem_bundle` is a CA certificate, the given CA certificate will be used
        for server certificate verification; otherwise the system CA
        certificates will be used
      </li>
      <li>
        • If `certificate` and `private_key` are set in `pem_bundle` or
        `pem_json`, client auth will be turned on for the connection
      </li>
    </ul>
    `pem_bundle` should be a PEM-concatenated bundle of a private key + client
    certificate, an issuing CA certificate, or both. `pem_json` should contain
    the same information; for convenience, the JSON format is the same as that
    output by the issue command from the PKI backend.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/cassandra/config/connection`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">hosts</span>
        <span class="param-flags">required</span>
        A set of comma-deliniated Cassandra hosts to connect to.
      </li>
      <li>
        <span class="param">username</span>
        <span class="param-flags">required</span>
        The username to use for superuser access.
      </li>
      <li>
        <span class="param">password</span>
        <span class="param-flags">required</span>
        The password corresponding to the given username.
      </li>
      <li>
        <span class="param">tls</span>
        <span class="param-flags">optional</span>
        Whether to use TLS when connecting to Cassandra.
      </li>
      <li>
        <span class="param">insecure_tls</span>
        <span class="param-flags">optional</span>
        Whether to skip verification of the server certificate when using TLS.
      </li>
      <li>
        <span class="param">pem_bundle</span>
        <span class="param-flags">optional</span>
        Concatenated PEM blocks containing a certificate and private key;
        a certificate, private key, and issuing CA certificate; or just a CA
        certificate.
      </li>
      <li>
        <span class="param">pem_json</span>
        <span class="param-flags">optional</span>
        JSON containing a certificate and private key;
        a certificate, private key, and issuing CA certificate; or just a CA
        certificate. For convenience format is the same as the output of the
        `issue` command from the `pki` backend; see [the pki documentation](https://www.vaultproject.io/docs/secrets/pki/index.html).
      </li>
      <li>
        <span class="param">protocol_version</span>
        <span class="param-flags">optional</span>
        The CQL protocol version to use. Defaults to 2.
      </li>
      <li>
        <span class="param">connect_timeout</span>
        <span class="param-flags">optional</span>
        The connection timeout to use. Defaults to 5 seconds.
      </li>
    </ul>
  </dd>

  <dt>Returns</dt>
  <dd>
    A `204` response code.
  </dd>
</dl>

### /cassandra/roles/
#### POST

<dl class="api">
  <dt>Description</dt>
  <dd>
    Creates or updates the role definition.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/cassandra/roles/<name>`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">creation_cql</span>
        <span class="param-flags">optional</span>
        The CQL statements executed to create and configure the new user.
        Must be semi-colon separated. The '{{username}}' and '{{password}}'
        values will be substituted; it is required that these parameters are
        in single quotes. The default creates a non-superuser user with
        no authorization grants.
      </li>
      <li>
        <span class="param">rollback_cql</span>
        <span class="param-flags">optional</span>
        The CQL statements executed to attempt a rollback if an error is
        encountered during user creation. The default is to delete the user.
        Must be semi-colon separated. The '{{username}}' and '{{password}}'
        values will be substituted; it is required that these parameters are
        in single quotes.
      </li>
      <li>
        <span class="param">lease</span>
        <span class="param-flags">optional</span>
        The lease value provided as a string duration
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
    Queries the role definition.
  </dd>

  <dt>Method</dt>
  <dd>GET</dd>

  <dt>URL</dt>
  <dd>`/cassandra/roles/<name>`</dd>

  <dt>Parameters</dt>
  <dd>
     None
  </dd>

  <dt>Returns</dt>
  <dd>

    ```javascript
    {
      "data": {
        "creation_cql": "CREATE USER...",
        "rollback_cql": "DROP USER...",
        "lease": "12h",
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
  <dd>`/cassandra/roles/<name>`</dd>

  <dt>Parameters</dt>
  <dd>
     None
  </dd>

  <dt>Returns</dt>
  <dd>
    A `204` response code.
  </dd>
</dl>

### /cassandra/creds/
#### GET

<dl class="api">
  <dt>Description</dt>
  <dd>
    Generates a new set of dynamic credentials based on the named role.
  </dd>

  <dt>Method</dt>
  <dd>GET</dd>

  <dt>URL</dt>
  <dd>`/cassandra/creds/<name>`</dd>

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
