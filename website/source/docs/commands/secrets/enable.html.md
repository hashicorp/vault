---
layout: "docs"
page_title: "secrets enable - Command"
sidebar_title: "<code>enable</code>"
sidebar_current: "docs-commands-secrets-enable"
description: |-
  The "secrets enable" command enables an secrets engine at a given path. If an
  secrets engine already exists at the given path, an error is returned. After
  the secrets engine is enabled, it usually needs configuration. The
  configuration varies by secrets engine.
---

# secrets enable

The `secrets enable` command enables an secrets engine at a given path. If an
secrets engine already exists at the given path, an error is returned. After the
secrets engine is enabled, it usually needs configuration. The configuration
varies by secrets engine.

By default, secrets engines are enabled at the path corresponding to their TYPE,
but users can customize the path using the `-path` option.

Some secrets engines persist data, some act as data pass-through, and some
generate dynamic credentials. The secrets engine will likely require
configuration after it is mounted. For details on the specific configuration
options, please see the [secrets engine
documentation](/docs/secrets/index.html).


## Examples

Enable the AWS secrets engine at "aws/":

```text
$ vault secrets enable aws
Success! Enabled the aws secrets engine at: aws/
```

Enable the SSH secrets engine at ssh-prod/:

```text
$ vault secrets enable -path=ssh-prod ssh
```

Enable the database secrets engine with an explicit maximum TTL of 30m:

```text
$ vault secrets enable -max-lease-ttl=30m database
```

Enable a custom plugin (after it is registered in the plugin registry):

```text
$ vault secrets enable -path=my-secrets my-plugin
```

For more information on the specific configuration options and paths, please see
the [secrets engine](/docs/secrets/index.html) documentation.

## Usage

The following flags are available in addition to the [standard set of
flags](/docs/commands/index.html) included on all commands.

- `-default-lease-ttl` `(duration: "")` - The default lease TTL for this secrets
  engine. If unspecified, this defaults to the Vault server's globally
  configured default lease TTL.

- `-description` `(string: "")` - Human-friendly description for the purpose of
  this engine.

- `-force-no-cache` `(bool: false)` - Force the secrets engine to disable
  caching. If unspecified, this defaults to the Vault server's globally
  configured cache settings. This does not affect caching of the underlying
  encrypted data storage.

- `-local` `(bool: false)` - Mark the secrets engine as local-only. Local
  engines are not replicated or removed by replication.

- `-max-lease-ttl` `(duration: "")` The maximum lease TTL for this secrets
  engine. If unspecified, this defaults to the Vault server's globally
  configured maximum lease TTL.

- `-path` `(string: "")` Place where the secrets engine will be accessible. This
  must be unique cross all secrets engines. This defaults to the "type" of the
  secrets engine.
