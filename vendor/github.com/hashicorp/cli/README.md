# Go CLI Library [![GoDoc](https://godoc.org/github.com/hashicorp/cli?status.png)](https://pkg.go.dev/github.com/hashicorp/cli)

cli is a library for implementing command-line interfaces in Go.
cli is the library that powers the CLI for
[Packer](https://github.com/hashicorp/packer),
[Consul](https://github.com/hashicorp/consul),
[Vault](https://github.com/hashicorp/vault),
[Terraform](https://github.com/hashicorp/terraform),
[Nomad](https://github.com/hashicorp/nomad), and more.

## Features

* Easy sub-command based CLIs: `cli foo`, `cli bar`, etc.

* Support for nested subcommands such as `cli foo bar`.

* Optional support for default subcommands so `cli` does something
  other than error.

* Support for shell autocompletion of subcommands, flags, and arguments
  with callbacks in Go. You don't need to write any shell code.

* Automatic help generation for listing subcommands.

* Automatic help flag recognition of `-h`, `--help`, etc.

* Automatic version flag recognition of `-v`, `--version`.

* Helpers for interacting with the terminal, such as outputting information,
  asking for input, etc. These are optional, you can always interact with the
  terminal however you choose.

* Use of Go interfaces/types makes augmenting various parts of the library a
  piece of cake.

## Example

Below is a simple example of creating and running a CLI

```go
package main

import (
	"log"
	"os"

	"github.com/hashicorp/cli"
)

func main() {
	c := cli.NewCLI("app", "1.0.0")
	c.Args = os.Args[1:]
	c.Commands = map[string]cli.CommandFactory{
		"foo": fooCommandFactory,
		"bar": barCommandFactory,
	}

	exitStatus, err := c.Run()
	if err != nil {
		log.Println(err)
	}

	os.Exit(exitStatus)
}
```

