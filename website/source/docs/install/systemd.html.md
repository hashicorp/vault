---
layout: "docs"
page_title: "Configure SystemD"
sidebar_current: "docs-install-systemd"
description: |-
  Running Vault with SystemD is cool and fun
---

# SystemD

systemd provides a number of security options and features to allow the more secure running of the Vault process.

Note: Be sure to read the systemd documentation carefully before changing options. This configuration is a guide only, and not guarenteed to match your requirements perfectly.

## An example SystemD config

```
[Unit]
Description=A Tool for Managing Secrets
Documentation=https://vaultproject.io/docs/
After=network.target
ConditionFileNotEmpty=/etc/vault/config.json

[Service]
User=vault
Group=vault
PrivateDevices=yes
PrivateTmp=yes
ProtectSystem=full
ProtectHome=read-only
SecureBits=keep-caps
Capabilities=CAP_IPC_LOCK+ep
CapabilityBoundingSet=CAP_SYSLOG CAP_IPC_LOCK
NoNewPrivileges=yes
ExecStart=/usr/local/bin/vault server -config=/etc/vault/config.json
KillSignal=SIGINT
TimeoutStopSec=30s
Restart=on-failure
StartLimitInterval=60s
StartLimitBurst=3

[Install]
WantedBy=multi-user.target
```

## `User` and `Group`

Vault is designed to run as an unprivileged user, and there is no reason to run Vault with root or Administrator privileges, which can expose the Vault process memory and allow access to Vault encryption keys. Running Vault as a regular user reduces its privilege.

## `ConditionFileNotEmpty`

This should be the path to your Vault configuration file, so to prevent Vault trying to run if the config file exists but is empty.

## `NoNewPrivileges`

This is the simplest, most effective way to ensure that a process and its children can never elevate privileges again, which is important for an secure application like Vault. When set to true, it ensures that the vault service process and all its children can never gain new privileges.

## `Capabilities` and `CapabilityBoundingSet`

Vault 0.8.0+ requires the `CAP_IPC_LOCK` permission, which it uses to lock memory via the `mlock` facility. This prevents confidential material being leaked into the swap sections on disk accidentally.

## `ProtectSystem`, `PrivateDevices` and `PrivateTmp`

Setting `ProtectSystem` to true allows the Vault services to choose to mount some filesystems read-only for their processes. Having it set to "true" will create a new mount namespace and mount the /usr and /boot directories read-only in it.

Combined with `PrivateDevices`, which when set to true sets up a new /dev mount for the executed processes and only adds API pseudo devices such as /dev/null, /dev/zero or /dev/random (as well as the pseudo TTY subsystem) to it, but no physical devices such as /dev/sda, system memory /dev/mem, system ports /dev/port and others.

When `PrivateTmp` is true, it will set up a new file system namespace for the executed processes and mounts private /tmp and /var/tmp directories inside it that is not shared by processes outside of the namespace. This is useful to secure access to temporary files of the process, but makes sharing between processes via /tmp or /var/tmp impossible. If this is enabled, all temporary files created by a service in these directories will be removed after the service is stopped.

The combination allows a way to securely turn off physical device access by the executed process.
