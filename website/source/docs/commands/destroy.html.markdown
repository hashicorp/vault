---
layout: "docs"
page_title: "Command: destroy"
sidebar_current: "docs-commands-destroy"
description: |-
  The `vault destroy` command is used to destroy the Vault-managed infrastructure.
---

# Command: destroy

The `vault destroy` command is used to destroy the Vault-managed
infrastructure.

## Usage

Usage: `vault destroy [options] [dir]`

Infrastructure managed by Vault will be destroyed. This will ask for
confirmation before destroying.

This command accepts all the flags that the
[apply command](/docs/commands/apply.html) accepts. If `-force` is
set, then the destroy confirmation will not be shown.
