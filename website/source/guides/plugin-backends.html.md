---
layout: "guides"
page_title: "Plugin Backends - Guides"
sidebar_current: "guides-plugin-backends"
description: |-
  Learn how to build, register, and mount a custom plugin backend.
---

# Introduction

Plugin backends utilize the [plugin system][plugin-system] to enable 
third-party secret and auth backends to be mounted. 

It is worth noting that even though [database backends][database-backend]
operate under the same underlying plugin mechanism, they are slightly different
in design than plugin backends demonstrated in this guide. The database backend 
manages multiple plugins under the same backend mount point, whereas plugin
backends are generic backends that function as either secret or auth backends. 

This guide provides steps to build, register, and mount non-database external
plugin backends.

## Setting up Vault

Set `plugin_directory` to the desired path in the Vault configuration file.
The path should exist and have proper lockdown on access permissions.

```
$ cat vault-config.hcl
...
plugin_directory="/etc/vault/vault_plugins"
...
```

## Build the Plugin Backend

Build the custom backend binary, and move it to the `plugin_directory` path.
In this guide, we will use `mock-plugin` that comes from Vault's 
`logical/plugin/mock` package.

```
$ ls .
main.go

$ ls ..
backend.go  backend_test.go  mock-plugin/  path_internal.go  path_kv.go

$ go build -o mock-plugin main.go

$ mv mock-plugin /etc/vault/vault_plugins
```

## Register the Plugin Into the Plugin Catalog

Start the Vault server. Find out the sha256 sum of the compiled plugin binary,
and use that to register the plugin into Vault's plugin catalog.

```
$ shasum -a 256 /etc/vault/vault_plugins/mock-plugin
2c071aafa1b30897e60b79643e77592cb9d1e8f803025d44a7f9bbfa4779d615  /etc/vault/vault_plugins/mock-plugin

$ vault sys/plugins/catalog/mock-plugin sha_256=2c071aafa1b30897e60b79643e77592cb9d1e8f803025d44a7f9bbfa4779d615 command=mock-plugin
Success! Data written to: sys/plugins/catalog/mock-plugin
```

## Mount the Plugin

```
$ vault mount -path=mock -plugin-name=mock-plugin plugin
Successfully mounted plugin 'mock-plugin' at 'mock'!

$ vault mounts
Path        Type       Accessor            Plugin       Default TTL  Max TTL  Force No Cache  Replication Behavior  Description
cubbyhole/  cubbyhole  cubbyhole_80ef4e30  n/a          n/a          n/a      false           local                 per-token private secret storage
mock/       plugin     plugin_10fc2cce     mock-plugin  system       system   false           replicated
secret/     kv         kv_ef2a14ec         n/a          system       system   false           replicated            key/value secret storage
sys/        system     system_e3a4cccd     n/a          n/a          n/a      false           replicated            system endpoints used for control, policy and debugging
```

## Perform operations on the mount

```
$ vault write mock/kv/foo value=bar
Key  	Value
---  	-----
value	bar
```

## Unmount the plugin

```
$ vault unmount mock
Successfully unmounted 'mock' if it was mounted

$ vault mounts
Path        Type       Accessor            Plugin  Default TTL  Max TTL  Force No Cache  Replication Behavior  Description
cubbyhole/  cubbyhole  cubbyhole_80ef4e30  n/a     n/a          n/a      false           local                 per-token private secret storage
secret/     kv         kv_ef2a14ec         n/a     system       system   false           replicated            key/value secret storage
sys/        system     system_e3a4cccd     n/a     n/a          n/a      false           replicated            system endpoints used for control, policy and debugging
```

[plugin-system]: /docs/internals/plugins.html
[database-backend]: /docs/secrets/databases/index.html