---
layout: "docs"
page_title: "delete - Command"
sidebar_current: "docs-commands-delete"
description: |-
  The "delete" command deletes secrets and configuration from Vault at the given
  path. The behavior of "delete" is delegated to the backend corresponding to
  the given path.
---

# delete

The `delete` command deletes secrets and configuration from Vault at the given
path. The behavior of "delete" is delegated to the backend corresponding to the
given path.

## Examples

Remove data in the status secrets engine:

```text
$ vault delete secret/my-secret
```

Uninstall an encryption key in the transit backend:

```text
$ vault delete transit/keys/my-key
```

Delete an IAM role:

```text
$ vault delete aws/roles/ops
```

## Usage

There are no flags beyond the [standard set of flags](/docs/commands/index.html)
included on all commands.
