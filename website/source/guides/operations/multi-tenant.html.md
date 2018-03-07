---
layout: "guides"
page_title: "Multi-Tenant Pattern - Guides"
sidebar_current: "guides-operations-multi-tenant"
description: |-
  This guide provides guidance in creating a multi-tenant environment.
---

# Multi-Tenant Pattern with Namespaces

Everything in Vault is path-based, and often use the terms `path` and `namespace`
interchangeably. The application namespace pattern is a useful construct for
providing Vault as a service to internal customers, giving them the ability to
leverage a multi-tenant Vault implementation with full agency to their
application's interactions with Vault.


## Reference Material

- [Static Secrets](/guides/secret-mgmt/static-secrets.html) guide
- [Policies](/guides/identity/policies.html) guide


## Estimated Time to Complete

10 minutes

## Personas

The end-to-end scenario described in this guide involves two personas:

- **`admin`** with privileged permissions
- **`user`** trusted user writing and reading secrets from Vault

## Challenge

When Vault is primarily used as a central location to manage secrets, it might
become necessary to design _Vault as a service_ to serve multiple organizations
within a company.  In a Vault as a service scenario, you want to create a
multi-tenant environment so that each team, organization, or app can have an
isolated namespace where they can perform all necessary tasks within their
dedicated namespace, and not interfering others.

![Vault as a service](/assets/images/vault-multi-tenant.png)

## Solution

As policies in Vault are controlled by path-based concepts, providing a
top-level application or team namespace allows for the creation of a relative
root of Vault control to contain a team's Vault utilization.

The creation of the top-level namespace is an **implicit action**. There is no
step in which you "mount" or create a namespace; however, the concept of the
namespace is instantiated by enabling Vault secret engines.

![Vault as a service](/assets/images/vault-multi-tenant-2.png)


## Prerequisites

To perform the tasks described in this guide, you need to have a Vault
environment.  Refer to the [Getting
Started](/intro/getting-started/install.html) guide to install Vault. Make sure
that your Vault server has been [initialized and
unsealed](/intro/getting-started/deploy.html).

### <a name="policy"></a>Policy requirements

-> **NOTE:** For the purpose of this guide, you can use **`root`** token to work
with Vault. However, it is recommended that root tokens are only used for just
enough initial setup or in emergencies. As a best practice, use tokens with
appropriate set of policies based on your role in the organization.

To perform all tasks demonstrated in this guide, your policy must include the
following permissions:

```shell
# Manage secret engines
path "sys/mounts/*" {
  capabilities = [ "create", "read", "update", "delete", "list", "sudo" ]
}

# List enabled secret engines
path "sys/mounts/" {
  capabilities = [ "list", "read" ]
}

# Create auth method
path "sys/auth/*" {
  capabilities = [ "create", "read", "update", "delete", "list", "sudo" ]  
}

# Write ACL policies
path "sys/policy/*" {
  capabilities = [ "create", "read", "update", "delete", "list" ]
}

# Create and manage userpass auth method
path "auth/userpass/*" {
  capabilities = [ "create", "read", "update", "delete", "list" ]
}
```

If you are not familiar with policies, complete the
[policies](/guides/identity/policies.html) guide.



## Steps

In this guide, you are going to perform the following steps:

1. [Enabling Secret Engines](#step1)
2. [Creating Policies](#step2)
3. [Enabling Authentication Methods](#step3)
4. [Verification - Writing and Reading Secrets](#step4)


### <a name="step1"></a>Step 1: Enabling Secret Engines
(**Persona:** admin)

In this step, you are going to create a top-level namespace (`path`) for an
organization.

#### CLI command

To enable a secret engine at a specific path:

```shell
$ vault secrets enable -path=<PATH> <TYPE>
````

**Example:**

```bash
# Enabling key-value secret engine at finance/ path
$ vault secrets enable -path=finance/ -description="Dedicated to the Finance org" kv

# Secret engine types can be mounted multiple times under different paths.
# This shows a second key-value secret engine mount
$ vault secrets enable -path=marketing/ -description="Dedicated to the Marketing org" kv

# List all secret engine mount locations
$ vault secrets list
Path          Type         Description
----          ----         -----------
cubbyhole/    cubbyhole    per-token private secret storage
finance/      kv           Dedicated to the Finance org
identity/     identity     identity store
marketing/    kv           Dedicated to the Marketing org
secret/       kv           key/value secret storage
sys/          system       system endpoints used for control, policy and debugging
```

#### API call using cURL

To enable a secret engine at a specific path via API:

```shell
$ curl --header "X-Vault-Token: <TOKEN>" \
       --request POST \
       --data <PARAMETERS> \
       <VAULT_ADDRESS>/v1/sys/mounts/<PATH>
```

Where `<TOKEN>` is your valid token, and `<PARAMETERS>` holds [configuration
parameters](/api/system/mounts.html#mount-secret-backend) of the secret engine.

**Example:**

```shell
# Enabling key-value secret engine at finance/ path
$ curl --header "X-Vault-Token: ..." \
       --request POST \
       --data '{"type":"kv", "description": "Dedicated to the Finance org"}' \
       https://vault.rocks/v1/sys/mounts/finance/

# Secret engine types can be mounted multiple times under different paths.
# This shows a second key-value secret engine mount
$ curl --header "X-Vault-Token: ..." \
       --request POST \
       --data '{"type":"kv", "description": "Dedicated to the Marketing org"}' \
       https://vault.rocks/v1/sys/mounts/marketing/

# List all secret engine mount locations
$ curl --header "X-Vault-Token: ..." \
      --request GET \
      https://vault.rocks/v1/sys/mounts | jq
{
  ...
  "finance/": {
    "accessor": "kv_b34cd82b",
    "config": {
      "default_lease_ttl": 0,
      "force_no_cache": false,
      "max_lease_ttl": 0,
      "plugin_name": ""
    },
    "description": "Dedicated for the Finance org",
    "local": false,
    "seal_wrap": false,
    "type": "kv"
  },
  "marketing/": {
    "accessor": "kv_90c8b5ec",
    "config": {
      "default_lease_ttl": 0,
      "force_no_cache": false,
      "max_lease_ttl": 0,
      "plugin_name": ""
    },
    "description": "Dedicated for the Marketing org",
  ...
  },
}
```


### <a name="step2"></a>Step 2: Creating Policies
(**Persona:** admin)

This is a policy for a namespace administrator which includes full permissions
under the path, as well as full permissions for namespaced policies under
`/sys/<PATH>`.

#### Author a policy file

`finance-admins.hcl`

```shell
# Full permissions on the finance path
path "finance" {
   capabilities = ["create", "read", "update", "delete", "list", "sudo"]
}

# Full permissions on the descendant paths under finance
path "finance/*" {
   capabilities = ["create", "read", "update", "delete", "list", "sudo"]
}

# Full permissions to manage policies for finance
path "sys/policy/finance/*" {
   capabilities = ["create", "read", "update", "delete", "list", "sudo"]
}
```

#### CLI command

To create policies:

```shell
$ vault policy write <POLICY_NAME> <POLICY_FILE>
```

**Example:**

```shell
# Create finance-admins policy
$ vault policy write finance-admins finance-admins.hcl
```

#### API call using cURL

To create a policy, use `/sys/policy` endpoint:

```shell
$ curl --header "X-Vault-Token: <TOKEN>" \
       --request PUT \
       --data <PAYLOAD> \
       <VAULT_ADDRESS>/v1/sys/policy/<POLICY_NAME>
```

Where `<TOKEN>` is your valid token, and `<PAYLOAD>` includes policy name and
stringfied policy.

**Example:**

```shell
# Create finance-admins policy
$ curl --request PUT --header "X-Vault-Token: ..." --data @payload.json \
    https://vault.rocks/v1/sys/policy/finance-admins

$ cat payload.json
{
  "policy": "path \"finance/*\" { capabilities = [\"create\", \"read\", \"update\", ... }"
}
```


### Additional Policy Example

Policy for the key-value secret engine for a "secret writer" without delete
permission would require administrative intervention to delete the secret. Also
includes a fine-grained conditional to prevent an update to the key named
"value" being set to an empty string. This would prevent an update action from
replacing the previous secret’s value key with an empty value.

#### Author a Policy File

`finance-secret-writer.hcl`:

```plaintext
path "finance/*" {
   capabilities = ["create", "update"]
   allowed_parameters = {
       "*" = []
   }
   denied_parameters = {
       "value" = [""]
   }
}
```

Policy for `kv` secret engine for machine consumption, limited only to the
ability to read secrets.

```plaintext
path "finance/*" {
   capabilities = ["read"]
}
```

Policy for `kv` secret engine for human consumption without protection against
deletion.

```plaintext
path "marketing/*" {
   capabilities = ["create", "update", "read", "delete" ]
}
```



### <a name="step3"></a>Step 3: Enabling Authentication Methods
(**Persona:** admin)

For the Vault clients to authenticate, authentication method(s) need to be
enabled and configured.

#### CLI command

To enable an auth method by executing the following command:

```shell
$ vault auth enable -path=<MOUNT_POINT> <AUTH_TYPE>
```

**Example:**

```shell
# Enable `userpass` auth method at `userpass` path
$ vault auth enable userpass
```


Let's associate the `finance-admins` policy to a user using the `userpass` auth
method.

```shell
# Create finance-admin user with finance-admins Policy
$ vault write auth/userpass/users/finance-admin password="p@ssw0rd" policies="finance-admins"

# Login as finance-admin
$ vault login -method=userpass username=finance-admin
Password (will be hidden):

Success! You are now authenticated. The token information displayed below
is already stored in the token helper. You do NOT need to run "vault login"
again. Future Vault requests will automatically use this token.

Key                    Value
---                    -----
token                  1f0a88b4-fdbb-6329-ab23-6cdde5360aea
token_accessor         2af2c52b-eed9-d73d-1537-d619d34f8946
token_duration         768h
token_renewable        true
token_policies         [finance-admins default]
token_meta_username    finance-admin
```

#### API call using cURL

To enable an auth method via API, execute the following command:

```shell
$ curl --header "X-Vault-Token: <TOKEN>" \
       --request POST \
       --data <PARAMETERS> \
       <VAULT_ADDRESS>/v1/sys/auth/<path>
```

Where `<TOKEN>` is your valid token, and `<PARAMETERS>` holds [configuration
parameters](/api/system/auth.html#mount-auth-backend) of the auth method.

**Example:**

```shell
# Enable `userpass` auth method at `userpass` path
$ curl --header "X-Vault-Token: ..." \
       --request POST \
       --data '{"type": "userpass"}' \
       https://vault.rocks/v1/sys/auth/userpass
```


Let's associate the `finance-admins` policy to a user using the `userpass` auth
method.

```shell
# Create finance-admin user with finance-admins Policy
$ curl --request POST --header "X-Vault-Token: ..." \
       --data '{"password": "p@ssw0rd", "policies": "finance-admins"}' \
       https://vault.rocks/v1/auth/userpass/users/finance-admin

# Login as finance-admin
$ curl --request POST --data '{"password": "p@ssw0rd"}' \
       https://vault.rocks/v1/auth/userpass/login/finance-admin | jq
{
 "request_id": "e17d8a48-7d43-af1c-c6d3-109c302a564d",
 "lease_id": "",
 "renewable": false,
 "lease_duration": 0,
 "data": null,
 "wrap_info": null,
 "warnings": null,
 "auth": {
   "client_token": "84ffc3de-aaf6-2acd-aec9-4b95847fe408",
   "accessor": "7ccd8023-3383-82d7-f0df-048e2b30c281",
   "policies": [
     "finance-admins",
     "default"
   ],
   "metadata": {
     "username": "finance-admin"
   },
   "lease_duration": 2764800,
   "renewable": true,
   "entity_id": "a98bb4ba-6603-78c4-17c0-b902422d5205"
 }
}
```


### <a name="step4"></a>Step 4: Verification - Writing and Reading Secrets
(**Persona:** user)

Now, you are logged in as `finance-admin` user.  Remember from [Step 2](#step2)
that the `finance-admins` policy permits all capabilities on `finance/` and its
descendant paths.

#### CLI command

To create key/value secrets:

```shell
$ vault write <PATH> <KEY>=VALUE>
```

The command to read secret is:

```shell
$ vault read secret/<PATH>
```

**Example:**

```shell
# Write secret
$ vault write finance/secrets admin_pass="S3cr\#t"
Success! Data written to: finance/secrets

# Read the secret back
$ vault read finance/secrets
Key                 Value
---                 -----
refresh_interval    768h
admin_pass          S3cr\#t
```

**Create a Policy**

The `finance-admins` policy permits the creation of policies under `finance`;
therefore, you should be able to run the following command successfully.

```shell
$ vault policy write finance/finance-secret-writer finance-secret-writer.hcl
Success! Uploaded policy: finance/finance-secret-writer
```


#### API call using cURL

Use `<PATH>` endpoint to create secrets:

```shell
$ curl --header "X-Vault-Token: <TOKEN>" \
       --request POST \
       --data <SECRETS> \
       <VAULT_ADDRESS>/v1/<PATH>
```

Use `<PATH>` endpoint to retrieve secrets:

```shell
$ curl --header "X-Vault-Token: <TOKEN>" \
       --request GET \
       <VAULT_ADDRESS>/v1/<PATH>
```

**Example:**

```shell
# Write secret
$ curl --header "X-Vault-Token: ..." \
       --request POST \
       --data '{"admin_pass": "S3cr\#t"}'
       https://vault.rocks/v1/finance/secrets

# Read the secret back
$ curl --header "X-Vault-Token: ..." \
       --request GET \
       https://vault.rocks/v1/finance/secrets | jq
{
 "request_id": "7b4e0329-6b4e-962b-7901-a9144916f4ae",
 "lease_id": "",
 "renewable": false,
 "lease_duration": 2764800,
 "data": {
   "admin_pass": "S3cr\#t"
 },
 "wrap_info": null,
 "warnings": null,
 "auth": null
}
```

**Create a Policy**

The `finance-admins` policy permits the creation of policies under `finance`;
therefore, you should be able to run the following command successfully.

```shell
# Create finance-secret-writer policy
$ curl --request PUT --header "X-Vault-Token: ..." --data @payload.json \
    https://vault.rocks/v1/sys/policy/finance/finance-secret-writer

$ cat payload.json
{
  "policy": "path \"finance/*\" { capabilities = [\"create\", \"update\"] allowed_parameters = ... }"
}
```


### Additional Discussion

You might want to create an orphan token to own the top-level
root/administrative concept for the app namespace to prevent it from being
revoked unexpectedly.

```shell
$ vault tokenn create -orphan -display-name="finance-root" -policy="finance-admins"
```

For more information about the orphan tokens, refer to:

- [Token Hierarchies and Orphan Tokens](/docs/concepts/tokens.html\#token-hierarchies-and-orphan-tokens)
- [Tokens and Leases](/guides/identity/lease.html#step5)


## Next steps

The use of [AppRole Pull Authentication](/guides/identity/authentication.html) is a good
use case to leverage the response wrapping. Go through the guide if you have not
done so.  To better understand the lifecycle of Vault tokens, proceed to [Tokens
and Leases](/guides/identity/lease.html) guide.
