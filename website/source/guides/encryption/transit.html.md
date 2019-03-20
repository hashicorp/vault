---
layout: "guides"
page_title: "Encryption as a Service - Guides"
sidebar_title: "Encryption as a Service"
sidebar_current: "guides-encryption-transit"
description: |-
  HashiCorp Vault's transit secrets engine handles cryptographic functions on data in-transit. It can also viewed as _encryption as a service_.
---

# Encryption as a Service: Transit Secrets Engine

Vault's `transit` secrets engine handles cryptographic functions on
data-in-transit. Vault doesn't store the data sent to the secrets engine, so it
can also be viewed as ***encryption as a service***.  

Although the `transit` secrets engine provides additional features (sign and
verify data, generate hashes and HMACs of data, and act as a source of random
bytes), its primary use case is to encrypt data. This relieves the burden of
proper encryption/decryption from application developers and pushes the burden
onto the operators of Vault.


## Reference Materials

- [Transit Secret Engine](/docs/secrets/transit/index.html)
- [Transit Secret Engine API](/api/secret/transit/index.html)
- [Transparent Data Encryption in the Modern Datacenter](https://www.hashicorp.com/blog/transparent-data-encryption-in-the-modern-datacenter)

~> **NOTE:** An [interactive
tutorial](https://www.katacoda.com/hashicorp/scenarios/vault-transit) is
also available if you do not have a Vault environment to perform the steps
described in this guide.


## Estimated Time to Complete

10 minutes


## Personas

The end-to-end scenario described in this guide involves two personas:

- **operator** with privileged permissions to manage the encryption keys
- **app** with un-privileged permissions encrypt/decrypt secrets via API


## Challenge

Think of the following scenario:

_Example Inc._ recently made headlines for a massive data breach which exposed
millions of their users' payment card accounts online. When they tracked down the
problem they found that a new HVAC system with management software had been put
into their data centers and had created vulnerabilities in their networks and
exposed ports and IPs to the databases publicly.

## Solutions

The `transit` secrets engine enables security teams to fortify data during
transit and at rest. So even if an intrusion occurs, your data is encrypted with
AES 256-bit CBC encryption (TLS in transit). Even if an attacker were able to
access the raw data, they would only have encrypted bits. This means attackers
would need to compromise multiple systems before exfiltrating data.

![Encryption as a Service](/img/vault-encryption.png)

This guide demonstrates the basics of the `transit` secrets engine.

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
# Enable transit secrets engine
path "sys/mounts/transit" {
  capabilities = [ "create", "read", "update", "delete", "list" ]
}

# To read enabled secrets engines
path "sys/mounts" {
  capabilities = [ "read" ]
}

# Manage the transit secrets engine
path "transit/*" {
  capabilities = [ "create", "read", "update", "delete", "list" ]
}
```

If you are not familiar with policies, complete the
[policies](/guides/identity/policies.html) guide.



## Steps

You will perform the following:

1. [Configure Transit Secret Engine](#step1)
1. [Encrypt Secrets](#step2)
1. [Decrypt a cipher-text](#step3)
1. [Rotate the Encryption Key](#step4)
1. [Update Key Configuration](#step5)


### <a name="step1"></a>Step 1: Configure Transit Secret Engine
(**Persona:** operator)

The `transit` secrets engine must be configured before it can perform its
operations.  This step is usually done by an **operator** or configuration
management tool.


#### CLI command

Enable the `transit` secret engine by executing the following command:

```plaintext
$ vault secrets enable transit
```

> By default, the secrets engine will mount at the name of the engine.  If you
wish to enable it at a different path, use the `-path` argument.

> **Example:** `vault secrets enable -path=encryption transit`

Now, create an encryption key ring named, `orders` by executing the following
command:

```plaintext
$ vault write -f transit/keys/orders
```


#### API call using cURL

Enable `transit` secret engine using `/sys/mounts` endpoint:

```plaintext
$ curl --header "X-Vault-Token: <TOKEN>" \
       --request POST \
       --data <PARAMETERS> \
       <VAULT_ADDRESS>/v1/sys/mounts/<PATH>
```

Where `<TOKEN>` is your valid token, and `<PARAMETERS>` holds [configuration
parameters](/api/system/mounts.html#enable-secrets-engine) of the secret engine.

**Example:**

The following example enables transit secret engine at `sys/mounts/transit`
path, and passed the secret engine type (`transit`) in the request payload.

```plaintext
$ curl --header "X-Vault-Token: ..." \
       --request POST \
       --data '{"type":"transit"}' \
       https://127.0.0.1:8200/v1/sys/mounts/transit
```

Now, create an encryption key ring named, `orders` using the `transit/keys`
endpoint:

```plaintext
$ curl --header "X-Vault-Token: ..." \
       --request POST \
       https://127.0.0.1:8200/v1/transit/keys/orders    
```


#### Web UI

Open a web browser and launch the Vault UI (e.g. http://127.0.0.1:8200/ui) and then login.

1. Select **Enable new engine** and select **Transit** from **Secrets engine type**
drop-down list.
  ![Enable new engine](/img/vault-secrets-enable.png)

1. Click **Enable Engine**.

1. Select **Create encryption key** and enter `orders` in the **Name** field.
  ![Create a key](/img/vault-transit-1.png)

1. Click **Create encryption key** to complete.

<br>

~> **NOTE:** Typically, you want to create an encryption key ring for each
application.



### <a name="step2"></a>Step 2: Encrypt Secrets
(**Persona:** operator)

Once the `transit` secrets engine has been configured, any client with a valid
token with proper permission can send data to encrypt.

Here, you are going to encrypt a plaintext, _"credit-card-number"_.

-> **NOTE:** You can pass non-text binary file such as a PDF or image.
When you encrypt a plaintext, it must be base64-encoded.


#### CLI command

To encrypt your secret, use the `transit/encrypt` endpoint:

```plaintext
$ vault write transit/encrypt/<key_ring_name>
```

Execute the following command to encrypt a plaintext:

```plaintext
$ vault write transit/encrypt/orders plaintext=$(base64 <<< "credit-card-number")

Key           Value
---           -----
ciphertext    vault:v1:cZNHVx+sxdMErXRSuDa1q/pz49fXTn1PScKfhf+PIZPvy8xKfkytpwKcbC0fF2U=
```

Vault does *NOT* store any of this data. The output you received is the
ciphertext. You can store this ciphertext at the desired location (e.g. MySQL
database) or pass it to another application.



#### API call using cURL

To encrypt your secret, use the [`transit/encrypt`
endpoint](/api/secret/transit/index.html#encrypt-data).

**Example:**

```shell
# Generate base64-encoded plaintext
$ base64 <<< "credit-card-number"
Y3JlZGl0LWNhcmQtbnVtYmVyCg==

# Pass the base64-encoded plaintext in the request payload
$ curl --header "X-Vault-Token: ..." \
       --request POST \
       --data '{"plaintext": "Y3JlZGl0LWNhcmQtbnVtYmVyCg=="}' \
       https://127.0.0.1:8200/v1/transit/encrypt/orders | jq
{
  "request_id": "f483d9b6-8132-782e-1665-ad432c2461ab",
  "lease_id": "",
  "renewable": false,
  "lease_duration": 0,
  "data": {
    "ciphertext": "vault:v1:/9hdQutaWpZR72s3+VSCLK1JNhV1wKM49hYVjh7RjmuxIy/OvshtgV4L4uVB+aQ="
  },
  ...
}
```

Vault does *NOT* store any of this data. The output you received is the
ciphertext. You can store this ciphertext at the desired location (e.g. MySQL
database) or pass it to another application.



#### Web UI

1. Select the **orders** encryption key.

1. Select **Key actions**.
  ![Key action](/img/vault-transit-2.png)

1. Make sure that **Encrypt** is selected under **TRANSIT ACTIONS**, and then
enter "credit-card-number" in the **Plaintext** field.
  ![Encrypt plaintext](/img/vault-transit-3.png)

1. Click **Encode to base64** to encode the plaintext.

1. Click **Encrypt**.
  Vault does *NOT* store any of this data. The output you received is the
ciphertext. You can click **Copy** to copy the resulting ciphertext and store it
at the desired location (e.g. MySQL database) or pass it to another application.
![Encrypt plaintext](/img/vault-transit-4.png)



### <a name="step3"></a>Step 3: Decrypt a cipher-text
(**Persona:** operator)

Any client with a valid token with proper permission can decrypt the ciphertext
generated by Vault. To decrypt the ciphertext, invoke the `transit/decrypt`
endpoint.


#### CLI command

Execute the following command to decrypt the ciphertext resulted in [Step
2](#step2):

```plaintext
$ vault write transit/decrypt/orders \
        ciphertext="vault:v1:cZNHVx+sxdMErXRSuDa1q/pz49fXTn1PScKfhf+PIZPvy8xKfkytpwKcbC0fF2U=" \

Key          Value
---          -----
plaintext    Y3JlZGl0LWNhcmQtbnVtYmVyCg==
```

The resulting data is base64-encoded.  To reveal the original plaintext, run the
following command:

```plaintext
$ base64 --decode <<< "Y3JlZGl0LWNhcmQtbnVtYmVyCg=="
credit-card-number
```


#### API call using cURL

Use the `transit/decrypt` endpoint to decrypt the ciphertext resulted in [Step
2](#step2):

**Example:**

```shell
# Pass the ciphertext in the request payload to decode
$ curl --header "X-Vault-Token: ..." \
       --request POST \
       --data '{"ciphertext": "Yvault:v1:/9hdQutaWpZR72s3+VSCLK1JNhV1wKM49hYVjh7RjmuxIy/OvshtgV4L4uVB+aQ="}' \
       https://127.0.0.1:8200/v1/transit/decrypt/orders | jq
{
   "request_id": "062d7998-8932-76f2-f96c-5938a55ff005",
   "lease_id": "",
   "renewable": false,
   "lease_duration": 0,
   "data": {
     "plaintext": "Y3JlZGl0LWNhcmQtbnVtYmVyCg=="
   },
   ...
}

# The resulting data is base64-encoded that it must be decoded to reveal the plaintext
$ base64 --decode <<< "Y3JlZGl0LWNhcmQtbnVtYmVyCg=="
credit-card-number
```


#### Web UI

1. Select the **orders** encryption key.

1. Select **Key actions**.

1. Make sure that **Decrypt** is selected under **TRANSIT ACTIONS**, and then
enter the ciphertext you wish to decrypt.
  ![Decrypt ciphertext](/img/vault-transit-5.png)

1. Click **Decrypt**.

1. The resulting data is base64-encoded. Click **Decode from base64** to reveal
the plaintext.



### <a name="step4"></a>Step 4: Rotate the Encryption Key
(**Persona:** operator)

One of the benefits of using the Vault `transit` secrets engine is its ability
to easily rotate the encryption keys. Keys can be rotated manually by a human,
or an automated process which invokes the key rotation API endpoint through
cron, a CI pipeline, a periodic Nomad batch job, Kubernetes Job, etc.

Vault maintains the versioned keyring and the operator can decide the minimum
version allowed for decryption operations. When data is encrypted using Vault,
the resulting ciphertext is prepended with the version of the key used to
encrypt it.

#### CLI command

To rotate the encryption key, invoke the `transit/keys/<key_ring_name>/rotate`
endpoint.

```plaintext
$ vault write -f transit/keys/orders/rotate
```

Let's encrypt another data:

```plaintext
$ vault write transit/encrypt/orders plaintext=$(base64 <<< "visa-card-number")
Key           Value
---           -----
ciphertext    vault:v2:45f9zW6cglbrzCjI0yCyC6DBYtSBSxnMgUn9B5aHcGEit71xefPEmmjMbrk3
```

Compare the ciphertexts from [Step 2](#step2).  

```
ciphertext    vault:v1:cZNHVx+sxdMErXRSuDa1q/pz49fXTn1PScKfhf+PIZPvy8xKfkytpwKcbC0fF2U=
```

Notice that the first ciphertext starts with "**`vault:v1:`**".  After rotating
the encryption key, the ciphertext starts with "**`vault:v2:`**".  This indicates
that the data gets encrypted using the latest version of the key after the
rotation.


Execute the following command to rewrap your cipertext from [Step 2](#step2)
with the latest version of the encryption key:

```plaintext
$ vault write transit/rewrap/orders \
        ciphertext="vault:v1:cZNHVx+sxdMErXRSuDa1q/pz49fXTn1PScKfhf+PIZPvy8xKfkytpwKcbC0fF2U="
Key           Value
---           -----
ciphertext    vault:v2:kChHZ9w4ILRfw+DzO53IZ8m5PyB2yp2/tKbub34uB+iDqtDRB+NLCPrpzTtJHJ4=
```

Notice that the resulting ciphertext now starts with "`vault:v2:`".  

This operation does not reveal the plaintext data. But Vault will decrypt the
value using the appropriate key in the keyring and then encrypted the resulting
plaintext with the newest key in the keyring.



#### API call using cURL

To rotate the encryption key, invoke the `transit/keys/<key_ring_name>/rotate`
endpoint.

```plaintext
$ curl --header "X-Vault-Token: ..." \
       --request POST
       https://127.0.0.1:8200/v1/transit/keys/orders/rotate
```

Let's encrypt another data:

```shell
# Generate base64-encoded plaintext
$ base64 <<< "visa-card-number"
dmlzYS1jYXJkLW51bWJlcgo=

# Pass the base64-encoded plaintext in the request payload
$ curl --header "X-Vault-Token: ..." \
       --request POST \
       --data '{"plaintext": "dmlzYS1jYXJkLW51bWJlcgo="}' \
       https://127.0.0.1:8200/v1/transit/encrypt/orders | jq
{
  ...
  "data": {
    "ciphertext": "vault:v2:et873RqkfLlS268LqYspVUnqhqZm0flNwhthe4ZzfcuZQab1TnirQ8/hMNYA"
  },
  ...
}
```

Compare the ciphertexts from [Step 2](#step2).  

```
ciphertext    vault:v1:cZNHVx+sxdMErXRSuDa1q/pz49fXTn1PScKfhf+PIZPvy8xKfkytpwKcbC0fF2U=
```

Notice that the first ciphertext starts with "**`vault:v1:`**".  After rotating
the encryption key, the ciphertext starts with "**`vault:v2:`**".  This indicates
that the data gets encrypted using the latest version of the key after the
rotation.


Execute the `transit/rewrap` endpoint to rewrap your cipertext from [Step 2](#step2)
with the latest version of the encryption key:

```plaintext
$ curl --header "X-Vault-Token: ..." \
       --request POST \
       --data '{"ciphertext": "vault:v1:/9hdQutaWpZR72s3+VSCLK1JNhV1wKM49hYVjh7RjmuxIy/OvshtgV4L4uVB+aQ="}' \
       https://127.0.0.1:8200/v1/transit/rewrap/orders | jq
{
 ...
 "data": {
   "ciphertext": "vault:v2:ykqXDP65tLVSrqxoNZh51gIobYQSwNGT+SbiD/2nl8rrhF2md+wplBGdXlhDzd4="
 },
 ...
```

Notice that the resulting ciphertext now starts with "`vault:v2:`".  

This operation does not reveal the plaintext data. But Vault will decrypt the
value using the appropriate key in the keyring and then encrypted the resulting
plaintext with the newest key in the keyring.


## <a name="step5"></a>Step 5: Update Key Configuration
(**Persona:** operator)

The operators can [update the encryption key
configuration](/api/secret/transit/index.html#update-key-configuration) to
specify the minimum version of ciphertext allowed to be decrypted, the minimum
version of the key that can be used to encrypt the plaintext, the key is allowed
to be deleted, etc.

This helps further tightening the data encryption rule.


#### CLI Command

Execute the key rotation command a few times to generate multiple versions of
the key:

```plaintext
$ vault write -f transit/keys/orders/rotate
```

Now, read the `orders` key information:

```plaintext
$ vault read transit/keys/orders

Key                       Value
---                       -----
...
keys                      map[6:1531439714 1:1531439594 2:1531439667 3:1531439714 4:1531439714 5:1531439714]
latest_version            6
min_decryption_version    1
min_encryption_version    0
...
```

In the example, the current version of the key is **6**. However, there is no
restriction about the minimum encryption key version, and any of the key
versions can decrypt the data (`min_decryption_version`).

Run the following command to enforce the use of the encryption key at version
**5** or later to decrypt the data.

```plaintext
$ vault write transit/keys/orders/config min_decryption_version=5
```

Now, verify the `orders` key configuration:

```plaintext
$ vault read transit/keys/orders

Key                       Value
---                       -----
allow_plaintext_backup    false
deletion_allowed          false
derived                   false
exportable                false
keys                      map[5:1531811719 6:1531811721]
latest_version            6
min_decryption_version    5
min_encryption_version    0
...
```


#### API call using cURL

Execute the `transit/keys/<key_ring_name>/rotate` endpoint a few times key
rotation command a few times to generate multiple versions of the key:

```plaintext
$ curl --header "X-Vault-Token: ..." \
       --request POST
       https://127.0.0.1:8200/v1/transit/keys/orders/rotate
```

Read the `transit/keys/orders` endpoint to review the `orders` key
detail:

```plaintext
$ curl --header "X-Vault-Token: ..." \
       https://127.0.0.1:8200/v1/transit/keys/orders | jq
{
  ...
   "keys": {
     "1": 1531804669,
      "2": 1531810236,
      "3": 1531811712,
      "4": 1531811715,
      "5": 1531811719,
      "6": 1531811721
   },
   "latest_version": 6,
   "min_decryption_version": 1,
   "min_encryption_version": 0,
   ...
```

In the example, the current version of the key is **6**. However, there is no
restriction about the minimum encryption key version, and any of the key
versions can decrypt the data (`min_decryption_version`).

Run the following command to enforce the use of the encryption key at version
**5** or later to decrypt the data.

```plaintext
$ curl --header "X-Vault-Token: ..." \
       --request POST
       --data '{"min_decryption_version": 5}'
       https://127.0.0.1:8200/v1/transit/keys/orders/config
```

Now, verify the `orders` key configuration:

```plaintext
$ curl --header "X-Vault-Token: ..." \
       https://127.0.0.1:8200/v1/transit/keys/orders | jq
{
  ...
  "keys": {
     "5": 1531811719,
     "6": 1531811721
   },
   "latest_version": 6,
   "min_decryption_version": 5,
   "min_encryption_version": 0,
   ...
```

<br>

-> **NOTE:** Notice that the output only displays two valid encryption key
versions (`5` and `6`).



## Next steps

[Transit Secrets Re-wrapping](/guides/encryption/transit-rewrap.html) guide
introduces a sample application which re-wraps data after rotating an encryption
key in the transit engine in Vault.
