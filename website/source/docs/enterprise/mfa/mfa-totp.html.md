---
layout: "docs"
page_title: "Vault Enterprise TOTP MFA"
sidebar_current: "docs-vault-enterprise-mfa-totp"
description: |-
  Vault Enterprise supports TOTP MFA type.
---

# TOTP MFA

This page demonstrates the TOTP MFA on ACL'd paths of Vault.

## Steps

### Configure TOTP MFA method

```
vault write sys/mfa/method/totp/my_totp issuer=Vault period=30 key_size=30 algorithm=SHA256 digits=6
```

### Create Secret

Create a secret to be accessed after validating MFA.

```
vault write secret/foo data="which can only be read after MFA validation"
```

### Create Policy

Create a policy that gives access to secret through the MFA method created
above.

#### Sample Payload

```hcl
path "secret/foo" {
    capabilities = ["read"]
    mfa_methods = ["my_totp"]
}
```

```
vault policy-write totp-policy payload.hcl
```

### Enable Auth Backend

MFA works only for tokens that have identity information on them. Tokens
created by logging in using authentication backends will have the associated
identity information. Let's create a user in the `userpass` backend and
authenticate against it.

```
vault auth-enable userpass
```

### Create User

```
vault write auth/userpass/users/testuser password=testpassword policies=totp-policy
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
token_policies          [default totp-policy]
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
policies                [default totp-policy]
renewable               true
ttl                     2764623
```

### Generate TOTP Method Secret on Entity

Let's generate a TOTP key using the `my_totp` configuration and store it in the
entity of the user. A barcode and a URL for the secret key will be returned by
the API. This should be distributed to the intended user to be able to generate
TOTP passcode.

```
vault write sys/mfa/method/totp/my_totp/admin-generate entity_id=307d6c16-6f5c-4ae7-46a9-2d153ffcbc63
```

```
Key     Value
---     -----
barcode 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
url     otpauth://totp/Vault:307d6c16-6f5c-4ae7-46a9-2d153ffcbc63?algorithm=SHA256&digits=6&issuer=Vault&period=30&secret=AQESPQUPHWYIXV7FGOMBYT3A2N4LQKEIRNKTSRCWTKVEW66L
```

Note that Vault's [TOTP secret backend](/docs/secrets/totp/index.html) can be leveraged to create TOTP passcodes.

### Login

Authenticate the CLI to use the newly created token.

```
vault auth 70f97438-e174-c03c-40fe-6bcdc1028d6c
```

### Read Secret

Read the secret by supplying the TOTP passcode.

```
vault read -mfa my_totp:146378 secret/foo
```

```
Key                     Value
---                     -----
refresh_interval        768h0m0s
data                    which can only be read after MFA validation
```
