---
layout: "guides"
page_title: "Identity: Entities and Groups - Guides"
sidebar_title: "Identity - Entities & Groups"
sidebar_current: "guides-identity-identity"
description: |-
  This guide demonstrates the commands to create entities, entity aliases, and
  groups.  For the purpose of the demonstration, userpass auth method will be
  used.  
---

# Identity - Entities and Groups

Vault supports multiple authentication methods and also allows enabling the same
type of authentication method on different mount paths. Each Vault client may
have multiple accounts with various identity providers that are enabled on the
Vault server.

Vault clients can be mapped as ***entities*** and their corresponding accounts
with authentication providers can be mapped as ***aliases***. In essence, each
entity is made up of zero or more aliases. Identity secrets engine internally
maintains the clients who are recognized by Vault.

## Reference Material

- [Identity Secrets Engine](/docs/secrets/identity/index.html)
- [Identity Secrets Engine (API)](/api/secret/identity/index.html)
- [External vs Internal Groups](/docs/secrets/identity/index.html#external-vs-internal-groups)

~> **NOTE:** An [interactive
tutorial](https://www.katacoda.com/hashicorp/scenarios/vault-identity) is
also available if you do not have a Vault environment to perform the steps
described in this guide.

## Estimated Time to Complete

10 minutes

## Personas

The steps described in this guide are typically performed by **operations**
persona.


## Challenge

Bob has accounts in both Github and LDAP.  Both Github and LDAP auth methods are
enabled on the Vault server that he can authenticate using either one of his
accounts. Although both accounts belong to Bob, there is no association between
the two accounts to set some common properties.

## Solution

Create an _entity_ representing Bob, and associate aliases representing each of
his accounts as the entity member. You can set additional policies and metadata
on the entity level so that both accounts can inherit.

When Bob authenticates using either one of his accounts, the entity identifier
will be tied to the authenticated token. When such tokens are put to use, their
entity identifiers are audit logged, marking a trail of actions performed by
specific users.


## Prerequisites

To perform the tasks described in this guide, you need to have a Vault
environment.  Refer to the [Getting
Started](/intro/getting-started/install.html) guide to install Vault. Make sure
that your Vault server has been [initialized and
unsealed](/intro/getting-started/deploy.html).


### Policy requirements

-> **NOTE:** For the purpose of this guide, you can use the **`root`** token to work
with Vault. However, it is recommended that root tokens are used for just
enough initial setup or in emergencies. As a best practice, use tokens with
an appropriate set of policies based on your role in the organization.

To perform all tasks demonstrated in this guide, your policy must include the
following permissions:

```shell
# Configure auth methods
path "sys/auth" {
  capabilities = [ "read", "list" ]
}

# Configure auth methods
path "sys/auth/*" {
  capabilities = [ "create", "update", "read", "delete", "list", "sudo" ]
}

# Manage userpass auth methods
path "auth/userpass/*" {
  capabilities = [ "create", "read", "update", "delete" ]
}

# Manage github auth methods
path "auth/github/*" {
  capabilities = [ "create", "read", "update", "delete" ]
}

# Display the Policies tab in UI
path "sys/policies" {
  capabilities = [ "read", "list" ]
}

# Create and manage ACL policies from UI
path "sys/policies/acl/*" {
  capabilities = [ "create", "read", "update", "delete", "list" ]
}

# Create and manage policies
path "sys/policy" {
  capabilities = [ "read", "list" ]
}

# Create and manage policies
path "sys/policy/*" {
  capabilities = [ "create", "read", "update", "delete", "list" ]
}

# List available secret engines to retrieve accessor ID
path "sys/mounts" {
  capabilities = [ "read" ]
}

# Create and manage entities and groups
path "identity/*" {
  capabilities = [ "create", "read", "update", "delete", "list" ]
}
```

If you are not familiar with policies, complete the
[policies](/guides/identity/policies.html) guide.


## Steps

In this lab, you are going to learn the API-based commands to create entities,
entity aliases, and groups.  For the purpose of the training, you are going to
leverage the userpass auth method.  The challenge exercise walks you through
creating an external group by mapping a GitHub group to an identity group.

1. [Create an Entity with Alias](#step1)
2. [Test the Entity](#step2)
3. [Create an Internal Group](#step3)
4. [Create an External Group](#step4)



### <a name="step1"></a>Step 1: Create an Entity with Alias

You are going to create a new entity with base policy assigned.  The entity
defines two entity aliases with each has a different policy assigned.

**Scenario:**  A user, Bob Smith at ACME Inc. happened to have two sets of
credentials: `bob` and `bsmith`.  He can authenticate with Vault using either
one of his accounts.  To manage his accounts and link them to identity `Bob
Smith` in QA team, you are going to create an entity for Bob.

![Entity Bob Smith](/img/vault-entity-1.png)

-> For the simplicity of this guide, you are going to work with the `userpass`
auth method.  However, in reality, the user `bob` might be a username exists in
Active Directory, and `bsmith` might be Bob's username in GitHub.

#### Scenario Policies

**`base.hcl`**

```hcl
path "secret/training_*" {
   capabilities = ["create", "read"]
}
```

**`test.hcl`**

```hcl
path "secret/test" {
   capabilities = [ "create", "read", "update", "delete" ]
}
```

**`team-qa.hcl`**

```hcl
path "secret/team-qa" {
   capabilities = [ "create", "read", "update", "delete" ]
}
```

~> **NOTE:** If you are running [K/V Secrets Engine v2](/api/secret/kv/kv-v2.html)
at `secret`, set the policies path accordingly: `secret/data/training_*`,
`secret/data/test`, and `secret/data/team-qa`.

Now, you are going to create `bob` and `bsmith` users with appropriate policies
attached.



#### CLI command

1. Create policies: `base`, `test`, and `team-qa`.

    ```shell
    # Create base policy
    $ vault policy write base base.hcl

    # Create test policy
    $ vault policy write test test.hcl

    # Create team-qa policy
    $ vault policy write team-qa team-qa.hcl

    # List all policies to verify that 'base', 'test' and 'team-qa' policies exist
    $ vault policy list
    base
    default
    team-qa
    test
    root
    ```

1. Enable the `userpass` auth method.

    ```plaintext
    $ vault auth enable userpass
    ```

1. Create a new user in userpass:
    - username: bob
    - password: training
    - policy: test

    ```plaintext
    $ vault write auth/userpass/users/bob password="training" policies="test"
    ```

1. Create another user in userpass:
    - username: bsmith
    - password: training
    - policy: team-qa

    ```plaintext
    $ vault write auth/userpass/users/bsmith password="training" policies="team-qa"
    ```

1. Execute the following command to discover the mount accessor for the userpass auth method:

    ```plaintext
    $ vault auth list -detailed
    Path                  Type        Accessor                ...
    ----                  ----        --------                ...  
    token/                token       auth_token_bec8530a     ...
    userpass/             userpass    auth_userpass_70eba76b  ...
    ```

    In the output, locate the **Accessor** value for `userpass`.

    Run the following command to store the userpass accessor value in a file named, `accessor.txt`.

    ```plaintext
    $ vault auth list -format=json | jq -r '.["userpass/"].accessor' > accessor.txt
    ```

1. Create an entity for `bob-smith`.

    ```plaintext
    $ vault write identity/entity name="bob-smith" policies="base" \
         metadata=organization="ACME Inc." \
         metadata=team="QA"

    Key        Value
    ---        -----
    aliases    <nil>
    id         631256b1-8523-9838-5501-d0a1e2cdad9c         
    ```

    -> Make a note of the generated entity ID (**`id`**).


1. Now, add the user `bob` to the `bob-smith` entity by creating an entity alias:

    ```plaintext
    $ vault write identity/entity-alias name="bob" \
         canonical_id=<entity_id> \
         mount_accessor=<userpass_accessor>
    ```

    The `<userpass_accessor>` value is stored in `accessor.txt`.

    **Example:**

    ```plaintext
    $ vault write identity/entity-alias name="bob" \
           canonical_id="631256b1-8523-9838-5501-d0a1e2cdad9c" \
           mount_accessor=$(cat accessor.txt)

    Key             Value
    ---             -----
    canonical_id    631256b1-8523-9838-5501-d0a1e2cdad9c
    id              873f7b12-dec8-c182-024e-e3f065d8a9f1
    ```

1. Repeat the step to add user `bsmith` to the `bob-smith` entity.

    **Example:**

    ```plaintext
    $ vault write identity/entity-alias name="bsmith" \
           canonical_id="631256b1-8523-9838-5501-d0a1e2cdad9c" \
           mount_accessor=$(cat accessor.txt)

    Key             Value
    ---             -----
    canonical_id    631256b1-8523-9838-5501-d0a1e2cdad9c
    id              55d46747-b99e-6a82-05f5-61bb60fd7d15
    ```

1. Review the entity details.

    ```plaintext
    $ vault read identity/entity/id/<entity_id>
    ```

    The output should include the entity aliases, metadata (organization, and
    team), and base policy.




#### API call using cURL

1. Create policies: `base`, `test`, and `team-qa`.

    To create a policy, use the `/sys/policy` endpoint:

    ```shell
    $ curl --header "X-Vault-Token: <TOKEN>" \
           --request PUT \
           --data <PAYLOAD> \
           <VAULT_ADDRESS>/v1/sys/policy/<POLICY_NAME>
    ```

    Where `<TOKEN>` is your valid token, and `<PAYLOAD>` includes the policy name and
    stringified policy.

    **Example:**

    ```shell
    # Create the API request payload, payload-1.json
    $ tee payload-1.json <<EOF
    {
      "policy": "path \"secret/training_*\" {\n capabilities = [\"create\", \"read\"]\n}"
    }
    EOF

    # Create base policy
    $ curl --header "X-Vault-Token: ..." \
           --request PUT \
           --data @payload-1.json \
           http://127.0.0.1:8200/v1/sys/policy/base

    # Create the API request payload, payload-2.json
    $ tee payload-2.json <<EOF
    {
      "policy": "path \"secret/test\" {\n capabilities = [ \"create\", \"read\", \"update\", \"delete\" ]\n }"
    }
    EOF

    # Create base policy
    $ curl --header "X-Vault-Token: ..." \
           --request PUT \
           --data @payload-2.json \
           http://127.0.0.1:8200/v1/sys/policy/test

    # Create the API request payload, payload-1.json
    $ tee payload-3.json <<EOF
    {
      "policy": "path \"secret/team-qa\" {\n capabilities = [ \"create\", \"read\", \"update\", \"delete\" ]\n }"
    }
    EOF

    # Create base policy
    $ curl --header "X-Vault-Token: ..." \
           --request PUT \
           --data @payload-3.json \
           http://127.0.0.1:8200/v1/sys/policy/team-qa

    # List all policies to verify that 'base', 'test' and 'team-qa' policies exist
    $ curl --header "X-Vault-Token: ..." \
           http://127.0.0.1:8200/v1/sys/policy | jq
    ```

1. Enable the `userpass` auth method.

    ```plaintext
    $ curl --header "X-Vault-Token: ..." \
           --request POST \
           --data '{"type": "userpass"}' \
           http://127.0.0.1:8200/v1/sys/auth/userpass
    ```

1. Create a new user in userpass:
    - username: bob
    - password: training
    - policy: test

    ```plaintext
    $ curl --header "X-Vault-Token: ..." \
           --request POST \
           --data '{"password": "training", "policies": "test"}' \
           http://127.0.0.1:8200/v1/auth/userpass/users/bob
    ```

1. Create another user in userpass:
    - username: bsmith
    - password: training
    - policy: team-qa

    ```plaintext
    $ curl --header "X-Vault-Token: ..." \
           --request POST \
           --data '{"password": "training", "policies": "team-qa"}' \
           http://127.0.0.1:8200/v1/auth/userpass/users/bsmith
    ```

1. Execute the following command to discover the mount accessor for the userpass
   auth method.

    ```plaintext
    $ curl --header "X-Vault-Token: ..." \
           http://127.0.0.1:8200/v1/sys/auth | jq
     {
       ...
       "userpass/": {
         "accessor": "auth_userpass_9b6cd254",
        ...
       },
       ...
    ```

    -> Make a note of the userpass accessor value (**`auth_userpass_XXXXX`**).

1. Create an entity for bob-smith.

    ```plaintext
    $ tee payload.json <<EOF
    {
      "name": "bob-smith",
      "metadata": {
        "organization": "ACME Inc.",
        "team": "QA"
      },
      "policies": ["base"]
    }
    EOF

    $ curl --header "X-Vault-Token: ..." \
           --request POST \
           --data @payload.json \
           http://127.0.0.1:8200/v1/identity/entity
    {
      "request_id": "4d4d340f-f4c9-0201-c87e-42cc140a383a",
      "lease_id": "",
      "renewable": false,
      "lease_duration": 0,
      "data": {
        "aliases": null,
        "id": "6ded4d31-481f-040b-11ad-c6db0cb4d211"
      },
      ...
    ```

    -> Make a note of the generated entity ID (**`id`**).

1. Now, add the user `bob` to the `bob-smith` entity by creating an entity alias.
In the request body, you need to pass the userpass name as `name`, the userpass 
accessor value as `mount_accessor`, and the entity id as `canonical_id`.

    **Example:**

    ```plaintext
    $ tee payload-bob.json <<EOF
    {
      "name": "bob",
      "canonical_id": "6ded4d31-481f-040b-11ad-c6db0cb4d211",
      "mount_accessor": "auth_userpass_9b6cd254"
    }
    EOF

    $ curl --header "X-Vault-Token: ..." \
           --request POST \
           --data @payload-bob.json \
           http://127.0.0.1:8200/v1/identity/entity-alias
    ```


1. Repeat the step to add user `bsmith` to the `bob-smith` entity.

    **Example:**

    ```plaintext
    $ tee payload-bsmith.json <<EOF
    {
      "name": "bsmith",
      "canonical_id": "6ded4d31-481f-040b-11ad-c6db0cb4d211",
      "mount_accessor": "auth_userpass_9b6cd254"
    }
    EOF

    $ curl --header "X-Vault-Token: ..." \
           --request POST \
           --data @payload-bsmith.json \
           http://127.0.0.1:8200/v1/identity/entity-alias
    ```

1. Review the entity details. (**NOTE:** Be sure to enter the entity ID matching
  your environment.)

    ```plaintext
    $ curl --header "X-Vault-Token: ..." \
           http://127.0.0.1:8200/v1/identity/entity/id/<ENTITY_ID>
    {
       "request_id": "cc0793bf-fafe-4b2c-fd82-88855712845c",
       "lease_id": "",
       "renewable": false,
       "lease_duration": 0,
       "data": {
         "aliases": [
           {
             "canonical_id": "6ded4d31-481f-040b-11ad-c6db0cb4d211",
             ...
             "mount_type": "userpass",
             "name": "bob"
           },
           {
             "canonical_id": "6ded4d31-481f-040b-11ad-c6db0cb4d211",
             ...
             "mount_type": "userpass",
             "name": "bsmith"
           }
         ],
         ...
    ```

    The `bob` and `bsmith` users should appear in the entity alias list.


#### Web UI

1. Open a web browser and launch the Vault UI (e.g. http://127.0.01:8200/ui)
and then login.

1. Click the **Policies** tab, and then select **Create ACL policy**.

1. Enter **`base`** in the **Name** field, and paste in the [`base.hcl` policy
rules](#scenario-policies) in the **Policy** text editor.

    ![Create Policy](/img/vault-policy-2.png)

1. Click **Create Policy** to complete.

1. Repeat the steps to create policies for **`test`** and **`team-qa`** as well.

    ![Create Policy](/img/vault-policy-1.png)

1. Click the **Access** tab, and select **Enable new method**.

1. Select **Username & Password** from the **Type** drop-down menu.

    ![Create Policy](/img/vault-auth-method-2.png)

1. Click **Enable Method**.

1. Click the Vault CLI shell icon (**`>_`**) to open a command shell.  Enter the
following command to create a new user, **`bob`**:

    ```plaintext
    $ vault write auth/userpass/users/bob password="training" policies="test"
    ```
    ![Create Policy](/img/vault-auth-method-3.png)

1. Enter the following command to create a new user, **`bsmith`**:

    ```plaintext
    $ vault write auth/userpass/users/bsmith password="training" policies="team-qa"
    ```
    ![Create Policy](/img/vault-auth-method-4.png)

1. Click the icon (**`>_`**) again to hide the shell.

1. From the **Access** tab, select **Entities** and then **Create entity**.

1. Populate the **Name**, **Policies** and **Metadata** fields as shown below:

    ![Create Policy](/img/vault-entity-4.png)

1. Click **Create**.

1. Select **Add alias**.  Enter **`bob`** in the **Name** field and select
**`userpass/ (userpass)`** from the **Auth Backend** drop-down list.

    ![Create Policy](/img/vault-entity-5.png)

1. Click **Create**.

1. Return to the **Entities** list.  Select **Add alias** from the **`bob-smith`**
entity menu.

    ![Create Policy](/img/vault-entity-6.png)

1. Enter **`bsmith`** in the **Name** field and select **`userpass/ (userpass)`** from the
**Auth Backend** drop-down list, and then click **Create**.




### <a name="step2"></a>Step 2: Test the Entity

To better understand how a token inherits the capabilities from the entity's
policy, you are going to test it by logging in as `bob`.

### CLI Command

First, login as `bob`.

```plaintext
$ vault login -method=userpass username=bob password=training

Key                    Value
---                    -----
token                  ac318416-0dc1-4311-67e4-b58381c86fde
token_accessor         79cced7b-51df-9523-920f-a1579687516b
token_duration         768h
token_renewable        true
token_policies         ["default" "test"]
identity_policies      ["base"]
policies               ["base" "default" "test"]
token_meta_username    bob
```

> Upon a successful authentication, a token will be returned. Notice that the
output displays **`token_policies`** and **`identity_policies`**. The generated
token has both `test` and `base` policies attached.

The `test` policy grants CRUD operations on the `secret/test` path.  
Test to make sure that you can write secrets in the path.

```plaintext
$ vault kv put secret/test owner="bob"
Success! Data written to: secret/test
```


Although the username `bob` does not have `base` policy attached, the token
inherits the capabilities granted in the base policy because `bob` is a member
of the `bob-smith` entity, and the entity has base policy attached.

Check to see that the bob's token inherited the capabilities.  

```plaintext
$ vault token capabilities secret/training_test
create, read
```

> The `base` policy grants create and read capabilities on
`secret/training_*` path; therefore, `bob` is permitted to run create and
read operations against any path starting with `secret/training_*`.


What about the `secret/team-qa` path?

```plaintext
$ vault token capabilities secret/team-qa
deny
```
￼
The user `bob` only inherits capability from its associating entity's policy.
The user can access the `secret/team-qa` path only if he logs in with
`bsmith` credentials.


~> Log back in with the token you used to configure the entity before proceed to
[Step 3](#step3).


#### API call using cURL

First, login as `bob`.

```plaintext
$ curl --request POST \
       --data '{"password": "training"}' \
       http://127.0.0.1:8200/v1/auth/userpass/login/bob
{
 ...
 "auth": {
   "client_token": "b3c2ac10-9f8f-4e64-9a1c-337236ba20f6",
   "accessor": "92204429-6555-772e-cf51-52492d7f1686",
   "policies": [
     "base",
     "default",
     "test"
   ],
   "token_policies": [
      "default",
      "test"
    ],
    "identity_policies": [
      "base"
    ],
   ...
```

> Upon a successful authentication, a token will be returned. Notice that the
output displays **`token_policies`** and **`identity_policies`**. The generated
token has both `test` and `base` policies attached.

The `test` policy grants CRUD operations on the `secret/test` path. Test
to make sure that you can write secrets in the path.

```plaintext
$ curl --header "X-Vault-Token: ..." \
       --request POST \
       --data '{"owner": "bob"}' \
       http://127.0.0.1:8200/v1/secret/test
```


Although the username `bob` does not have `base` policy attached, the token
inherits the capabilities granted in the base policy because `bob` is a member
of the `bob-smith` entity, and the entity has base policy attached.

Check to see that the bob's token inherited the capabilities.  

```plaintext
$ curl --header "X-Vault-Token: ..." \
         --request POST \
         --data '{"paths": ["secret/training_test"]}'
         http://127.0.0.1:8200/v1/sys/capabilities-self | jq
{
 "secret/training_test": [
   "create",
   "read"
 ],
 ...
```

> The `base` policy grants create and read capabilities on
`secret/training_*` path; therefore, `bob` is permitted to run create and
read operations against any path starting with `secret/training_*`.


What about the `secret/team-qa` path?

```plaintext
$ curl --header "X-Vault-Token: ..." \
       --request POST \
       --data '{"paths": ["secret/team-qa"]}'
       http://127.0.0.1:8200/v1/sys/capabilities-self | jq
{
 "secret/team-qa": [
   "deny"
 ],
 ...
```
￼
The user `bob` only inherits capability from its associating entity's policy.
The user can access the `secret/team-qa` path only if he logs in with
`bsmith` credentials.


!> **NOTE:** Log back in with the token you used to configure the entity before proceed to
[Step 3](#step3).


### <a name="step3"></a>Step 3: Create an Internal Group

Now, you are going to create an internal group named, **`engineers`**.  Its
member is `bob-smith` entity that you created in [Step 1](#step1).

![Entity Bob Smith](/img/vault-entity-3.png)

The group policy, `team-eng` defines the following: **`team-eng.hcl`**

```plaintext
path "secret/team/eng" {
  capabilities = [ "create", "read", "update", "delete"]
}
```

#### CLI Command

1. Create a new policy named, `team-eng`:

    ```plaintext
    $ vault policy write team-eng ./team-eng.hcl
    ```

1. Create an internal group named, `engineers` and add `bob-smith` entity as a
group member and attach `team-eng`.

    ```plaintext
    $ vault write identity/group name="engineers" \
          policies="team-eng" \
          member_entity_ids=<entity_id> \
          metadata=team="Engineering" \
          metadata=region="North America"
    ````
    Where `<entity_id>` is the value you copied at [Step 1](#step1).

    **Example:**

    ```plaintext
    $ vault write identity/group name="engineers" \
          policies="team-eng" \
          member_entity_ids="631256b1-8523-9838-5501..."  \
          metadata=team="Engineering" \
          metadata=region="North America"
    Key     Value
    ---     -----
    id      81bdac90-284a-7b8c-6289-5fa7693bcb4a
    name    engineers
    ```

Now, when you login as `bob` or `bsmith`, its generated token inherits the
group-level policy, **`team-eng`**. You can perform similar tests demonstrated
in [Step 2](#step2) to verify that.


#### API call using cURL

1. Create a new policy named, `team-eng`:

    ```shell
    # API request payload containing stringified policy
    $ tee payload.json <<EOF
    {
      "policy": "path \"secret/team/eng\" {\n capabilities = [\"create\", \"read\", \"delete\", \"update\"]\n }"
    }
    EOF

    # Create base policy
    $ curl --header "X-Vault-Token: ..." \
           --request PUT \
           --data @payload-1.json \
           http://127.0.0.1:8200/v1/sys/policy/team-eng
    ```


1. Create an internal group named, `engineers` and add `bob-smith` entity as a
group member and attach `team-eng`.

    ```shell
    # API request msg payload.  Be sure to replace <ENTITY_ID> with correct value
    $ tee payload-group.json <<EOF
    {
      "name": "engineers",
      "policies": ["team-eng"],
      "member_entity_ids": ["<ENTITY_ID>"],
      "metadata": {
        "team": "Engineering",
        "region": "North America"
      }
    }
    EOF

    # Use identity/group endpoint
    $ curl --header "X-Vault-Token: ..." \
           --request PUT \
           --data @payload-group.json \
           http://127.0.0.1:8200/v1/identity/group | jq
    {
       "request_id": "2b6eefd6-67a6-31c7-dbc3-11c1c132e2cf",
       "lease_id": "",
       "renewable": false,
       "lease_duration": 0,
       "data": {
         "id": "d62157aa-b5f6-b6fe-aa40-0ffc54defc41",
         "name": "engineers"
       },
       ...
    ```

Now, when you login as `bob` or `bsmith`, its generated token inherits the
group-level policy, **`team-eng`**. You can perform similar tests demonstrated
in [Step 2](#step2) to verify that.


#### Web UI

1. Click the **Policies** tab, and then select **Create ACL policy**.

1. Enter **`team-eng`** in the **Name** field, and paste in the [`team-eng.hcl` policy
rules](#step3) in the **Policy** text editor, and then click **Create Policy**.

1. Click the **Access** tab and select **Entities**.

1. Select the **`bob-smith`** entity and copy its **ID** displayed under the
**Details** tab.

1. Now, click **Groups** from the left navigation, and select **Create group**.

1. Enter the group information as shown below.

    ![Group](/img/vault-entity-7.png)

    ~> **NOTE:** Make sure to enter the `bob-smith` entity **ID** you copied in the
    **Member Entity IDs** field.

1. Click **Create**.

Now, when you login as `bob` or `bsmith`, its generated token inherits the
group-level policy, **`team-eng`**. You can perform similar tests demonstrated
in [Step 3](#step3) to verify that.

<br>

> **Summary:** By default, Vault creates an internal group. When you create an
internal group, you specify the ***group members*** rather than ***group
alias***. Group _aliases_ are mapping between Vault and external identity providers
(e.g. LDAP, GitHub, etc.).  Therefore, you define group aliases only when you
create **external** groups.  For internal groups, you specify `member_entity_ids`
and/or `member_group_ids`.



### <a name="step4"></a>Step 4: Create an External Group

It is common for organizations to enable auth methods such as LDAP, Okta and
perhaps GitHub to handle the Vault user authentication, and individual user's
group memberships are defined within those identity providers.

In order to manage the group-level authorization, you can create an external
group to link Vault with the external identity provider (auth provider) and
attach appropriate policies to the group.

#### Example Scenario

Any user who belongs to **`training`** team in GitHub organization,
**`example-inc`** are permitted to perform all operations against the
`secret/education` path.

**NOTE:** This scenario assumes that the GitHub organization, `example-inc`
exists as well as `training` team within the organization.

### CLI Command

```shell
# Write a new policy file
# If you are running KV v2, set the path to "secret/data/education" instead
$ tee education.hcl <<EOF
path "secret/education" {
  capabilities = [ "create", "read", "update", "delete", "list" ]
}
EOF

# Create a new policy named 'education'
$ vault policy write education education.hcl

# Enable GitHub auth method
$ vault auth enable github

# Retrieve the mount accessor for the GitHub auth method and save it in accessor.txt
$ vault auth list -format=json | jq -r '.["github/"].accessor' > accessor.txt

# Configure to point to your GitHub organization (e.g. hashicorp)
$ vault write auth/github/config organization=example-inc

# Create an external group named, "education"
# Be sure to copy the generated group ID
$ vault write identity/group name="education" \
       policies="education" \
       type="external" \
       metadata=organization="Product Education"

# Create a group alias where canonical_id is the group ID
# 'name' is the actual GitHub team name (NOTE: Use slugified team name.)
$ vault write identity/group-alias name="training" \
       mount_accessor=$(cat accessor.txt) \
       canonical_id="<group_ID>"
```



#### API call using cURL

```shell
# API request payload containing stringfied policy
# If you are running KV v2, set the path to "secret/data/education" instead
$ tee payload-pol.json <<EOF
{
  "policy": "path \"secret/education\" {\n capabilities = [\"create\", \"read\", \"delete\", \"update\", \"list\"]\n }"
}
EOF

# Create education policy
$ curl --header "X-Vault-Token: ..." \
       --request PUT \
       --data @payload-pol.json \
       http://127.0.0.1:8200/v1/sys/policy/education

# Enable GitHub Auth Method at github
$ curl --header "X-Vault-Token: ..." \
       --request POST \
       --data '{"type": "github"}' \
       http://127.0.0.1:8200/v1/sys/auth/github

# Configure GitHub auth method by setting organization
$ curl --header "X-Vault-Token: ..." \
       --request POST \
       --data '{"organization": "example-inc"}' \
       http://127.0.0.1:8200/v1/auth/github/config

# Get the github accessor value (**`auth_github_XXXXX`**)
$ curl --header "X-Vault-Token: ..." \
      http://127.0.0.1:8200/v1/sys/auth | jq
{
  ...
  "userpass/": {
    "accessor": "auth_github_91010f60",
   ...
  },
  ...
}

# API request msg payload to create an external group  
$ tee payload-edu.json <<EOF
{
   "name": "education",
   "policies": ["education"],
   "type": "external",
   "metadata": {
     "organization": "Product Education"
   }
}
EOF

# Create an external group named, "education"
# Be sure to copy the group ID (id)
$ curl --header "X-Vault-Token: ..." \
       --request POST \
       --data @payload-edu.json \
       http://127.0.0.1:8200/v1/identity/group | jq
{
   "request_id": "a8161086-13db-f982-4216-7d996eae3fd9",
   "lease_id": "",
   "renewable": false,
   "lease_duration": 0,
   "data": {
     "id": "ea18cb62-2478-d370-b726-a77d1700de80",
     "name": "education"
   },
  ...

# API request msg payload to create a group aliases, training
$ tee payload-training.json <<EOF
{
  "canonical_id": "<GROUP_ID>",
  "mount_accessor": "auth_github_XXXXX",
  "name": "training"
}
EOF

# Create 'training' group alias
$ curl --header "X-Vault-Token: ..." \
       --request POST \
       --data @payload-training.json \
       http://127.0.0.1:8200/v1/identity/group-alias | jq
```

#### Web UI

1. Click the **Policies** tab, and then select **Create ACL policy**.

1. Enter **`education`** in the **Name** field, and enter the following policy
in the **Policy** text editor, and then click **Create Policy**. (**NOTE:** If
you are running KV v2, set the path to **`secret/data/education`** instead.)

    ```plaintext
    path "secret/education" {
      capabilities = [ "create", "read", "update", "delete", "list" ]
    }
    ```

1. Click the **Access** tab and select **Auth Methods**.

1. Select **Enable new method**.

1. Select **GitHub** from the **Type** drop-down menu, and then enter
**`example-inc`** in the **Organization** field.  

1. Click **Enable Method**.

1. Click the **Access** tab and select **Groups**.

1. Select **Create group**. Enter the group information as shown below.

    ![Create Policy](/img/vault-entity-9.png)

1. Click **Create**.

1. Select **Add alias** and enter **`training`** in the **Name** field.  Select
**github/ (github)** from the **Auth Backend** drop-down list.

    ![Create Policy](/img/vault-entity-10.png)

1. Click **Create**.

<br>

> **Summary:** At this point, any GitHub user who belongs to `training`
team within the `example-inc` organization can authenticate with Vault. The
generated token for the user has `education` policy attached.


## Next steps

Now that you have learned about managing user identity using entities and
groups, read the [AppRole Pull
Authentication](/guides/identity/authentication.html) guide to learn how apps or
machines can authenticate with Vault.
