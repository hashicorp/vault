# Vault Plugin: Azure Auth Backend

This is a standalone backend plugin for use with [Hashicorp Vault](https://www.github.com/hashicorp/vault).
This plugin allows for Azure Managed Service Identities to authenticate with Vault.

**Please note**: We take Vault's security and our users' trust very seriously. If you believe you have found a security issue in Vault, _please responsibly disclose_ by contacting us at [security@hashicorp.com](mailto:security@hashicorp.com).

## Quick Links

- [Vault Website](https://www.vaultproject.io)
- [Vault Project Github](https://www.github.com/hashicorp/vault)
- [Azure Auth Docs](https://developer.hashicorp.com/vault/docs/auth/azure)
- [Azure Auth API Docs](https://developer.hashicorp.com/vault/api-docs/auth/azure)

## Getting Started

This is a [Vault plugin](https://developer.hashicorp.com/vault/docs/plugins)
and is meant to work with Vault. This guide assumes you have already installed Vault
and have a basic understanding of how Vault works.

Otherwise, first read this guide on how to [get started with
Vault](https://developer.hashicorp.com/vault/tutorials/getting-started/getting-started-install).

To learn specifically about how plugins work, see documentation on [Vault plugins](https://developer.hashicorp.com/vault/docs/plugins).

## Security Model

The current authentication model requires providing Vault with a token generated using Azure's Managed Service Identity, which can be used to make authenticated calls to Azure. This token should not typically be shared, but in order for Azure to be treated as a trusted third party, Vault must validate something that Azure has cryptographically signed and that conveys the identity of the token holder.

## Usage

Please see [documentation for the plugin](https://developer.hashicorp.com/vault/docs/auth/azure)
on the Vault website.

This plugin is currently built into Vault and by default is accessed
at `auth/azure`. To enable this in a running Vault server:

```sh
$ vault auth enable azure
Successfully enabled 'azure' at 'azure'!
```

To see all the supported paths, see the [Azure auth backend docs](https://developer.hashicorp.com/vault/docs/auth/azure).

## Developing

If you wish to work on this plugin, you'll first need
[Go](https://www.golang.org) installed on your machine.

### Build Plugin

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

Put the plugin binary into a location of your choice. This directory
will be specified as the [`plugin_directory`](https://developer.hashicorp.com/vault/docs/configuration#plugin_directory)
in the Vault config used to start the server. It may also be specified
via [`-dev-plugin-dir`](https://developer.hashicorp.com/vault/docs/commands/server#dev-plugin-dir)
if running Vault in dev mode.

```hcl
# config.hcl
plugin_directory = "path/to/plugin/directory"
...
```

### Register Plugin

Start a Vault server with this config file:

```sh
$ vault server -dev -config=path/to/config.hcl ...
...
```

Or start a Vault server in dev mode:

```sh
$ vault server -dev -dev-root-token-id=root -dev-plugin-dir="path/to/plugin/directory"
```

Once the server is started, register the plugin in the Vault server's [plugin catalog](https://developer.hashicorp.com/vault/docs/plugins/plugin-architecture#plugin-catalog):

```sh
$ SHA256=$(openssl dgst -sha256 bin/vault-plugin-auth-azure | cut -d ' ' -f2)
$ vault plugin register \
        -sha256=$SHA256 \
        -command="vault-plugin-auth-azure" \
        auth azure-plugin
...
Success! Data written to: sys/plugins/catalog/azure-plugin
```

Finally, enable the auth method to use this plugin:

```sh
$ vault auth enable azure-plugin
...

Successfully enabled 'plugin' at 'azure-plugin'!
```

### Azure Environment Setup

A Terraform [configuration](bootstrap/terraform) is included in this repository that
automates provisioning of Azure resources necessary to configure and authenticate
using the auth method. By default, the resources are created in `westus2`. See 
[variables.tf](bootstrap/terraform/variables.tf) for the available variables.

Before applying the Terraform configuration, you'll need to:

1. [Authenticate](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs#authenticating-to-azure)
   the Terraform provider to Azure
2. Provide an SSH public key for access to the Azure VM via the `TF_VAR_ssh_public_key_path`
   variable (defaults to `~/.ssh/id_rsa.pub`)

The Terraform configuration will create:

* A service principal with necessary role assignments
* A virtual network, subnet, and security group with only SSH access from your local 
  machine's public IP address
* A linux virtual machine instance

To provision the Azure resources, run the following:

```sh
$ make setup-env   
```

The `local_environment_setup.sh` file will be created in the `bootstrap/terraform`
directory as a result of running `make setup-env`. This file contains environment
variables needed to configure the auth method. The values can also be accessed
via `terraform output`.

To access the virtual machine via SSH:

```sh
ssh adminuser@${VM_IP_ADDRESS}
```

Once you're finished with plugin development, you can run the following to
destroy the Azure resources:

```sh
$ make teardown-env   
```

### Configure Plugin

A [scripted configuration](bootstrap/configure.sh) of the plugin is provided in
this repository. You can use the script or manually configure the auth method
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
$ PLUGIN_NAME=vault-plugin-auth-azure \
  PLUGIN_DIR=$GOPATH/vault-plugins \
  PLUGIN_PATH=local-auth-azure \
  make configure
```

### Tests

If you are developing this plugin and want to verify it is still
functioning, we recommend running the tests.

To run the tests, invoke `make test`:

```sh
$ make test
```

You can also specify a `TESTARGS` variable to filter tests like so:

```sh
$ make test TESTARGS='--run=TestConfig'
```
