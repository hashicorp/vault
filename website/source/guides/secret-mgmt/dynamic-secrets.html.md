---
layout: "guides"
page_title: "Secret as a Service - Guides"
sidebar_current: "guides-secret-mgmt-dataynamic-secrets"
description: |-
  Vault can dynamically generate secrets on--dataemand for some systems.
---

# Secret as a Service: Dynamic Secrets

Vault can generate secrets on-demand for some systems. For example, when an app
needs to access an Amazon S3 bucket, it asks Vault for AWS credentials. Vault
will generate an AWS credential granting permissions to access the S3 bucket. In
addition, Vault will automatically revoke this credential after the TTL is
expired.

The [Getting Started](/intro/getting-started/dynamic-secrets.html) guide walks
you through the generation of dynamic AWS credentials.

## Reference Material

- [Getting Started - Dynamic Secrets](/intro/getting-started/dynamic-secrets.html)
- [Database Backends](/docs/secrets/databases/index.html)
- [Role API](/api/secret/databases/index.html#create-role)

## Estimated Time to Complete

10 minutes

## Personas

The end-to-end scenario described in this guide involves two personas:

- **`admin`** with privileged permissions to configure secret backends
- **`apps`** read the secrets from Vault

## Challenge

Data protection is a top priority which means that the database credential
rotation is a critical part of any data protection initiative. Each role has a
different set of permissions granted to access the database. When a system is
attacked by hackers, continuous credential rotation becomes necessary and needs
to be automated.


## Solution

Applications ask Vault for database credential rather than setting them as
environment variables. The administrator specifies the TTL of the database
credentials to enforce its validity so that they are automatically revoked when
they are no longer used.

![Dynamic Secret Workflow](/assets/images/vault-dynamic-secrets.png)

Each app instance can get unique credentials that they don't have to share. By
making those credentials to be short-lived, you reduced the change of the secret
to being compromised. If an app was compromised, the credentials used by the app
can be revoked rather than changing more global set of credentials.

## Prerequisites

To perform the tasks described in this guide, you need to have a Vault
environment.  Refer to the [Getting
Started](/intro/getting-started/install.html) guide to install Vault. Make sure
that your Vault server has been [initialized and
unsealed](/intro/getting-started/deploy.html).

### PostgreSQL

This guide requires that you have PostgreSQL that you can connect to,
and have a database named **`myapp`**. You can download and install
[PostgreSQL](https://www.postgresql.org/download/) locally, or connect to a
remote host.

**Example on Ubuntu:**

```shell
# Install PostgreSQL
$ sudo apt-get install -y postgresql postgresql-contrib

# Initialize PostgreSQL
$ sudo postgresql-setup initdb

# Switch to postgres user
$ su - postgres

# Create myapp database
$ psql -U postgres -c 'CREATE DATABASE myapp;'
```

[PostgreSQL Wiki](https://wiki.postgresql.org/wiki/First_steps) gives you a
summary of basic commands to get started.

### Policy requirements

-> **NOTE:** For the purpose of this guide, you can use **`root`** token to work
with Vault. However, it is recommended that root tokens are only used for just
enough initial setup or in emergencies. As a best practice, use tokens with
appropriate set of policies based on your role in the organization.

To perform all tasks demonstrated in this guide, your policy must include the
following permissions:

```shell
# Mount secret backends
path "sys/mounts/*" {
  capabilities = [ "create", "read", "update", "delete", "list" ]
}

# Configure the database backend and create roles
path "database/*" {
  capabilities = [ "create", "read", "update", "delete", "list" ]
}

# Write ACL policies
path "sys/policy/*" {
  capabilities = [ "create", "read", "update", "delete", "list" ]
}

# Manage tokens for verification
path "auth/token/create" {
  capabilities = [ "create", "read", "update", "delete", "list", "sudo" ]
}
```

If you are not familiar with policies, complete the
[policies](/guides/configuration/policies.html) guide.

## Steps

In this guide, you are going to configure PostgreSQL secret backend, and create
a read-only database role. The Vault generated PostgreSQL credentials will only
have read permission.

1. [Mount the database secret backend](#step1)
2. [Configure PostgreSQL backend](#step2)
3. [Create a role](#step3)
4. [Request PostgreSQL credentials](#step4)
5. [Validation](#validation)

Step 1 through 3 need to be performed by an `admin` user.  Step 4 describes
the commands that an `app` runs to get a database credentials from Vault.


### <a name="step1"></a>Step 1: Mount the database secret backend
(**Persona:** admin)

As most of the secret backends, the [database backend](/docs/secrets/databases/index.html)
must be mounted.

#### CLI command

To mount a database backend:

```shell
$ vault secrets enable <PATH>
```

**Example:**

```shell
$ vault secrets enable database
```

**NOTE:** In this guide, the database backend is mounted at the `/database path` in
Vault.  However, it is possible to mount your secret backends at any location.

#### API call using cURL

Mount `database` secret backend using `/sys/mounts` endpoint:

```shell
$ curl --header "X-Vault-Token: <TOKEN>" \
       --request POST \
       --data <PARAMETERS> \
       <VAULT_ADDRESS>/v1/sys/mounts/<PATH>
```

Where `<TOKEN>` is your valid token, and `<PARAMETERS>` holds [configuration
parameters](/api/system/mounts.html#mount-secret-backend) of the backend.

**Example:**

The following example mounts database backend at `sys/mounts/database` path, and
passed the backend type ("database") in the request payload.

```shell
$ curl --header "X-Vault-Token: ..." \
       --request POST \
       --data '{"type":"database"}' \
       https://$ vault.rocks/v1/sys/mounts/database
```

**NOTE:** It is possible to mount your database secret backends at any location.


### <a name="step1"></a>Step 2: Configure PostgreSQL backend
(**Persona:** admin)

The PostgreSQL backend needs to be configured with valid credentials. It is very
common to give Vault the **root** credentials and let Vault manage the auditing
and lifecycle credentials; it's much better than having one person manage the
credentials.


The following command configures the database secret backend using
`postgresql-database-plugin` where the database connection URL is
`postgresql://root:rootpassword@localhost:5432/myapp`.  The allowed role is
`readonly` which you will create in [Step 3](#step3).

**NOTE:** If your
database connection URL is different from this example, be sure to replace the
command with correct URL to match your environment.


#### CLI command

**Example:**

```shell
$ vault write database/config/postgresql plugin_name=postgresql-database-plugin \
  allowed_roles=readonly connection_url=postgresql://root:rootpassword@localhost:5432/myapp
```


#### API call using cURL

**Example:**

```shell
$ curl --header "X-Vault-Token: ..." --request POST --data @payload.json \
    https://vault.rocks/v1/database/config/postgresql

$ cat payload.json
{
	"plugin_name": "postgresql-database-plugin",
	"allowed_roles": "readonly",
	"connection_url": "postgresql://root:rootpassword@localhost:5432/myapp"
}
```

### <a name="step3"></a>Step 3: Create a role
(**Persona:** admin)

In [Step 2](#step2), you configured the PostgreSQL backend by passing **`readonly`** role
as an allowed member. The next step is to define the `readonly` role. A role is
a logical name that maps to a policy used to generate credentials.

-> Vault does not know what kind of PostgreSQL users you want to create. So,
supply the information in SQL to create desired users.

**Example:** `readonly.sql`

```plaintext
CREATE ROLE "{{name}}" WITH LOGIN PASSWORD '{{password}}' VALID UNTIL '{{expiration}}';
GRANT SELECT ON ALL TABLES IN SCHEMA public TO "{{name}}";
```

The values within the `{{<value>}}` will be filled in by Vault. Notice that
**`VALID_UNTIL`** clause. This tells PostgreSQL to revoke the credentials even
if Vault is offline or unable to communicate with it.


#### CLI command

**Example:**

```shell
$ vault write database/roles/readonly db_name=postgresql creation_statements=@readonly.sql \
    default_ttl=1h max_ttl=24h
```

The above command creates a role named, `readonly` with default TTL of 1
hour, and max TTL of the credential is set to 24 hours. The `readonly.sql`
statement is passed as the role creation statement.

#### API call using cURL

**Example:**

```shell
$ curl --header "X-Vault-Token: ..." --request POST --data @payload.json \
    https://vault.rocks/v1/database/roles/readonly

$ cat payload.json
{
	"db_name": "postgres",
	"creation_statements": "CREATE ROLE '{{name}}' WITH LOGIN PASSWORD '{{password}}' VALID UNTIL '{{expiration}}';
  GRANT SELECT ON ALL TABLES IN SCHEMA public TO '{{name}}';",
	"default_ttl": "1h",
	"max_ttl": "24h"
}
```

The `db_name`, `creation_statements`, `default_ttl`, and `max_ttl` are set in
the `role-payload.json`.


### <a name="step4"></a>Step 4: Request PostgreSQL credentials
(**Persona:** apps)

Now, you are switching to [`apps` persona](#personas). To get a new set of
PostgreSQL credentials, the client app needs to be able to **read** from the
`readonly` role endpoint. Therefore, the app's token must have a policy granting
the read permission.

`apps-policy.hcl`

```shell
# Get credentials from the database backend
path "database/creds/readonly" {
  capabilities = [ "read" ]
}
```

#### CLI command

First create an `apps` policy, and generate a token so that you can authenticate
as an `apps` persona.

**Example:**

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

Use the returned token to perform the remaining.

**NOTE:** [AppRole Pull Authentication](/guides/configuration/authentication.html) guide
demonstrates more sophisticated way of generating a token for your apps.

```shell
# Authenticate with Vault using the generated token first
$ vault login e4bdf7dc-cbbf-1bb1-c06c-6a4f9a826cf2
Successfully authenticated! You are now logged in.
token: e4bdf7dc-cbbf-1bb1-c06c-6a4f9a826cf2
token_duration: 2764277
token_policies: [apps default]

# Invoke the vault command
$ vault read database/creds/readonly

Key            	Value
---            	-----
lease_id       	database/creds/readonly/4b5c6e82--dataf88-0dec-c0cb-f07eee8f0329
lease_duration 	1h0m0s
lease_renewable	true
password       	A1a-4urzp0wu92r5s1q0
username       	v-token-readonly-9x3qrw452wwz4w6421xt-1515625519
```

**NOTE:** Re-run the command and notice that Vault returns a different set of
credentials each time. This means that each app instance can acquire a unique
set of credentials.


#### API call using cURL

First create an `apps` policy, and generate a token so that you can authenticate
as an `app` persona.

```shell
# Payload to pass in the API call
$ cat payload.json
{
  "policy": "path \"database/creds/readonly\" {capabilities = [ \"read\" ]}"
}

# Create "apps" policy
$ curl --header "X-Vault-Token: ..." --request PUT \
       --data @payload.json \
       https://vault.rocks/v1/sys/policy/apps

# Generate a new token with apps policy
$ curl --header "X-Vault-Token: ..." --request POST \
       --data '{"policies": ["apps"]}' \
       https://vault.rocks/v1/auth/token/create | jq
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

Be sure to use the returned token to perform the remaining.

**NOTE:** [AppRole Pull Authentication](/guides/configuration/authentication.html) guide
demonstrates more sophisticated way of generating a token for your apps.

```shell
$ curl --header "X-Vault-Token: 1c97b03a-6098-31cf-9d8b-b404e52dcb4a" \
       --request GET \
       https://vault.rocks/v1/database/creds/readonly | jq
{
  "request_id": "e0e5a6c1-5e69-5cf3-c9d2-020af192de36",
  "lease_id": "database/creds/readonly/7aa462ab-98cb-fdcb-b226-f0a0d37644cc",
  "renewable": true,
  "lease_duration": 3600,
  "data": {
    "password": "A1a-2680ut032xqt16tq",
    "username": "v-token-readonly-6s4su6z93472x0r2787t-1515625742"
  },
  "wrap_info": null,
  "warnings": null,
  "auth": null
}
```

### Validation

(1) Generate a new set of credentials.

```plaintext
$ vault read database/creds/readonly

Key            	Value
---            	-----
lease_id       	database/creds/readonly/3e8174da-6ca0-143b-aa8c-4c238aa02809
lease_duration 	1h0m0s
lease_renewable	true
password       	A1a-w2xv2zsq4r5ru940
username       	v-token-readonly-48rt0t36sxp4wy81x8x1-1515627434
```

The generated username is `v-token-readonly-48rt0t36sxp4wy81x8x1-1515627434`.

(2) Connect to the postgres as an admin user, and run the following psql commands.

```plaintext
$ psql -U postgres

postgres > \du
                                                       List of roles
                    Role name                     |                         Attributes                         | Member of
--------------------------------------------------+------------------------------------------------------------+-----------
 postgres                                         | Superuser, Create role, Create DB, Replication, Bypass RLS | {}
 v-token-readonly-48rt0t36sxp4wy81x8x1-1515627434 | Password valid until 2018-01-11 00:37:14+00                | {}

postgres > \q
 ```

The `\du` command lists all users. You should be able to verify that the username generated by Vault exists.

(3) Renew the lease for this credential by passing its **`lease_id`**.

```plaintext
$ vault renew database/creds/readonly/3e8174da-6ca0-143b-aa8c-4c238aa02809

Key            	Value
---            	-----
lease_id       	database/creds/readonly/3e8174da-6ca0-143b-aa8c-4c238aa02809
lease_duration 	1h0m0s
lease_renewable	true
```

(4) Revoke the generated credentials.

```plaintext
$ vault lease revoke database/creds/readonly/3e8174da-6ca0-143b-aa8c-4c238aa02809
```

**NOTE:** If you run the command with **`-prefix=true`** flag, it revokes all
secrets under `database/creds/readonly`.

Now, when you check the list of users in PostgreSQL, none of the Vault generated
user name exists.


## Next steps

This guide discussed how to generate credentials on--dataemand so that the access
credentials no longer need to be written to disk. Next, learn about the
[Tokens and Leases](/guides/configuration/lease.html) so that you can control the lifecycle of
those credentials.
