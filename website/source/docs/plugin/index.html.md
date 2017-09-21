---
layout: "docs"
page_title: "Custom Plugin Backends"
sidebar_current: "docs-plugin"
description: |-
  Plugin backends are mountable backends that are implemented unsing Vault's plugin system.
---

# Custom Plugin Backends

Plugin backends are the components in Vault that can be implemented separately from Vault's
builtin backends. These backends can be either authentication or secrets engines.

Detailed information regarding the plugin system can be found in the
[internals documentation](https://www.vaultproject.io/docs/internals/plugins.html).

# Enabling/Disabling Plugin Backends

Before a plugin backend can be mounted, it needs to be registered via the
[plugin catalog](https://www.vaultproject.io/docs/internals/plugins.html#plugin-catalog). After
the plugin is registered, it can be mounted by specifying the registered plugin name:

```text
$ vault secrets enable -path=my-secrets -plugin-name=passthrough-plugin plugin
Successfully mounted plugin 'passthrough-plugin' at 'my-secrets'!
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
