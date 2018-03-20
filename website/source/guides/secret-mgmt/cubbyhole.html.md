---
layout: "guides"
page_title: "Cubbyhole Response Wrapping - Guides"
sidebar_current: "guides-secret-mgmt-cubbyhole"
description: |-
  Vault provides a capability to wrap Vault response and store it in a
  "cubbyhole" where the holder of the one-time use wrapping token can unwrap to
  uncover the secret.
---

# Cubbyhole

The term _cubbyhole_ comes from an Americanism where you get a "locker" or "safe
place" to store your belongings or valuables. In Vault, cubbyhole is your
"locker".  All secrets are namespaced under **your token**. If that token
expires or is revoked, all the secrets in its cubbyhole are revoked as well.

It is not possible to reach into another token's cubbyhole even as the root
user. This is the key difference between the cubbyhole and the key/value secret
engine. The secrets in the key/value secret engine are accessible to any token for as
long as its policy allows it.


## Reference Material

- [Cubbyhole](/docs/secrets/cubbyhole/index.html)
- [Response Wrapping](/docs/concepts/response-wrapping.html)

## Estimated Time to Complete

10 minutes

## Personas

The end-to-end scenario described in this guide involves two personas:

- **`admin`** with privileged permissions to create tokens
- **`apps`** trusted entity retrieving secrets from Vault

## Challenge

In order to tightly manage the secrets, you set the scope of who can do what
using the [Vault policy](/docs/concepts/policies.html) and attach that to
tokens, roles, entities, etc.

Think of a case where you have a trusted entity (Chef, Jenkins, etc.) which
reads secrets from Vault. This trusted entity must obtain a token. If the
trusted entity or its host machine was rebooted, it must re-authenticate with
Vault using a valid token.

How can you securely distribute the initial token to the trusted entity?

## Solution

Use Vault's **cubbyhole response wrapping** where the initial token is stored in
the cubbyhole secret engine. The wrapped secret can be unwrapped using the
single-use wrapping token. Even the user or the system created the initial token
won't see the original value. The wrapping token is short-lived and can be
revoked just like any other tokens so that the risk of unauthorized access can
be minimized.

## Prerequisites

To perform the tasks described in this guide, you need to have a Vault
environment.  Refer to the [Getting
Started](/intro/getting-started/install.html) guide to install Vault. Make sure
that your Vault server has been [initialized and
unsealed](/intro/getting-started/deploy.html).

### <a name="policy"></a>Policy requirements

-> **NOTE:** For the purpose of this guide, you can use **`root`** token to work
with Vault. However, it is recommended that root tokens are only used for just
enough initial setup or in emergencies. As a best practice, use tokens with
appropriate set of policies based on your role in the organization.

To perform all tasks demonstrated in this guide, your policy must include the
following permissions:

```shell
# Manage tokens
path "auth/token/*" {
  capabilities = [ "create", "read", "update", "delete", "sudo" ]
}

# Write ACL policies
path "sys/policy/*" {
  capabilities = [ "create", "read", "update", "delete", "list" ]
}

# Manage secret/dev secret engine - for Verification test
path "secret/dev" {
  capabilities = [ "create", "read", "update", "delete", "list" ]
}
```

If you are not familiar with policies, complete the
[policies](/guides/identity/policies.html) guide.

## Steps

Think of a scenario where apps read secrets from Vault. The `apps` need:

- Policy granting "read" permission on the specific path (`secret/dev`)
- Valid tokens to interact with Vault

![Response Wrapping Scenario](/assets/images/vault-cubbyhole.png)

Setting the appropriate policies and token generation are done by the `admin`
persona. For the `admin` to distribute the initial token to the app securely, it
uses cubbyhole response wrapping. In this guide, you perform the following:

1. [Create and wrap a token](#step1)
2. [Unwrap the secret](#step2)

**NOTE:** This guide demonstrates how the response wrapping works. To learn more
about reading and writing secrets in Vault, refer to the [Static
Secret](/guides/secret-mgmt/static-secrets.html) guide.

### <a name="step1"></a>Step 1: Create and wrap a token
(**Persona:** admin)

To solve the [challenge](#challenge) addressed in this guide:

1. More privileged token (`admin`) wraps a secret only the expecting client can
read
2. The receiving client (`app`) unwraps the secret to obtain the token

When the response to `vault token create` request is wrapped, Vault inserts the
generated token into the cubbyhole of a single-use token, returning that
single-use wrapping token. Retrieving the secret requires an unwrap operation
against this wrapping token.

In this scenario, an [admin user](#personas) creates a token using response
wrapping. To perform the steps in this guide, first create a policy for the app.

`apps-policy.hcl`:

```shell
# For testing, read-only on secret/dev path
path "secret/dev" {
  capabilities = [ "read" ]
}
```

#### CLI command

First create an `apps` policy:

```shell
$ vault policy write apps apps-policy.hcl
Policy 'apps' written.
```

To create a token using response wrapping:

```shell
$ vault token create -policy=<POLICY_NAME> -wrap-ttl=<WRAP_TTL>
```

Where the `<WRAP_TTL>` can be either an integer number of seconds or a string
duration of seconds (15s), minutes (20m), or hours (25h).

**Example:**

Generate a token for `apps` persona using response wrapping with TTL of 120
seconds.

```shell
$ vault token create -policy=apps -wrap-ttl=120

Key                          	Value
---                          	-----
wrapping_token:              	9ac59985-094f-a2de-aed8-bf688e436fbc
wrapping_token_ttl:          	2m0s
wrapping_token_creation_time:	2018-01-10 00:47:54.970185208 +0000 UTC
wrapping_token_creation_path:	auth/token/create
wrapped_accessor:            	195763a9-3f26-1fcf-6a1a-ee0a11e76cb1
```

The response is the wrapping token; therefore, the admin user does not even see
the generated token from the `token create` command.


#### API call using cURL

First create an `apps` policy using `sys/policy` endpoint:

```shell
$ curl --header "X-Vault-Token: <TOKEN>" \
       --request PUT \
       --data <PAYLOAD> \
       <VAULT_ADDRESS>/v1/sys/policy/<POLICY_NAME>
```

Where `<TOKEN>` is your valid token, and `<PAYLOAD>` includes policy name and
stringified policy.

**Example:**

```shell
# Request payload
$ cat payload.json
{
  "policy": "path \"secret/dev\" { capabilities = [ \"read\" ] }"
}

# API call to create a policy named, "apps"
$ curl --header "X-Vault-Token: ..." --request PUT --data @payload.json \
       https://vault.rocks/v1/sys/policy/apps
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

To wrap the response of token create request:

```shell
$ curl --header "X-Vault-Wrap-TTL: 120" \
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
    "ttl": 120,
    "creation_time": "2018-01-10T01:43:38.025351336Z",
    "creation_path": "auth/token/create",
    "wrapped_accessor": "44e8253c-65b4-1690-1bf1-7902a7a6b2aa"
  },
  "warnings": null,
  "auth": null
}  
```

This API call generates a token for `apps` persona using response wrapping with
TTL of 60 seconds. The admin user does not even see the generated token.


### <a name="step2"></a>Step 2: Unwrap the secret
(**Persona:** apps)

The `apps` persona receives a wrapping token from the `admin`.  In order for the
`apps` to acquire a valid token to read secrets from `secret/dev` path, it must
run the unwrap operation using this token.

-> **NOTE:** If a client has been expecting delivery of a response-wrapping
token and none arrives, this may be due to an attacker intercepting the token
and then preventing it from traveling further. This should cause an alert to
trigger an immediate investigation.

The following tasks will be performed to demonstrate the client operations:

1. Create a token with **`default`** policy
2. Authenticate with Vault using this `default` token (less privileged token)
3. Unwrap the secret to obtain more privileged token (**`apps`** persona token)
4. Verify that you can read `secret/dev` using the `apps`token


#### CLI command

First, create a token with `default` policy:

```shell
# Create a token with `default` policy
$ vault token create -policy=default
Key            	Value
---            	-----
token          	4522b2e8-27fe-bdc5-b932-d982f3166c6c
token_accessor 	96108f48-7475-6190-b058-769a2e5ebc8e
token_duration 	768h0m0s
token_renewable	true
token_policies 	[default]

# Authenticate using the generated token
$ vault login 4522b2e8-27fe-bdc5-b932-d982f3166c6c
Successfully authenticated! You are now logged in.
token: 4522b2e8-27fe-bdc5-b932-d982f3166c6c
token_duration: 2764729
token_policies: [default]

# Verify that you do NOT have a permission to read secret/dev
$ vault read secret/dev
Error reading secret/dev: Error making API request.

URL: GET http://<VAULT_ADDRESS>/v1/secret/dev
Code: 403. Errors:

* permission denied
```

The command to unwrap the wrapped secret is:

```shell
$ vault unwrap <WRAPPING_TOKEN>
```
Or

```shell
$ VAULT_TOKEN=<WRAPPING_TOKEN> vault unwrap
```

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

Verify that this token has `apps` policy attached.

Once the client acquired the token, future requests can be made using this
token.

```shell
$ vault login 7bb915b2-8a44-48b0-a71d-72b590252016

$ vault read secret/dev
No value found at secret/dev
```

#### API call using cURL

First, create a token with `default` policy:

```shell
# Create a new token default policy
$ curl --header "X-Vault-Token: ..." --request POST \
     --data '{"policies": "default"}' \
     https://vault.rocks/v1/auth/token/create | jq
{
  ...
  "auth": {
    "client_token": "5fe14760-b0fd-22dc-403c-14a05003b67f",
    "accessor": "e709610e-916e-f7e3-b93b-41f4dfdca7a0",
    "policies": [
      "default"
    ],
    ...
 }
}

# Verify that you can NOT read secret/dev using default token
$ curl --header "X-Vault-Token: 5fe14760-b0fd-22dc-403c-14a05003b67f" \
       --request GET \
       https://vault.rocks/v1/secret/dev | jq
{
 "errors": [
   "permission denied"
 ]
}
```

Now, unwrap the secret using `/sys/wrapping/unwrap` endpoint:

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

Since there is no data in `secret/dev`, it returns an empty array.

## Additional Discussion

The `cubbyhole` secret engine provides your own private secret storage space
where no one else can read (including `root`). This comes handy when you want to
store a password tied to your username that should not be shared with anyone.

The cubbyhole secret engine is mounted at the **`cubbyhole/`** prefix by
default. The secrets you store in the `cubbyhole/` path are tied to your token
and all tokens are permitted to read and write to the `cubbyhole` secret engine
by the [`default`](/docs/concepts/policies.html#default-policy) policy.

```shell
...
# Allow a token to manage its own cubbyhole
path "cubbyhole/*" {
    capabilities = ["create", "read", "update", "delete", "list"]
}
...
```

To test the cubbyhole secret engine, perform the following steps. (NOTE: Keep
using the `apps` token from [Step 2](#step2) to ensure that you are logged in with
non-root token.)

#### CLI command

Commands to write and read secrets to the `cubbyhole` secret engine:

```shell
# Write key-value pair(s) in your cubbyhole
$ vault write cubbyhole/<PATH> <KEY>=<VALUE>

# Read values from your cubbyhole
$ vault write cubbyhole/<PATH>
```

**Example:**

Write secrets under `cubbyhole/private/` path, and read it back.

```shell
# Write "token" to cubbyhole/private/access-token path
$ vault write cubbyhole/private/access-token token="123456789abcdefg87654321"
Success! Data written to: cubbyhole/private/access-token

# Read value from cubbyhole/private/access-token path
$ vault read cubbyhole/private/access-token
Key  	Value
---  	-----
token	123456789abcdefg87654321
```

Now, try to access the secret using the `root` token. It should NOT return the
secret.

```shell
$ VAULT_TOKEN=<ROOT_TOKEN> vault read cubbyhole/private/access-token

No value found at cubbyhole/private/access-token
```

#### API call using cURL

The API to work with the `cubbyhole` secret engine is very similar to `secret` secret engine:

```shell
$ curl --header "X-Vault-Token: <TOKEN>" \
       --request POST \
       --data <SECRETS> \
       <VAULT_ADDRESS>/v1/cubbyhole/<PATH>
```

**Example:**

Write secrets under `cubbyhole/private/` path, and read it back.

```shell
# Write "token" to cubbyhole/private/access-token path
$ curl --header "X-Vault-Token: e095129f-123a-4fef-c007-1f6a487cfa78" --request POST \
       --data '{"token": "123456789abcdefg87654321"}' \
       https://vault.rocks/v1/cubbyhole/private/access-token

# Read value from cubbyhole/private/access-token path
$ curl --header "X-Vault-Token: e095129f-123a-4fef-c007-1f6a487cfa78" --request GET \
       https://vault.rocks/v1/cubbyhole/private/access-token  | jq
{
 "request_id": "b2ff9f04-7a72-7eb0-672f-225b5eb652df",
 "lease_id": "",
 "renewable": false,
 "lease_duration": 0,
 "data": {
   "token": "123456789abcdefg87654321"
 },
 "wrap_info": null,
 "warnings": null,
 "auth": null
}
```

Now, try to access the secret using the `root` token. It should NOT return the
secret.

```shell
$ curl --header "X-Vault-Token: root" --request GET \
       https://vault.rocks/v1/cubbyhole/private/access-token  | jq
{
 "errors": []
}
```

Also, refer to [Cubbyhole Secret Engine (API)](/api/secret/cubbyhole/index.html).


## Next steps

The use of [AppRole Pull Authentication](/guides/identity/authentication.html) is a good
use case to leverage the response wrapping. Go through the guide if you have not
done so.  To better understand the lifecycle of Vault tokens, proceed to [Tokens
and Leases](/guides/identity/lease.html) guide.
