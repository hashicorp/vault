---
layout: "guides"
page_title: "Seal Wrap / FIPS 140-2 - Guides"
sidebar_title: "Seal Wrap / FIPS 140-2"
sidebar_current: "guides-operations-seal-wrap"
description: |-
  In this guide,
---


# Seal Wrap / FIPS 140-2

~> **Enterprise Only:** Vault's HSM auto-unseal and Seal Wrap features are a
part of _Vault Enterprise_.

***Vault Enterprise*** integrates with [HSM platforms](/docs/enterprise/hsm/index.html)
to opt-in automatic [unsealing](/docs/concepts/seal.html#unsealing).
HSM integration provides three pieces of special functionality:

- **Master Key Wrapping**: Vault protects its master key by transiting it through
the HSM for encryption rather than splitting into key shares
- **Automatic Unsealing**: Vault stores its encrypted master key in storage,
allowing for automatic unsealing
- **Seal Wrapping** to provide FIPS KeyStorage-conforming functionality for Critical Security Parameters

![Unseal with HSM](/img/vault-hsm-autounseal.png)

In some large organizations, there is a fair amount of complexity in designating
key officers, who might be available to unseal Vault installations as the most
common pattern is to deploy Vault immutably. As such automating unseal using an
HSM provides a simplified yet secure way of unsealing Vault nodes as they get
deployed.

Vault pulls its encrypted master key from storage and transit it through the
HSM for decryption via  **PKCS \#11 API**. Once the master key is decrypted,
Vault uses the master key to decrypt the encryption key to resume with Vault
operations.


## Reference Material

- [HashiCorp + AWS: Integrating CloudHSM with Vault Enterprise](https://www.hashicorp.com/resources/hashicorp-and-aws-integrating-cloudhsm-with-vault-e) webinar
- [Seal Wrap documentation](/docs/enterprise/sealwrap/index.html)
- [Vault Configuration - pkcs11 Seal](/docs/configuration/seal/pkcs11.html)
- [Vault Enterprise HSM Support](/docs/enterprise/hsm/index.html)
- [NIST SC-12: Cryptographic Key Establishment and Management](https://nvd.nist.gov/800-53/Rev4/control/SC-12)
- [NIST SC-13: Cryptographic Protection](https://nvd.nist.gov/800-53/Rev4/control/SC-13)


## Estimated Time to Complete

10 minutes


## Challenge

The Federal Information Processing Standard (FIPS) 140-2 is a U.S. Government
computer security standard used to accredit cryptography modules. If your
product or service does not follow FIPS' security requirements, it may
complicate your ability to operate with U.S. Government data.

Aside from doing business with U.S. government, your organization may care about
FIPS which approves various cryptographic ciphers for hashing, signature, key
exchange, and encryption for security.


## Solution

Integrate Vault with FIPS 140-2 certified HSM and enable the ***Seal Wrap***
feature to protect your data.

Vault encrypts secrets using 256-bit AES in GCM mode with a randomly generated
nonce prior to writing them to its persistent storage. By enabling seal wrap,
Vault wraps your secrets with **an extra layer of encryption** leveraging the
HSM encryption and decryption.

![Seal Wrap](/img/vault-seal-wrap.png)

#### Benefits of the Seal Wrap:

- Conformance with FIPS 140-2 directives on Key Storage and Key Transport as [certified by Leidos](/docs/enterprise/sealwrap/index.html#fips-140-2-compliance)
- Supports FIPS level of security equal to HSM
  * For example, if you use Level 3 hardware encryption on an HSM, Vault will be
  using FIPS 140-2 Level 3 cryptography
- Allows Vault to be deployed in high security [GRC](https://en.wikipedia.org/wiki/Governance,_risk_management,_and_compliance)
environments (e.g. PCI-DSS, HIPAA) where FIPS guidelines important for external audits
- Pathway for Vault's use in managing  Department of Defense's (DOD) or North
Atlantic Treaty Organization (NATO) military secrets


## Prerequisites

This intermediate operations guide assumes that you have:

- A [supported HSM](/docs/enterprise/hsm/index.html) cluster to be integrated
  with Vault
- Vault Enterprise Premium



## Steps

This guide walks you through the following steps:

1. [Configure HSM Auto-unseal](#step1)
1. [Enable Seal Wrap](#step2)
1. [Test the Seal Wrap Feature](#step3)



### <a name="step1"></a>Step 1: Configure HSM Auto-unseal

When a Vault server is started, it normally starts in a sealed state where a
quorum of existing unseal keys is required to unseal it. By integrating Vault
with HSM, your Vault server can be automatically unsealed by the trusted HSM key
provider.

#### Task 1: Write a Vault configuration file

To integrate your Vault Enterprise server with an HSM cluster, the configuration
file must define the [`PKCS11 seal` stanza](/docs/configuration/seal/pkcs11.html)
providing necessary connection information.


**Example: `config-hsm.hcl`**

```shell
# Provide your AWS CloudHSM cluster connection information
seal "pkcs11" {
  lib = "/opt/cloudhsm/lib/libcloudhsm_pkcs11.so"
  slot = "1"
  pin = "vault:Password1"
  key_label = "hsm_demo"
  hmac_key_label = "hsm_hmac_demo"
  generate_key = "true"
}

# Configure the storage backend for Vault
storage "file" {
  path = "/tmp/vault"
}

# Addresses and ports on which Vault will respond to requests
listener "tcp" {
  address     = "127.0.0.1:8200"
  tls_disable = 1
}

ui = true
```

> **NOTE:** For the purpose of this guide, the storage backend is set to the
local file system (`/tmp/vault`) to make the verification step easy.

The example configuration defines the following in its **`seal`** stanza:

- **`lib`** is set to the path to the PKCS \#11 library on the virtual machine
  where Vault Enterprise is installed
- **`slot`** should be set to the slot number to use
- **`pin`** is the PKCS \#11 PIN for login
- **`key_label`** defines the label of the key you want to use
- **`hmac_key_label`** defines the label of the key you want to use for HMACing.
  (NOTE: HMAC is optional and only used for mechanisms that do not support
  authenticated data.)
- **`generate_key`** is set to `true`.  If no existing key with the label
  specified by `key_label` can be found at Vault initialization time, Vault
  generates a key

~> **IMPORTANT:** Having Vault generate its own key is the easiest way to get up
and running, but for security, Vault marks the key as **non-exportable**. If
your HSM key backup strategy requires the key to be exportable, you should
generate the key yourself.  Refer to the [key generation  attributes](/docs/configuration/seal/pkcs11.html#vault-key-generation-attributes).




#### Task 2: Initialize your Vault Enterprise server

Start the Vault server with your Vault configuration file. For example, if your
configuration file is located at `/home/ec2-user/config-hsm.hcl`, the command
would look like:

```plaintext
$ vault server -config=/home/ec2-user/config-hsm.hcl

  SDK Version: 2.03
==> Vault server configuration:

   HSM PKCS#11 Version: 2.40
           HSM Library: Cavium PKCS#11 Interface
   HSM Library Version: 1.0
   HSM Manufacturer ID: Cavium Networks
              HSM Type: pkcs11
                   Cgo: enabled
            Listener 1: tcp (addr: "127.0.0.1:8200", cluster address: "127.0.0.1:8201", tls: "disabled")
             Log Level: info
                 Mlock: supported: true, enabled: false
               Storage: file
               Version: Vault v0.10.1+ent.hsm
           Version Sha: 0e628142d6b6e5cabfdb9680a6d669d38f15574f

==> Vault server started! Log data will stream in below:
```

<br>

In another terminal, set the `VAULT_ADDR` environment variable, and [initialize]
(/intro/getting-started/deploy.html#initializing-the-vault) your Vault server.

**Example:**

```shell
# Set the VAULT_ADDR environment variable
$ export VAULT_ADDR="http://127.0.0.1:8200"

# Initialize Vault
$ vault operator init

Recovery Key 1: 2bU2wOfmyMqYcsEYo4Mo9q4s/KAODgHHjcmZmFOo+XY=
Initial Root Token: 8d726c6b-98ba-893f-23d5-be3d2fec480e

Success! Vault is initialized

Recovery key initialized with 1 key shares and a key threshold of 1. Please
securely distribute the key shares printed above.
```

There is only a single master key created which is encrypted by the HSM using
PKCS \#11, and then placed in the storage. When Vault needs to be unsealed, it
grabs the HSM encrypted master key from the storage, round trips it through the
HSM to decrypt the master key.

~> **NOTE:** When Vault is initialized while using an HSM, rather than unseal
keys being returned to the operator, **recovery keys** are returned. These are
generated from an internal recovery key that is [split via Shamir's Secret
Sharing](/docs/enterprise/hsm/behavior.html#initialization), similar to Vault's
treatment of unseal keys when running without an HSM. Some Vault operations such
as generation of a root token require these recovery keys.

Login to the Vault using the generated root token to verify.

```plaintext
$ vault login 8d726c6b-98ba-893f-23d5-be3d2fec480e
```

#### Task 3: Verification

Stop and restart the Vault server and then verify its status:

```plaintext
$ vault status

Key                      Value
---                      -----
Recovery Seal Type       shamir
Sealed                   false
Total Recovery Shares    1
Threshold                1
Version                  0.10.1+ent.hsm
Cluster Name             vault-cluster-80556565
Cluster ID               40316cdd-3d42-ec36-e7b0-6a7a0684568c
HA Enabled               false
```

The `Sealed` status is **`false`** which means that the Vault was automatically
unsealed upon its start. You can proceed with Vault operations.


### <a name="step2"></a>Step 2: Enable Seal Wrap

-> **NOTE:** For FIPS 140-2 compliance, seal wrap requires FIPS
140-2 Certified HSM which is supported by _Vault Enterprise Premium_.

For some values, seal wrapping is **always enabled** including the recovery key, any
stored key shares, the master key, the keyring, and more. When working with the
key/value secret engine, you can enable seal wrap to wrap all data.


### CLI command

Check the enabled secret engines.

```plaintext
$ vault secrets list -format=json
{
  ...
  "secret/": {
    "type": "kv",
    "description": "key/value secret storage",
    "accessor": "kv_75820543",
    "config": {
      "default_lease_ttl": 0,
      "max_lease_ttl": 0,
      "force_no_cache": false
    },
    "options": {
      "version": "1"
    },
    "local": false,
    "seal_wrap": false
  },
  ...
```

Notice that the `seal_wrap` parameter is set to **`false`**.

> For the purpose of comparing seal wrapped data against unwrapped data, enable
additional key/value secret engine at the `secret2/` path.

```shell
# Pass the '-seal-wrap' flag when you enable the KV workflow
$ vault secrets enable -path=secret2/ -version=1 -seal-wrap kv
```

The above command enabled [key/value version 1](/docs/secrets/kv/kv-v1.html) with
seal wrap feature enabled.

```plaintext
$ vault secrets list -format=json
{
  ...
  "secret2/": {
    "type": "kv",
    "description": "",
    "accessor": "kv_bdd74241",
    "config": {
      "default_lease_ttl": 0,
      "max_lease_ttl": 0,
      "force_no_cache": false
    },
    "options": {
      "version": "1"
    },
    "local": false,
    "seal_wrap": true
  },
  ...
```

Notice that the `seal_wrap` parameter is set to **`true`** at `secret2/`.


#### API call using cURL

Check the enabled secret engines.

```plaintext
$ curl --header "X-Vault-Token: ..." \
       http://127.0.0.1:8200/v1/sys/mounts | jq
...
  "secret/": {
    "accessor": "kv_f05b8b9c",
    "config": {
      "default_lease_ttl": 0,
      "force_no_cache": false,
      "max_lease_ttl": 0
    },
    "description": "key/value secret storage",
    "local": false,
    "options": {
      "version": "2"
    },
    "seal_wrap": false,
    "type": "kv"
  },
...
```

Notice that the `seal_wrap` parameter is set to **`false`**.

> For the purpose of comparing seal wrapped data against unwrapped data, enable
additional key/value secret engine at the `secret2/` path.

```shell
# Set the seal_wrap parameter to true in the request payload
$ tee payload.json <<EOF
{
  "type": "kv",
  "options": {
    "version": "1"
  },
  "seal_wrap": true
}
EOF

# Enable kv secret engine at secret2/
$ curl --header "X-Vault-Token: ..." \
       --request POST \
       --data @payload.json \
       http://127.0.0.1:8200/v1/sys/mounts/secret2
```

The above command enabled [key/value version 1](/docs/secrets/kv/kv-v1.html) with
seal wrap feature enabled.

```plaintext
$ curl --header "X-Vault-Token: ..." \
       http://127.0.0.1:8200/v1/sys/mounts | jq
{
...
  "secret2/": {
    "accessor": "kv_724c81c9",
    "config": {
      "default_lease_ttl": 0,
      "force_no_cache": false,
      "max_lease_ttl": 0
    },
    "description": "",
    "local": false,
    "options": {
      "version": "1"
    },
    "seal_wrap": true,
    "type": "kv"
  },
  ...
}
```

Notice that the `seal_wrap` parameter is set to **`true`** at `secret2/`.


#### Web UI

Open a web browser and launch the Vault UI (e.g. `http://127.0.0.1:8200/ui`) and
then login.

![Enable Secret Engine](/img/vault-seal-wrap-2.png)

> For the purpose of comparing seal wrapped data against unwrapped data, enable
additional key/value secret engine at the `secret2/` path.

Select **Enable new engine**.  

- Enter **`secret2`** in the path field
- Select **Version 1** for KV version
- Select the check box for **Seal Wrap**

![Enable Secret Engine](/img/vault-seal-wrap-3.png)

Click **Enable Engine**.


### <a name="step3"></a>Step 3: Test the Seal Wrap Feature

In this step, you are going to:

1. Write some test data
1. [View the encrypted secrets](#view-the-encrypted-secrets)

#### CLI command

Write a secret at `secret/unwrapped`.

```shell
# Write a key named 'password' with its value 'my-long-password'
$ vault kv put secret/unwrapped password="my-long-password"

# Read the path to verify
$ vault kv get secret/unwrapped
====== Data ======
Key         Value
---         -----
password    my-long-password
```

Write the same secret at `secret2/wrapped`.

```shell
# Write a key named 'password' with its value 'my-long-password'
$ vault kv put secret2/wrapped password="my-long-password"

# Read the path to verify
$ vault kv get secret2/wrapped
====== Data ======
Key         Value
---         -----
password    my-long-password
```
Using a valid token, you can write and read secrets the same way
regardless of the seal wrap.



#### API call using cURL

Write a secret at `secret/unwrapped`.

```shell
# Create a payload
$ tee payload.json <<EOF
{
  "data": {
    "password": "my-long-password"
  }
}
EOF

# Write secret at 'secret/unwrapped'
$ curl --header "X-Vault-Token: ..." \
       --request POST \
       --data @payload.json \
       http://127.0.0.1:8200/v1/secret/unwrapped

# Read the path to verify
$ curl --header "X-Vault-Token: ..." \
       http://127.0.0.1:8200/v1/secret/unwrapped | jq
{
   "request_id": "cef0b061-3860-6a3c-8bba-1fdbf7be7bf1",
   "lease_id": "",
   "renewable": false,
   "lease_duration": 2764800,
   "data": {
     "data": {
       "password": "my-long-password"
     }
   },
   "wrap_info": null,
   "warnings": null,
   "auth": null
}       
```

Write the same secret at `secret2/wrapped`.

```shell
# Write the same secret at 'secret2/wrapped'
$ curl --header "X-Vault-Token: ..." \
       --request POST \
       --data @payload.json \
       http://127.0.0.1:8200/v1/secret2/wrapped

# Read the path to verify
$ curl --header "X-Vault-Token: ..." \
      http://127.0.0.1:8200/v1/secret2/wrapped | jq
{
  "request_id": "cc9ce831-2ca7-7115-a899-ce33767bdc2c",
  "lease_id": "",
  "renewable": false,
  "lease_duration": 2764800,
  "data": {
    "data": {
      "password": "my-long-password"
    }
  },
  "wrap_info": null,
  "warnings": null,
  "auth": null
}
```

Using a valid token, you can write and read secrets the same way regardless of
the seal wrap.

#### Web UI

Select **secret** and click **Create secret**.

![Enable Secret Engine](/img/vault-seal-wrap-4.png)

Enter the following:

- PATH FOR THIS SECTE: `unwrapped`
- key: `password`
- value: `my-long-password`

![Enable Secret Engine](/img/vault-seal-wrap-5.png)

Click **Save**.

Repeat the same step for **secret2** to write the same secret at the
`secret2/wrapped` path.

![Enable Secret Engine](/img/vault-seal-wrap-6.png)
Click **Save**.


Using a valid token, you can write and read secrets the same way
regardless of the seal wrap.


#### View the encrypted secrets

Remember that the Vault server was configured to use the local file system
(`/tmp/vault`) as its storage backend in this example.

```shell
# Configure the storage backend for Vault
storage "file" {
  path = "/tmp/vault"
}
```

SSH into the machine where the Vault server is running, and check the stored
values in the `/tmp/vault` directory.

```plaintext
$ cd /tmp/vault/logical
```

Under the `/tmp/vault/logical` directory, there are two sub-directories. One
maps to `secret/` and another maps to `secret2/` although you cannot tell by
the folder names.


View the secret at rest.

```shell
# One of the directory maps to secret/unwrapped
$ cd 2da357cd-55f2-7eed-c46e-c477b70bed18

# View its content - password value is encrypted
$ cat _unwrapped
{"Value":"AAAAAQICk547prhuhMiBXLq2lx8ZkMpSB3p+GKHAwuMhKrZGSeqsFevMS6YoqTVlbvpU9B4zWPZ2HA
SeNZ3YMw=="}

# Another directory maps to secret2/wrapped
$ cd ../5bcea44d-28a3-87af-393b-c6d398fe41d8

# View its content - password value is encrypted
$ cat _wrapped
{"Value":"ClBAg9oN7zBBaDBZcsilDAyGkL7soPe7vBA5+ADADuyzo8GuHZHb9UFN2nF1h0OpKEgCIkG3JNHcXt
tZqCi6szcuNBgF3pwhWGwB4FREM3b5CRIQYK7239Q92gRGrcBBeZD6ghogEtSBDmZJBahk7n4lIYF3X4iBqmwZgH
Vo4lzWur7rzncgASofCIIhENEEGghoc21fZGVtbyINaHNtX2htYWNfZGVtb3M="}
```

Secrets are encrypted regardless; however, the seal-wrapped value is
significantly longer despite the fact that both values are the same,
`my-long-password`.


~> When Vault's Seal Wrap feature is used with a FIPS 140-2 certified HSM, Vault
will store Critical Security Parameters (CSPs) in a manner that is compliant
with KeyStorage and KeyTransit requirements.



## Next steps

This guide used the local file system as the storage backend to keep it simple.
To learn more about making your Vault cluster highly available, read the [Vault
HA with Consul](/guides/operations/vault-ha-consul.html) guide.
