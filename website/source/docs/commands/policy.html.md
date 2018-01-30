---
layout: "docs"
page_title: "policy - Command"
sidebar_current: "docs-commands-policy"
description: |-
  The "policy" command groups subcommands for interacting with policies. Users
  can Users can write, read, and list policies in Vault.
---

# policy

The `policy` command groups subcommands for interacting with policies. Users can
Users can write, read, and list policies in Vault.

For more information, please see the [policy
documentation](/docs/concepts/policies.html).

## Examples

List all enabled policies:

```text
$ vault policy list
```

Create a policy named "my-policy" from contents on local disk:

```text
$ vault policy write my-policy ./my-policy.hcl
```

Delete the policy named my-policy:

```text
$ vault policy delete my-policy
```

## Usage

```text
Usage: vault policy <subcommand> [options] [args]

  # ...

Subcommands:
    delete    Deletes a policy by name
    list      Lists the installed policies
    read      Prints the contents of a policy
    write     Uploads a named policy from a file
```

For more information, examples, and usage about a subcommand, click on the name
of the subcommand in the sidebar.
