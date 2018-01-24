---
layout: "guides"
page_title: "Cubbyhole Response Wrapping - Guides"
sidebar_current: "guides-cubbyhole"
description: |-
  Vault provides a capability to wrap Vault response and store it in a
  "cubbyhole" where the holder of the one-time use wrapping token can unwrap to
  uncover the secret.
---

# Cubbyhole

The term _cubbyhole_ comes from an Americnaism where you get a "locker" or "safe
place" to store your belongings or valuables. In Vault, cubbyhole is your
"locker".  All secrets are namespaced under **your token**. If that token
expires or is revoked, all the secrets in its cubbyhole are revoked as well.

It is not possible to reach into another token's cubbyhole even as the root
user. This is the key difference between the cubbyhole and the key/value secret
backend. The secrets in the key/value backends are accessible to any token for as
long as its policy allows it.


## Reference Material

- [Cubbyhole](/docs/secrets/cubbyhole/index.html)
- [Response Wrapping](/docs/concepts/response-wrapping.html)

## Estimated Time to Complete

10 minutes

## Personas

The end-to-end scenario described in this guide involves two personas:

- **`admin`** with privileged permissions to create tokens and policies
- **`app`** is the receiving client of a wrapped response

## Challenge

In order to tightly manage the secrets, you set the scope of who can do what
using the [Vault policy](/docs/concepts/policies.html) and attach that to
tokens, roles, entities, etc.

How to securely distribute the initial token to a machine or app?

## Solution

Use Vault's **cubbyhole response wrapping** where the initial token is stored in
the cubbyhole backend. The wrapped secret can be unwrap using the single use
wrapping token. Even the user or the system created the initial token won't see
the original value. The wrapping token is short-lived and can be revoked just
like any other tokens so that the risk of unauthorized access can be minimized.

## Prerequisites

To perform the tasks described in this guide, you need to have a Vault
environment.  Refer to the [Getting
Started](/intro/getting-started/install.html) guide to install Vault. Make sure
that your Vault server has been [initialized and
unsealed](/intro/getting-started/deploy.html).

### <a name="policy"></a>Policy requirements

To perform all tasks demonstrated in this guide, you need to be able to
authenticate with Vault as an [**`admin`** user](#personas).

-> **NOTE:** For the purpose of this guide, you can use **`root`** token to work
with Vault. However, it is recommended that root tokens are only used for just
enough initial setup or in emergencies. As a best practice, use tokens with
appropriate set of policies based on your role in the organization.

The `admin` policy must include the following permissions:

```shell
# Manage tokens
path "sys/auth/token/*" {
  capabilities = [ "create", "read", "update", "delete" ]
}

# Write ACL policies
path "sys/policy/*" {
  capabilities = [ "create", "read", "update", "delete", "list" ]
}

# Access cubbyhole backend
path "cubbyhole/private/*" {
	capabilities = ["create", "read", "update", "delete", "list"]
}
```

## Steps

To distribute the initial token to an app using cubbyhole response wrapping, you
perform the following tasks:

1. [Create and wrap a token](#step1)
2. [Unwrap the secret](#step2)

### <a name="step1"></a>Step 1: Create and wrap a token
(**Persona:** admin)

When the response to `vault token-create` request is wrapped, Vault inserts the
generated token it into the cubbyhole of a single-use token, returning that
single-use wrapping token. Retrieving the secret requires an unwrap operation
against this wrapping token.

In this scenario, an [admin user](#personas) creates a token using response wrapping. To perform the steps in this guide, first create a policy for the app.

`app-policy.hcl`:

```shell
# Unwrap the token
path "sys/wrapping/unwrap" {
  capabilities = [ "create", "read" ]
}

# For testing, read-only on secret/dev path
path "secret/dev" {
  capabilities = [ "read" ]
}
```

#### CLI command

First create an `apps` policy:

```shell
$ vault policy-write apps apps-policy.hcl
Policy 'apps' written.
```

To create a token using response wrapping:

```shell
$ vault token-create -policy=<POLICY_NAME> -wrap-ttl=<WRAP_TTL>
```

Where the `<WRAP_TTL>` is a numeric string indicating the TTL of the response.

**Example:**

Generate a token for `app` persona using response wrapping with TTL of 60
seconds.

```shell
$ vault token-create -policy=apps -wrap-ttl=60s

Key                          	Value
---                          	-----
wrapping_token:              	9ac59985-094f-a2de-aed8-bf688e436fbc
wrapping_token_ttl:          	1m0s
wrapping_token_creation_time:	2018-01-10 00:47:54.970185208 +0000 UTC
wrapping_token_creation_path:	auth/token/create
wrapped_accessor:            	195763a9-3f26-1fcf-6a1a-ee0a11e76cb1
```

The response is the wrapping token; therefore, the admin user does not even see
the generated token from the `token-create` command.

#### API call using cURL

First create an `apps` policy:

```shell
$ curl --header "X-Vault-Token: ..." --request PUT \
       --data @payload.json \
       https://vault.rocks/v1/sys/policy/apps

$ cat payload.json
{
  "policy": "path \"sys/wrapping/unwrap\" { capabilities = [ \"create\", \"read\" ] ... }"
}
```

Response wrapping is per-request and is triggered by providing to Vault the
desired TTL for a response-wrapping token for that request. This is set using
the **`X-Vault-Wrap-TTL`** header in the request and can be either an integer
number of seconds or a string duration.

```shell
$ curl --header "X-Vault-Wrap-TTL: <TTL>" \
       --header "X-Vault-Token: <TOKEN>" \
       --request <HTTP_VERB> \
       --data '<PARAMETERS>' \
       <VAULT_ADDRESS>/v1/<ENDPOINT>
```

Where `<TTL>` can be either an integer number of seconds or a string duration of
seconds (15s), minutes (20m), or hours (25h).

**Example:**

To wrap the response of token-create request:

```shell
$ curl --header "X-Vault-Wrap-TTL: 60s" \
       --header "X-Vault-Token: ..." \
       --request POST \
       --data '{"policies":["apps"]}' \
       https://vault.rocks/v1/auth/token/create | jq
{
  "request_id": "",
  "lease_id": "",
  "renewable": false,
  "lease_duration": 0,
  "data": null,
  "wrap_info": {
    "token": "e095129f-123a-4fef-c007-1f6a487cfa78",
    "ttl": 60,
    "creation_time": "2018-01-10T01:43:38.025351336Z",
    "creation_path": "auth/token/create",
    "wrapped_accessor": "44e8253c-65b4-1690-1bf1-7902a7a6b2aa"
  },
  "warnings": null,
  "auth": null
}  
```

Generate a token for `app` persona using response wrapping with TTL of 60
seconds. The admin user does not even see the generated token.


### <a name="step2"></a>Step 2: Unwrap the secret
(**Persona:** app)

The response is wrapped by a wrapping token, and retrieving it requires an
unwrap operation against this token.

-> **NOTE:** If a client has been expecting delivery of a response-wrapping
token and none arrives, this may be due to an attacker intercepting the token
and then preventing it from traveling further. This should cause an alert to
trigger an immediate investigation.

#### CLI command

To unwrap the secret:

```shell
$ vault unwrap <WRAPPING_TOKEN>
```
Or

```shell
$ VAULT_TOKEN=<WRAPPING_TOKEN> vault unwrap
```

In this scenario, the wrapped secret is a Vault token. Therefore, it probably
makes better sense to use the second option.

**Example:**

```shell
$ VAULT_TOKEN=9ac59985-094f-a2de-aed8-bf688e436fbc vault unwrap

Key            	Value
---            	-----
token          	7bb915b2-8a44-48b0-a71d-72b590252016
token_accessor 	195763a9-3f26-1fcf-6a1a-ee0a11e76cb1
token_duration 	768h0m0s
token_renewable	true
token_policies 	[apps default]
```

Once the client acquired the token, future requests can be made using this
token.

```shell
$ vault auth 7bb915b2-8a44-48b0-a71d-72b590252016

$ vault read secret/dev
```

#### API call using cURL

To unwrap the secret, use `/sys/wrapping/unwrap` endpoint:

```shell
$ curl --header "X-Vault-Token: <WRAPPING_TOKEN>" \
       --request POST \
       <VAULT_ADDRESS>/v1/sys/wrapping/unwrap
```

**Example:**

```shell
$ curl --header "X-Vault-Token: e095129f-123a-4fef-c007-1f6a487cfa78" \
       --request POST \
       https://vault.rocks/v1/sys/wrapping/unwrap | jq
{
  "request_id": "d704435d-c1cf-b8a3-52f6-ec50bc8246c4",
  "lease_id": "",
  "renewable": false,
  "lease_duration": 0,
  "data": null,
  "wrap_info": null,
  "warnings": null,
  "auth": {
    "client_token": "af5f7682-aa55-fa37-5039-ee116df56600",
    "accessor": "19b5407e-b304-7cde-e946-54942325d3c1",
    "policies": [
      "apps",
      "default"
    ],
    "metadata": null,
    "lease_duration": 2764800,
    "renewable": true
  }
}
```

Once the client acquired the token, future requests can be made using this
token.

```plaintext
$ curl --header "X-Vault-Token: af5f7682-aa55-fa37-5039-ee116df56600" \
       --request GET \
       https://vault.rocks/v1/secret/dev | jq
{
  "errors": []
}
```

## Additional Discussion

Similar to the key/value secret backend, the cubbyhole backend is mounted at the
**`cubbyhole/`** prefix by default. The secrets you store in the `cubbyhole/` path
are tied to your token and only accessible by you.

To test the cubbyhole secret backend, perform the following steps.

First, create `tester` policy which grants permissions on the path under `cubbyhole/private/` prefix.  

```shell
$ vault policy-write tester tester.hcl

$ cat tester.hcl
path "cubbyhole/private/*" {
	capabilities = ["create", "read", "update", "delete", "list"]
}
```

Create a token attached to the `tester` policy, and then authenticate using the
token.

```shell
$ vault token-create -policy=tester
Key            	Value
---            	-----
token          	2ba26888-b531-1626-3598-01ea4aa383bb
token_accessor 	28cbd05c-31a3-0aaa-4dca-838a9aafe4cb
token_duration 	768h0m0s
token_renewable	true
token_policies 	[default tester]

$ unset VAULT_TOKEN

$ vault auth 2ba26888-b531-1626-3598-01ea4aa383bb
Successfully authenticated! You are now logged in.
token: 2ba26888-b531-1626-3598-01ea4aa383bb
token_duration: 2764651
token_policies: [default tester]
```

You should be able to write secrets under `cubbyhole/private/` path, and read it
back.

```shell
$ vault write cubbyhole/private/access-token token="123456789abcdefg87654321"
Success! Data written to: cubbyhole/private/access-token

$ vault read cubbyhole/private/access-token
Key  	Value
---  	-----
token	123456789abcdefg87654321
```

Now, try to access the secret using the root token, you shouldn't be able to
read.

```shell
$ VAULT_TOKEN=<ROOT_TOKEN> vault read cubbyhole/private/access-token

No value found at cubbyhole/private/access-token
```

Also, refer to [Cubbyhole Secret Backend HTTP API](/api/secret/cubbyhole/index.html).


## Next steps
