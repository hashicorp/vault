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

~> This guide highlights the use of ACL templating which was introduced in
**Vault 0.11**.

## Reference Material

- [Policies](/docs/concepts/policies.html) documentation
- [Policy API](/api/system/policy.html) documentation
- [Identity Secrets Engine](/docs/secrets/identity/index.html)
- [Identity: Entities and Groups](/guides/identity/identity.html) guide


## Estimated Time to Complete

10 minutes


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

> **NOTE:** This feature leverages [Vault
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

- Each user can perform all operations on their allocated key/value secret path (`user-kv/<user_name>`)

- Each group can have its own key/value secret path where all operations can
be performed by the group members (`group-kv/<group_name>`)

- Each group can update the group information such as metadata about the group

<br>

You are going to perform the following:

1. [Write templated ACL policies](#step1)
1. [Deploy your policies](#step2)
1. [Create entities and groups](#step3)
1. [Test the ACL templating](#step4)


### <a name="step1"></a>Step 1: Write templated ACL policies

Policy authors can pass in a policy path containing templating delimiters which
is double curly braces (**`{{<parameter>}}`**).

**Example:**

```hcl
path "auth/userpass/users/{{entity.name}}/password" {
  capabilities = [ "update" ]
}
```

#### Available Templating Parameters

| Parameter expression    | Description                                        |
|-------------------------|----------------------------------------------------|
| `entity.id`               | Identity entity ID                             
| `entity.metadata.<key>`   | Entity metadata value for `<key>`
| `group.id`                | Identity Group ID
| `group.<group_id>.id`     | The group ID of a particular group (`<group_id>`)
| `group.<group_id>.name`   | The name of a particular group (`<group_id>`)
| `group.<group_id>.metadata.<key>` | The metadata value for `<key>` for a particular group (`<group_id>`)

Identity groups are not directly attached to a token and an entity can be
associated with multiple groups.  Therefore, in order to reference a group, the
exact group ID must be provided (e.g.
`group.59f001d5-dd49-6d63-51e4-357c1e7a4d44.name`).

~> The **`group.id`** parameter enumerates all groups associated with an entity
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

**`user-secrets.hcl`**

```hcl
# Grant permissions on user specific path
path "user-kv/{{entity.name}}/*" {
	capabilities = [ "create", "update", "read", "delete", "list" ]
}

# Grant permissions on the group specific path that this user is a part of
path "group-kv/{{group.id}}/*" {
	capabilities = [ "create", "update", "read", "delete", "list" ]
}
```









### <a name="step2"></a>Step 2: Deploy your policies

#### CLI command

```plaintext
$ vault policy write user-tmpl user-tmpl.hcl
```


#### API call using cURL

To create a policy, use the `/sys/policies/acl` endpoint:

```shell
$ curl --header "X-Vault-Token: <TOKEN>" \
       --request PUT \
       --data <PAYLOAD> \
       <VAULT_ADDRESS>/v1/sys/policies/acl/<POLICY_NAME>
```

Where `<TOKEN>` is your valid token, and `<PAYLOAD>` includes the policy name and
stringified policy.

**Example:**

```shell
# API request payload
$ tee admin-payload.json <<EOF
{
  "policy": "path \"auth/*\" {\n capabilities = [\"create\", \"read\", \"update\", ... }"
}
EOF

# Create admin policy
$ curl --header "X-Vault-Token: ..." \
       --request PUT
       --data @payload.json \
       http://127.0.0.1:8200/v1/sys/policies/acl/user-tmpl
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

1. Repeat the steps to create a policy for **`acct_manager`**.




## Next steps

To learn about Sentinel policies, refer to the [Sentinel
Policies](/guides/identity/sentinel.html) guide.
