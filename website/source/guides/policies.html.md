---
layout: "guides"
page_title: "Policies - Guides"
sidebar_current: "guides-governance"
description: |-
  Policies in Vault control what a user can access.
---

# Policies

Use policy to govern the behavior of clients by specifying the access privilege
(_authorization_).

When you first initialize Vault, the
[**`root`**](/docs/concepts/policies.html#root-policy) policy gets created by
default. The `root` policy is a special policy that gives superuser access to
_everything_ in Vault. This allows the superuser to set up initial policies,
tokens, etc.

In addition, there is another build-in policy,
[**`default`**](/docs/concepts/policies.html#default-policy) gets created. The
`default` policy is attached to all tokens and provides common permissions.

Everything in Vault is path based, and write policies to grant or forbid access
to certain paths and operations in Vault. Empty policy grants **no permission**
in the system.


### HashiCorp Configuration Language (HCL)

Policies written in [HCL](https://github.com/hashicorp/hcl) format are often
referred as **_ACL Policy_**. [Sentinel](https://www.hashicorp.com/sentinel) is
another framework for policy which is available in [Vault
Enterprise](/docs/enterprise/index.html).  Since Sentinel is an enterprise-only
feature, this guide focuses on writing ACL policies.

**NOTE:** HCL is JSON compatible; therefore, JSON can be used as completely
valid input.

## Reference Material

- [Policies](/docs/concepts/policies.html#default-policy) documentation
- [Policy API](/api/system/policy.html) documentation
- [Getting Started guide](/intro/getting-started/policies.html) on policies

## Estimated Time to Complete

10 minutes

## Challenge

Since Vault centrally secure, store, and access control secrets across distributed
infrastructure and applications, it is critical to control permissions before
any user or machine can gain access.


## Solution

Restrict the use of root policy, and write fine-grained policies to practice
**least privileged**. For example, if an app gets AWS credentials from Vault,
write policy grants to `read` from AWS secret backend but not to `delete`, etc.

Policies are attached to tokens and roles to enforce client permissions on
Vault.


## Prerequisites

To perform the tasks described in this guide, you need to have a Vault
environment.  Refer to the [Getting
Started](/intro/getting-started/install.html) guide to install Vault.

Make sure that your Vault server has been [initialized and
unsealed](/intro/getting-started/deploy.html).

## Steps

This guide demonstrates basic policy authoring and management tasks.

1. [Write ACL policies](#step1)
2. [View existing policies](#step2)
3. [Check capabilities of a token](#step3)


### <a name="step1"></a>Step 1: Write ACL policies

ACL policies are defined for each path:

```plaintext
path "<PAHT>" {
  capabilities = [ "<LIST_OF_CAPABILITIES>" ]
}
```

-> The path can have a wildcard ("`*`") specifying at the end to allow for
namespacing. For example, "`secret/training_*`" grants permissions on any
path starts with "`secret/training_`" (e.g. `secret/training_vault`).

Define one or more capabilities on each path to control operations that are
permitted.

| Capability      | Associated HTTP verbs  |
| --------------- |------------------------|
| create          | POST/PUT               |
| read            | GET                    |
| update          | POST/PUT               |
| delete          | DELETE                 |
| list            | LIST


#### Policy requirements

First step in writing policies is to gather policy requirements. As an exercise,
assume that the following are the policy requirements defined for `developers`
and `devops`.

Developers must be able to:

- Perform all operations on `secret/apps/*` path
- Obtain AWS credentials for `devops` role
- Obtain database credentials for `devops` role

DevOps must be able to:

- View and manage leases
- Mount and manage secret backends

#### Example policy for Developers

`dev-pol.hcl`

```shell
# Full permissions on secret/apps/* path
path "secret/apps/*" {
  capabilities = ["create", "read", "update", "delete", "list"]
}

# Read new credentials from aws secret backend
path "aws/creds/devops" {
  capabilities = ["read"]
}

# Read new credentials from database secret backend
path "database/creds/devops" {
  capabilities = ["read"]  
}
```

#### Example policy for DevOps

`devops-pol.hcl`

```shell
# Permissions to manage leases
path "sys/leases/*" {
  capabilities = ["create", "read", "update", "delete", "list"]
}

# Permissions to create and manage secret backends
path "sys/mounts/*" {
  capabilities = ["create", "read", "update", "delete", "list"]
}

# Move already mounted backend to a new endpoint
path "sys/remount/*" {
  capabilities = ["create", "update"]
}
```

Once the ACL policies are written, create policies.

#### CLI command

Create new ACL policies:

```shell
$ vault write dev-pol dev-pol.hcl

$ vault write devops-pol devops-pol.hcl
```

**NOTE:** To update an existing policy, simply re-run the same command by
passing your modified policy (`*.hcl`).

#### API call using cURL

Before begin, create the following environment variables for your convenience:

- **VAULT_ADDR** is set to your Vault server address
- **VAULT_TOKEN** is set to your Vault token

**Example:**

```plaintext
$ export VAULT_ADDR=http://127.0.0.1:8201

$ export VAULT_TOKEN=0c4d13ba-9f5b-475e-faf2-8f39b28263a5
```

Now, create new ACL policies using API:

```shell
# Create dev-pol policy
$ curl -X PUT -H "X-Vault-Token: $VAULT_TOKEN" -d @dev-payload.json \
    $VAULT_ADDR/v1/sys/policies/acl/dev-pol

$ cat dev-payload.json
{
  "policy": "path \"secret/apps/*\" { capabilities = [\"create\", \"read\", \"update\" ... }"
}

# Create devops-pol policy
$ curl -X PUT -H "X-Vault-Token: $VAULT_TOKEN" -d @devops-payload.json \
    $VAULT_ADDR/v1/sys/policies/acl/devops-pol

$ cat devops-payload.json
{
  "policy": "path \"sys/leases/*\" { capabilities = [\"create\", \"read\", \"update\" ... }"
}
```

**NOTE:** To update an existing policy, simply re-run the same command by
passing your modified policy in the request payload (`*.json`).



### <a name="step2"></a>Step 2: View existing policies

Make sure that you see the policies you created in [Step 1](#step1).

#### CLI command

The following command lists existing policies:

```shell
vault policies
```

To view a specific policy:

```shell
vault read sys/policy/<POLICY_NAME>
```

**Example:**

```shell
vault read sys/policy/default

Key  	Value
---  	-----
name 	default
rules	# Allow tokens to look up their own properties
path "auth/token/lookup-self" {
    capabilities = ["read"]
}

# Allow tokens to renew themselves
path "auth/token/renew-self" {
    capabilities = ["update"]
}

# Allow tokens to revoke themselves
path "auth/token/revoke-self" {
    capabilities = ["update"]
}
...
```


#### API call using cURL

To list existing ACL policies, use the `/sys/policy` endpoint.

```plaintext
curl -X LIST -H "X-Vault-Token: $VAULT_TOKEN" $VAULT_ADDR/v1/sys/policy | jq
```

-> NOTE: You can also use `/sys/policies` endpoint which is used to manage
ACL, RGP, and EGP policies in Vault (RGP and EGP policies are enterprise-only
features). To list policies, invoke `/sys/policies/acl` endpoint.

To read a specific policy, the endpoint path should be
`/sys/policy/<POLICY_NAME>`.

**Example:**

```plaintext
curl -X GET -H "X-Vault-Token: $VAULT_TOKEN" $VAULT_ADDR/v1/sys/policy/default | jq

{
  "name": "default",
  "rules": "\n# Allow tokens to look up their own properties\npath \"auth/token/lookup-self\" ...",
  "request_id": "1379d18d-e5e3-dbd2-8f3a-0e446c016a23",
  "lease_id": "",
  "renewable": false,
  "lease_duration": 0,
  ...
```


### <a name="step3"></a>Step 3: Check capabilities of a token

Use the `/sys/capabilities` endpoint to fetch the capabilities of a token on a
given path. This helps to verify what operations are granted based on the
policies attached to the token.

#### CLI command

The command is:

```shell
vault capabilities <TOKEN> <PATH>
```

**Example:**

```plaintext
vault capabilities a59c0d41-8df7-ba8e-477e-9bfb394f28a0 secret/apps

Capabilities: [create delete list read update]
```

In the absence of token, it returns capabilities of current token invoking this
command.

```shell
vault capabilities secret/apps
```

#### API call using cURL

Use the `sys/capabilities` endpoint.

**Example:**

```plaintext
$ curl -X POST -H "X-Vault-Token: $VAULT_TOKEN" -d @payload.json \
    $VAULT_ADDR/v1/sys/capabilities

$ cat payload.json
{
  "token": "a59c0d41-8df7-ba8e-477e-9bfb394f28a0",
  "path": "secret/apps"
}
```

To check current token's capabilities permitted on a path, use
`sys/capabilities-self` endpoint.

```plaintext
curl -X POST -H "X-Vault-Token: $VAULT_TOKEN" -d '{"path":"secret/apps"}' \
    $VAULT_ADDR/v1/sys/capabilities-self
```



## Next steps

In this guide, you learned how to write policies based on given policy
requirements. Next, [AppRole Pull Authentication](/guides/authentication.html)
guide demonstrates how to associate policies to a role.
