---
layout: "guides"
page_title: "One-Time SSH Password - Guides"
sidebar_title: "One-Time SSH Password"
sidebar_current: "guides-secret-mgmt-ssh-otp"
description: |-
  The one-time SSH password secrets engine allows Vault to issue a one-time
  password (OTP) every time a client wants to SSH into a remote host using a
  helper command on the remote host to perform verification.
---

# Vault SSH Secrets Engine: One-Time SSH Password

In a distributed cloud environment, tenant and system is increasingly
important part of the online security. If an attacker gains access to your
virtual machines, they can get control of most running applications, local data
as well as its connected machines and systems.  

The Vault SSH secrets engine provides secure authentication and authorization
for access to machines via the SSH protocol. It supports [signed SSH
certificate](/docs/secrets/ssh/signed-ssh-certificates.html) and one-time SSH
password modes. This guide demonstrates the one-time SSH password mode.


## Reference Material

- [One-Time SSH Passwords](/docs/secrets/ssh/one-time-ssh-passwords.html)
- [SSH Secrets Engine (API)](/api/secret/ssh/index.html)
- [Vault SSH Helper](https://github.com/hashicorp/vault-ssh-helper)

## Estimated Time to Complete

10 minutes


## Personas

The end-to-end scenario described in this guide involves two personas:

- **`operations`** with privileged permissions to setup SSH secrets engine
- **`client`** trusted entity to request SSH OTP from Vault

## Challenge

By default, SSH servers use password authentication with optional public key
authentication. If any user on the system has a fairly weak password, this
allows an attacker to hijack the SSH connection.


## Solution

Vault can create a one-time password (OTP) for SSH authentication on a network
every time a client wants to SSH into a remote host using a helper command on
the remote host to perform verification.

![SSH OTP Workflow](/img/vault-ssh-otp-1.png)

An authenticated client requests an OTP from the Vault server. If the client is
authorized, Vault issues and returns an OTP. The client uses this OTP during the
SSH authentication to connect to the desired target host.

When the client establishes an SSH connection, the OTP is received by the
***Vault helper*** which validates the OTP with the Vault server. The Vault
server then deletes this OTP, ensuring that it is **only used once**.

~> Since the Vault server is contacted during SSH connection establishment, every
login attempt and the correlating Vault lease information is logged to the audit
secrets engine.


## Prerequisites

To perform the tasks described in this guide, you need to have a Vault
environment.  Refer to the [Getting
Started](/intro/getting-started/install.html) guide to install Vault.

### Policy requirements

-> **NOTE:** For the purpose of this guide, you can use **`root`** token to work
with Vault. However, it is recommended that root tokens are only used for
initial setup or in emergencies. As a best practice, use tokens with
appropriate set of policies based on your role in the organization.

To perform all tasks demonstrated in this guide, your policy must include the
following permissions:

```shell
# To view in Web UI
path "sys/mounts" {
  capabilities = [ "read", "update" ]
}

# To configure the SSH secrets engine
path "ssh/*" {
  capabilities = [ "create", "read", "update", "delete", "list" ]
}

# To enable secret engines
path "sys/mounts/*" {
  capabilities = [ "create", "read", "update", "delete" ]
}
```

If you are not familiar with policies, complete the
[policies](/guides/identity/policies.html) guide.


## Steps

You will perform the following:

1. [Install vault-ssh-helper](#step1)
1. [Setup the SSH Secrets Engine](#step2)
1. [Request an OTP](#step3)
1. [Establish an SSH session](#step4)


### <a name="step1"></a>Step 1: Install vault-ssh-helper
(**Persona:** operations)

The SSH secrets engine uses [Vault SSH
Helper](https://github.com/hashicorp/vault-ssh-helper) to verify the OTP used
during the SSH authentication. Therefore, the helper agent must be installed
onto every ***target host***.  

1. Download and install the latest version of `vault-ssh-helper` from
[releases.hashicorp.com](https://releases.hashicorp.com/vault-ssh-helper).

    **Example:**

    ```shell
    # Download the vault-ssh-helper
    $ wget https://releases.hashicorp.com/vault-ssh-helper/0.1.4/vault-ssh-helper_0.1.4_linux_amd64.zip

    # Unzip the vault-ssh-helper in /user/local/bin
    $ sudo unzip -q vault-ssh-helper_0.1.4_linux_amd64.zip -d /usr/local/bin

    # Make sure that vault-ssh-helper is executable
    $ sudo chmod 0755 /usr/local/bin/vault-ssh-helper

    # Set the usr and group of vault-ssh-helper to root
    $ sudo chown root:root /usr/local/bin/vault-ssh-helper
    ```

1. Create a Vault SSH Helper configuration file,
**`/etc/vault-ssh-helper.d/config.hcl`**.

    ```hcl
    vault_addr = "<VAULT_ADDRESS>"
    ssh_mount_point = "ssh"
    ca_cert = "/etc/vault-ssh-helper.d/vault.crt"
    tls_skip_verify = false
    allowed_roles = "*"
    ```

    Where the **`<VAULT_ADDRESS>`** is the address of the Vault server generating the OTP.

    **Example:**

    ```hcl
    vault_addr = "https://198.51.100.10:8200"
    ssh_mount_point = "ssh"
    ca_cert = "/etc/vault-ssh-helper.d/vault.crt"
    tls_skip_verify = false
    allowed_roles = "*"
    ```

    > Refer to the [documentation](https://github.com/hashicorp/vault-ssh-helper#vault-ssh-helper-configuration)
    for the entire list of configuration properties.


1. Modify the **`/etc/pam.d/sshd`** file as follows:

    ```shell
    # PAM configuration for the Secure Shell service

    # Standard Un*x authentication.
    #@include common-auth
    auth requisite pam_exec.so quiet expose_authtok log=/tmp/vaultssh.log /usr/local/bin/vault-ssh-helper -dev -config=/etc/vault-ssh-helper.d/config.hcl
    auth optional pam_unix.so not_set_pass use_first_pass nodelay

    ...
    ```

    **NOTE:** `common-auth` is the standard Linux authentication module which is
    commented out in favor of using our custom configuration.

    > Refer to the
    [documentation](https://github.com/hashicorp/vault-ssh-helper#pam-configuration)
    for details about these parameter settings.


1. Modify the **`/etc/ssh/sshd_config`** file.

    ```plaintext
    ChallengeResponseAuthentication yes
    PasswordAuthentication no
    UsePAM yes
    ```

    This enables the keyboard-interactive authentication and PAM authentication
    modules. The password authentication is disabled.

1. Restart the SSH service:

    ```plaintext
    $ sudo systemctl restart sshd
    ```

~> This step must be performed on **all** target hosts that you wish to connect
using the Vault's one-time SSH password.



### <a name="step2"></a>Step 2: Setup the SSH Secrets Engine
(**Persona:** operations)

On the ***Vault server***, you must enable the SSH secrets engine before you can
perform the operation. Then you are going to create a role named,
`otp_key_role`.

#### CLI command

First, enable the SSH secrets engine.

```plaintext
$ vault secrets enable ssh
```

Next, create a role.

```plaintext
$ vault write ssh/roles/otp_key_role key_type=otp \
        default_user=ubuntu \
        cidr_list=0.0.0.0/0
```

This creates `otp_key_role` with `ubuntu` as its default username for which a
credential will be generated.

#### API call using cURL

Enable `ssh` secret engine using `/sys/mounts` endpoint:

```plaintext
$ curl --header "X-Vault-Token: <TOKEN>" \
       --request POST \
       --data <PARAMETERS> \
       <VAULT_ADDRESS>/v1/sys/mounts/<PATH>
```

Where `<TOKEN>` is your valid token, and `<PARAMETERS>` holds [configuration
parameters](/api/system/mounts.html#enable-secrets-engine) of the secret engine.

**Example:**

The following example enables SSH secret engine at `sys/mounts/ssh`
path, and passed the secret engine type (`ssh`) in the request payload.

```plaintext
$ curl --header "X-Vault-Token: ..." \
       --request POST \
       --data '{"type":"ssh"}' \
       https://127.0.0.1:8200/v1/sys/mounts/ssh
```

Now, create a role using the `ssh/roles/otp_key_role` endpoint.

```plaintext
$ tee payload.json <<EOF
{
  "key_type": "otp",
  "default_user": "ubuntu",
  "cidr_list": "0.0.0.0/0"
}
EOF

$ curl --header "X-Vault-Token: ..." \
       --request POST \
       --data @payload.json
       https://127.0.0.1:8200/v1/ssh/roles/otp_key_role   
```

This creates `otp_key_role` with `ubuntu` as its default username for which a
credential will be generated.


#### Web UI

Open a web browser and launch the Vault UI (e.g. `http://127.0.0.1:8200/ui`) and
then login.

1. Select **Enable new engine** and select **SSH** from **Secrets engine type**
drop-down list.

1. Click **Enable Engine**.

1. Select **Create role**.

1. Enter **`otp_key_role`** in the **Role name** field, select **otp** from the
**Key type** drop-down list, and then enter **`ubuntu`** in the **Default user**
field.

1. Select **More options** to expand the optional parameter fields, and then
enter **`0.0.0.0/0`** in the **CIDR list** field.

![Create Role](/img/vault-ssh-otp-2.png)

1. Click **Create role**.

<br>

~> **NOTE:** Its recommended to create individual roles for each username to
ensure absolute isolation between usernames. This is required for Dynamic Key
type and OTP type. For the purpose of this guide, the **`cidr_list`** will be
set to `0.0.0.0/0`. For production, be sure to set this parameter to be as
granular as you can since there is no need to keep this role open to the world.


### <a name="step3"></a>Step 3: Request an OTP
(**Persona:** client)

The client must have the following permission to request an OTP for
`otp_key_role`.

```hcl
path "ssh/creds/otp_key_role" {
  capabilities = [ "update" ]
}
```

#### CLI command

To generate an OTP credential for an IP of the remote host belongs to the
`otp_key_role`:

```plaintext
$ vault write ssh/creds/otp_key_role ip=<REMOTE_HOST_IP>
```

**Example:**

```plaintext
$ vault write ssh/creds/otp_key_role ip=192.0.2.10
Key                Value
---                -----
lease_id           ssh/creds/otp_key_role/234bb081-d22e-3762-3ae5-744110ea4d0a
lease_duration     768h
lease_renewable    false
ip                 192.0.2.10
key                f1cb47ad-6255-0be8-6bd8-5c4b3b01c8df
key_type           otp
port               22
username           ubuntu
```

The **`key`** value is the OTP to use during SSH authentication.


#### API call using cURL

To generate an OTP credential for an IP of the remote host belongs to the
`otp_key_role`:

```plaintext
$ curl --header "X-Vault-Token: ..." \
       --request POST \
       --data '{"ip": "<REMOTE_HOST_IP>"}'
       https://127.0.0.1:8200/v1/ssh/creds/otp_key_role   
```

**Example:**

```plaintext
$ curl --header "X-Vault-Token: ..." \
       --request POST \
       --data '{"ip": "192.0.2.10"}'
       https://127.0.0.1:8200/v1/ssh/creds/otp_key_role  | jq
{
   "request_id": "83327e6b-cffb-ce77-c4e1-74a60dd69436",
   "lease_id": "ssh/creds/otp_key_role/aebd080e-51a1-0fc6-ef5c-1e1c95ca729a",
   "renewable": false,
   "lease_duration": 2764800,
   "data": {
     "ip": "192.0.20.10",
     "key": "6e472878-721a-b066-2cec-1bed0127ad44",
     "key_type": "otp",
     "port": 22,
     "username": "ubuntu"
   },
   ...
}
```

The **`key`** value is the OTP to use during SSH authentication.


#### Web UI

To generate an OTP credential for an IP of the remote host belongs to the
`otp_key_role`:

1. Select **`ssh`** under **Secrets Engines**.

1. Select **otp_key_role** and enter **`ubuntu`** in the **Username** field, and
enter the target host's IP address (e.g. `192.0.2.10`) in the **IP Address**
field.

1. Click **Generate**.

1. Click **Copy credentials**.  This copies the OTP (`key` value).


### <a name="step4"></a>Step 4: Establish an SSH session

Simply, use the OTP generated at [Step 3](#step3) to authenticate:

```plaintext
$ ssh ubuntu@192.0.2.10
Password: <Enter OTP>
```

<br>
**NOTE:** If [`sshpass`](https://gist.github.com/arunoda/7790979) is installed,
you can create a new OTP and SSH into the target host with single line of CLI
command:

```plaintext
$ vault ssh -role otp_key_role -mode otp -strict-host-key-checking=no ubuntu@192.0.2.10
```


## Next steps

Read the [Signed SSH
Certificates](/docs/secrets/ssh/signed-ssh-certificates.html) documentation if
you want to use signed SSH certificate so that the clients can use their local
SSH keys to connect to the target hosts.
