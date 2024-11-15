# Vault Plugin: AliCloud Auth Backend [![Build Status](https://travis-ci.org/hashicorp/vault-plugin-auth-alicloud.svg?branch=master)](https://travis-ci.org/hashicorp/vault-plugin-auth-alicloud)

This is a standalone backend plugin for use with [Hashicorp Vault](https://www.github.com/hashicorp/vault).
This plugin allows authentication to Vault using Resource Access Management (RAM).

**Please note**: We take Vault's security and our users' trust very seriously. If you believe you have found a security issue in Vault, _please responsibly disclose_ by contacting us at [security@hashicorp.com](mailto:security@hashicorp.com).

## Quick Links
    - Vault Website: https://www.vaultproject.io
    - AliCloud Auth Docs: https://www.vaultproject.io/docs/auth/alicloud.html
    - Main Project Github: https://www.github.com/hashicorp/vault

## Getting Started

This is a [Vault plugin](https://www.vaultproject.io/docs/internals/plugins.html)
and is meant to work with Vault. This guide assumes you have already installed Vault
and have a basic understanding of how Vault works.

Otherwise, first read this guide on how to [get started with Vault](https://www.vaultproject.io/intro/getting-started/install.html).

To learn specifically about how plugins work, see documentation on [Vault plugins](https://www.vaultproject.io/docs/internals/plugins.html).

## Security Model

This authentication model places Vault in the middle of a call between a client and AliCloud's "GetCallerIdentity" method. Based on AliCloud's response, it grants an access token based on pre-configured roles.

## Usage

Please see [documentation for the plugin](https://www.vaultproject.io/docs/auth/alicloud.html)
on the Vault website.

This plugin is currently built into Vault and by default is accessed
at `auth/alicloud`. To enable this in a running Vault server:

```sh
$ vault auth enable alicloud
Successfully enabled 'alicloud' at 'alicloud'!
```

To see all the supported paths, see the [AliCloud auth backend docs](https://www.vaultproject.io/docs/auth/alicloud.html).

## Developing

If you wish to work on this plugin, you'll first need
[Go](https://www.golang.org) installed on your machine.

For local dev first make sure Go is properly installed, including
setting up a [GOPATH](https://golang.org/doc/code.html#GOPATH).
Next, clone this repository into
`$GOPATH/src/github.com/hashicorp/vault-plugin-auth-alicloud`.

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

```hcl
plugin_directory = "path/to/plugin/directory"
```

Start a Vault server with this config file:
```sh
$ vault server -config=path/to/config.hcl ...
...
```

Once the server is started, register the plugin in the Vault server's [plugin catalog](https://developer.hashicorp.com/vault/docs/plugins/plugin-architecture#plugin-catalog):

```sh
$ vault plugin register \
        -sha256=<SHA256 Hex value of the plugin binary> \
        -command="vault-plugin-auth-alicloud" \
        auth \
        alicloud
```

Note you should generate a new sha256 checksum if you have made changes
to the plugin. Example using openssl:

```sh
openssl dgst -sha256 $GOPATH/vault-plugin-auth-alicloud
...
SHA256(.../go/bin/vault-plugin-auth-alicloud)= 896c13c0f5305daed381952a128322e02bc28a57d0c862a78cbc2ea66e8c6fa1
```

Enable the auth plugin backend using the AliCloud auth plugin:

```sh
$ vault auth enable -plugin-name='alicloud' plugin
...

Successfully enabled 'plugin' at 'alicloud'!
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

**Warning:** The acceptance tests create/destroy/modify *real resources*,
which may incur real costs in some cases. In the presence of a bug,
it is technically possible that broken backends could leave dangling
data behind. Therefore, please run the acceptance tests at your own risk.
At the very least, we recommend running them in their own private
account for whatever backend you're testing.

Acceptance tests require the following environment variables.
```sh
export VAULT_AUTH_ACC_TEST_ROLE_ARN=<myrolearn>
export VAULT_AUTH_ACC_TEST_ACCESS_KEY_ID=<myaccesskeyid>
export VAULT_AUTH_ACC_TEST_SECRET_KEY=<mysecretkey>
export VAULT_ACC=1
```

To run the acceptance tests, invoke `make testacc`:

```sh
$ make testacc
```

