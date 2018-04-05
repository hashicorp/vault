---
layout: "guides"
page_title: "Vault Auto-unseal using AWS KMS - Guides"
sidebar_current: "guides-autounseal-aws-kms"
description: |-
  In this guide, we'll show an example of how to use Terraform to provision an
  instance that can utilize an encryption key from AWS Key Management Services
  to unseal Vault.
---


# Vault Auto Unseal using AWS Key Management Service

~> **Enterprise Only:** Vault replication feature is a part of _Vault Enterprise_.

When a Vault server is started, it starts in a
[***sealed***](/docs/concepts/seal.html) state and it does not know how to
decrypt data. Before any operation can be performed on the Vault, it must be
unsealed. Unsealing is the process of constructing the master key necessary to
decrypt the data encryption key.

![Unseal with Shamir's Secret Sharing](/assets/images/vault-autounseal.png)

This guide demonstrates an example of how to use Terraform to provision an
instance that can utilize an encryption key from [AWS Key Management Services
(KMS)](https://aws.amazon.com/kms/) to unseal Vault.

## Reference Material

- [Vault Enterorise Auto Unseal](/docs/enterprise/auto-unseal/index.html)
- [Configuration: `awskms` Seal](/docs/configuration/seal/awskms.html)


## Estimated Time to Complete

10 minutes

## Personas

The steps described in this guide are typically performed by **operations**
persona.

## Challenge

Vault unseal operation requires a quorum of existing unseal keys split by
Shamir's Secret sharing algorithm. This is done so that the "_keys to the
kingdom_" won't fall into one person's hand.  However, this process is manual
and can become painful when you have many Vault clusters as there are now
many different key holders with many different keys.

## Solution

Vault Enterprise supports opt-in automatic unsealing via cloud technologies such
Amazon KMS or Google Cloud KMS. This feature enables operators to delegate the
unsealing process to trusted cloud providers to ease operations in the event of
partial failure and to aid in the creation of new or ephemeral clusters.

![Unseal with Shamir's Secret Sharing](/assets/images/vault-autounseal-2.png)

## Prerequisites

This guide assumes the following:   

- Access to **Vault Enterprise 0.9.0 or later** which supports AWS KMS as an unseal mechanism
- A URL to download Vault Enterprise from (an Amazon S3 bucket will suffice)
- AWS account for provisioning cloud resources
- [Terraform installed](https://www.terraform.io/intro/getting-started/install.html)
and basic understanding of its usage

### Download demo assets

Clone or download the demo assets from the
[hashicorp/vault-guides](https://github.com/hashicorp/vault-guides/tree/master/operations/aws-kms-unseal/terraform) GitHub repository to perform the steps described in this guide.


## Steps

This guide demonstrates how to implement and use the Auto Unseal feature using
AWS KMS. Included is a Terraform configuration that has the following:   

* Ubuntu 16.04 LTS with Vault Enterprise    
* An instance profile granting the Amazon EC2 instance to an AWS KMS key
* Vault configured with access to an AWS KMS key   

In this guide, you are going to perform the following steps:

1. [Provision the Cloud Resources](#step-1-provision-the-cloud-resources)
1. [Verification](#step-2-verification)
1. [Clean Up](#step-3-clean-up)


### Step 1: Provision the Cloud Resources

**Task 1:** Be sure to set your working directory to where the
`/aws-kms-unseal/terraform` folder is located.

Terraform files should be located in this directory:

```bash
~/git/vault-guides/operations/aws-kms-unseal/terraform$ tree
.
├── README.md
├── instance-profile.tf
├── instance.tf
├── main.tf
├── ssh-key.tf
├── terraform.tfvars.example
├── userdata.tpl
└── variables.tf
```

**Task 2:** Set your AWS credentials as environment variables:

```bash
$ export AWS_ACCESS_KEY_ID = "<YOUR_AWS_ACCESS_KEY_ID>"

$ export AWS_SECRET_ACCESS_KEY = "<YOUR_AWS_SECRET_ACCESS_KEY>"
```

Specify your Vault Enterprise URL in a file named **`terraform.tfvars`**.

An example is provided (`terraform.tfvars.example`):

```plaintext
vault_url = "http://s3.amazonaws.com/some/path/to/vault-enterprise.zip"
```

**Task 3:** Perform a `terraform init` to pull down the necessary provider
resources. Then `terraform plan` to verify your changes and the resources that
will be created. If all looks good, then perform a `terraform apply` to
provision the resources. The Terraform output will display the public IP address
to SSH into your server.

```plaintext
$ terraform init
Initializing provider plugins...
...
Terraform has been successfully initialized!


$ terraform plan
...
Plan: 15 to add, 0 to change, 0 to destroy.


$ terraform apply
...
Apply complete! Resources: 15 added, 0 changed, 0 destroyed.

Outputs:

connections = Connect to Vault via SSH   ssh ubuntu@192.0.2.1 -i private.key
Vault Enterprise web interface  http://192.0.2.1:8200/ui
```

**NOTE:** Outputs will contain instructions to connect to the server via SSH as
shown in the example above.

### Step 2: Verification

SSH into the provisioned EC2 instance.

```plaintext
$ ssh ubuntu@192.0.2.1 -i private.key

The authenticity of host '34.201.166.196 (34.201.166.196)' can't be established.
ECDSA key fingerprint is SHA256:B3FBKHAxBP/oW84GN74EOV+XPeVC00juipfgMPgo5Kc.
Are you sure you want to continue connecting (yes/no)? yes
```
When you are prompted, enter "yes" to continue.


```plaintext
$ vault status
Error checking seal status: Error making API request.

URL: GET http://127.0.0.1:8200/v1/sys/seal-status
Code: 400. Errors:

* server is not yet initialized

$ vault init -stored-shares=1 -recovery-shares=1 -recovery-threshold=1 -key-shares=1 -key-threshold=1
Recovery Key 1: oOxAQfxcZitjqZfF3984De8rUckPeahQDUvmJ1A4JrQ=
Initial Root Token: 54c4dbe3-d45b-79d9-18d0-602831a6a991

Vault initialized successfully.

Recovery key initialized with 1 keys and a key threshold of 1. Please
securely distribute the above keys.

$ systemctl stop vault
root@ip-192-168-100-100:~# vault status
Error checking seal status: Get http://127.0.0.1:8200/v1/sys/seal-status: dial tcp 127.0.0.1:8200: getsockopt: connection refused

$ systemctl start vault
root@ip-192-168-100-100:~# vault status
Type: shamir
Sealed: false
Key Shares: 1
Key Threshold: 1
Unseal Progress: 0
Unseal Nonce:
Version: 0.9.0+prem.hsm
Cluster Name: vault-cluster-01cf6f33
Cluster ID: fb787d8a-b882-fee8-b461-445320cde311

High-Availability Enabled: false

$ vault auth 54c4dbe3-d45b-79d9-18d0-602831a6a991
Successfully authenticated! You are now logged in.
token: 54c4dbe3-d45b-79d9-18d0-602831a6a991
token_duration: 0
token_policies: [root]


$ cat /etc/vault.d/vault.hcl
storage "file" {
  path = "/opt/vault"
}
listener "tcp" {
  address     = "0.0.0.0:8200"
  tls_disable = 1
}
seal "awskms" {
  kms_key_id = "d7c1ffd9-8cce-45e7-be4a-bb38dd205966"
}
ui=true
```


### Step 3: Clean Up

Once completed, execute the following commands to clean up:

```plaintext
$ terraform destroy -force

$ rm -rf .terraform terraform.tfstate* private.key
```


## Next steps

Once you have a Vault environment setup, the next step is to write policies.
Read [Policies](/guides/identity/policies.html) to learn how to write policies
to govern the behavior of clients.
