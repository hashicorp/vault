---
layout: "guides"
page_title: "Static Secrets - Guides"
sidebar_current: "guides-static-secrets"
description: |-
  Vault supports generating new unseal keys as well as rotating the underlying
  encryption keys. This guide covers rekeying and rotating Vault's encryption
  keys.
---

# Static Secrets

Vault can be used to store any secrets in a secure manner.  The secrets may be
SSL certificates and keys for your organization's domain, credentials to connect
to a corporate database server, etc. Storing such sensitive information in a
plaintext is not desirable. This guide demonstrates the use case of Vault as a
Secret Storage.


## Reference Material

- [Key/Value Secret Backend](/docs/secrets/kv/index.html)
- [Key/Value Secret Backend API](/apikey/secret/kv/index.html)
- [Client libraries](/apikey/libraries.html) for Vault API for commonly used languages

## Estimated Time to Complete

10 minutes

## Challenge

Consider the following situations:

- Developers use a single admin account to access a third-party app
  (e.g. Splunk) and anyone who knows the user ID and password can log in as an
  admin
- SSH keys to connect to remote machines are shared and stored as a plaintext
- API keys to invoke external system APIs are stored as a plaintext
- An app integrates with LDAP, and its configuration information is in a
  plaintext

Organizations often seek an uniform solution to store any sensitive information
securely.

## Solution

Leverage Vault as a centralized secret storage to secure any sensitive
information. Vault encrypts these secrets using 256-bit AES in GCM mode with a
randomly generated nonce prior to writing them to its persistent storage. The
storage backend never sees the unencrypted value, so gaining access to the raw
storage isn't enough to access your secrets.


## Prerequisites

To perform the tasks described in this guide, you need to have a Vault
environment.  Refer to the [Getting
Started](/intro/getting-started/install.html) guide to install Vault. Make sure
that your Vault server has been [initialized and
unsealed](/intro/getting-started/deploy.html).


## Steps

This guide demonstrates the basic steps to store secrets using Vault.  The
scenario here is to store the following secrets:

- API key (Google API)
- Root certificate of a production database (MySQL)

To store your API key within the configured physical storage for Vault, use the
key/value secret backend via **`secret/`** prefixed.

-> Key/Value secret backend passes through any operation back to the configured
storage backend for Vault. For example, if your Vault server is configured with
Consul as its storage backend, a "read" operation turns into a read from Consul
at the same path.



You will perform the following:

1. [Store the Google API key](#step1)
2. [Store the root certificate for MySQL](#step2)
3. [Retrieve the secrets](#step3)

### <a name="step1"></a>Step 1: Store the Google API key

If the secret path convention is **`secret/<OWNER>/apikey/<APP>`**, store the
Google API key in `secret/eng/apikey/Googl`. If you have an API key for New Relic
owned by the DevOps team, the path may look like
`secret/devops/apikey/New_Relic`.

#### CLI command

To create key/value secrets:

```shell
$ vault write secret/<PATH> <KEY>=VALUE>
```

The `<PATH>` can be anything you want it to be, and your organization should
decide on the naming convention that makes most sense.

**Example:**

```shell
$ vault write secret/eng/apikey/Google key=AAaaBBccDDeeOTXzSMT1234BB_Z8JzG7JkSVxI
Success! Data written to: secret/eng/apikey/Google
```

The secret key is "key" and its value is "AAaaBBccDDeeOTXzSMT1234BB_Z8JzG7JkSVxI" in
this example.

#### API call using cURL

Use `/secret/<PATH>` endpoint to create secrets:

```shell
$ curl --header "X-Vault-Token: <TOKEN>" \
       --request POST \
       --data <SECRETS> \
       <VAULT_ADDRESS>/v1/secret/<PATH>
```

Where `<TOKEN>` is your valid token, `<SECRETS>` is the key-value pair(s) of your
secrets, and `secret/<PATH>` is the path to your secrets.

**Example:**

```shell
$ curl --header "X-Vault-Token: ..." --request POST \
       --data '{"key": "AAaaBBccDDeeOTXzSMT1234BB_Z8JzG7JkSVxI"}' \
       https://vault.rocks/v1/secret/eng/apikey/Google
```

The secret key is "key" and its value is
"AAaaBBccDDeeOTXzSMT1234BB_Z8JzG7JkSVxI" in this example.


### <a name="step2"></a>Step 2: Store the root certificate for MySQL

For the purpose of this guide, generate a new self-sign certificate using
[OpenSSL](https://www.openssl.org/source/).

```shell
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
convention is **`secret/<ENVIRONMENT>/cert/<SYSTEM>`**. To store the root certificate
for production MySQL, the path becomes `secret/staging/cert/postgres`.

**Example:**

```shell
$ vault write secret/prod/cert/mysql cert=@cert.pem
```

This example reads the root certificate from a PEM file from the disk, and store it under
`secret/prod/cert/mysql` path.
> **NOTE:** Any value begins with "@" is loaded from a file.


#### API call using cURL

To perform the same task using the Vault API, pass the token in the request header.

**Example:**

```shell
$ curl --header "X-Vault-Token: ..." \
       --request POST \
       --data @cert.pem \
       https://vault.rocks/v1/secret/prod/cert/mysql
```
> **NOTE:** Any value begins with "@" is loaded from a file.


### <a name="step3"></a>Step 3: Retrieve the secrets

Retrieving the secret from Vault is simple.  

#### CLI command

```shell
$ vault read secret/<PATH>
```

**Example:**

```shell
$ vault read secret/eng/apikey/Google
Key             	Value
---             	-----
refresh_interval	768h0m0s
key             	AAaaBBccDDeeOTXzSMT1234BB_Z8JzG7JkSVxI
```

To return the key value alone, pass `-field=key` as an argument.

```shell
$ vault read -field=key secret/eng/apikey/Google
AAaaBBccDDeeOTXzSMT1234BB_Z8JzG7JkSVxI
```

#### Root certificate example:

```shell
$ vault read -field=cert secret/prod/cert/mysql
-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEA6E2Uq0XqreZISgVMUu9pnoMsq+OoK1PI54rsA9vtDE6wiRk0GWhf5vD4DGf1
...
```

#### API call using cURL

**Example:**

```shell
$ curl --header "X-Vault-Token: ..." --request GET \
     https://vault.rocks/v1/secret/eng/apikey/Google | jq
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

```shell
$ curl --header "X-Vault-Token: ..." --request GET \
     https://vault.rocks/v1/secret/eng/apikey/Google | jq ".data.key"
```

#### Root certificate example:

```shell
$ curl --header "X-Vault-Token: ..." --request GET \
     https://vault.rocks/v1/secret/prod/cert/mysql | jq ".data.cert"
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

```shell
$ vault write secret/eng/apikey/Google key=-

AAaaBBccDDeeOTXzSMT1234BB_Z8JzG7JkSVxI
<Ctrl+d>
```

#### Option 2: Read the secret from a file

Using the Google API key example, you can create a file containing the key (apikey.txt):

```text
{
  "key": "AAaaBBccDDeeOTXzSMT1234BB_Z8JzG7JkSVxI"
}
```

The CLI command would look like:

```shell
$ vault write secret/eng/apikey/Google @apikey.txt
```

#### Option 3: Disable all vault command history

Sometimes, you may not even want the vault command itself to appear in history
at all.  The Option 1 and Option 2 prevents the secret to appear in the history;
however, the vault command, `vault write secret/eng/apikey/Google` will appear
in history.

In bash:

```shell
$ export HISTIGNORE="&:vault"
```

**NOTE:** This prevents the use of the Up arrow key for command history as well.


### Q: How do I save multiple values?

The two examples introduced in this guide only had a single key-value pair.  You can pass multiple values in the command.

```shell
$ vault write secret/dev/config/mongodb url=foo.example.com:35533 db_name=users \
 username=admin password=pa$$w0rd
```

Or, read the secret from a file:

```shell
$ vault write secret/dev/config/mongodb @mongodb.txt

$ cat mongodb.txt
{
  "url": "foo.example.com:35533",
  "db_name": "users",
  "username": "admin",
  "password": "pa$$w0rd"
}
```


## Next steps
