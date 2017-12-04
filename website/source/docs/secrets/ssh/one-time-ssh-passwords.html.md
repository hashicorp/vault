---
layout: "docs"
page_title: "One-Time SSH Passwords (OTP) - SSH Secret Backend"
sidebar_current: "docs-secrets-ssh-one-time-ssh-passwords"
description: |-
  The One-Time SSH Password (OTP) SSH secret backend type allows a Vault server
  to issue a One-Time Password every time a client wants to SSH into a remote
  host using a helper command on the remote host to perform verification.
---

# One-Time SSH Passwords

The One-Time SSH Password (OTP) SSH secret backend type allows a Vault server to
issue a One-Time Password every time a client wants to SSH into a remote host
using a helper command on the remote host to perform verification.

An authenticated client requests credentials from the Vault server and, if
authorized, is issued an OTP. When the client establishes an SSH connection to
the desired remote host, the OTP used during SSH authentication is received by
the Vault helper, which then validates the OTP with the Vault server. The Vault
server then deletes this OTP, ensuring that it is only used once.

Since the Vault server is contacted during SSH connection establishment, every
login attempt and the correlating Vault lease information is logged to the audit
backend.

See [Vault-SSH-Helper](https://github.com/hashicorp/vault-ssh-helper) for
details on the helper.

This page will show a quick start for this backend. For detailed documentation
on every path, use `vault path-help` after mounting the backend.

### Drawbacks

The main concern with the OTP backend type is the remote host's connection to
Vault; if compromised, an attacker could spoof the Vault server returning a
successful request. This risk can be mitigated by using TLS for the connection
to Vault and checking certificate validity; future enhancements to this backend
may allow for extra security on top of what TLS provides.

### Mount the backend

```text
$ vault mount ssh
Successfully mounted 'ssh' at 'ssh'!
```

### Create a Role

Create a role with the `key_type` parameter set to `otp`. All of the machines
represented by the role's CIDR list should have helper properly installed and
configured.

```text
$ vault write ssh/roles/otp_key_role \
    key_type=otp \
    default_user=username \
    cidr_list=x.x.x.x/y,m.m.m.m/n
Success! Data written to: ssh/roles/otp_key_role
```

### Create a Credential

Create an OTP credential for an IP of the remote host that belongs to
`otp_key_role`.

```text
$ vault write ssh/creds/otp_key_role ip=x.x.x.x
Key            	Value
lease_id       	ssh/creds/otp_key_role/73bbf513-9606-4bec-816c-5a2f009765a5
lease_duration 	600
lease_renewable	false
port           	22
username       	username
ip             	x.x.x.x
key            	2f7e25a2-24c9-4b7b-0d35-27d5e5203a5c
key_type       	otp
```

### Establish an SSH session

```text
$ ssh username@localhost
Password: <Enter OTP>
username@ip:~$
```

### Automate it!

A single CLI command can be used to create a new OTP and invoke SSH with the
correct parameters to connect to the host.

```text
$ vault ssh -role otp_key_role username@x.x.x.x
OTP for the session is `b4d47e1b-4879-5f4e-ce5c-7988d7986f37`
[Note: Install `sshpass` to automate typing in OTP]
Password: <Enter OTP>
```

The OTP will be entered automatically using `sshpass` if it is installed.

```text
$ vault ssh -role otp_key_role -strict-host-key-checking=no username@x.x.x.x
username@<IP of remote host>:~$
```

Note: `sshpass` cannot handle host key checking. Host key checking can be
disabled by setting `-strict-host-key-checking=no`.

## API

The SSH secret backend has a full HTTP API. Please see the
[SSH secret backend API](/api/secret/ssh/index.html) for more
details.
