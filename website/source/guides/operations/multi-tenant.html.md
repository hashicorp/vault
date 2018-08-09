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


## Challenge

When Vault is primarily used as a central location to manage secrets, multiple
organizations within a company need to be able to manage their secrets in Vault
without interfering each other. This means that a company needs a ***Vault as a
Service*** model allowing each organization (tenant) to be restricted from
accessing, editing, or creating policies and secrets engines outside of their
tenant scope.

![Multi-Tenant](/assets/images/vault-multi-tenant.png)

## Solution

Create an ACL namespace dedicated to each team, organization, or app where they
can perform all necessary tasks within their tenant namespace.  

Each namespace has its own:

- Policies
- Mounts
- Tokens
- Identity entities and groups

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

# Write ACL policies
path "sys/policy/*" {
  capabilities = [ "create", "read", "update", "delete", "list" ]
}
```

If you are not familiar with policies, complete the
[policies](/guides/identity/policies.html) guide.


## Steps

**Scenario:**

In this guide, you are going to perform the following steps:

1. [Create ACL Namespaces](#step1)
1. Write Policies
1. Working in a namespace



### <a name="step1"></a>Step 1: Create ACL Namespaces


#### CLI command

To create a new namespace, run: **`vault namespace create <namespace_name>`**

Once a namespace is created, you can do one of the two things:

* Set **`VAULT_NAMESPACE`** so that all subsequent CLI commands will be
executed against that particular namespace

    ```plaintext
    $ export VAULT_NAMESPACE=<namespace_name>
    ```

* Specify the target namespace with **`-namespace`** flag

    ```plaintext
    $ vault policy write -namespace=<namespace_name> <policy_name> <policy_file>
    ```

<br>

Create a namespace dedicated to the **`finance`** team as well as **`education`** organizations:

```plaintext
$ vault namespace create finance

$ vault namespace create education
````

Now, create child namespaces called `training` and `certification` under the
`education` namespace:

```plaintext
$ vault namespace create -namespace=education training

$ vault namespace create -namespace=education certification
````

List the created namespaces:

```plaintext
$ vault namespace list
education/
finance/

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

1. Create a namespace dedicated to the **`finance`** team:

    ```plaintext
    $ curl --header "X-Vault-Token: ..." \
           --request POST \
           http://127.0.0.1:8200/v1/sys/namespaces/finance
    ```

1. Similarly, create a top-level namespace for the **`education`** organization:

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
           "finance/",
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

This is a policy for a namespace administrator which includes full permissions
under the path, as well as full permissions for namespaced policies.

#### Author a policy file

`finance-admins.hcl`

```shell
# Full permissions on the finance path
path "finance/*" {
   capabilities = ["create", "read", "update", "delete", "list", "sudo"]
}
```

#### CLI command

To create policies:

```shell
# Create finance-admins policy
$ vault policy write finance-admins finance-admins.hcl
```

#### API call using cURL

To create a policy, use `/sys/policy` endpoint:

```shell
# Create finance-admins policy
$ curl --request PUT --header "X-Vault-Token: ..." --data @payload.json \
    https://vault.rocks/v1/sys/policy/finance-admins

$ cat payload.json
{
  "policy": "path \"finance/*\" { capabilities = [\"create\", \"read\", \"update\", ... }"
}
```


## Next steps
