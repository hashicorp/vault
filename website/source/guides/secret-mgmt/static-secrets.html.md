---
layout: "guides"
page_title: "Static Secrets - Guides"
sidebar_title: "Static Secrets"
sidebar_current: "guides-secret-mgmt-static-secrets"
description: |-
  Vault supports generating new unseal keys as well as rotating the underlying
  encryption keys. This guide covers rekeying and rotating Vault's encryption
  keys.
---

# Static Secrets - Key/Value Secret Engine

Vault can be used to store any secret in a secure manner.  The secrets may be
SSL certificates and keys for your organization's domain, credentials to connect
to a corporate database server, etc. Storing such sensitive information in
plaintext is not desirable. This guide demonstrates the use case of Vault as a
Secret Storage.


## Reference Material

- [Key/Value Secret Engine](/docs/secrets/kv/index.html)
- [Key/Value Secret Engine API](/api/secret/kv/index.html)
- [Client libraries](/api/libraries.html) for Vault API for commonly used languages

## Estimated Time to Complete

10 minutes

## Personas

The end-to-end scenario described in this guide involves two personas:

- **`devops`** with privileged permissions to write secrets
- **`apps`** reads the secrets from Vault

## Challenge

Consider the following situations:

- Developers use a single admin account to access a third-party app
  (e.g. Splunk) and anyone who knows the user ID and password can log in as an
  admin
- SSH keys to connect to remote machines are shared and stored as a plaintext
- API keys to invoke external system APIs are stored as a plaintext
- An app integrates with LDAP, and its configuration information is in a
  plaintext

Organizations often seek an uniform workflow to securely store this sensitive
information.

## Solution

Vault as centralized secret storage to secure any sensitive information. Vault
encrypts these secrets using 256-bit AES in GCM mode with a randomly generated
nonce prior to writing them to its persistent storage. The storage backend never
sees the unencrypted value, so gaining access to the raw storage isn't enough to
access your secrets.

~> **NOTE:** This guide demonstrates secret management using [v2 of the KV
secret engine](/docs/secrets/kv/kv-v2.html).

## Prerequisites

To perform the tasks described in this guide, you need to have a Vault
environment.  Refer to the [Getting
Started](/intro/getting-started/install.html) guide to install Vault. Make sure
that your Vault server has been [initialized and
unsealed](/intro/getting-started/deploy.html).

### Policy requirements

-> **NOTE:** For the purpose of this guide, you can use **`root`** token to work
with Vault. However, it is recommended that root tokens are only used for just
enough initial setup or in emergencies. As a best practice, use tokens with
appropriate set of policies based on your role in the organization.

To perform all tasks demonstrated in this guide, your policy must include the
following permissions:

```shell
# Write and manage secrets in key/value secret engine
path "secret/*" {
  capabilities = [ "create", "read", "update", "delete", "list" ]
}

# Create policies to permit apps to read secrets
path "sys/policy/*" {
  capabilities = [ "create", "read", "update", "delete", "list" ]
}

# Create tokens for verification & test
path "auth/token/create" {
  capabilities = [ "create", "update", "sudo" ]
}
```

If you are not familiar with policies, complete the
[policies](/guides/identity/policies.html) guide.


## Steps

This guide demonstrates the basic steps to store secrets using Vault.  The
scenario here is to store the following secrets:

- API key (Google API)
- Root certificate of a production database (MySQL)

To store your API key within the configured physical storage for Vault, use the
**key/value** secret engine via **`secret/`** prefixed path.

-> Key/Value secret engine passes through any operation back to the configured
storage backend for Vault. For example, if your Vault server is configured with
Consul as its storage backend, a "read" operation turns into a read from Consul
at the same path.


You will perform the following:

1. [Enable KV Secret Engine v2](#step1)
1. [Store the Google API key](#step2)
1. [Store the root certificate for MySQL](#step3)
1. [Generate a token for apps](#step4)
1. [Retrieve the secrets](#step5)

![Personas Introduction](/img/vault-static-secrets.png)

Step 1 through 4 are performed by `devops` persona.  Step 5 describes the
commands that `apps` persona runs to read secrets from Vault.

### <a name="step1"></a>Step 1: Enable KV Secret Engine v2
(**Persona:** devops)

Currently, when you start the Vault server in [**dev
mode**](/intro/getting-started/dev-server.html#starting-the-dev-server), it
automatically enables **v2** of the KV secret engine at **`secret/`**. If you
start the Vault server in non-dev mode, the default is v1.

If you are running the server in **dev** mode, skip to [Step 2](#step2).
Otherwise, you must perform one of the following:

- Option 1: Upgrade the v1 of KV secret engine to v2
- Option 2: Enable the KV secret engine v2 at a different path


#### CLI command

Option 1: To upgrade from **v1** to **v2**:

```plaintext
$ vault kv enable-versioning secret/
```
<br>
Option 2: To enable the KV secret engine v2 at **`secret_v2/`**:

```plaintext
$ vault secrets enable -path=secret_v2/ kv-v2
```

Or

```plaintext
$ vault secrets enable -path=secret_v2/ -version=2 kv
```


#### API call using cURL

Option 1: To upgrade from **v1** to **v2**:

```plaintext
$ curl --header "X-Vault-Token: <TOKEN>" \
       --request POST \
       --data @payload.json \
       <VAULT_ADDRESS>/v1/sys/mounts/secret/tune
```

Where `<TOKEN>` is your valid token, and `<VAULT_ADDRESS>` is where your vault
server is running. The `payload.json` includes the version information.


**Example:**

```plaintext
$ cat payload.json
{
  "options": {
      "version": "2"
  }
}

$ curl --header "X-Vault-Token: ..." \
       --request POST \
       --data @payload.json \
       http://127.0.0.1:8200/v1/sys/mounts/secret/tune
```

<br>
Option 2: To enable the KV secret engine v2 at **`secret_v2/`**:


```plaintext
$ curl --header "X-Vault-Token: ..." \
       --request POST \
       --data '{"type":"kv-v2"}' \
       https://127.0.0.1:8200/v1/sys/mounts/secret_v2
```


<br>

~> **NOTE:** This guide assumes that you are working with KV secret engine
**v2** which is mounted at **`secret/`**.


### <a name="step2"></a>Step 2: Store the Google API key
(**Persona:** devops)

Everything after the **`secret/`** path is a key-value pair to write to the
secret engine. You can specify multiple values. If the value has a space, you
need to surround it with quotes. Having keys with spaces is permitted, but
strongly discouraged because it can lead to unexpected client-side behavior.

Let's assume that the path convention in your organization is
**`secret/<OWNER>/apikey/<APP>`** for API keys. To store the Google API key used
by the engineering team, the path would be `secret/eng/apikey/Google`. If you
have an API key for New Relic owned by the DevOps team, the path would look like
`secret/devops/apikey/New_Relic`.

#### CLI command

To create key/value secrets:

```plaintext
$ vault kv put secret/<PATH> <KEY>=VALUE>
```

The `<PATH>` can be anything you want it to be, and your organization should
decide on the naming convention that makes most sense.

**Example:**

```plaintext
$ vault kv put secret/eng/apikey/Google key=AAaaBBccDDeeOTXzSMT1234BB_Z8JzG7JkSVxI
Success! Data written to: secret/eng/apikey/Google
```

The secret key is "key" and its value is "AAaaBBccDDeeOTXzSMT1234BB_Z8JzG7JkSVxI" in
this example.

#### API call using cURL

Use `/secret/data/<PATH>` endpoint to create secrets:

```plaintext
$ curl --header "X-Vault-Token: <TOKEN>" \
       --request POST \
       --data @payload.json \
       <VAULT_ADDRESS>/v1/secret/data/<PATH>
```

Where `<TOKEN>` is your valid token, and `secret/data/<PATH>` is the path to
your secrets. The [`payload.json`](/api/secret/kv/kv-v2.html#parameters-2)
contains the parameters to invoke the endpoint.

**Example:**

```plaintext
$ tee payload.json <<EOF
{
  "data": {
    "key": "AAaaBBccDDeeOTXzSMT1234BB_Z8JzG7JkSVxI"
  }
}
EOF

$ curl --header "X-Vault-Token: ..." --request POST \
       --data @payload.json \
       http://127.0.0.1:8200/v1/secret/data/eng/apikey/Google
```

The secret key is "key" and its value is
"AAaaBBccDDeeOTXzSMT1234BB_Z8JzG7JkSVxI" in this example.


### <a name="step3"></a>Step 3: Store the root certificate for MySQL
(**Persona:** devops)

For the purpose of this guide, generate a new self-sign certificate using
[OpenSSL](https://www.openssl.org/source/).

```plaintext
$ openssl req --request509 -sha256 -nodes -newkey rsa:2048 -keyout selfsigned.key -out cert.pem
```

Generated `cert.pem` file:

```plaintext
-----BEGIN CERTIFICATE-----
MIIDSjCCAjICCQC47CQCg4u0kDANBgkqhkiG9w0BAQsFADBnMQswCQYDVQQGEwJV
UzELMAkGA1UECAwCQ0ExFjAUBgNVBAcMDVNhbiBGcmFuY2lzY28xFDASBgNVBAMM
C2V4YW1wbGUuY29tMR0wGwYJKoZIhvcNAQkBFg5tZUBleGFtcGxlLmNvbTAeFw0x
ODAxMTcwMTMzNThaFw0xODAyMTYwMTMzNThaMGcxCzAJBgNVBAYTAlVTMQswCQYD
VQQIDAJDQTEWMBQGA1UEBwwNU2FuIEZyYW5jaXNjbzEUMBIGA1UEAwwLZXhhbXBs
ZS5jb20xHTAbBgkqhkiG9w0BCQEWDm1lQGV4YW1wbGUuY29tMIIBIjANBgkqhkiG
9w0BAQEFAAOCAQ8AMIIBCgKCAQEA1cPTXpnOeUXU4tgblNLSS2rcA7eIqzc6gnMY
Sh76WxOaN8VncyJw89/28QYOSYeWRn4fYywbPhHpFmrY6+1gW/8y0+Yoj7TL2Mvs
5m1ZH9eOS6kcnnX/lr+HCfJpTHokKk/Vxr0/p6agkdZq0OYMPAmiuw1M4afd5abm
8s5R99b4DgQyNvRYJp+JMddz2cM8t2AKQH4rq2NEf/GBHqHpHKmaxTyX5Rh7zg/g
WJQ/DjxUVLpbRy+soiUJTZzamrO0iu9fcww+1Q4TZsMWizA4ChQFI7uegKkZ2Alv
SNItsv01FQH3IB7pNWuna3IXXY789R0Qp0Ha5ScryVc9syg4cQIDAQABMA0GCSqG
SIb3DQEBCwUAA4IBAQBtUcuwL0rS/uhk4v53ALF+ryRoLF93wT+O9KOvK15Pi1dX
oZ9yxu5GOGi59womsrDs1vNrBuIQNVQ69dbUYu1LkhgQGDUWQb8JpCp++WHWTIeP
YTJ5C/Q1B3rXeQrVWPvO0bMCig+/G5DGtzZmKWMQGHhfOvSwrkA58YAwjC+rqexl
skA+hQ2JiU4bzIxvlPLBOUA/p+TgUKtdzPY3lxyDO2p7+8ZD56B0PoW87zNJYRcu
VdSr7er8UkUr5nVjcw/6MJeptmx6QaiHgTUSFf2HjFfzsBa/IY1VGr/8bOII+IFN
iYQTLBNG0/q/PZGeMX/RHxmCzZz/7wE0CDPMLbyf
-----END CERTIFICATE-----
```

**NOTE:** If you don't have OpenSSL, simply copy the above certificate and
save it as `cert.pem`.


#### CLI command

The command is basically the same as the Google API key example. Now, the path
convention for certificates is **`secret/<ENVIRONMENT>/cert/<SYSTEM>`**. To
store the root certificate for production MySQL, the path becomes
`secret/prod/cert/mysql`.

**Example:**

```plaintext
$ vault kv put secret/prod/cert/mysql cert=@cert.pem
```

This example reads the root certificate from a PEM file from the disk, and store
it under `secret/prod/cert/mysql` path.
> **NOTE:** Any value begins with "@" is loaded from a file.


#### API call using cURL

To perform the same task using the Vault API, pass the token in the request header.

**Example:**

```plaintext
$ curl --header "X-Vault-Token: ..." \
       --request POST \
       --data @cert.pem \
       http://127.0.0.1:8200/v1/secret/data/prod/cert/mysql
```
> **NOTE:** Any value begins with "@" is loaded from a file.


### <a name="step4"></a>Step 4: Generate a token for apps
(**Persona:** devops)

To read the secrets, `apps` persona needs "read" permit on those secret engine
paths. In this scenario, the `apps` policy must include the following:

**Example:** `apps-policy.hcl`

```shell
# Read-only permit
path "secret/data/eng/apikey/Google" {
  capabilities = [ "read" ]  
}

# Read-only permit
path "secret/data/prod/cert/mysql" {
  capabilities = [ "read" ]
}
```


#### CLI command

First create `apps` policy, and generate a token so that you can authenticate
as an `apps` persona, and read the secrets.

```shell
# Create "apps" policy
$ vault policy write apps apps-policy.hcl
Policy 'apps' written.

# Create a new token with app policy
$ vault token create -policy="apps"
Key            	Value
---            	-----
token          	e4bdf7dc-cbbf-1bb1-c06c-6a4f9a826cf2
token_accessor 	54700b7e--data828-a6c4-6141-96e71e002bd7
token_duration 	768h0m0s
token_renewable	true
token_policies 	[apps default]
```

Now, `apps` can use this token to read the secrets.


#### API call using cURL

First create an `apps` policy, and generate a token so that you can authenticate
as an `app` persona.

**Example:**

```shell
# Payload to pass in the API call
$ tee payload.json <<EOF
{
  "policy": "path \"secret/data/eng/apikey/Google\" { capabilities = [ \"read\" ] ...}"
}
EOF

# Create "apps" policy
$ curl --header "X-Vault-Token: ..." --request PUT \
       --data @payload.json \
       http://127.0.0.1:8200/v1/sys/policy/apps

# Generate a new token with apps policy
$ curl --header "X-Vault-Token: ..." --request POST \
       --data '{"policies": ["apps"]}' \
       http://127.0.0.1:8200/v1/auth/token/create | jq
{
 "request_id": "e1737bc8-7e51-3943-42a0-2dbd6cb40e3e",
 "lease_id": "",
 "renewable": false,
 "lease_duration": 0,
 "data": null,
 "wrap_info": null,
 "warnings": null,
 "auth": {
   "client_token": "1c97b03a-6098-31cf-9d8b-b404e52dcb4a",
   "accessor": "b10a3eb7-15fe-1924-600e-403cfda34c28",
   "policies": [
     "apps",
     "default"
   ],
   "metadata": null,
   "lease_duration": 2764800,
   "renewable": true,
   "entity_id": ""
 }
}
```

Now, `apps` can use this token to read the secrets.

**NOTE:** For the purpose of this guide, you created a policy for `apps`
persona, and generated a token for it. However, in a real world, you may have
a dedicated `policy author`, or `admin` to write policies. Also, the
consumer of the API key may be different from the consumer of the root
certificate. Then each persona would have a policy based on what it needs to
access.

![Personas Introduction](/img/vault-static-secrets2.png)



### <a name="step5"></a>Step 5: Retrieve the secrets
(**Persona:** apps)

Using the token from [Step 4](#step4), read the Google API key, and root certificate for
MySQL.


#### CLI command

The command to read secret is:

```plaintext
$ vault kv get secret/<PATH>
```

**Example:**

```shell
# Authenticate with Vault using the generated token first
$ vault login e4bdf7dc-cbbf-1bb1-c06c-6a4f9a826cf2
Successfully authenticated! You are now logged in.
token: e4bdf7dc-cbbf-1bb1-c06c-6a4f9a826cf2
token_duration: 2764277
token_policies: [apps default]

# Read the API key
$ vault kv get secret/eng/apikey/Google
Key             	Value
---             	-----
refresh_interval	768h0m0s
key             	AAaaBBccDDeeOTXzSMT1234BB_Z8JzG7JkSVxI
```

To return the key value alone, pass `-field=key` as an argument.

```plaintext
$ vault kv get -field=key secret/eng/apikey/Google
AAaaBBccDDeeOTXzSMT1234BB_Z8JzG7JkSVxI
```

#### Root certificate example:

The command is basically the same:

```plaintext
$ vault kv get -field=cert secret/prod/cert/mysql
-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEA6E2Uq0XqreZISgVMUu9pnoMsq+OoK1PI54rsA9vtDE6wiRk0GWhf5vD4DGf1
...
```

#### API call using cURL

Use `secret/` endpoint to retrieve secrets from key/value secret engine:

```plaintext
$ curl --header "X-Vault-Token: <TOKEN_FROM_STEP4>" \
       --request Get \
       <VAULT_ADDRESS>/v1/secret/data/<PATH>
```

**Example:**

Read the Google API key.

```plaintext
$ curl --header "X-Vault-Token: 1c97b03a-6098-31cf-9d8b-b404e52dcb4a" \
       --request GET \
       http://127.0.0.1:8200/v1/secret/data/eng/apikey/Google | jq
{
"request_id": "5a2005ac-1149-2275-cab3-76cee71bf524",
"lease_id": "",
"renewable": false,
"lease_duration": 2764800,
"data": {
  "key": "AAaaBBccDDeeOTXzSMT1234BB_Z8JzG7JkSVxI"
},
"wrap_info": null,
"warnings": null,
"auth": null
}
```

**NOTE:** This example uses `jq` to parse the JSON output.

Retrieve the key value with `jq`:

```plaintext
$ curl --header "X-Vault-Token: 1c97b03a-6098-31cf-9d8b-b404e52dcb4a" \
       --request GET \
       http://127.0.0.1:8200/v1/secret/data/eng/apikey/Google | jq ".data.key"
```

#### Root certificate example:

```plaintext
$ curl --header "X-Vault-Token: 1c97b03a-6098-31cf-9d8b-b404e52dcb4a" \
       --request GET \
       http://127.0.0.1:8200/v1/secret/data/prod/cert/mysql | jq ".data.cert"
```

## Additional Discussion

### Q: How do I enter my secrets without appearing in history?

As a precaution, you may wish to avoid passing your secret as a part of the CLI
command so that the secret won't appear in the history file.  Here are a few
techniques you can use.

#### Option 1: Use a dash "-"

An easy technique is to use a dash "-" and then press Enter. This allows you to
enter the secret in a new line. After entering the secret, press **`Ctrl+d`** to
end the pipe and write the secret to the Vault.

```plaintext
$ vault kv put secret/eng/apikey/Google key=-

AAaaBBccDDeeOTXzSMT1234BB_Z8JzG7JkSVxI
<Ctrl+d>
```

#### Option 2: Read the secret from a file

Using the Google API key example, you can create a file containing the key (apikey.txt):

```plaintext
{
  "key": "AAaaBBccDDeeOTXzSMT1234BB_Z8JzG7JkSVxI"
}
```

The CLI command would look like:

```plaintext
$ vault kv put secret/eng/apikey/Google @apikey.txt
```

#### Option 3: Disable all vault command history

Sometimes, you may not even want the vault command itself to appear in history
at all.  The Option 1 and Option 2 prevents the secret to appear in the history;
however, the vault command, `vault kv put secret/eng/apikey/Google` will appear
in history.

In bash:

```plaintext
$ export HISTIGNORE="&:vault*"
```

**NOTE:** This prevents the use of the Up arrow key for command history as well.


### Q: How do I save multiple values?

The two examples introduced in this guide only had a single key-value pair.  You
can pass multiple values in the command.

```plaintext
$ vault kv put secret/dev/config/mongodb url=foo.example.com:35533 db_name=users \
 username=admin password=passw0rd
```

Or, read the secret from a file:

```plaintext
$ tee mongodb.txt <<EOF
{
    "url": "foo.example.com:35533",
    "db_name": "users",
    "username": "admin",
    "password": "pa$$w0rd"
}
EOF

$ vault kv put secret/dev/config/mongodb @mongodb.txt
```

## Next steps

This guide introduced the CLI commands and API endpoints to read and write
secrets in key/value secret engine. To keep it simple, the `devops` persona
generated a token for `apps`.  Read [AppRole Pull
Authentication](/guides/identity/authentication.html) guide to learn about
programmatically generate a token for apps.
