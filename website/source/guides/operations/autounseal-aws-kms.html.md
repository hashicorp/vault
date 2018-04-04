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
[sealed](/docs/concepts/seal.html) state and must be unsealed before you can
work with it.  Normally, the unsealing process


Vault Enterprise supports opt-in automatic unsealing via cloud technologies such
Amazon KMS or Google Cloud KMS. This feature enables operators to delegate the
unsealing process to trusted cloud providers to ease operations in the event of
partial failure and to aid in the creation of new or ephemeral clusters.

This feature enables operators to delegate the unsealing process to trusted
cloud providers to ease operations in the event of partial failure and to aid in
the creation of new or ephemeral clusters.




In this guide, we'll show an example of how to use Terraform to provision an
instance that can utilize an encryption key from AWS Key Management Services
(KMS) to unseal Vault.

## Overview
Vault unseal operation either requires either a number of people who each possess a shard of a key, split by Shamir's Secret sharing algorithm, or protection of the master key via an HSM or cloud key management services (Google CKMS or AWS KMS).

This guide has a guide on how to implement and use this feature in AWS. Included is a Terraform configuration that has the following features:  
* Ubuntu 16.04 LTS with Vault Enterprise (0.9.0+prem.hsm).   
* An instance profile granting the AWS EC2 instance to a KMS key.   
* Vault configured with access to a KMS key.   


## Prerequisites

This guide assumes the following:   

1. Access to Vault Enterprise > 0.9.0 which supports AWS KMS as an unseal mechanism.
1. A URL to download Vault Enterprise from (an S3 bucket will suffice).
1. AWS account for provisioning cloud resources.
1. Terraform installed, and basic understanding of its usage


## Usage
Instructions assume this location as a working directory, as well as AWS credentials exposed as environment variables

1. Set Vault Enterprise URL in a file named terraform.tfvars (see terraform.tfvars.example)
1. Perform the following to provision the environment

```
terraform init
terraform plan
terraform apply
```

Outputs will contain instructions to connect to the server via SSH


```
# vault status
Error checking seal status: Error making API request.

URL: GET http://127.0.0.1:8200/v1/sys/seal-status
Code: 400. Errors:

* server is not yet initialized

# vault init -stored-shares=1 -recovery-shares=1 -recovery-threshold=1 -key-shares=1 -key-threshold=1
Recovery Key 1: oOxAQfxcZitjqZfF3984De8rUckPeahQDUvmJ1A4JrQ=
Initial Root Token: 54c4dbe3-d45b-79d9-18d0-602831a6a991

Vault initialized successfully.

Recovery key initialized with 1 keys and a key threshold of 1. Please
securely distribute the above keys.

# systemctl stop vault
root@ip-192-168-100-100:~# vault status
Error checking seal status: Get http://127.0.0.1:8200/v1/sys/seal-status: dial tcp 127.0.0.1:8200: getsockopt: connection refused

# systemctl start vault
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

# vault auth 54c4dbe3-d45b-79d9-18d0-602831a6a991
Successfully authenticated! You are now logged in.
token: 54c4dbe3-d45b-79d9-18d0-602831a6a991
token_duration: 0
token_policies: [root]


# cat /etc/vault.d/vault.hcl
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




Once complete perform the following to clean up

```
terraform destroy -force
rm -rf .terraform terraform.tfstate* private.key

```





## Next steps

Read [Production Hardening](/guides/operations/production.html) to learn best
practices for a production hardening deployment of Vault.
