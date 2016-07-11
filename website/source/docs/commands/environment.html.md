---
layout: "docs"
page_title: "Environment"
sidebar_current: "docs-commands-environment"
description: |-
  Vault's behavior can be modified by certain environment variables.
---

# Environment variables

The Vault CLI will read the following environment variables to set
behavioral defaults. These can be overridden in all cases using
command-line arguments; see the command-line help for details.

The following table describes them:

<table>
  <tr>
    <th>Variable name</th>
    <th>Value</th>
  </tr>
  <tr>
    <td><tt>VAULT_TOKEN</tt></td>
    <td>The Vault authentication token.  If not specified, the token located in <tt>$HOME/.vault-token</tt> will be used if it exists.</td>
  </tr>
  <tr>
    <td><tt>VAULT_ADDR</tt></td>
    <td>The address of the Vault server.</td>
  </tr>
  <tr>
    <td><tt>VAULT_ADVERTISE_ADDR</tt></td>
    <td>The advertised address of the server to use for client request forwarding when running in High Availability mode.</td>
  <tr>
    <td><tt>VAULT_CACERT</tt></td>
    <td>Path to a PEM-encoded CA cert file to use to verify the Vault server SSL certificate.</td>
  </tr>
  <tr>
    <td><tt>VAULT_CAPATH</tt></td>
    <td>Path to a directory of PEM-encoded CA cert files to verify the Vault server SSL certificate.  If <tt>VAULT_CACERT</tt> is specified, its value will take precedence.</td>
  </tr>
  <tr>
    <td><tt>VAULT_CLIENT_CERT</tt></td>
    <td>Path to a PEM-encoded client certificate for TLS authentication to the Vault server.</td>
  </tr>
  <tr>
    <td><tt>VAULT_CLIENT_KEY</tt></td>
    <td>Path to an unencrypted PEM-encoded private key matching the client certificate.</td>
  </tr>
  <tr>
    <td><tt>VAULT_MAX_RETRIES</tt></td>
    <td>The maximum number of retries when a `5xx` error code is encountered. Default is `2`, for three total tries; set to `0` or less to disable retrying.</td>
  </tr>
  <tr>
    <td><tt>VAULT_SKIP_VERIFY</tt></td>
    <td>If set, do not verify Vault's presented certificate before communicating with it.  Setting this variable is not recommended except during testing.</td>
  </tr>
  <tr>
    <td><tt>VAULT_TLS_SERVER_NAME</tt></td>
    <td>If set, use the given name as the SNI host when connecting via TLS.</td>
  </tr>
</table>
