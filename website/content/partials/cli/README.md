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
- to make it easier to include and maintain parameters shared across commands in
  the same family.
- to make it easier to include and maintain parameters shared across command
  families.
- to make it easier to include information about standard flags on the command
  pages.


## Directory structure

partials/cli/<command-family>          partials specific to a command family
partials/cli/<command-family>/args     command-family arguments
partials/cli/<command-family>/flags    command-family flags
partials/cli/<command-family>/options  command-family options

partials/cli/shared          partials for parameters shared across some, but not all, command families
partials/cli/shared/args     shared arguments (does not exist yet)
partials/cli/shared/flags    shared flags
partials/cli/shared/options  shared options (does not exist yet)

partials/global-settings        partials for standard/global parameters
partials/global-settings/flags  global flags (e.g., `-header`)
partials/global-settings/env    global environment variables (e.g., `VAULT_LICENSE`)
partials/global-settings/both   parameters that exits as flags and variables

## Partial templates

- Use the parameter name as the file name and "NAME" in the anchor definition,
  even if the use of dashes or underscores is inconsistent with other parameters
  or partial names. For example, if the flag is `-my_weird_flag`, make the
  partial filename `my_weird_flag.mdx` and the anchor ID
  `COMMAND-flag-my_weird_flag`.
- If the parameter is shared across command families, but not applicable to **all**
  command families, it belongs under `partials/cli/shared`
- If the parameter is a flag with a cooresponding environment variable but
  **does not** apply to all commands, talk with a technical writer before
  creating your partials.
- If the parameter is required, use `<required>` for the default entry.
- Include `-` as part of the name for flag names **except for anchor IDs**.
- Use `=` in example text for options
- Omit `=` in example text for flags

### Template 1 - command-specific parameters

Use the following template for parameters that exist as command-exclusively
arguments, flags, or options.

-- Template (start) --

<a id="COMMAND_ROOT-[arg | option | flag]-NAME" />


**`NAME (TYPE : DEFAULT)`**

DESCRIPTION

**Example**: `EXAMPLE_OF_VALID_USE`

-- Template (end) --


### Template 2 - shared parameters

Use the following template for parameters that exist as arguments, flags, or
options that are not global but are shared across more than one command family.

-- Template (start) --

<a id="shared-[arg | option | flag]-NAME" />

**`NAME (TYPE : DEFAULT)`**

DESCRIPTION

**Example**: `EXAMPLE_OF_VALID_USE`

-- Template (end) --
