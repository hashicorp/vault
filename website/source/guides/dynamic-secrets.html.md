---
layout: "guides"
page_title: "Secret as a Service - Guides"
sidebar_current: "guides-dynamic-secrets"
description: |-
  Vault can dynamically generate secrets on-demand for some systems.
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

## Challenge

Data protection is a top priority which means that the database credential
rotation is a critical part of the effort. Each role has a different set of
permissions granted to access the database. When the system is constantly under
attack by hackers, continuous credential rotation becomes necessary and needs to
be automated.


## Solution

The application asks Vault for database credential rather than setting them as
environment variables. The administrator specifies the TTL of the database
credentials to enforce its validity so that they are automatically revoked when
they are no longer used.

![Dynamic Secret Workflow](assets/images/vault-dynamic-secrets.png)

Each app instance can get unique credentials that they don't have to share. By
making those credentials to be short-lived, you reduced the change of the secret
to being compromised. If an app was compromised, the credentials used by the app
can be revoked rather than changing more global set of credentials.

## Prerequisites

To perform the tasks described in this guide, you need to have a Vault
environment.  Refer to the [Getting
Started](/intro/getting-started/install.html) guide to install Vault.

Make sure that your Vault server has been [initialized and
unsealed](/intro/getting-started/deploy.html).


#### PostgreSQL

This guide assumes that you have [PostgreSQL
installed](https://www.postgresql.org/download/), and have a database named
`myapp` created.

**Example on Ubuntu:**

```shell
# Install PostgreSQL
$ sudo apt-get install -y postgresql postgresql-contrib

# Switch to postgres user
$ su - postgres

# Create myapp database
$ psql -U postgres -c 'CREATE DATABASE myapp;'
```



## Steps

In this guide, you are going to configure PostgreSQL secret backend, and create
a read-only database role. The Vault generated PostgreSQL credentials will only
have read permission.

1. [Mount the database secret backend](#step1)
2. [Configure the PostgreSQL backend](#step2)
3. [Create a role](#step3)
4. [Generate PostgreSQL credentials](#step4)

Be sure to follow the [Validation](#validation) to test the outcomes.

### <a name="step1"></a>Step 1: Mount the database secret backend

As most of the secret backends, the [database backend](/docs/secrets/databases/index.html)
must be mounted.

#### CLI command

```shell
vault mount database
```

#### API call using cURL

Before begin, create the following environment variables for your convenience:

- **VAULT_ADDR** is set to your Vault server address
- **VAULT_TOKEN** is set to your Vault token

**Example:**

```plaintext
$ export VAULT_ADDR=http://127.0.0.1:8201

$ export VAULT_TOKEN=0c4d13ba-9f5b-475e-faf2-8f39b28263a5
```

Now, mount the `database` backend using API:

```text
$ curl -X POST -H "X-Vault-Token: $VAULT_TOKEN" -d @postgres.json \
    $VAULT_ADDR/v1/sys/mounts/database

$ cat postgres.json

{
	"type": "database",
	"description": "Database secret backend"
}
```


### <a name="step1"></a>Step 2: Configure the PostgreSQL backend

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
vault write database/config/postgresql plugin_name=postgresql-database-plugin \
  allowed_roles=readonly connection_url=postgresql://root:rootpassword@localhost:5432/myapp
```


#### API call using cURL

**Example:**

```text
$ curl -X POST -H "X-Vault-Token: $VAULT_TOKEN" -d @postgres-config.json \
    $VAULT_ADDR/v1/database/config/postgresql

$ cat postgres-config.json
{
	"plugin_name": "postgresql-database-plugin",
	"allowed_roles": "readonly",
	"connection_url": "postgresql://root:rootpassword@localhost:5432/myapp"
}
```

### <a name="step3"></a>Step 3: Create a role

In Step 2, you configured the PostgreSQL backend by passing **`readonly`** role
as an allowed member. The next step is to define this `readonly` role. A role is
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

```plaintext
vault write database/roles/readonly db_name=postgresql creation_statements=@readonly.sql \
    default_ttl=1h max_ttl=24h
```

The above command creates a role named, `readonly` with default TTL of 1
hour, and max TTL of the credential is set to 24 hours. The `readonly.sql`
statement is passed as the role creation statement.

#### API call using cURL

**Example:**

```text
$ curl -X POST -H "X-Vault-Token: $VAULT_TOKEN" -d @role-payload.json \
    $VAULT_ADDR/v1/database/roles/readonly

$ cat role-payload.JSON

{
	"db_name": "postgres",
	"creation_statements": "CREATE ROLE '{{name}}' WITH LOGIN PASSWORD '{{password}}' VALID UNTIL '{{expiration}}';
  GRANT SELECT ON ALL TABLES IN SCHEMA public TO '{{name}}';",
	"default_ttl": "1h",
	"max_ttl": "24h"
}
```

The `db_name`, `creation_statements`, `default_ttl`, and `max_ttl` are set in the `role-payload.json`.


### <a name="step4"></a>Step 4: Generate PostgreSQL credentials

To generate a new set of PostgreSQL credentials, simply **read** from the `readonly` role endpoint.

**NOTE:** Typically, an administrator performs [Step 1](#step1) through
[Step3](step3).  In order for the client app to get the database credentials
from Vault, the client's policy must include the following:

```text
path "database/creds/readonly" {
  capabilities = ["read"]
}
```

#### CLI command

```shell
vault read database/creds/readonly

Key            	Value
---            	-----
lease_id       	database/creds/readonly/4b5c6e82-df88-0dec-c0cb-f07eee8f0329
lease_duration 	1h0m0s
lease_renewable	true
password       	A1a-4urzp0wu92r5s1q0
username       	v-token-readonly-9x3qrw452wwz4w6421xt-1515625519
```

**NOTE:** Re-run the command and notice that Vault returns a different set of
credentials each time. This means that each app instance can acquire a unique
set of credentials.


#### API call using cURL

```plaintext
curl -X GET -H "X-Vault-Token: $VAULT_TOKEN" $VAULT_ADDR/v1/database/creds/readonly | jq

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
vault read database/creds/readonly

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
vault renew database/creds/readonly/3e8174da-6ca0-143b-aa8c-4c238aa02809

Key            	Value
---            	-----
lease_id       	database/creds/readonly/3e8174da-6ca0-143b-aa8c-4c238aa02809
lease_duration 	1h0m0s
lease_renewable	true
```

(4) Revoke the generated credentials.

```plaintext
vault revoke database/creds/readonly/3e8174da-6ca0-143b-aa8c-4c238aa02809
```

**NOTE:** If you run the command with **`-prefix=true`** flag, it revokes all
secrets under `database/creds/readonly`.

Now, when you check the list of users in PostgreSQL, none of the Vault generated
user name exists.


## Next steps

This guide discussed how to generate credentials on-demand so that the access
credentials no longer need to be written to disk. Next, learn about the
[Tokens and Leases](/guides/lease.html) so that you can control the lifecycle of
those credentials.
