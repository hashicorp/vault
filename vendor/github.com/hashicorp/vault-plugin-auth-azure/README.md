# Vault Plugin: Azure Auth Backend [![Build Status](https://travis-ci.org/hashicorp/vault-plugin-auth-azure.svg?branch=master)](https://travis-ci.org/hashicorp/vault-plugin-auth-azure)


This is a standalone backend plugin for use with [Hashicorp Vault](https://www.github.com/hashicorp/vault).
This plugin allows for Azure Managed Service Identities to authenticate with Vault.

**Please note**: We take Vault's security and our users' trust very seriously. If you believe you have found a security issue in Vault, _please responsibly disclose_ by contacting us at [security@hashicorp.com](mailto:security@hashicorp.com).

## Quick Links
    - Vault Website: https://www.vaultproject.io
    - Azure Auth Docs: https://www.vaultproject.io/docs/auth/azure.html
    - Main Project Github: https://www.github.com/hashicorp/vault


## Getting Started

This is a [Vault plugin](https://www.vaultproject.io/docs/internals/plugins.html)
and is meant to work with Vault. This guide assumes you have already installed Vault
and have a basic understanding of how Vault works.

Otherwise, first read this guide on how to [get started with Vault](https://www.vaultproject.io/intro/getting-started/install.html).

To learn specifically about how plugins work, see documentation on [Vault plugins](https://www.vaultproject.io/docs/internals/plugins.html).

## Security Model

The current authentication model requires providing Vault with a token generated using Azure's Managed Service Identity, which can be used to make authenticated calls to Azure. This token should not typically be shared, but in order for Azure to be treated as a trusted third party, Vault must validate something that Azure has cryptographically signed and that conveys the identity of the token holder.

## Usage

Please see [documentation for the plugin](https://www.vaultproject.io/docs/auth/azure.html)
on the Vault website.

This plugin is currently built into Vault and by default is accessed
at `auth/azure`. To enable this in a running Vault server:

```sh
$ vault auth enable azure
Successfully enabled 'azure' at 'azure'!
```

To see all the supported paths, see the [Azure auth backend docs](https://www.vaultproject.io/docs/auth/azure.html).

## Developing

If you wish to work on this plugin, you'll first need
[Go](https://www.golang.org) installed on your machine.

For local dev first make sure Go is properly installed, including
setting up a [GOPATH](https://golang.org/doc/code.html#GOPATH).
Next, clone this repository into
`$GOPATH/src/github.com/hashicorp/vault-plugin-auth-azure`.
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
will be specified as the [`plugin_directory`](https://www.vaultproject.io/docs/configuration/index.html#plugin_directory)
in the Vault config used to start the server.

```json
...
plugin_directory = "path/to/plugin/directory"
...
```

Start a Vault server with this config file:
```sh
$ vault server -config=path/to/config.json ...
...
```

Once the server is started, register the plugin in the Vault server's [plugin catalog](https://www.vaultproject.io/docs/internals/plugins.html#plugin-catalog):

```sh
$ vault write sys/plugins/catalog/azure \
        sha_256=<expected SHA256 Hex value of the plugin binary> \
        command="vault-plugin-auth-azure"
...
Success! Data written to: sys/plugins/catalog/azure
```

Note you should generate a new sha256 checksum if you have made changes
to the plugin. Example using openssl:

```sh
openssl dgst -sha256 $GOPATH/vault-plugin-auth-azure
...
SHA256(.../go/bin/vault-plugin-auth-azure)= 896c13c0f5305daed381952a128322e02bc28a57d0c862a78cbc2ea66e8c6fa1
```

Enable the auth plugin backend using the Azure auth plugin:

```sh
$ vault auth enable -plugin-name='azure' plugin
...

Successfully enabled 'plugin' at 'azure'!
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
