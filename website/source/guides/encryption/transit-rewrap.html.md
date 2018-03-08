---
layout: "guides"
page_title: "Transit Secrets Re-rapping - Guides"
sidebar_current: "guides-encryption-rewrap"
description: |-
  The goal of this guide is to demonstrate one possible way to re-wrap data after
  rotating an encryption key in the transit engine in Vault.
---

# Transit Secrets Engine

In addition to being able to store secrets, Vault can encrypt/decrypt data that
is stored elsewhere. The primary use of this is to allow applications to encrypt
their data while still storing it in their primary data store. Vault does not
store the data.

The [`transit` secret engine](/docs/secrets/transit/index.html) handles
cryptographic functions on data-in-transit, and often referred to as
**_Encryption as a Service_** (EaaS). Both small amounts of arbitrary data, and
large files such as images, can be protected with the transit engine.  This EaaS
function can augment or eliminate the need for Transparent Data Encryption (TDE)
with databases to encrypt the contents of a bucket, volume, and disk, etc.  

![Encryption as a Service](/assets/images/vault-encryption.png)

## Encryption Key Rotation

One of the benefits of using the Vault EaaS is its ability to easily rotate the
encryption keys. Keys can be rotated manually by a human, or an automated
process which invokes the key rotation API endpoint through `cron`, a CI
pipeline, a periodic Nomad batch job, Kubernetes Job, etc.

The goal of this guide is to demonstrate an example for re-wrapping data after
rotating an encryption key in the transit engine in Vault.


## Reference Material

- [Transit Secret Engine](/docs/secrets/transit/index.html)
- [Transit Secret Engine API](/api/secret/transit/index.html)
- [Transparent Data Encryption in the Modern Datacenter](https://www.hashicorp.com/blog/transparent-data-encryption-in-the-modern-datacenter)


## Estimated Time to Complete

30 minutes


## Personas

The end-to-end scenario described in this guide involves two personas:

- **security engineer** with privileged permissions to manage the encryption keys
- **app** with un-privileged permissions rewraps secrets via API


## Challenge

Vault maintains the versioned keyring and the operator can decide
the minimum version allowed for decryption operations.  When data is
encrypted using Vault, the resulting ciphertext is prepended with the version of
the key used to encrypt it.  

The following example shows data that was encrypted using the fourth version of
a particular encryption key:

```
vault:v4:ueizdCqCJ/YhowQSvmJyucnLfIUMd4S/nLTpGTcz64HXoY69dwOrqerFzOlhqg==
```

For example, an organization could decide that a key should be rotated _once a
week_, and that the minimum version allowed to decrypt records is the current
version as well as the previous two versions.  If the current version is five,
then Vault would decrypt records that were sent to it with the following
prefixes:

- vault:**v5**:lkjasfdlkjafdlkjsdflajsdf==
- vault:**v4**:asdfas9pirapirteradr33vvv==
- vault:**v3**:ouoiujarontoiue8987sdjf^1==

In this example, what would happen if you send Vault data that was encrypted
with the first or second version of the key (`vault:v1:...` or `vault:v2:...`)?  

Vault would refuse to decrypt the data as the key used is less than the minimum
key version allowed.


## Solution

Luckily, Vault provides an easy way of re-wrapping encrypted data when a key is
rotated.  Using the rewrap API endpoint, a non-privileged Vault entity can send
data encrypted with an older version of the key to have it re-encrypted with the
latest version. The application performing the re-wrapping never interacts with
the decrypted data. The process of rotating the encryption key and rewrapping
records could (and should) be completely **automated**. Records could be updated
slowly over time to lessen database load, or all at once at the time of
rotation.  The exact implementation will depend heavily on the needs of each
particular organization or application.


## Prerequisites

To perform the tasks described in this guide, you need to have a Vault
environment.  Refer to the [Getting
Started](/intro/getting-started/install.html) guide to install Vault. Make sure
that your Vault server has been [initialized and
unsealed](/intro/getting-started/deploy.html).

The following tools are required in order to successfully run the sample
application provided in this guide:

- [.NET Core](https://www.microsoft.com/net/download)
- [Docker](https://docs.docker.com/install/)

Download the sample application code from
[vault-guides](https://github.com/hashicorp/vault-guides/tree/master/secrets/transit/vault-transit-rewrap-example)
repository to perform the steps described in this guide.

The `vault-transit-rewrap-example` contains the following:

```bash
.
├── AppDb.cs
├── DBHelper.cs
├── Program.cs
├── README.md
├── Record.cs
├── VaultClient.cs
├── WebHelper.cs
└── rewrap_example.csproj
```


### Policy requirements

-> **NOTE:** For the purpose of this guide, you can use the **`root`** token to work
with Vault. However, it is recommended that root tokens are only used for just
enough initial setup or in emergencies. As a best practice, use tokens with
an appropriate set of policies based on your role in the organization.

To perform all tasks demonstrated in this guide, your policy must include the
following permissions:

```shell
# Manage transit secret engine
path "transit/keys/*" {
  capabilities = [ "create", "read", "update", "delete", "list", "sudo" ]
}

# Enable transit secret engine
path "sys/mounts/transit" {
  capabilities = [ "create", "update" ]
}

# Write ACL policies
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

This guide introduces a sample _.Net_ application which automates the
re-wrapping of the data using the latest encryption key.

For the purpose of this guide, a MySQL database runs locally using Docker.
However, these steps would work for an existing MySQL database by supplying the
proper network information to your environment.

You are going to perform the following steps:

1. [Test database setup (Docker)](#step1)
1. [Enable the transit secret engine](#step2)
1. [Generate a new token for sample app](#step3)
1. [Run the sample application](#step4)
1. [Rotate the encryption keys](#step5)
1. [Re-wrapping data programmatically](#step6)


### <a name="step1"></a>Step 1: Test database setup (Docker)

You need a database to test with.  You can create one to test with easily using
Docker:

```bash
# Pull the latest mysql container image
docker pull mysql/mysql-server:5.7

# Create a directory for our data (change the following line if running on Windows)
mkdir ~/rewrap-data

# Run the container.  The following command creates a database named 'my_app',
# specifies the root user password as 'root', and adds a user named vault
docker run --name mysql-rewrap \
  -p 3306:3306 \
  -v ~/rewrap-data/var/lib/mysql \
  -e MYSQL_ROOT_PASSWORD=root \
  -e MYSQL_ROOT_HOST=% \
  -e MYSQL_DATABASE=my_app \
  -e MYSQL_USER=vault \
  -e MYSQL_PASSWORD=vaultpw \
  -d mysql/mysql-server:5.7
```

### <a name="step2"></a>Step 2: Enable the transit secret engine
(**Persona:** security engineer)

#### CLI command

Enable the `transit` secret engine by executing the following command:

```bash
$ vault secrets enable transit
```

Create an encryption key to use for transit named, "my_app_key".

```bash
$ vault write -f transit/keys/my_app_key
```

#### API call using cURL

Enable the `transit` secret engine via API, use the `/sys/mounts` endpoint:

```bash
$ curl --header "X-Vault-Token: <TOKEN>" \
       --request POST \
       --data <PARAMETERS> \
       <VAULT_ADDRESS>/v1/sys/mounts/transit
```

Where `<TOKEN>` is your valid token, and `<PARAMETERS>` holds [configuration
parameters](/api/system/mounts.html#enable-secrets-engine) of the secret engine.

To crate a new encryption key, use the `transit/keys/<key_name>` endpoint:

```bash
$ curl --header "X-Vault-Token: <TOKEN>" \
       --request POST \
       --data <PARAMETERS> \
       <VAULT_ADDRESS>/v1/transit/keys/<KEY_NAME>
```

Where `<PARAMETERS>` holds [configuration
parameters](/api/secret/transit/index.html#create-key) to specify the key type.

**Example:**

```shell
$ curl --header "X-Vault-Token: ..." \
       --request POST \
       --data '{"type": "transit"}' \
       https://localhost:8200/v1/sys/mounts/transit
```

The above example passes the **type** (`transit`) in the request payload which
at the `sys/mounts/transit` endpoint.


Next, create an encryption key to use for transit named, "my_app_key".

```bash
$ curl --header "X-Vault-Token: ..." \
       --request POST \
       https://localhost:8200/v1/transit/keys/my_app_key
```


### <a name="step3"></a>Step 3: Generate a new token for sample app
(**Persona:** security engineer)

Before generating a token, create a limited scope policy named, "**rewrap_example**"
for the sample application.  

The ACL policy (`rewrap_example.hcl`) looks as follows:

```shell
path "transit/keys/my_app_key" {
  capabilities = ["read"]
}

path "transit/rewrap/my_app_key" {
  capabilities = ["update"]
}

# This last policy is needed to seed the database as part of the example.  
# It can be omitted if seeding is not required
path "transit/encrypt/my_app_key" {
  capabilities = ["update"]
}
```

#### CLI command

Create `rewrap_example` policy:

```shell
$ vault policy write rewrap_example ./rewrap_example.hcl
```

Finally, create a token to use the `rewrap_example` policy:

```shell
$ vault token create -policy=rewrap_example
```

**Example:**

```shell
$ vault token create -policy=rewrap_example
Key                Value
---                -----
token              68396128-82d8-002e-f289-1106944fee9f
token_accessor     75f05f43-6a5f-2eb1-5bb8-0de3c7cf0996
token_duration     768h
token_renewable    true
token_policies     [default rewrap_example]
```

The generated token is what the sample application uses to connect to Vault.


#### API call using cURL

To create a policy via API, use the `/sys/policy` endpoint:

```plaintext
$ curl --request PUT --header "X-Vault-Token: ..." \
       --data @payload.json \
       https://localhost:8200/v1/sys/policy/rewrap_example

$ cat payload.json
{
  "policy": "path \"transit/keys/my_app_key\" { capabilities = [\"read\"] } path \"transit/rewrap/my_app_key\" ... }"
}
```

Finally, create a token to use the `rewrap_example` policy:

```plaintext
$ curl --header "X-Vault-Token: ..." --request POST  \
       --data '{ "policies": ["rewrap_example"] }' \
       https://localhost:8200/v1/auth/token/create | jq
{
 "request_id": "da8bde73-99ab-b435-a344-fb963b3a599f",
 "lease_id": "",
 "renewable": false,
 "lease_duration": 0,
 "data": null,
 "wrap_info": null,
 "warnings": null,
 "auth": {
   "client_token": "997107d5-9049-8b4b-8f39-a33a458fd02d",
   "accessor": "d400076b-4143-4d63-7473-a8cc52c73ba3",
   "policies": [
     "default",
     "rewrap-example"
   ],
   "metadata": null,
   "lease_duration": 2764800,
   "renewable": true,
   "entity_id": ""
 }
}
```

The generated token is what the sample application uses to connect to Vault.


### <a name="step4"></a>Step 4: Run the sample application
(**Persona:** app)

You are now ready to run the app. Be sure to [download](#prerequisites) the
sample application code before beginning.

**Sample application**

| File                   | Description                                       |
|------------------------|---------------------------------------------------|
| Program.cs             | Starting point of this sample app (the `Main()` method) is in this file.  It reads the environment variable values, connects to Vault and the MySQL database.  If the `user_data` table does not exist, it creates it.   |
| DBHelper.cs            | Defines a method to create the `user_data` table if it does not exist. Finds and updates records that need to be rewrapped with the new key.    |
| AppDb.cs               | Connects to the MySQL database.                   |
| Record.cs              | Sample data record template.                      |
| VaultClient.cs         | Defines methods necessary to rewrap transit data. |
| WebHelper.cs           | Helper code to seed the initial table schema.     |
| rewrap_example.csproj  | Project file for this sample app.                 |


The sample app retrieves the user token, Vault address, and the name of the
transit key through environment variables. Be sure to supply the token created
in [Step 3](#step3):

```bash
$ VAULT_TOKEN=<APP_TOKEN> \
     VAULT_ADDR=<VAULT_ADDRESS> \
     VAULT_TRANSIT_KEY=my_app_key \
     SHOULD_SEED_USERS=true \
     dotnet run
```

> If you need to seed test data you can do so by including the
`SHOULD_SEED_USERS=true`.  

**Example:**

```bash
$ VAULT_TOKEN=$TOKEN VAULT_ADDR=http://localhost:8200 VAULT_TRANSIT_KEY=my_app_key SHOULD_SEED_USERS=true dotnet run

Connecting to Vault server...
Created (if not exist) my_app DB
Create (if not exist) user_data table
Seeded the database...
Moving rewrap...
Current Key Version: 5
Found 0 records to rewrap.
```

You can inspect the contents of the database with:

```bash
$ docker exec -it mysql-rewrap mysql -uroot -proot
...
mysql> DESC user_data;
mysql> SELECT * FROM user_data WHERE dob LIKE "vault:v1%" limit 10;
...data...
```

### <a name="step5"></a>Step 5: Rotate the encryption keys
(**Persona:** security engineer)

The encryption key (`my_app_key`) can be rotated easily.

#### CLI command

To rotate the key, you write to the `transit/keys/<KEY_NAME>/rotate` path.

```bash
$ vault write -f transit/keys/my_app_key/rotate
Success! Data written to: transit/keys/my_app_key/rotate
```

Run the command a few times to generate several versions of the encryption key
for testing.

To view the key information:

```bash
$ vault read transit/keys/my_app_key
Key                       Value
---                       -----
allow_plaintext_backup    false
deletion_allowed          false
derived                   false
exportable                false
keys                      map[5:1519623974 6:1519623980 1:1519620952 2:1519623255 3:1519623285 4:1519623603]
latest_version            6
min_decryption_version    1
min_encryption_version    0
name                      my_app_key
supports_decryption       true
supports_derivation       true
supports_encryption       true
supports_signing          false
type                      aes256-gcm96
```

You can see that in the above example the current version of the key is six.
There is no restriction about a minimum encryption key version, and any of the key
versions can decrypt the data (`min_decryption_version`).

Let's enforce the use of the encryption key at version five or later to decrypt
data.

```shell
# replace '5' with the appropriate version
$ vault write transit/keys/my_app_key/config min_decryption_version=5

# Verify the changes were successful
$ vault read transit/keys/my_app_key
Key                       Value
---                       -----
allow_plaintext_backup    false
deletion_allowed          false
derived                   false
exportable                false
keys                      map[5:1519623974 6:1519623980]
latest_version            6
min_decryption_version    5
min_encryption_version    0
name                      my_app_key
supports_decryption       true
supports_derivation       true
supports_encryption       true
supports_signing          false
type                      aes256-gcm96
```


#### API call using cURL

To rotate the encryption key via API, use the `transit/keys/<KEY_NAME>/rotate` endpoint:

```bash
$ curl --request POST --header "X-Vault-Token: ..." \
       https://localhost:8200/v1/transit/keys/my_app_key/rotate
```

Run the command a few times to generate several versions of the encryption key
for testing.

```shell
# Verify the changes were successful
$ curl --request GET --header "X-Vault-Token: ..." \
      https://localhost:8200/v1/transit/keys/my_app_key | jq
{
  "request_id": "ed13436a-4816-2f51-0552-6a001e823548",
  "lease_id": "",
  "renewable": false,
  "lease_duration": 0,
  "data": {
    "allow_plaintext_backup": false,
    "deletion_allowed": false,
    "derived": false,
    "exportable": false,
    "keys": {
      "1": 1519620952,
      "2": 1519623255,
      "3": 1519623285,
      "4": 1519623603,
      "5": 1519623974,
      "6": 1519623980
    },
    "latest_version": 6,
    "min_decryption_version": 1,
    "min_encryption_version": 0,
    "name": "my_app_key",
    ...
  },
  ...
}
```

You can see that in the above example the current version of the key is six.
There is no restriction about the minimum encryption key version, and any of the key
versions can decrypt the data (`min_decryption_version`).

Let's enforce the use of the encryption key at version five or later to decrypt the
data.

```shell
$ curl --request POST --header "X-Vault-Token: ..." \
       --data '{ "min_decryption_version": 5 }'
       https://localhost:8200/v1/transit/keys/my_app_key/config

# Verify the changes were successful
$ curl --request GET --header "X-Vault-Token: ..." \
       https://localhost:8200/v1/transit/keys/my_app_key | jq
{
 ...
 "data": {
   ...
   "keys": {
     "5": 1519623974,
     "6": 1519623980
   },
   "latest_version": 6,
   "min_decryption_version": 5,
   "min_encryption_version": 0,
   "name": "my_app_key",
   ...
  },
}
```


### <a name="step6"></a>Step 6: Programmatically re-wrap the data
(**Persona:** app)

Now you have records in the database and you have updated our minimum key
version. You can run the application again and should see it update records as
appropriate. Remember you can inspect records using the MySQL shell (see above).

**Example:**

```shell
$ VAULT_TOKEN=2616214b-6868-3589-b443-0330d7915882 VAULT_ADDR=http://localhost:8200 \
          VAULT_TRANSIT_KEY=my_app_key SHOULD_SEED_USERS=true dotnet run
Connecting to Vault server...
Created (if not exist) my_app DB
Create (if not exist) user_data table
Seeded the database...
Current Key Version: 6
Found 3500 records to rewrap.
Wrapped another 10 records: 10 so far...
Wrapped another 10 records: 20 so far...
Wrapped another 10 records: 30 so far...
...
```

#### Validation

The application has now re-wrapped all records with the latest key.  You can
verify this by running the application again, or by inspecting the records using the
MySQL client.

```bash
$ docker exec -it mysql-rewrap mysql -uroot -proot
...
mysql> DESC user_data;
mysql> SELECT * FROM user_data WHERE dob LIKE "vault:v1%" limit 10;
Empty set (0.00 sec)

mysql> SELECT * FROM user_data WHERE dob LIKE "vault:v6%" limit 10;
...data...
```

### Conclusion

An application similar to this could be scheduled via cron, run periodically as
a [Nomad batch
job](https://www.nomadproject.io/docs/job-specification/periodic.html), or
executed in a variety of other ways.  You could also modify it to re-wrap a
limited number of records at a time so as to not put undue strain on the
database.  The final implementation should be based upon the needs and design
goals specific to each organization or application.  


## Next Steps

Since the main focus of this guide was to programmatically rewrap your secrets
using the latest encryption key, the token used by the sample application was
generated manually. In a production environment, you'll want to pass the token in
a more secure manner.  Refer to the [Cubbyhole Response
Wrapping](/guides/secret-mgmt/cubbyhole.html) guide to wrap the token so that only the
expecting app can unwrap to obtain the token.

Also, refer to the [AppRole Pull
Authentication](/guides/identity/authentication.html) to generate tokens for
apps using the AppRole auth method.
