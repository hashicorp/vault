---
layout: "guides"
page_title: "Secure Introduction of Vault Clients - Guides"
sidebar_title: "Secure Introduction of Vault Clients"
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

If you can securely get the first secret from an originator to a consumer, all
subsequent secrets transmitted between this originator and consumer can be
authenticated with the trust established by the successful distribution and user
of that first secret. Getting the first secret to the consumer, is the ***secure
introduction*** challenge.

![Secure Introduction](/img/vault-secure-intro-1.png)

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

![Auth Method](/img/vault-auth-method.png)

There are three basic approaches to securely authenticate a secret consumer:

- [Platform Integration](#platform-integration)
- [Trusted Orchestrator](#trusted-orchestrator)
- [Vault Agent](#vault-agent)


## Platform Integration

In the **Platform Integration** model, Vault trusts the underlying platform
(e.g. AliCloud, AWS, Azure, GCP) which assigns a token or cryptographic identity
(such as IAM token, signed JWT) to virtual machine, container, or serverless
function.

Vault uses the provided identifier to verify the identity of the client by
interacting with the underlying platform. After the client identity is verified,
Vault returns a token to the client that is bound to their identity and policies
that grant access to secrets.

![Platform Integration](/img/vault-secure-intro-2.png)

For example, suppose we have an application running on a virtual machine in AWS
EC2. When that instance is started, an IAM token is provided via the machine
local metadata URL. That IAM token is provided to Vault, as part of the AWS Auth
Method, to login and authenticate the client. Vault uses that token to query the
AWS API and verify the token validity and fetch additional metadata about the
instance (Account ID, VPC ID, AMI, Region, etc). These properties are used to
determine the identity of the client and to distinguish between different roles
(e.g. a Web server versus an API server).

Once validated and assigned to a role, Vault generates a token that is
appropriately scoped and returns it to the client. All future requests from the
client are made with the associated token, allowing Vault to efficiently
authenticate the client and check for proper authorizations when consuming
secrets.

![Vault AWS EC2 Authentication Flow](/img/vault-aws-ec2-auth-flow.png)


### Use Case

When the client app is running on a VM hosted on a supported cloud platform, you
can leverage the corresponding auth method to authenticate with Vault.

### Reference Materials:

- [AWS Auth Method](/docs/auth/aws.html)
- [Azure Auth Method](/docs/auth/azure.html)
- [GCP Auth Method](/docs/auth/gcp.html)


## Trusted Orchestrator

In the **Trusted Orchestrator** model, you have an _orchestrator_ which is
already authenticated against Vault with privileged permissions. The
orchestrator launches new applications and inject a mechanism they can use to
authenticate (e.g. AppRole, PKI cert, token, etc) with Vault.

![Trusted Orchestrator](/img/vault-secure-intro-3.png)

For example, suppose [Terraform](https://www.terraform.io/) is being used as a
trusted orchestrator. This means Terraform already has a Vault token, with
enough capabilities to generate new tokens or create new mechanisms to
authenticate such as an AppRole. Terraform can interact with platforms such as
VMware to provision new virtual machines. VMware does not provide a
cryptographic identity, so a platform integration isn't possible. Instead,
Terraform can provision a new AppRole credential, and SSH into the new machine
to inject the credentials. Terraform is creating the new credential in Vault,
and making that credential available to the new resource. In this way, Terraform
is acting as a trusted orchestrator and extending trust to the new machine. The
new machine, or application running on it, can use the injected credentials to
authenticate against Vault.

![AppRole auth method workflow](/img/vault-secure-intro-4.png)


### Use Case

When you are using an orchestrator tool such as Chef to launch applications,
this model can be applied regardless of where the applications are running.

### Reference Materials:

- [AppRole Auth Method](/docs/auth/approle.html)
  - [AppRole Pull Authentication](/guides/identity/authentication.html)
  - [AppRole with Terraform and Chef Demo](/guides/identity/approle-trusted-entities.html)
- [TLS Certificates Auth Method](/docs/auth/cert.html)
- [Token Auth Method](/docs/auth/token.html)
  - [Cubbyhole Response Wrapping](/guides/secret-mgmt/cubbyhole.html)


## Vault Agent

Vault agent is a client daemon which automates the workflow of client login and
token refresh. It can be used with either [platform
integration](#platform-integration) or [trusted
orchestrator](#trusted-orchestrator) approaches.

#### Vault agent auto-auth:

- Automatically authenticates to Vault for those [supported auth
methods](/docs/agent/autoauth/methods/index.html)
- Keeps token renewed (re-authenticates as needed) until the renewal is no
longer allowed
- Designed with robustness and fault tolerance

![Vault Agent](/img/vault-secure-intro-5.png)

To leverage this feature, run the vault binary in agent mode (`vault agent
-config=<config_file>`) on the client. The agent configuration file must specify
the auth method and [sink](/docs/agent/autoauth/sinks/index.html) locations
where the token to be written.

When the agent is started, it will attempt to acquire a Vault token using the
auth method specified in the agent configuration file.  On successful
authentication, the resulting token is written to the sink locations.
Optionally, this token can be response-wrapped or encrypted. Whenever the
current token value changes, the agent writes to the sinks. If authentication
fails, the agent waits for a while and then retry.

The client can simply retrieve the token from the sink and connect to Vault
using the token. This simplifies client integration since the Vault agent
handles the login and token refresh logic.

### Reference Materials:

- [Streamline Secrets Management with Vault Agent and Vault 0.11](https://youtu.be/zDnIqSB4tyA)
- [Vault Agent documentation](/docs/agent/index.html)
- [Auto-Auth documentation](/docs/agent/autoauth/index.html)


## Next steps

When a [platform integration](#platform-integration) is available that should be
preferred, as it is generally the simpler solution and works independent of the
orchestration mechanism. For a [trusted orchestrator](#trusted-orchestrator),
specific documentation for that orchestrator should be consulted on Vault
integration.
