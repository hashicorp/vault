---
layout: "docs"
page_title: "policy write - Command"
sidebar_current: "docs-commands-policy-write"
description: |-
  The "policy write" command uploads a policy with name NAME from the contents
  of a local file PATH or stdin. If PATH is "-", the policy is read from stdin.
  Otherwise, it is loaded from the file at the given path on the local disk.
---

# policy write

The `policy write` command uploads a policy with name NAME from the contents of
a local file PATH or stdin. If PATH is "-", the policy is read from stdin.
Otherwise, it is loaded from the file at the given path on the local disk.

For details on the policy syntax, please see the [policy
documentation](/docs/concepts/policies.html).

## Examples

Upload a policy named "my-policy" from "/tmp/policy.hcl" on the local disk:

```text
$ vault policy write my-policy /tmp/policy.hcl
```

Upload a policy from stdin:

```text
$ cat my-policy.hcl | vault policy write my-policy -
```

## Usage

There are no flags beyond the [standard set of flags](/docs/commands/index.html)
included on all commands.
