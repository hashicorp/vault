---
layout: docs
page_title: Consul - Secrets Engines
description: The Consul secrets engine for Vault generates tokens for Consul dynamically.
---

# Consul Secrets Engine

@include 'x509-sha1-deprecation.mdx'

The Consul secrets engine generates [Consul](https://www.consul.io) API tokens
dynamically based on Consul ACL policies.

## Setup

Most secrets engines must be configured in advance before they can perform their
functions. These steps are usually completed by an operator or configuration
management tool.

1.  Enable the Consul secrets engine:

    ```shell-session
    $ vault secrets enable consul
    Success! Enabled the consul secrets engine at: consul/
    ```

    By default, the secrets engine will mount at the name of the engine. To
    enable the secrets engine at a different path, use the `-path` argument.

1. Vault can bootstrap the ACL system of your Consul cluster if it has
   not already been done. In this case, you only need the address of your
   Consul cluster to configure the Consul secret engine:

    ```text
    $ vault write consul/config/access \
        address=127.0.0.1:8500
    Success! Data written to: consul/config/access
    ```

    If you have already bootstrapped the ACL system of your Consul cluster, you
    will need to give Vault a management token:

    In Consul versions below 1.4, acquire a [management token][consul-mgmt-token] from Consul, using the
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

    For Consul 1.4 and above, use the command line to generate a token with the appropriate policy:

    ```shell-session
    $ CONSUL_HTTP_TOKEN="<management-token>" consul acl token create -policy-name="global-management"
    AccessorID:   865dc5e9-e585-3180-7b49-4ddc0fc45135
    SecretID:     ef35f0f1-885b-0cab-573c-7c91b65a7a7e
    Description:
    Local:        false
    Create Time:  2018-10-22 17:40:24.128188 -0700 PDT
    Policies:
        00000000-0000-0000-0000-000000000001 - global-management
    ```

1.  Configure Vault to connect and authenticate to Consul:

    ```shell-session
    $ vault write consul/config/access \
        address="127.0.0.1:8500" \
        token="7652ba4c-0f6e-8e75-5724-5e083d72cfe4"
    Success! Data written to: consul/config/access
    ```

1.  Configure a role that maps a name in Vault to a Consul ACL policy. Depending on your Consul version,
    you will either provide a policy document and a token_type, a list of policies or roles, or a set of
    service or node identities.
    When users generate credentials, they are generated against this role.

    For Consul versions below 1.4, the policy must be base64-encoded. The policy language is
    [documented by Consul](https://www.consul.io/docs/security/acl/acl-legacy).

    Write a policy and proceed to link it to the role:

    ```shell-session
    $ vault write consul/roles/my-role policy="$(base64 <<< 'key "" { policy = "read" }')"
    Success! Data written to: consul/roles/my-role
    ```

    For Consul versions 1.4 and above, [generate a policy in Consul](https://www.consul.io/docs/guides/acl.html),
    and proceed to link it to the role:

    ```shell-session
    $ vault write consul/roles/my-role policies="readonly"
    Success! Data written to: consul/roles/my-role
    ```

    For Consul versions 1.5 and above, [generate a role in Consul](https://www.consul.io/api/acl/roles) and
    proceed to link it to the role, or [attach a Consul service identity](https://www.consul.io/commands/acl/token/create#service-identity) to the role:

    ```shell-session
    $ vault write consul/roles/my-role consul_roles="api-server"
    Success! Data written to: consul/roles/my-role
    ```

    ```shell-session
    $ vault write consul/roles/my-role service_identities="myservice:dc1,dc2"
    Success! Data written to: consul/roles/my-role
    ```

    For Consul versions 1.8 and above, [attach a Consul node identity](https://www.consul.io/commands/acl/token/create#node-identity) to the role.

    ```shell-session
    $ vault write consul/roles/my-role node_identities="server-1:dc1"
    Success! Data written to: consul/roles/my-role
    ```

    -> **Token lease duration:** If you do not specify a value for `ttl` (or `lease` for Consul versions below 1.4) the
    tokens created using Vault's Consul secrets engine are created with a Time To Live (TTL) of 30 days. You can change
    the lease duration by passing `-ttl=<duration>` to the command above with "duration" being a string with a time
    suffix like "30s" or "1h".

1.  For Enterprise users, you may further limit a role's access by adding the optional parameters `consul_namespace` and/or
    `partition`. Please refer to Consul's [namespace documentation](https://www.consul.io/docs/enterprise/namespaces) and
    [admin partition documentation](https://www.consul.io/docs/enterprise/admin-partitions) for further information about
    these features.

    For Consul versions 1.7 and above, link a Consul namespace to the role:

    ```shell-session
    $ vault write consul/roles/my-role consul_roles="namespace-management" consul_namespace="ns1"
    Success! Data written to: consul/roles/my-role
    ```

    For Consul version 1.11 and above, link an admin partition to a role:

    ```shell-session
    $ vault write consul/roles/my-role consul_roles="admin-management" partition="admin1"
    Success! Data written to: consul/roles/my-role
    ```

## Usage

After the secrets engine is configured and a user/machine has a Vault token with
the proper permission, it can generate credentials.

Generate a new credential by reading from the `/creds` endpoint with the name
of the role:

```shell-session
$ vault read consul/creds/my-role
Key                 Value
---                 -----
lease_id            consul/creds/my-role/b2469121-f55f-53c5-89af-a3ba52b1d6d8
lease_duration      768h
lease_renewable     true
accessor            c81b9cf7-2c4f-afc7-1449-4e442b831f65
consul_namespace    ns1
local               false
partition           admin1
token               642783bf-1540-526f-d4de-fe1ac1aed6f0
```

!> **Expired token rotation:** Once a token's TTL expires, then Consul operations will no longer be allowed with it.
This requires you to have an external process to rotate tokens. At this time, the recommended approach for operators
is to rotate the tokens manually by creating a new token using the `vault read consul/creds/my-role` command. Once
the token is synchronized with Consul, apply the token to the agents using the Consul API or CLI.

## Tutorial

Refer to [Administer Consul Access Control Tokens with
Vault](https://learn.hashicorp.com/tutorials/consul/vault-consul-secrets) for a
step-by-step tutorial.

## API

The Consul secrets engine has a full HTTP API. Please see the
[Consul secrets engine API](/api-docs/secret/consul) for more
details.

[consul-mgmt-token]: https://www.consul.io/docs/agent/http/acl.html#acl_create
