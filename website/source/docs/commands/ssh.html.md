---
layout: "docs"
page_title: "ssh - Command"
sidebar_current: "docs-commands-ssh"
description: |-
  The "ssh" command establishes an SSH connection with the target machine using
  credentials obtained from an SSH secrets engine.
---

# ssh

The `ssh` command establishes an SSH connection with the target machine.

This command uses one of the SSH secrets engines to authenticate and
automatically establish an SSH connection to a host. This operation requires
that the SSH secrets engine is mounted and configured.

The user must have `ssh` installed locally - this command will exec out to it
with the proper commands to provide an "SSH-like" consistent experience.

## Examples

SSH using the OTP mode (requires [sshpass](https://linux.die.net/man/1/sshpass)
for full automation):

```text
$ vault ssh -mode=otp -role=my-role user@1.2.3.4
```

SSH using the CA mode:

```text
$ vault ssh -mode=ca -role=my-role user@1.2.3.4
```

SSH using CA mode with host key verification:

```text
$ vault ssh \
    -mode=ca \
    -role=my-role \
    -host-key-mount-point=host-signer \
    -host-key-hostnames=example.com \
    user@example.com
```

For step-by-step guides and instructions for each of the available SSH
auth methods, please see the corresponding [SSH secrets
engine](/docs/secrets/ssh/index.html).

## Usage

The following flags are available in addition to the [standard set of
flags](/docs/commands/index.html) included on all commands.

### Output Options

- `-field` `(string: "")` - Print only the field with the given name. Specifying
  this option will take precedence over other formatting directives. The result
  will not have a trailing newline making it ideal for piping to other processes.

- `-format` `(string: "table")` - Print the output in the given format. Valid
  formats are "table", "json", or "yaml". This can also be specified via the
  `VAULT_FORMAT` environment variable.

### SSH Options

- `-mode` `(string: "")` - Name of the role to use to generate the key.

- `-mount-point` `(string: "ssh/")` - Mount point to the SSH secrets engine.

- `-no-exec` `(bool: false)` - Print the generated credentials, but do not
  establish a connection.

- `-role` `(string: "")` - Name of the role to use to generate the key.

- `-strict-host-key-checking` `(string: "")` - Value to use for the SSH
  configuration option "StrictHostKeyChecking". The default is ask. This can
  also be specified via the `VAULT_SSH_STRICT_HOST_KEY_CHECKING` environment
  variable.

- `-user-known-hosts-file` `(string: "~/.ssh/known_hosts")` - Value to use for
  the SSH configuration option "UserKnownHostsFile". This can also be specified
  via the `VAULT_SSH_USER_KNOWN_HOSTS_FILE` environment variable.

### CA Mode Options

- `-host-key-hostnames` `(string: "*")` - List of hostnames to delegate for the
  CA. The default value allows all domains and IPs. This is specified as a
  comma-separated list of values. This can also be specified via the
  `VAULT_SSH_HOST_KEY_HOSTNAMES` environment variable.

- `-host-key-mount-point` `(string: "")` - Mount point to the SSH
  secrets engine where host keys are signed. When given a value, Vault will
  generate a custom "known_hosts" file with delegation to the CA at the provided
  mount point to verify the SSH connection's host keys against the provided CA.
  By default, host keys are validated against the user's local "known_hosts"
  file. This flag forces strict key host checking and ignores a custom user
  known hosts file. This can also be specified via the
  `VAULT_SSH_HOST_KEY_MOUNT_POINT` environment variable.

- `-private-key-path` `(string: "~/.ssh/id_rsa")` - Path to the SSH private key
  to use for authentication. This must be the corresponding private key to
  `-public-key-path`.

- `-public-key-path` `(string: "~/.ssh/id_rsa.pub")` - Path to the SSH public
  key to send to Vault for signing.
