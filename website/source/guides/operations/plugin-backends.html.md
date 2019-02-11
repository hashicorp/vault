---
layout: "guides"
page_title: "Plugin Backends - Guides"
sidebar_title: "Building Plugin Backends"
sidebar_current: "guides-operations-plugin-backends"
description: |-
  Learn how to build, register, and mount a custom plugin backend.
---

# Introduction

Plugin backends utilize the [plugin system][plugin-system] to enable third-party
secrets engines and auth methods.

It is worth noting that even though [database secrets engines][database-backend]
operate under the same underlying plugin mechanism, they are slightly different
in design than plugin backends demonstrated in this guide. The database secrets
engine manages multiple plugins under the same backend mount point, whereas
plugin backends are kv backends that function as either secret or auth methods.

This guide provides steps to build, register, and mount non-database external
plugin backends.

## Setup Vault

Set `plugin_directory` to the desired path in the Vault configuration file.
The path should exist and have proper lockdown on access permissions.

```hcl
# /etc/vault/config.d/plugins.hcl
plugin_directory = "/etc/vault/vault_plugins"
```

If the Vault server is already running, you will need to tell it to reload its
configuration by sending SIGHUP. If you stop and start the Vault server, you
will need to unseal it again.

## Compile Plugin

Build the custom binary, and move it inside the `plugin_directory` path
configured above. This guide uses `mock-plugin` that comes from Vault's
[`logical/plugin/mock`](https://github.com/hashicorp/vault/tree/master/logical/plugin/mock/mock-plugin) package.

Download the source (you would probably use your own plugin):

```sh
$ go get -f -u -d github.com/hashicorp/vault
# ...
$ cd $GOPATH/src/github.com/hashicorp/vault/logical/plugin/mock/mock-plugin
```

Compile the plugin:

```sh
$ go build -o my-mock-plugin
```

Put the plugin in the directory:

```sh
$ mv my-mock-plugin /etc/vault/vault_plugins
```

Alternatively, if you wanted a custom version of a plugin built into Vault, such as AppRole:

```sh
$ cd $GOPATH/src/github.com/hashicorp/vault/builtin/credential/approle/cmd/approle
$ go build
$ mv approle /etc/vault/vault_plugins
```

## Register in Plugin Catalog

Calculate the SHA256 sum of the compiled plugin binary, and use that to register
the plugin into Vault's plugin catalog:

```sh
$ shasum -a 256 /etc/vault/vault_plugins/my-mock-plugin
2c071aafa1b30897e60b79643e77592cb9d1e8f803025d44a7f9bbfa4779d615  /etc/vault/vault_plugins/my-mock-plugin

$ vault write sys/plugins/catalog/secret/my-mock-plugin \
    sha256=2c071aafa1b30897e60b79643e77592cb9d1e8f803025d44a7f9bbfa4779d615 \
    command=my-mock-plugin
Success! Data written to: sys/plugins/catalog/secret/my-mock-plugin
```

## Enable Plugin

Enabling the plugin varies depending on if it's a secrets engine or auth method:

```sh
$ vault secrets enable -path=my-secrets-plugin my-mock-plugin
Success! Enabled the my-mock-plugin plugin at: my-secrets-plugin/
```

If you try to mount this particular plugin as an auth method instead of a
secrets engine, you will get an error:

```sh
$ vault auth enable -path=my-auth-plugin my-mock-plugin
# ...
* cannot mount 'my-mock-plugin' of type 'secret' as an auth method
```

## Perform Operations

Each plugin responds to read, write, list, and delete as its own behavior.

```text
$ vault write my-secrets-plugin/kv/foo value=bar
Key      Value
---      -----
value    bar
```

## Disable Plugin

When you are done using the plugin, disable it.

```text
$ vault secrets disable my-secrets-plugin
Success! Disabled the secrets engine (if it existed) at: my-secrets-plugin/
```

[plugin-system]: /docs/internals/plugins.html
[database-backend]: /docs/secrets/databases/index.html
