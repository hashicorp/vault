---
layout: "docs"
page_title: "policy fmt - Command"
sidebar_current: "docs-commands-policy-fmt"
description: |-
  The "policy fmt" formats a local policy file to the policy specification. This
  command will overwrite the file at the given PATH with the properly-formatted
  policy file contents.
---

# policy fmt

The `policy fmt` formats a local policy file to the policy specification. This
command will overwrite the file at the given PATH with the properly-formatted
policy file contents.

## Examples

Format the local file "my-policy.hcl":

```text
$ vault policy fmt my-policy.hcl
```

## Usage

There are no flags beyond the [standard set of flags](/docs/commands/index.html)
included on all commands.
