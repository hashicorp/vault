---
layout: "docs"
page_title: "operator - Command"
sidebar_current: "docs-commands-operator"
description: |-
  The "operator" command groups subcommands for operators interacting with
  Vault. Most users will not need to interact with these commands.
---

# operator

The `operator` command groups subcommands for operators interacting with Vault.
Most users will not need to interact with these commands.

## Examples

Initialize a new Vault cluster:

```text
$ vault operator init
Unseal Key 1: sP/4C/fwIDjJmHEC2bi/1Pa43uKhsUQMmiB31GRzFc0R
Unseal Key 2: kHkw2xTBelbDFIMEgEC8NVX7NDSAZ+rdgBJ/HuJwxOX+
Unseal Key 3: +1+1ZnkQDfJFHDZPRq0wjFxEuEEHxDDOQxa8JJ/AYWcb
Unseal Key 4: cewseNJTLovmFrgpyY+9Hi5OgJlJgGGCg7PZyiVdPwN0
Unseal Key 5: wyd7rMGWX5fi0k36X4e+C4myt5CoTmJsHJ0rdYT7BQcF

Initial Root Token: 6662bb4a-afd0-4b6b-faad-e237fb564568

# ...
```

Force a Vault to resign leadership in a cluster:

```text
$ vault operator step-down
Success! Stepped down: https://vault.rocks
```

Rotate Vault's underlying encryption key:

```text
$ vault operator rotate
Success! Rotated key

Key Term        2
Install Time    01 Jan 07 12:30 UTC
```

## Usage

```text
Usage: vault operator <subcommand> [options] [args]

  # ...

Subcommands:
    generate-root    Generates a new root token
    init             Initializes a server
    key-status       Provides information about the active encryption key
    rekey            Generates new unseal keys
    rotate           Rotates the underlying encryption key
    seal             Seals the Vault server
    step-down        Forces Vault to resign active duty
    unseal           Unseals the Vault server
```

For more information, examples, and usage about a subcommand, click on the name
of the subcommand in the sidebar.
