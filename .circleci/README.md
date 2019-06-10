# CircleCI config

This directory contains both the source code (under `./config/`)
and the generated single-file `config.yml`
which defines the CircleCI workflows for this project.

The Makefile in this directory generates the `./config.yml`
in CircleCI 2.0 syntax,
from the tree rooted at `./config/`,
which contains files in CircleCI 2.1 syntax.
CircleCI supports [generating a single config file from many],
using the `$ circleci config pack` command.
It also supports [expanding 2.1 syntax to 2.0 syntax]
using the `$ circleci config process` command.

[generating a single config file from many]: https://circleci.com/docs/2.0/local-cli/#packing-a-config
[expanding 2.1 syntax to 2.0 syntax]: https://circleci.com/docs/2.0/local-cli/#processing-a-config

## Prerequisites

You will need the [CircleCI CLI tool] installed and working,
at least version `0.1.5607`.

```
$ circleci version
0.1.5607+f705856
```

NOTE: It is recommended to [download this tool directly from GitHub Releases].
Do not install it using Homebrew, as this version cannot be easily updated.
It is also not recommended to pipe curl to bash (which CircleCI recommend) for security reasons!

[CircleCI CLI tool]: https://circleci.com/docs/2.0/local-cli/
[download this tool directly from GitHub Releases]: https://github.com/CircleCI-Public/circleci-cli/releases

## How to make changes

Before making changes, be sure to understand the layout
of the `./config/` file tree, as well as circleci 2.1 syntax.
See the [Syntax and layout] section below.

To update the config, you should edit, add or remove files
in the `./config/` directory,
and then run `make ci-config`.
If that's successful,
you should then commit every `*.yml` file in the tree rooted in this directory.
That is: you should commit both the source under `./config/`
and the generated file `./config.yml` at the same time, in the same commit.
Do not edit the `./config.yml` file directly, as you will lose your changes
next time `make ci-config` is run.

[Syntax and layout]: #syntax-and-layout

### Verifying `./config.yml`

To check whether or not the current `./config.yml` is up to date with the source,
and whether it is valid, run `$ make ci-verify`.
Note that `$ make ci-verify` should be run in CI,
as well as by a local git commit hook,
to ensure we never commit files that are invalid or out of date.

#### Example shell session

```sh
$ make ci-config
config.yml updated 
$ git add -A . # The -A makes sure to include deletions/renames etc.
$ git commit -m "ci: blah blah blah"
Changes detected in .circleci/, running 'make -C .circleci ci-verify'
--> Generated config.yml is up to date!
--> Config file at config.yml is valid.
```

### Syntax and layout

It is important to understand the layout of the config directory.
Read the documentation on [packing a config] for a full understanding
of how multiple YAML files are merged by the circleci CLI tool.

[packing a config]: https://circleci.com/docs/2.0/local-cli/#packing-a-config

Here is an example file tree (with comments added afterwards):

```sh
$ tree . 
.
├── Makefile
├── README.md # This file.
├── config    # The source code for config.yml is rooted here.
│   ├── @config.yml # Files beginning with @ are treated specially by `circleci config pack`
│   ├── commands    # Subdirectories of config become top-level keys.
│   │   └── go_test.yml # Filenames (minus .yml) become top-level keys under their parent (in this case "commands").
│   │                   # The contents of go_test.yml therefore are placed at: .commands.go_test:
│   └── jobs        # jobs also becomes a top-level key under config...
│       ├── build-go-dev.yml # ...and likewise filenames become keys under their parent.
│       ├── go-mod-download.yml
│       ├── install-ui-dependencies.yml
│       ├── test-go-race.yml
│       ├── test-go.yml
│       └── test-ui.yml
└── config.yml # The generated file in 2.0 syntax.
```

About those `@` files... Preceding a filename with `@`
indicates to `$ circleci config pack` that the contents of this YAML file
should be at the top-level, rather than underneath a key named after their filename.
This naming convention is unfortunate as it breaks autocompletion in bash,
but there we go.

### Why not just use YAML references?

YAML references only work within a single file,
this is because `circleci config pack` is not a text-level packer,
but rather stitches together the structures defined in each YAML
file according to certain rules.
Therefore it must parse each file separately,
and YAML references are handled by the parser.
