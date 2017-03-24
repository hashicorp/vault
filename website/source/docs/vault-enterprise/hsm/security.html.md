---
layout: "docs"
page_title: "Vault Enterprise HSM Security Details"
sidebar_current: "docs-vault-enterprise-hsm-security"
description: |-
  Recommendations to ensure the security of a Vault Enterprise HSM deployment.

---

# Vault Enterprise HSM Security Details

This page provides information to help ensure that a Vault HSM deployment is
performed as securely as possible.

## PKCS#11 Authentication

PKCS#11 authentication occurs via a slot number and PIN. In practice, because
the PIN is not required to be numeric (and some HSMs require more complex
PINs), this behaves like a username and password.

Like a username and password, these values should be protected. If they are
stored in Vault's configuration file, read access to the file should be tightly
controlled to appropriate users. (Vault's configuration file should always have
tight write controls.) Rather than storing these values into Vault's
configuration file, they can also be supplied via the environment; see the
[Configuration](/docs/vault-enterprise/hsm/configuration.html) page for more details.

The attack surface of stolen PKCS#11 credentials depends highly on the
individual HSM, but generally speaking, it should be assumed that if an
attacker can see these credentials and has access to a machine on which Vault
is running, the attacker will be able to access the HSM key protecting Vault's
master key. Therefore, it is extremely important that access to the machine on
which Vault is running is also tightly controlled.

## Recovery Key Shares Protection

Recovery key shares should be protected in the same way as your organization
would protect key shares for the cryptographic barrier. As a quorum of recovery
key shares can be used with the `generate-root` feature to generate a new root
token, and root tokens can do anything within Vault, PGP encryption should
always be used to protect the returned recovery key shares and the recovery
share holders should be highly trusted individuals.
