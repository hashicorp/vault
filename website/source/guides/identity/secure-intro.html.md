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
therefore, only the ***trusted entities*** should have an access to your secrets.

If you can securely get the first secret from an originator to a consumer,
all subsequent secrets transmitted between this originator and consumer can be
authenticated with the trust established by the successful distribution and user
of that first secret.  

![Secure Introduction](/assets/images/vault-secure-intro-1.png)

Tokens are the core method for authentication within Vault which means that the
secret consumer must first acquire a token.


## Challenge

How does a secret consumer prove that it is the legitimate recipient for a
secret so that it can acquire a token?

How can you avoid persisting raw token values during our secure
introduction?  

## Solution: Secure Introduction Approach

There are two common practices: [platform integration](#platform-integration)
and [trusted orchestrator](#trusted-orchestrator).  


### Platform Integration

![](/assets/images/vault-secure-intro-2.png)

In this model, you have a 3-legged trust model. The client needs to provide
something to Vault to prove its identity. Vault trusts the underlying platform
(e.g. AWS), which provides something to the application (e.g. an IAM token),
which is provided to Vault to complete the chain. For each platform, there is a
slightly different token, but its all the same basic mechanism.




### Trusted Orchestrator

![](/assets/images/vault-approle-workflow2.png)

In this model, you have an orchestrator which is already authenticated against
Vault with a high level of access. The orchestrator is launching new
applications and injecting a mechanism they can use to authenticate (e.g.
AppRole, PKI cert, token, etc).















It would be useful for us to capture these two patterns at a high level and talk about how they can both be used, and link to the various Auth methods that makes sense for given platforms or orchestrators.






















## Next steps

Read the [_AppRole with Terraform and
Chef_](/guides/identity/approle-trusted-entities.html) guide to better
understand the role of trusted entities using Terraform and Chef as an example.

To learn more about response wrapping, go to the [Cubbyhole Response
Wrapping](/guides/secret-mgmt/cubbyhole.html) guide.
