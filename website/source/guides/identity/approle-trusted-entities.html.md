---
layout: "guides"
page_title: "AppRole With Terraform & Chef - Guides"
sidebar_current: "guides-identity-approle-tf-chef"
description: |-
  This guide discusses the concepts necessary to help users
  understand Vault's AppRole authentication pattern and how to use it to
  securely introduce a Vault authentication token to a target server,
  application, container, etc.x
---

# Vault AppRole with Terraform and Chef Demo

In the [AppRole Pull
Authentication](/guides/identity/authentication.html#advanced-feature) guide,
the question of how best to deliver the Role ID and Secret ID were brought up,
and the role of trusted entities (Terraform, Chef, Nomad, Kubernetes, etc.) was
mentioned.

![AppRole auth method workflow](/assets/images/vault-approle-workflow2.png)

This _intermediate_ Vault guide aims to provide a **simple**, **end-to-end**
example of how to use Vault's [AppRole authentication
method](/docs/auth/approle.html), along with Terraform and Chef, to address the
challenge of **_secure introduction_** of an initial token to a target server or
application.

The purpose of this guide is to provide a step-by-step instruction to reproduce
the working implementation of the demo introduced in the [Delivering Secret
Zero: Vault AppRole with Terraform and
Chef](https://www.hashicorp.com/resources/delivering-secret-zero-vault-approle-terraform-chef)
webinar.

[![YouTube](/assets/images/vault-approle-youtube.png)](https://youtu.be/OIcIzFWjThM)

-> **NOTE:** This is a proof of concept and **NOT SUITABLE FOR PRODUCTION USE**.


## Reference Material

- [AppRole Auth Method](/docs/auth/approle.html)
- [Authenticating Applications with HashiCorp Vault AppRole](https://www.hashicorp.com/blog/authenticating-applications-with-vault-approle)
- [Delivering Secret Zero: Vault AppRole with Terraform and
Chef](https://www.hashicorp.com/resources/delivering-secret-zero-vault-approle-terraform-chef)

## Estimated Time to Complete

20 minutes

## challenge

The goal of the AppRole authentication method is to provide a mechanism for the
secure introduction of secrets to target systems (servers, applications,
containers, etc.).

What systems within our environment do we trust to handle or deliver secrets to
our target system?


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

The following resources are required to perform this demo:

- An Amazon S3 bucket
- An IAM user with administrator permissions (to be able to create additional
  IAM policies, instance profiles)

### Demo assets

Clone or download the demo assets from the
[hashicorp/vault-guides](https://github.com/hashicorp/vault-guides/tree/master/identity/vault-chef-approle)
GitHub repository to perform the steps described in this guide.

The following assets can be found in the repository:

- Chef cookbook (**`/chef`**): A sample cookbook with a recipe that installs NGINX
and demonstrates Vault Ruby Gem functionality used to interact with Vault APIs.
- Terraform configurations (**`/terraform-aws`**):
    - **`/terraform-aws/mgmt-node`**: Configuration to set up a management
    server running both Vault and Chef Server, for demo purposes.
    - **`/terraform-aws/chef-node`**: Configuration to set up a Chef node and
    bootstrap it with the Chef Server, passing in Vault's AppRole RoleID and the
    appropriate Chef run-list.
- Vault configuration (**`/vault`**): Data scripts used to configure the
appropriate mounts and policies in Vault for this demo.


## Steps

Provisioning for this project happens in 2 phases:

![AppRole auth method workflow](/assets/images/vault-approle-tf-chef.png)

- [Phase 1 - Provision our Vault plus Chef Server](#phase1)
    - [Step 1: Provision the Vault and Chef Server](#step-1-provision-the-vault--chef-server)
    - [Step 2: Initialize and Unseal Vault](#step-2-initialize-and-unseal-vault)
    - [Step 3: AppRole Setup](#step-3-approle-setup)
    - [Step 4: Configure Tokens for Terraform and Chef](#step-4-configure-tokens-for-terraform-and-chef)
    - [Step 5: Put the SecretID Token Into a Chef Data Bag](#step-5-put-the-secretid-token-into-a-chef-data-bag)
    - [Step 6: Write Some Secrets](#step-6-write-some-secrets)
- [Phase 2 - Provision our Chef Node to Show AppRole Login](#phase2)


## <a name="phase1"></a>Phase 1: Provision our Vault & Chef Server

### Step 1: Provision the Vault and Chef Server

This provides a quick and simple Vault and Chef Server configuration to help you
get started.  

-> NOTE: Demo purpose only and **_NOT SUITABLE FOR PRODUCTION USE!!_**

In this phase, we use Terraform to spin up a server (and associated AWS
resources) with both Vault and Chef Server installed. Once this server is up and
running, we'll complete the appropriate configuration steps in Vault to set up
our AppRole and tokens for use in the demo.

~> If using _Terraform Enterprise_, [create a
Workspace](https://www.terraform.io/docs/enterprise/getting-started/workspaces.html)
for this repo and set the appropriate Terraform/Environment variables using the
`terraform.tfvars.example` file as a reference. Follow the instructions in the
documentation to perform the appropriate setup in Terraform Enterprise.

#### Using Terraform Open Source:

1. Change the working directory (`cd`) to
`identity/vault-chef-approle/terraform-aws/mgmt-node`.

2. Make sure to update the `terraform.tfvars.example` file accordingly and
rename it to `terraform.tfvars`.

3. Perform a `terraform init` to pull down the necessary provider resources.
Then `terraform plan` to verify your changes and the resources that will be
created. If all looks good, then perform a `terraform apply` to provision the
resources. - The Terraform output will display the public IP address to SSH into
your server.

4. Initial setup of the Chef Server takes several minutes. Once you can SSH into
your mgmt server, run `tail -f /var/log/tf-user-data.log` to see when the
initial configuration is complete. When you see
`.../var/lib/cloud/instance/scripts/part-001: Complete`, you'll know that
initial setup is complete.

Once the user-data script has completed, you'll see the following subfolders in
your home directory: - `/home/ubuntu/vault-chef-approle-demo`: root of our
project

- `/home/ubuntu/vault-chef-approle-demo/chef`: root of our Chef app; this is
where our `knife` configuration is located [`.chef/knife.rb`]

- `/home/ubuntu/vault-chef-approle-demo/scripts`: there's a
`vault-approle-setup.sh` script located here to help automate the setup of
Vault, or you can follow along in the rest of this README to configure Vault
manually

### Step 2: Initialize and Unseal Vault

Before moving on, let's set our working environment variables:

```
export VAULT_ADDR=http://127.0.0.1:8200
export VAULT_SKIP_VERIFY=true
```

To set up Vault manually, continue below. Otherwise, run the
`/home/ubuntu/demo_setup.sh` script to get up and running, and skip to "Phase 2
[Provision our Chef Node to Show AppRole Login]"

1. Before we can do anything in Vault, we need to initialize and unseal it.
We'll take a bit of a shortcut here... **_DON'T DO THIS IN PRODUCTION!!!_**

```
curl \
--silent \
--request PUT \
--data '{"secret_shares": 1, "secret_threshold": 1}' \
${VAULT_ADDR}/v1/sys/init | tee \
>(jq -r .root_token > /home/ubuntu/vault-chef-approle-demo/root-token) \
>(jq -r .keys[0] > /home/ubuntu/vault-chef-approle-demo/unseal-key)

vault operator unseal $(cat /home/ubuntu/vault-chef-approle-demo/unseal-key)

export VAULT_TOKEN=$(cat /home/ubuntu/vault-chef-approle-demo/root-token)
```

### Step 3: AppRole Setup

In the next few steps, we'll be creating a number of policies and tokens within
Vault. Below is a table that summarizes them:

| Policy         | Description | Token Attachment     |
|--------------------|-------------|------------------------|
| `app-1-secret-read` | Sets the policy for the final token that will be delivered via the AppRole login | None. This will be delivered to the client upon AppRole login |
| `app-1-approle-roleid-get` | Sets the policy for the token that we'll give to Terraform to deliver the RoleID (only) | `roleid-token` |
| `terraform-token-create`   | The Terraform Vault provider doesn't use the token supplied to it directly. This is to prevent the token from being exposed in Terraform's state file. Instead, the Token given to Terraform needs to have the capability to create child tokens with short TTLs. See [here] (https://www.terraform.io/docs/providers/vault/index.html#token) for more info | `roleid-token` |
| `app-1-approle-secretid-create` | Sets the policy for the token that we'll store in the Chef Data Bag. This will only be able to pull our AppRole's SecretID | `secretid-token` |



These setup steps will only need to be performed upon initial creation of an
AppRole, and would typically be done by a Vault administrator.

Now that we have Vault unsealed, we can begin to set up our policies, AppRole
auth method, and tokens.

1. Set up our AppRole policy. This is the policy that will be attached to
_secret zero_ which we are delivering to our app:

```bash
# Policy to apply to AppRole token
tee app-1-secret-read.json <<EOF
{"policy":"path \"secret/app-1\" {capabilities = [\"read\", \"list\"]}"}
EOF

# Write the policy
curl \
    --silent \
    --location \
    --header "X-Vault-Token: $VAULT_TOKEN" \
    --request PUT \
    --data @app-1-secret-read.json \
    $VAULT_ADDR/v1/sys/policy/app-1-secret-read
```

2. Enable the AppRole authentication method:

```bash
# Enable AppRole auth backend
tee approle.json <<EOF
{
  "type": "approle",
  "description": "Demo AppRole auth backend"
}
EOF

curl \
    --silent \
    --location \
    --header "X-Vault-Token: $VAULT_TOKEN" \
    --request POST \
    --data @approle.json \
    $VAULT_ADDR/v1/sys/auth/approle
```

3. Configure the AppRole:

```bash
# AppRole backend configuration
tee app-1-approle-role.json <<EOF
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

curl \
    --silent \
    --location \
    --header "X-Vault-Token: $VAULT_TOKEN" \
    --request POST \
    --data @app-1-approle-role.json \
    $VAULT_ADDR/v1/auth/approle/role/app-1
```

### Step 4: Configure Tokens for Terraform and Chef

At this point, we're ready to configure the policies and tokens that we'll give
to Terraform and Chef. Remember, the point here is that we are giving each
system a _limited_ token that is only able to pull either the `RoleID` or
`SecretID`, _but not both_.

1. Create the token that we'll give to Terraform that will allow it to pull the
`RoleID` from Vault:

```bash
# Policy to get RoleID
tee app-1-approle-roleid-get.json <<EOF
{"policy":"path \"auth/approle/role/app-1/role-id\" {capabilities = [\"read\"]}"}
EOF

# Write the policy
curl \
    --silent \
    --location \
    --header "X-Vault-Token: $VAULT_TOKEN" \
    --request PUT \
    --data @app-1-approle-roleid-get.json \
    $VAULT_ADDR/v1/sys/policy/app-1-approle-roleid-get

# For Terraform
# See: https://www.terraform.io/docs/providers/vault/index.html#token
tee terraform-token-create.json <<EOF
{"policy":"path \"/auth/token/create\" {capabilities = [\"update\"]}"}
EOF

# Write the policy
curl \
    --silent \
    --location \
    --header "X-Vault-Token: $VAULT_TOKEN" \
    --request PUT \
    --data @terraform-token-create.json \
    $VAULT_ADDR/v1/sys/policy/terraform-token-create

# Configure token for RoleID
tee roleid-token-config.json <<EOF
{
  "policies": [
    "app-1-approle-roleid-get",
    "terraform-token-create"
  ],
  "metadata": {
    "user": "chef-demo"
  },
  "ttl": "720h",
  "renewable": true
}
EOF

# Get token
curl \
    --silent \
    --location \
    --header "X-Vault-Token: $VAULT_TOKEN" \
    --request POST \
    --data @roleid-token-config.json \
    $VAULT_ADDR/v1/auth/token/create > roleid-token.json
```

2. The token and associated metadata will be written out to the file
`roleid-token.json`. The `client_token` value is what we'll give to Terraform.
The file should look similar to the following:

```json
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
    "metadata": null,
    "lease_duration": 2592000,
    "renewable": true,
    "entity_id": ""
  }
}
```

3. Let's do the same for Chef, but for the `SecretID`:

```bash
# Policy to get SecretID
tee app-1-approle-secretid-create.json <<EOF
{"policy":"path \"auth/approle/role/app-1/secret-id\" {capabilities = [\"update\"]}"}
EOF

# Write the policy
curl \
    --silent \
    --location \
    --header "X-Vault-Token: $VAULT_TOKEN" \
    --request PUT \
    --data @app-1-approle-secretid-create.json \
    $VAULT_ADDR/v1/sys/policy/app-1-approle-secretid-create

# Configure token for SecretID
tee secretid-token-config.json <<EOF
{
  "policies": [
    "app-1-approle-secretid-create"
  ],
  "metadata": {
    "user": "chef-demo"
  },
  "ttl": "720h",
  "renewable": true
}
EOF

# Get token
curl \
    --silent \
    --location \
    --header "X-Vault-Token: $VAULT_TOKEN" \
    --request POST \
    --data @secretid-token-config.json \
    $VAULT_ADDR/v1/auth/token/create > secretid-token.json
```

4. Similarly to above, the file should look like this:

```json
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
    "metadata": null,
    "lease_duration": 2592000,
    "renewable": true,
    "entity_id": ""
  }
}
```

### Step 5: Put the SecretID Token Into a Chef Data Bag

At this point, we have a token to give Terraform (which we'll do in Phase 2) and
one to give to our Chef Server. For the sake of simplicity, we can just put it
in a Data Bag... and this is OK because, again, this token can _only_ retrieve
`SecretID`s which are useless without a corresponding `RoleID`.

1. Let's create our Chef Data Bag and dump the `SecretID` token in there, along
with the rest of the metadata, because... why not? :-)

```bash
cd /home/ubuntu/vault-chef-approle-demo/chef/

# Use the path for where you created this file in the previous step
# We're just adding an 'id' field to the file as that's a required field for data bags
cat /home/ubuntu/secretid-token.json | jq --arg id approle-secretid-token '. + {id: $id}' > secretid-token.json

knife data bag create secretid-token
knife data bag from file secretid-token secretid-token.json
knife data bag list
knife data bag show secretid-token
knife data bag show secretid-token approle-secretid-token
```

2. The last step should show the following output:

```
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

### Step 6: Write Some Secrets

1. Finally, let's write some dummy data to show that we can read "stuff" from Vault on
our target app:

```bash
# Write some demo secrets
tee demo-secrets.json <<'EOF'
{
  "username": "app-1-user",
  "password": "$up3r$3cr3t!"
}
EOF

curl \
    --silent \
    --location \
    --header "X-Vault-Token: $VAULT_TOKEN" \
    --request POST \
    --data @demo-secrets.json \
    $VAULT_ADDR/v1/secret/app-1
```

2. We can verify the data, just to be safe:

```bash
curl \
    --silent \
    --location \
    --header "X-Vault-Token: $VAULT_TOKEN" \
    --request GET \
    $VAULT_ADDR/v1/secret/app-1 | jq
```

And you should see the following:

```json
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

At this point, just about all the pieces are in place. Remember, these setup
steps will only need to be performed upon initial creation of an AppRole, and
would typically be done by a Vault administrator.

## <a name="phase2"></a>Phase 2: Provision our Chef Node to Show AppRole Login

To complete the demo, we'll now run our `chef-node` Terraform configuration to
see how everything talks to each other. First, some final setup...

1. Open another terminal window/tab (on your host machine, not the `mgmt-node`)
and `cd` into the `identity/vault-chef-approle/terraform-aws/chef-node`
directory.

2. Update the `terraform.tfvars.example` file accordingly and rename to
`terraform.tfvars`:

    * Update the `vault_address` and `chef_server_address` variables with the IP
    address of our `mgmt-node` from above.
    * Update the `vault_token` variable with the `RoleID` token from Step 4.2
    above.

3. Perform a `terraform init` to pull down the necessary provider resources. Then
`terraform plan` to verify your changes and the resources that will be created.
If all looks good, then perform a `terraform apply` to provision the resources.

    - The Terraform output will display the public IP address to SSH into your
    server.

At this point, Terraform will perform the following actions:

- Pull a `RoleID` from our Vault server
- Provision an AWS instance
- Write the `RoleID` to the AWS instance as an environment variable
- Run the Chef provisioner to bootstrap the AWS instance with our Chef Server
- Run our Chef recipe which will install Nginx, perform our AppRole login, get
our secrets, and output them to our `index.html` file

4. Once Terraform completes the `apply`, it will output the public IP address of
our new server. We can plug that IP address into a browser to see the output. It
should look similar to the following:

```
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

The following is a curated list of webinars/blogs/etc. that add additional
context to fill out the concepts discussed in the webinar and demonstrated in
the code:

- Webinar for this demo repo: [Delivering Secret Zero: Vault AppRole with Terraform and Chef](https://www.hashicorp.com/resources/delivering-secret-zero-vault-approle-terraform-chef)
- [Jeff Mitchell: Managing Secrets in a Container Environment](https://www.youtube.com/watch?v=skENC9aXgco)
- [Seth Vargo: Using HashiCorp's Vault with Chef](https://www.hashicorp.com/blog/using-hashicorps-vault-with-chef)
- [Seth Vargo & JJ Asghar: Manage Secrets with Chef and HashiCorps Vault](https://blog.chef.io/2016/12/12/manage-secrets-with-chef-and-hashicorps-vault/)
    - [Associated Github repo](https://github.com/sethvargo/vault-chef-webinar)
- [Alan Thatcher: Vault AppRole Authentication](http://blog.alanthatcher.io/vault-approle-authentication/)
- [Alan Thatcher: Integrating Chef and HashiCorp Vault](http://blog.alanthatcher.io/integrating-chef-and-hashicorp-vault/)
- [Vault Ruby Client](https://github.com/hashicorp/vault-ruby)
- https://github.com/hashicorp-guides/vault-approle-chef (README will eventually be merged with this document)



## Next steps

Watch the video recording of the [Delivering Secret Zero: Vault AppRole with
Terraform and
Chef](https://www.hashicorp.com/resources/delivering-secret-zero-vault-approle-terraform-chef)
webinar which talks about the usage of AppRole with Terraform and Chef as its trusted entities.
