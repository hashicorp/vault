---
layout: "docs"
page_title: "RabbitMQ - Secrets Engines"
sidebar_current: "docs-secrets-rabbitmq"
description: |-
  The RabbitMQ secrets engine for Vault generates user credentials to access RabbitMQ.
---

# RabbitMQ Secrets Engine

The RabbitMQ secrets engine generates user credentials dynamically based on
configured permissions and virtual hosts. This means that services that need to
access a virtual host no longer need to hardcode credentials.

With every service accessing the messaging queue with unique credentials,
auditing is much easier when questionable data access is discovered. Easily
track issues down to a specific instance of a service based on the RabbitMQ
username.

Vault makes use both of its own internal revocation system as well as the
deleting RabbitMQ users when creating RabbitMQ users to ensure that users become
invalid within a reasonable time of the lease expiring.

## Setup

Most secrets engines must be configured in advance before they can perform their
functions. These steps are usually completed by an operator or configuration
management tool.

1. Enable the RabbitMQ secrets engine:

    ```text
    $ vault secrets enable rabbitmq
    Success! Enabled the rabbitmq secrets engine at: rabbitmq/
    ```

    By default, the secrets engine will mount at the name of the engine. To
    enable the secrets engine at a different path, use the `-path` argument.

1. Configure the credentials that Vault uses to communicate with RabbitMQ to
generate credentials:

    ```text
    $ vault write rabbitmq/config/connection \
        connection_uri="http://localhost:15672" \
        username="admin" \
        password="password"
    Success! Data written to: rabbitmq/config/connection
    ```

    It is important that the Vault user have the administrator privilege to
    manager users.


1. Configure a role that maps a name in Vault to virtual host permissions:

    ```text
    $ vault write rabbitmq/roles/my-role \
        vhosts='{"/":{"write": ".*", "read": ".*"}}'
    Success! Data written to: rabbitmq/roles/my-role
    ```

    By writing to the `roles/my-role` path we are defining the `my-role` role.
    This role will be created by evaluating the given `vhosts` and `tags`
    statements. By default, no tags and no virtual hosts are assigned to a role.
    You can read more about [RabbitMQ management tags][rmq-perms].

## Usage

After the secrets engine is configured and a user/machine has a Vault token with
the proper permission, it can generate credentials.

1. Generate a new credential by reading from the `/creds` endpoint with the name
of the role:

    ```text
    $ vault read rabbitmq/creds/my-role
    Key                Value
    ---                -----
    lease_id           rabbitmq/creds/my-role/37d70d04-f24d-760a-e06e-b9b21087f0f4
    lease_duration     768h
    lease_renewable    true
    password           a98af72b-b6c9-b4b1-fe37-c73a572befed
    username           token-590f1fe2-1094-a4d6-01a7-9d4ff756a085
    ```

    Using ACLs, it is possible to restrict using the rabbitmq secrets engine
    such that trusted operators can manage the role definitions, and both users
    and applications are restricted in the credentials they are allowed to read.

## API

The RabbitMQ secrets engine has a full HTTP API. Please see the
[RabbitMQ secrets engine API](/api/secret/rabbitmq/index.html) for more
details.

[rmq-perms]: https://www.rabbitmq.com/management.html#permissions
