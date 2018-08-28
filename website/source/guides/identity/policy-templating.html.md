---
layout: "guides"
page_title: "ACL Policy Templating - Guides"
sidebar_current: "guides-identity-policy-templating"
description: |-
  As of 0.11, ACL policies supports templating to allow non-static policy paths.
---

# ACL Policy Templating

Vault operates on a **secure by default** standard, and as such, an empty policy
grants **no permissions** in the system. Therefore, policies must be created to
govern the behavior of clients and instrument Role-Based Access Control (RBAC)
by specifying access privileges (_authorization_).

Since everything in Vault is path based that the policy authors must be aware of
all existing paths as well as paths to be created.

The [Policies](/guides/identity/policies.html) guide walks you through the
creation of ACL policies in Vault.

> This guide highlights the use of ACL templating which was introduced in
**Vault 0.11**.

## Reference Material

- [Policies](/docs/concepts/policies.html) documentation
- [Policy API](/api/system/policy.html) documentation
- [Identity Secrets Engine](/docs/secrets/identity/index.html)
- [Identity: Entities and Groups](/guides/identity/identity.html) guide


## Estimated Time to Complete

5 - 10 minutes


## Challenge

The only way to specify non-static paths in ACL policies was to use globs (`*`)
at the end of its paths.  

```hcl
path "transit/keys/*" {
  capabilities = [ "read" ]
}

path "secret/webapp_*" {
  capabilities = [ "create", "read", "update", "delete", "list" ]
}
```

This makes many management and delegation tasks to be challenging. For example,
allowing a user to change their own password by invoking the
"**`auth/userpass/users/<user_name>/password`**" endpoint can require either a policy
for _every user_ or Sentinel which is a part of Vault Enterprise.

## Solution

As of **Vault 0.11**, ACL templating capability is available to allow a subset
of user information to be used within ACL policy paths.

**NOTE:** This feature leverages [Vault
Identities](/docs/secrets/identity/index.html) to inject values in ACL policy
paths.

## Prerequisites

To perform the tasks described in this guide, you need to have an environment
with **Vault 0.11** or later.  Refer to the [Getting
Started](/intro/getting-started/install.html) guide to install Vault.

Alternatively, you can use the [Vault
Playground](https://www.katacoda.com/hashicorp/scenarios/vault-playground)
environment.

-> This guide assumes that you know how to create ACL policies.  If you don't,
go through the  [Interactive Policy
tutorial](https://www.katacoda.com/hashicorp/scenarios/vault-policies) or
[Policies](/guides/identity/policies.html) guide first.


### Policy requirements

Since this guide demonstrates the creation of an **`admin`** policy, log in with
**`root`** token if possible. Otherwise, refer to the policy requirement in the
[Policies](/guides/identity/policies.html#policy-requirements) guide.


## Steps

Assuming that the following requirements exist:

- Each user can perform all operations on their allocated key/value secret path (**`user-kv/<user_name>`**)

- Each group can have its own key/value secret path where all operations can
be performed by the group members (**`group-kv/<group_name>`**)

- Each group can update the group information such as metadata about the group

<br>

You are going to perform the following:

1. [Write templated ACL policies](#step1)
1. [Deploy your policy](#step2)
1. [Setup an entity and a group](#step3)
1. [Test the ACL templating](#step4)


### <a name="step1"></a>Step 1: Write templated ACL policies

Policy authors can pass in a policy path containing templating delimiters which
is double curly braces (**`{{<parameter>}}`**).

**Example:**

```hcl
path "auth/userpass/users/{{identity.entity.name}}/password" {
  capabilities = [ "update" ]
}
```

#### Available Templating Parameters

| Parameter expression    | Description                                        |
|-------------------------|----------------------------------------------------|
| `entity.id`               | Identity entity ID                             
| `entity.name`             | Identity entity name
| `entity.metadata.<key>`   | Entity metadata value for `<key>`
| `group.id`                | Identity group ID
| `group.name`              | Identity group name
| `group.<group_id>.id`     | The group ID of a particular group (`<group_id>`)
| `group.<group_id>.name`   | The name of a particular group (`<group_id>`)
| `group.<group_id>.metadata.<key>` | The metadata value for `<key>` for a particular group (`<group_id>`)

Identity groups are not directly attached to a token and an entity can be
associated with multiple groups.  Therefore, in order to reference a group, the
exact group ID must be provided (e.g.
`identity.group.59f001d5-dd49-6d63-51e4-357c1e7a4d44.name`).

~> The **`identity.group.id`** parameter enumerates all groups associated with an entity
and add a path for each group that the entity is a member of. This allows
functionality like group-based storage areas for every group where gaining
access to those storage areas is as easy as being added to new groups.


**Example:** A token is attached to the entity ID. The entity's group membership
information is not attached to the token. The example below shows that the
entity belongs to two different groups.  

```plaintext
$ vault token lookup

Key                  Value
---                  -----
accessor             aab457ae-a9bc-06bb-60e5-895c8c378bbe
creation_time        1534364708
creation_ttl         768h
display_name         userpass-bob
entity_id            04480b3e-b796-8217-740e-cad118e5365a
...

$ vault read identity/entity/id/04480b3e-b796-8217-740e-cad118e5365a

Key                    Value
---                    -----
aliases                [map[canonical_id:04480b3e-b796-8217-740e-cad118e5365a ...
creation_time          2018-08-15T20:17:43.980450493Z
direct_group_ids       [59f001d5-dd49-6d63-51e4-357c1e7a4d44 3528c7dd-953b-8b2f-c1c6-5943e3ad1a4a]
disabled               false
group_ids              [59f001d5-dd49-6d63-51e4-357c1e7a4d44 3528c7dd-953b-8b2f-c1c6-5943e3ad1a4a]
...
```

#### Example ACL policy   

**`user-tmpl.hcl`**

```hcl
# Grant permissions on user specific path
path "user-kv/{{entity.name}}/*" {
	capabilities = [ "create", "update", "read", "delete", "list" ]
}
```

**`group-tmpl.hcl`**

```hcl
# Grant permissions on the group specific path that this user is a part of
path "group-kv/{{group.name}}/*" {
	capabilities = [ "create", "update", "read", "delete", "list" ]
}

# Group member can update the group information
path "identity/group/id/{{group.id}}" {
  capabilities = [ "update", "read" ]
}
```


### <a name="step2"></a>Step 2: Deploy your policy

#### CLI command

```shell
# Create user-tmpl policy
$ vault policy write user-tmpl user-tmpl.hcl

# Create group-tmpl policy
$ vault policy write group-tmpl group-tmpl.hcl
```


#### API call using cURL

To create a policy, use the `/sys/policies/acl` endpoint:

```plaintext
$ curl --header "X-Vault-Token: <TOKEN>" \
       --request PUT \
       --data <PAYLOAD> \
       <VAULT_ADDRESS>/v1/sys/policies/acl/<POLICY_NAME>
```

Where `<TOKEN>` is your valid token, and `<PAYLOAD>` includes the policy name and
stringified policy.

**Example:**

```shell
# API request payload for user-tmpl
$ tee payload_user.json <<EOF
{
  "policy": "path "user-kv/{{identity.entity.name}}/*" {\n capabilities = [ "create", "update", "read", "delete", "list" ]\n } ..."
}
EOF

# Create user-tmpl policy
$ curl --header "X-Vault-Token: ..." \
       --request PUT
       --data @payload_user.json \
       http://127.0.0.1:8200/v1/sys/policies/acl/user-tmpl

# API request payload for group-tmpl
$ tee payload_group.json <<EOF
{
   "policy": "path "group-kv/{{identity.group.id}}/*" {\n capabilities = [ "create", "update", "read", "delete", "list" ]\n }"
}
EOF

# Create group-tmpl policy
$ curl --header "X-Vault-Token: ..." \
       --request PUT
       --data @payload_group.json \
       http://127.0.0.1:8200/v1/sys/policies/acl/group-tmpl
```

#### Web UI

Open a web browser and launch the Vault UI (e.g. http://127.0.0.1:8200/ui) and
then login.

1. Click the **Policies** tab, and then select **Create ACL policy**.

1. Toggle **Upload file**, and click **Choose a file** to select your
**`user-tmpl.hcl`** file you authored at [Step 1](#step1).

    ![Create Policy](/assets/images/vault-ctrl-grp-2.png)

    This loads the policy and sets the **Name** to be `user-tmpl`.

1. Click **Create Policy** to complete.

1. Repeat the steps to create **`group-tmpl`** policy.


### <a name="step3"></a>Step 3: Setup an entity and a group

Let's create an entity, **`bob_smith`** with a user **`bob`** as its entity
alias.  Also, create a group, **`education`** and add the **`bob_smith`** entity
as its group member.

![Entity & Group](/assets/images/vault-acl-templating.png)

-> This step only demonstrates CLI commands and Web UI to create
entities and groups.  Refer to the [Identity - Entities and
Groups](/guides/identity/identity.html) guide if you need the full details.


The following command uses [`jq`](https://stedolan.github.io/jq/download/) tool
to parse JSON output.

```shell
# Enable userpass
$ vault auth enable userpass

# Create a user, bob
$ vault write auth/userpass/users/bob password="training"

# Retrieve the userpass mount accessor and save it in a file named, accessor.txt
$ vault auth list -format=json | jq -r '.["userpass/"].accessor' > accessor.txt  

# Create bob_smith entity and save the identity ID in the entity_id.txt
$ vault write -format=json identity/entity name="bob_smith" policies="user-tmpl" \
        | jq -r ".data.id" > entity_id.txt

# Add an entity alias for the bob_smith entity
$ vault write identity/entity-alias name="bob" \
       canonical_id=$(cat entity_id.txt) \
       mount_accessor=$(cat accessor.txt)

# Finally, create education group and add bob_smith entity as a member
$ vault write identity/group name="education" \
      policies="group-tmpl" \
      member_entity_ids=$(cat entity_id.txt)  \
      | jq -r ".data.id" > group_id.txt
```

#### Web UI

1. Click the **Access** tab, and select **Enable new method**.

1. Select **Username & Password** from the **Type** drop-down menu.

1. Click **Enable Method**.

1. Click the Vault CLI shell icon (**`>_`**) to open a command shell.  Enter the
following command to create a new user, **`bob`**:

    ```plaintext
    $ vault write auth/userpass/users/bob password="training"
    ```
    ![Create Policy](/assets/images/vault-ctrl-grp-3.png)

1. Click the icon (**`>_`**) again to hide the shell.

1. From the **Access** tab, select **Entities** and then **Create entity**.

1. Enter **`bob_smith`** in the **Name** field and enter **`user-tmpl`** in the
**Policies** filed.

1. Click **Create**.

1. Select **Add alias**.  Enter **`bob`** in the **Name** field and select
**`userpass/ (userpass)`** from the **Auth Backend** drop-down list.

1. Select the **`bob_smith`** entity and copy its **ID** displayed under the
**Details** tab.

1. Click **Groups** from the left navigation, and select **Create group**.

1. Enter **`education`** in the **Name**, and enter **`group-tmpl`**
in the **Policies** fields.

1. Enter the `bob_smith` entity ID in the **Member Entity IDs** field, and
then click **Create**.


### <a name="step4"></a>Step 4: Test the ACL templating

#### CLI Command

1. Enable key/value secrets engine at `user-kv` and `group-kv` paths.

    ```plaintext
    $ vault secrets enable -path=user-tmpl -version=1 kv

    $ vault secrets enable -path=group-tmpl -version=1 kv
    ```

1. Log in as **`bob`**.

    ```plaintext
    $ vault login -method=userpass username="bob" password="training"

    Key                    Value
    ---                    -----
    token                  5f2b2594-f0b4-0a7b-6f51-767345091dcc
    token_accessor         78b652dd-4320-f18f-b882-0732b7ae9ac9
    token_duration         768h
    token_renewable        true
    token_policies         ["default"]
    identity_policies      ["group-tmpl" "user-tmpl"]
    policies               ["default" "group-tmpl" "user-tmpl"]
    token_meta_username    bob
    ```

1. Try writing some secrets at `user-kv/bob_smith` path.

    ```plaintext
    $ vault kv put user-kv/bob_smith/apikey webapp="12344567890"
    Success! Data written to: user-kv/bob_smith/apikey
    ```

1. Now, try accessing the `group-kv/education` path.

    ```plaintext
    $ vault kv put group-kv/education/creds password="12344567890"
    Success! Data written to: group-kv/education/creds
    ```

1. Verify that you can update the group information by adding metadata.

    ```plaintext
    $ vault write identity/group/id/$(cat group_id.txt) metadata=region="US West" policies="group-tmpl"
    ```

#### API call using cURL

1. Log in as **`bob`**.

    ```plaintext
    $ curl --request POST \
           --data '{"password": "training"}' \
           http://127.0.0.1:8200/v1/auth/userpass/login/bob | jq
    ```

    Copy the generated **`client_token`** value.

1. Try writing some secrets at `user-kv/bob_smith` path.

    ```plaintext
    $ curl --header "X-Vault-Token: <bob_client_token>" \
           --request POST \
           --data '{"webapp": "12344567890"}' \
           http://127.0.0.1:8200/v1/user-kv/bob_smith/apikey
    ```

1. Now, try accessing the `group-kv/education` path.

    ```plaintext
    $ curl --header "X-Vault-Token: <bob_client_token>" \
           --request POST \
           --data '{"password": "12344567890"}' \
           http://127.0.0.1:8200/v1/group-kv/education/creds
    ```

1. Verify that you can update the group information by adding metadata.

    ```plaintext
    $ tee payload.json <<EOF
    {
      "metadata": {
        "password": "123456789"
      },
      "policies": "group-tmpl"
    }
    EOF

    $ curl --header "X-Vault-Token: <bob_client_token>" \
           --request POST \
           --data '{"password": "12344567890"}' \
           http://127.0.0.1:8200/v1/group-kv/education/creds
    ```



## Next steps

To learn about Sentinel policies, refer to the [Sentinel
Policies](/guides/identity/sentinel.html) guide.
