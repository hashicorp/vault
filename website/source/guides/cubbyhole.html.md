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
environment.  You can follow the [Getting Started][getting-started] guide to
[install Vault][install-vault]. Alternatively, if you are familiar with
[Vagrant](https://www.vagrantup.com/), you can spin up a
[HashiStack](https://github.com/hashicorp/vault-guides/tree/master/provision/hashistack/vagrant)
virtual machine.

Make sure that your Vault server has been [initialized and unsealed][initialize].

[getting-started]: /intro/getting-started/install.html
[install-vault]: /intro/getting-started/install.html
[initialize]: /intro/getting-started/deploy.html

## Steps

To distribute the initial token to an app using cubbyhole response wrapping, you
perform the following tasks:

1. Create and wrap a token
2. Unwrap the secret

### Step 1: Create and wrap a token

When the response to `vault token-create` request is wrapped, Vault inserts the
generated token it into the cubbyhole of a single-use token, returning that
single-use wrapping token. Retrieving the secret requires an unwrap operation
against this wrapping token.

#### CLI command

```shell
vault token-create -policy=<POLICY_NAME> -wrap-ttl=<WRAP_TTL>
```

Where the `<WRAP_TTL>` is a numeric string indicating the TTL of the response.

**Example:**

```shell
vault token-create -policy=app-policy -wrap-ttl=60s

Key                          	Value
---                          	-----
wrapping_token:              	9ac59985-094f-a2de-aed8-bf688e436fbc
wrapping_token_ttl:          	1m0s
wrapping_token_creation_time:	2018-01-10 00:47:54.970185208 +0000 UTC
wrapping_token_creation_path:	auth/token/create
wrapped_accessor:            	195763a9-3f26-1fcf-6a1a-ee0a11e76cb1
```

#### API call using cURL

Response wrapping is per-request and is triggered by providing to Vault the
desired TTL for a response-wrapping token for that request. This is set using
the **`X-Vault-Wrap-TTL`** header in the request and can be either an integer
number of seconds or a string duration.

**Example:**

```text
curl -X POST -H "X-Vault-Token: $VAULT_TOKEN" -H "X-Vault-Wrap-TTL: 60s" \
  -d '{"policies":["app-policy"]}' $VAULT_ADDR/v1/auth/token/create | jq

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

### Step 2: Unwrap the secret

The client uses the wrapping token to unwrap the secret.

**NOTE:**
If a client has been expecting delivery of a response-wrapping token and none
arrives, this may be due to an attacker intercepting the token and then
preventing it from traveling further. This should cause an alert to trigger an
immediate investigation.

#### CLI command

```text
vault unwrap <WRAPPING_TOKEN>
```
Or

```text
VAULT_TOKEN=<WRAPPING_TOKEN> vault unwrap
```

**Example:**

```shell
$ vault unwrap 9ac59985-094f-a2de-aed8-bf688e436fbc

Key            	Value
---            	-----
token          	7bb915b2-8a44-48b0-a71d-72b590252016
token_accessor 	195763a9-3f26-1fcf-6a1a-ee0a11e76cb1
token_duration 	768h0m0s
token_renewable	true
token_policies 	[app-policy default]
```

#### API call using cURL

To enable the AppRole auth backend via API:

```text
curl -X POST -H "X-Vault-Token: $WRAPPING_TOKEN" $VAULT_ADDR/v1/sys/wrapping/unwrap | jq

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
      "app-policy",
      "default"
    ],
    "metadata": null,
    "lease_duration": 2764800,
    "renewable": true
  }
}
```


## Reference Content

Similar to the key/value secret backend, the cubbyhole backend is mounted at the
**`cubbyhole/`** prefix by default. The secrets you store in the `cubbyhole/` path
are tied to your token and only accessible by you.

To test the cubbyhole secret backend, perform the following steps.

First, create `tester` policy which grants permissions on the path under `cubbyhole/private/` prefix.  

```text
$ vault policy-write tester tester.hcl

$ cat tester.hcl
path "cubbyhole/private/*" {
	capabilities = ["create", "read", "update", "delete", "list"]
}
```

Create a token attached to the `tester` policy, and then authenticate using the
token.

```text
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

```text
$ vault write cubbyhole/private/access-token token="123456789abcdefg87654321"
Success! Data written to: cubbyhole/private/access-token

$ vault read cubbyhole/private/access-token
Key  	Value
---  	-----
token	123456789abcdefg87654321
```

Now, try to access the secret using the root token, you shouldn't be able to
read.

```text
VAULT_TOKEN=<ROOT_TOKEN> vault read cubbyhole/private/access-token

No value found at cubbyhole/private/access-token
```

Also, refer to [Cubbyhole Secret Backend HTTP API](/api/secret/cubbyhole/index.html).


## Next steps
