---
layout: "docs"
page_title: "Custom Plugin Backends"
sidebar_current: "docs-plugin"
description: |-
  Plugin backends are mountable backends that are implemented unsing Vault's plugin system.
---

# Custom Plugin Backends

Plugin backends are the components in Vault that can be implemented separately from Vault's
builtin backends. These backends can be either authentication or secret backends.

Detailed information regarding the plugin system can be found in the
[internals documentation](https://www.vaultproject.io/docs/internals/plugins.html).

# Mounting/unmounting Plugin Backends

Before a plugin backend can be mounted, it needs to be registered via the
[plugin catalog](https://www.vaultproject.io/docs/internals/plugins.html#plugin-catalog). After
the plugin is registered, it can be mounted by specifying the registered plugin name:

```
$ vault mount -path=my-secrets -plugin-name=passthrough-plugin plugin
Successfully mounted plugin 'passthrough-plugin' at 'my-secrets'!
```

Listing mounts will display backends that are mounted as plugins, along with the
name of plugin backend that is mounted:

```
$ vault mounts
Path         Type       Accessor            Plugin              Default TTL  Max TTL  Force No Cache  Replication Behavior  Description
my-secrets/  plugin     plugin_deb84140     passthrough-plugin  system       system   false           replicated
...
```

Unmounting a plugin backend is the identical to unmounting internal backends:

```
$ vault unmount my-secrets
```
