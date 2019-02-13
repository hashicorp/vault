---
layout: "guides"
page_title: "Vault Deployment Guide - Guides"
sidebar_current: "guides-operations-deployment-guide"
description: |-
  This deployment guide covers the steps required to install and
  configure a single HashiCorp Vault cluster as defined in the
  Vault Reference Architecture.
ea_version: 1.0
---

# Vault Deployment Guide

This deployment guide covers the steps required to install and configure a single HashiCorp Vault cluster as defined in the [Vault Reference Architecture](/guides/operations/reference-architecture.html).

Below are instructions for installing and configuring Vault on Linux hosts running the systemd system and service manager.

## Reference Material

This deployment guide is designed to work in combination with the [Vault Reference Architecture](/guides/operations/reference-architecture.html). Although not a strict requirement to follow the Vault Reference Architecture, please ensure you are familiar with the overall architecture design; for example installing Vault on multiple physical or virtual (with correct anti-affinity) hosts for high-availability and using Consul for the HA and storage backend.

During the installation of Vault you should also review and apply the recommendations provided in the [Vault Production Hardening](/guides/operations/production.html) guide.

## Overview

To provide a highly-available single cluster architecture, we recommend Vault be deployed to more than one host, as shown in the [Vault Reference Architecture](/guides/operations/reference-architecture.html), and connected to a Consul cluster for persistent data storage.

![Reference Diagram](/img/vault-ref-arch-2-02305ae7.png)

The below setup steps should be completed on all Vault hosts.

- [Download Vault](#download-vault)
- [Install Vault](#install-vault)
- [Configure systemd](#configure-systemd)
- [Configure Consul](#configure-consul)
- [Configure Vault](#configure-vault)
- [Start Vault](#start-vault)

## Download Vault

Precompiled Vault binaries are available for download at [https://releases.hashicorp.com/vault/](https://releases.hashicorp.com/vault/) and Vault Enterprise binaries are available for download by following the instructions made available to HashiCorp Vault customers.

You should perform checksum verification of the zip packages using the SHA256SUMS and SHA256SUMS.sig files available for the specific release version. HashiCorp provides [a guide on checksum verification](https://www.hashicorp.com/security.html) for precompiled binaries.

```text
VAULT_VERSION="0.10.3"
curl --silent --remote-name https://releases.hashicorp.com/vault/${VAULT_VERSION}/vault_${VAULT_VERSION}_linux_amd64.zip
curl --silent --remote-name https://releases.hashicorp.com/vault/${VAULT_VERSION}/vault_${VAULT_VERSION}_SHA256SUMS
curl --silent --remote-name https://releases.hashicorp.com/vault/${VAULT_VERSION}/vault_${VAULT_VERSION}_SHA256SUMS.sig
```

## Install Vault

Unzip the downloaded package and move the `vault` binary to `/usr/local/bin/`. Check `vault` is available on the system path.

```text
unzip vault_${VAULT_VERSION}_linux_amd64.zip
sudo chown root:root vault
sudo mv vault /usr/local/bin/
vault --version
```

The `vault` command features opt-in autocompletion for flags, subcommands, and arguments (where supported). Enable autocompletion.

```text
vault -autocomplete-install
complete -C /usr/local/bin/vault vault
```

Give Vault the ability to use the mlock syscall without running the process as root. The mlock syscall prevents memory from being swapped to disk.

```text
sudo setcap cap_ipc_lock=+ep /usr/local/bin/vault
```

Create a unique, non-privileged system user to run Vault.

```text
sudo useradd --system --home /etc/vault.d --shell /bin/false vault
```

## Configure systemd

Systemd uses [documented sane defaults](https://www.freedesktop.org/software/systemd/man/systemd.directives.html) so only non-default values must be set in the configuration file.

Create a Vault service file at /etc/systemd/system/vault.service.

```text
sudo touch /etc/systemd/system/vault.service
```

Add the below configuration to the Vault service file:

```text
[Unit]
Description="HashiCorp Vault - A tool for managing secrets"
Documentation=https://www.vaultproject.io/docs/
Requires=network-online.target
After=network-online.target
ConditionFileNotEmpty=/etc/vault.d/vault.hcl

[Service]
User=vault
Group=vault
ProtectSystem=full
ProtectHome=read-only
PrivateTmp=yes
PrivateDevices=yes
SecureBits=keep-caps
AmbientCapabilities=CAP_IPC_LOCK
Capabilities=CAP_IPC_LOCK+ep
CapabilityBoundingSet=CAP_SYSLOG CAP_IPC_LOCK
NoNewPrivileges=yes
ExecStart=/usr/local/bin/vault server -config=/etc/vault.d/vault.hcl
ExecReload=/bin/kill --signal HUP $MAINPID
KillMode=process
KillSignal=SIGINT
Restart=on-failure
RestartSec=5
TimeoutStopSec=30
StartLimitIntervalSec=60
StartLimitBurst=3

[Install]
WantedBy=multi-user.target
```

The following parameters are set for the `[Unit]` stanza:

- [`Description`](https://www.freedesktop.org/software/systemd/man/systemd.unit.html#Description=) - Free-form string describing the vault service
- [`Documentation`](https://www.freedesktop.org/software/systemd/man/systemd.unit.html#Documentation=) - Link to the vault documentation
- [`Requires`](https://www.freedesktop.org/software/systemd/man/systemd.unit.html#Requires=) - Configure a requirement dependency on the network service
- [`After`](https://www.freedesktop.org/software/systemd/man/systemd.unit.html#Before=) - Configure an ordering dependency on the network service being started before the vault service
- [`ConditionFileNotEmpty`](https://www.freedesktop.org/software/systemd/man/systemd.unit.html#ConditionArchitecture=) - Check for a non-zero sized configuration file before vault is started

The following parameters are set for the `[Service]` stanza:

- [`User`, `Group`](https://www.freedesktop.org/software/systemd/man/systemd.exec.html#User=) - Run vault as the vault user
- [`ProtectSystem`, `ProtectHome`, `PrivateTmp`, `PrivateDevices`](https://www.freedesktop.org/software/systemd/man/systemd.exec.html#Sandboxing) - Sandboxing settings to improve the security of the host by restricting vault privileges and access
- [`SecureBits`, `Capabilities`, `CapabilityBoundingSet`, `AmbientCapabilities`](http://man7.org/linux/man-pages/man7/capabilities.7.html) - Configure the capabilities of the vault process
- [`NoNewPrivileges`](https://www.freedesktop.org/software/systemd/man/systemd.exec.html#NoNewPrivileges=) - Prevent vault and any child process from gaining new privileges
- [`ExecStart`](https://www.freedesktop.org/software/systemd/man/systemd.service.html#ExecStart=) - Start vault with the `server` argument and path to the configuration file
- [`ExecReload`](https://www.freedesktop.org/software/systemd/man/systemd.service.html#ExecReload=) - Send vault a HUP signal to trigger a configuration reload in vault
- [`KillMode`](https://www.freedesktop.org/software/systemd/man/systemd.kill.html#KillMode=) - Treat vault as a single process
- [`KillSignal`](https://www.freedesktop.org/software/systemd/man/systemd.kill.html#KillSignal=) - Send SIGINT signal when shutting down vault
- [`Restart`](https://www.freedesktop.org/software/systemd/man/systemd.service.html#RestartSec=) - Restart vault ([in a sealed state](/docs/concepts/seal.html)) unless it returned a clean exit code
- [`RestartSec`](https://www.freedesktop.org/software/systemd/man/systemd.service.html#RestartSec=) - Restart vault after 5 seconds of it being considered 'failed'
- [`TimeoutStopSec`](https://www.freedesktop.org/software/systemd/man/systemd.service.html#TimeoutStopSec=) - Wait 30 seconds for a clean stop before sending a SIGKILL signal
- [`StartLimitIntervalSec`, `StartLimitBurst`](https://www.freedesktop.org/software/systemd/man/systemd.unit.html#StartLimitIntervalSec=interval) - Limit vault to three start attempts in 60 seconds

The following parameters are set for the `[Install]` stanza:

- [`WantedBy`](https://www.freedesktop.org/software/systemd/man/systemd.unit.html#WantedBy=) - Creates a weak dependency on vault being started by the multi-user run level

## Configure Consul

When using Consul as the storage backend for Vault, we recommend using Consul's [ACL system](https://www.consul.io/docs/guides/acl.html) to restrict access to the path where Vault stores data. This access restriction is an added security measure in addition to the [encryption Vault uses to protect data](/docs/internals/architecture.html) written to the storage backend.

The Consul website provides documentation on [bootstrapping the ACL system](https://www.consul.io/docs/guides/acl.html#bootstrapping-acls), generating a management token and using that token to add some initial tokens for Consul agents, UI access etc. You should complete the bootstrapping section of the Consul documentation before continuing with this guide.

Vault requires a Consul token with specific policy to limit the requests Vault can make to Consul endpoints.

On a host running a Consul agent, and using a Consul management token, create a Consul client token with specific policy for Vault:

```text
CONSUL_TOKEN="6609e426-1aeb-4b0d-c302-3a7568fbc1f9"
curl \
    --request PUT \
    --header "X-Consul-Token: ${CONSUL_TOKEN}" \
    --data \
'{
  "Name": "Vault Token",
  "Type": "client",
  "Rules": "node \"\" { policy = \"write\" } service \"vault\" { policy = \"write\" } agent \"\" { policy = \"write\" }  key \"vault\" { policy = \"write\" } session \"\" { policy = \"write\" } "
}' http://127.0.0.1:8500/v1/acl/create
```

The response includes the value you will use as the `token` parameter value in Vault's storage stanza configuration. An example response:

```json
{"ID":"fe3b8d40-0ee0-8783-6cc2-ab1aa9bb16c1"}
```

## Configure Vault

Vault uses [documented sane defaults](/docs/configuration) so only non-default values must be set in the configuration file.

Create a configuration file at /etc/vault.d/vault.hcl:

```text
sudo mkdir --parents /etc/vault.d
sudo touch /etc/vault.d/vault.hcl
sudo chown --recursive vault:vault /etc/vault.d
sudo chmod 640 /etc/vault.d/vault.hcl
```

### Listener stanza

The `listener` stanza configures the addresses and ports on which Vault will respond to requests.

Add the below configuration to the Vault configuration file:

```hcl
listener "tcp" {
  address       = "0.0.0.0:8200"
  tls_cert_file = "/path/to/fullchain.pem"
  tls_key_file  = "/path/to/privkey.pem"
}
```

The following parameters are set for the `tcp` listener stanza:

- [`address`](/docs/configuration/listener/tcp.html#address) `(string: "127.0.0.1:8200")` - Changing from the loopback address to allow external access to the Vault UI
- [`tls_cert_file`](/docs/configuration/listener/tcp.html#tls_cert_file) `(string: <required-if-enabled>, reloads-on-SIGHUP)` - Must be set when using TLS
- [`tls_key_file`](/docs/configuration/listener/tcp.html#tls_key_file) `(string: <required-if-enabled>, reloads-on-SIGHUP)` - Must be set when using TLS

[More information about tcp listener configuration](/docs/configuration/listener/tcp.html).

~> Vault should always be configured to use TLS to provide secure communication between clients and the Vault cluster. This requires a certificate file and key file be installed on each Linux host running Vault. The certificate file and key file must have permissions allowing the vault user/group to read them.

### Seal stanza

This is an __ENTERPRISE__ feature.

If you are deploying [Vault Enterprise](https://www.hashicorp.com/products/vault), you can include `seal` stanza configuration to specify the seal type to use for additional data protection, such as using HSM or Cloud KMS solutions to encrypt and decrypt the Vault master key. This stanza is optional, and if this is not configured, Vault will use the Shamir algorithm to cryptographically split the master key.

If you are deploying Vault Enterprise, you should review the [seal configuration section](/docs/configuration/seal/index.html) of our documentation.

An example PKCS #11 compatible HSM example is:

``` hcl
seal "pkcs11" {
  lib            = "/usr/vault/lib/libCryptoki2_64.so"
  slot           = "0"
  pin            = "AAAA-BBBB-CCCC-DDDD"
  key_label      = "vault-hsm-key"
  hmac_key_label = "vault-hsm-hmac-key"
}
```

### Storage stanza

The `storage` stanza configures the storage backend, which represents the location for the durable storage of Vault's data.

Add the below configuration to the Vault configuration file:

```hcl
storage "consul" {
  token = "{{ consul_token }}"
}
```

The following parameters are set for the `consul` storage stanza:

- [`token`](/docs/configuration/storage/consul.html#token) `(string: "")` - Specify the Consul ACL token with permission to read and write from `/vault` in Consul's key-value store

[More information about consul storage configuration](/docs/configuration/storage/consul.html).

~> Vault should always be configured to use a Consul token with a restrictive ACL policy to read and write from `/vault` in Consul's key-value store. This follows the principal of least privilege, ensuring Vault is unable to access Consul key-value data stored outside of the `/vault` path.

### Telemetry stanza

The `telemetry` stanza specifies various configurations for Vault to publish metrics to upstream systems.

If you decide to configure Vault to publish telemtery data, you should review the [telemetry configuration section](/docs/configuration/telemetry.html) of our documentation.

### High Availability Parameters

The `api_addr` parameter configures the API address used in high availability scenarios, when client redirection is used instead of request forwarding. Client redirection is the fallback method used when request forwarding is turned off or there is an error performing the forwarding. As such, a redirect address is always required for all HA setups.

This parameter value defaults to the `address` specified in the `listener` stanza, but Vault will log a `[WARN]` message if it is not explicitly configured.

Add the below configuration to the Vault configuration file:

```hcl
api_addr = "{{ full URL to Vault API endpoint }}"
```

[More information about high availability configuration](/docs/configuration/#high-availability-parameters).

### Vault UI

Vault features a web-based user interface, allowing you to easily create, read, update, and delete secrets, authenticate, unseal, and more using a graphical user interface, rather than the CLI or API.

Vault should not be deployed in a public internet facing environment, so enabling the Vault UI is typically of benefit to provide a more familiar experience to administrators who are not as comfortable working on the command line, or who do not have alternative access.

Optionally, add the below configuration to the Vault configuration file to enable the Vault UI:

```hcl
ui = true
```

[More information about configuring the Vault UI](/docs/configuration/ui/index.html).

## Start Vault

Enable and start Vault using the systemctl command responsible for controlling systemd managed services. Check the status of the vault service using systemctl.

```text
sudo systemctl enable vault
sudo systemctl start vault
sudo systemctl status vault
```

## Next Steps

- Read [Production Hardening](/guides/operations/production.html) to learn best
  practices for a production hardening deployment of Vault.
