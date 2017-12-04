---
layout: "docs"
page_title: "RabbitMQ Secret Backend"
sidebar_current: "docs-secrets-rabbitmq"
description: |-
  The RabbitMQ secret backend for Vault generates user credentials to access RabbitMQ.
---

# RabbitMQ Secret Backend

Name: `rabbitmq`

The RabbitMQ secret backend for Vault generates user credentials dynamically
based on configured permissions and virtual hosts. This means that services
that need to access a virtual host no longer need to hardcode credentials:
they can request them from Vault, and use Vault's leasing mechanism to
more easily roll users.

Additionally, it introduces a new ability: with every service accessing the
messaging queue with unique credentials, it makes auditing much easier when
questionable data access is discovered: you can track it down to the specific
instance of a service based on the RabbitMQ username.

Vault makes use both of its own internal revocation system as well as the
deleting RabbitMQ users when creating RabbitMQ users to ensure that users
become invalid within a reasonable time of the lease expiring.

This page will show a quick start for this backend. For detailed documentation
on every path, use `vault path-help` after mounting the backend.

## Quick Start

The first step to using the RabbitMQ backend is to mount it. Unlike the
`kv` backend, the `rabbitmq` backend is not mounted by default.

```text
$ vault mount rabbitmq
Successfully mounted 'rabbitmq' at 'rabbitmq'!
```

Next, Vault must be configured to connect to the RabbitMQ. This is done by
writing the RabbitMQ management URI, RabbitMQ management administrator user,
and the user's password.

```text
$ vault write rabbitmq/config/connection \
    connection_uri="http://localhost:15672" \
    username="admin" \
    password="password"
```

In this case, we've configured Vault with the URI "http://localhost:15672",
user "admin", and password "password" connecting to a local RabbitMQ
management instance. It is important that the Vault user have the
administrator privilege to manager users.

Optionally, we can configure the lease settings for credentials generated
by Vault. This is done by writing to the `config/lease` key:

```
$ vault write rabbitmq/config/lease ttl=3600 max_ttl=86400
Success! Data written to: rabbitmq/config/lease
```

This restricts each credential to being valid or leased for 1 hour
at a time, with a maximum use period of 24 hours. This forces an
application to renew their credentials at least hourly, and to recycle
them once per day.

The next step is to configure a role. A role is a logical name that maps
to tags and virtual host permissions used to generated those credentials.
For example, lets create a "readwrite" virtual host role:

```text
$ vault write rabbitmq/roles/readwrite \
    vhosts='{"/":{"write": ".*", "read": ".*"}}'
Success! Data written to: rabbitmq/roles/readwrite
```

By writing to the `roles/readwrite` path we are defining the `readwrite` role.
This role will be created by evaluating the given `vhosts` and `tags` statements.
By default, no tags and no virtual hosts are assigned to a role. You can read more
about RabbitMQ management tags [here](https://www.rabbitmq.com/management.html#permissions).
Configure, write, and read permissions are granted per virtual host.

To generate a new set of credentials, we simply read from that role.
Vault is now configured to create and manage credentials for RabbitMQ!

```text
$ vault read rabbitmq/creds/readwrite
lease_id       rabbitmq/creds/readwrite/2740df96-d1c2-7140-c406-77a137fa3ecf
lease_duration 3600
lease_renewable	true
password       e1b6c159-ca63-4c6a-3886-6639eae06c30
username       root-4b95bf47-281d-dcb5-8a60-9594f8056092
```

By reading from the `creds/readwrite` path, Vault has generated a new
set of credentials using the `readwrite` role configuration. Here we
see the dynamically generated username and password, along with a one
hour lease.

Using ACLs, it is possible to restrict using the rabbitmq backend such
that trusted operators can manage the role definitions, and both
users and applications are restricted in the credentials they are
allowed to read.

If you get stuck at any time, simply run `vault path-help rabbitmq` or with a
subpath for interactive help output.

## API

The RabbitMQ secret backend has a full HTTP API. Please see the
[RabbitMQ secret backend API](/api/secret/rabbitmq/index.html) for more
details.
