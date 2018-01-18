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

## Challenge

Think of a scenario where a DevOps team wants to configure Jenkins to read
secrets from Vault so that it can inject the secrets to app's environment
variables at deployment time. For example, Jenkins configures
`WORDPRESS_DB_HOST`, `WORDPRESS_DB_USER`, `WORDPRESS_DB_PASSWORD`, and
`WORDPRESS_DB_NAME` environment variables when deploying WordPress.

To retrieve database connection information from Vault, Jenkins must
authenticate with Vault first.   

![Vault communication](/assets/images/vault-approle.png)

As a user, you can authenticate with Vault using your LDAP credentials, and
Vault would generate a token once it's verified. How can an app programmatically
request a token so that it can interact with Vault?  


## Solution

Enable **AppRole** auth backend so that the Jenkins
server can obtain a Vault token with appropriate policies attached.

## Prerequisites

To perform the tasks described in this guide, you need to have a Vault
environment.  Refer to the [Getting
Started](/intro/getting-started/install.html) guide to install Vault.

Make sure that your Vault server has been [initialized and
unsealed](/intro/getting-started/deploy.html).

Complete the [policies](/guides/policies.html) guide so that you have policies
to work with.


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


### <a name="step1"></a>Step 1: Enable AppRole auth backend

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

When you enabled AppRole auth backend, it gets mounted at the
**`/auth/approle`** path. In this example, you are going to create a role for
Jenkins server.

-> Policies created in the [policies](/guides/policies.html) guide are
referenced in this step.

#### CLI command

To create a role for machines or apps, run:

```plaintext
vault write auth/approle/role/<ROLE_NAME> [args]
```

There are a number of [parameters](/api/auth/approle/index.html) that you can
set. Most importantly, you want to set policies for the role. You are going to
create a role operates in **pull** mode in this example.

**Example:**

The following example creates a role named `jenkins` with `dev-pol` and
`devops-pol` policies attached.

```plaintext
$ vault write auth/approle/role/jenkins policies="dev-pol,devops-pol"

$ vault read auth/approle/role/jenkins

  Key               	Value
  ---               	-----
  bind_secret_id    	true
  bound_cidr_list
  period            	0
  policies          	[dev-pol devops-pol]
  secret_id_num_uses	0
  secret_id_ttl     	0
  token_max_ttl     	0
  token_num_uses    	0
  token_ttl         	0
```

**NOTE:** If desired, the token use limits, TTL and settings can be specified
during the role creation.

#### API call using cURL

**Example:**

```text
curl -X POST -H "X-Vault-Token:$VAULT_TOKEN" -d '{"policies":"dev-pol,devops-pol"}' \
    $VAULT_ADDR/v1/auth/approle/role/jenkins
```


### <a name="step3"></a>Step 3: Get Role ID and Secret ID

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

To get a Vault token for the `jenkins` role, you need to pass the role ID and
secret ID you obtained previously to the client (in this case, Jenkins).

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
  token_policies      	[dev-pol devops-pol]
  token_meta_role_name	"jenkins"
```

#### API call using cURL

Notice that the following API call passes the role ID and secret ID in the
request payload. Upon a successful login, a token will be returned.

**Example:**

```plaintext
$ cat jenkins.json
  {
    "role_id": "675a50e7-cfe0-be76-e35f-49ec009731ea",
    "secret_id": "ed0a642f-2acf-c2da-232f-1b21300d5f29"
  }

$ curl -X POST -d @jenkins.json $VAULT_ADDR/v1/auth/approle/login | jq
```

-> Once the client acquired the token, the future requests can be made using
that token.


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
