---
layout: "docs"
page_title: "Consul - Secrets Engines"
sidebar_title: "Consul"
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

2. In Consul versions below 1.4, acquire a [management token][consul-mgmt-token] from Consul, using the
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
For Consul 1.4 and above, use the command line to generate a token with the appropiate policy:

   ```sh
   $ CONSUL_HTTP_TOKEN=d54fe46a-1f57-a589-3583-6b78e334b03b consul acl token create -policy-name=global-management
   AccessorID:   865dc5e9-e585-3180-7b49-4ddc0fc45135
   SecretID:     ef35f0f1-885b-0cab-573c-7c91b65a7a7e
   Description:
   Local:        false
   Create Time:  2018-10-22 17:40:24.128188 -0700 PDT
   Policies:
       00000000-0000-0000-0000-000000000001 - global-management
   ```

3. Configure Vault to connect and authenticate to Consul:

    ```text
    $ vault write consul/config/access \
        address=127.0.0.1:8500 \
        token=7652ba4c-0f6e-8e75-5724-5e083d72cfe4
    Success! Data written to: consul/config/access
    ```

4. Configure a role that maps a name in Vault to a Consul ACL policy. Depending on your Consul version, 
you will either provide a policy document and a token_type, or a set of policies.
When users generate credentials, they are generated against this role. For Consul versions below 1.4:

    ```text
    $ vault write consul/roles/my-role policy=$(base64 <<< 'key "" { policy = "read" }')
    Success! Data written to: consul/roles/my-role
    ```
The policy must be base64-encoded. The policy language is [documented by Consul](https://www.consul.io/docs/internals/acl.html).

For Consul versions 1.4 and above, [generate a policy in Consul](https://www.consul.io/docs/guides/acl.html), and procede to link it 
to the role:
    ```text
    $ vault write consul/roles/my-role policies=readonly
    Success! Data written to: consul/roles/my-role
    ```

## Usage

After the secrets engine is configured and a user/machine has a Vault token with
the proper permission, it can generate credentials.

Generate a new credential by reading from the `/creds` endpoint with the name
of the role:

```text
$ vault read consul/creds/my-role
Key                Value
---                -----
lease_id           consul/creds/my-role/b2469121-f55f-53c5-89af-a3ba52b1d6d8
lease_duration     768h
lease_renewable    true
token              642783bf-1540-526f-d4de-fe1ac1aed6f0
```

When using Consul 1.4, the response will include the accessor for the token

```text
$ vault read consul/creds/my-role
Key                Value
---                -----
lease_id           consul/creds/my-role/7miMPnYaBCaVWDS9clNE0Nv3
lease_duration     768h
lease_renewable    true
accessor           6d5a0348-dffe-e87b-4266-2bec03800abb
token              bc7a42c0-9c59-23b4-8a09-7173c474dc42
```
## API

The Consul secrets engine has a full HTTP API. Please see the
[Consul secrets engine API](/api/secret/consul/index.html) for more
details.

[consul-mgmt-token]: https://www.consul.io/docs/agent/http/acl.html#acl_create
