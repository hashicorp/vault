---
layout: "guides"
page_title: "Policies - Guides"
sidebar_current: "guides-configuration-policies"
description: |-
  Policies in Vault control what a user can access.
---

# Policies

In Vault, use policies to govern the behavior of clients and instrument
Role-Based Access Control (RBAC) by specifying access privileges
(_authorization_).

When you first initialize Vault, the
[**`root`**](/docs/concepts/policies.html#root-policy) policy gets created by
default. The `root` policy is a special policy that gives superuser access to
_everything_ in Vault. This allows the superuser to set up initial policies,
tokens, etc.

In addition, there is another built-in policy,
[**`default`**](/docs/concepts/policies.html#default-policy) gets created. The
`default` policy is attached to all tokens and provides common permissions.

Everything in Vault is path based, and admins write policies to grant or forbid
access to certain paths and operations in Vault. Vault operates on a **secure by
default** standard, and as such as empty policy grants **no permission** in the
system.


### HashiCorp Configuration Language (HCL)

Policies written in [HCL](https://github.com/hashicorp/hcl) format are often
referred as **_ACL Policies_**. [Sentinel](https://www.hashicorp.com/sentinel) is
another framework for policy which is available in [Vault
Enterprise](/docs/enterprise/index.html).  Since Sentinel is an enterprise-only
feature, this guide focuses on writing ACL policies as a foundation.

**NOTE:** HCL is JSON compatible; therefore, JSON can be used as completely
valid input.

## Reference Material

- [Policies](/docs/concepts/policies.html#default-policy) documentation
- [Policy API](/api/system/policy.html) documentation
- [Getting Started guide](/intro/getting-started/policies.html) on policies

## Estimated Time to Complete

10 minutes

## Personas

The scenario described in this guide introduces the following personas:

- **`root`** sets up initial policies for `admin`
- **`admin`** is empowered with managing a Vault infrastructure for a team or
organizations
- **`provisioner`** configures secret backends and creates policies for
client apps


## Challenge

Since Vault centrally secure, store, and access control secrets across
distributed infrastructure and applications, it is critical to control
permissions before any user or machine can gain access.


## Solution

Restrict the use of root policy, and write fine-grained policies to practice
**least privileged**. For example, if an app gets AWS credentials from Vault,
write policy grants to `read` from AWS secret backend but not to `delete`, etc.

Policies are attached to tokens and roles to enforce client permissions on
Vault.


## Prerequisites

To perform the tasks described in this guide, you need to have a Vault
environment.  Refer to the [Getting
Started](/intro/getting-started/install.html) guide to install Vault. Make sure
that your Vault server has been [initialized and
unsealed](/intro/getting-started/deploy.html).

### Policy requirements

Since this guide demonstrates the creation of an **`admin`** policy, log in with
**`root`** token if possible. Otherwise, make sure that you have the following
permissions:

```shell
# Manage auth backends broadly across Vault
path "auth/*"
{
  capabilities = ["create", "read", "update", "delete", "list", "sudo"]
}

# List, create, update, and delete auth backends
path "sys/auth/*"
{
  capabilities = ["create", "read", "update", "delete", "sudo"]
}

# To list policies - Step 3
path "sys/policy"
{
  capabilities = ["read"]
}

# Create and manage ACL policies broadly across Vault
path "sys/policy/*"
{
  capabilities = ["create", "read", "update", "delete", "list", "sudo"]
}

# List, create, update, and delete key/value secrets
path "secret/*"
{
  capabilities = ["create", "read", "update", "delete", "list", "sudo"]
}

# Manage and manage secret backends broadly across Vault.
path "sys/mounts/*"
{
  capabilities = ["create", "read", "update", "delete", "list", "sudo"]
}

# Read health checks
path "sys/health"
{
  capabilities = ["read", "sudo"]
}

# To perform Step 4
path "sys/capabilities"
{
  capabilities = ["create", "update"]
}

# To perform Step 4
path "sys/capabilities-self"
{
  capabilities = ["create", "update"]
}
```


## Steps

The basic workflow of creating policies is:

![Policy Creation Workflow](/assets/images/vault-policy-authoring-workflow.png)

This guide demonstrates basic policy authoring and management tasks.

1. [Write ACL policies in HCL format](#step1)
2. [Create policies](#step2)
3. [View existing policies](#step3)
4. [Check capabilities of a token](#step4)


### <a name="step1"></a>Step 1: Write ACL policies in HCL format

Remember, empty policy grants **no permission** in the system. Therefore, ACL
policies are defined for each path.

```shell
path "<PATH>" {
  capabilities = [ "<LIST_OF_CAPABILITIES>" ]
}
```

-> The path can have a wildcard ("`*`") specifying at the end to allow for
namespacing. For example, "`secret/training_*`" grants permissions on any
path starts with "`secret/training_`" (e.g. `secret/training_vault`).

Define one or more [capabilities](/docs/concepts/policies.html#capabilities) on each path to control operations that are
permitted.

| Capability      | Associated HTTP verbs  |
| --------------- |------------------------|
| create          | POST/PUT               |
| read            | GET                    |
| update          | POST/PUT               |
| delete          | DELETE                 |
| list            | LIST


#### Policy requirements

First step in creating policies is to **gather policy requirements**.

**Example:**

**`admin`** is a type of user empowered with managing a Vault infrastructure for
a team or organizations. Empowered with sudo, the Administrator is focused on
configuring and maintaining the health of Vault cluster(s) as well as
providing bespoke support to Vault users.

`admin` must be able to:

- Mount and manage auth backends broadly across Vault
- Mount and manage secret backends broadly across Vault
- Create and manage ACL policies broadly across Vault
- Read system health check

**`provisioner`** is a type of user or service that will be used by an automated
tool (e.g. Terraform) to provision and configure a namespace within a Vault
secret backend for a new Vault user to access and write secrets.

`provisioner` must be able to:

- Mount and manage auth backends
- Mount and manage secret backends
- Create and manage ACL policies


Now, you are ready to author policies to fulfill the requirements.

#### Example policy for admin

`admin-policy.hcl`

```shell
# Manage auth backends broadly across Vault
path "auth/*"
{
  capabilities = ["create", "read", "update", "delete", "list", "sudo"]
}

# List, create, update, and delete auth backends
path "sys/auth/*"
{
  capabilities = ["create", "read", "update", "delete", "sudo"]
}

# List existing policies
path "sys/policy"
{
  capabilities = ["read"]
}

# Create and manage ACL policies broadly across Vault
path "sys/policy/*"
{
  capabilities = ["create", "read", "update", "delete", "list", "sudo"]
}

# List, create, update, and delete key/value secrets
path "secret/*"
{
  capabilities = ["create", "read", "update", "delete", "list", "sudo"]
}

# Manage and manage secret backends broadly across Vault.
path "sys/mounts/*"
{
  capabilities = ["create", "read", "update", "delete", "list", "sudo"]
}

# Read health checks
path "sys/health"
{
  capabilities = ["read", "sudo"]
}
```

#### Example policy for provisioner

`provisioner-policy.hcl`

```shell
# Manage auth backends broadly across Vault
path "auth/*"
{
  capabilities = ["create", "read", "update", "delete", "list", "sudo"]
}

# List, create, update, and delete auth backends
path "sys/auth/*"
{
  capabilities = ["create", "read", "update", "delete", "sudo"]
}

# List existing policies
path "sys/policy"
{
  capabilities = ["read"]
}

# Create and manage ACL policies
path "sys/policy/*"
{
  capabilities = ["create", "read", "update", "delete", "list"]
}

# List, create, update, and delete key/value secrets
path "secret/*"
{
  capabilities = ["create", "read", "update", "delete", "list"]
}
```

### <a name="step2"></a>Step 2: Create policies

Now, create `admin` and `provisioner` policies in Vault.

#### CLI command

To create policies:

```shell
$ vault policy write <POLICY_NAME> <POLICY_FILE>
```

**Example:**

```shell
# Create admin policy
$ vault policy write admin admin-policy.hcl

# Create provisioner policy
$ vault policy write provisioner provisioner-policy.hcl
```

**NOTE:** To update an existing policy, simply re-run the same command by
passing your modified policy (`*.hcl`).

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

Now, create `admin` and `provisioner` policies:

```shell
# Create admin policy
$ curl --request PUT --header "X-Vault-Token: ..." --data @admin-payload.json \
    https://vault.rocks/v1/sys/policy/admin

$ cat admin-payload.json
{
  "policy": "path \"auth/*\" { capabilities = [\"create\", \"read\", \"update\", ... }"
}

# Create provisioner policy
$ curl --request PUT --header "X-Vault-Token: ..." --data @provisioner-payload.json \
    https://vault.rocks/v1/sys/policy/provisioner

$ cat provisioner-payload.json
{
  "policy": "path \"auth/*\" { capabilities = [\"create\", \"read\", \"update\", ... }"
}
```

-> NOTE: You can also use `/sys/policies` endpoint which is used to manage
ACL, RGP, and EGP policies in Vault (RGP and EGP policies are enterprise-only
features). To list policies, invoke `/sys/policies/acl` endpoint.

**NOTE:** To update an existing policy, simply re-run the same command by
passing your modified policy in the request payload (`*.json`).



### <a name="step3"></a>Step 3: View existing policies

Make sure that you see the policies you created in [Step 2](#step2).

#### CLI command

The following command lists existing policies:

```shell
$ vault policy list
```

To view a specific policy:

```shell
$ vault policy read <POLICY_NAME>
```

**Example:**

```shell
# Read admin policy
$ vault policy read admin

# Mount and manage auth backends broadly across Vault
path "auth/*"
{
  capabilities = ["create", "read", "update", "delete", "list", "sudo"]
}

path "sys/auth/*"
{
  capabilities = ["create", "read", "update", "delete", "sudo"]
}

# Create and manage ACL policies broadly across Vault
path "sys/policy/*"
{
  capabilities = ["create", "read", "update", "delete", "list", "sudo"]
}
...
```

#### API call using cURL

To list existing ACL policies, use the `/sys/policy` endpoint.

```shell
$ curl --request LIST --header "X-Vault-Token: ..." https://vault.rocks/v1/sys/policy | jq
```

To read a specific policy, the endpoint path should be
`/sys/policy/<POLICY_NAME>`.

**Example:**

Read the admin policy:

```plaintext
$ curl --request GET --header "X-Vault-Token: ..." https://vault.rocks/v1/sys/policy/admin | jq
{
  "name": "admin",
  "rules": "# Mount and manage auth backends broadly across Vault\npath \"auth/*\"\n{\n  ...",
  "request_id": "e8151bf3-8136-fef9-428b-1506042350cf",
  "lease_id": "",
  "renewable": false,
  "lease_duration": 0,
  "data": {
  ...
```

### <a name="step4"></a>Step 4: Check capabilities of a token

Use the `/sys/capabilities` endpoint to fetch the capabilities of a token on a
given path. This helps to verify what operations are granted based on the
policies attached to the token.

#### CLI command

The command is:

```shell
$ vault token capabilities <TOKEN> <PATH>
```

**Example:**

First, create a token attached to `admin` policy:

```shell
$ vault token create -policy="admin"
Key            	Value
---            	-----
token          	79ecdd41-9bac-1ac7-1ee4-99fbce796221
token_accessor 	39b5e8b5-7bbf-6c6d-c536-ba79d3a80dd5
token_duration 	768h0m0s
token_renewable	true
token_policies 	[admin default]
```

Now, fetch the capabilities of this token on `sys/auth/approle` path.

```plaintext
$ vault token capabilities 79ecdd41-9bac-1ac7-1ee4-99fbce796221 sys/auth/approle
Capabilities: [create delete read sudo update]
```

The result should match the policy rule you wrote on `sys/auth/*` path. You can
repeat the steps to generate a token for `provisioner` and check its
capabilities on paths.


In the absence of token, it returns capabilities of current token invoking this
command.

```shell
$ vault token capabilities sys/auth/approle
Capabilities: [root]
```

#### API call using cURL

Use the `sys/capabilities` endpoint.

**Example:**

First, create a token attached to `admin` policy:

```shell
$ curl --request POST --header "X-Vault-Token: ..." --data '{ "policies":"admin" }' \
       https://vault.rocks/v1/auth/token/create
{
 "request_id": "870ef38c-1401-7beb-633c-ff09cca3db68",
 "lease_id": "",
 "renewable": false,
 "lease_duration": 0,
 "data": null,
 "wrap_info": null,
 "warnings": null,
 "auth": {
   "client_token": "9f3a9fbb-4e1a-87c3-9d4d-ee4d96d40af1",
   "accessor": "f8a269c0-153a-c1ea-ae97-e7e964814392",
   "policies": [
     "root"
   ],
   "metadata": null,
   "lease_duration": 0,
   "renewable": false,
   "entity_id": ""
 }
}
```

Now, fetch the capabilities of this token on `sys/auth/approle` path.

```shell
# Request payload
$ cat payload.json
{
  "token": "9f3a9fbb-4e1a-87c3-9d4d-ee4d96d40af1",
  "path": "sys/auth/approle"
}

$ curl --request POST --header "X-Vault-Token: ..." --data @payload.json \
    https://vault.rocks/v1/sys/capabilities
{
  "capabilities": [
    "create",
    "delete",
    "read",
    "sudo",
    "update"
  ],
  "request_id": "03f9d5e2-7e8a-4cd3-b9e9-034c058d3d06",
  "lease_id": "",
  "renewable": false,
  "lease_duration": 0,
  "data": {
    "capabilities": [
      "create",
      "delete",
      "read",
      "sudo",
      "update"
    ]
  },
  "wrap_info": null,
  "warnings": null,
  "auth": null
}
```

The result should match the policy rule you wrote on `sys/auth/*` path. You can
repeat the steps to generate a token for `provisioner` and check its
capabilities on paths.

To check current token's capabilities permitted on a path, use
`sys/capabilities-self` endpoint.

```plaintext
$ curl --request POST --header "X-Vault-Token: ..." --data '{"path":"sys/auth/approle"}' \
    https://vault.rocks/v1/sys/capabilities-self
```


## Next steps

In this guide, you learned how to write policies based on given policy
requirements. Next, [AppRole Pull Authentication](/guides/configuration/authentication.html)
guide demonstrates how to associate policies to a role.
