---
layout: "guides"
page_title: "ACL Policy Path Templating - Guides"
sidebar_current: "guides-identity-policy-templating"
description: |-
  As of 0.11, ACL policies support templating to allow non-static policy paths.
---

# ACL Policy Path Templating

Vault operates on a **secure by default** standard, and as such, an empty policy
grants **no permissions** in the system. Therefore, policies must be created to
govern the behavior of clients and instrument Role-Based Access Control (RBAC)
by specifying access privileges (_authorization_).

Since everything in Vault is path based, policy authors must be aware of
all existing paths as well as paths to be created.

The [Policies](/guides/identity/policies.html) guide walks you through the
creation of ACL policies in Vault.

-> This guide highlights the use of ACL templating which was introduced in
**Vault 0.11**.

## Reference Material

- [Templated Policies](/docs/concepts/policies.html#templated-policies)
- [Policy API](/api/system/policy.html)
- [Identity: Entities and Groups](/guides/identity/identity.html)
- [Streamline Secrets Management with Vault Agent and Vault 0.11](https://youtu.be/zDnIqSB4tyA?t=24m37s)

~> **NOTE:** An [interactive
tutorial](https://www.katacoda.com/hashicorp/scenarios/vault-policy-templating)
is also available if you do not have a Vault environment to perform the steps
described in this guide.

## Estimated Time to Complete

10 minutes

## Challenge

The only way to specify non-static paths in ACL policies was to use globs (`*`)
at the end of paths.

```hcl
path "transit/keys/*" {
  capabilities = [ "read" ]
}

path "secret/webapp_*" {
  capabilities = [ "create", "read", "update", "delete", "list" ]
}
```

This makes many management and delegation tasks challenging. For example,
allowing a user to change their own password by invoking the
`auth/userpass/users/<user_name>/password` endpoint can require either a policy
for _every user_ or requires the use of Sentinel which is a part of [Vault
Enterprise](/docs/enterprise/sentinel/index.html).

## Solution

As of **Vault 0.11**, ACL templating capability is available to allow a subset
of user information to be used within ACL policy paths.

-> **NOTE:** This feature leverages [Vault
Identities](/docs/secrets/identity/index.html) to inject values into ACL policy
paths.

## Prerequisites

To perform the tasks described in this guide, you need to have an environment
with **Vault 0.11** or later. Refer to the [Getting
Started](/intro/getting-started/install.html) guide to install Vault.

Alternately, you can use the [Vault
Playground](https://www.katacoda.com/hashicorp/scenarios/vault-playground)
environment.

~> This guide assumes that you know how to create ACL policies. If you don't,
go through the interactive [Policy
tutorial](https://www.katacoda.com/hashicorp/scenarios/vault-policies) or
[Policies](/guides/identity/policies.html) guide first.

### Policy requirements

Since this guide demonstrates the creation of an `admin` policy, log in with the
`root` token if possible. Otherwise, refer to the policy requirement in the
[Policies](/guides/identity/policies.html#policy-requirements) guide.


## Steps

Assume that the following policy requirements were given:

- Each _user_ can perform all operations on their allocated key/value secret
 path (`user-kv/data/<user_name>`)

- The education _group_ has a dedicated key/value secret store for each region where
 all operations can be performed by the group members
 (`group-kv/data/education/<region>`)

- The _group_ members can update the group information such as metadata about
 the group (`identity/group/id/<group_id>`)

In this guide, you are going to perform the following steps:

1. [Write templated ACL policies](#step1)
1. [Deploy your policy](#step2)
1. [Setup an entity and a group](#step3)
1. [Test the ACL templating](#step4)

### <a name="step1"></a>Step 1: Write templated ACL policies

Policy authors can pass in a policy path containing double curly braces as
templating delimiters: `{{<parameter>}}`.


#### Available Templating Parameters

|                                    Name                                |                                    Description                               |
| :--------------------------------------------------------------------- | :--------------------------------------------------------------------------- |
| `identity.entity.id`                                                   | The entity's ID                                                              |
| `identity.entity.name`                                                 | The entity's name                                                            |
| `identity.entity.metadata.<<metadata key>>`                            | Metadata associated with the entity for the given key                        |
| `identity.entity.aliases.<<mount accessor>>.id`                        | Entity alias ID for the given mount                                          |
| `identity.entity.aliases.<<mount accessor>>.name`                      | Entity alias name for the given mount                                        |
| `identity.entity.aliases.<<mount accessor>>.metadata.<<metadata key>>` | Metadata associated with the alias for the given mount and metadata key      |
| `identity.groups.ids.<<group id>>.name`                                | The group name for the given group ID                                        |
| `identity.groups.names.<<group name>>.id`                              | The group ID for the given group name                                        |
| `identity.groups.ids.<<group id>>.metadata.<<metadata key>>`           | Metadata associated with the group for the given key                         |
| `identity.groups.names.<<group name>>.metadata.<<metadata key>>`       | Metadata associated with the group for the given key                         |


-> **NOTE:** Identity groups are not directly attached to a token and an entity
can be associated with multiple groups. Therefore, in order to reference a
group, the **group ID** or **group name** must be provided (e.g.
`identity.groups.ids.59f001d5-dd49-6d63-51e4-357c1e7a4d44.name`).

Example:

This policy allows users to change their own password given that the username
and password are defined in the `userpass` auth method.

```hcl
path "auth/userpass/users/{{identity.entity.aliases.auth_userpass_6671d643.name}}/password" {
  capabilities = [ "update" ]
}
```

#### Write the following policies:

User template (`user-tmpl.hcl`)

```hcl
# Grant permissions on user specific path
path "user-kv/data/{{identity.entity.name}}/*" {
	capabilities = [ "create", "update", "read", "delete", "list" ]
}

# For Web UI usage
path "user-kv/metadata" {
  capabilities = ["list"]
}
```

Group template (`group-tmpl.hcl`)

```hcl
# Grant permissions on the group specific path
# The region is specified in the group metadata
path "group-kv/data/education/{{identity.groups.names.education.metadata.region}}/*" {
	capabilities = [ "create", "update", "read", "delete", "list" ]
}

# Group member can update the group information
path "identity/group/id/{{identity.groups.names.education.id}}" {
  capabilities = [ "update", "read" ]
}

# For Web UI usage
path "group-kv/metadata" {
  capabilities = ["list"]
}

path "identity/group/id" {
  capabilities = [ "list" ]
}
```


### <a name="step2"></a>Step 2: Deploy your policy

- [CLI command](#step2-cli)
- [API call using cURL](#step2-api)
- [Web UI](#step2-ui)

#### <a name="step2-cli"></a>CLI command

```shell
# Create the user-tmpl policy
$ vault policy write user-tmpl user-tmpl.hcl

# Create the group-tmpl policy
$ vault policy write group-tmpl group-tmpl.hcl
```


#### <a name="step2-api"></a>API call using cURL

To create a policy, use the `/sys/policies/acl` endpoint:

```sh
$ curl --header "X-Vault-Token: <TOKEN>" \
       --request PUT \
       --data <PAYLOAD> \
       <VAULT_ADDRESS>/v1/sys/policies/acl/<POLICY_NAME>
```

Where `<TOKEN>` is your valid token, and `<PAYLOAD>` includes the policy name and
stringified policy.

Example:

```shell
# API request payload for user-tmpl
$ tee payload_user.json <<EOF
{
  "policy": "path "user-kv/data/{{identity.entity.name}}/*" {\n capabilities = [ "create", "update", "read", "delete", "list" ]\n } ..."
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
   "policy": "path "group-kv/data/{{identity.group.id}}/*" {\n capabilities = [ "create", "update", "read", "delete", "list" ]\n }"
}
EOF

# Create group-tmpl policy
$ curl --header "X-Vault-Token: ..." \
       --request PUT
       --data @payload_group.json \
       http://127.0.0.1:8200/v1/sys/policies/acl/group-tmpl
```

#### <a name="step2-ui"></a>Web UI

Open a web browser and launch the Vault UI (e.g. http://127.0.0.1:8200/ui) and
then login.

1. Click the **Policies** tab, and then select **Create ACL policy**.

1. Toggle **Upload file**, and click **Choose a file** to select the
   `user-tmpl.hcl` file you wrote at [Step 1](#step1).

    ![Create Policy](/img/vault-ctrl-grp-2.png)

    This loads the policy and sets the **Name** to `user-tmpl`.

1. Click the **Create Policy** button.

1. Repeat the steps to create the `group-tmpl` policy.


### <a name="step3"></a>Step 3: Setup an entity and a group

Let's create an entity, **`bob_smith`** with a user **`bob`** as its entity
alias. Also, create a group, **`education`** and add the **`bob_smith`** entity
as its group member.

![Entity & Group](/img/vault-acl-templating.png)

-> This step only demonstrates CLI commands and Web UI to create
entities and groups. Refer to the [Identity - Entities and
Groups](/guides/identity/identity.html) guide if you need the full details.

- [CLI command](#step3-cli)
- [Web UI](#step3-ui)

#### <a name="step3-cli"></a>CLI command

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
# Save the generated group ID in the group_id.txt file
$ vault write -format=json identity/group name="education" \
      policies="group-tmpl" \
      metadata=region="us-west" \
      member_entity_ids=$(cat entity_id.txt)  \
      | jq -r ".data.id" > group_id.txt
```

#### <a name="step3-ui"></a>Web UI

1. Click the **Access** tab, and select **Enable new method**.

1. Select **Username & Password** from the **Type** drop-down menu.

1. Click **Enable Method**.

1. Click the Vault CLI shell icon (**`>_`**) to open a command shell. Enter the
following command to create a new user, **`bob`**.

    ```plaintext
    $ vault write auth/userpass/users/bob password="training"
    ```
    ![Create Policy](/img/vault-ctrl-grp-3.png)

1. Click the icon (**`>_`**) again to hide the shell.

1. From the **Access** tab, select **Entities** and then **Create entity**.

1. Enter **`bob_smith`** in the **Name** field and enter **`user-tmpl`** in the
**Policies** filed.

1. Click **Create**.

1. Select **Add alias**. Enter **`bob`** in the **Name** field and select
**`userpass/ (userpass)`** from the **Auth Backend** drop-down list.

1. Select the **`bob_smith`** entity and copy its **ID** displayed under the
**Details** tab.

1. Click **Groups** from the left navigation, and select **Create group**.

1. Enter **`education`** in the **Name**, and enter **`group-tmpl`** in the
**Policies** fields. Under **Metadata**, enter **`region`** as a key and
**`us-west`** as the key value. Enter the `bob_smith` entity ID in the **Member
Entity IDs** field.
    ![Group](/img/vault-acl-templating-2.png)

1. Click **Create**.


### <a name="step4"></a>Step 4: Test the ACL templating

- [CLI command](#step4-cli)
- [API call using cURL](#step4-api)
- [Web UI](#step4-ui)

#### <a name="step4-cli"></a>CLI Command

1. Enable key/value v2 secrets engine at `user-kv` and `group-kv` paths.

    ```plaintext
    $ vault secrets enable -path=user-kv kv-v2

    $ vault secrets enable -path=group-kv kv-v2
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

1. Remember that `bob` is a member of the `bob_smith` entity; therefore, the
"`user-kv/data/{{identity.entity.name}}/*`" expression in the `user-tmpl` policy
translates to "**`user-kv/data/bob_smith/*`**". Let's test!

    ```plaintext
    $ vault kv put user-kv/bob_smith/apikey webapp="12344567890"
    Key              Value
    ---              -----
    created_time     2018-08-30T18:28:30.845345444Z
    deletion_time    n/a
    destroyed        false
    version          1
    ```

1. The region was set to `us-west` for the `education` group that the
`bob_smith` belongs to. Therefore, the
"`group-kv/data/education/{{identity.groups.names.education.metadata.region}}/*`"
expression in the `group-tmpl` policy translates to
"**`group-kv/data/education/us-west/*`**". Let's verify.

    ```plaintext
    $ vault kv put group-kv/education/us-west/db_cred password="ABCDEFGHIJKLMN"
    Key              Value
    ---              -----
    created_time     2018-08-30T18:29:02.023749491Z
    deletion_time    n/a
    destroyed        false
    version          1
    ```

1. Verify that you can update the group information. The `group-tmpl` policy
permits "update" and "read" on the
"`identity/group/id/{{identity.groups.names.education.id}}`" path. In [Step
2](#step2), you saved the `education` group ID in the `group_id.txt` file.

    ```plaintext
    $ vault write identity/group/id/$(cat group_id.txt) \
            policies="group-tmpl" \
            metadata=region="us-west" \
            metadata=contact_email="james@example.com"
    ```

    Read the group information to verify that the data has been updated.

    ```plaintext
    $ vault read identity/group/id/$(cat group_id.txt)

    Key                  Value
    ---                  -----
    alias                map[]
    creation_time        2018-08-29T20:38:49.383960564Z
    id                   d6ee454e-915a-4bef-9e43-4ffd7762cd4c
    last_update_time     2018-08-29T22:52:42.005544616Z
    member_entity_ids    [1a272450-d147-c3fd-63ae-f16b65b5ee02]
    member_group_ids     <nil>
    metadata             map[contact_email:james@example.com region:us-west]
    modify_index         3
    name                 education
    parent_group_ids     <nil>
    policies             [group-tmpl]
    type                 internal
    ```

#### <a name="step4-api"></a>API call using cURL

1. Enable key/value v2 secrets engine at `user-kv` and `group-kv` paths.

    ```plaintext
    $ tee payload.json <<EOF
    {
      "type": "kv",
      "options": {
        "version": "2"
      }
    }
    EOF

    $ curl --header "X-Vault-Token: ..." \
           --request POST \
           --data @payload.json \
           https://127.0.0.1:8200/v1/sys/mounts/user-kv

    $ curl --header "X-Vault-Token: ..." \
           --request POST \
           --data @payload.json \
           https://127.0.0.1:8200/v1/sys/mounts/group-kv
    ```

1. Log in as **`bob`**.

    ```plaintext
    $ curl --request POST \
           --data '{"password": "training"}' \
           http://127.0.0.1:8200/v1/auth/userpass/login/bob | jq
    ```

    Copy the generated **`client_token`** value for `bob`.

1. Remember that `bob` is a member of the `bob_smith` entity; therefore, the
"`user-kv/data/{{identity.entity.name}}/*`" expression in the `user-tmpl` policy
translates to "**`user-kv/data/bob_smith/*`**". Let's test!

    ```plaintext
    $ curl --header "X-Vault-Token: <bob_client_token>" \
           --request POST \
           --data '{ "data": {"webapp": "12344567890"} }' \
           http://127.0.0.1:8200/v1/user-kv/data/bob_smith/apikey
    ```

1. The region was set to `us-west` for the `education` group that the
`bob_smith` belongs to. Therefore, the
"`group-kv/data/education/{{identity.groups.names.education.metadata.region}}/*`"
expression in the `group-tmpl` policy translates to
"**`group-kv/data/education/us-west/*`**". Let's verify.

    ```plaintext
    $ curl --header "X-Vault-Token: <bob_client_token>" \
           --request POST \
           --data '{ "data": {"password": "ABCDEFGHIJKLMN"} }' \
           http://127.0.0.1:8200/v1/group-kv/data/education/us-west/db_cred
    ```

1. Verify that you can update the group information. The `group-tmpl` policy
permits "update" and "read" on the
"`identity/group/id/{{identity.groups.names.education.id}}`" path.

    ```plaintext
    $ tee group_info.json <<EOF
    {
      "metadata": {
        "region": "us-west",
        "contact_email": "james@example.com"
      },
      "policies": "group-tmpl"
    }
    EOF

    $ curl --header "X-Vault-Token: <bob_client_token>" \
           --request POST \
           --data @group_info.json \
           http://127.0.0.1:8200/v1/identity/group/id/<education_group_id>
    ```

    Where the group ID is the ID returned in [Step 2](#step2). (NOTE: If you performed
    Step 2 using the CLI commands, the group ID is stored in the `group_id.txt`
    file. If you performed the tasks via Web UI, copy the `education` group ID
    from UI.)

    Read the group information to verify that the data has been updated.

    ```plaintext
    $ curl --header "X-Vault-Token: <bob_client_token>" \
           http://127.0.0.1:8200/v1/identity/group/id/<education_group_id>
    ```


#### <a name="step4-ui"></a>Web UI

1. In **Secrets** tab, select **Enable new engine**.

1. Select the radio-button for **KV**, and then click **Next**.

1. Enter **`user-kv`** in the path field, and then select **2** for KV
version.

1. Click **Enable Engine**.

1. Return to **Secrets** and then select **Enable new engine** again.

1. Select the radio-button for **KV**, and then click **Next**.

1. Enter **`group-kv`** in the path field, and then select **2** for KV
version.

1. Click **Enable Engine**.

1. Now, sign out as the current user so that you can log in as `bob`. ![Sign
off](/img/vault-acl-templating-3.png)

1. In the Vault sign in page, select **Username** and then enter **`bob`** in
the **Username** field, and **`training`** in the **Password** field.

1. Click **Sign in**.

1. Remember that `bob` is a member of the `bob_smith` entity; therefore, the
"`user-kv/data/{{identity.entity.name}}/*`" expression in the `user-tmpl` policy
translates to "**`user-kv/data/bob_smith/*`**". Select **`user-kv`** secrets
engine, and then select **Create secret**.

1. Enter **`bob_smith/apikey`** in the **PATH FOR THIS SECRET** field,
**`webapp`** in the key field, and **`12344567890`** in its value field.

1. Click **Save**. You should be able to perform this successfully.

1. The region was set to `us-west` for the `education` group that the
`bob_smith` belongs to. Therefore, the
"`group-kv/data/education/{{identity.groups.names.education.metadata.region}}/*`"
expression in the `group-tmpl` policy translates to
"**`group-kv/data/education/us-west/*`**".  From the **Secrets** tab, select
**`group-kv`** secrets engine, and then select **Create secret**.

1. Enter **`education/us-west/db_cred`** in the **PATH FOR THIS SECRET** field.
Enter **`password`** in the key field, and **`ABCDEFGHIJKLMN`** in its value
field.

1. Click **Save**. You should be able to perform this successfully.

1. To verify that you can update the group information which is allowed by the
"`identity/group/id/{{identity.groups.names.education.id}}`" expression in the
`group-tmpl` policy, select the **Access** tab.

1. Select **Groups**, and then **`education`**.

1. Select **Edit group**. Add a new metadata where the key is
**`contact_email`** and its value is **`james@example.com`**.

1. Click **Save**. The group metadata should be successfully updated.


## Next steps

To learn about Sentinel policies to implement finer-grained policies, refer to
the [Sentinel Policies](/guides/identity/sentinel.html) guide.
