---
layout: "docs"
page_title: "kv put - Command"
sidebar_title: "<code>put</code>"
sidebar_current: "docs-commands-kv-put"
description: |-
  The "put" command writes the data to the given path in the key-value store. The data can be of
  any type.
---

# kv put

The `kv put` command writes the data to the given path in the key-value store. The data can be of
  any type. 

## Examples

To write data to a given path in the key-value store:

```text
$  vault kv put secret/foo bar=baz
```
To write data from a file on disk to a given path in the key-value store:

```text
$ vault kv put secret/foo @data.json
```

To write data from stdin to a given path in the key-value store:

```text
$ echo "abcd1234" | vault kv put secret/foo bar=-
```

To perform a Check-And-Set operation, specify the -cas flag with the
  appropriate version number corresponding to the key you want to perform
  the CAS operation on:

```text
$ vault kv put -cas=1 secret/foo bar=baz
```


## Usage

```text 
$ vault kv put [options] KEY [DATA]


Common Options:

  -cas=<int>
      Specifies to use a Check-And-Set operation. If not set the write will be
      allowed. If set to 0 a write will only be allowed if the key doesn’t
      exist. If the index is non-zero the write will only be allowed if
      the key’s current version matches the version specified in the cas
      parameter. The default is -1.
```
### Output Options

- `-field` `(string: "field_name")`
  Print only the field with the given name. Specifying this option will
  take precedence over other formatting directives. The result will not
  have a trailing newline making it idea for piping to other processes.
  
- `-format` `(string: "table")` - Print the output in the given format. Valid
  formats are "table", "json", or "yaml". This can also be specified via the
  `VAULT_FORMAT` environment variable.
