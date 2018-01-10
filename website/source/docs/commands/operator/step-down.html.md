---
layout: "docs"
page_title: "operator step-down - Command"
sidebar_current: "docs-commands-operator-step-down"
description: |-
  The "operator step-down" forces the Vault server at the given address to step
  down from active duty.
---

# operator step-down

The `operator step-down` forces the Vault server at the given address to step
down from active duty. While the affected node will have a delay before
attempting to acquire the leader lock again, if no other Vault nodes acquire the
lock beforehand, it is possible for the same node to re-acquire the lock and
become active again.

## Examples

Force a Vault server to step down as the leader:

```text
$ vault operator step-down
Success! Stepped down: http://127.0.0.1:8200
```

## Usage

There are no flags beyond the [standard set of flags](/docs/commands/index.html)
included on all commands.
