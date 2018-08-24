---
layout: "guides"
page_title: "Secure Multi-Tenancy with Namepaces - Guides"
sidebar_title: "Multi-Tenant: Namespaces"
sidebar_current: "guides-operations-multi-tenant"
description: |-
  This guide provides guidance in creating a multi-tenant environment.
---

# Secure Multi-Tenancy with Namespaces

~> **Enterprise Only:** The namespaces feature is a part of _Vault Enterprise Pro_.

Everything in Vault is path-based, and often uses the terms `path` and
`namespace` interchangeably. The application namespace pattern is a useful
construct for providing Vault as a service to internal customers, giving them
the ability to implement secure multi-tenancy within Vault in order to provide
isolation and ensure teams can self-manage their own environments.


## Reference Material

- [Namespaces](/docs/enterprise/namespaces/index.html)
- [Streamline Secrets Management with Vault Agent and Vault 0.11](https://youtu.be/zDnIqSB4tyA)
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
organizations within a company may need to be able to manage their secrets in a
self-serving manner. This means that a company needs to implement a ***Vault as
a Service*** model allowing each organization (tenant) to manage their own
secrets and policies. The most importantly, tenants should be restricted to work
only within their tenant scope.

![Multi-Tenant](/img/vault-multi-tenant.png)

## Solution

Create a **namespace** dedicated to each team, organization, or app where
they can perform all necessary tasks within their tenant namespace.  

Each namespace can have its own:

- Policies
- Auth Methods
- Secret Engines
- Tokens
- Identity entities and groups

~> Tokens are locked to a namespace or child-namespaces.  Identity groups can
pull in entities and groups from _other_ namespaces.

## Prerequisites

To perform the tasks described in this guide, you need to have a **Vault
Enterprise** environment.  

-> **NOTE:** The creation of namespaces should be performed by a user with a
highly privileged token such as **`root`** to set up isolated environments for
each organization, team, or application.


## Steps

**Scenario:** In this guide, you are going to create a namespace dedicated to
the Education organization which has Training and Certification teams. Delegate
operational tasks to the team admins so that the Vault cluster operators won't
have to be involved.

![Scenario](/img/vault-multi-tenant-2.png)

In this guide, you are going to perform the following steps:

1. [Create namespaces](#step1)
1. [Write policies](#step2)
1. [Setup entities and groups](#step3)
1. [Test the organization admin user](#step4)
1. [Test the team admin user](#step5)
1. [Audit ambient credentials](#step6)


### <a name="step1"></a>Step 1: Create namespaces
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

1. Select **Namespaces** and then click **Create a namespace**.

1. Enter **`education`** in the **Path** field.

1. Click **Save**.

1. To create child namespaces, select the down-arrow on the upper left corner of
the UI, and select **education** under **CURRENT NAMESPACE**.

    ![NS Selection](/img/vault-multi-tenant-1.png)

1. Under the **Access** tab, select **Namespaces** and then click **Create a namespace**.

1. Enter **`training`** in the **Path** field, and click **Save**.

1. Select **Create a namespace** again, and then enter **`certification`** in
the **Path** field, and click **Save**.



### <a name="step2"></a>Step 2: Write Policies
(**Persona:** operations)

In this scenario, there is an organization-level administrator who is a
superuser within the scope of the **`education`** namespace.  Also, there is a
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

-> Also, refer to the [Additional Discussion](#policy-with-namespaces) section
to learn more about policy authoring with namespaces.


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
       https://127.0.0.1:8200/v1/sys/policies/acl/edu-admin

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
       https://127.0.0.1:8200/v1/education/training/sys/policies/acl/training-admin
```


#### Web UI

1. In the Web UI, make sure that the **CURRENT NAMESPACE** is set to
**education** in the upper left menu.

1. Click the **Policies** tab, and then select **Create ACL policy**.

1. Toggle **Upload file** sliding switch, and click **Choose a file** to select your
**`edu-admin.hcl`** file you authored.  This loads the policy and sets the
**Name** to be `edu-admin`.

1. Click **Create Policy** to complete.

1. Set the **CURRENT NAMESPACE** to be **education/training** in the upper left
menu.
    ![Namespace](/img/vault-multi-tenant-6.png)

1. In the **Policies** tab, select **Create ACL policy**.

1. Toggle **Upload file**, and click **Choose a file** to select your
**`training-admin.hcl`** file you authored.  

1. Click **Create Policy**.


### <a name="step3"></a>Step 3: Setup entities and groups
(**Persona:** operations)

Bob who is an organization-level administrator (superuser) has two accounts:
**`bob`** and **`bsmith`**. You will create an entity, **Bob Smith** to
associate those two accounts.

Also, you are going to create a group for the team-level administrator, **Team
Admin**, and add Bob Smith entity as a group member so that Bob can inherit the
`training-admin` policy to manage the child namespace if he ever has to take
over.

![Entities and Groups](/img/vault-multi-tenant-3.png)

-> This step only demonstrates CLI commands and Web UI to create
entities and groups.  Refer to the [Identity - Entities and
Groups](/guides/identity/identity.html) guide if you need the full details.
Also, read the [Additional Discussion](#additional-discussion) section for
an example of setting up external groups.

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

# Create a group, "Training Admin" in education/training namespace
$ vault write -namespace=education/training identity/group \
        name="Training Admin" policies="training-admin" \
        member_entity_ids=$(cat entity_id.txt)

# Enable userpass auth method in training namespace
$ vault auth enable -namespace=education/training userpass

# Create a user 'bsmith'
$ vault write -namespace=education/training \
        auth/userpass/users/bsmith password="password"

# Get the mount accessor for userpass auth method and save it in accessor2.txt file
$ vault auth list -namespace=education/training -format=json \
        | jq -r '.["userpass/"].accessor' > accessor2.txt

# Add 'bsmith' to Bob Smith entity as its alias
$ vault write -namespace=education identity/entity-alias name="bsmith" \
        canonical_id=$(cat entity_id.txt) mount_accessor=$(cat accessor2.txt)        
```


#### Web UI

1. In the Web UI, make sure that the **CURRENT NAMESPACE** is set to
**education** in the upper left menu.

1. Click the **Access** tab, and select **Enable new method**.

1. Select **Username & Password** from the **Type** drop-down menu.

1. Click **Enable Method**.

1. Click the Vault CLI shell icon (**`>_`**) to open a command shell.  Enter the
following command to create a new user, **`bob`**.

    ```plaintext
    vault write auth/userpass/users/bob password="password"
    ```
    ![Create Policy](/img/vault-multi-tenant-4.png)

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

1. Now, set the **CURRENT NAMESPACE** to **education/training**.
    ![Namespace](/img/vault-multi-tenant-6.png)

1. In the **Access** tab, select **Groups**, and select **Create group**.

1. Paste in the entity ID in the **Member Entity IDs** field you copied.

1. Enter **`Training Admin`** in the **Name** field, **`training-admin`** in the
**Policies** field, and click **Create**.

1. Click the **Access** tab, and select **Enable new method**.

1. Select **Username & Password** from the **Type** drop-down menu.

1. Click **Enable Method**.  Copy the mount accessor value which you will user later.
    ![Namespace](/img/vault-multi-tenant-8.png)

1. Click the Vault CLI shell icon (**`>_`**) to open a command shell.  Enter the
following command to create a new user, **`bsmith`**.

    ```plaintext
    vault write auth/userpass/users/bsmith password="password"
    ```

1. Set the **CURRENT NAMESPACE** back to **education**.

1. In the command shell, enter the following command.  Be sure to replace the
`<Bob_Smith_entity_id>` with the value you copied at step 13, and
`<mount_accessor>` with the value you copied at step 20.

    ```plaintext
    vault write identity/entity-alias name="bsmith" \
            canonical_id=<Bob_Smith_entity_id> mount_accessor=<mount_accessor>
    ```


### <a name="step4"></a>Step 4: Test the organization admin user
(**Persona:** org-admin)

#### CLI Command

Log in as **`bob`** into the `education` namespace:

```plaintext
$ vault login -namespace=education -method=userpass username="bob" password="password"

Key                    Value
---                    -----
token                  5ai0qpQeCdRHALzEY4Q8sW.28dk2
token_accessor         9xXQmdx6Aq6zw1KX4gpzb.28dk2
token_duration         768h
token_renewable        true
token_policies         ["default"]
identity_policies      ["edu-admin"]
policies               ["default" "edu-admin"]
token_meta_username    bob
```

Notice that the user, `bob` only has `default` policy attached to his token
(`token_policies`); however, he inherited the `edu-admin` policy from the `Bob
Smith` entity (`identity_policies`).

Test to make sure that `bob` can create a namespace, enable secrets engine, and
whatever else that you want to verify.

```shell
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

```plaintext
$ curl --header "X-Vault-Namespace: education" \
       --request POST \
       --data '{"password": "password"}' \
       http://127.0.0.1:8200/v1/auth/userpass/login/bob | jq
{
   ...
   "auth": {
     "client_token": "5ai0qpQeCdRHALzEY4Q8sW.28dk2",
     "accessor": "9xXQmdx6Aq6zw1KX4gpzb.28dk2",
     "policies": [
        "default",
        "edu-admin"
      ],
      "token_policies": [
        "default"
      ],
      "identity_policies": [
        "edu-admin"
      ],
      "external_namespace_policies": {
        "9dKXw": [
          "training-admin"
        ]
      },
     "metadata": {
       "username": "bob"
     },
     ...
   }
}
```

Notice that the user, `bob` only has `default` policy attached to his token
(`token_policies`); however, he inherited the `edu-admin` policy from the `Bob
Smith` entity (`identity_policies`). Also, `training-admin` policy is listed
under `external_namespace_policies` due to its membership to the Training Admin
group in `education/training` namespace.

Verify that `bob` can perform the operations permitted by the `edu-admin`
policy.

```shell
# Create a new namespace called 'web-app'
# Be sure to use generated bob's client token
$ curl --header "X-Vault-Token: 5ai0qpQeCdRHALzEY4Q8sW.28dk2" \
       --request POST \
       http://127.0.0.1:8200/v1/education/sys/namespaces/web-app

# Enable key/value v2 secrets engine at edu-secret
$ curl --header "X-Vault-Token: 5ai0qpQeCdRHALzEY4Q8sW.28dk2" \
       --request POST \
       --data '{"type": "kv-v2"}' \
       http://127.0.0.1:8200/v1/education/sys/mounts/edu-secret
```


#### Web UI

1. Open a web browser and launch the Vault UI (e.g. http://127.0.01:8200/ui). If
you are already logged in, sign out.

1. At the **Sign in to Vault**, set the **Namespace** to **`education`**.

1. Select the **Userpass** tab, and enter **`bob`** in the **Username** field,
and **`password`** in the **Password** field.
    ![Login](/img/vault-multi-tenant-5.png)

1. Click **Sign in**.  Notice that the CURRENT NAMESPACE is set to **education**
in the upper left corner of the UI.

1. To add a new namespace, select **Access**.

1. Select **Namespaces** and then click **Create a namespace**.

1. Enter **`web-app`** in the **Path** field, and then click **Save**.

1. Select **Secrets**, and then **Enable new engine**.

1. Select **KV** from the **Secrets engine type** drop-down list, and enter
**`edu-secret`** in the **Path** field.

1. Click **Enable Engine** to finish.


### <a name="step5"></a>Step 5: Test the team admin user
(**Persona:** team-admin)

#### CLI Command

Log in as **`bsmith`** into the **`education/training`** namespace:

```plaintext
$ vault login -namespace=education/training -method=userpass username="bsmith" password="password"

Key                    Value
---                    -----
token                  5YNNjDDl6D8iW3eGQIlU0q.9dKXw
token_accessor         6TVkDhdvEQXO2JaD64TVLv.9dKXw
token_duration         768h
token_renewable        true
token_policies         ["default"]
identity_policies      ["training-admin"]
policies               ["default" "training-admin"]
token_meta_username    bsmith
```

Notice that the user, `bsmith` inherited the `training-admin` policy from the
`Training Admin` group (`training_admin`) which `Bob Smith` entity is a member
of.

Verify that `bsmith` can perform the operations permitted by the
`training-admin` policy.

```shell
# Set the target namespace as an env variable
$ export VAULT_NAMESPACE="education/training"

# Create a new namespace called 'vault-training'
$ vault namespace create vault-training
Success! Namespace created at: education/training/vault-training/

# Enable key/value v1 secrets engine at team-secret
$ vault secrets enable -path=team-secret -version=1 kv
Success! Enabled the kv secrets engine at: team-secret/
```

When you are done testing, unset the VAULT_NAMESPACE environment variable.  

```plaintext
$ unset VAULT_NAMESPACE
```

#### API call using cURL

Log in as **`bsmith`** into the `education` namespace:

```plaintext
$ curl --header "X-Vault-Namespace: education/training" \
       --request POST \
       --data '{"password": "password"}' \
       http://127.0.0.1:8200/v1/auth/userpass/login/bsmith | jq
{
   ...
   "auth": {
     "client_token": "5YNNjDDl6D8iW3eGQIlU0q.9dKXw",
      "accessor": "6TVkDhdvEQXO2JaD64TVLv.9dKXw",
      "display_name": "education-training-auth-userpass-bsmith",
      "policies": [
        "default",
        "training-admin"
      ],
      "token_policies": [
        "default"
      ],
      "identity_policies": [
        "training-admin"
      ],
      "external_namespace_policies": {
        "28dk2": [
          "edu-admin"
        ]
      },
      "metadata": {
        "username": "bsmith"
      },
     ...
   }
}
```

Notice that the user, `bsmith` inherited the `training-admin` policy from the
`Training Admin` group which `Bob Smith` entity is a member of.  Also,
`edu-admin` policy is listed under `external_namespace_policies`.

Verify that `bsmith` can perform the operations permitted by the
`training-admin` policy.

```shell
# Create a new namespace called 'vault-training'
# Be sure to use generated bsmith's client token
$ curl --header "X-Vault-Token: 5YNNjDDl6D8iW3eGQIlU0q.9dKXw" \
       --request POST \
       http://127.0.0.1:8200/v1/education/training/sys/namespaces/web-app

# Enable key/value v1 secrets engine at team-secret
$ curl --header "X-Vault-Token: 5YNNjDDl6D8iW3eGQIlU0q.9dKXw" \
       --request POST \
       --data '{"type": "kv"}' \
       http://127.0.0.1:8200/v1/education/training/sys/mounts/edu-secret
```

#### Web UI

1. Open a web browser and launch the Vault UI (e.g. http://127.0.01:8200/ui). If
you are already logged in, sign out.

1. At the **Sign in to Vault**, set the **Namespace** to
**`education/training`**.

1. Select the **Userpass** tab, and enter **`bsmith`** in the **Username**
field, and **`password`** in the **Password** field.

1. Click **Sign in**.  

1. To add a new namespace, select **Access**.

1. Select **Namespaces** and then click **Create a namespace**.

1. Enter **`vault-training`** in the **Path** field, and then click **Save**.

1. Select **Secrets**, and then **Enable new engine**.

1. Select **KV** from the **Secrets engine type** drop-down list, and enter
**`team-secret`** in the **Path** field.

1. Click **Enable Engine** to finish.

### <a name="step6"></a>Step 6: Audit ambient credentials
(**Persona:** operator)

Many auth and secrets providers, such as AWS, Azure, GCP, and AliCloud, use ambient
credentials to authenticate API calls. For example, AWS may:

1. Use an access key and secret key configured in Vault.
1. If not present, check for environment variables such as "AWS_ACCESS_KEY_ID" and "AWS_SECRET_ACCESS_KEY".
1. If not present, load credentials configured in "~/.aws/credentials".
1. If not present, use instance metadata.

This becomes a problem if these ambient credentials are not intended to be used within
a particular namespace. 

For example, suppose that your Vault server is running on an 
AWS EC2 instance. You give the owner of a namespace a particular set of permissions to
use for AWS. However, that owner _does not_ configure them. So, Vault falls back to using
the credentials available in instance metadata, leading to a privilege escalation.

To handle this:

- Ensure no environment variables are available that could grant a privilege escalation.
- Ensure that any privileges granted through instance metadata (in this example) or other 
ambient identity info represent a _loss_ of privilege.
- Directly configure the correct credentials in namespaces, and restrict access to that
endpoint so credentials can't later be edited to use ambient credentials.

<br>

~> **Summary:** As this guide demonstrated, each namespace you created behaves
as an **isolated** Vault environment. Once you sign into a namespace, there is
no visibility into other namespaces regardless of its hierarchical relationship.
Tokens, policies, and secrets engines are tied to its namespace; therefore, each
client must acquire a valid token for each namespace to access their secrets.


## Additional Discussion

For the simplicity, this guide used the username and password (`userpass`) auth
method which was enabled in the education namespace.  However, most likely, your
organization uses LDAP auth method which is enabled in the **root** namespace
instead.

Here are the steps to create the "Training Admin" group as described in this
guide using the LDAP auth method enabled in the root namespace.

1. Enable and configure the desired auth method (e.g. LDAP) in the root
namespace.

    ```plaintext
    $ vault auth enable ldap

    $ vault write auth/ldap/config \
            url="ldap://ldap.example.com" \
            userdn="ou=Users,dc=example,dc=com" \
            groupdn="ou=Groups,dc=example,dc=com" \
            groupfilter="(&(objectClass=group)(member:1.2.840.113556.1.4.1941:={{.UserDN}}))" \
            groupattr="cn" \
            upndomain="example.com" \
            certificate=@ldap_ca_cert.pem \
            insecure_tls=false \
            starttls=true
    ```

1. Create an _external_ group in the root namespace.

    ```shell
    # Get the mount accessor for ldap auth method and save it in accessor.txt file
    $ vault auth list -format=json \
            | jq -r '.["ldap/"].accessor' > accessor.txt

    # Create an external group and save the generated group ID in group_id.txt
    $ vault write -format=json identity/group name="training_admin_root" \
            type="external" \
            | jq -r ".data.id" > group_id.txt    

    # Create a group alias - assuming that the group name in LDAP is "ops_training"
    $ vault write -format=json identity/group-alias name="ops_training" \
            mount_accessor=$(cat accessor.txt) \
            canonical_id=$(cat group_id.txt)             
    ```

1. In the `education/training` namespace, create an _internal_ group which has
the external group (`training_admin_root`) as its member.

    ```plaintext
    $ vault write -namespace=education/training identity/group \
            name="Training Admin" \
            policies="training-admin" \
            member_group_ids=$(cat group_id.txt)
    ```

### Policy with namespaces

In this guide, you created policies in each namespace (`education` and
`education/training`).  Therefore, you did not have to specify the target
namespace in the policy paths.  

If you want to create policies in the root namespace to control `education` and
`education/training` namespaces, prepend the namespace in the paths.

For example:

```hcl
# Manage policies in the 'education' namespace
path "education/sys/policies/*" {
   capabilities = ["create", "read", "update", "delete", "list", "sudo"]
}

# Manage tokens in the 'education' namespace
path "education/auth/token/*" {
   capabilities = ["create", "read", "update", "delete", "list", "sudo"]
}

# Manage policies under 'education/training' namespace
path "education/training/sys/policies/*" {
   capabilities = ["create", "read", "update", "delete", "list", "sudo"]
}

# Manage tokens in the 'education/training' namespace
path "education/training/auth/token/*" {
   capabilities = ["create", "read", "update", "delete", "list", "sudo"]
}
...
```

In [Step 2](#step2), you deployed the `training-admin` policy in the
`education/training` namespace. The path is relative to the working namespace.
So, if you want to create the `training-admin` policy in the **`education`**
namespace instead, the paths starts with `training/` rather than
`education/training/`.

```hcl
# Manage namespaces
path "training/sys/namespaces/*" {
   capabilities = ["create", "read", "update", "delete", "list", "sudo"]
}

# Manage policies via API
path "training/sys/policies/*" {
   capabilities = ["create", "read", "update", "delete", "list", "sudo"]
}

# Manage policies via CLI
path "training/sys/policy/*" {
  capabilities = ["create", "read", "update", "delete", "list", "sudo"]
}
...
```

~> **NOTE:** Important to remember that tokens are local to the namespace.
Therefore, you need a valid token for the namespace you want to operate in. The
token created in the `education` namespace is not valid in the
`education/training` namespace. This is so that each namespace is completely
isolated from one another to ensure a secure multi-tenant environment.



## Next steps

Refer to the [Sentinel Policies](/guides/identity/sentinel.html) guide if you
need to write policies that allow you to embed finer control over the user
access across those namespaces.
