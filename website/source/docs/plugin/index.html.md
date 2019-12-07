---
layout: "docs"
page_title: "Custom Plugin Backends"
sidebar_title: "Plugin Backends"
sidebar_current: "docs-plugin"
description: |-
  Plugin backends are mountable backends that are implemented unsing Vault's plugin system.
---

# Custom Plugin Backends

Plugin backends are the components in Vault that can be implemented separately from Vault's
builtin backends. These backends can be either authentication or secrets engines.

The [`api_addr`][api_addr] must be set in order for the plugin process to establish
communication with the Vault server during mount time. If the storage backend
has HA enabled and supports automatic host address detection (e.g. Consul),
Vault will automatically attempt to determine the `api_addr` as well.

Detailed information regarding the plugin system can be found in the
[internals documentation](https://www.vaultproject.io/docs/internals/plugins.html).

# Enabling/Disabling Plugin Backends

Before a plugin backend can be mounted, it needs to be registered via the
[plugin catalog](https://www.vaultproject.io/docs/internals/plugins.html#plugin-catalog). After
the plugin is registered, it can be mounted by specifying the registered plugin name:

```text
$ vault secrets enable -path=my-secrets passthrough-plugin
Success! Enabled the passthrough-plugin secrets engine at: my-secrets/
```

Listing secrets engines will display secrets engines that are mounted as
plugins:

```text
$ vault secrets list
Path         Type       Accessor            Plugin              Default TTL  Max TTL  Force No Cache  Replication Behavior  Description
my-secrets/  plugin     plugin_deb84140     passthrough-plugin  system       system   false           replicated
```

Disabling a plugin backend is the identical to disabling internal secrets engines:

```text
$ vault secrets disable my-secrets
```

# Upgrading Plugins

Vault executes plugin binaries when they are configured and roles established
around them. The binary cannot be modified or replaced while running, so
upgrades cannot be performed by simply swapping the binary and updating the hash
in the plugin catalog.

Instead, you can restart or reload a plugin with the
`sys/plugins/reload/backend` [API][plugin_reload_api]. Follow these steps to
replace or upgrade a Vault plugin binary:

1. Register plugin_v1 to the catalog
2. Mount the plugin backend
3. Register plugin_v2 to the catalog under the same plugin name, but with
updated command to run plugin_v2 and updated sha256 of plugin_v2
4. Trigger a plugin reload with sys/plugins/reload/backend to reload all mounted
backends using that plugin, or just a subset of the mounts using either the
plugin or mounts parameter.

Until step 4, the mount will still use plugin_v1, and when the reload is
triggered, Vault will kill plugin_v1’s process and start a plugin_v2 process.

[api_addr]: /docs/configuration/index.html#api_addr
[plugin_reload_api]: /api/system/plugins-reload-backend.html
