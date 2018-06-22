---
layout: "guides"
page_title: "PKI Secrets Engine Demo - Guides"
sidebar_current: "guides-secret-mgmt-pki"
description: |-
  The PKI secrets engine generates dynamic X.509 certificates. With this secrets
  engine, services can get certificates without going through the usual manual
  process of generating a private key and CSR, submitting to a CA, and waiting
  for a verification and signing process to complete. Vault's built-in
  authentication and authorization mechanisms provide the verification
  functionality.
---

# PKI Secrets Engine

Use Vault to create X509 certificates for usage in MTLS or other arbitrary PKI
encryption.   While this can be used to create web server certificates.  If
users do not import the CA chains, the browser will complain about self-signed
certificates.

Creating PKI certificates is generally a cumbersome process using traditional
tools like openssl or even more advanced frameworks like CFSSL.   These tools
also require a human component to verify certificate distribution meets
organizational security policies.   

The purpose of this guide is to provide the instruction to reproduce the working
implementation demo introduced in the [Streamline Certificate Management with HashiCorp Vault](https://www.hashicorp.com/resources/streamline-certificate-management-with-vault)
webinar.

[![YouTube](/assets/images/vault-pki-demo-1.png)](https://youtu.be/k8FXTeFCp90)


## Reference Material

- [Streamline Certificate Management with HashiCorp Vault](https://www.hashicorp.com/resources/streamline-certificate-management-with-vault)
- [RFC 5280 Internet X.509 Public Key Infrastructure Certificate and Certificate
Revocation List (CRL) Profile](https://tools.ietf.org/html/rfc5280)
- [OpenSSL x509 Man Pages](https://www.openssl.org/docs/man1.1.0/apps/x509.html)
- [PKI (Certificates) Secrets Engine](/docs/secrets/pki/index.html)
- [PKI Secrets Engine (API)](/api/secret/pki/index.html)


## Estimated Time to Complete

15 minutes


## Challenge

Traditional PKI process workflow looks like:

1. Create Certificates Signing Request (CSR)
    - Generate public and private keys
    - Sign CSR with your private key
1. Submit your CSR to Certificate Authority (CA)
1. CA usually signs the CSR and returns the public key which is your certificate

If revocation is needed, you have to update your locate Certificate Revocation
List (CRL) or Online Certificate Status Protocol (OCSP).


Typical SSH login methods:

- Host-based
- Username and Password
- SSH keys

The host must be well secured.

## Solution

Vault PKI secrets engine makes this a lot simpler.  The PKI secrets engine can
be an Intermediate-Only certificate authority which potentially allows for
higher levels of security.

1. Store CA outside the Vault (air gapped)
1. Create CSRs for the intermediates
1. Sign CSR outside Vault and import intermediate
1. Issue leaf certificates from the Intermediate CA


## Prerequisites

The following resources are required to perform the demo described in this guide:

* Linux or iOS (Windows support in the future)
* [VirtualBox](https://www.virtualbox.org/wiki/Downloads) (tested with `5.2.6`)
* [Vagrant](https://www.vagrantup.com/downloads.html) (tested with `2.0.2`)


### Download the demo assets

Clone or download the demo assets from the
[hashicorp/vault-guides](https://github.com/hashicorp/vault-guides/tree/master/secrets/pki/vagrant)
GitHub repository to perform the steps described in this guide.

```plaintext
$ git clone https://github.com/hashicorp/vault-guides.git
$ cd vault-guides/secrets/pki/vagrant
```

## Steps

-> **NOTE:** This guide leverages Vegrant to setup a
[Vault](https://www.vagrantup.com/downloads.html) environment. Make sure that
you have installed Vagrant as well as VirtualBox as specified in the
[Prerequisites](#prerequisites).


In this guide, you perform the following:

1. [Provision a demo virtual machine](#step1)
1. [Run demo \#1](#step2)
1. [Run demo \#2](#step3)
1. [Run demo \#3](#step4)
1. [Run demo \#4](#step5)
1. [Run demo \#5](#step6)
1. [Tear down the demo environment](#step7)


### <a name="step1"></a>Step 1: Provision a demo virtual machine

1. First, change your working directory to `pki/vagrant/demo`.

    ```plaintext
    $ cd vault-guides/secrets/pki/vagrant/demo
    ```

1. Locate the **`default_env.sh`** file in which a number of environment variables
are listed.  Some of them have default values set with a matching variable name
prefixed with "`DEFAULT_`".

    ```shell
    # These are all the defaults for any environment variables below.  Setting environment variables before accessing this set of defaults will override anything set here.
    DEFAULT_ROOT_DOMAIN=hashidemos.com

    # Software versions.   These are probably the only values that should change over time
    DEFAULT_VAULT_VERSION=0.10.1
    DEFAULT_CONSUL_VERSION=1.0.6
    DEFAULT_CONSUL_TEMPLATE_VERSION=0.19.4
    ...

    VAULT_VERSION=${VAULT_VERSION:-$DEFAULT_VAULT_VERSION}
    ...

    CONSUL_VERSION=${CONSUL_VERSION:-$DEFAULT_CONSUL_VERSION}
    ...

    CONSUL_TEMPLATE_VERSION=${CONSUL_TEMPLATE_VERSION:-$DEFAULT_CONSUL_TEMPLATE_VERSION}
    ...
    ```

    The default values should be sufficient enough to perform this guide. If you
    wish to overwrite any of the values, set the values of the environment
    variables without the "`DEFAULT_`" prefix.

    **Example:**

    ```shell
    # These are all the defaults for any environment variables below.  Setting environment variables before accessing this set of defaults will override anything set here.
    DEFAULT_ROOT_DOMAIN=hashidemos.com

    # Software versions.   These are probably the only values that should change over time
    DEFAULT_VAULT_VERSION=0.10.1
    DEFAULT_CONSUL_VERSION=1.0.6
    DEFAULT_CONSUL_TEMPLATE_VERSION=0.19.4
    ...

    # Overwrite the default of 0.10.1
    VAULT_VERSION=0.10.3
    ...
    # Overwrite the default of 1.0.6
    CONSUL_VERSION=1.1.0
    ...
    # Overwrite the default of 0.19.4
    CONSUL_TEMPLATE_VERSION=0.19.5
    ...
    ```

1. Run the `run_me_first.sh` script in the
`vault-guides/secrets/pki/vagrant` directory to create the customized
provisioners.

    The user-data and provisioning scripts rely on the environment variables you've
    set (or left default) from the previous step.   

    ```plaintext
    $ cd ..
    $ ./run_me_first.sh
    ```

1. Now, spin up a vagrant image by executing the `vagrant` command.

    ```plaintext
    $ vagrant up

    ...
    ==> core-01: Exporting NFS shared folders...
    ==> core-01: Preparing to edit /etc/exports. Administrator privileges will be required...
    Password:
    ```

    ~> **NOTE:** Vagrant will need to make some updates to `/etc/exports`.  This
    will require **sudo** access that you will be prompted to enter your
    password.  The same will be required during the `vagrant destroy` as well.

1. Open **three terminal sessions** and connect to the virtual machine
provisioned by Vagrant:

    ```plaintext
    $ vagrant ssh
    ```

    > If you have a terminal like [iTerm2](https://www.iterm2.com/downloads.html),
    Screen, [tmux](https://github.com/tmux/tmux/wiki) or any other kind of terminal
    window management tool, open up three sessions as shown.   

![Vagrant Sessions](/assets/images/vault-pki-demo.png)


### <a name="step2"></a>Step 2: Run demo \#1

In one of the terminal sessions, execute the `demo1_bootstrap_environment.sh`
script in the `/demo` directory.

```plaintext
$ cd /demo
$ sudo ./demo1_bootstrap_environment.sh
```

This demo script bootstraps a Vault development server to create an initial set
of certificates that will later be used for securing Vault and Consul
communicates over TLS.

-> A short description of each task gets displayed along with the vault command.
Review the command and then press **Return** or **Enter** key on your keyboard
to execute each command.

<br>
Commands getting executed in this demo are:

```shell
# Enabling CA Certificate PKI Secret Engine
vault secrets enable -path=pki_root pki

# Tuning CA Certificate PKI Secret Engine
vault secrets tune -max-lease-ttl=87648h pki_root

# Generating CA Root certificate. JSON output found at /var/tmp/root_ca_output.json
curl -s --header "X-Vault-Token: 0bf5d133-5b15-711b-b19e-110db6994d18" \
      --request POST \
      --data @/var/tmp/input-params.json \
      http://172.18.0.3:8200/v1/pki_root/root/generate/internal > /var/tmp/root_ca_output.json

# Updating CRL And CA information
vault write pki_root/config/urls \
      issuing_certificates=http://172.18.0.3:8200/v1/pki_root/ca \
      crl_distribution_points=http://172.18.0.3:8200/v1/pki_root/crl     

# Enabling PKI Secret Engine at pki_int_main
vault secrets enable -path=pki_int_main pki

# Tuning PKI Secret Engine pki_int_main
vault secrets tune -max-lease-ttl=43824h pki_int_main

# Generating Intermediate CSR. JSON output found at /var/tmp/intermediate_csr_output.json CSR located at /var/tmp/pki_int_main.csr
curl -s --header "X-Vault-Token: 0bf5d133-5b15-711b-b19e-110db6994d18" \
        --request POST \
        --data @/var/tmp/input-params.json \
        http://172.18.0.3:8200/v1/pki_int_main/intermediate/generate/internal \
        > /var/tmp/intermediate_csr_output.json

# Generating Intermediate PEM cert. JSON output found at /var/tmp/intermediate_csr_output.json
vault write -format=json pki_root/root/sign-intermediate \
      csr=@/var/tmp/docker/pki_int_main.csr \
      format=pem_bundle  > /var/tmp/intermediate_csr_output.json

# Store the certificate in /var/tmp/pki_int_main.pem
jq -r '.data.certificate' /var/tmp/intermediate_csr_output.json > /var/tmp/pki_int_main.pem

vault write pki_int_main/intermediate/set-signed certificate=@/var/tmp/docker/pki_int_main.pem

# Updating CRL And CA information for pki_int_main
vault write pki_int_main/config/urls \
      issuing_certificates=http://172.18.0.3:8200/v1/pki_int_main/ca \
      crl_distribution_points=http://172.18.0.3:8200/v1/pki_int_main/crl

# Roles are the entities allowed to create certificates
# Creating role hashidemos-com in pki_int_main
vault write pki_int_main/roles/hashidemos-com \
      allowed_domains=hashidemos.com \
      allow_subdomains=true \
      max_ttl=43824h \
      allow_any_name=true \
      generate_lease=true \
      enforce_hostnames=false

# Creating role consul-dev-hashidemos-com in pki_int_main
vault write pki_int_main/roles/consul-dev-hashidemos-com \
      allowed_domains=dev.hashidemos.com \
      allow_subdomains=true \
      max_ttl=43824h \
      allow_any_name=true \
      generate_lease=true \
      enforce_hostnames=false

# Creating role vault-dev-hashidemos-com in pki_int_main
vault write pki_int_main/roles/vault-dev-hashidemos-com \
      allowed_domains=dev.hashidemos.com \
      allow_subdomains=true \
      max_ttl=43824h \
      allow_any_name=true \
      generate_lease=true \
      enforce_hostnames=false

# Creating cert for hashidemos.com. JSON Output: /var/tmp/hashidemos.com_cert.json
vault write -format=json pki_int_main/issue/hashidemos-com \
      common_name=hashidemos.com \
      ttl=168h > /var/tmp/hashidemos.com_cert.json

# Creating cert for consul1.dev.hashidemos.com. JSON Output: /var/tmp/consul1.dev.hashidemos.com_cert.json
vault write -format=json pki_int_main/issue/consul-dev-hashidemos-com \
      common_name=consul1.dev.hashidemos.com \
      ttl=168h > /var/tmp/consul1.dev.hashidemos.com_cert.json

# Creating cert for vault1.dev.hashidemos.com. JSON Output: /var/tmp/vault1.dev.hashidemos.com_cert.json
vault write -format=json pki_int_main/issue/vault-dev-hashidemos-com \
      common_name=vault1.dev.hashidemos.com \
      ttl=168h > /var/tmp/vault1.dev.hashidemos.com_cert.json
```

When you finish, the development server is stopped and a new non-development
server will be started to run the rest of the demo.



### <a name="step3"></a>Step 3: Run demo \#2

Now, run the `demo2_short_ttls.sh` script. This script creates a short lived TTL
(default 60s) and immediately begins verify the cert with openssl. You should
notice that after the TTL expires, the cert becomes invalid.

```plaintext
$ cd /demo
$ sudo ./demo2_short_ttls.sh
```

-> A short description of each task gets displayed along with the vault command.
Review the command and then press **Return** or **Enter** key on your keyboard
to execute each command.

<br>
Commands getting executed in this demo are:

```shell
# Creating role short-ttl-hashidemos-com in pki_int_main
vault write pki_int_main/roles/short-ttl-hashidemos-com \
      allowed_domains=dev.hashidemos.com \
      allow_subdomains=true \
      max_ttl=43824h \
      allow_any_name=true \
      generate_lease=true \
      enforce_hostnames=false

# Creating policy to assign to token(s)
cat > /var/tmp/policy-short-ttl-hashidemos-com << EOF
path "pki_int_main/issue*" {
    capabilities = ["create","update"]
}

path "pki_root/cert/ca" {
    capabilities = ["read"]
}

path "auth/token/renew" {
    capabilities = ["update"]
}

path "auth/token/renew-self" {
    capabilities = ["update"]
}
EOF


#  Create short-ttl-hashidemos-com policy
vault policy write short-ttl-hashidemos-com /var/tmp/docker/policy-short-ttl-hashidemos-com

# Creating token with single policy into /var/tmp/token-short-ttl-hashidemos-com
vault token create -policy=short-ttl-hashidemos-com \
      -field=token -ttl=120s > /var/tmp/token-short-ttl-hashidemos-com

# Creating cert for shortttl.dev.hashidemos.com. JSON Output: /var/tmp/shortttl.dev.hashidemos.com_cert.json
vault write -format=json pki_int_main/issue/short-ttl-hashidemos-com \
      common_name=shortttl.dev.hashidemos.com \
      ttl=60 > /var/tmp/shortttl.dev.hashidemos.com_cert.json
```

~> **NOTE:** Keep the script running in this window.


### <a name="step4"></a>Step 4: Run demo \#3

Run the `demo3_renew_lease.sh` in **another** terminal session window.

```plaintext
$ cd /demo
$ sudo ./demo3_renew_lease.sh
```
This script renews the certificate that's not invalid every 30 seconds.

```shell
# Creating cert for shortttl.dev.hashidemos.com. JSON Output: /var/tmp/shortttl.dev.hashidemos.com_cert.json
vault write -format=json pki_int_main/issue/short-ttl-hashidemos-com \
      common_name=shortttl.dev.hashidemos.com \
      ttl=60 > /var/tmp/shortttl.dev.hashidemos.com_cert.json
```

You should now see the certificate is now **valid** in the window running
`demo2_short_ttls.sh`.

![Vagrant Sessions](/assets/images/vault-pki-demo-2.png)


~> **NOTE:** Keep the script running in this window.



### <a name="step5"></a>Step 5: Run demo \#4

Run `demo4_revoke_cert.sh` in another window.

```plaintext
$ cd /demo
$ sudo ./demo4_revoke_cert.sh
```

This script creates a CRL entry in the CRL endpoint registered within
Vault.   

```plaintext
Revoking from pki_int_main serial 30:22:5b:73:f0:e1:35:f3:5b:eb:51:ab:ad:e9:a2:00:b1:65:2b:be
{"request_id":"bee24fc0-2043-f3e9-1a82-8b290f7b0ea0","lease_id":"","renewable":false,"lease_duration":0,"data":{"revocation_time":1529455943,"revocation_time_rfc3339":"2018-06-20T00:52:23.141146857Z"},"wrap_info":null,"warnings":null,"auth":null}
Checking CRL
Certificate Revocation List (CRL):
        Version 2 (0x1)
    Signature Algorithm: sha256WithRSAEncryption
        Issuer: /CN=hashidemos.com Intermediate 1
        Last Update: Jun 20 00:52:23 2018 GMT
        Next Update: Jun 23 00:52:23 2018 GMT
        CRL extensions:
            X509v3 Authority Key Identifier:
                keyid:BD:D1:C0:06:7C:91:83:81:D1:53:B4:F7:02:67:41:A2:B2:3E:3B:85

Revoked Certificates:
    Serial Number: 30225B73F0E135F35BEB51ABADE9A200B1652BBE
        Revocation Date: Jun 20 00:52:23 2018 GMT
    Signature Algorithm: sha256WithRSAEncryption
         19:9f:25:2d:6d:2d:4f:aa:f0:83:d5:ae:d7:ba:9e:31:89:01:
         ...
         27:a4:09:3d
```

You will notice that the CRL has the serial number of the certificate you
revoked into.  However, if you're looking at the terminal window which is
running `demo2_short_ttls.sh`, you should notice that the certificate is still
valid. This is a reminder that when you're validating certificates, make sure to
check the CRL to ensure it hasn't been revoked.

There's another point to this however. Using short-lived TTLs is a powerful
pattern that often obviates the need to use CRLs because the certificate will
expire in such a short period of time to reduce any kind of damage caused by a
leaked certificate.   


If you wait until the renewal runs in the `demo3_renew_lease.sh` terminal, you
can run this command again.  Notice the change in serial number for the revoked
certificate.


### <a name="step6"></a>Step 6: Run demo \#5

Stop the **`demo2_short_ttls.sh`** and **`demo3_renew_lease.sh`** that are
running by pressing **Ctrl+C**, and then execute the
**`demo5_consul_template_renew.sh`** script in one of the terminals.

```plaintext
$ sudo ./demo5_consul_template_renew.sh
```

This demo script combines the concepts of `demo2_short_ttls.sh` and
`demo3_renew_lease.sh` together using Consul Teamplate as the structure for
renewal.  Consul-template has the ability to watch certain Vault endpoints and
dump the secrets to a file location of your choosing.   It's smart enough to
know to renew any tokens that it's using that are expiring.   

-> A short description of each task gets displayed along with the vault command.
Review the command and then press **Return** or **Enter** key on your keyboard
to execute each command.

<br>
Commands getting executed in this demo are:

```shell
# Write policy file which gets attached to token(s)
cat > /var/tmp/policy-short-ttl-hashidemos-com << EOF
path "pki_int_main/issue*" {
    capabilities = ["create","update"]
}

path "pki_root/cert/ca" {
    capabilities = ["read"]
}

path "auth/token/renew" {
    capabilities = ["update"]
}

path "auth/token/renew-self" {
    capabilities = ["update"]
}
EOF

# Create short-ttl-hashidemos-com policy
vault policy write short-ttl-hashidemos-com /var/tmp/docker/policy-short-ttl-hashidemos-com

# Creating token with single policy into /var/tmp/token-short-ttl-hashidemos-com
vault token create -policy=short-ttl-hashidemos-com -field=token \
      -ttl=120s > /var/tmp/token-short-ttl-hashidemos-com

# Creating Consul Teamplate template, cert.tmpl
cat /etc/consul_templates/pki/shortttl.dev.hashidemos.com/cert.tmpl
{{- with secret "pki_int_main/issue/short-ttl-hashidemos-com" "common_name=shortttl.dev.hashidemos.com" "ttl=120s" -}}
{{ .Data.certificate }}{{ end }}

# Creating Consul Teamplate template, ca.tmpl
cat /etc/consul_templates/pki/shortttl.dev.hashidemos.com/ca.tmpl
{{ with secret "pki_root/cert/ca" -}}
{{ .Data.certificate }}{{ end }}
{{ with secret "pki_int_main/issue/short-ttl-hashidemos-com" "common_name=shortttl.dev.hashidemos.com" "ttl=120s" -}}
{{ .Data.issuing_ca }}{{ end }}

# Creating Consul Teamplate template, key.tmpl
cat /etc/consul_templates/pki/shortttl.dev.hashidemos.com/key.tmpl
{{ with secret "pki_int_main/issue/short-ttl-hashidemos-com" "common_name=shortttl.dev.hashidemos.com" "ttl=120s" -}}
{{ .Data.private_key }}{{ end }}

# Creating Consul Teamplate template, serial.tmpl
cat /etc/consul_templates/pki/shortttl.dev.hashidemos.com/serial.tmpl
{{ with secret "pki_int_main/issue/short-ttl-hashidemos-com" "common_name=shortttl.dev.hashidemos.com" "ttl=120s" -}}
{{ .Data.serial_number }}{{ end }}

# Creating Consul Teamplate config file, consul_template.tmpl
cat /etc/consul_templates/pki/shortttl.dev.hashidemos.com/consul_template.tmpl
  vault {
    address = "https://vault1.dev.hashidemos.com:8200"
    token = "652ded21-c06d-371a-b0c6-f27e25537b83"
    ssl {
      cert    = "/etc/certs/vault1.dev.hashidemos.com_crt.pem"
      key     = "/etc/certs/vault1.dev.hashidemos.com_key.pem"
      ca_cert = "/etc/certs/vault1.dev.hashidemos.com_ca_chain_full.pem"
    }
  }
  template {
    source      = "/etc/consul_templates/pki/shortttl.dev.hashidemos.com/cert.tmpl"
    destination = "/var/tmp/shortttl.dev.hashidemos.com_crt.pem"
  }

  template {
    source      = "/etc/consul_templates/pki/shortttl.dev.hashidemos.com/key.tmpl"
    destination = "/var/tmp/shortttl.dev.hashidemos.com_key.pem"
  }

  template {
    source      = "/etc/consul_templates/pki/shortttl.dev.hashidemos.com/ca.tmpl"
    destination = "/var/tmp/shortttl.dev.hashidemos.com_ca_chain_full.pem"
  }

  template {
    source      = "/etc/consul_templates/pki/shortttl.dev.hashidemos.com/serial.tmpl"
    destination = "/var/tmp/shortttl.dev.hashidemos.com.serial"
  }

# Restoring demo prompt waits
# Testing certificate validity using openssl
01:02:45: Verifying Certificate: [root@core-01] /opt/bin/consul-template -config=/etc/consul_templates/pki/shortttl.dev.hashidemos.com/consul_template.tmpl

Certificate /var/tmp/shortttl.dev.hashidemos.com_crt.pem is invalid
01:02:50: Verifying Certificate: Certificate /var/tmp/shortttl.dev.hashidemos.com_crt.pem is valid
01:02:55: Verifying Certificate: Certificate /var/tmp/shortttl.dev.hashidemos.com_crt.pem is valid
...
```

After running this for 2-3 minutes (or however long you want), stop the demo by pressing **Ctrl+C**.

<br>

-> **NOTE:** Vault doesn't have the concept of blocking queries like Consul does.
This means that it doesn't have the capabilities of noticing changes to a
particular secret you are monitoring.  However, in the case of PKI certs,
[Consul Template](https://github.com/hashicorp/consul-template) will renew based on the lease for the certificate. Keep this
in mind of you're trying to use Consul Template for automating retrieval of
other Vault secrets.


### <a name="step7"></a>Step 7: Tear down the demo environment

Clean up your environment.

```plaintext
$ exit
$ vagrant destroy
```

~> **NOTE:** Vagrant will need to make some updates to `/etc/exports`.  This
will require **sudo** access that you will be prompted to enter your
password during the `vagrant destroy`.



## Next steps

The use of [AppRole Pull Authentication](/guides/identity/authentication.html) is a good
use case to leverage the response wrapping. Go through the guide if you have not
done so.  To better understand the lifecycle of Vault tokens, proceed to [Tokens
and Leases](/guides/identity/lease.html) guide.
