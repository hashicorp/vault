# vault-plugin-database-couchbase

A [Vault](https://www.vaultproject.io) plugin for Couchbase

This project uses the database plugin interface introduced in Vault version 0.7.1.

The plugin supports the generation of static and dynamic user roles and root credential rotation.

## Build

To build this package for any platform you will need to clone this repository and cd into the repo directory and `go build -o couchbase-database-plugin ./cmd/couchbase-database-plugin/`. To test `go test` will execute a set of basic tests against against the docker.io/couchbase/server-sandbox:6.5.0 couchbase database image. To test against different sandbox images, for example 5.5.1, set the `COUCHBASE_VERSION=5.5.1` environment variable. If you want to run the tests against a local couchbase installation or an already running couchbase container, set the environment variable `COUCHBASE_HOST` before executing. **Note** you will need to align the Administrator username, password and bucket_name with the pre-set values in the `couchbase_test.go` file. Set VAULT_ACC to execute all of the tests. A subset of tests can be run using the command `go test -run TestDriver/Init` for example.

## Installation

The Vault plugin system is documented on the [Vault documentation site](https://www.vaultproject.io/docs/internals/plugins.html).

You will need to define a plugin directory using the `plugin_directory` configuration directive, then place the
`vault-plugin-database-couchbase` executable generated above, into the directory.

**Please note:** Versions v0.2.0 onwards of this plugin are incompatible with Vault versions before 1.6.0 due to an update of the database plugin interface.

Sample commands for registering and starting to use the plugin:

```bash
$ SHA256=$(shasum -a 256 plugins/couchbase-database-plugin | cut -d' ' -f1)

$ vault secrets enable database

$ vault write sys/plugins/catalog/database/couchbase-database-plugin sha256=$SHA256 \
        command=couchbase-database-plugin
```

At this stage you are now ready to initialize the plugin to connect to couchbase cluster using unencrypted or encrypted communications.

Prior to initializing the plugin, ensure that you have created an administration account. Vault will use the user specified here to create/update/revoke database credentials. That user must have the appropriate permissions to perform actions upon other database users.

### Unencrypted plugin initialization

```bash
$ vault write database/config/insecure-couchbase plugin_name="couchbase-database-plugin" \
        hosts="localhost" username="Administrator" password="password" \
        bucket_name="travel-sample" \ # only needed for pre-6.5.0 clusters
        allowed_roles="insecure-couchbase-admin-role,insecure-couchbase-*-bucket-role,static-account"

# You should consider rotating the admin password. Note that if you do, the new password will never be made available
# through Vault, so you should create a vault-specific database admin user for this.
$ vault write -force database/rotate-root/insecure-couchbase

 ```

Note: If you want to connect the plugin to a couchbase cluster prior to version 6.5.0 you will also have to supply an existing bucket (bucket_name="travel-sample") or the command will fail with the error message **"error verifying connection: error in Connection waiting for cluster: unambiguous timeout"**.

### Encrypted plugin initialization

The example here uses the self signed CA certificate that comes with the out of the box couchbase cluster installation and is not suitable for real production use where commercial grade certificates should be obtained.

```bash
$ BASE64PEM=$(curl -X GET http://Administrator:Admin123@127.0.0.1:8091/pools/default/certificate|base64 -w0)

$ vault write database/config/secure-couchbase plugin_name="couchbase-database-plugin" \
      hosts="couchbases://localhost" username="Administrator" password="password" \
      tls=true base64pem=${BASE64PEM} \
      bucket_name="travel-sample" \ # only needed for pre-6.5.0 clusters
      allowed_roles="secure-couchbase-admin-role,secure-couchbase-*-bucket-role,static-account"

# You should consider rotating the admin password. Note that if you do, the new password will never be made available
# through Vault, so you should create a vault-specific database admin user for this.
$ vault write -force database/rotate-root/secure-couchbase
```

### Dynamic Role Creation

When you create roles, you need to provide a JSON string containing the Couchbase RBAC roles which are documented [here](https://docs.couchbase.com/server/6.5/learn/security/roles.html). From Couchbase 6.5 groups are supported and the creation statement can contain just roles or just groups or a mixture of the two. **Note** to use a group, it must have been created in the database previously.

```bash
# if a creation_statement is not provided the user account will default to read only admin, '{"roles":[{"role":"ro_admin"}]}'
$ vault write database/roles/insecure-couchbase-admin-role db_name=insecure-couchbase \
        default_ttl="5m" max_ttl="1h" creation_statements='{"roles":[{"role":"admin"}],"groups":["Supervisor"]}'

$ vault write database/roles/insecure-couchbase-travel-sample-bucket-role db_name=insecure-couchbase \
        default_ttl="5m" max_ttl="1h" creation_statements='{"roles":[{"role":"bucket_full_access","bucket_name":"travel-sample"}]}'
Success! Data written to: database/roles/insecure-couchbase-travel-sample-bucket-role
```

If you create a role that uses groups on a pre 6.5 couchbase server it will be successful, but when you try to generate credentials
you will receive the error **rpc error: code = Unknown desc = {"errors":{"groups":"Unsupported key"}} ...**

To retrieve the credentials for the dynamic accounts

```bash

$ vault read database/creds/insecure-couchbase-admin-role
Key                Value
---                -----
lease_id           database/creds/insecure-couchbase-admin-role/KJ7CTmpFni6U6BCDJ14HcmDm
lease_duration     5m
lease_renewable    true
password           A1a-yCSH5rAh8QAkCzwu
username           v-token-insecure-couchbase-admin-role-yA2hgb0tfewf

$ vault read database/creds/insecure-couchbase-travel-sample-bucket-role
Key                Value
---                -----
lease_id           database/creds/insecure-couchbase-travel-sample-bucket-role/OzHdfkIZdeY9p8kjdWur512j
lease_duration     5m
lease_renewable    true
password           A1a-0yTIuO4q0dCvphz1
username           v-token-insecure-couchbase-travel-sample-bucket-role-iN5

```

### Static Role Creation

In order to use static roles, the user must already exist in the Couchbase security settings. The example below assumes that there is an existing user with the name "vault-edu". If the user does not exist you will receive the following error.

```bash
* 1 error occurred:
        * error setting credentials: rpc error: code = Unknown desc = user not found | {"unique_id":"74f229fd-b3b3-4036-9673-312adae094bb","endpoint":"http://localhost:8091"}
```

```bash
$ vault write database/static-roles/static-account db_name=insecure-couchbase \
        username="vault-edu" rotation_period="5m"
Success! Data written to: database/static-roles/static-account
````

To retrieve the credentials for the vault-edu user

```bash
$ vault read database/static-creds/static-account
Key                    Value
---                    -----
last_vault_rotation    2020-06-15T14:32:16.682130141-05:00
password               A1a-09ApRvglZY1Usdjp
rotation_period        5m
ttl                    30s
username               vault-edu
```

## Developing

You can run `make dev-vault` in the root of the repo to start up a development vault server and automatically register a local build of the plugin. You will need to have a built `vault` binary available in your `$PATH` to do so.

### Acceptance tests

Run `make testacc`.