---
layout: "docs"
page_title: "operator migrate - Command"
sidebar_title: "<code>migrate</code>"
sidebar_current: "docs-commands-operator-migrate"
description: |-
  The "operator migrate" command copies data between storage backends to facilitate
  migrating Vault between configurations. It operates directly at the storage
  level, with no decryption involved.
---

# operator migrate

The `operator migrate` command copies data between storage backends to facilitate
migrating Vault between configurations. It operates directly at the storage
level, with no decryption involved. Keys in the destination storage backend will
be overwritten, and the destination should _not_ be initialized prior to the
migrate operation. The source data is not modified, with the exception of a small lock
key added during migration.

This is intended to be an offline operation to ensure data consistency, and Vault
will not allow starting the server if a migration is in progress.

## Examples

Migrate all keys:

```text
$ vault operator migrate -config migrate.hcl

2018-09-20T14:23:23.656-0700 [INFO ] copied key: data/core/seal-config
2018-09-20T14:23:23.657-0700 [INFO ] copied key: data/core/wrapping/jwtkey
2018-09-20T14:23:23.658-0700 [INFO ] copied key: data/logical/fd1bed89-ffc4-d631-00dd-0696c9f930c6/31c8e6d9-2a17-d98f-bdf1-aa868afa1291/archive/metadata
2018-09-20T14:23:23.660-0700 [INFO ] copied key: data/logical/fd1bed89-ffc4-d631-00dd-0696c9f930c6/31c8e6d9-2a17-d98f-bdf1-aa868afa1291/metadata/5kKFZ4YnzgNfy9UcWOzxxzOMpqlp61rYuq6laqpLQDnB3RawKpqi7yBTrawj1P
...
```

Migration is done in a consistent, sorted order. If the migration is halted or
exits before completion (e.g. due to a connection error with a storage backend),
it may be resumed from an arbitrary key prefix:

```text
$ vault operator migrate -config migrate.hcl -start "data/logical/fd"
```

## Configuration

The `operator migrate` command uses a dedicated configuration file to specify the source
and destination storage backends. The format of the storage stanzas is identical
to that used to [configure Vault](/docs/configuration/storage/index.html),
with the only difference being that two stanzas are required: `storage_source` and `storage_destination`.

```hcl
storage_source "mysql" {
  username = "user1234"
  password = "secret123!"
  database = "vault"
}

storage_destination "consul" {
  address = "127.0.0.1:8500"
  path    = "vault"
}
```

## Usage

The following flags are available for the `operator migrate` command.

- `-config` `(string: <required>)` - Path to the migration configuration file.

- `-start` `(string: "")` - Migration starting key prefix. Only keys at or after this value will be copied.

- `-reset` - Reset the migration lock. A lock file is added during migration to prevent
  starting the Vault server or another migration. The `-reset` option can be used to
  remove a stale lock file if present.
