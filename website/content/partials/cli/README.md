
We document CLI command arguments, options, and flags as partials:

- as a first step toward templatizing and autogenerating the CLI command pages.
- to make it easier to include and maintain elements shared across commands in
  the same family.
- to make it easier to include and maintain elements shared across command
  families.
- to make it easier to include information about standard flags on the command
  pages.


Partial template for CLI elements (required elements use <required> in place of
a default value):

<a id="COMMAND-[arg | option | flag]-NAME" />

**`-NAME (TYPE : DEFAULT)`**

DESCRIPTION

**Example**: `EXAMPLE_OF_VALID_USE`
