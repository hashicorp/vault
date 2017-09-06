---
layout: "docs"
page_title: "Cassandra Secret Backend"
sidebar_current: "docs-secrets-cassandra"
description: |-
  The Cassandra secret backend for Vault generates database credentials to access Cassandra.
---

# Cassandra Secret Backend

Name: `cassandra`

~> **Deprecation Note:** This backend is deprecated in favor of the
combined databases backend added in v0.7.1. See the documentation for
the new implementation of this backend at
[Cassandra Database Plugin](/docs/secrets/databases/cassandra.html).

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
Unlike the `kv` backend, the `cassandra` backend is not mounted by default.

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

The Cassandra secret backend has a full HTTP API. Please see the
[Cassandra secret backend API](/api/secret/cassandra/index.html) for more
details.
