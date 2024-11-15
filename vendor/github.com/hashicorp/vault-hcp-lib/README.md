# Vault Library: HCP Vault Library

The HCP Vault library is a standalone backend library for use with [Hashicorp
Vault](https://www.github.com/hashicorp/vault).

Please note: We take Vault's security and our users' trust very seriously. If
you believe you have found a security issue in Vault, please responsibly
disclose by contacting us at [security@hashicorp.com](mailto:security@hashicorp.com).

## Quick Links

- [Vault Website](https://developer.hashicorp.com/vault)
- [Vault Project GitHub](https://www.github.com/hashicorp/vault)

## Getting Started

This is a Vault Library and is meant to work with Vault. This guide assumes you have already installed
Vault and have a basic understanding of how Vault works.

Otherwise, first read this guide on how to [get started with
Vault](https://developer.hashicorp.com/vault/tutorials/getting-started/getting-started-install).

## Usage

The HCP Vault library is built into the Vault binary and accessed through the Vault CLI client.

To connect and authenticate to your HCP Vault, use `vault hcp connect`:

```sh
$ vault hcp connect
```

The Vault CLI authenticates users and machines to HCP Vault using a provided credential or interactively with an HCP token generated through browser login. On a successful authentication, the CLI caches the returned HCP token and current HCP Vault address 

By default, the Vault CLI uses interactive authentication and directs users to the HCP login page.

Non-interactive authentication requires service principal credentials
previously generated through the HCP portal. The provided credential
must have sufficient permission to access the organization, project, and
 HCP Vault cluster.
 
 For example, to connect with a client ID and secret:

```sh
$ vault hcp connect -client-id=client-id-value -secret-id=secret-id-value
```

You can also target specific organizations, projects, and clusters by providing the relevant identification:

```sh
$ vault hcp connect           \
  -client-id=client-id-value  \
  -secret-id=secret-id-value  \
  -organization-id=org-UUID   \
  -project-id=proj-UUID       \
  -cluster-id=cluster-name
```

To clean HCP credentials from the cache use the `disconnect` subcommand:

```sh
$ vault hcp disconnect
```

For more information about supported subcommands and options, refer to the [Vault CLI documentation](https://add-documentation-here).

## How to contribute

Thanks for considering contributing to this project. Unfortunately, HashiCorp does not currently accept new contributions for this project.

## License

This code is released under the Mozilla Public License 2.0. Please see [LICENSE](https://github.com/hashicorp/terraform-aws-hcp-consul/blob/main/LICENSE) for more details.
