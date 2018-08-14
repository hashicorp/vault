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
- **`org-admin`** is the organization-level administrator
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
- Identity entities and groups

~> Tokens are locked to a namespace or child-namespaces.  Identity groups can
pull in entities and groups from _other_ namespaces.

## Prerequisites

To perform the tasks described in this guide, you need to have a **Vault
Enterprise** environment.  

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

# Manage entities and groups
path "identity/*" {
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
1. [Setup entities and groups](#step3)
1. [Test the organization admin user](#step4)


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


#### Web UI

1. Open a web browser and launch the Vault UI (e.g. http://127.0.01:8200/ui)
and then login.

1. Select **Access**.

1. Select **Namespaces** and then click **Add a namespace**.

1. Enter **`education`** in the **Path** field.

1. Click **Save**.

1. Select **Add a namespace** again, and then enter **`education/training`** in
the **Path** field.

1. Click **Save**.

1. Select **Add a namespace** again, and then enter
**`education/certification`** in the **Path** field.

1. Click **Save**.



### <a name="step2"></a>Step 2: Write Policies
(**Persona:** operations)

In this scenario, there is an organization-level administrator who is a super
user within the scope of the **`education`** namespace.  Also, there is a
team-level administrator for **`training`** and **`certification`**.


#### Policy for education admin

**Requirements:**

- Create and manage namespaces
- Create and manage policies
- Enable and manage secret engines
- Create and manage entities and groups
- Manage tokens

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

- Create and manage child-namespaces
- Create and manage policies
- Enable and manage secret engines

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

Now, let's deploy the policies!


#### CLI command

To target a specific namespace, you can do one of the following:

* Set **`VAULT_NAMESPACE`** so that all subsequent CLI commands will be
executed against that particular namespace

    ```plaintext
    $ export VAULT_NAMESPACE=<namespace_name>
    $ vault policy write <policy_name> <policy_file>
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

* Prepend the API endpoint with namespace name (e.g.
**`<namespace_name>`**`/sys/policies/acl`)


Create **`edu-admin`** and **`training-admin`** policies.

```shell
# Create a request payload
$ tee edu-payload.json <<EOF
{
  "policy": "path \"sys/namespaces/education/*\" {\n  capabilities = [\"create\", \"read\", \"update\", \"delete\", \"list\", \"sudo\"]\n } ... "
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
 "policy": "path \"sys/namespaces/education/training/*\" {\n  capabilities = [\"create\", \"read\", \"update\", \"delete\", \"list\", \"sudo\"]\n  }  ... "
}
EOF

# Create training-admin policy under 'education/training' namespace
# This example directs the target namespace in the API endpoint
$ curl --header "X-Vault-Token: ..." \
       --request PUT \
       --data @training-payload.json \
       https://vault.rocks/v1/education/training/sys/policies/acl/training-admin
```


#### Web UI

1. In the Web UI, select **education** for the **CURRENT NAMESPACE** in the
upper left menu.

1. Click the **Policies** tab, and then select **Create ACL policy**.

1. Toggle **Upload file**, and click **Choose a file** to select your
**`edu-admin.hcl`** file you authored.  This loads the policy and sets the
**Name** to be `edu-admin`.

1. Click **Create Policy** to complete.

1. Set the **CURRENT NAMESPACE** to be **education/training** in the upper left
menu.

1. In the **Policies** tab, select **Create ACL policy**.

1. Toggle **Upload file**, and click **Choose a file** to select your
**`training-admin.hcl`** file you authored.  

1. Click **Create Policy**.


### <a name="step3"></a>Step 3: Setup entities and groups
(**Persona:** operations)

In this step, you are going to create an entity, Bob Smith who is an
organization-level administrator.  Also, you are going to create a group for
team-level administrator, Team Admin, and add  Bob Smith as a member so that he
inherits the `training-admin` policy as well.

![Entities and Groups](/assets/images/vault-multi-tenant-3.png)

-> **NOTE:** If you are not familiar with entities and groups, refer to the
[Identity - Entities and
Groups](http://localhost:4567/guides/identity/identity.html) guide.

#### CLI Command

```shell
# First, you need to enable userpass auth method
$ vault auth enable -namespace=education userpass

# Create a user 'bob'
$ vault write -namespace=education \
        auth/userpass/users/bob password="password"

# Create an entity for Bob Smith with 'edu-admin' policy attached
# Save the generated entity ID in entity_id.txt file
$ vault write -namespace=education -format=json identity/entity name="Bob Smith" \
        policies="edu-admin" | jq -r ".data.id" > entity_id.txt

# Get the mount accessor for userpass auth method and save it in accessor.txt file
$ vault auth list -namespace=education -format=json \
        | jq -r '.["userpass/"].accessor' > accessor.txt

# Create an entity alias for Bob Smith to attach 'bob'
$ vault write -namespace=education identity/entity-alias name="bob" \
        canonical_id=$(cat entity_id.txt) mount_accessor=$(cat accessor.txt)

# Create a group, "team-admin" in education/training namespace
$ vault write -namespace=education/training identity/group name="Training Admin" \
        policies="training-admin" member_entity_ids=$(cat entity_id.txt)
```


#### API call using cURL

```shell
# Enable the `userpass` auth method
$ curl --header "X-Vault-Token: ..." \
       --header "X-Vault-Namespace: education" \
       --request POST \
       --data '{"type": "userpass"}' \
       http://127.0.0.1:8200/v1/sys/auth/userpass

# Create a user 'bob'
$ curl --header "X-Vault-Token: ..." \
       --header "X-Vault-Namespace: education" \
       --request POST \
       --data '{"password": "password"}' \
       http://127.0.0.1:8200/v1/auth/userpass/users/bob

# Create an entity for Bob Smith with 'edu-admin' policy attached
# Copy the generated entity ID
$ curl --header "X-Vault-Token: ..." \
       --header "X-Vault-Namespace: education" \
       --request POST \
       --data '{"name": "Bob Smith", "policies": "edu-admin"}' \
       http://127.0.0.1:8200/v1/identity/entity
{
   ...
   "data": {
     "aliases": null,
     "id": "6ded4d31-481f-040b-11ad-c6db0cb4d211"
   },
   ...
}

# Get the mount accessor for userpass auth method
$ curl --header "X-Vault-Token: ..." \
       http://127.0.0.1:8200/v1/sys/auth | jq
{
 ...
 "userpass/": {
   "accessor": "auth_userpass_9b6cd254",
  ...
}

# Create the API request message payload
$ tee payload-bob.json <<EOF
{
  "name": "bob",
  "canonical_id": "6ded4d31-481f-040b-11ad-c6db0cb4d211",
  "mount_accessor": "auth_userpass_9b6cd254"
}
EOF

# Create an entity alias for Bob Smith to attach 'bob'
$ curl --header "X-Vault-Token: ..." \
       --header "X-Vault-Namespace: education" \
       --request POST \
       --data @payload-bob.json \
       http://127.0.0.1:8200/v1/identity/entity-alias

# Create a group, "team-admin" in education/training namespace
# API request msg payload.  Be sure to enter the correct Bob Smith entity ID
$ tee payload-group.json <<EOF
{
  "name": "Training Admin",
  "policies": "training-admin",
  "member_entity_ids": "6ded4d31-481f-040b-11ad-c6db0cb4d211"
}
EOF

# Use identity/group endpoint
$ curl --header "X-Vault-Token: ..." \
       --header "X-Vault-Namespace: education/training" \
       --request PUT \
       --data @payload-group.json \
       http://127.0.0.1:8200/v1/identity/group | jq
{
   ...
   "data": {
     "id": "d62157aa-b5f6-b6fe-aa40-0ffc54defc41",
     "name": "Training Admin"
   },
   ...
}
```


#### Web UI

1. In the Web UI, select **education** for the **CURRENT NAMESPACE** in the
upper left menu.

1. Click the **Access** tab, and select **Enable new method**.

1. Select **Username & Password** from the **Type** drop-down menu.

1. Click **Enable Method**.

1. Click the Vault CLI shell icon (**`>_`**) to open a command shell.  Enter the
following command to create a new user, **`bob`**:

    ```plaintext
    $ vault write education/auth/userpass/users/bob password="password"
    ```
    ![Create Policy](/assets/images/vault-multi-tenant-4.png)

1. Click the icon (**`>_`**) again to hide the shell.

1. From the **Access** tab, select **Entities** and then **Create entity**.

1. Enter **`Bob Smith`** in the **Name** field, and **`edu-admin`** in the
**Policies** field.   

1. Click **Create**.

1. Select **Add alias**.  Enter **`bob`** in the **Name** field and select
**`userpass/ (userpass)`** from the **Auth Backend** drop-down list.

1. Click **Create**.

1. Click the **Access** tab and select **Entities**.

1. Select the **`bob-smith`** entity and copy its **ID** displayed under the
**Details** tab.

1. Now, click **Groups** from the left navigation, and select **Create group**.

1. Enter **`Training Admin`** in the **Name** field, **`training-admin`** in the
**Policies** field, and finally paste in the entity ID in the **Member Entity
IDs** field.

1. Click **Create**.


### <a name="step4"></a>Step 4: Test the organization admin user
(**Persona:** org-admin)

#### CLI Command

Log in as **`bob`** into the `education` namespace:

```shell
# Login as 'bob'
$ vault login -namespace=education -method=userpass username="bob" password="password"

Key                    Value
---                    -----
token                  52e6fb8e-7f1f-8cf4-47c5-fff2c932e2ee
token_accessor         a7d01e20-1bac-98b9-a40d-349ad1868e31
token_duration         768h
token_renewable        true
token_policies         ["default"]
identity_policies      ["edu-admin" "training-admin"]
policies               ["default" "edu-admin" "training-admin"]
token_meta_username    bob

# Set the target namespace as an env variable
$ export VAULT_NAMESPACE="education"

# Create a new namespace called 'web-app'
$ vault namespace create web-app
Success! Namespace created at: education/web-app/

# Enable key/value v2 secrets engine at edu-secret
$ vault secrets enable -path=edu-secret kv-v2
Success! Enabled the kv-v2 secrets engine at: edu-secret/
```

Optionally, you can create new policies to test that `bob` can perform the
operations as expected.  When you are done testing, unset the VAULT_NAMESPACE
environment variable.  

```plaintext
$ unset VAULT_NAMESPACE
```

#### API call using cURL

Log in as **`bob`** into the `education` namespace:

```shell
# Log in as bob
$ curl --header "X-Vault-Namespace: education" \
       --request POST \
       --data '{"password": "password"}' \
       http://127.0.0.1:8200/v1/auth/userpass/login/bob | jq
{
   ...
   "auth": {
     "client_token": "e639b1ec-a2df-9645-169a-fb50d1a19c01",
     "accessor": "a462f0f9-dd3d-5f52-de8a-180f23d5cb05",
     "policies": [
       "default",
       "edu-admin",
       "training-admin"
     ],
     "token_policies": [
       "default"
     ],
     "identity_policies": [
       "edu-admin",
       "training-admin"
     ],
     "metadata": {
       "username": "bob"
     },
     ...
   }
}

# Create a new namespace called 'web-app'
# Be sure to use generated bob's client token
$ curl --header "X-Vault-Token: e639b1ec-a2df-9645-169a-fb50d1a19c01" \
       --request POST \
       http://127.0.0.1:8200/v1/education/sys/namespaces/web-app

# Enable key/value v2 secrets engine at edu-secret
$ curl --header "X-Vault-Token: e639b1ec-a2df-9645-169a-fb50d1a19c01" \
       --request POST \
       --data '{"type": "kv-v2"}' \
       http://127.0.0.1:8200/v1/education/sys/mounts/edu-secret
```


#### Web UI

1. Open a web browser and launch the Vault UI (e.g. http://127.0.01:8200/ui). If
you are already logged in, log out.

1. At the **Sign in to Vault**, set the **Namespace** to **`education`**.

1. Select the **Userpass** tab, and enter **`bob`** in the **Username** field,
and **`password`** in the **Password** field.

1. Click **Sign in**.  Notice that the CURRENT NAMESPACE is set to **education**
in the upper left corner of the UI.

1. To add a new namespace, select **Access**.

1. Select **Namespaces** and then click **Add a namespace**.

1. Enter **`web-app`** in the **Path** field, and then click **Save**.

1. Select **Secrets**.

1. Select **Enable new engine**.

1. Select **KV** from the **Secrets engine type** drop-down list, and enter
**`edu-secret`** in the **Path** field.

1. Click **Enable Engine** to finish.


## Next steps

Refer to the [Sentinel Policies](/guides/identity/sentinel.html) guide if you
need to write policies that allow you to embed finer control over the user
access across those namespaces. 
