---
layout: "docs"
page_title: "SSH - Secrets Engines"
sidebar_current: "docs-secrets-ssh"
description: |-
  The Vault SSH secrets engine provides secure authentication and authorization
  for access to machines via the SSH protocol. There are multiple modes to the
  Vault SSH secrets engine including signed SSH certificates, dynamic SSH keys,
  and one-time passwords.
---

# SSH Secrets Engine

Name: `ssh`

The Vault SSH secrets engine provides secure authentication and authorization
for access to machines via the SSH protocol. The Vault SSH secrets engine helps
manage access to machine infrastructure, providing several ways to issue SSH
credentials.

The Vault SSH secrets engine supports the following modes. Each mode is
individually documented on its own page.

- [Signed SSH Certificates](/docs/secrets/ssh/signed-ssh-certificates.html)
- [One-time SSH Passwords](/docs/secrets/ssh/one-time-ssh-passwords.html)
- [Dynamic SSH Keys](/docs/secrets/ssh/dynamic-ssh-keys.html) <sup>DEPRECATED</sup>

All guides assume a basic familiarity with the SSH protocol.

## API

The SSH secrets engine has a full HTTP API. Please see the
[SSH secrets engine API](/api/secret/ssh/index.html) for more
details.
