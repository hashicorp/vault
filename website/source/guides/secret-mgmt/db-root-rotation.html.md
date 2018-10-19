---
layout: "guides"
page_title: "DB Root Credential Rotation - Guides"
sidebar_title: "DB Root Credential Rotation"
sidebar_current: "guides-secret-mgmt-db-root-rotation"
description: |-
  Vault enables the combined database secret engines to automate the rotation of
  root credentials.
---

# Database Root Credential Rotation

## Database Secrets Engine

Vault's [database secrets engine](/docs/secrets/databases/index.html) provides a
centralized workflow for managing credentials for various database systems. By
leveraging this, every service instance gets a unique set of database
credentials instead of sharing one.  Having those credentials tied directly to
each service instance and live only for the life of the service, any abnormal
access pattern can be mapped to a specific service instance and its credential
can be revoked immediately.

This reduces the manual tasks performed by the database administrator and make
the access to the database to be more efficient and secure.

The [Secret as a Service: Dynamic Secrets](/guides/secret-mgmt/dynamic-secrets.html)
guide demonstrates the primary workflow.

## Reference Material

- [Secret as a Service: Dynamic Secrets](/guides/secret-mgmt/dynamic-secrets.html)
- [Database Secret Engine (API)](/api/secret/databases/index.html#rotate-root-credentials)
- [PostgreSQL Database Plugin HTTP API](/api/secret/databases/postgresql.html)

## Estimated Time to Complete

10 minutes


## Challenge

Because Vault is managing the database credentials on behalf of the database
administrator, it must also be given a set of highly privileged credentials
which can grant and revoke access to the database system.  Therefore, it is very
common to give Vault the **root** credentials.

However, these credentials are often long-lived and never change once configured
on Vault. This may violate the _Governance, Risk and Compliance_ (GRC)
surrounding that data stored in the database.


## Solution

Use the Vault's **`/database/rotate-root/:name`** API endpoint to rotate the
root credentials stored for the database connection.

![DB Root Credentials](/img/vault-db-root-rotation.png)

~> **Best Practice:** Use this feature to rotate the root credentials
immediately after the initial configuration of each database.

## Prerequisites

To perform the tasks described in this guide, you need to have a Vault
environment.  Refer to the [Getting
Started](/intro/getting-started/install.html) guide to install Vault. Make sure
that your Vault server has been [initialized and
unsealed](/intro/getting-started/deploy.html).

### PostgreSQL

This guide requires that you have a PostgreSQL server to connect to. If you
don't have one, install [PostgreSQL](https://www.postgresql.org/download/).

- Refer to the [PostgreSQL documentation](https://www.postgresql.org/docs/online-resources/) for details
- [PostgreSQL Wiki](https://wiki.postgresql.org/wiki/First_steps) gives you a
summary of basic commands to get started.

### Policy requirements

-> **NOTE:** For the purpose of this guide, you can use **`root`** token to work
with Vault. However, it is recommended that root tokens are only used for just
enough initial setup or in emergencies. As a best practice, use tokens with
appropriate set of policies based on your role in the organization.

To perform all tasks demonstrated in this guide, your policy must include the
following permissions:

```shell
# Mount secret engines
path "sys/mounts/*" {
  capabilities = [ "create", "read", "update", "delete", "list" ]
}

# Configure the database secret engine and create roles
path "database/*" {
  capabilities = [ "create", "read", "update", "delete", "list" ]
}
```

If you are not familiar with policies, complete the
[policies](/guides/identity/policies.html) guide.


## Steps

Using Vault, you can easily rotate the root credentials for your database
through the **`database/rotate-root/:name`** endpoint.

This guide demonstrates the overall workflow to manage the database credentials
including the root.

You are going to perform the following:

1. [Enable the database secret engine](#step1)
1. [Configure PostgreSQL secret engine](#step2)
1. [Verify the configuration (***Optional***)](#step3)
1. [Rotate the root credentials](#step4)


### <a name="step1"></a>Step 1: Enable the database secret engine

#### CLI command

Enable a database secret engine:

```plaintext
$ vault secrets enable database
```

**NOTE:** This example enables the database secret engine at the **`/database`**
path in Vault.  

#### API call using cURL

To enable a database secret engine, use the **`/sys/mounts`** endpoint.

```plaintext
$ curl --header "X-Vault-Token: ..." \
       --request POST \
       --data '{"type":"database"}' \
       https://127.0.0.1:8200/v1/sys/mounts/database
```

**NOTE:** This example mounts database secret engine at **`/database`**,
and passes the secret engine type ("`database`") in the request payload.

### <a name="step2"></a>Step 2: Configure PostgreSQL secret engine

In the [Secret as a Service: Dynamic
Secrets](/guides/secret-mgmt/dynamic-secrets.html#step2) guide, the PostgreSQL
plugin was configured with its root credentials embedded in the `connection_url`
(`root` and `rootpassword`) as below:

```plaintext
$ vault write database/config/postgresql \
      plugin_name=postgresql-database-plugin \
      allowed_roles="*" \
      connection_url=postgresql://root:rootpassword@postgres.host.address:5432/postgres
```

The username and password can be templated using the format,
**`{{<field-name>}}`**.

~> In order to leverage the database root credential rotation feature, you must
use the templated credentials: **`{{username}}`** and **`{{password}}`**.


#### CLI command

```plaintext
$ vault write database/config/postgresql \
     plugin_name=postgresql-database-plugin \
     connection_url="postgresql://{{username}}:{{password}}@postgres.host.address:5432/postgres" \
     allowed_roles="*" \
     username="root" \
     password="rootpassword"
```

Notice that the `connection_url` value contains the templated credentials, and
`username` and `password` parameters are also passed to initiate the connection.


Create a role, `readonly`:

```shell
# Create readonly.sql defining the role privilege  
$ tee readonly.sql <<EOF
CREATE ROLE "{{name}}" WITH LOGIN PASSWORD '{{password}}' VALID UNTIL '{{expiration}}';
GRANT SELECT ON ALL TABLES IN SCHEMA public TO "{{name}}";
EOF

# Create a role named 'readonly' with TTL of 1 hr
$ vault write database/roles/readonly db_name=postgresql \
         creation_statements=@readonly.sql \
         default_ttl=1h max_ttl=24h
```

#### API call using cURL

```plaintext
$ tee payload.json <<EOF
{
    "plugin_name": "postgresql-database-plugin",
    "connection_url": "postgresql://{{username}}:{{password}}@postgres.host.address:5432/postgres",
    "allowed_roles": "readonly",
    "username": "root",
    "password": "rootpassword"
}
EOF

$ curl --header "X-Vault-Token: ..."
       --request POST \
       --data @payload.json \
       http://127.0.0.1:8200/v1/database/config/postgresql
```

Notice that the `connection_url` value contains the templated credentials, and
`username` and `password` parameters are also passed to initiate the connection.


Create a role, `readonly`:

```shell
# Create the request payload
$ tee payload.json <<EOF
{
	"db_name": "postgres",
	"creation_statements": "CREATE ROLE '{{name}}' WITH LOGIN PASSWORD '{{password}}' VALID UNTIL '{{expiration}}';
  GRANT SELECT ON ALL TABLES IN SCHEMA public TO '{{name}}';",
	"default_ttl": "1h",
	"max_ttl": "24h"
}
EOF

$ curl --header "X-Vault-Token: ..."
       --request POST \
       --data @payload.json \
       http://127.0.0.1:8200/v1/database/roles/readonly
```


### <a name="step3"></a>Step 3: Verify the configuration (***Optional***)

Before rotate the root credentials, make sure that the secret engine was
configured correctly.

#### CLI command

```shell
# Get a new set of credentials
$ vault read database/creds/readonly
Key                Value
---                -----
lease_id           database/creds/readonly/999c43f0-f79e-ba90-24a8-4de5af33a2e9
lease_duration     1h
lease_renewable    true
password           A1a-u7wxtrpx09xp40yq
username           v-root-readonly-x6q809467q98yp4yx4z4-1525378026e


# Make sure that you can connect to the database using the Vault generated credentials
$ psql -h postgres.host.address -p 5432 \
       -U v-root-readonly-x6q809467q98yp4yx4z4-1525378026e postgres
Password for user v-root-readonly-x6q809467q98yp4yx4z4-1525378026:

postgres=> \du
Role name                                       |                         Attributes                         | Member of
------------------------------------------------+------------------------------------------------------------+----------
postgres                                        | Superuser, Create role, Create DB, Replication, Bypass RLS | {}
v-root-readonly-x6q809467q98yp4yx4z4-1525378026 | Password valid until 2018-05-03 21:07:11+00                | {}

postgres=> \q
```


#### API call using cURL

```shell
# Get a new set of credentials
$ curl --header "X-Vault-Token: 1c97b03a-6098-31cf-9d8b-b404e52dcb4a" \
       http://127.0.0.1:8200/v1/database/creds/readonly | jq
{
   "request_id": "527970fd-f5e8-4de5-d4ed-1b7970eaef0b",
   "lease_id": "database/creds/readonly/ac79265e-668c-242f-4f67-1dae33da094c",
   "renewable": true,
   "lease_duration": 3600,
   "data": {
     "password": "A1a-0tr8u15y0us2u08v",
     "username": "v-root-readonly-x7v65y1xuprzxv9vpt80-1525378873"
   },
   "wrap_info": null,
   "warnings": null,
   "auth": null
}

# Make sure that you can connect to the database using the Vault generated credentials
$ psql -h postgres.host.address -p 5432 \
       -U v-root-readonly-x6q809467q98yp4yx4z4-1525378026e postgres
Password for user v-root-readonly-x6q809467q98yp4yx4z4-1525378026:

postgres=> \du
Role name                                       |                         Attributes                         | Member of
------------------------------------------------+------------------------------------------------------------+----------
postgres                                        | Superuser, Create role, Create DB, Replication, Bypass RLS | {}
v-root-readonly-x6q809467q98yp4yx4z4-1525378026 | Password valid until 2018-05-03 21:07:11+00                | {}
v-root-readonly-x7v65y1xuprzxv9vpt80-1525378873 | Password valid until 2018-05-03 21:21:18+00                | {}

postgres=> \q
```
<br>

This confirms that the Vault successfully connected to your PostgreSQL server
and created a new user based on the privilege defined by `readonly.sql`.

> The user credentials generated by the Vault has a limited TTL based on your
configuration (`default_ttl`). In addition, you can revoke them if necessary.


### <a name="step4"></a>Step 4: Rotate the root credentials

Vault provides an API endpoint to easily rotate the root database credentials.

#### CLI command

```plaintext
$ vault write -force database/rotate-root/postgresql
```

#### API call using cURL

```plaintext
$ curl --header "X-Vault-Token: ..." \
       --request POST \
       http://127.0.0.1:8200/v1/database/rotate-root/postgresql
```

This is all you need to do.

---
<br>

Repeat [Step 3](#step3) to verify that Vault continues to generate database
credentials after the root credential rotation.

To verify that the root credential was rotated:

```plaintext
$ psql -h postgres.host.address -p 5432 -U root postgres
Password for user root:
```

Entering the initial password (e.g. `rootpassword`) will ***not*** work since
the password was rotated by the Vault.

You can invoke the **`database/rotate-root/:name`** endpoint periodically to
secure the root credential.


~> **NOTE:** Once the root credential was rotated, only the Vault knows the new
root password. This is the same for all root database credentials given to Vault.
Therefore, you should create a separate superuser dedicated to the Vault usage
which is not used for other purposes.  


## Next steps

In this guide, you learned how to rotate the root database credentials.

Read the [AppRole Pull Authentication](/guides/identity/authentication.html)
guide to learn about generating a client token for your app so that it can
request database credentials from Vault.
