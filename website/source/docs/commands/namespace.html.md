---
layout: "docs"
page_title: "namespace - Command"
sidebar_title: "<code>namespace</code>"
sidebar_current: "docs-commands-namespace"
description: |-
  The "namespace" command groups subcommands for interacting with namespaces.
---

# namespace

The `namespace` command groups subcommands for interacting with namespaces.

## Examples

List all namespaces:

```text
$ vault namespace list
```

Create a namespace at the path `ns1/`:

```text
$ vault namespace create ns1/
```

Delete the namespace at path `ns1/`:

```text
$ vault namespace delete ns1/
```

Lookup the namespace information at path `ns1/`:

```text
$ vault namespace lookup ns1/
```

## Usage

```text
Usage: vault namespace <subcommand> [options] [args]

  This command groups subcommands for interacting with Vault namespaces.
  These set of subcommands operate on the context of the namespace that the
  current logged in token belongs to.

Subcommands:
    create    Create a new namespace
    delete    Delete an existing namespace
    list      List child namespaces
    lookup    Look up an existing namespace
```

For more information, examples, and usage about a subcommand, click on the name
of the subcommand in the sidebar.
