---
layout: "docs"
page_title: "Consul - Secrets Engines"
sidebar_current: "docs-secrets-consul"
description: |-
  The Consul secrets engine for Vault generates tokens for Consul dynamically.
---

# Consul Secrets Engine

The Consul secrets engine generates [Consul](https://www.consul.io) API tokens
dynamically based on Consul ACL policies.

## Setup

Most secrets engines must be configured in advance before they can perform their
functions. These steps are usually completed by an operator or configuration
management tool.

1. Enable the Consul secrets engine:

    ```text
    $ vault secrets enable consul
    Success! Enabled the consul secrets engine at: consul/
    ```

    By default, the secrets engine will mount at the name of the engine. To
    enable the secrets engine at a different path, use the `-path` argument.

1. Acquire a [management token][consul-mgmt-token] from Consul, using the
`acl_master_token` from your Consul configuration file or another management
token:

    ```sh
    $ curl \
        --header "X-Consul-Token: my-management-token" \
        --request PUT \
        --data '{"Name": "sample", "Type": "management"}' \
        https://consul.rocks/v1/acl/create
    ```

    Vault must have a management type token so that it can create and revoke ACL
    tokens. The response will return a new token:

    ```json
    {
      "ID": "7652ba4c-0f6e-8e75-5724-5e083d72cfe4"
    }
    ```

1. Configure Vault to connect and authenticate to Consul:

    ```text
    $ vault write consul/config/access \
        address=127.0.0.1:8500 \
        token=7652ba4c-0f6e-8e75-5724-5e083d72cfe4
    Success! Data written to: consul/config/access
    ```

1. Configure a role that maps a name in Vault to a Consul ACL policy.
When users generate credentials, they are generated against this role:

    ```text
    $ vault write consul/roles/my-role policy=$(base64 <<< 'key "" { policy = "read" }')
    Success! Data written to: consul/roles/my-role
    ```

    The policy must be base64-encoded. The policy language is [documented by
    Consul](https://www.consul.io/docs/internals/acl.html).

## Usage

After the secrets engine is configured and a user/machine has a Vault token with
the proper permission, it can generate credentials.

1. Generate a new credential by reading from the `/creds` endpoint with the name
of the role:

    ```text
    $ vault read consul/creds/readonly
    Key                Value
    ---                -----
    lease_id           consul/creds/my-role/b2469121-f55f-53c5-89af-a3ba52b1d6d8
    lease_duration     768h
    lease_renewable    true
    token              642783bf-1540-526f-d4de-fe1ac1aed6f0
    ```

## API

The Consul secrets engine has a full HTTP API. Please see the
[Consul secrets engine API](/api/secret/consul/index.html) for more
details.

[consul-mgmt-token]: https://www.consul.io/docs/agent/http/acl.html#acl_create
