# Vault Plugin Database Redis ElastiCache

This is a standalone [Database Plugin](https://www.vaultproject.io/docs/secrets/databases) for use with [Hashicorp
Vault](https://www.github.com/hashicorp/vault).

This plugin supports exclusively AWS ElastiCache for Redis. [Redis Enterprise](https://github.com/RedisLabs/vault-plugin-database-redis-enterprise) 
and [Redis Open Source](https://github.com/fhitchen/vault-plugin-database-redis) use different plugins.

Please note: We take Vault's security and our users' trust very seriously. If
you believe you have found a security issue in Vault, please responsibly
disclose by contacting us at [security@hashicorp.com](mailto:security@hashicorp.com).


## Quick Links

- [Vault Website](https://www.vaultproject.io)
- [Plugin System](https://www.vaultproject.io/docs/plugins)


## Compatibility

The plugin is automatically available in Vault starting from [v1.12.0](https://developer.hashicorp.com/vault/docs/release-notes/1.12.0).
If you are using a previous version of Vault, you can [build & register the binary yourself](https://developer.hashicorp.com/vault/docs/plugins/plugin-management)
similar to how a custom external plugin operates.


## Getting Started

This is a [Vault plugin](https://www.vaultproject.io/docs/plugins)
and is meant to work with Vault. This guide assumes you have already installed
Vault and have a basic understanding of how Vault works.

Otherwise, first read this guide on how to [get started with
Vault](https://www.vaultproject.io/intro/getting-started/install.html).


## Development

If you wish to work on this plugin, you'll first need
[Go](https://www.golang.org) installed on your machine (version 1.17+ recommended)

Make sure Go is properly installed, including setting up a [GOPATH](https://golang.org/doc/code.html#GOPATH).

To run the tests locally you will need to have write permissions to an [ElastiCache for Redis](https://aws.amazon.com/elasticache/redis/) instance. 
A small Terraform project is included to provision one for you if needed. More details in the [Environment Set Up](#environment-set-up) section.

## Building

If you're developing for the first time, run `make bootstrap` to install the
necessary tools. Bootstrap will also update repository name references if that
has not been performed ever before.

```sh
$ make bootstrap
```

To compile a development version of this plugin, run `make` or `make dev`.
This will put the plugin binary in the `bin` and `$GOPATH/bin` folders. `dev`
mode will only generate the binary for your platform and is faster:

```sh
$ make dev
```

## Tests

### Environment Set Up

To test the plugin, you need access to an Elasticache for Redis Cluster. 
A Terraform project is included for convenience to initialize a new cluster if needed.
If not already available, you can install Terraform by using [this documentation](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-files.html).

The setup script tries to find and use available AWS credentials from the environment. You can configure AWS credentials using [this documentation](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-files.html).
Or if you prefer you can edit the provider defined ./bootstrap/terraform/elasticache.tf with your desired set of credentials.

Note that resources created via the Terraform project cost a small amount of money per hour.

To set up the test cluster:

```sh
$ make setup-env
...
Apply complete! Resources: 4 added, 0 changed, 0 destroyed.
```

Set the `create_aws_user` variable to `false` to skip creating an IAM user for
plugin management:

```sh
$ make setup-env TF_VAR_create_aws_user=false
...
Apply complete! Resources: 4 added, 0 changed, 0 destroyed.
```

### Environment Teardown

The test cluster created via the setup-env command can be destroyed using the teardown-env command.

```sh
$ make teardown-env
...
Destroy complete! Resources: 4 destroyed.
```

### Testing Manually

Put the plugin binary into a location of your choice. This directory
will be specified as the [`plugin_directory`](https://www.vaultproject.io/docs/configuration#plugin_directory)
in the Vault config used to start the server.

```hcl
# config.hcl
plugin_directory = "path/to/plugin/directory"
...
```

Start a Vault server with this config file:

```sh
$ vault server -dev -config=path/to/config.hcl ...
...
```

Once the server is started, register the plugin in the Vault server's [plugin catalog](https://www.vaultproject.io/docs/plugins/plugin-architecture#plugin-catalog):

```sh
$ SHA256=$(openssl dgst -sha256 $GOPATH/vault-plugin-database-redis-elasticache | cut -d ' ' -f2)
$ vault plugin register -sha256=$SHA256 database vault-plugin-database-redis-elasticache
...
Success! Data written to: sys/plugins/catalog/database/vault-plugin-database-redis-elasticache
```

Enable the database engine to use this plugin:

```sh
$ vault secrets enable database
...

Success! Enabled the database secrets engine at: database/
```

Once the database engine is enabled you can configure an ElastiCache instance:

```sh
$ vault write database/config/redis-mydb \
        plugin_name="vault-plugin-database-redis-elasticache" \
        username=$USERNAME \
        password=$PASSWORD \
        url=$URL \
        region=$REGION
...

Success! Data written to: database/config/redis-mydb
```

Configure a static role:

```sh
$ vault write database/static-roles/redis-myrole \
        db_name="redis-mydb" \
        username="my-elasticache-username" \
        rotation_period=5m
...

Success! Data written to: database/roles/redis-myrole
```

Retrieve your first set of static credentials:

```sh
$ vault read database/static-creds/redis-myrole
Key                    Value
---                    -----
last_vault_rotation    2022-09-06T12:15:33.958413491-04:00
password               PASSWORD
rotation_period        5m
ttl                    4m55s
username               my-elasticache-username
```


### Automated Tests

To run the tests, invoke `make test`:

```sh
$ make test
```

You can also specify a `TESTARGS` variable to filter tests like so:

```sh
$ make test TESTARGS='-run=TestConfig'
```

### Acceptance Tests

The majority of tests must communicate with an existing ElastiCache instance. See the [Environment Set Up](#environment-set-up) section for instructions on how to prepare a test cluster.

Some environment variables are required to run tests expecting to communicate with an ElastiCache cluster.
The username and password should be valid IAM access key and secret key with read and write access to the ElastiCache cluster used for testing. They may also be inferred from the usual `AWS_*` environment variables.
The URL should be the complete configuration endpoint including the port, for example: `vault-plugin-elasticache-test.id.xxx.use1.cache.amazonaws.com:6379`.

```sh
$ export TEST_ELASTICACHE_ACCESS_KEY_ID="AWS ACCESS KEY ID"
$ export TEST_ELASTICACHE_SECRET_ACCESS_KEY="AWS SECRET ACCESS KEY"
$ export TEST_ELASTICACHE_URL="vault-plugin-elasticache-test.id.xxx.use1.cache.amazonaws.com:6379"
$ export TEST_ELASTICACHE_REGION="us-east-1"
$ export TEST_ELASTICACHE_USER="vault-test"

$ make testacc
```

You can also specify a `TESTARGS` variable to filter tests like so:

```sh
$ make testacc TESTARGS='-run=TestConfig'
```
