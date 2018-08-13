---
layout: "guides"
page_title: "Multi-Tenant Pattern - Guides"
sidebar_current: "guides-operations-multi-tenant"
description: |-
  This guide provides guidance in creating a multi-tenant environment.
---

# Multi-Tenant Pattern with ACL Namespaces

~> **Enterprise Only:** ACL Namespace feature is a part of _Vault Enterprise Pro_.

Everything in Vault is path-based, and often use the terms `path` and
`namespace` interchangeably. The application namespace pattern is a useful
construct for providing Vault as a service to internal customers, giving them
the ability to leverage a multi-tenant Vault implementation with full agency to
their application's interactions with Vault.


## Reference Material

- [Vault Deployment Reference Architecture](/guides/operations/reference-architecture.html)
- [Policies](/guides/identity/policies.html) guide


## Estimated Time to Complete

10 minutes


## Personas

The scenario described in this guide introduces the following personas:

- **`operations`** is the cluster-level administrator with privileged policies
- **`edu-admin`** is the organization-level administrator
- **`team-admin`** is the team-level administrator

## Challenge

When Vault is primarily used as a central location to manage secrets, multiple
organizations within a company need to be able to manage their secrets in
self-serving manner. This means that a company needs to implement a ***Vault as
a Service*** model allowing each organization (tenant) to manage their own
secrets and policies. The most importantly, tenants should be restricted to work
only within their tenant scope.

![Multi-Tenant](/assets/images/vault-multi-tenant.png)

## Solution

Create an **ACL namespace** dedicated to each team, organization, or app where
they can perform all necessary tasks within their tenant namespace.  

Each namespace can have its own:

- Policies
- Mounts
- Tokens

> Tokens are locked to a namespace or child-namespaces.  Identity groups can
pull in entities and groups from other namespaces.

## Prerequisites

To perform the tasks described in this guide, you need to have a Vault
Enterprise environment.  

### <a name="policy"></a>Policy requirements

-> **NOTE:** For the purpose of this guide, you can use **`root`** token to work
with Vault. However, it is recommended that root tokens are only used for just
enough initial setup or in emergencies. As a best practice, use tokens with
appropriate set of policies based on your role in the organization.

To perform all tasks demonstrated in this guide, your policy must include the
following permissions:

```shell
# Create ACL namespaces
path "sys/namespaces/*" {
  capabilities = [ "create", "read", "update", "delete", "list", "sudo" ]
}

# Write policies
path "sys/policies/*" {
  capabilities = [ "create", "read", "update", "delete", "list" ]
}
```

If you are not familiar with policies, complete the
[policies](/guides/identity/policies.html) guide.


## Steps

**Scenario:** In this guide, you are going to create a namespace dedicated to
_Education_ organization which has _Training_ and _Certification_ teams.
Delegate policy management to the team admins so that the cluster operator
won't have to be involved.

![Scenario](/assets/images/vault-multi-tenant-2.png)

In this guide, you are going to perform the following steps:

1. [Create ACL Namespaces](#step1)
1. [Write Policies](#step2)
1. [Setup users](#step3)
1. Working in a namespace


### <a name="step1"></a>Step 1: Create ACL Namespaces
(**Persona:** operations)

#### CLI command

To create a new namespace, run: **`vault namespace create <namespace_name>`**

1. Create a namespace dedicated to the **`education`** organizations:

    ```plaintext
    $ vault namespace create education
    ````

1. Create child namespaces called `training` and `certification` under the
`education` namespace:

    ```plaintext
    $ vault namespace create -namespace=education training

    $ vault namespace create -namespace=education certification
    ````

1. List the created namespaces:

    ```plaintext
    $ vault namespace list
    education/

    $ vault namespace list -namespace=education
    certification/
    training/
    ```

#### API call using cURL

To create a new namespace, invoke **`sys/namespaces`** endpoint:

```shell
$ curl --header "X-Vault-Token: <TOKEN>" \
       --request POST \
       <VAULT_ADDRESS>/v1/sys/namespaces/<NS_NAME>
```

Where `<TOKEN>` is your valid token, and `<NS_NAME>` is the desired namespace
name.

1. Create a namespace for the **`education`** organization:

    ```plaintext
    $ curl --header "X-Vault-Token: ..." \
           --request POST \
           http://127.0.0.1:8200/v1/sys/namespaces/education
    ```

1. Now, create a child namespace called **`training`** and **`certification`**
under `education`. To do so, pass the top-level namespace name in the
**`X-Vault-Namespace`** header.

    ```shell
    # Create a training namespace under education
    # NOTE: Top-level namespace is in the API endpoint
    $ curl --header "X-Vault-Token: ..." \
           --header "X-Vault-Namespace: education" \
           --request POST \
           http://127.0.0.1:8200/v1/education/sys/namespaces/training

    # Create a certification namespace under education
    # NOTE: Pass the top-level namespace in the header
    $ curl --header "X-Vault-Token: ..." \
           --header "X-Vault-Namespace: education" \
           --request POST \
           http://127.0.0.1:8200/v1/sys/namespaces/certification
    ```

1. List the created namespaces:

    ```plaintext
    $ curl --header "X-Vault-Token: ..." \
           --request LIST
           http://127.0.0.1:8200/v1/sys/namespaces | jq
    {
       ...
       "data": {
         "keys": [
           "education/"
         ]
       },
       ...


    $ curl --header "X-Vault-Token: ..." \
           --request LIST
           http://127.0.0.1:8200/v1/education/sys/namespaces | jq
     {
       ...
       "data": {
         "keys": [
           "certification/",
           "training/"
         ]
       },
       ...
    ```


### <a name="step2"></a>Step 2: Write Policies
(**Persona:** operations)

In this scenario, there is an organization-level administrator who is a super
user within the scope of the **`education`** namespace.  Also, there is a
team-level administrator for **`training`** and **`certification`**.


#### Policy for education admin

**Requirements:**

- Create and manage namespaces for `education` namespace and its child namespaces
- Create and manage policies for `education` namespace and its child namespaces
- Enable and manage secret engines for `education` namespace and its child namespaces
- Create and manage entities and groups for `education` namespace and its child namespaces
- Manage tokens for `education` namespace and its child namespaces

**`edu-admin.hcl`**

```shell
# Manage namespaces
path "sys/namespaces/*" {
   capabilities = ["create", "read", "update", "delete", "list", "sudo"]
}

# Manage policies via API
path "sys/policies/*" {
   capabilities = ["create", "read", "update", "delete", "list", "sudo"]
}

# Manage policies via CLI
path "sys/policy/*" {
   capabilities = ["create", "read", "update", "delete", "list", "sudo"]
}

# List policies via CLI
path "sys/policy" {
   capabilities = ["read", "update", "list"]
}

# Enable and manage secrets engines
path "sys/mounts/*" {
   capabilities = ["create", "read", "update", "delete", "list"]
}

# List available secret engines
path "sys/mounts" {
  capabilities = [ "read" ]
}

# Create and manage entities and groups
path "identity/*" {
   capabilities = ["create", "read", "update", "delete", "list"]
}

# Manage tokens
path "auth/token/*" {
   capabilities = ["create", "read", "update", "delete", "list", "sudo"]
}
```


#### Policy for training admin

**Requirements:**


- Create and manage child-namespaces under `education/training` namespace
- Create and manage policies under `education/training` namespace
- Enable and manage secret engines under `education/training` namespace

**`training-admin.hcl`**

```shell
# Manage namespaces
path "sys/namespaces/*" {
   capabilities = ["create", "read", "update", "delete", "list", "sudo"]
}

# Manage policies via API
path "sys/policies/*" {
   capabilities = ["create", "read", "update", "delete", "list", "sudo"]
}

# Manage policies via CLI
path "sys/policy/*" {
  capabilities = ["create", "read", "update", "delete", "list", "sudo"]
}

# List policies via CLI
path "sys/policy" {
  capabilities = ["read", "update", "list"]
}


# Enable and manage secrets engines
path "sys/mounts/*" {
   capabilities = ["create", "read", "update", "delete", "list"]
}

# List available secret engines
path "sys/mounts" {
  capabilities = [ "read" ]
}
```


#### CLI command

To target a specific namespace, you can do one of the following:

* Set **`VAULT_NAMESPACE`** so that all subsequent CLI commands will be
executed against that particular namespace

    ```plaintext
    $ export VAULT_NAMESPACE=<namespace_name>
    ```

* Specify the target namespace with **`-namespace`** flag

    ```plaintext
    $ vault policy write -namespace=<namespace_name> <policy_name> <policy_file>
    ```

Since you have to deploy policies onto "`education`" and "`education/training`"
namespaces, use "`-namespace`" flag instead of environment variable.

Create **`edu-admin`** and **`training-admin`** policies.

```shell
# Create edu-admin policy under 'education' namespace
$ vault policy write -namespace=education edu-admin edu-admin.hcl

# Create training-admin policy under 'education/training' namespace
$ vault policy write -namespace=education/training training-admin training-admin.hcl
```

#### API call using cURL

To target a specific namespace, you can do one of the following:

* Pass the target namespace in the **`X-Vault-Namespace`** header

* Prepend the API endpoint with namespace name (e.g. `<namespace_name>/sys/policies/acl`)


Create **`edu-admin`** and **`training-admin`** policies.

```shell
# Create a request payload
$ tee edu-payload.json <<EOF
{
  "policy": "path \"sys/namespaces/education/*\" { capabilities = [\"create\", \"read\", \"update\", \"delete\", \"list\", \"sudo\"] } ... "
}
EOF

# Create edu-admin policy under 'education' namespace
$ curl --header "X-Vault-Token: ..." \
       --header "X-Vault-Namespace: education" \
       --request PUT \
       --data @edu-payload.json \
       https://vault.rocks/v1/sys/policies/acl/edu-admin

# Create a request payload
$ tee training-payload.json <<EOF
{
 "policy": "path \"sys/namespaces/education/training/*\" { capabilities = [\"create\", \"read\", \"update\", \"delete\", \"list\", \"sudo\"] }  ... "
}
EOF

# Create training-admin policy under 'education/training' namespace
# Notice that the target namespace is in the API endpoint on this example
$ curl --header "X-Vault-Token: ..." \
       --request PUT \
       --data @training-payload.json \
       https://vault.rocks/v1/education/training/sys/policies/acl/training-admin
```


### <a name="step3"></a>Step 3: Setup entities and groups
(**Persona:** operations)

Create organization-level admin user, `bob`, and team-level admin user, `jenn`.

-> **NOTE:** For the purpose of this guide, you are going to use the username &
password (`userpass`) auth method.

#### CLI Command

```shell
# First, you need to enable auth method
$ vault auth enable -namespace=education userpass

# Create a user 'bob' with both 'edu-admin' & 'training-admin' policies attached
$ vault write -namespace=education \
        auth/userpass/users/bob password="training" \
         policies="edu-admin, training-admin"

# First, you need to enable auth method
$ vault auth enable -namespace=training userpass

# Create a user 'jenn' with 'training-admin' policy attached
$ vault write -namespace=education/training \
        auth/userpass/users/jenn password="training" \
        policies="training-admin"
```


#### API call using cURL

```shell
# Enable the `userpass` auth method
$ curl --header "X-Vault-Token: ..." \
       --header "X-Vault-Namespace: education" \
       --request POST \
       --data '{"type": "userpass"}' \
       http://127.0.0.1:8200/v1/sys/auth/userpass

# Create a user 'bob' with both 'edu-admin' & 'training-admin' policies attached
$ curl --header "X-Vault-Token: ..." \
       --header "X-Vault-Namespace: education" \
       --request POST \
       --data '{"password": "training", "policies": ["edu-admin", "training-admin"]}' \
       http://127.0.0.1:8200/v1/auth/userpass/users/bob

# Enable the `userpass` auth method
$ curl --header "X-Vault-Token: ..." \
      --header "X-Vault-Namespace: education/training" \
      --request POST \
      --data '{"type": "userpass"}' \
      http://127.0.0.1:8200/v1/sys/auth/userpass

# Create a user 'jenn' with 'training-admin' policy attached
$ curl --header "X-Vault-Token: ..." \
      --header "X-Vault-Namespace: education/training" \
      --request POST \
      --data '{"password": "training", "policies": "training-admin"}' \
      http://127.0.0.1:8200/v1/auth/userpass/users/jenn
```

















### <a name="step4"></a>Step 4: Test organization admin user
(**Persona:** edu-admin)

#### CLI Command

Log in as **`bob`** into the `education` namespace:

```plaintext
$ vault login -namespace=education -method=userpass username="bob" password="training"

Key                    Value
---                    -----
token                  9ed2c0bf-65e2-a713-f0fa-2650b034cacf
token_accessor         d827c99c-35aa-64a1-78e8-7e3a4fd9741f
token_duration         768h
token_renewable        true
token_policies         ["default" "edu-admin" "training-admin"]
identity_policies      []
policies               ["default" "edu-admin" "training-admin"]
token_meta_username    bob
```


#### API call using cURL

Log in as **`bob`** into the `education` namespace:

```plaintext
$ curl --header "X-Vault-Namespace: education" \
       --request POST \
       --data '{"password": "training"}' \
       http://127.0.0.1:8200/v1/auth/userpass/login/bob | jq
{
  ...
  "auth": {
    "client_token": "e2ddcc38-4955-1b81-e9db-0d1091162172",
    "accessor": "eaff9242-fac6-f26f-729e-14822acd98f7",
    "policies": [
      "default",
      "edu-admin",
      "training-admin"
    ],
    "token_policies": [
      "default",
      "edu-admin",
      "training-admin"
    ],
    "metadata": {
      "username": "bob"
    },
    ...
}
```









### <a name="step5"></a>Step 5: Test training admin user
(**Persona:** training-admin)

#### CLI Command

Log in as **`jenn`** into the `education/training` namespace:

```plaintext
vault login -namespace=education/training -method=userpass username="jenn" password="training"

Key                    Value
---                    -----
token                  1187670c-aed7-1e5c-3e4e-079a607df5c6
token_accessor         c10acc61-c737-f0d4-9c97-34c83a4fa1f7
token_duration         768h
token_renewable        true
token_policies         ["default" "training-admin"]
identity_policies      []
policies               ["default" "training-admin"]
token_meta_username    jenn
```


#### API call using cURL

Log in as **`jenn`** into the `education/training` namespace:

```plaintext
$ curl --header "X-Vault-Namespace: education/training" \
       --request POST \
       --data '{"password": "training"}' \
       http://127.0.0.1:8200/v1/auth/userpass/login/jenn | jq
{
  ...
  "auth": {
    "client_token": "1187670c-aed7-1e5c-3e4e-079a607df5c6",
    "accessor": "c10acc61-c737-f0d4-9c97-34c83a4fa1f7",
    "policies": [
      "default",
      "training-admin"
    ],
    "token_policies": [
      "default",
      "training-admin"
    ],
    "metadata": {
      "username": "jenn"
    },
    ...
}
```




## Next steps
