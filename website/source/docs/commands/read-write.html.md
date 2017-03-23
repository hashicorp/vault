---
layout: "docs"
page_title: "Reading and Writing Data"
sidebar_current: "docs-commands-readwrite"
description: |-
  The Vault CLI can be used to read, write, and delete secrets. This page documents how to do this.
---

# Reading and Writing Data with the CLI

The Vault CLI can be used to read, write, and delete data from Vault.
This data might be raw secrets, it might be configuration for
a backend, etc. Whatever it is, the interface to read and write data
to Vault is the same.

To determine what paths can be used to read and write data,
please use the built-in [help system](/docs/commands/help.html)
to discover the paths.

## Writing Data

To write data to Vault, you use `vault write`. It is very easy to use:

```
$ vault write secret/password \
    value=itsasecret
...
```

The above writes a value to `secret/password`. As mentioned in the getting
started guide, multiple values can also be written:

```
$ vault write secret/password \
    value=itsasecret \
    username=something
...
```

For the `secret/` backend, the key/value pairs are arbitrary and can be
anything. For other backends, they're generally more strict, and the
help system can tell you what data to send to Vault.

In addition to writing key/value pairs, Vault can write from a variety
more sources.

#### stdin

`vault write` can read data to write from stdin by using "-" as the value.
If you use "-" as the entire argument, then Vault expects to read a JSON
object from stdin. The example below is equivalent to the first example
above.

```
$ echo -n '{"value":"itsasecret"}' | vault write secret/password -
...
```

You can also add more values in addition to "-" on the command-line.
Depending on their ordering will determine if they overwrite the values
from stdin: if they're after the "-" (positionally on the command-line),
then they will overwrite it, otherwise the values in stdin will overwrite
the command line values.

In addition to reading full JSON objects, Vault can read just a JSON
value. The example below is also identical to the previous example.

```
$ echo -n "itsasecret" | vault write secret/password value=-
...
```

#### Files

`vault write` can read data from files as well. The usage is very similar
to stdin as documented above, but the syntax is `@filename`. Example:

```
$ cat data.json
{ "value": "itsasecret" }

$ vault write secret/password @data.json
...
```

And, just like stdin, you can also specify just values:

```
$ cat data.txt
itsasecret

$ vault write secret/password value=@data.txt
```

Unlike stdin, you can specify multiple files, repeat files, etc. all
on the command line. Reading from files is very useful for complex data.

## Reading Data

Data can be read using `vault read`. This command is very simple:

```
$ vault read secret/password
Key             	Value
---             	-----
refresh_interval	768h0m0s
value           	itsasecret
```

You can use the `-format` flag to get various different formats out
from the command. Some formats are easier to use in different environments
than others.

You can also use the `-field` flag to extract an individual field
from the secret data.

```
$ vault read -field=value secret/password
itsasecret
```

