---
layout: "docs"
page_title: "kv patch - Command"
sidebar_title: "<code>patch</code>"
sidebar_current: "docs-commands-kv-patch"
description: |-
  The "patch" command writes the data to the given path in the key-value store. The data can be of
  any type. Patch can also be used to add another field to an existing path in the key-value store.
---

# kv delete

The `kv patch` command writes the data to the given path in the key-value store. The data can be of
  any type. Patch can also be used to add another field to an existing path in the key-value store.

## Examples

To write data to a given path in the key-value store:

```text
$  vault kv patch secret/foo bar=baz

Key              Value
---              -----
created_time     2019-01-12T22:50:59.754692971Z
deletion_time    n/a
destroyed        false
version          1
```

To add another field to a given path in the key-value store:

```text
$ vault kv patch secret/foo bar2=baz2
#...

====== Metadata ======
Key              Value
---              -----
created_time     2019-01-12T22:50:10.908178624Z
deletion_time    n/a
destroyed        false
version          2

==== Data ====
Key     Value
---     -----
bar     baz
bar2    baz2
```

## Usage

```text 
$  vault kv patch [options] KEY [DATA]
```

### Output Options

- `-field` `(string: "field_name")`
  Print only the field with the given name. Specifying this option will
  take precedence over other formatting directives. The result will not
  have a trailing newline making it idea for piping to other processes.
  
- `-format` `(string: "table")` - Print the output in the given format. Valid
  formats are "table", "json", or "yaml". This can also be specified via the
  `VAULT_FORMAT` environment variable.
