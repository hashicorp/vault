---
layout: "docs"
page_title: "operator generate-root - Command"
sidebar_current: "docs-commands-operator-generate-root"
description: |-
  The "operator generate-root" command generates a new root token by combining a
  quorum of share holders.
---

# operator generate-root

The `operator generate-root` command generates a new root token by combining a
quorum of share holders. One of the following must be provided to start the root
token generation:

- A base64-encoded one-time-password (OTP) provided via the `-otp` flag. Use the
  `-generate-otp` flag to generate a usable value. The resulting token is XORed
  with this value when it is returned. Use the `-decode` flag to output the
  final value.

- A file containing a PGP key or a
  [keybase](/docs/concepts/pgp-gpg-keybase.html) username in the `-pgp-key`
  flag. The resulting token is encrypted with this public key.

An unseal key may be provided directly on the command line as an argument to the
command. If key is specified as "-", the command will read from stdin. If a TTY
is available, the command will prompt for text.

Please see the [generate root guide](/guides/configuration/generate-root.html) for
step-by-step instructions.

## Examples

Generate an OTP code for the final token:

```text
$ vault operator generate-root -generate-otp
```

Start a root token generation:

```text
$ vault operator generate-root -init -otp="..."
```

Enter an unseal key to progress root token generation:

```text
$ vault operator generate-root -otp="..."
```


## Usage

The following flags are available in addition to the [standard set of
flags](/docs/commands/index.html) included on all commands.

### Output Options

- `-format` `(string: "table")` - Print the output in the given format. Valid
  formats are "table", "json", or "yaml". This can also be specified via the
  `VAULT_FORMAT` environment variable.

### Command Options

- `-cancel` `(bool: false)` - Reset the root token generation progress. This
  will discard any submitted unseal keys or configuration.

- `-decode` `(string: "")` - Decode and output the generated root token. This
  option requires the `-otp` flag be set to the OTP used during initialization.

- `-generate-otp` `(bool: false)` - Generate and print a high-entropy
  one-time-password (OTP) suitable for use with the "-init" flag.

- `-init` `(bool: false)` - Start a root token generation. This can only be done
  if there is not currently one in progress.

- `-nonce` `(string; "")`- Nonce value provided at initialization. The same
  nonce value must be provided with each unseal key.

- `-otp` `(string: "")` - OTP code to use with `-decode` or `-init`.

- `-pgp-key` `(keybase or pgp)`- Path to a file on disk containing a binary or
  base64-encoded public GPG key. This can also be specified as a Keybase
  username using the format "keybase:<username>". When supplied, the generated
  root token will be encrypted and base64-encoded with the given public key.

- `-status` `(bool: false)` - Print the status of the current attempt without
  providing an unseal key. The default is false.
