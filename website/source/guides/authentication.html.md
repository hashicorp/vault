---
layout: "guides"
page_title: "AppRole Pull Authentication - Guides"
sidebar_current: "guides-authentication"
description: |-
  Authentication is a process in Vault by which user or machine-supplied
  information is verified to create a token with pre-configured policy.
---

# Authentication

Before a client can interact with Vault, it must authenticate against an [**auth
backend**](/docs/auth/index.html) to acquire a token. This token has policies attached so
that the behavior of the client can be governed.

Since tokens are the core method for authentication within Vault, there is a
**token** auth backend (often refer as **_token store_**). This is a special
auth backend responsible for creating and storing tokens.

### Auth Backends

Auth backends performs authentication to verify the user or machine-supplied
information. Some of the supported auth backends are targeted towards users
while others are targeted toward machines or apps. For example,
[**LDAP**](/docs/auth/ldap.html) auth backend enables user authentication using
an existing LDAP server while [**AppRole**](/docs/auth/approle.html) auth
backend is recommended for machines or apps.

[Getting Started](/intro/getting-started/authentication.html) guide walks you
through how to enable GitHub auth backend for user authentication.

This guide focuses on generating tokens for machines or apps by enabling
[**AppRole**](/docs/auth/approle.html) auth backend.


## Reference Material

- [Getting Started](/intro/getting-started/authentication.html)
- [Auth Backends](docs/auth/index.html)
- [GitHub Auth APIs](/api/auth/github/index.html)


## Estimated Time to Complete

10 minutes

## Personas

The end-to-end scenario described in this guide involves two personas:

- **`admin`** with privileged permissions to configure an auth backend
- **`app`** is the consumer of secrets stored in Vault

## Challenge

Think of a scenario where a DevOps team wants to configure Jenkins to read
secrets from Vault so that it can inject the secrets to app's environment
variables (e.g. `MYSQL_DB_HOST`) at deployment time.

Instead of hardcoding secrets in each build script as a plaintext, Jenkins
retrieves secrets from Vault.

![Vault communication](/assets/images/vault-approle.png)

As a user, you can authenticate with Vault using your LDAP credentials, and
Vault generates a token. This token has policies granting you to perform
appropriate operations.

How can a Jenkins server programmatically request a token so that it can read
secrets from Vault?  


## Solution

Enable **AppRole** auth backend so that the Jenkins server can obtain a Vault
token with appropriate policies attached. Since each AppRole has attached
policies, you can write fine-grained policies limiting which app can access
which path.  


## Prerequisites

To perform the tasks described in this guide, you need to have a Vault
environment.  Refer to the [Getting
Started](/intro/getting-started/install.html) guide to install Vault. Make sure
that your Vault server has been [initialized and
unsealed](/intro/getting-started/deploy.html).

### <a name="policy"></a>Policy requirements

To perform all tasks demonstrated in this guide, you need to be able to
authenticate with Vault as an **`admin`** user. The `admin` user's policy must
include the following permissions:

```shell
# Mount the AppRole auth backend
path "sys/auth/approle/*" {
  capabilities = [ "create", "read", "update", "delete" ]
}

# Create and manage roles
path "auth/approle/*" {
  capabilities = [ "create", "read", "update", "delete", "list" ]
}

# Write ACL policies
path "sys/policy/*" {
  capabilities = [ "create", "read", "update", "delete", "list" ]
}

# Write test data
path "secret/mysql/*" {
  capabilities = [ "create", "read", "update", "delete", "list" ]
}
```

-> **NOTE:** For the purpose of this guide, you can use **`root`** token. However,
Vault team recommends that root tokens are only used for just enough initial
setup or in emergencies. As a best practice, use tokens with appropriate
set of policies based on your role in the organization.

If you are not familiar with policies, complete the
[policies](/guides/policies.html) guide.


## Steps

Vault supports a number of authentication backends, and most auth backends must
be enabled first including AppRole.

The overall workflow is:
![AppRole auth backend workflow](assets/images/vault-approle-workflow.png)

In this guide, you are going to perform the following steps:

1. [Enable AppRole auth backend](#step1)
2. [Create a role with policy attached](#step2)
3. [Get Role ID and Secret ID](#step3)
4. [Login with Role ID & Secret ID](#step4)
5. [Read secrets using the AppRole token](#step5)

Step 1 through 3 need to be performed by an `admin` user.  Step 4 and 5 describe
the commands that an `app` runs to get a token and read secrets from Vault.


### <a name="step1"></a>Step 1: Enable AppRole auth backend
(**Persona:** admin)

[AppRole](/docs/auth/approle.html) is an authentication mechanism within Vault
to allow machines or apps to acquire a token to interact with Vault. It uses
**Role ID** and **Secret ID** for login.

#### CLI command

```shell
vault auth-enable approle
```

#### API call using cURL

Before begin, create the following environment variables for your convenience:

- **VAULT_ADDR** is set to your Vault server address
- **VAULT_TOKEN** is set to your Vault token

**Example:**

```plaintext
$ export VAULT_ADDR=http://127.0.0.1:8201

$ export VAULT_TOKEN=0c4d13ba-9f5b-475e-faf2-8f39b28263a5
```

Now, enable the AppRole auth backend via API:

```text
curl -X POST -H "X-Vault-Token: $VAULT_TOKEN" --data '{"type": "approle"}' \
    $VAULT_ADDR/v1/sys/auth/approle
```


### <a name="step2"></a>Step 2: Create a role with policy attached
(**Persona:** admin)

When you enabled AppRole auth backend, it gets mounted at the
**`/auth/approle`** path. In this example, you are going to create a role for
Jenkins server.

The scenario in this guide requires the **`app`** persona (Jenkins) to have the
following policy (`jenkins-pol.hcl`):

```shell
# Login with AppRole
path "auth/approle/login" {
  capabilities = [ "create", "read" ]
}

# Read test data
path "secret/mysql/*" {
  capabilities = [ "read" ]
}
```

#### CLI command

Before creating a role, create `jenkins` policy:

```shell
vault policy-write jenkins jenkins-pol.hcl
```

The command to create a new AppRole:

```plaintext
vault write auth/approle/role/<ROLE_NAME> [args]
```

There are a number of [parameters](/api/auth/approle/index.html) that you can
set. Most importantly, you want to specify the policies when you create a role.


**Example:**

The following example creates a role named `jenkins` with `jenkins` policy
attached. (NOTE: This example creates a role operates in [**pull**
mode](/docs/auth/approle.html).)

```plaintext
$ vault write auth/approle/role/jenkins policies="jenkins"

$ vault read auth/approle/role/jenkins

  Key               	Value
  ---               	-----
  bind_secret_id    	true
  bound_cidr_list
  period            	0
  policies          	[jenkins]
  secret_id_num_uses	0
  secret_id_ttl     	0
  token_max_ttl     	0
  token_num_uses    	0
  token_ttl         	0
```

**NOTE:** To attach multiple policies, pass the policy names as a comma
separated string.

```shell
vault write auth/approle/role/jenkins policies="jenkins,anotherpolicy"
````

#### API call using cURL

Before creating a role, create `jenkins` policy:

```text
$ curl -X PUT -H "X-Vault-Token: $VAULT_TOKEN" -d @payload.json \
    $VAULT_ADDR/v1/sys/policy/jenkins

$ cat payload.json
{
  "policy": "path \"auth/approle/login\" {  capabilities = [ \"create\", \"read\" ] } ... }"
}
```

There are a number of [parameters](/api/auth/approle/index.html) that you can
set. Most importantly, you want to specify the policies when you create a role.

You are going to create a role operates in [**pull** mode](/docs/auth/approle.html)
in this example.

**Example:**

The following example creates a role named `jenkins` with `jenkins` policy
attached. (NOTE: This example creates a role operates in [**pull**
mode](/docs/auth/approle.html).)

```text
$ curl -X POST -H "X-Vault-Token:$VAULT_TOKEN" -d '{"policies":"jenkins"}' \
    $VAULT_ADDR/v1/auth/approle/role/jenkins

$ curl -X GET -H "X-Vault-Token: $VAULT_TOKEN" $VAULT_ADDR/v1/auth/approle/role/jenkins | jq

{
  "request_id": "b18054ad-1ab5-8d83-eeed-193d97026ee7",
  "lease_id": "",
  "renewable": false,
  "lease_duration": 0,
  "data": {
    "bind_secret_id": true,
    "bound_cidr_list": "",
    "period": 0,
    "policies": [
      "jenkins"
    ],
    "secret_id_num_uses": 0,
    "secret_id_ttl": 0,
    "token_max_ttl": 0,
    "token_num_uses": 0,
    "token_ttl": 0
  },
  "wrap_info": null,
  "warnings": null,
  "auth": null
}
```

**NOTE:** To attach multiple policies, pass the policy names as a comma
separated string.

```shell
$ curl -X POST -H "X-Vault-Token:$VAULT_TOKEN" -d '{"policies":"jenkins,anotherpolicy"}' \
    $VAULT_ADDR/v1/auth/approle/role/jenkins
````

### <a name="step3"></a>Step 3: Get Role ID and Secret ID
(**Persona:** admin)

Since the example created a `jenkins` role which operates in pull mode, Vault
will generate the Secret ID. Similarly to tokens, you can set properties such as
usage-limit, TTLs, and expirations on the secret IDs.

#### CLI command

Now, you need to fetch the Role ID and Secret ID of a role.

**Example:**

```shell
$ vault read auth/approle/role/jenkins/role-id
  Key    	Value
  ---    	-----
  role_id	675a50e7-cfe0-be76-e35f-49ec009731ea

$ vault write -f auth/approle/role/jenkins/secret-id
  Key               	Value
  ---               	-----
  secret_id         	ed0a642f-2acf-c2da-232f-1b21300d5f29
  secret_id_accessor	a240a31f-270a-4765-64bd-94ba1f65703c
```

To list existing Secret IDs after creation:

```text
vault list auth/approle/role/jenkins/secret-id
```

#### API call using cURL

**Example:**

```text
$ curl -X GET -H "X-Vault-Token:$VAULT_TOKEN" $VAULT_ADDR/v1/auth/approle/role/jenkins/role-id | jq

$ curl -X POST -H "X-Vault-Token:$VAULT_TOKEN" $VAULT_ADDR/v1/auth/approle/role/jenkins/secret-id | jq
```

To list existing secret IDs after creation:

```text
curl -X LIST -H "X-Vault-Token: $VAULT_TOKEN" $VAULT_ADDR/v1/auth/approle/role/jenkins/secret-id | jq
```


### <a name="step4"></a>Step 4: Login with Role ID & Secret ID
(**Persona:** app)

To get a Vault token for the `jenkins` role, you need to pass the role ID and
secret ID you obtained previously to the client (in this case, Jenkins).

-> Refer to the [Advanced Features](#advanced-features) section for securely
distributing the role ID and secret ID to the client app.

#### CLI command

**Example:**

```text
vault write auth/approle/login role_id="675a50e7-cfe0-be76-e35f-49ec009731ea" \
  secret_id="ed0a642f-2acf-c2da-232f-1b21300d5f29"

  Key                 	Value
  ---                 	-----
  token               	eeaf890e-4b0f-a687-4190-c75b1d6d70bc
  token_accessor      	fcee5d4e-7281-8bb0-2901-e743c52e0502
  token_duration      	768h0m0s
  token_renewable     	true
  token_policies      	[jenkins]
  token_meta_role_name	"jenkins"
```

Now you have a client token with `default` and `jenkins` policies attached.


#### API call using cURL

Notice that the following API call passes the role ID and secret ID in the
request payload.

**Example:**

```plaintext
$ cat jenkins.json
  {
    "role_id": "675a50e7-cfe0-be76-e35f-49ec009731ea",
    "secret_id": "ed0a642f-2acf-c2da-232f-1b21300d5f29"
  }

$ curl -X POST -d @jenkins.json $VAULT_ADDR/v1/auth/approle/login | jq

{
  "request_id": "fccae32b-1e6a-9a9c-7666-f5cb07805c1e",
  "lease_id": "",
  "renewable": false,
  "lease_duration": 0,
  "data": null,
  "wrap_info": null,
  "warnings": null,
  "auth": {
    "client_token": "3e7dd0ac-8b3e-8f88-bb37-a2890455ca6e",
    "accessor": "375c077e-bf02-a09b-c864-63d7f967e86b",
    "policies": [
      "default",
      "jenkins"
    ],
    "metadata": {
      "role_name": "jenkins"
    },
    "lease_duration": 2764800,
    "renewable": true,
    "entity_id": "54e0b765-6daf-0ff5-70b9-32c0d491f473"
  }
}
```

Now you have a client token with `default` and `jenkins` policies attached.


### <a name="step5"></a>Step 5: Read secrets using the AppRole token
(**Persona:** app)

Once the client acquired the token, the future requests can be made using
that token.

#### CLI command

**Example:**

You can pass the `client_token` returned in [Step 4](#step4) as a part of the
CLI command.

```shell
VAULT_TOKEN=3e7dd0ac-8b3e-8f88-bb37-a2890455ca6e vault read secret/mysql/webapp
No value found at secret/mysql/webapp
```

Alternatively, you can first authenticate with Vault using the `client_token`.

```shell
$ vault auth 3e7dd0ac-8b3e-8f88-bb37-a2890455ca6e
Successfully authenticated! You are now logged in.
token: 3e7dd0ac-8b3e-8f88-bb37-a2890455ca6e
token_duration: 2762013
token_policies: [default jenkins]

$ vault read secret/mysql/webapp
No value found at secret/mysql/webapp
```

Since there is no value in the `secret/mysql/webapp`, it returne "no value
found" message.

**Optional:** Using the `admin` user's token, you can store some secrets in the
`secret/mysql/webapp` backend.

```shell
$ vault write secret/dev/config/mongodb @mysqldb.txt

$ cat mysqldb.txt
{
  "url": "foo.example.com:35533",
  "db_name": "users",
  "username": "admin",
  "password": "pa$$w0rd"
}
```

Now, try to read secrets from `secret/mysql/webapp` using `client_token` again.
This time, it should return the values you just created.


#### API call using cURL

You can now pass the `client_token` returned in [Step 4](#step4) in the
**`X-Vault-Token`** header.

**Example:**

```plaintext
curl -X GET -H "X-Vault-Token: 3e7dd0ac-8b3e-8f88-bb37-a2890455ca6e" \
    $VAULT_ADDR/v1/secret/mysql/webapp | jq

{
  "errors": []
}
```

Since there is no value in the `secret/mysql/webapp`, it returns an empty array.

**Optional:** Using the `admin` user's token, you can store some secrets in the
`secret/mysql/webapp` backend.

```shell
$ curl -X POST -H "X-Vault-Token: $VAULT_TOKEN" --data @mysqldb.txt \
    $VAULT_ADDR/v1/secret/mysql/webapp

$ cat mysqldb.text
{
  "url": "foo.example.com:35533",
  "db_name": "users",
  "username": "admin",
  "password": "p@ssw0rd"
}
```

Now, try to read secrets from `secret/mysql/webapp` using `client_token` again.
This time, it should return the values you just created.




## Advanced Features

To keep the Secret ID confidential, use [**response
wrapping**](/docs/concepts/response-wrapping.html) so that the only expected
client can unwrap the Secret ID.

In the previous step, you executed the following command to retrieve the Secret
ID:

```text
vault write -f auth/approle/role/jenkins/secret-id
```

Instead, use response wrapping:

```text
vault write -wrap-ttl=60s -f auth/approle/role/jenkins/secret-id

Key                          	Value
---                          	-----
wrapping_token:              	9bbe23b7-5f8c-2aec-83dc-e97e94a2e632
wrapping_accessor:           	cb5bdc8f-0cdb-35ff-0e68-9de57a79c3bf
wrapping_token_ttl:          	1m0s
wrapping_token_creation_time:	2018-01-08 21:29:38.826611 -0800 PST
wrapping_token_creation_path:	auth/approle/role/jenkins/secret-id
```

The client app uses this `wrapping_token` to unwrap and obtain the Secret ID.

```text
VAULT_TOKEN=9bbe23b7-5f8c-2aec-83dc-e97e94a2e632 vault unwrap

Key               	Value
---               	-----
secret_id         	575f23e4-01ad-25f7-2661-9c9bdbb1cf81
secret_id_accessor	7d8a40b7-a6fd-a634-579b-b7d673ff86fb
```

To retrieve the Secret ID alone, you can use `jq` as follow:

```text
VAULT_TOKEN=2577044d-cf86-a065-e28f-e2a14ea6eaf7 vault unwrap -format=json | jq -r ".data.secret_id"

b07d7a47-1d0d-741d-20b4-ae0de7c6d964
```


## Next steps

To learn more about response wrapping, go to [Cubbyhole Response
Wrapping](/guides/cubbyhole.html) guide.
