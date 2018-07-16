---
layout: "guides"
page_title: "Secure Introduction of Vault Clients - Guides"
sidebar_current: "guides-identity-secure-intro"
description: |-
  This introductory guide walk through the mechanism of Vault clients to
  authenticate with Vault. There are two approaches at a high-level: platform
  integration, and trusted orchestrator.
---

# Secure Introduction of Vault Clients

A _secret_ is something that will elevate the risk if exposed to unauthorized
entities and results in undesired consequences (e.g. unauthorized data access);
therefore, only the ***trusted entities*** should have an access to your
secrets.

If you can securely get the first secret from an originator to a consumer,
all subsequent secrets transmitted between this originator and consumer can be
authenticated with the trust established by the successful distribution and user
of that first secret.  

![Secure Introduction](/assets/images/vault-secure-intro-1.png)

The Vault authentication process verifies the secret consumer's identity and
then generate a **token** to associate with that identity.
[Tokens](/docs/concepts/tokens.html) are the core method for authentication
within Vault which means that the secret consumer must first acquire a valid
token.


## Challenge

How does a secret consumer (an application or machine) prove that it is the
legitimate recipient for a secret so that it can acquire a token?

How can you avoid persisting raw token values during our secure
introduction?  

## Secure Introduction Approach

Vault's auth methods perform authentication of its client and assigning a set of
policies which defines the permitted operations for the client.

![Auth Method](/assets/images/vault-auth-method.png)

There are two basic patterns to securely authenticate a secret consumer:
[platform integration](#platform-integration) and [trusted
orchestrator](#trusted-orchestrator).  


### Platform Integration

In the **Platform Integration** model, Vault trusts the underlying platform
(e.g. AWS, Azure, GCP) which assigns an identifier to its cloud resources (e.g.
an IAM token, instance ID, JWT). The Vault client (secret consumer)
authenticates with Vault using its platform provided identifier. Once its
identity was successfully validated against the platform, Vault returns an
initial token to the client with a set of configured policies attached.

![Platform Integration](/assets/images/vault-secure-intro-2.png)

**Use Case**

When the client app is running on a VM hosted on a supported cloud platform, you
can leverage the corresponding auth method to authenticate with Vault.

**Reference Materials:**

- [AWS Auth Method](/docs/auth/aws.html)
- [Azure Auth Method](/docs/auth/azure.html)
- [GCP Auth Method](/docs/auth/azure.html)

### Trusted Orchestrator

In the **Trusted Orchestrator** model, you have an _orchestrator_ which is
already authenticated against Vault with privileged permissions. The
orchestrator launches new applications and inject a mechanism they can use to
authenticate (e.g. AppRole, PKI cert, token, etc) with Vault.

![Trusted Orchestrator](/assets/images/vault-secure-intro-3.png)

**Use Case**

When you are using an orchestrator tool such as Chef to launch applications,
this model can be applied regardless of where the applications are running.

**Reference Materials:**

- [AppRole Auth Method](/docs/auth/approle.html)
  - [AppRole Pull Authentication](/guides/identity/authentication.html)
  - [AppRole with Terraform and Chef Demo](/guides/identity/approle-trusted-entities.html)
- [TLS Certificates Auth Method](/docs/auth/cert.html)
- [Token Auth Method](/docs/auth/token.html)
  - [Cubbyhole Response Wrapping](/guides/secret-mgmt/cubbyhole.html)



## Next steps

Read the reference materials listed for secure introduction model best suited
for your use case.
