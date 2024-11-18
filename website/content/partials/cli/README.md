## Terminology

non-runnable    a command with no side-effects other than printing help text.
command tree    a hierarchical graph of CLI commands where internal nodes are
                non-runnable and leaf nodes represent runnable CLI commands.
root command    the root node in a tree or subtree of CLI commands. For example,
                `plugin` is the root node for all plugin commands and
                `plugin runtime` is the root node for runtime commands.
command family - the top-most root command for a collection of CLI commands.
                 For example `audit` or `plugin`.

### Exceptions :(

The `agent` family of commands is malformed. Rather than having a root node
(`agent`) with two subcommands (`agent start` and `agent generate-config`), the
root command is runnable.


## Why partials?

We document CLI command arguments, options, and flags as partials:

- as a first step toward templatizing and autogenerating the CLI command pages.
- to make it easier to include and maintain elements shared across commands in
  the same family.
- to make it easier to include and maintain elements shared across command
  families.
- to make it easier to include information about standard flags on the command
  pages.


## Directory structure

partials/cli/<command-family>          partials specific to a command family
partials/cli/<command-family>/args     command-family arguments
partials/cli/<command-family>/flags    command-family flags
partials/cli/<command-family>/options  command-family options

partials/cli/shared          partials for elements shared across command families
partials/cli/shared/args     shared arguments (does not exist yet)
partials/cli/shared/flags    shared flags
partials/cli/shared/options  shared options (does not exist yet)

partials/global-settings      partials for standard/global elements
partials/global-settings/flags  global flags (e.g., `-header`)
partials/global-settings/env    global environment variables (e.g., `VAULT_LICENSE`)
partials/global-settings/both  elements that exits as flags and variables

## Partial templates

- If the element is shared across command families, but not applicable to **all**
  command families, it belongs under `partials/cli/shared`
- If the element is a flag with a cooresponding environment variable but **does not**
  apply to all commands, talk with a technical writer before creating your
  partials.
- If the element is required, use `<required>` for the default entry.
- Include `-` as part of the name for flag names **except for anchor IDs**.
- Use `=` in example text for options
- Omit `=` in example text for flags

### Template 1

Use the following template for elements that exist exclusively as arguments,
flags, or options:

<a id="COMMAND_ROOT-[arg | option | flag]-NAME" />

**`NAME (TYPE : DEFAULT)`**

DESCRIPTION

**Example**: `EXAMPLE_OF_VALID_USE`
