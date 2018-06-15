---
layout: "guides"
page_title: "PKI Secrets Engine - Guides"
sidebar_current: "guides-secret-mgmt-pki"
description: |-
  The PKI secrets engine generates dynamic X.509 certificates. With this secrets
  engine, services can get certificates without going through the usual manual
  process of generating a private key and CSR, submitting to a CA, and waiting
  for a verification and signing process to complete. Vault's built-in
  authentication and authorization mechanisms provide the verification
  functionality.
---

# Vault PKI Secrets Engine

Use Vault to create X509 certificates for usage in mTLS or other arbitrary PKI
encryption.   While this can be used to create web server certificates.  If
users do not import the CA chains, the browser will complain about self-signed
certificates.

Creating PKI certificates is generally a cumbersome process using traditional
tools like openssl or even more advanced frameworks like CFSSL.   These tools
also require a human component to verify certificate distribution meets
organizational security policies.   

## Reference Material

- [RFC 5280 Internet X.509 Public Key Infrastructure Certificate and Certificate
Revocation List (CRL) Profile](https://tools.ietf.org/html/rfc5280)
- [OpenSSL x509 Man Pages](https://www.openssl.org/docs/man1.1.0/apps/x509.html)

## Estimated Time to Complete

15 minutes

## Personas

The end-to-end scenario described in this guide involves three personas:

- **Operator**: The operator will be responsible for creating the vault/consul
clusters.  In the absence of dedicated security teams, they may also be
responsible for creating all auth methods, ACLs, roles and secret engine mount
points.

- **Developer**: The developer will use the PKI Intermediate issue endpoints to
issue leaf certs. They will most likely use API endpoints to issue certs
directly into applications.

- **InfoSec**: InfoSec will create ACLs, roles and tie together auth methods to
policies. They may also be responsible for enabling the various PKI engines for
root and intermediate CAs.   


## Challenge

Using a vagrant image:

  * Start a minimal Vault/Consul cluster
  * Bootstrap the cluster to be secured with PKI certificates for internode
  communication
  * Create a short-lived TTL PKI Cert.  See it expire
  * Renew the certificate using scripts
  * Revoke the certificate, verify CRL entries
  * Use consul-template to create certs & renew them automatically


## Solution

Run a set of provisioning scripts against your Vault cluster.  This will create:

* A PKI CA
* A PKI Intermediate CA
* Roles and Policies to access the Intermediate CA
* Lower privileged token to use for authentication
* Leaf cert(s) from the Intermediate CA


## Prerequisites

### Vagrant tool

This guide leverages Vegrant to setup a
[Vault](https://www.vagrantup.com/downloads.html) environment. The following
resources are required to perform this demo:

* Linux or Mac system (Windows support in the future)
* [VirtualBox](https://www.virtualbox.org/wiki/Downloads) installed (tested with 5.2.6)
* [Vagrant](https://www.vagrantup.com/downloads.html) installed (tested with 2.0.2)


### Download demo assets

Clone or download the demo assets from the
[hashicorp/vault-guides](https://github.com/hashicorp/vault-guides/tree/master/secrets/pki/vagrant)
GitHub repository to perform the steps described in this guide.

```plaintext
$ git clone https://github.com/hashicorp/vault-guides.git
$ cd vault-guides/secrets/pki/vagrant
```

## Steps

In this guide, you perform the following:

1. [](#step1)
2. [](#step2)



### Step 2: Change to pki directory
```
cd vault-guides/secrets/pki/vagrant
```

### Step 3 (optional): Modify environment variables

There are a number of default environment variables listed in
demo/default_env.sh. For the most part, these variables should suffice for this
guide, but you can override them to customize the guide.   The DEFAULT_*
variables have default settings, just remove DEFAULT_ to actual env var.

Likely envs to override:

```shell
ROOT_DOMAIN=hashidemos.com
VAULT_VERSION=0.10.1
CONSUL_VERSION=1.0.6
CONSUL_TEMPLATE_VERSION=0.19.4
```

### Step 4: Update user-data and provisioning scripts

The user-data and provisioning scripts rely on the environment variables you've
set (or left default) from the previous step.   Run this to create the
customized provisioners:

```plaintext
$ ./run_me_first.sh
```

### Step 5: Start vagrant image

```plaintext
$ vagrant up
```
**NOTE**

Vagrant will need to make some updates to /etc/exports.  This will require sudo
access, so when prompted for your password, this is what it's for.  The same
will be required during destroy to remove the entries

If you're using this guide over time, you may see messages like the following:

```plaintext
==> core-01: A newer version of the box 'coreos-stable' for provider 'virtualbox' is
==> core-01: available! You currently have version '1745.4.0'. The latest is version
==> core-01: '1745.5.0'. Run `vagrant box update` to update.
```

It's best to run the update so CoreOS functions properly.

```plaintext
$ vagrant box update
$ vagrant destroy
$ vagrant up
```


### Step 6: Connect to vagrant with three terminal windows

If you have a terminal like iTerm2 or using Screen/Tmux or any other kind of
terminal window management, open up three sessions to vagrant.   I like to tile
one big window on one side with the other side split horizontally with two
smaller windows

__Diagram created with [Asciiflow](http://asciiflow.com/)__

```plaintext
+--+--+
|  |  |
|  +--+
|  |  |
+--+--+

```

```plaintext
$ vagrant ssh
```

### Step 7: Run demo1_bootstrap_environment.sh

**NOTE**
This guide uses a helper script called demo-magic.sh to prevent users needing to
type/copypasta commands.   For most of the demos, you will need to press
"Enter/Return" after the command is displayed.    There are some demos which
change the behavior to auto-run.   A message will be output when those changes
occur.

You can change this behavior by updating PROMPT_TIMEOUT environment variable to
a non-zero number.  0 means wait for button press, any number is the number of
seconds to wait until proceeding.

```plaintext
$ sudo /demo/demo1_bootstrap_environment.sh
```

This will bootstrap a Vault Dev server to create an initial set of certificates
that will later be used for securing Vault and Consul communicates with TLS.
This instance will be stopped and a new non-dev server will be started to run
the rest of the demo.

### Step 8: Run demo2_short_ttls.sh

This demonstration creates a short lived TTL (Default 60s) and immediately
begins verify the cert with openssl.   You should notice that after the TTL
expires, the cert becomes invalid.   Keep this running in this window.

```plaintext
$ sudo /demo/demo2_short_ttls.sh
```

### Step 9: Run demo3_renew_lease.sh in another window

This demonstration renews the certificate that's not invalid on a regular basis.
You should now see the certificate is valid in the window running demo2.  Leave
this running in this window:

```plaintext
$ sudo /demo/demo3_renew_lease.sh
```

### Step 10: Run demo4_revoke_cert.sh in another window

This demonstration creates a CRL entry in the CRL endpoint registered within
Vault.   You will notice that the CRL has the serial number of the certificate
you revoked into.  However, if you're looking at the demo2 window, you should
notice that the certificate is still valid.  ???  

This is a reminder that when you're validating certificates, make sure to check
the CRL to ensure it hasn't been revoked!

There's another point to this however.   Using short-lived TTLs is a powerful
pattern that often obviates the need to use CRLs because the certificate will
expire in such a short period of time to reduce any kind of damage caused by a
leaked certificate.   

```plaintext
$ sudo /demo/demo4_revoke_cert.sh
```

If you wait until the renewal runs in demo3 window, you can run this command
again.  Notice the change in serial number for the revoked certificate.

### Step 11: Run demo5_consul_template_renew.sh

Before running this, stop the demo2 and demo3 windows (Ctrl-C).   I usually do
this in the big window.

This demonstration combines the concepts of demo2 and demo3 together using
consul-template as the structure for renewal.  Consul-template has the ability
to watch certain Vault endpoints and dump the secrets to a file location of your
choosing.   It's smart enough to know to renew any tokens that it's using that
are expiring.   

```plaintext
$ sudo /demo/demo5_consul_template_renew.sh
```

After running this for 2-3 minutes (or however long you want), stop the demo
using Ctrl-c.

**NOTE**

Vault doesn't have the concept of blocking queries like Consul does.   This
means that it doesn't have the capabilities of noticing changes to a particular
secret you are monitoring.  However, in the case of PKI certs, consul-template
will renew based on the lease for the certificate.   Keep this in mind of you're
trying to use consul-template for automating retrieval of other Vault secrets.

### Step 12: Tear down the guide

Congratulations!  You've taken a step into the wonderful world of PKI
Provisioning with Vault!   To be complete, clean up your guide environment.
Also, remember that vagrant will ask for your sudo password to remove NFS
entries from /etc/exports:

```plaintext
$ exit
$ vagrant destroy
```










## Next steps

The use of [AppRole Pull Authentication](/guides/identity/authentication.html) is a good
use case to leverage the response wrapping. Go through the guide if you have not
done so.  To better understand the lifecycle of Vault tokens, proceed to [Tokens
and Leases](/guides/identity/lease.html) guide.
