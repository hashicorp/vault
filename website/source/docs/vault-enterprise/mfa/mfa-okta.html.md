---
layout: "docs"
page_title: "Vault Enterprise Okta MFA"
sidebar_current: "docs-vault-enterprise-mfa-okta"
description: |-
  Vault Enterprise supports Okta MFA type.
---

# Okta MFA

This page demonstrates the Okta MFA on ACL'd paths of Vault.

## Steps

### Enable Auth Backend

```
vault auth-enable userpass
```

### Fetch Mount Accessor

```
vault auth -methods
```

```
Path       Type      Accessor                Default TTL  Max TTL  Replication Behavior  Description
...
userpass/  userpass  auth_userpass_54b8e339  system       system   replicated
```


### Configure Okta MFA method

```
vault write sys/mfa/method/okta/my_okta mount_accessor=auth_userpass_54b8e339 org_name="dev-262775" api_token="0071u8PrReNkzmATGJAP2oDyIXwwveqx9vIOEyCZDC"
```

### Create Policy

Create a policy that gives access to secret through the MFA method created
above.

#### Sample Payload

```hcl
path "secret/foo" {
    capabilities = ["read"]
    mfa_methods = ["my_okta"]
}
```

```
vault policy-write okta-policy payload.hcl
```

### Create User

MFA works only for tokens that have identity information on them. Tokens
created by logging in using authentication backends will have the associated
identity information. Let's create a user in the `userpass` backend and
authenticate against it.


```
vault write auth/userpass/users/testuser password=testpassword policies=okta-policy
```

### Create Login Token

```
vault write auth/userpass/login/testuser password=testpassword
```

```
Key                     Value
---                     -----
token                   70f97438-e174-c03c-40fe-6bcdc1028d6c
token_accessor          a91d97f4-1c7d-6af3-e4bf-971f74f9fab9
token_duration          768h0m0s
token_renewable         true
token_policies          [default okta-policy]
token_meta_username     "testuser"
```

Note that the CLI is not authenticated with the newly created token yet, we did
not call `vault auth`, instead we used the login API to simply return a token.

### Fetch Entity ID From Token

Caller identity is represented by the `entity_id` property of the token.

```
vault token-lookup 70f97438-e174-c03c-40fe-6bcdc1028d6c
```

```
Key                     Value
---                     -----
accessor                a91d97f4-1c7d-6af3-e4bf-971f74f9fab9
creation_time           1502245243
creation_ttl            2764800
display_name            userpass-testuser
entity_id               307d6c16-6f5c-4ae7-46a9-2d153ffcbc63
expire_time             2017-09-09T22:20:43.448543132-04:00
explicit_max_ttl        0
id                      70f97438-e174-c03c-40fe-6bcdc1028d6c
issue_time              2017-08-08T22:20:43.448543003-04:00
meta                    map[username:testuser]
num_uses                0
orphan                  true
path                    auth/userpass/login/testuser
policies                [default okta-policy]
renewable               true
ttl                     2764623
```

### Login

Authenticate the CLI to use the newly created token.

```
vault auth 70f97438-e174-c03c-40fe-6bcdc1028d6c
```

### Read Secret

Reading the secret will trigger an Okta push. This will be a blocking call until
the push notification is either approved or declined.

```
vault read secret/foo
```

```
Key                     Value
---                     -----
refresh_interval        768h0m0s
data                    which can only be read after MFA validation
```
