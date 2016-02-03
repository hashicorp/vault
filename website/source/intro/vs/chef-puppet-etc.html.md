---
layout: "intro"
page_title: "Vault vs. Chef, Puppet, etc."
sidebar_current: "vs-other-chef"
description: |-
  Comparison between Vault and configuration management solutions such as Chef, Puppet, etc.
---

# Vault vs. Chef, Puppet, etc.

A big part of configuring software is setting up secrets: configuring a
web application to talk to a service, configuring the credentials of a
database, etc. Because of this, configuration management systems all face
a problem of safely storing these secrets.

Chef, Puppet, etc. all solve this in a similar way: single-key
encrypted storage. Chef has encrypted data bags, Puppet has encrypted
Hiera, and so on. The encrypted data is always one secret (a password,
a key, etc.) away from being decrypted, and this secret is generally
not well protected since in an elastic environment, every server needs
to somehow get this secret to decrypt the data. Additionally, access to
the encrypted data isn't always logged, so if there is an intrusion, it
isn't clear what data has been accessed and by who.

Vault is not tied to any specific configuration management system. You can
read secrets from configuration management, but you can also use the API
directly to read secrets from applications. This means that configuration
management requires fewer secrets, and in many cases doesn't ever have to
persist them to disk.

Vault encrypts the data onto physical storage and requires multiple
keys to even read it. If an attacker were to gain access to the physical
encrypted storage, it couldn't be read without multiple keys which are generally
distributed to multiple individuals. This is known as _unsealing_, and happens
once whenever Vault starts.

For an unsealed Vault, every interaction is logged in via the audit backends.
Even erroneous requests (invalid access tokens, for example) are logged.
To access any data, an access token is required. This token is usually
associated with an identity coming from a system such as GitHub, LDAP, etc.
This identity is also written to the audit log.

Access tokens can be given fine-grained control over what secrets can be
accessed. It is rare to have a single key that can access all secrets. This
makes it easier to have fine-grained access for consumers of Vault.

For tips on how to integrate Vault using configuration management, please see
[Using HashiCorp's Vault with Chef](https://www.hashicorp.com/blog/using-hashicorp-vault-with-chef.html).
Although this post is about Chef, the principles can be broadly applied to many
of the tools listed here.
