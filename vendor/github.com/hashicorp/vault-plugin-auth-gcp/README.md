# Vault Plugin: Google Cloud Platform Auth Backend

This is a standalone backend plugin for use with [Hashicorp Vault](https://www.github.com/hashicorp/vault).
This plugin allows for various GCP entities to authenticate with Vault.

Currently, this plugin supports login for:
- IAM service accounts

**Please note**: We take Vault's security and our users' trust very seriously. If you believe you have found a security issue in Vault, _please responsibly disclose_ by contacting us at [security@hashicorp.com](mailto:security@hashicorp.com).

## Quick Links
    - Vault Website: https://www.vaultproject.io
    - GCP Auth BE Docs: https://www.vaultproject.io/docs/auth/gcp.html
    - Main Project Github: https://www.github.com/hashicorp/vault


## Getting Started

This is a [Vault plugin](https://www.vaultproject.io/docs/internals/plugins.html)
and is meant to work with Vault. This guide assumes you have already installed Vault
and have a basic understanding of how Vault works.

Otherwise, first read this guide on how to [get started with Vault](https://www.vaultproject.io/intro/getting-started/install.html).

To learn specifically about how plugins work, see documentation on [Vault plugins](https://www.vaultproject.io/docs/internals/plugins.html).

### Usage

Please see [documentation for the plugin](https://www.vaultproject.io/docs/auth/gcp.html)
on the Vault website.

This plugin is currently built into Vault and by default is accessed
at `auth/gcp`. To enable this in a running Vault server:

```sh
$ vault auth-enable 'gcp'
Successfully enabled 'gcp' at 'gcp'!
```

To see all the supported paths, see the [GCP auth backend docs](https://www.vaultproject.io/docs/auth/gcp.html).

## Developing

If you wish to work on this plugin, you'll first need
[Go](https://www.golang.org) installed on your machine
(version 1.8+ is *required*).

For local dev first make sure Go is properly installed, including
setting up a [GOPATH](https://golang.org/doc/code.html#GOPATH).
Next, clone this repository into
`$GOPATH/src/github.com/hashicorp/vault-gcp-auth-plugin`.
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
$ vault write sys/plugins/catalog/mygcpplugin \
        sha_256=<expected SHA256 Hex value of the plugin binary> \
        command="vault-plugin-auth-gcp"
...
Success! Data written to: sys/plugins/catalog/mygcpplugin
```

Note you should generate a new sha256 checksum if you have made changes
to the plugin. Example using openssl:

```sh
openssl dgst -sha256 $GOPATH/vault-plugin-gcp-auth
...
SHA256(.../go/bin/vault-plugin-auth-gcp)= 896c13c0f5305daed381952a128322e02bc28a57d0c862a78cbc2ea66e8c6fa1
```

Any name can be substituted for the plugin name "mygcpplugin". This
name will be referenced in the next step, where we enable the auth
plugin backend using the GCP auth plugin:

```sh
$ vault auth-enable -plugin-name='mygcpplugin' -path='gcp' plugin
...

Successfully enabled 'plugin' at 'gcp'!
```

### Testing

To run tests, type `make test`. Note: this requires Docker to be installed. If
this exits with exit status 0, then everything is working!

```sh
$ make test
...
```

If you're developing a specific package, you can run tests for just that
package by specifying the `TEST` variable. For example below, only
`vault` package tests will be run.

```sh
$ make test TEST=./vault
...
```

#### Acceptance Tests

This plugin has comprehensive [acceptance tests](https://en.wikipedia.org/wiki/Acceptance_testing)
covering most of the features of this auth backend.

If you are developing this plugin and want to verify it is still
functioning (and you haven't broken anything else), we recommend
running the acceptance tests.

Acceptance tests typically require other environment variables to be set for
things such as access keys. The test itself should error early and tell
you what to set, so it is not documented here.

**Warning:** The acceptance tests create/destroy/modify *real resources*,
which may incur real costs in some cases. In the presence of a bug,
it is technically possible that broken backends could leave dangling
data behind. Therefore, please run the acceptance tests at your own risk.
At the very least, we recommend running them in their own private
account for whatever backend you're testing.

To run the acceptance tests, invoke `make testacc`:

```sh
$ make testacc
```

You can also specify a `TESTARGS` variable to filter tests like so:

```sh
$ make testacc TESTARGS='--run=TestConfig'
```
