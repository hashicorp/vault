---
layout: "docs"
page_title: "Vault Secure Introduction Client Configuration"
sidebar_current: "docs-vault-enterprise-vsi-configuration"
description: |-
  Configuration details for Vault Secure Introduction Client.

---

# Vault Secure Introduction Client Configuration

The Vault Secure Introduction client has a flexible configuration system that
allows combining settings from CLI flags, environment variables, and a
configuration file, in that order of preference.

Generally speaking, values with a source at a higher order of preference
override those from a source with a lower order of preference. The main
difference is in specification of servers, which is detailed further in this
document.

The available directives from the configuration file, as well as their
environment variable and CLI flag counterparts, follow.

## The Configuration File

The configuration file is in [HCL](https://github.com/hashicorp/hcl) format and
contains top-level directives as well as directive blocks. An example:

```javascript
environment "aws" {
}

vault {
  address = "https://vault.service.consul:8200"
  mount_path = "auth/aws"
}

serve "file" {
  path = "/ramdisk/vault-token"
}
```

The configuration file can be specified either with the `config` CLI flag
or by simply providing the path to the configuration file as the command's
argument.

All directives take string arguments unless explicitly specified otherwise.

## Top-Level Directives

 * `environment`: A block with information about he environment under which the
   client is running. If not specified, the client will attempt to
   automatically discover its environment. The type may also be specified by
   the `VAULT_SI_ENVIRONMENT` environment variable and the `environment` CLI
   flag. Additional configuration key/value pairs can be passed in via the
   `envconfig` CLI flag, which can be specified multiple times.
 * `nonce_path`: If the client should save and load its nonce, the path where
   the nonce should be stored. May also be specified by the
   `VAULT_SI_NONCE_PATH` environment variable and the `nonce-path` CLI flag.
   _This is a security-sensitive directive._
 * `vault`: A block with Vault server information, detailed below.
 * `serve`: One or more blocks containing serving information, detailed below.

## Environment Block Directives

In the configuration file, the type is specified as the block's key, with
optional additional string values specified inside:

```hcl
environment "aws" {
   role = "prod"
}
```

The behavior with respect to these additional string values is
environment-specific. Currently, all environments simply round-trip any given
values to the Vault login endpoint. In the example above using the AWS
environment, the final set of values given to the login endpoint would be a
`pkcs7` key coming from the environment, as well as a `role` key with value
`prod` coming from the extra environment configuration.

## Vault Block Directives

 * `address`: The address of the vault server, including scheme. May also be
   specified by the `VAULT_ADDR` environment variable and the `address` CLI
   flag.
 * `mount_path`: The mount path of the authentication backend in Vault. If not
   set, defaults to a value specific to the running environment (e.g. for AWS,
   it will default to `auth/aws`. May also be specified by the
   `VAULT_SI_MOUNT_PATH` environment variable and the `mount-path` CLI flag.
 * `tls_skip_verify`: A boolean indicating whether to skip verification of the
   certificate provided by the Vault server. May also be specified by the
   `VAULT_SKIP_VERIFY` environment variable or the `tls-skip-verify` CLI flag.
   _This is a security-sensitive directive._
 * `ca_cert`: A file containing a PEM-encoded X.509 CA certificate to use in
   the validation chain of the Vault server's TLS certificate. May also be
   specified by the `VAULT_CACERT` environment variable or the `ca-cert` CLI
   flag.
 * `ca_path`: A directory containing PEM-encoded X.509 CA certificates to use
   in the validation chain of the Vault server's TLS certificate. May also be
   specified by the `VAULT_CAPATH` environment variable or the `ca-path` CLI
   flag.

## Serve Block Directives

In the configuration file, serve blocks can be one of two types:

 * Named serve blocks specify a name (`serve "myserver" {...`). The type of server
   must be specified by the `type` directive within the block.
 * Anonymous serve blocks, rather than specify a name, specify the type of
   server (`serve "file" {...`).

On the CLI, serve blocks are specified in one of two formats:

 * `<type>:<key>=<value>`: specifies a key/value configuration directive for an
   anonymous server with the given type.
 * `<type>:<name>:<key>=<value>`: specifies a key/value configuration directive
   for a named server with the given type.

The merging rules for CLI and configuration are as follows:

 * Each anonymous serve block in the configuration file stands alone, using
   only the directives contained in the block.
 * Each type of anonymous server specified in CLI flags CLI stands alone, with
   the key/value configuration directives merged per-type. As a result, there
   can only be one anonymous server per type specified in CLI flags. These are
   not merged with any anonymous server specified in the configuration file.
 * Key/value configuration directives for named serve blocks are merged between
   the CLI and configuration file.

### File Type Serve Block Directives

 * `path`: The path to a file on disk where the token should be written. To
   avoid any possible issues during writing, the token will first be written to
   a temporary file in the same directory, then atomically renamed to the given
   path. The token will always be written with permissions `0640`; directory
   permissions for this location should ensure access only to appropriate
   readers.
