---
layout: "docs"
page_title: "secrets - Command"
sidebar_current: "docs-commands-secrets"
description: |-
  The "secrets" command groups subcommands for interacting with Vault's secrets
  engines.
---

# secrets

The `secrets` command groups subcommands for interacting with Vault's secrets
engines. Each secrets engine behaves differently. Please see the documentation
for more information.

Some secrets engines persist data, some act as data pass-through, and some
generate dynamic credentials. The secrets engine will likely require
configuration after it is mounted. For details on the specific configuration
options, please see the [secrets engine
documentation](/docs/secrets/index.html).

## Examples

Enable a secrets engine:

```text
$ vault secrets enable database
Success! Enabled the database secrets engine at: database/
```

List all secrets engines:

```text
$ vault secrets list
Path          Type         Description
----          ----         -----------
cubbyhole/    cubbyhole    per-token private secret storage
database/     database     n/a
secret/       kv           key/value secret storage
sys/          system       system endpoints used for control, policy and debugging
```

Move a secrets engine to a new path:

```text
$ vault secrets move database/ db-prod/
Success! Moved secrets engine database/ to: db-prod/
```

Tune a secrets engine:

```text
$ vault secrets tune -max-lease-ttl=30m db-prod/
Success! Tuned the secrets engine at: db-prod/
```

Disable a secrets engine:

```text
$ vault secrets disable db-prod/
Success! Disabled the secrets engine (if it existed) at: db-prod/
```

## Usage

```text
Usage: vault secrets <subcommand> [options] [args]

  # ...

Subcommands:
    disable    Disable a secrets engine
    enable     Enable a secrets engine
    list       List enabled secrets engines
    move       Move a secrets engine to a new path
    tune       Tune a secrets engine configuration
```

For more information, examples, and usage about a subcommand, click on the name
of the subcommand in the sidebar.
