---
layout: "guides"
page_title: "AppRole With Terraform & Chef - Guides"
sidebar_title: "AppRole with Terraform and Chef"
sidebar_current: "guides-identity-approle-tf-chef"
description: |-
  This guide discusses the concepts necessary to help users
  understand Vault's AppRole authentication pattern and how to use it to
  securely introduce a Vault authentication token to a target server,
  application, container, etc.x
---

# Vault AppRole with Terraform and Chef Demo

In the [AppRole Pull
Authentication](/guides/identity/authentication.html#advanced-features) guide,
the question of how best to deliver the Role ID and Secret ID were brought up,
and the role of trusted entities (Terraform, Chef, Nomad, Kubernetes, etc.) was
mentioned.

![AppRole auth method workflow](/img/vault-approle-workflow2.png)

This _intermediate_ Vault guide aims to provide a **simple**, **end-to-end**
example of how to use Vault's [AppRole authentication
method](/docs/auth/approle.html), along with Terraform and Chef, to address the
challenge of the **_secure introduction_** of an initial token to a target
system.

The purpose of this guide is to provide the instruction to reproduce the working
implementation demo introduced in the [Delivering Secret Zero: Vault AppRole
with Terraform and
Chef](https://www.hashicorp.com/resources/delivering-secret-zero-vault-approle-terraform-chef)
webinar.

[![YouTube](/img/vault-approle-youtube.png)](https://youtu.be/OIcIzFWjThM)

-> **NOTE:** This is a proof of concept and **NOT SUITABLE FOR PRODUCTION USE**.


## Reference Material

- [AppRole Auth Method](/docs/auth/approle.html)
- [Authenticating Applications with HashiCorp Vault AppRole](https://www.hashicorp.com/blog/authenticating-applications-with-vault-approle)
- [Delivering Secret Zero: Vault AppRole with Terraform and
Chef](https://www.hashicorp.com/resources/delivering-secret-zero-vault-approle-terraform-chef)

## Estimated Time to Complete

20 minutes

## Challenge

The goal of the AppRole authentication method is to provide a mechanism for the
secure introduction of secrets to target systems (servers, applications,
containers, etc.).

The question becomes what systems within our environment do we trust to handle
or deliver the `RoleID` and `SecretID` to our target systems.


## Solution

Use _Trusted Entities_ to deliver the AppRole authentication values. For
example, use Terraform to deliver your `RoleID` or embed it into your AMI or
Dockerfile. Then you might use Jenkins or Chef to obtain the
[response-wrapped](/guides/secret-mgmt/cubbyhole.html) `SecretID` and deliver it
to the target system.

AppRole allows us to securely introduce the authentication token to the target
system by preventing any single system from having full access to an
authentication token that does not belong to. This helps us maintain the
security principles of **least privilege** and **non-repudiation**.

The important thing to note here is that regardless of what systems are
considered as _Trusted Entities_, the same pattern applies.

For example:

- With Chef, you might use the [Vault Ruby Gem](https://github.com/hashicorp/vault-ruby)
  for simplified interaction with Vault APIs
- Terraform provides a Vault provider: [Provider: Vault - Terraform by HashiCorp](https://www.terraform.io/docs/providers/vault/index.html)
- For Jenkins, you might use the Vault CLI or APIs directly, as described here:
  [Reading Vault Secrets in your Jenkins pipeline](http://nicolas.corrarello.com/general/vault/security/ci/2017/04/23/Reading-Vault-Secrets-in-your-Jenkins-pipeline.html)


## Prerequisites

This guide assumes that you are proficient enough to perform basic Terraform
tasks. If you are not familiar with Terraform, refer to the [online
documentation](https://www.terraform.io/intro/getting-started/install.html).

The following AWS resources are required to perform this demo:

- An [Amazon S3 bucket](https://docs.aws.amazon.com/AmazonS3/latest/gsg/CreatingABucket.html)
- An [IAM user credential with administrator permissions](https://docs.aws.amazon.com/IAM/latest/UserGuide/id_users_change-permissions.html)
(to be able to create additional IAM policies and instance profiles)

### Download demo assets

Clone or download the demo assets from the
[hashicorp/vault-guides](https://github.com/hashicorp/vault-guides/tree/master/identity/vault-chef-approle)
GitHub repository to perform the steps described in this guide.

The following assets can be found in the repository:

- Chef cookbook (**`/chef/cookbooks`**): A sample cookbook with a recipe that installs NGINX
and demonstrates Vault Ruby Gem functionality used to interact with Vault APIs.
- Terraform configurations (**`/terraform-aws`**):
    - **`/terraform-aws/mgmt-node`**: Configuration to set up a management
    server running both Vault and Chef Server, for demo purposes.
    - **`/terraform-aws/chef-node`**: Configuration to set up a Chef node and
    bootstrap it with the Chef Server, passing in Vault's AppRole RoleID and the
    appropriate Chef run-list.
- Vault configuration (**`/scripts`**): Data scripts used to configure the
appropriate mounts and policies in Vault for this demo.


## Steps

The scenario in this guide uses Terraform and Chef as trusted entities to
deliver `RoleID` and `SecretID`.

![AppRole auth method workflow](/img/vault-approle-tf-chef.png)

For the simplicity of the demonstration, both Vault and Chef are installed on
the same node. Terraform provisions the node which contains the `RoleID` as an
environment variable. Chef pulls the `SecretID` from Vault.


Provisioning for this demo happens in 2 phases:

- [Phase 1 - Provision our Vault plus Chef Server](#phase1)
    - [Step 1: Provision the Vault and Chef Server](#step-1-provision-the-vault-and-chef-server)
    - [Step 2: Initialize and Unseal Vault](#step-2-initialize-and-unseal-vault)
    - [Step 3: AppRole Setup](#step-3-approle-setup)
    - [Step 4: Configure Tokens for Terraform and Chef](#step-4-configure-tokens-for-terraform-and-chef)
    - [Step 5: Save the Token in a Chef Data Bag](#step-5-save-the-token-in-a-chef-data-bag)
    - [Step 6: Write Secrets](#step-6-write-secrets)
- [Phase 2 - Provision our Chef Node to Show AppRole Login](#phase2)


## <a name="phase1"></a>Phase 1: Provision our Vault & Chef Server

### Step 1: Provision the Vault and Chef Server

This provides a quick and simple Vault and Chef Server configuration to help you
get started.  

**NOTE:** This is done for demonstration purpose and **NOT a recommended
practice** for production.

In this phase, you use Terraform to spin up a server (and associated AWS
resources) with both Vault and Chef Server installed. Once this server is up and
running, you'll complete the appropriate configuration steps in Vault to set up
our AppRole and tokens for use in the demo.

~> If using _Terraform Enterprise_, [create a
Workspace](https://www.terraform.io/docs/enterprise/getting-started/workspaces.html)
for this repo and set the appropriate Terraform/Environment variables using the
`terraform.tfvars.example` file as a reference. Follow the instructions in the
documentation to perform the appropriate setup in Terraform Enterprise.

#### Using Terraform Open Source:

**Task 1:** Change the working directory (`cd`) to
`identity/vault-chef-approle/terraform-aws/mgmt-node`.

```shell
.
├── main.tf
├── outputs.tf
├── templates
│   └── userdata-mgmt-node.tpl
├── terraform.tfvars.example
└── variables.tf
```

**Task 2:** Update the `terraform.tfvars.example` file to match your account and
rename it to `terraform.tfvars`.

At minimum, replace the following variable with appropriate values:

- **`s3_bucket_name`**
- **`vpc_id`**
- **`subnet_id`**
- **`key_name`**
- **`ec2_pem`**

> NOTE: If your VPC, subnet and EC2 key pair were created on a region other than
`us-east-1`, be sure to set the **`aws_region`** value to match your chosen region.

**Task 3:** Perform a `terraform init` to pull down the necessary provider resources.
Then `terraform plan` to verify your changes and the resources that will be
created. If all looks good, then perform a `terraform apply` to provision the
resources. The Terraform output will display the public IP address to SSH into
your server.

```plaintext
$ terraform init
Initializing provider plugins...
...
Terraform has been successfully initialized!


$ terraform plan
...
Plan: 5 to add, 0 to change, 0 to destroy.


$ terraform apply
...
Apply complete! Resources: 5 added, 0 changed, 0 destroyed.

Outputs:
vault-public-ip = 192.0.2.0
```

The Terraform output will display the public IP address to SSH into
your server.

For example:

```plaintext
$ ssh -i "/path/to/EC2/private_key.pem" ubuntu@192.0.2.0
```

**Task 4:** Initial setup of the Chef server takes several minutes. Once you can
SSH into your mgmt server, run `tail -f /var/log/tf-user-data.log` to see when
the initial configuration is complete.

```plaintext
$ tail -f /var/log/tf-user-data.log
```

When you see the following message, the initial setup is complete.

```plaintext
+ echo '2018/03/27 21:53:06 /var/lib/cloud/instance/scripts/part-001: Complete'
```

You can find the following subfolders in
your home directory:

- **`/home/ubuntu/vault-chef-approle-demo`**: root of our project
- **`/home/ubuntu/vault-chef-approle-demo/chef`**: root of our Chef app; this is
where our `knife` configuration is located (`.chef/knife.rb`)
- **`/home/ubuntu/vault-chef-approle-demo/scripts`**: there's a
`vault-approle-setup.sh` script located here to help automate the setup of
Vault, or you can follow along in the rest of this README to configure Vault
manually

### Step 2: Initialize and Unseal Vault

Before moving on, set your working environment variables in your mgmt server:

```plaintext
$ export VAULT_ADDR=http://127.0.0.1:8200
$ export VAULT_SKIP_VERIFY=true
```

Before you can do anything in Vault, you need to initialize and unseal it.
Perform ***one*** of the following:

- **Option 1:** Run the `/home/ubuntu/demo_setup.sh` script to get up and running, and proceed to
[Phase 2 - Provision our Chef Node to Show AppRole Login](#phase2).
- **Option 2:** Continue onto [Step 3: AppRole Setup](#step-3-approle-setup) to
set up the demo environment ***manually***.


### Step 3: AppRole Setup

First, initialize and unseal the Vault server using a shortcut.

~> This is a convenient shortcut for demo. **_DO NOT DO THIS IN PRODUCTION!!!_**

Refer to the [online documentation for initializing and unsealing](/intro/getting-started/deploy.html#initializing-the-vault) Vault for more details.

```shell
# Initialize the Vault server and write out the unseal keys and root token into files
$ curl --silent
       --request PUT \
       --data '{"secret_shares": 1, "secret_threshold": 1}' \
       ${VAULT_ADDR}/v1/sys/init | tee \
       >(jq -r .root_token > /home/ubuntu/vault-chef-approle-demo/root-token) \
       >(jq -r .keys[0] > /home/ubuntu/vault-chef-approle-demo/unseal-key)

# Unseal vault
$ vault operator unseal $(cat /home/ubuntu/vault-chef-approle-demo/unseal-key)

# Set the root token to VAULT_TOKEN env var
$ export VAULT_TOKEN=$(cat /home/ubuntu/vault-chef-approle-demo/root-token)
```

In the next few steps, you will create a number of policies and tokens within
Vault. Below is a table that summarizes them:

| Policy         | Description | Token Attachment     |
|--------------------|-------------|------------------------|
| `app-1-secret-read` | Sets the policy for the final token that will be delivered via the AppRole login | None. This will be delivered to the client upon AppRole login |
| `app-1-approle-roleid-get` | Sets the policy for the token that you'll give to Terraform to deliver the RoleID (only) | `roleid-token` |
| `terraform-token-create`   | The Terraform Vault provider doesn't use the token supplied to it directly. This is to prevent the token from being exposed in Terraform's state file. Instead, the Token given to Terraform needs to have the capability to create child tokens with short TTLs. See [here] (https://www.terraform.io/docs/providers/vault/index.html#token) for more info | `roleid-token` |
| `app-1-approle-secretid-create` | Sets the policy for the token that you'll store in the Chef Data Bag. This will only be able to pull our AppRole's SecretID | `secretid-token` |



These setups only need to be performed upon initial creation of an AppRole, and
would typically be done by a Vault administrator.

Now that you have your Vault server unsealed, you can begin to set up necessary
policies, AppRole auth method, and tokens.

#### Task 1: Set up our AppRole policy

This is the policy that will be attached to _secret zero_ which you are
delivering to our application (**app-1**).

**API call using cURL**

```bash
# Policy to apply to AppRole token
$ tee app-1-secret-read.json <<EOF
{"policy":"path \"secret/app-1\" {capabilities = [\"read\", \"list\"]}"}
EOF

# Create the app-1-secret-read policy in Vault
$ curl --silent \
       --location \
       --header "X-Vault-Token: $VAULT_TOKEN" \
       --request PUT \
       --data @app-1-secret-read.json \
       $VAULT_ADDR/v1/sys/policy/app-1-secret-read
```

<br>
**CLI command**

```bash
# Policy to apply to AppRole token
$ tee app-1-secret-read.hcl <<EOF
path "secret/app-1" {
  capabilities = ["read", "list"]
}
EOF

# Create the app-1-secret-read policy in Vault
$ vault policy write app-1-secret-read app-1-secret-read.hcl
```


#### Task 2: Enable the AppRole authentication method

**API call using cURL**

```bash
# Payload for invoking sys/auth API endpoint
$ tee approle.json <<EOF
{
  "type": "approle",
  "description": "Demo AppRole auth method"
}
EOF

# Enable AppRole auth backend
$ curl --silent \
       --location \
       --header "X-Vault-Token: $VAULT_TOKEN" \
       --request POST \
       --data @approle.json \
       $VAULT_ADDR/v1/sys/auth/approle
```
<br>
**CLI command**

```plaintext
$ vault auth enable -description="Demo AppRole auth method" approle
```

#### Task 3: Configure the AppRole

Now, you are going to create an AppRole role named, **app-1**.

**API call using cURL**

```bash
# Payload containing AppRole auth method configuration
# TTL is set to 10 minutes, and Max TTL to be 30 minutes
$ tee app-1-approle-role.json <<EOF
{
    "role_name": "app-1",
    "bind_secret_id": true,
    "secret_id_ttl": "10m",
    "secret_id_num_uses": "1",
    "token_ttl": "10m",
    "token_max_ttl": "30m",
    "period": 0,
    "policies": [
        "app-1-secret-read"
    ]
}
EOF

# AppRole backend configuration
$ curl --silent \
       --location \
       --header "X-Vault-Token: $VAULT_TOKEN" \
       --request POST \
       --data @app-1-approle-role.json \
       $VAULT_ADDR/v1/auth/approle/role/app-1
```
<br>
**CLI command**

```bash
# TTL is set to 10 minutes, and Max TTL to be 30 minutes
$ vault write auth/approle/role/app-1 policies="app-1-secret-read" token_ttl="10m" token_max_ttl="30m"
```


### Step 4: Configure Tokens for Terraform and Chef

Now, you're ready to configure the policies and tokens to Terraform and Chef to
interact with Vault. Remember, the point here is that you are giving each system
a _limited_ token that is only able to pull either the `RoleID` or `SecretID`,
_but not both_.

![AppRole auth method workflow](/img/vault-approle-tf-chef-2.png)

#### Task 1: Create a policy and token for Terraform
Create a token with appropriate policies allowing Terraform to pull
the `RoleID` from Vault:

**API call using cURL**

```bash
# Policy file granting to retrieve RoleID from Vault
$ tee app-1-approle-roleid-get.hcl <<EOF
{"policy":"path \"auth/approle/role/app-1/role-id\" {capabilities = [\"read\"]}"}
EOF

# Create the app-1-approle-roleid-get policy in Vault
$ curl --silent \
       --location \
       --header "X-Vault-Token: $VAULT_TOKEN" \
       --request PUT \
       --data @app-1-approle-roleid-get.hcl \
       $VAULT_ADDR/v1/sys/policy/app-1-approle-roleid-get

# For Terraform
# See: https://www.terraform.io/docs/providers/vault/index.html#token
# Policy granting to create tokens required by Terraform
$ tee terraform-token-create.hcl <<EOF
{"policy":"path \"/auth/token/create\" {capabilities = [\"update\"]}"}
EOF

# Create the app-1-approle-roleid-get policy in Vault
$ curl --silent \
       --location \
       --header "X-Vault-Token: $VAULT_TOKEN" \
       --request PUT \
       --data @terraform-token-create.hcl \
       $VAULT_ADDR/v1/sys/policy/terraform-token-create

# Payload to configure token for Terraform to pull RoleID
$ tee roleid-token-config.json <<EOF
{
  "policies": [
    "app-1-approle-roleid-get",
    "terraform-token-create"
  ],
  "meta": {
    "user": "terraform-demo"
  },
  "ttl": "720h",
  "renewable": true
}
EOF

# Get token and save it in roleid-token.json
$ curl --silent \
       --location \
       --header "X-Vault-Token: $VAULT_TOKEN" \
       --request POST \
       --data @roleid-token-config.json \
       $VAULT_ADDR/v1/auth/token/create > roleid-token.json
```

The token and associated metadata will be written out to the file
`roleid-token.json`. The `client_token` value is what you'll give to Terraform.
The file should look similar to the following:

```plaintext
$ cat roleid-token.json | jq
{
  "request_id": "2e1d05eb-988d-4cf7-7b6a-d2668de31536",
  "lease_id": "",
  "renewable": false,
  "lease_duration": 0,
  "data": null,
  "wrap_info": null,
  "warnings": null,
  "auth": {
    "client_token": "6a7ad093-42ab-885e-3d67-6d51a5583da6",
    "accessor": "f6170506-ee0f-5a59-8478-e0aac2d3259f",
    "policies": [
      "app-1-approle-roleid-get",
      "default",
      "terraform-token-create"
    ],
    "metadata": {
      "user": "terraform-demo"
    },
    "lease_duration": 2592000,
    "renewable": true,
    "entity_id": ""
  }
}
```
<br>
**CLI command**

```bash
# Policy file granting to retrieve RoleID from Vault
$ tee app-1-approle-roleid-get.hcl <<EOF
path "auth/approle/role/app-1/role-id" {
  capabilities = [ "read" ]
}
EOF

# Create the app-1-approle-roleid-get policy in Vault
$ vault policy write app-1-approle-roleid-get app-1-approle-roleid-get.hcl

# For Terraform
# See: https://www.terraform.io/docs/providers/vault/index.html#token
# Policy granting to create tokens required by Terraform
$ tee terraform-token-create.hcl <<EOF
path "auth/token/create" {
  capabilities = [ "update" ]
}
EOF

# Create the app-1-approle-roleid-get policy in Vault
$ vault policy write terraform-token-create terraform-token-create.hcl

# Get token and save it in roleid-token.txt
$ vault token create -policy="app-1-approle-roleid-get" -policy="terraform-token-create" \
      -metadata="user"="terraform-user" > roleid-token.txt
```

The token and associated metadata will be written out to the file
`roleid-token.txt`. The `token` value is what you'll give to Terraform.
The file should look similar to the following:

```plaintext
$ cat roleid-token.txt
Key                Value
---                -----
token              2600aeda-6385-c163-7171-543b1e1fabcf
token_accessor     6ef835e3-4948-8c61-1e89-3625ca31fd84
token_duration     768h
token_renewable    true
token_policies     [app-1-approle-roleid-get default terraform-token-create]
token_meta_user    terraform-demo
```

#### Task 2: Create a policy and token for Chef
Create a token with appropriate policies allowing Chef to pull the `SecretID`
from Vault:

**API call using cURL**

```bash
# Policy file granting to retrieve SecretID
$ tee app-1-approle-secretid-create.hcl <<EOF
{"policy":"path \"auth/approle/role/app-1/secret-id\" {capabilities = [\"update\"]}"}
EOF

# Create the app-1-approle-secretid-create policy in Vault
$ curl --silent \
    --location \
    --header "X-Vault-Token: $VAULT_TOKEN" \
    --request PUT \
    --data @app-1-approle-secretid-create.hcl \
    $VAULT_ADDR/v1/sys/policy/app-1-approle-secretid-create

# Payload to invoke auth/token/create endpoint
$ tee secretid-token-config.json <<EOF
{
  "policies": [
    "app-1-approle-secretid-create"
  ],
  "meta": {
    "user": "chef-demo"
  },
  "ttl": "720h",
  "renewable": true
}
EOF

# Get token for Chef to get SecretID from Vault and store it in secretid-token.json
$ curl --silent \
       --location \
       --header "X-Vault-Token: $VAULT_TOKEN" \
       --request POST \
       --data @secretid-token-config.json \
       $VAULT_ADDR/v1/auth/token/create > secretid-token.json
```

The resulting file should look like this:

```plaintext
$ cat secretid-token.json | jq
{
  "request_id": "6f6ad8a1-fedb-b838-60ce-87999f01aff6",
  "lease_id": "",
  "renewable": false,
  "lease_duration": 0,
  "data": null,
  "wrap_info": null,
  "warnings": null,
  "auth": {
    "client_token": "cdfdb7a0-d7a6-3769-927d-0ace297726ea",
    "accessor": "88e8aaca-1584-4881-3368-d9cb5cd7ddae",
    "policies": [
      "app-1-approle-secretid-create",
      "default"
    ],
    "metadata": {
      "user": "chef-demo"
    },
    "lease_duration": 2592000,
    "renewable": true,
    "entity_id": ""
  }
}
```

<br>
**CLI command**

```bash
# Policy file granting to retrieve SecretID
$ tee app-1-approle-secretid-create.hcl <<EOF
path "auth/approle/role/app-1/secret-id" {
  capabilities = [ "update" ]
}
EOF

# Create the app-1-approle-secretid-create policy in Vault
$ vault policy write app-1-approle-secretid-create app-1-approle-secretid-create.hcl

# Get token for Chef to get SecretID from Vault and store it in secretid-token.txt
$ vault token create -policy="app-1-approle-secretid-create" \
      -metadata="user"="chef-demo" > secretid-token.txt
```

The resulting file should look like this:

```plaintext
$ cat secretid-token.txt
Key                Value
---                -----
token              20d69183-59cb-c953-6dea-34f5f1bbe5f7
token_accessor     a17e7a43-c14a-b96a-0014-9149d218e74a
token_duration     768h
token_renewable    true
token_policies     [app-1-approle-secretid-create default]
token_meta_user    chef-demo
```



### Step 5: Save the Token in a Chef Data Bag

At this point, you have a client token generated for Terraform and another for
Chef server to log into Vault. For the sake of simplicity, you can put the
Chef's client token (`secretid-token.json`) in a [Data
Bag](https://docs.chef.io/data_bags.html) which is fine because this token can
***only*** retrieve `SecretID` from Vault which is not much of a use without a
corresponding `RoleID`.

Now, create a Chef Data Bag and put the `SecretID` token (`secretid-token.json`)
along with the rest of its metadata.

```bash
$ cd /home/ubuntu/vault-chef-approle-demo/chef/

# Use the path for where you created this file in the previous step
# You're just adding an 'id' field to the file as that's a required field for data bags
$ cat /home/ubuntu/secretid-token.json | jq --arg id approle-secretid-token '. + {id: $id}' > secretid-token.json

$ knife data bag create secretid-token

$ knife data bag from file secretid-token secretid-token.json

$ knife data bag list

$ knife data bag show secretid-token

$ knife data bag show secretid-token approle-secretid-token
```

The last step should show the following output:

```plaintext
$ knife data bag show secretid-token approle-secretid-token
WARNING: Unencrypted data bag detected, ignoring any provided secret options.
auth:
  accessor:       88e8aaca-1584-4881-3368-d9cb5cd7ddae
  client_token:   cdfdb7a0-d7a6-3769-927d-0ace297726ea
  entity_id:
  lease_duration: 2592000
  metadata:
  policies:
    app-1-approle-secretid-create
    default
  renewable:      true
data:
id:             approle-secretid-token
lease_duration: 0
lease_id:
renewable:      false
request_id:     6f6ad8a1-fedb-b838-60ce-87999f01aff6
warnings:
wrap_info:
```

### Step 6: Write Secrets

Let's write some test data in the `secret/app-1` path so that the target app
will have some secret to retrieve from Vault at a later step.

**API call using cURL**

```bash
# Write some demo secrets
$ tee demo-secrets.json <<'EOF'
{
  "username": "app-1-user",
  "password": "$up3r$3cr3t!"
}
EOF

$ curl --silent \
       --location \
       --header "X-Vault-Token: $VAULT_TOKEN" \
       --request POST \
       --data @demo-secrets.json \
       $VAULT_ADDR/v1/secret/app-1

# Verify that you can read back the data:
$ curl --silent \
       --location \
       --header "X-Vault-Token: $VAULT_TOKEN" \
       --request GET \
       $VAULT_ADDR/v1/secret/app-1 | jq
{
  "request_id": "1f73c7ee-27fa-bad0-9c77-b330eef1ea88",
  "lease_id": "",
  "renewable": false,
  "lease_duration": 2764800,
  "data": {
    "password": "$up3r$3cr3t!",
    "username": "app-1-user"
  },
  "wrap_info": null,
  "warnings": null,
  "auth": null
}
```
<br>
**CLI command**

```bash
# Write some demo secrets
$ vault write secret/app-1 username="app-1-user" password="\$up3r\$3cr3t!"

# Verify that you can read back the data:
$ vault read secret/app-1
Key                 Value
---                 -----
refresh_interval    768h
password            $up3r$3cr3t!
username            app-1-user
```


-> At this point, just about all the pieces are in place. Remember, these setup
steps will only need to be performed upon initial creation of an AppRole, and
would typically be done by a Vault administrator.



## <a name="phase2"></a>Phase 2: Provision our Chef Node to Show AppRole Login

To complete the demo, run the **`chef-node`** Terraform configuration to see how
everything talks to each other.

#### Task 1: Change the working directory

Open another terminal on your host machine (**not** the `mgmt-node`)
and `cd` into the `identity/vault-chef-approle/terraform-aws/chef-node`
directory:

```plaintext
$ cd identity/vault-chef-approle/terraform-aws/chef-node
```

#### Task 2: Update terraform.tfvars.example

Replace the variable values in `terraform.tfvars.example` to match your
environment and save it as `terraform.tfvars` like you have done at [Step 1](#step-1-provision-the-vault-and-chef-server).

Note the following:

* Update the **`vault_address`** and **`chef_server_address`** variables with
the IP address of our `mgmt-node` from above.
* Update the **`vault_token`** variable with the `RoleID` token from **Task 1**
in [Step 4](#step-4-configure-tokens-for-terraform-and-chef).  
  - If you ran the `demo-setup.sh` script (_Option 1_), retrieve the
  `client_token` in the `/home/ubuntu/vault-chef-approle-demo/roleid-token.json`
  file:

```plaintext
$ cat ~/vault-chef-approle-demo/roleid-token.json | jq ".auth.client_token"
```


#### Task 3: Run Terraform
Perform a `terraform init` to pull down the necessary provider
resources. Then `terraform plan` to verify your changes and the resources that
will be created. If all looks good, then perform a `terraform apply` to
provision the resources.

The Terraform output will display the public IP address to SSH into your
server.

> **NOTE:** If the `terraform apply` fails with "`io: read/write on closed pipe`"
error, this is a [known
issue](https://github.com/hashicorp/terraform/issues/17638) with Terraform
0.11.4 and 0.11.5.  Please try again with another Terraform version.

At this point, Terraform will perform the following actions:

- Pull a `RoleID` from our Vault server
- Provision an AWS instance
- Write the `RoleID` to the AWS instance as an environment variable
- Run the Chef provisioner to bootstrap the AWS instance with our Chef Server
- Run our Chef recipe which will install NGINX, perform our AppRole login, get
our secrets, and output them to our `index.html` file

![AppRole auth method workflow](/img/vault-approle-tf-chef-3.png)

The Chef recipe can be found at
`identity/vault-chef-approle/chef/cookbooks/vault_chef_approle_demo/recipes/default.rb`.

```shell
...
# Configure address for Vault Gem
Vault.address = ENV['VAULT_ADDR']

# Get AppRole RoleID from our environment variables (delivered via Terraform)
var_role_id = ENV['APPROLE_ROLEID']

# Get Vault token from data bag (used to retrieve the SecretID)
vault_token_data = data_bag_item('secretid-token', 'approle-secretid-token')

# Set Vault token (used to retrieve the SecretID)
Vault.token = vault_token_data['auth']['client_token']

# Get AppRole SecretID from Vault
var_secret_id = Vault.approle.create_secret_id('app-1').data[:secret_id]
...
```

#### Task 4: Verification
Once Terraform completes the `apply` operation, it will output the public IP
address of our new server. You can plug that IP address into a browser to see
the output. It should look similar to the following:

```plaintext
Role ID:
f6286b97-246e-9fb4-4d9f-0c9465451851

Secret ID:
72f4b60c-26d0-d947-5026-153943174831

AppRole Token:
d11d81e4-0ba1-fefc-03f8-e5f06793b60d

Read Our Secrets:
{:password=>"$up3r$3cr3t!", :username=>"app-1-user"}
```

## Additional References

The following is a curated list of webinars, blogs and GitHub repositories that
add additional context to fill out the concepts discussed in the webinar and
demonstrated in the code:

- [Managing Secrets in a Container Environment by Jeff Mitchell](https://www.youtube.com/watch?v=skENC9aXgco)
- [Using HashiCorp's Vault with Chef written by Seth Vargo](https://www.hashicorp.com/blog/using-hashicorps-vault-with-chef)
- [Manage Secrets with Chef and HashiCorps Vault by Seth Vargo & JJ Asghar](https://blog.chef.io/2016/12/12/manage-secrets-with-chef-and-hashicorps-vault/)
    - [Associated GitHub repository](https://github.com/sethvargo/vault-chef-webinar)
- [Vault AppRole Authentication written by Alan Thatcher](http://blog.alanthatcher.io/vault-approle-authentication/)
- [Integrating Chef and HashiCorp Vault written by Alan Thatcher](http://blog.alanthatcher.io/integrating-chef-and-hashicorp-vault/)
- [Vault Ruby Client](https://github.com/hashicorp/vault-ruby)


## Next Steps

Watch the video recording of the [Delivering Secret Zero: Vault AppRole with
Terraform and
Chef](https://www.hashicorp.com/resources/delivering-secret-zero-vault-approle-terraform-chef)
webinar which talks about the usage of AppRole with Terraform and Chef as its
trusted entities.
