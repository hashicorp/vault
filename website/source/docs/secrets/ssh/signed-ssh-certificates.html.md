---
layout: "docs"
page_title: "Signed SSH Certificates - SSH Secret Backend"
sidebar_current: "docs-secrets-ssh-signed-ssh-certificates"
description: |-
  The signed SSH certificates is the simplest and most powerful in terms of
  setup complexity and in terms of being platform agnostic. When using this
  type, an SSH CA signing key is generated or configured at the backend's mount.
  This key will be used to sign other SSH keys.
---

# Signed SSH Certificates

The signed SSH certificates is the simplest and most powerful in terms of setup
complexity and in terms of being platform agnostic. By leveraging Vault's
powerful CA capabilities and functionality built into OpenSSH, clients can SSH
into target hosts using their own local SSH keys.

In this section, the term "**client**" refers to the person or machine
performing the SSH operation. The "**host**" refers to the target machine. If
this is confusing, substitute "client" with "user".

This page will show a quick start for this backend. For detailed documentation
on every path, use `vault path-help` after mounting the backend.

## Client Key Signing

Before a client can request their SSH key be signed, the Vault SSH backend must
be configured. Usually a Vault administrator or security team performs these
steps. It is also possible to automate these actions using a configuration
management tool like Chef, Puppet, Ansible, or Salt.

### Signing Key &amp; Role Configuration

The following steps are performed in advance by a Vault administrator, security
team, or configuration management tooling.

1. Mount the backend. Like all secret backends in Vault, the SSH secret backend
must be mounted before use.

    ```text
    $ vault mount -path=ssh-client-signer ssh
    Successfully mounted 'ssh' at 'ssh-client-signer'!
    ```

    This mounts the SSH backend at the path "ssh-client-signer". It is possible
    to mount the same secret backend multiple times using different `-path`
    arguments. The name "ssh-client-signer" is not special - it can be any name,
    but this documentation will assume "ssh-client-signer".

1. Configure Vault with a CA for signing client keys using the `/config/ca`
endpoint. If you do not have an internal CA, Vault can generate a keypair for
you.

    ```text
    $ vault write ssh-client-signer/config/ca generate_signing_key=true
    Key             Value
    ---             -----
    public_key      ssh-rsa AAAAB3NzaC1yc2EA...
    ```

    If you already have a keypair, specify the public and private key parts as
    part of the payload:

    ```text
    $ vault write ssh-client-signer/config/ca \
        private_key="..." \
        public_key="..."
    ```

    Regardless of whether it is generated or uploaded, the client signer public
    key is accessible via the API at the `/public_key` endpoint.

1. Add the public key to all target host's SSH configuration. This process can
be manual or automated using a configuration management tool. The public key is
accessible via the API and does not require authentication.

    ```text
    $ curl -o /etc/ssh/trusted-user-ca-keys.pem https://vault.rocks/v1/ssh-client-signer/public_key
    ```

    ```text
    $ vault read -field=public_key ssh-client-signer/config/ca > /etc/ssh/trusted-user-ca-keys.pem
    ```

    Add the path where the public key contents are stored to the SSH
    configuration file as the `TrustedUserCAKeys` option.

    ```text
    # /etc/ssh/sshd_config
    # ...
    TrustedUserCAKeys /etc/ssh/trusted-user-ca-keys.pem
    ```

    Restart the SSH service to pick up the changes.

1. Create a named Vault role for signing client keys.

    Because of the way some SSH certificate features are implemented, options
    are passed as a map. The following example adds the `permit-pty` extension
    to the certificate.

    ```text
    $ vault write ssh-client-signer/roles/my-role -<<"EOH"
    {
      "allow_user_certificates": true,
      "allowed_users": "*",
      "default_extensions": [
        {
          "permit-pty": ""
        }
      ],
      "key_type": "ca",
      "default_user": "ubuntu",
      "ttl": "30m0s"
    }
    EOH
    ```

### Client SSH Authentication

The following steps are performed by the client (user) that wants to
authenticate to machines managed by Vault. These commands are usually run from
the client's local workstation.

1. Locate or generate the SSH public key. Usually this is `~/.ssh/id_rsa.pub`.
If you do not have an SSH keypair, generate one:

    ```text
    $ ssh-keygen -t rsa -C "user@example.com"
    ```

1. Ask Vault to sign your **public key**. This file usually ends in `.pub` and
the contents begin with `ssh-rsa ...`.

    ```text
    $ vault write ssh-client-signer/sign/my-role \
        public_key=@$HOME/.ssh/id_rsa.pub

    Key             Value
    ---             -----
    serial_number   c73f26d2340276aa
    signed_key      ssh-rsa-cert-v01@openssh.com AAAAHHNzaC1...
    ```

    The result will include the serial and the signed key. This signed key is
    another public key.

    To customize the signing options, use a JSON payload:

    ```text
    $ vault write ssh-client-signer/sign/my-role -<<"EOH"
    {
      "public_key": "ssh-rsa AAA...",
      "valid_principals": "my-user",
      "key_id": "custom-prefix",
      "extension": {
        "permit-pty": ""
      }
    }
    EOH
    ```

1. Save the resulting signed, public key to disk. Limit permissions as needed.

    ```text
    $ vault write -field=signed_key ssh-client-signer/sign/my-role \
        public_key=@$HOME/.ssh/id_rsa.pub > signed-cert.pub
    ```

    If you are saving the certificate directly beside your SSH keypair, suffix
    the name with `-cert.pub` (`~/.ssh/id_rsa-cert.pub`). With this naming
    scheme, OpenSSH will automatically use it during authentication.

1. (Optional) View enabled extensions, principals, and metadata of the signed
key.

    ```text
    $ ssh-keygen -Lf ~/.ssh/signed-cert.pub
    ```

1. SSH into the host machine using the signed key. You must supply both the
signed public key from Vault **and** the corresponding private key as
authentication to the SSH call.

    ```text
    $ ssh -i signed-cert.pub -i ~/.ssh/id_rsa username@10.0.23.5
    ```

## Host Key Signing

For an added layers of security, we recommend enabling host key signing. This is
used in conjunction with client key signing to provide an additional integrity
layer. When enabled, the SSH agent will verify the target host is valid and
trusted before attempting to SSH. This will reduce the probability of a user
accidentally SSHing into an unmanaged or malicious machine.

### Signing Key Configuration

1. Mount the backend. For the most security, mount at a different path from the
client signer.

    ```text
    $ vault mount -path=ssh-host-signer ssh
    Successfully mounted 'ssh' at 'ssh-host-signer'!
    ```

1. Configure Vault with a CA for signing host keys using the `/config/ca`
endpoint. If you do not have an internal CA, Vault can generate a keypair for
you.

    ```text
    $ vault write ssh-host-signer/config/ca generate_signing_key=true
    Key             Value
    ---             -----
    public_key      ssh-rsa AAAAB3NzaC1yc2EA...
    ```

    If you already have a keypair, specify the public and private key parts as
    part of the payload:

    ```text
    $ vault write ssh-host-signer/config/ca \
        private_key="..." \
        public_key="..."
    ```

    Regardless of whether it is generated or uploaded, the host signer public
    key is accessible via the API at the `/public_key` endpoint.

1. Extend host key certificate TTLs.

    ```text
    $ vault mount-tune -max-lease-ttl=87600h ssh-host-signer
    ```

1. Create a role for signing host keys. Be sure to fill in the list of allowed
domains, set `allow_bare_domains`, or both.

    ```text
    $ vault write ssh-host-signer/roles/hostrole \
        key_type=ca \
        ttl=87600h \
        allow_host_certificates=true \
        allowed_domains="localdomain,example.com" \
        allow_subdomains=true
    ```

1. Sign the host's SSH public key.

    ```text
    $ vault write ssh-host-signer/sign/hostrole \
        cert_type=host \
        public_key=@/etc/ssh/ssh_host_rsa_key.pub
    Key             Value
    ---             -----
    serial_number   3746eb17371540d9
    signed_key      ssh-rsa-cert-v01@openssh.com AAAAHHNzaC1y...
    ```

1. Set the resulting signed certificate as `HostCertificate` in the SSH
configuration on the host machine.

    ```text
    $ vault write -field=signed_key ssh-host-signer/sign/hostrole \
        cert_type=host \
        public_key=@/etc/ssh/ssh_host_rsa_key.pub > /etc/ssh/ssh_host_rsa_key-cert.pub
    ```

    Set permissions on the certificate to be `0640`:

    ```text
    $ chmod 0640 /etc/ssh/ssh_host_rsa_key-cert.pub
    ```

    Add host key and host certificate to the SSH configuration file.

    ```text
    # /etc/ssh/sshd_config
    # ...

    # For client keys
    TrustedUserCAKeys /etc/ssh/trusted-user-ca-keys.pem

    # For host keys
    HostKey /etc/ssh/ssh_host_rsa_key
    HostCertificate /etc/ssh/ssh_host_rsa_key-cert.pub
    ```

    Restart the SSH service to pick up the changes.

### Client-Side Host Verification

1. Retrieve the host signing CA public key to validate the host signature of
target machines.

    ```text
    $ curl https://vault.rocks/v1/ssh-host-signer/public_key
    ```

    ```text
    $ vault read -field=public_key ssh-host-signer/config/ca
    ```

1. Add the resulting public key to the `known_hosts` file with authority.

    ```text
    # ~/.ssh/known_hosts
    @cert-authority *.example.com ssh-rsa AAAAB3NzaC1yc2EAAA...
    ```

1. SSH into target machines as usual.

## Troubleshooting

When initially configuring this type of key signing, enable `VERBOSE` SSH
logging to help annotate any errors in the log.

```text
# /etc/ssh/sshd_config
# ...
LogLevel VERBOSE
```

Restart SSH after making these changes.

By default, SSH logs to `/var/log/auth.log`, but so do many other things. To
extract just the SSH logs, use the following:

```sh
$ tail -f /var/log/auth.log | grep --line-buffered "sshd"
```

If you are unable to make a connection to the host, the SSH server logs may
provide guidance and insights.

### Name is not a listed principal

If the `auth.log` displays the following messages:

```text
# /var/log/auth.log
key_cert_check_authority: invalid certificate
Certificate invalid: name is not a listed principal
```

The certificate does not permit the username as a listed principal for
authenticating to the system. This is most likely due to an OpenSSH bug (see
[known issues](#known-issues) for more information). This bug does not respect
the `allowed_users` option value of "\*". Here are ways to work around this
issue:

1. Set `default_user` in the role. If you are always authenticating as the same
user, set the `default_user` in the role to the username you are SSHing into the
target machine:

    ```text
    $ vault write ssh/roles/my-role -<<"EOH"
    {
      "default_user": "YOUR_USER",
      // ...
    }
    EOH
    ```

1. Set `valid_principals` during signing. In situations where multiple users may
be authenticating to SSH via Vault, set the list of valid principles during key
signing to include the current username:

    ```text
    $ vault write ssh-client-signer/sign/my-role -<<"EOH"
    {
      "valid_principals": "my-user"
      // ...
    }
    EOH
    ```


### No Prompt After Login

If you do not see a prompt after authenticating to the host machine, the signed
certificate may not have the `permit-pty` extension. There are two ways to add
this extension to the signed certificate.

- As part of the role creation

    ```text
    $ vault write ssh-client-signer/roles/my-role -<<"EOH"
    {
      "default_extensions": [
        {
          "permit-pty": ""
        }
      ]
      // ...
    }
    EOH
    ```

- As part of the signing operation itself:

    ```text
    $ vault write ssh-client-signer/sign/my-role -<<"EOH"
    {
      "extension": {
        "permit-pty": ""
      }
      // ...
    }
    EOH
    ```

### No Port Forwarding

If port forwarding from the guest to the host is not working, the signed
certificate may not have the `permit-port-forwarding` extension. Add the
extension as part of the role creation or signing process to enable port
forwarding. See [no prompt after login](#no-prompt-after-login) for examples.

```json
{
  "default_extensions": [
    {
      "permit-port-forwarding": ""
    }
  ]
}
```

### No X11 Forwarding

If X11 forwarding from the guest to the host is not working, the signed
certificate may not have the `permit-X11-forwarding` extension. Add the
extension as part of the role creation or signing process to enable X11
forwarding. See [no prompt after login](#no-prompt-after-login) for examples.

```json
{
  "default_extensions": [
    {
      "permit-X11-forwarding": ""
    }
  ]
}
```

### No Agent Forwarding

If agent forwarding from the guest to the host is not working, the signed
certificate may not have the `permit-agent-forwarding` extension. Add the
extension as part of the role creation or signing process to enable agent
forwarding. See [no prompt after login](#no-prompt-after-login) for examples.

```json
{
  "default_extensions": [
    {
      "permit-agent-forwarding": ""
    }
  ]
}
```

### Known Issues

- On SELinux-enforcing systems, you may need to adjust related types so that the
  SSH daemon is able to read it. For example, adjust the signed host certificate
  to be an `sshd_key_t` type.

- On some versions of SSH, you may get the following error:

    ```text
    no separate private key for certificate
    ```

    This is a bug introduced in OpenSSH version 7.2 and fixed in 7.5. See
    [OpenSSH bug 2617](https://bugzilla.mindrot.org/show_bug.cgi?id=2617) for
    details.

## API

The SSH secret backend has a full HTTP API. Please see the
[SSH secret backend API](/api/secret/ssh/index.html) for more
details.
