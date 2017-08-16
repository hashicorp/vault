---
layout: "guides"
page_title: "Production Hardening - Guides"
sidebar_current: "guides-production-hardening"
description: |-
  This guide provides guidance on best practices for a production hardened deployment of HashiCorp Vault.
---

# Production Hardening

This guide provides guidance on best practices for a production hardened
deployment of Vault.  The recommendations are based on the [security
model](/docs/internals/security.html) and focus on defense in depth.

~> **Apply When Possible!** This guide is meant to provide guidance for an
_ideal_ deployment of Vault, not to document requirements.  It is entirely
possible to use Vault without applying any of the following recommendations.
These are best practice recommendations that should be applied when possible
and practical.

# Recommendations

* **End-to-End TLS**. Vault should always be used with TLS in production. If
  intermediate load balancers or reverse proxies are used to front Vault, they
  should _not_ terminate TLS. This way traffic is always encrypted in transit
  to Vault and minimizes risks introduced by intermediate layers.

* **Single Tenancy**. Vault should be the only main process running on a
  machine. This reduces the risk that another process running on the same
  machine is compromised and can interact with Vault. Similarly, running on
  bare metal should be preferred to a VM, and a VM preferred to a container.
  This reduces the surface area introduced by additional layers of abstraction
  and other tenants of the hardware. Both VM and container based deployments
  work, but should be avoided when possible to minimize risk.

* **Firewall traffic**. Vault listens on well known ports, use a local firewall
  to restrict all incoming and outgoing traffic to Vault and essential system
  services like NTP. This includes restricting incoming traffic to permitted
  subnets and outgoing traffic to services Vault needs to connect to, such as
  databases.

* **Disable SSH / Remote Desktop**. When running a Vault as a single tenant
  application, users should never access the machine directly. Instead, they
  should access Vault through its API over the network. Use a centralized
  logging and telemetry solution for debugging. Be sure to restrict access to
  logs as need to know.

* **Disable Swap**. Vault encrypts data in transit and at rest, however it must
  still have sensitive data in memory to function. Risk of exposure should be
  minimized by disabling swap to prevent the operating system from paging
  sensitive data to disk. Vault attempts to ["memory lock" to physical memory
  automatically](/docs/configuration/index.html#disable_mlock), but disabling
  swap adds another layer of defense.

* **Don't Run as Root**. Vault is designed to run as an unprivileged user, and
  there is no reason to run Vault with root or Administrator privileges, which
  can expose the Vault process memory and allow access to Vault encryption
  keys. Running Vault as a regular user reduces its privilege. Configuration
  files for Vault should have permissions set to restrict access to only the
  Vault user.

* **Turn Off Core Dumps**. A user or administrator that can force a core dump
  and has access to the resulting file can potentially access Vault encryption
  keys. Preventing core dumps is a platform-specific process; on Linux setting
  the resource limit `RLIMIT_CORE` to `0` disables core dumps. This can be
  performed by process managers and is also exposed by various shells; in Bash
  `ulimit -c 0` will accomplish this.

* **Immutable Upgrades**. Vault relies on an external storage backend for
  persistence, and this decoupling allows the servers running Vault to be
  managed immutably. When upgrading to new versions, new servers with the
  upgraded version of Vault are brought online. They are attached to the same
  shared storage backend and unsealed. Then the old servers are destroyed. This
  reduces the need for remote access and upgrade orchestration which may
  introduce security gaps.

* **Avoid Root Tokens**. Vault provides a root token when it is first
  initialized. This token should be used to setup the system initially,
  particularly setting up authentication backends so that users may
  authenticate. We recommend treating Vault [configuration as
  code](https://www.hashicorp.com/blog/codifying-vault-policies-and-configuration/),
  and using version control to manage policies. Once setup, the root token
  should be revoked to eliminate the risk of exposure. Root tokens can be
  [generated when needed](/guides/generate-root.html), and should be
  revoked as soon as possible.

* **Enable Auditing**. Vault supports several auditing backends. Enabling
  auditing provides a history of all operations performed by Vault and provides
  a forensics trail in the case of misuse or compromise. Audit logs [securely
  hash](/docs/audit/index.html) any sensitive data, but access should still be
  restricted to prevent any unintended disclosures.

* **Upgrade Frequently**. Vault is actively developed, and updating frequently
  is important to incorporate security fixes and any changes in default
  settings such as key lengths or cipher suites. Subscribe to the [Vault
  mailing list](https://groups.google.com/forum/#!forum/vault-tool) and [GitHub
  CHANGELOG](https://github.com/hashicorp/vault/blob/master/CHANGELOG.md) for
  updates.

* **Configure SELinux / AppArmor**. Using additional mechanisms like SELinux
  and AppArmor can help provide additional layers of security when using Vault.
  While Vault can run on many operating systems, we recommend Linux due to the
  various security primitives mentioned here.

* **Restrict Storage Access**. Vault encrypts all data at rest, regardless of
  which storage backend is used. Although the data is encrypted, an [attacker
  with arbitrary control](/docs/internals/security.html) can cause data
  corruption or loss by modifying or deleting keys. Access to the storage
  backend should be restricted to only Vault to avoid unauthorized access or
  operations.
