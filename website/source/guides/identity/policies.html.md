---
layout: "guides"
page_title: "Policies - Guides"
sidebar_title: "Policies"
sidebar_current: "guides-identity-policies"
description: |-
  Policies in Vault control what a user can access.
---

# Policies

In Vault, we use policies to govern the behavior of clients and instrument
Role-Based Access Control (RBAC) by specifying access privileges
(_authorization_).

When you first initialize Vault, the
[**`root`**](/docs/concepts/policies.html#root-policy) policy gets created by
default. The `root` policy is a special policy that gives superuser access to
_everything_ in Vault. This allows the superuser to set up initial policies,
tokens, etc.

In addition, another built-in policy,
[**`default`**](/docs/concepts/policies.html#default-policy), is created. The
`default` policy is attached to all tokens and provides common permissions.

Everything in Vault is path based, and admins write policies to grant or forbid
access to certain paths and operations in Vault. Vault operates on a **secure by
default** standard, and as such, an empty policy grants **no permissions** in the
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

~> **NOTE:** An [interactive
tutorial](https://www.katacoda.com/hashicorp/scenarios/vault-policies) is
also available if you do not have a Vault environment to perform the steps
described in this guide.


## Estimated Time to Complete

10 minutes

## Personas

The scenario described in this guide introduces the following personas:

- **`root`** sets up initial policies for `admin`
- **`admin`** is empowered with managing a Vault infrastructure for a team or
organizations
- **`provisioner`** configures secret engines and creates policies for
client apps


## Challenge

Since Vault centrally secures, stores, and controls access to secrets across
distributed infrastructure and applications, it is critical to control
permissions before any user or machine can gain access.


## Solution

Restrict the use of root policy, and write fine-grained policies to practice
**least privileged**. For example, if an app gets AWS credentials from Vault,
write policy grants to `read` from AWS secret engine but not to `delete`, etc.

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
# Manage auth methods broadly across Vault
path "auth/*"
{
  capabilities = ["create", "read", "update", "delete", "list", "sudo"]
}

# Create, update, and delete auth methods
path "sys/auth/*"
{
  capabilities = ["create", "update", "delete", "sudo"]
}

# List auth methods
path "sys/auth"
{
  capabilities = ["read"]
}

# Create and manage ACL policies via CLI
path "sys/policy/*"
{
  capabilities = ["create", "read", "update", "delete", "list", "sudo"]
}

# Create and manage ACL policies via API
path "sys/policies/acl/*"
{
  capabilities = ["create", "read", "update", "delete", "list", "sudo"]
}

# To list policies - Step 3
path "sys/policy"
{
  capabilities = ["read"]
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

# List, create, update, and delete key/value secrets
path "secret/*"
{
  capabilities = ["create", "read", "update", "delete", "list", "sudo"]
}

# Manage secret engines broadly across Vault
path "sys/mounts/*"
{
  capabilities = ["create", "read", "update", "delete", "list", "sudo"]
}

# List existing secret engines
path "sys/mounts"
{
  capabilities = ["read"]
}

# Read health checks
path "sys/health"
{
  capabilities = ["read", "sudo"]
}
```


## Steps

The basic workflow of creating policies is:

![Policy Creation Workflow](/img/vault-policy-authoring-workflow.png)

This guide demonstrates basic policy authoring and management tasks.

1. [Write ACL policies in HCL format](#step1)
1. [Create policies](#step2)
1. [View existing policies](#step3)
1. [Check capabilities of a token](#step4)


### <a name="step1"></a>Step 1: Write ACL policies in HCL format

Remember, an empty policy grants **no permission** in the system. Therefore, ACL
policies are defined for each path.

```shell
path "<PATH>" {
  capabilities = [ "<LIST_OF_CAPABILITIES>" ]
}
```

-> The path can have a wildcard ("`*`") at the end to allow for
namespacing. For example, "`secret/training_*`" grants permissions on any
path starting with "`secret/training_`" (e.g. `secret/training_vault`).

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

The first step in creating policies is to **gather policy requirements**.

**Example:**

**`admin`** is a type of user empowered with managing a Vault infrastructure for
a team or organizations. Empowered with sudo, the Administrator is focused on
configuring and maintaining the health of Vault cluster(s) as well as
providing bespoke support to Vault users.

`admin` must be able to:

- Enable and manage auth methods broadly across Vault
- Manage the key/value secret engines at `secret/` path
- Create and manage ACL policies broadly across Vault
- Read system health check

**`provisioner`** is a type of user or service that will be used by an automated
tool (e.g. Terraform) to provision and configure a namespace within a Vault
secret engine for a new Vault user to access and write secrets.

`provisioner` must be able to:

- Enable and manage auth methods
- Manage the key/value secret engines at `secret/` path
- Create and manage ACL policies


Now, you are ready to author policies to fulfill these requirements.

#### Example policy for admin

`admin-policy.hcl`

```shell
# Manage auth methods broadly across Vault
path "auth/*"
{
  capabilities = ["create", "read", "update", "delete", "list", "sudo"]
}

# Create, update, and delete auth methods
path "sys/auth/*"
{
  capabilities = ["create", "update", "delete", "sudo"]
}

# List auth methods
path "sys/auth"
{
  capabilities = ["read"]
}

# List existing policies via CLI
path "sys/policy"
{
  capabilities = ["read"]
}

# Create and manage ACL policies via CLI
path "sys/policy/*"
{
  capabilities = ["create", "read", "update", "delete", "list", "sudo"]
}

# Create and manage ACL policies via API
path "sys/policies/acl/*"
{
  capabilities = ["create", "read", "update", "delete", "list", "sudo"]
}

# List, create, update, and delete key/value secrets
path "secret/*"
{
  capabilities = ["create", "read", "update", "delete", "list", "sudo"]
}

# Manage secret engines broadly across Vault
path "sys/mounts/*"
{
  capabilities = ["create", "read", "update", "delete", "list", "sudo"]
}

# List existing secret engines
path "sys/mounts"
{
  capabilities = ["read"]
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
# Manage auth methods broadly across Vault
path "auth/*"
{
  capabilities = ["create", "read", "update", "delete", "list", "sudo"]
}

# Create, update, and delete auth methods
path "sys/auth/*"
{
  capabilities = ["create", "update", "delete", "sudo"]
}

# List auth methods
path "sys/auth"
{
  capabilities = ["read"]
}

# List existing policies via CLI
path "sys/policy"
{
  capabilities = ["read"]
}

# Create and manage ACL policies via CLI
path "sys/policy/*"
{
  capabilities = ["create", "read", "update", "delete", "list", "sudo"]
}

# Create and manage ACL policies via API
path "sys/policies/acl/*"
{
  capabilities = ["create", "read", "update", "delete", "list", "sudo"]
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

To create a policy, use the
[`sys/policies/acl`](/api/system/policies.html#create-update-acl-policy)
endpoint:

```shell
$ curl --header "X-Vault-Token: <TOKEN>" \
       --request PUT \
       --data <PAYLOAD> \
       <VAULT_ADDRESS>/v1/sys/policies/acl/<POLICY_NAME>
```

Where `<TOKEN>` is your valid token, and `<PAYLOAD>` includes the policy name and
stringified policy.

-> **NOTE:** To create ACL policies, you can use the
[`sys/policy`](/api/system/policy.html) endpoint as well.

**Example:**

```shell
# Create the API request payload. Use stringified policy expression.
$ tee admin-payload.json <<EOF
{
  "policy": "# Manage auth methods broadly across Vault\npath \"auth/*\"\n{\n  capabilities = [\"create\", \"read\", \"update\", \"delete\", \"list\", \"sudo\"]\n}\n\n# List, create, update, and delete auth methods\npath \"sys/auth/*\"\n{\n  capabilities = [\"create\", \"read\", \"update\", \"delete\", \"sudo\"]\n}\n\n# List auth methods\npath \"sys/auth\"\n{\n  capabilities = [\"read\"]\n}\n\n# List existing policies\npath \"sys/policies\"\n{\n  capabilities = [\"read\"]\n}\n\n# Create and manage ACL policies broadly across Vault\npath \"sys/policies/*\"\n{\n  capabilities = [\"create\", \"read\", \"update\", \"delete\", \"list\", \"sudo\"]\n}\n\n# List, create, update, and delete key/value secrets\npath \"secret/*\"\n{\n  capabilities = [\"create\", \"read\", \"update\", \"delete\", \"list\", \"sudo\"]\n}\n\n# Manage and manage secret engines broadly across Vault.\npath \"sys/mounts/*\"\n{\n  capabilities = [\"create\", \"read\", \"update\", \"delete\", \"list\", \"sudo\"]\n}\n\n# List existing secret engines.\npath \"sys/mounts\"\n{\n  capabilities = [\"read\"]\n}\n\n# Read health checks\npath \"sys/health\"\n{\n  capabilities = [\"read\", \"sudo\"]\n}"
}
EOF

# Create admin policy
$ curl --header "X-Vault-Token: ..." \
       --request PUT \
       --data @admin-payload.json \
       http://127.0.0.1:8200/v1/sys/policies/acl/admin

# Create the API requset payload for creating provisioner policy
$ tee provisioner-payload.json <<EOF
{
  "policy": "# Manage auth methods broadly across Vault\npath \"auth/*\"\n{\n  capabilities = [\"create\", \"read\", \"update\", \"delete\", \"list\", \"sudo\"]\n}\n\n# List, create, update, and delete auth methods\npath \"sys/auth/*\"\n{\n  capabilities = [\"create\", \"read\", \"update\", \"delete\", \"sudo\"]\n}\n\n# List existing policies\npath \"sys/policy\"\n{\n  capabilities = [\"read\"]\n}\n\n# Create and manage ACL policies\npath \"sys/policy/*\"\n{\n  capabilities = [\"create\", \"read\", \"update\", \"delete\", \"list\"]\n}\n\n# List, create, update, and delete key/value secrets\npath \"secret/*\"\n{\n  capabilities = [\"create\", \"read\", \"update\", \"delete\", \"list\"]\n}"
}
EOF

# Create provisioner policy
$ curl --header "X-Vault-Token: ..." \
       --request PUT \
       --data @provisioner-payload.json \
       http://127.0.0.1:8200/v1/sys/policies/acl/provisioner
```

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

# Mount and manage auth methods broadly across Vault
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

To list existing ACL policies, use the `sys/policies/acl` endpoint.

```shell
$ curl --request LIST --header "X-Vault-Token: ..." http://127.0.0.1:8200/v1/sys/policies/acl | jq
```

To read a specific policy, the endpoint path should be
`sys/policies/acl/<POLICY_NAME>`.

-> **NOTE:** To read existing ACL policies, you can use the `sys/policy`
endpoint as well.

**Example:**

```shell
# Read the admin policy
$ curl --header "X-Vault-Token: ..." http://127.0.0.1:8200/v1/sys/policies/acl/admin | jq
{
  "request_id": "3f826e5c-70a0-2998-8082-fe34c67c59d1",
  "lease_id": "",
  "renewable": false,
  "lease_duration": 0,
  "data": {
    "name": "admin",
    "policy": "# Manage auth methods broadly across Vault\npath \"auth/*\"\n{\n  capabilities = [\"create\", \"read\" ...
  },
  "wrap_info": null,
  "warnings": null,
  "auth": null
}
```

### <a name="step4"></a>Step 4: Check capabilities of a token

This step shows how to print out the permitted capabilities of a token on a
path. This can help verifying what operations are granted based on the policies
attached to the token.

#### CLI command

The command is:

```shell
$ vault token capabilities <TOKEN> <PATH>
```

**Example:**

First, create a token attached to `admin` policy.

```shell
$ vault token create -policy="admin"
Key                  Value
---                  -----
token                2sHGlAHNj36LpqQ2Zevl2Owi
token_accessor       4G4UIsQOMwifg7vMLqf6QIc3
token_duration       768h
token_renewable      true
token_policies       ["admin" "default"]
identity_policies    []
policies             ["admin" "default"]
```

Now, fetch the capabilities of this token on the `sys/auth/approle` path.

```plaintext
$ vault token capabilities 2sHGlAHNj36LpqQ2Zevl2Owi sys/auth/approle
create, delete, read, sudo, update
```

The result should match the policy rule you wrote on the `sys/auth/*` path. You
can repeat the steps to generate a token for `provisioner` and check its
capabilities on paths.

In the absence of a token, it returns the capabilities of the current token
invoking this command.

```shell
$ vault token capabilities sys/auth/approle
root
```

#### API call using cURL

Use the `sys/capabilities` endpoint.

**Example:**

First, create a token attached to the `admin` policy:

```shell
$ curl --request POST --header "X-Vault-Token: ..." --data '{ "policies":"admin" }' \
       http://127.0.0.1:8200/v1/auth/token/create
{
   "request_id": "bd9b3216-f7e6-610c-4861-38b9112a1821",
   "lease_id": "",
   "renewable": false,
   "lease_duration": 0,
   "data": null,
   "wrap_info": null,
   "warnings": null,
   "auth": {
     "client_token": "3xlduc1vGMD7vKeGLyONAxdS",
     "accessor": "FOoNv0YJSCqtPVCpW03qVeKd",
     "policies": [
       "admin",
       "default"
     ],
     "token_policies": [
       "admin",
       "default"
     ],
     "metadata": null,
     "lease_duration": 2764800,
     "renewable": true,
     "entity_id": ""
   }
}
```

Now, fetch the capabilities of this token on the `sys/auth/approle` path.

```shell
# Request payload
$ tee payload.json <<EOF
{
  "token": "3xlduc1vGMD7vKeGLyONAxdS",
  "path": "sys/auth/approle"
}
EOF

$ curl --request POST --header "X-Vault-Token: ..." \
       --data @payload.json \
       http://127.0.0.1:8200/v1/sys/capabilities | jq
{
  "sys/auth/approle": [
    "create",
    "delete",
    "read",
    "sudo",
    "update"
  ],
  "capabilities": [
    "create",
    "delete",
    "read",
    "sudo",
    "update"
  ],
  ...
}
```

The result should match the policy rule you wrote on the `sys/auth/*` path. You can
repeat the steps to generate a token for `provisioner` and check its
capabilities on paths.

To check the current token's capabilities permitted on a path, use
the `sys/capabilities-self` endpoint.

```plaintext
$ curl --request POST --header "X-Vault-Token: ..." \
       --data '{"path":"sys/auth/approle"}' \
       http://127.0.0.1:8200/v1/sys/capabilities-self
```


## Next steps

In this guide, you learned how to write policies based on given policy
requirements. Next, the [AppRole Pull Authentication](/guides/identity/authentication.html)
guide demonstrates how to associate policies to a role.

To learn about Sentinel policies, refer to the [Sentinel
Policies](/guides/identity/sentinel.html) guide.
