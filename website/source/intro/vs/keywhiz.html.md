---
layout: "intro"
page_title: "Vault vs. Keywhiz"
sidebar_current: "vs-other-keywhiz"
description: |-
  Comparison between Vault and Keywhiz.
---

# Vault vs. Keywhiz

Keywhiz is a secret management solution built by Square. Keywhiz
has a client/server architecture. Clients of Keywhiz make use of
a FUSE filesystem to expose secrets as files on disk, and use Unix
file permissions for access control. Underneath, the Keywhiz clients
use mutual TLS (mTLS) to authenticate with a Keywhiz server, which
serves secrets.

Vault similarly is designed as a comprehensive secret management
solution. The client interaction with Vault is much more flexible,
both for authentication and usage of secrets. Vault supports [mTLS
authentication](/docs/auth/cert.html) along with
many [other mechanisms](/docs/auth/index.html).
The goal being to make it easy to authenticate as a machine for programtic
access and as a human for operator usage.

Vault exposes secrets via an API and not over a FUSE filesystem. The
[ACL system](/docs/concepts/policies.html) is used
to protect secrets and gate access, and depends on server side enforcement
instead of Unix permissions on the clients. All auditing is also done
server side using [audit backends](/docs/audit/index.html).

Keywhiz focuses on storage and distribution of secrets and decouples
rotation, and expects external systems to be used for periodic key rotation.
Vault instead supports dynamic secrets, generating credentials on-demand for
fine-grained security controls, auditing, and non-repudiation. Key rotation
is a first class concern for Vault, so that no external system needs to be used.

Lastly Vault forces a mandatory lease contract with clients. All secrets read
from Vault have an associated lease which enables operators to audit key usage,
perform key rolling, and ensure automatic revocation. Vault provides multiple
revocation mechansims to give operators a clear "break glass" procedure after
a potential compromise.

