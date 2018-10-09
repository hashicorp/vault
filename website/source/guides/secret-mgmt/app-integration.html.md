---
layout: "guides"
page_title: "Direct Application Integration - Guides"
sidebar_title: "Direct Application Integration"
sidebar_current: "guides-secret-mgmt-app-integration"
description: |-
  This guide demonstrates the use of Consul Template and Envconsul tools. To
  understand the difference between the two tools, you are going to retrieve the
  same information from Vault.
---

# Direct Application Integration

A modern system requires access to a multitude of secrets: database credentials,
API keys for external services, credentials for service-oriented architecture
communication, etc.  Vault steps in to provide a centralized secret management
system.  The next step is to decide how your applications acquire the secrets
from Vault.

This guide introduces ***Consul Template*** and ***Envconsul*** to help you
determine if these tools speed up the integration of your applications once
secrets are securely managed by Vault.

-> **NOTE:** Both [Consul Template](https://github.com/hashicorp/consul-template)
and  [Envconsul](https://github.com/hashicorp/consul-template) are open source
tools.

### Consul Template

Despite its name, Consul Template does **not** require a Consul cluster to
operate. It retrieves secrets from Vault and manages the acquisition and renewal
lifecycle.


### Envconsul

Envconsul launches a subprocess which dynamically populates environment
variables from secrets read from Vault.  Your applications then read those
environment variables.  Despite its name, Envconsul does **not** require a
Consul cluster to operate.  It enables flexibility and portability for
applications across systems.


## Reference Material

- [Consul Template](https://github.com/hashicorp/consul-template)
- [Envconsul](https://github.com/hashicorp/consul-template)
- [Secret as a Service: Dynamic Secrets](/guides/secret-mgmt/dynamic-secrets.html)


## Estimated Time to Complete

10 minutes


## Challenge

If your application code or script contains some secrets (e.g. database
credentials), it makes a good sense to manage the secrets using Vault. However,
it means that your application will need to retrieve the secrets at runtime.
Does that mean the application developers must make some code change?

Is there an easy way to retrieve the secrets from Vault and populate the
application code or script with secrets as needed?


## Solution

Both ***Consul Template*** and ***Envconsul*** provide first-class support for
Vault.  Leveraging these tools can minimize the level of changes introduced to
your applications. Depending on the current application design, you may not need
to make minimal to no code change.


## Prerequisites

To perform the tasks described in this guide, you need:

- A [Vault environment](/intro/getting-started/install.html)
- [Consul Template](https://releases.hashicorp.com/consul-template/)
- [Envconsul](https://releases.hashicorp.com/envconsul/)
- [PostgreSQL](#postgresql)


### PostgreSQL

This guide uses the database secrets engine to demonstrate the use of Consul
Template and Envconsul.  Therefore, you need a
[PostgreSQL](https://www.postgresql.org/download/) server to connect to.

~> Complete the [Secret as a Service: Dynamic
Secrets](/guides/secret-mgmt/dynamic-secrets.html) guide first if you are not
familiar with `database` secrets engine.



### Policy requirements

-> **NOTE:** For the purpose of this guide, you can use **`root`** token to work
with Vault. However, it is recommended that root tokens are only used for just
enough initial setup or in emergencies. As a best practice, use tokens with
appropriate set of policies based on your role in the organization.

To perform all tasks demonstrated in this guide, your policy must include the
following permissions:

```shell
# Enable database secrets engines at "database/" path
path "sys/mounts/database" {
  capabilities = [ "create", "read", "update", "delete", "list" ]
}

# Configure the database secret engine and create roles
path "database/*" {
  capabilities = [ "create", "read", "update", "delete", "list" ]
}

# Write ACL policies
path "sys/policy/*" {
  capabilities = [ "create", "read", "update", "delete", "list" ]
}

# Manage tokens for Consul Template & Envconsul to use
path "auth/token/create" {
  capabilities = [ "create", "read", "update", "delete", "list", "sudo" ]
}
```

If you are not familiar with policies, complete the
[policies](/guides/identity/policies.html) guide.



## Steps

This guide demonstrates the use of Consul Template and Envconsul tools. To
understand the difference between the two tools, you are going to retrieve the
same information from Vault.

1. [Setup Database Secrets Engine](#step1)
1. [Generate Client Token](#step2)
1. [Use Consul Template to Populate DB Credentials](#step3)
1. [Use Envconsul to Retrieve DB Credentials](#step4)


### <a name="step1"></a>Step 1: Setup Database Secrets Engine

In this step, you are going to enable and configure the `database` secrets
engine using `postgresql-database-plugin` where the database connection URL is
`postgresql://root:rootpassword@localhost:5432/myapp`.  

> **NOTE:** Your database connection URL is most likely different from this
example. Be sure to use the correct [connection URL]
(http://localhost:4567/api/secret/databases/postgresql.html#connection_url) to
match your environment.

~> Refer to the [PostgreSQL Database Secrets
Engine](/docs/secrets/databases/postgresql.html) documentation or [Secret as a
Service: Dynamic Secrets](/guides/secret-mgmt/dynamic-secrets.html) guide if you
are not familiar with `database` secrets engine. The detailed description of
working with `database` secrets engine is out of scope for this guide.


#### CLI command

```shell
# First, enable the database secrets engine
$ vault secrets enable database

# Configure the secret engine with appropriate parameter values
$ vault write database/config/postgresql
      plugin_name=postgresql-database-plugin \
      allowed_roles=* \
      connection_url=postgresql://root:rootpassword@localhost:5432/myapp

# Create readonly.sql to define a role permission in SQL
$ tee readonly.sql <<EOF
CREATE ROLE "{{name}}" WITH LOGIN PASSWORD '{{password}}' VALID UNTIL '{{expiration}}';
GRANT SELECT ON ALL TABLES IN SCHEMA public TO "{{name}}";
EOF

# Create a role, "readonly"
$ vault write database/roles/readonly db_name=postgresql creation_statements=@readonly.sql \
    default_ttl=1h max_ttl=24h
```


#### API call using cURL

```shell
# Enable `database` secret engine using `/sys/mounts` endpoint
$ curl --header "X-Vault-Token: ..." \
       --request POST \
       --data '{"type":"database"}' \
       https://127.0.0.1:8200/v1/sys/mounts/database

# Specify the database connection URL according to your environment
$ tee payload.json <<EOF
{
	"plugin_name": "postgresql-database-plugin",
	"allowed_roles": "*",
	"connection_url": "postgresql://root:rootpassword@localhost:5432/myapp"
}
EOF

# Configure the database secrets engine by passing the request payload
$ curl --header "X-Vault-Token: ..." \
       --request POST \
       --data @payload.json \
       http://127.0.0.1:8200/v1/database/config/postgresql

# Create the request payload to create a role
$ tee payload.json <<EOF
{
  "db_name": "postgres",
  "creation_statements": ["CREATE ROLE \"{{name}}\" WITH LOGIN PASSWORD '{{password}}' VALID UNTIL '{{expiration}}';
   GRANT SELECT ON ALL TABLES IN SCHEMA public TO \"{{name}}\";"],
  "default_ttl": "1h",
  "max_ttl": "24h"
}
EOF

# Create a role named, readonly
$ curl --header "X-Vault-Token: ..." \
       --request POST \
       --data @payload.json \
       http://127.0.0.1:8200/v1/database/roles/readonly
```


### <a name="step2"></a>Step 2: Generate Client Token

Consul Template tool itself is a Vault client. Therefore, it must have a valid
token with policies permitting it to retrieve secrets from `database` secret
engine you just configured in [Step 1](#step1).

First, create a policy definition file, `db_creds.hcl`.  This policy allows read
operation on the `database/creds/readonly` path to obtain the dynamically
generated username and password to access the PostgreSQL database. In addition,
the policy allows renewal of the lease if necessary.

**`db_creds.hcl`**:

```hcl
path "database/creds/readonly" {
  capabilities = [ "read" ]
}

path "/sys/leases/renew" {
  capabilities = [ "update" ]
}
```

Now, create a policy named, `db_creds` and generate a token with this policy
attached.


### CLI Command

```shell
# Create a `db_creds` policy
$ vault policy write db_creds db_creds.hcl

# Create a token with db_creds policy:
$ vault token create -policy="db_creds"
Key                  Value
---                  -----
token                89956bf1-6f4d-435d-4cf3-7496e9520a87
token_accessor       319eddff-42a1-eb2b-801e-dd8a0c0b07b4
token_duration       768h
token_renewable      true
token_policies       ["db_creds" "default"]
identity_policies    []
policies             ["db_creds" "default"]
```

**NOTE:** This is the token that Consul Template uses to talk to Vault.
Copy the **`token`** value and proceed to [Step 3](#step3).


#### API call using cURL

```shell
# Create an API request payload
$ tee payload.json <<EOF
{
  "policy": "path \"database/creds/readonly\" {\n capabilities = [ \"read\" ]\n } \n path \"sys/leases/renew\" {\n capabilities = [ \"update\" ] \n}"
}
EOF

# Create db_creds policy
$ curl --header "X-Vault-Token: ..." \
       --request PUT \
       --data @payload.json \
       http://127.0.0.1:8200/v1/sys/policy/db_creds

# Generate a new token with db_creds policy
$ curl --header "X-Vault-Token: ..." \
       --request POST \
       --data '{"policies": ["db_creds"]}' \
       http://127.0.0.1:8200/v1/auth/token/create | jq
{
  ...
   "auth": {
     "client_token": "37413eca-96aa-1d47-d09d-d1cad322c419",
     "accessor": "a05aa3ce-b5d6-9e82-da7f-181d78d475e4",
     "policies": [
       "db_creds",
       "default"
     ],
     "token_policies": [
       "db_creds",
       "default"
     ],
     ...
}
```

**NOTE:** This is the token that Consul Template uses to talk to Vault.
Copy the **`client_token`** value and proceed to [Step 3](#step3).


### <a name="step3"></a>Step 3: Use Consul Template to Populate DB Credentials

Assume that your application requires PostgreSQL database credentials to read
data.  Its configuration file, **`config.yml`** looks like:

```plaintext
username: "<DB_USRENAME>"
password: "<DB_PASSWORD>"
database: "myapp"
```

To have Consul Template to populate the `<DB_USRENAME>` and `<DB_PASSWORD>`, you
need to create a template file with Consul Template [templating
language](https://github.com/hashicorp/consul-template#templating-language).


1. Create a template file by replacing the username and password with
   Consul Template syntax and save it as **`config.yml.tpl`**. The file should
   contain the following:

    ```plaintext
    ---
    {{- with secret "database/creds/readonly" }}
    username: "{{ .Data.username }}"
    password: "{{ .Data.password }}"
    database: "myapp"
    {{- end }}
    ```

    -> **NOTE:** This template reads secrets from `database/creds/readonly`
    path in Vault. Set the `username` parameter value to "`.Data.username`" of
    the secret output. Similarly, set the `password` to "`.Data.password`" value.

1. Execute the `consul-template` command to populate `config.yml` file.

    The Consul Template command is: `consul-template -template="<input_file>:<output_file>"`

    The input file is the **`config.yml.tpl`** and specify the desired output
    file name to be **`config.yml`**:

    ```plaintext
    $ VAULT_TOKEN=<token> consul-template -template="config.yml.tpl:config.yml" -once
    ```

    While `<token>` is the token you copied at [Step 2](#step2).


1. Open the generated **`config.yml`** file to verify its content. It should
   look similar to:

    ```plaintext
    $ cat config.yml
    ---
    username: "v-token-readonly-tu17xrtz345uz643980r-1527630039"
    password: "A1a-7s0z9y223x2rp6v9"
    database: "myapp"
    ```

    The `username` and `password` were retrieved from Vault and populated in the
     `config.yml` file.


> **Summary**: You need to create a templated version of your application
scripts to leverage Consul Template.  However, it requires minimum effort to do
so in comparison to writing an application which invokes Vault API to accomplish
the same.



### <a name="step4"></a>Step 4: Use Envconsul to Retrieve DB Credentials

Create a file named, **`app.sh`** containing the following:

```hcl
#!/usr/bin/env bash

cat <<EOT
My connection info is:

username: "${DATABASE_CREDS_READONLY_USERNAME}"
password: "${DATABASE_CREDS_READONLY_PASSWORD}"
database: "my-app"
EOT
```

The main difference here is that the `app.sh` is reading ***environment
variables*** to set `username` and `password` values; therefore, no templating
is involved.

-> Notice that the environment variable name is derived from the secret _path_
with key name.

Run the Envconsul tool using the Vault token you generated at [Step 2](#step2).

```plaintext
$ VAULT_TOKEN=<token> envconsul -upcase -secret database/creds/readonly ./app.sh

My connection info is:

username: "v-token-readonly-ww1tq33s7z5uprpxxy68-1527631219"
password: "A1a-u54wut0v605qwz95"
database: "my-app"
```

The output should display the `username` and `password` populated.  

The `-upcase` flag tells Envconsul to convert all environment variable keys to
uppercase.  Otherwise, the default uses lowercase (e.g. `database_creds_readonly_username`).


> **Summary:** If your application is designed to read secrets from environment
variables, Envconsul requires minimal to no code change to integrate with Vault.



## Next steps

If the integration option is to directly invoke Vault API within your
application, refer to the [_AppRole Pull
Authentication_](/guides/identity/authentication.html) guide to learn about the
AppRole auth method which is designed for applications.
