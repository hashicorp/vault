# Vault Plugin: Azure Secrets Backend

This is a standalone backend plugin for use with [Hashicorp Vault](https://www.github.com/hashicorp/vault).
This plugin generates revocable, time-limited Service Principals for Microsoft Azure.

**Please note**: We take Vault's security and our users' trust very seriously. If you believe you have found a security issue in Vault, _please responsibly disclose_ by contacting us at [security@hashicorp.com](mailto:security@hashicorp.com).

## Quick Links
- [Vault Website](https://developer.hashicorp.com/vault/docs)
- [Azure Secrets Docs](https://developer.hashicorp.com/vault/docs/secrets/azure)
- [Vault Github Project](https://www.github.com/hashicorp/vault)

## Getting Started

This is a [Vault plugin](https://developer.hashicorp.com/vault/docs/plugins)
and is meant to work with Vault. This guide assumes you have already installed Vault
and have a basic understanding of how Vault works.

Otherwise, first read this guide on how to [get started with Vault](https://developer.hashicorp.com/vault/tutorials/getting-started/getting-started-install).

To learn specifically about how plugins work, see documentation on [Vault plugins](https://developer.hashicorp.com/vault/docs/plugins).

## Usage

Please see [documentation for the plugin](https://developer.hashicorp.com/vault/docs/secrets/azure)
on the Vault website.

This plugin is currently built into Vault and by default is accessed
at `azure`. To enable this in a running Vault server:

```sh
$ vault secrets enable azure
Success! Enabled the azure secrets engine at: azure/
```

## Developing

If you wish to work on this plugin, you'll first need
[Go](https://www.golang.org) installed on your machine
(version 1.17+ is *required*).

For local dev first make sure Go is properly installed, including
setting up a [GOPATH](https://golang.org/doc/code.html#GOPATH).
Next, clone this repository into
`$GOPATH/src/github.com/hashicorp/vault-plugin-secrets-azure`.
You can then download any required build tools by bootstrapping your
environment:

```sh
$ make bootstrap
```

To compile a development version of this plugin, run `make` or `make dev`.
This will put the plugin binary in the `bin` and `$GOPATH/bin` folders. `dev`
mode will only generate the binary for your platform and is faster:

```sh
$ make
$ make dev
```

Put the plugin binary into a location of your choice. This directory
will be specified as the [`plugin_directory`](https://developer.hashicorp.com/vault/docs/configuration#plugin_directory)
in the Vault config used to start the server.

```json
plugin_directory = "path/to/plugin/directory"
```

Start a Vault server with this config file:
```sh
$ vault server -config=path/to/config.json ...
...
```

Once the server is started, register the plugin in the Vault server's [plugin catalog](https://developer.hashicorp.com/vault/docs/plugins/plugin-architecture#plugin-catalog):

```sh

$ vault plugin register \
        -sha256=<SHA256 Hex value of the plugin binary> \
        -command="vault-plugin-secrets-azure" \
        secret \
        azure
...
Success! Data written to: sys/plugins/catalog/azure
```

Note you should generate a new sha256 checksum if you have made changes
to the plugin. Example using openssl:

```sh
openssl dgst -sha256 $GOPATH/vault-plugin-secrets-azure
...
SHA256(.../go/bin/vault-plugin-secrets-azure)= 896c13c0f5305daed381952a128322e02bc28a57d0c862a78cbc2ea66e8c6fa1
```

Enable the plugin backend using the secrets enable plugin command:

```sh
$ vault secrets enable -plugin-name='azure' plugin
...

Successfully enabled 'plugin' at 'azure'!
```

### Azure Environment Setup

A Terraform [configuration](bootstrap/terraform) is included in this repository that
automates provisioning of Azure resources necessary to configure the secrets engine.
By default, the resources are created in `westus2`. See [variables.tf](bootstrap/terraform/variables.tf) 
for the available variables.

Before applying the Terraform configuration, you'll need to:

1. [Authenticate](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs#authenticating-to-azure)
   the Terraform provider to Azure

The Terraform configuration will create:

* An app registration with necessary API permissions
* A service principal with necessary role assignments

To provision the Azure resources, run the following:

```sh
$ make setup-env   
```

The `local_environment_setup.sh` file will be created in the `bootstrap/terraform`
directory as a result of running `make setup-env`. This file contains environment
variables needed to configure the secrets engine. The values can also be accessed
via `terraform output`.

Once you're finished with plugin development, you can run the following to
destroy the Azure resources:

```sh
$ make teardown-env   
```

### Configure Plugin

A [scripted configuration](bootstrap/configure.sh) of the plugin is provided in
this repository. You can use the script or manually configure the secrets engine
using documentation.

To apply the scripted configuration, first source the environment variables generated by
the Azure environment setup:

```sh
$ source ./bootstrap/terraform/local_environment_setup.sh
```

Next, run the `make configure` target to register, enable, and configure the plugin with
your local Vault instance. You can specify the plugin name, plugin directory, and mount
path. Default values from the Makefile will be used if arguments aren't provided.

```sh
$ PLUGIN_NAME=vault-plugin-secrets-azure \
  PLUGIN_DIR=$GOPATH/vault-plugins \
  PLUGIN_PATH=local-secrets-azure \
  make configure
```

#### Tests

If you are developing this plugin and want to verify it is still
functioning (and you haven't broken anything else), we recommend
running the tests.

To run the tests, invoke `make test`:

```sh
$ make test
```

You can also specify a `TESTARGS` variable to filter tests like so:

```sh
$ make test TESTARGS='--run=TestConfig'
```

#### Acceptance Tests

This repository contains acceptance tests that interact with real Azure resources. There
are acceptance tests written in Go and [bats](https://bats-core.readthedocs.io/en/stable).

##### Go

To run the Go acceptance tests, run the following:

```sh
$ make testacc 
```

##### Bats

Acceptance tests requires Azure access, and the following to be installed:
- [Docker](https://docs.docker.com/get-docker/)
- [Azure CLI](https://docs.microsoft.com/en-us/cli/azure/install-azure-cli)
- [Terraform](https://learn.hashicorp.com/tutorials/terraform/install-cli)
- [bats](https://bats-core.readthedocs.io/en/stable)

_You will need to be properly logged in to Azure with your subscription set. See
['Azure Provider: Authenticating using the Azure CLI'](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/guides/azure_cli)_
for more information.

```sh
$ make test-acceptance AZURE_TENANT_ID=<your_tenant_id>
```

Setting `WITH_DEV_PLUGIN=` will use the provided builtin plugin. The default behavior is to build and register
the plugin from the working directory.

```sh
$ make test-acceptance AZURE_TENANT_ID=<your_tenant_id>
```

Running tests against Vault Enterprise requires a valid license, and specifying an enterprise docker image:

```sh
$ make test-acceptance AZURE_TENANT_ID=<your_tenant_id> \
  VAULT_LICENSE=........ \
  VAULT_IMAGE=hashicorp/vault-enterprise:latest
```

The `test-acceptance` make target also accepts the following environment based directives:

* `TESTS_FILTER`: a regex of Bats tests to run, useful when you only want to run a subset of the tests.
