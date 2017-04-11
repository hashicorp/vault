---
layout: "docs"
page_title: "Security Model"
sidebar_current: "docs-internals-security"
description: |-
  Learn about the security model of Vault.
---

# Security Model

Due to the nature of Vault and the confidentiality of data it is managing,
the Vault security model is very critical. The overall goal of Vault's security
model is to provide [confidentiality, integrity, availability, accountability,
authentication](https://en.wikipedia.org/wiki/Information_security).

This means that data at rest and in transit must be secure from eavesdropping
or tampering. Clients must be appropriately authenticated and authorized
to access data or modify policy. All interactions must be auditable and traced
uniquely back to the origin entity. The system must be robust against intentional
attempts to bypass any of its access controls.

# Threat Model

The following are the various parts of the Vault threat model:

* Eavesdropping on any Vault communication. Client communication with Vault
  should be secure from eavesdropping as well as communication from Vault to
  its storage backend.

* Tampering with data at rest or in transit. Any tampering should be detectable
  and cause Vault to abort processing of the transaction.

* Access to data or controls without authentication or authorization. All requests
  must be proceeded by the applicable security policies.

* Access to data or controls without accountability. If audit logging
  is enabled, requests and responses must be logged before the client receives
  any secret material.

* Confidentiality of stored secrets. Any data that leaves Vault to rest in the
  storage backend must be safe from eavesdropping. In practice, this means all
  data at rest must be encrypted.

* Availability of secret material in the face of failure. Vault supports
  running in a highly available configuration to avoid loss of availability.

The following are not parts of the Vault threat model:

* Protecting against arbitrary control of the storage backend. An attacker
  that can perform arbitrary operations against the storage backend can
  undermine security in any number of ways that are difficult or impossible to protect
  against. As an example, an attacker could delete or corrupt all the contents
  of the storage backend causing total data loss for Vault. The ability to control
  reads would allow an attacker to snapshot in a well-known state and rollback state
  changes if that would be beneficial to them.

* Protecting against the leakage of the existence of secret material. An attacker
  that can read from the storage backend may observe that secret material exists
  and is stored, even if it is kept confidential.

* Protecting against memory analysis of a running Vault. If an attacker is able
  to inspect the memory state of a running Vault instance then the confidentiality
  of data may be compromised.

# External Threat Overview

Given the architecture of Vault, there are 3 distinct systems we are concerned with
for Vault. There is the client, which is speaking to Vault over an API. There is Vault
or the server more accurately, which is providing an API and serving requests. Lastly,
there is the storage backend, which the server is utilizing to read and write data.

There is no mutual trust between the Vault client and server. Clients use
[TLS](https://en.wikipedia.org/wiki/Transport_Layer_Security) to verify the identity
of the server and to establish a secure communication channel. Servers require that
a client provides a client token for every request which is used to identify the client.
A client that does not provide their token is only permitted to make login requests.

The storage backends used by Vault are also untrusted by design. Vault uses a security
barrier for all requests made to the backend. The security barrier automatically encrypts
all data leaving Vault using the [Advanced Encryption Standard (AES)](https://en.wikipedia.org/wiki/Advanced_Encryption_Standard)
cipher in the [Galois Counter Mode (GCM)](https://en.wikipedia.org/wiki/Galois/Counter_Mode).
The nonce is randomly generated for every encrypted object. When data is read from the
security barrier the GCM authentication tag is verified during the decryption process to detect
any tampering.

Depending on the backend used, Vault may communicate with the backend over TLS
to provide an added layer of security. In some cases, such as a file backend this
is not applicable. Because storage backends are untrusted, an eavesdropper would
only gain access to encrypted data even if communication with the backend was intercepted.

# Internal Threat Overview

Within the Vault system, a critical security concern is an attacker attempting
to gain access to secret material they are not authorized to. This is an internal
threat if the attacker is already permitted some level of access to Vault and is
able to authenticate.

When a client first authenticates with Vault, a credential backend is used to
verify the identity of the client and to return a list of associated ACL policies.
This association is configured by operators of Vault ahead of time. For example,
GitHub users in the "engineering" team may be mapped to the "engineering" and "ops"
Vault policies. Vault then generates a client token which is a randomly generated
UUID and maps it to the policy list. This client token is then returned to the client.

On each request a client provides this token. Vault then uses it to check that the token
is valid and has not been revoked or expired, and generates an ACL based on the associated
policies. Vault uses a strict default deny or whitelist enforcement. This means unless
an associated policy allows for a given action, it will be denied. Each policy specifies
a level of access granted to a path in Vault. When the policies are merged (if multiple
policies are associated with a client), the highest access level permitted is used.
For example, if the "engineering" policy permits read/update access to the "eng/" path,
and the "ops" policy permits read access to the "ops/" path, then the user gets the
union of those. Policy is matched using the most specific defined policy, which may be
an exact match or the longest-prefix match glob pattern.

Certain operations are only permitted by "root" users, which is a distinguished
policy built into Vault. This is similar to the concept of a root user on a Unix system
or an Administrator on Windows. Although clients could be provided with root tokens
or associated with the root policy, instead Vault supports the notion of "sudo" privilege.
As part of a policy, users may be granted "sudo" privileges to certain paths, so that
they can still perform security sensitive operations without being granted global
root access to Vault.

Lastly, Vault supports using a [Two-man rule](https://en.wikipedia.org/wiki/Two-man_rule) for
unsealing using [Shamir's Secret Sharing technique](https://en.wikipedia.org/wiki/Shamir's_Secret_Sharing).
When Vault is started, it starts in an _sealed_ state. This means that the encryption key
needed to read and write from the storage backend is not yet known. The process of unsealing
requires providing the master key so that the encryption key can be retrieved. The risk of distributing
the master key is that a single malicious actor with access to it can decrypt the entire
Vault. Instead, Shamir's technique allows us to split the master key into multiple shares or parts.
The number of shares and the threshold needed is configurable, but by default Vault generates
5 shares, any 3 of which must be provided to reconstruct the master key.

By using a secret sharing technique, we avoid the need to place absolute trust in the holder
of the master key, and avoid storing the master key at all. The master key is only
retrievable by reconstructing the shares. The shares are not useful for making any requests
to Vault, and can only be used for unsealing. Once unsealed the standard ACL mechanisms
are used for all requests.

To make an analogy, a bank puts security deposit boxes inside of a vault.
Each security deposit box has a key, while the vault door has both a combination and a key.
The vault is encased in steel and concrete so that the door is the only practical entrance.
The analogy to Vault, is that the cryptosystem is the steel and concrete protecting the data.
While you could tunnel through the concrete or brute force the encryption keys, it would be
prohibitively time consuming. Opening the bank vault requires two-factors: the key and the combination.
Similarly, Vault requires multiple shares be provided to reconstruct the master key.
Once unsealed, each security deposit boxes still requires the owner provide a key, and similarly
the Vault ACL system protects all the secrets stored.
