---
layout: "docs"
page_title: "policy delete - Command"
sidebar_current: "docs-commands-policy-delete"
description: |-
  The "policy delete" command deletes the policy named NAME in the Vault server.
  Once the policy is deleted, all tokens associated with the policy are affected
  immediately.
---

# policy delete

The `policy delete` command deletes the policy named NAME in the Vault server.
Once the policy is deleted, all tokens associated with the policy are affected
immediately.

Note that it is not possible to delete the "default" or "root" policies. These
are built-in policies.

## Examples

Delete the policy named "my-policy":

```text
$ vault policy delete my-policy
```

## Usage

There are no flags beyond the [standard set of flags](/docs/commands/index.html)
included on all commands.
