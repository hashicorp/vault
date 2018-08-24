---
layout: "guides"
page_title: "Control Groups - Guides"
sidebar_title: "Control Groups"
sidebar_current: "guides-identity-control-groups"
description: |-
  Vault Enterprise has a support for Control Group Authorization which adds
  additional authorization factors to be required before satisfying a request.
---

# Control Groups

~> **Enterprise Only:** Control Groups is a part of _Vault Enterprise Premium_.

Control Groups add additional authorization factors to be required before
processing requests to increase the governance, accountability, and security of
your secrets.  When a control group is required for a request, the requesting
client receives the [wrapping token](/docs/concepts/response-wrapping.html) in
return. Only when all authorizations are satisfied, the wrapping token can be
used to unwrap the requested secrets.   


## Reference Material

- [Vault Enterprise Control Group Support](/docs/enterprise/control-groups/index.html)
- [Policies](http://localhost:4567/docs/concepts/policies.html)
- [Identity Groups](/docs/secrets/identity/index.html)
- [Control Group API](/api/system/control-group.html)
- [Sentinel Policies](/docs/enterprise/sentinel/index.html)

## Estimated Time to Complete

10 minutes


## Personas

The end-to-end scenario described in this guide involves three personas:

- **`admin`** with privileged permissions to create policies and identities
- **processor** with permission to approve secret access
- **controller** with limited permission to access secrets


## Challenge

In order to operate in EU, a company must abide by the [General Data Protection
Regulation (GDPR)](https://www.eugdpr.org/) as of May 2018.  The regulation
enforces two or more controllers jointly determine the purposes and means of
processing ([Chapter 4: Controller and
Processor](https://gdpr-info.eu/chapter-4/)).

Consider the following scenarios:

- Anytime an authorized user requests to read data at "`EU_GDPR_data/orders/*`",
at least two people from the _Security_ group must approve to ensure that the
user has a valid business reason for requesting the data.

- Anytime a database configuration is updated, it requires that one person from
the _DBA_ and one person from _Security_ group must approve it.


## Solution

Use ***Control Groups*** in your policies to implement dual controller
authorization required.


## Prerequisites

To perform the tasks described in this guide, you need to have a ***Vault
Enterprise*** environment.  

This guide assumes that you have some hands-on experience with [ACL
policies](/docs/concepts/policies.html) as well as
[Identities](/docs/secrets/identity/index.html).  If you are not familiar,
go through the following guides first:

- [Policies](/guides/identity/policies.html)
- [Identity - Entities & Groups](/guides/identity/identity.html)

### Policy requirements

Since this guide demonstrates the creation of policies, log in with a highly
privileged token such as **`root`**.  
Otherwise, required permissions to perform
the steps in this guide are:

```shell
# Create and manage ACL policies via CLI
path "sys/policy/*"
{
  capabilities = ["create", "read", "update", "delete", "list", "sudo"]
}

# Create and manage ACL policies via Web UI
path "sys/policies/acl/*"
{
  capabilities = ["create", "read", "update", "delete", "list", "sudo"]
}

# To enable secret engines
path "sys/mounts/*" {
  capabilities = [ "create", "read", "update", "delete" ]
}

# Setting up test data
path "EU_GDPR_data/*"
{
  capabilities = ["create", "read", "update", "delete", "list"]
}

# Manage userpass auth method
path "auth/userpass/*"
{
  capabilities = ["create", "read", "update", "delete", "list"]
}

# List, create, update, and delete auth methods
path "sys/auth/*"
{
  capabilities = ["create", "read", "update", "delete"]
}

# Create and manage entities and groups
path "identity/*" {
  capabilities = [ "create", "read", "update", "delete", "list" ]
}
```


## Steps

The scenario in this guide is that a user, **`Bob Smith`** has
_read-only_ permission on the "**`EU_GDPR_data/orders/*`**" path; however,
someone in the **`acct_manager`** group must approve it before he can actually
read the data.

As a member of the **`acct_manager`** group, **`Ellen Wright`** can authorize
Bob's request.

![Scenario](/img/vault-ctrl-grp-1.png)

You are going to perform the following:

1. [Implement a control group](#step1)
1. [Deploy the policies](#step2)
1. [Setup entities and a group](#step3)
1. [Verification](#step4)
1. [ACL Policies vs. Sentinel Policies](#step5)


-> Step 1, 2 and 3 are the tasks need to be performed by administrators or
operators who have the privileges to create policies and configure entities and
groups.


### <a name="step1"></a>Step 1: Implement a control group
(**Persona:** admin)

1. Author a policy named, **`read-gdpr-order.hcl`**.

    Bob needs "`read`" permit on "`EU_GDPR_data/orders/*`":

    ```hcl
    path "EU_GDPR_data/orders/*" {
    	capabilities = [ "read" ]
    }
    ```

    Now, add control group to this policy:

    ```hcl
    path "EU_GDPR_data/orders/*" {
    	capabilities = [ "read" ]

    	control_group = {
    		factor "authorizer" {
    			identity {
    				group_names = [ "acct_manager" ]
    				approvals = 1
    			}
    		}
    	}
    }
    ```

    For the purpose of this guide, the number of **`approvals`** is set to
    **`1`** to keep it simple and easy to test. Any member of the identity
    group, **`acct_manager`** can approve the read request. Although this
    example has only one factor (`authorizer`), you can add as many factor
    blocks as you need.

1. Now, write another policy for the **`acct_manager`** group named
**`acct_manager.hcl`**.

    ```hcl
    # To approve the request
    path "sys/control-group/authorize" {
        capabilities = ["create", "update"]
    }

    # To check control group request status
    path "sys/control-group/request" {
        capabilities = ["create", "update"]
    }
    ```

    > **NOTE:** The important thing here is that the authorizer (`acct_manager`)
    must have `create` and `update` permission on the
    **`sys/control-group/authorize`** endpoint so that they can approve the request.


1. Enable key/value secrets engine at **`EU_GDPR_data`** and write some mock data:

    ```shell
    # Enable kv-v1 at EU_GDPR_data
    $ vault secrets enable -path=EU_GDPR_data -version=1 kv

    # Write some mock data
    $ vault kv put EU_GDPR_data/orders/acct1 order_number="12345678" product_id="987654321"
    ```

### <a name="step2"></a>Step 2: Deploy the policies
(**Persona:** admin)

Deploy the `read-gdpr-order` and `acct_manager` policies that you wrote.

#### CLI command

```shell
# Create read-gdpr-order policy
$ vault policy write read-gdpr-order read-gdpr-order.hcl

# Create acct_manager policy
$ vault policy write acct_manager acct_manager.hcl
```


#### API call using cURL

```shell
# Construct API request payload to create read-gdpr-read policy
$ tee payload-1.json <<EOF
{
  "policy": "path \"EU_GDPR_data/orders/*\" {capabilities = [ \"read\" ]control_group = {factor \"authorizer\" ..."
}
EOF

# Create read-gdpr-order policy
$ curl --header "X-Vault-Token: ..." \
       --request PUT \
       --data @payload-1.json \
       http://127.0.0.1:8200/v1/sys/policies/acl/read-gdpr-order

# Construct API request payload to create acct_manager policy
$ tee payload-2.json <<EOF
{
 "policy": "path \"sys/control-group/authorize\" {capabilities = [\"create\", \"update\"]} ..."
}
EOF

# Create acct_manager policy
$ curl --header "X-Vault-Token: ..." \
      --request PUT \
      --data @payload-2.json \
      http://127.0.0.1:8200/v1/sys/policies/acl/acct_manager
```

#### Web UI

Open a web browser and launch the Vault UI (e.g. http://127.0.0.1:8200/ui) and
then login.

1. Click the **Policies** tab, and then select **Create ACL policy**.

1. Toggle **Upload file**, and click **Choose a file** to select your **`read-gdpr-order.hcl`** file you authored at [Step 1](#step1).

    ![Create Policy](/img/vault-ctrl-grp-2.png)

    This loads the policy and sets the **Name** to be `read-gdpr-order`.

1. Click **Create Policy** to complete.

1. Repeat the steps to create a policy for **`acct_manager`**.



### <a name="step3"></a>Step 3: Setup entities and a group
(**Persona:** admin)

-> This step only demonstrates CLI commands and Web UI to create
entities and groups.  Refer to the [Identity - Entities and
Groups](/guides/identity/identity.html) guide if you need the full details.

Now you have policies, let's create a user, **`bob`** and an **`acct_manager`**
group with **`ellen`** as a group member.

> **NOTE:** For the purpose of this guide, use `userpass` auth method to create
user `bob` and `ellen` so that the scenario can be easily tested.

#### CLI command

The following command uses [`jq`](https://stedolan.github.io/jq/download/) tool
to parse JSON output.

```shell
# Enable userpass
$ vault auth enable userpass

# Create a user, bob
$ vault write auth/userpass/users/bob password="training"

# Create a user, ellen
$ vault write auth/userpass/users/ellen password="training"

# Retrieve the userpass mount accessor and save it in a file named, accessor.txt
$ vault auth list -format=json | jq -r '.["userpass/"].accessor' > accessor.txt  

# Create Bob Smith entity and save the identity ID in the entity_id_bob.txt
$ vault write -format=json identity/entity name="Bob Smith" policies="read-gdpr-order" \
        metadata=team="Processor" \
        | jq -r ".data.id" > entity_id_bob.txt

# Add an entity alias for the Bob Smith entity
$ vault write identity/entity-alias name="bob" \
       canonical_id=$(cat entity_id_bob.txt) \
       mount_accessor=$(cat accessor.txt)

# Create Ellen Wright entity and save the identity ID in the entity_id_ellen.txt
$ vault write -format=json identity/entity name="Ellen Wright" policies="default" \
        metadata=team="Acct Controller" \
        | jq -r ".data.id" > entity_id_ellen.txt

# Add an entity alias for the Ellen Wright entity
$ vault write identity/entity-alias name="ellen" \
       canonical_id=$(cat entity_id_ellen.txt) \
       mount_accessor=$(cat accessor.txt)

# Finally, create acct_manager group and add Ellen Wright entity as a member
$ vault write identity/group name="acct_manager" \
      policies="acct_manager" \
      member_entity_ids=$(cat entity_id_ellen.txt)  
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
    ![Create Policy](/img/vault-ctrl-grp-3.png)

1. Enter the following command to create a new user, **`ellen`**:

    ```plaintext
    $ vault write auth/userpass/users/ellen password="training"
    ```

1. Click the icon (**`>_`**) again to hide the shell.

1. From the **Access** tab, select **Entities** and then **Create entity**.

1. Populate the **Name**, **Policies** and **Metadata** fields as shown below.

    ![Create Entity](/img/vault-ctrl-grp-7.png)

1. Click **Create**.

1. Select **Add alias**.  Enter **`bob`** in the **Name** field and select
**`userpass/ (userpass)`** from the **Auth Backend** drop-down list.

1. Return to the **Entities** tab and then **Create entity**.

1. Populate the **Name**, **Policies** and **Metadata** fields as shown below.

    ![Create Entity](/img/vault-ctrl-grp-4.png)

1. Click **Create**.

1. Select **Add alias**.  Enter **`ellen`** in the **Name** field and select
**`userpass/ (userpass)`** from the **Auth Backend** drop-down list.

1. Click **Create**.

1. Select the **`Ellen Wright`** entity and copy its **ID** displayed under the
**Details** tab.

1. Click **Groups** from the left navigation, and select **Create group**.

1. Enter **`acct_manager`** in the **Name**, and again enter **`acct_manager`**
in the **Policies** fields.

1. Enter the `Ellen Wright` entity ID in the **Member Entity IDs** field, and
then click **Create**.


### <a name="step4"></a>Step 4: Verification
(**Persona:** bob and ellen)

Now, let's see how the control group works.

#### CLI Command

1. Log in as **`bob`**.

    ```plaintext
    $ vault login -method=userpass username="bob" password="training"
    ```

1. Request to read "`EU_GDPR_data/orders/acct1`":

    ```plaintext
    $ vault kv get EU_GDPR_data/orders/acct1

    Key                              Value
    ---                              -----
    wrapping_token:                  1f1411bc-2f18-551a-5e58-0fe44432e9a5
    wrapping_accessor:               bbb4deef-e06d-9b2a-64a9-56f815c69ee7
    wrapping_token_ttl:              24h
    wrapping_token_creation_time:    2018-08-08 09:36:32 -0700 PDT
    wrapping_token_creation_path:    EU_GDPR_data/orders/acct1
    ```

    The response includes `wrapping_token` and `wrapping_accessor`.
    Copy this **`wrapping_accessor`** value.

1. Now, a member of `acct_manager` must approve this request.  Log in as
**`ellen`** who is a member of `acct_manager` group.

    ```plaintext
    $ vault login -method=userpass username="ellen" password="training"
    ```

1. As a user, `ellen`, you can check and authorize bob's request using the
following commands.

    ```shell
    # To check the current status
    $ vault write sys/control-group/request accessor=<wrapping_accessor>

    # To approve the request
    $ vault write sys/control-group/authorize accessor=<wrapping_accessor>
    ```

    **Example:**

    ```shell
    # Check the current status
    $ vault write sys/control-group/request accessor=bbb4deef-e06d-9b2a-64a9-56f815c69ee7
    Key               Value
    ---               -----
    approved          false
    authorizations    <nil>
    request_entity    map[name:Bob Smith id:38700386-723d-3d65-43b7-4fb44d7e6c30]
    request_path      EU_GDPR_data/orders/acct1

    # Approve the request
    $ vault write sys/control-group/authorize accessor=bbb4deef-e06d-9b2a-64a9-56f815c69ee7
    Key         Value
    ---         -----
    approved    true
    ```

    Now, the `approved` status is `true`.

1. Since the control group requires one approval from a member of `acct_manager`
group, the condition has been met.  Log back in as `bob` and unwrap the secret.

    **Example:**

    ```shell
    # Log back in as bob - you can use the bob's token: vault login <bob_token>
    $ vault login -method=userpass username="bob" password="training"

    # Unwrap the secrets by passing the wrapping_token
    $ vault unwrap 1f1411bc-2f18-551a-5e58-0fe44432e9a5
    Key                 Value
    ---                 -----
    refresh_interval    768h
    order_number        12345678
    product_id          987654321
    ```


#### API call using cURL

1. Log in as **`bob`**.

    ```plaintext
    $ curl --request POST \
           --data '{"password": "training"}' \
           http://127.0.0.1:8200/v1/auth/userpass/login/bob | jq
    ```

    Copy the generated **`client_token`** value.

1. Request to `EU_GDPR_data/orders/acct1`:

    ```plaintext
    $ curl --header "X-Vault-Token: <bob_client_token>" \
           http://127.0.0.1:8200/v1/EU_GDPR_data/orders/acct1 | jq
    {
       ...
       "wrap_info": {
         "token": "20a2f2b3-8bea-4e16-980b-82724dcdc38b",
         "accessor": "9910cb38-600c-29d8-1c39-764a1c89a481",
         "ttl": 86400,
         "creation_time": "2018-08-08T10:13:06-07:00",
         "creation_path": "EU_GDPR_data/orders/acct1"
       },
       ...
    }
    ```

    The response includes **`wrap_info`** instead of the actual data.
    Copy the **`accessor`** value.

1. Now, a member of `acct_manager` must approve this request.  Log in as
**`ellen`** who is a member of `acct_manager` group.

    ```plaintext
    $ curl --request POST \
           --data '{"password": "training"}' \
           http://127.0.0.1:8200/v1/auth/userpass/login/ellen | jq
    ```

    Copy the generated **`client_token`** value.

1. As a user, `ellen`, you can check the current status and then authorize bob's
request. (NOTE: Be sure to replace `<accessor>` with the `accessor` value you
copied earlier.)

    ```shell
    # To check the current status using sys/control-group/request endpoint
    $ curl --header "X-Vault-Token: <ellen_client_token>" \
           --request POST \
           --data '{"accessor": "<accessor>"}' \
           http://127.0.0.1:8200/v1/sys/control-group/request | jq
    {
       ...
       "data": {
         "approved": false,
         "authorizations": null,
         "request_entity": {
           "id": "38700386-723d-3d65-43b7-4fb44d7e6c30",
           "name": "Bob Smith"
         },
         "request_path": "EU_GDPR_data/orders/acct1"
       },
       ...
    }

    # Now, authorize the request using sys/control-group/authorize endpoint
    $ curl --header "X-Vault-Token: <ellen_client_token>" \
             --request POST \
             --data '{"accessor": "<accessor>"}' \
             http://127.0.0.1:8200/v1/sys/control-group/authorize | jq
    {
       ...
       "data": {
         "approved": true
       },
       ...
    }
    ```

    Now, the `approved` status is `true`.

1. The `bob` user should be able to unwrap the secrets.

    ```plaintext
    $ curl --header "X-Vault-Token: <bob_client_token>" \
           --request POST \
           --data '{"token": "<wrapping_token>"}' \
           http://127.0.0.1:8200/v1/sys/wrapping/unwrap | jq
    {
       ...
       "data": {
         "order_number": "12345678",
         "product_id": "987654321"
       },
       ...
    }
    ```

#### Web UI

The user, **`ellen`** can approve the data access request via UI.


1. Open the Vault sign in page in a web browser (e.g.
http://127.0.0.1:8200/ui/vault/auth?with=userpass).  In the **Userpass** tab,
enter **`ellen`** in the **Username** field, and **`training`** in the
**Password** field.

1. Click **Sign in**.

1. Select the **Access** tab, and then **Control Groups**.

1. Enter the **`wrapping_accessor`** value in the **Accessor** field and click
**Lookup**. ![Control Groups](/img/vault-ctrl-grp-5.png)

1. _Awaiting authorization_ message displays. ![Control
Groups](/img/vault-ctrl-grp-6.png)

1. Click **Authorize**. The message changes to "_Thanks! You have given
authorization_."

> Bob needs to request data access via CLI or API. Once the access request was
approved, use the CLI or API to unwrap the secrets.


### <a name="step5"></a>Step 5: ACL Policy vs. Sentinel Policy

Although the [**`read-gdpr-order.hcl`**](#step1) was written as ACL policy, you
can implement Control Groups in either ACL or Sentinel policies.

Using Sentinel, the same policy may look something like:

```hcl
import "controlgroup"

control_group = func() {
  numAuthzs = 0
  for controlgroup.authorizations as authz {
      if "acct_manager" in authz.groups.by_name {
         numAuthzs = numAuthzs + 1
      }
  }
  if numAuthzs >= 1 {
      return true
  }
  return false
}

main = rule {
   control_group()
}
```

Deploy this policy as an Endpoint Governing Policy attached to
"**`EU_GDPR_data/orders/*`**" path.

-> Refer to the [Sentinel
Properties](/docs/enterprise/sentinel/properties.html#control-group-properties)
documentation for the list of available properties associated with control
groups.


## Next steps

To protect your secrets, it may become necessary to write finer-grained
policies to introspect different aspects of incoming requests.  If you have not
already done so, read [Sentinel](https://docs.hashicorp.com/sentinel/)
documentation to learn more about what you can accomplish writing policies as a
code.
