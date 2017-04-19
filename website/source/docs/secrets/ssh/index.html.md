---
layout: "docs"
page_title: "SSH Secret Backend"
sidebar_current: "docs-secrets-ssh"
description: |-
  The SSH secret backend for Vault generates signed SSH certificates, dynamic SSH keys or One-Time-Passwords.
---

# SSH Secret Backend

Name: `ssh`

Vault SSH backend tries to solve the problem of managing access to machine
infrastructure by providing different ways to issue SSH credentials.

The backend issues in 3 types of credentials: CA signed keys, Dynamic keys and
OTP keys. Read and carefully understand all the types before choosing the one
which best suits your needs. In relation to the dynamic key and OTP key type,
the CA key signing is the simplest and most powerful in terms of setup
complexity and in terms of being platform agnostic.

This page will show a quick start for this backend. For detailed documentation
on every path, use `vault path-help` after mounting the backend.

----------------------------------------------------
## I. CA Key Type

When using this type, an SSH CA signing key is generated or configured at the
backend's mount. This key will be used to sign other SSH keys. The private half
of the signing key always stays within Vault and the public half is exposed via
the API. Each mount of this backend represents a unique signing key pair. It is
recommended that the host keys and client keys are signed using different
mounts of this backend.

### Mount a backend's instance for signing host keys

```text
vault mount -path ssh-host-signer ssh
Successfully mounted 'ssh' at 'ssh-host-signer'!
```

### Mount a backend's instance for signing client keys

```text
vault mount -path ssh-client-signer ssh
Successfully mounted 'ssh' at 'ssh-client-signer'!
```

### Configure the host CA certificate

```text
vault write -f ssh-host-signer/config/ca
Key             Value
---             -----
public_key      ssh-rsa
AAAAB3NzaC1yc2EAAAADAQABAAACAQDB7lbQ4ub8/7AbfnF4zntCxPGijXzdo50MrUz0I1Vz/n98pBPyrUEcNqFeCr6Qs9kRL8J8Bu48W/DtSxAiScIIlzyezulH+wxdZnCN3sxDUhFbTUosWstSt6P5BZbJsMTpfHfuYbOexewX9ljw636fUrRtU3/0Emq2FkOcKEv9v7z+Qwdaom+8R6CvRA0Kit0yyh6DXD+IjPnOmhYFoNH6ZhPU2xJ+7M16vLpCmH6OO1krHlDK+ZUvVKD2ZjfmbegoYEV/9FETSzVnQ/pOrP0SdwpmGhei7VO7hQwa6zOWIah66TYYdaIopHFha60rZPybliZFAi/CcUO7NYldq7Q2ty68RKp+Im2ROPbXLKYc2N+PThWwbxA1yQ7shvgSy/JEbAcYQuDLGn6tA6n5sXpvFcJe+vlf75w+laXIPUXTVHr2g1O5Irzwunpg03Pa+Kab8Hp8uq61XfTfNtCwHYr/K5YO27rHcpMA/gBKkOFmBbcI3VlUbi266HiyFANz0JMObzXS5hEGm90NevMD2nEFvZHYHx4hXRtoox7Bh4CCj2fUPaThnuTvv1hLBgryWxwDXL7OwHGx6wd/dpXo/1fG1DRUotfVcJP7k7FG9MUAB6tOR2RCuXEaTF6bv3tCpbsdzzzyhs/aIa6IeYRV4HUo3nW5g/IRT+6F/MLd63Cq5w==
```

The returned host CA public key should be added to `known_hosts` file in all
client machines with a `@cert-authority <domain>` prefix, where `<domain>` can
be any value (single host, wildcard domain, wildcard) SSH accepts to define
hosts (e.g. `@cert-authority *.example.com ssh-rsa AAAA...`). This public key
can also be retrieved using `vault read ssh-host-signer/config/ca`.

### Configure the client CA certificate

```text
vault write -f ssh-client-signer/config/ca
Key             Value
---             -----
public_key      ssh-rsa
AAAAB3NzaC1yc2EAAAADAQABAAACAQDoYiTm0QrK+NGkv8H/Sq1Q9l2PQwxY1ROSPpRqOoFAgCM7vtEZSGchm8j3/Da6Vj84sUg+09N+dwMZ6z3HG57dBWyfHGCbEt+fIJtRcCV7KkQ4xktBdapq5PychBg0bYIP0+JWXe7yOAR6s+DWJgUUHwmYiY0GlrzoztvtYcU/pvj4zWY/EWcCAYRiZThqJOJ9rDzhEzmVHVYelZ8dlehKZsUOQ3dY5fCHgWJ5/IKSfGUI1w1h/3ENAToBKInbW2J/Wy0mbsfUJhha84zgqSeHd4gZ9YqaSpzYffVvDz6Ud+KtPWt4ejJDyvJVFAVnVIvVUCkuXVefBlC10Ww6l52DXex/rdxQW+46ZNV5zlFOUi4VWFPT6bOlIEj1hgHG0cfFKARQXcVlp03xTi5DRJp3jlBf/QLeUy4kZEFUdQl+gkVLwySrOh+1kHwE15JQWpNHxfxWzYnNnjokT2RCW2ANO4Y8vfMPCVM5OKSVNISHm3UurNkJNTAomK/ElnXQgP23UcVgbF3mNRBPx9U3/+9HHkXcqDefSI7XBC7npIBFdetCxFt3nW/cX4l1zfq6ZExCVde+/40V1U+TOoLZVLIZZzEuA40J4H99oK7QN5/MtCGaSr7kgAx+7ymfSY37rXB6wOVi5HTBlhyh+VFGU+Wmy+q5sxfr6+ZLGwoGpvNh0Q==
```

The returned client CA public key should be saved in a file that is added to
the `TrustedUserCAKeys` list in the `sshd_config` file on the host machine:

```text
cat /etc/ssh/sshd_config
...
...
TrustedUserCAKeys /etc/ssh/trusted-user-ca-keys.pem
```

The client CA public key can also be retrieved later using `vault read
ssh-client-signer/config/ca`.

### Allow host certificate to have longer TTLs

```text
vault mount-tune -max-lease-ttl=87600h ssh-host-signer
Successfully tuned mount 'ssh-host-signer'!
```

### Create a role to sign host keys

You'll want to add in some allowed domains and either the `allow_subdomains`
flag, the `allow_bare_domains` flag, or both.

```text
vault write ssh-host-signer/roles/hostrole ttl=87600h allow_host_certificates=true key_type=ca allowed_domains="localdomain,example.com" allow_subdomains=true
Success! Data written to: ssh-host-signer/roles/hostrole
```

### Create a role to sign client keys

Because of the way that some SSH certificate features are implemented, some
common options to this call must be passed in as a map. Vault's CLI does not
currently have a native way to pass in map values, so the easiest way to do
this is to create a JSON file and use it as the input. In the following example
this is used to add the `permit-pty` extension to the certificate.

```json
{
  "allow_user_certificates": true,
  "allowed_users": "*",
  "default_extensions": [
    {
      "permit-pty": ""
    }
  ],
  "key_type": "ca",
  "default_user": "icecream",
  "ttl": "30m0s"
}
```

```text
vault write ssh-client-signer/roles/clientrole @clientrole.json
Success! Data written to: ssh-client-signer/roles/clientrole
```

### Sign the host key

You can generate a new key, or you can sign one of the host keys that was
autogenerated when SSH was first started up, e.g.
`/etc/ssh/ssh_host_rsa_key.pub`.

```text
cat /etc/ssh/ssh_host_rsa_key.pub | vault write ssh-host-signer/sign/hostrole public_key=- cert_type=host
Key             Value
---             -----
serial_number   3746eb17371540d9
signed_key      ssh-rsa-cert-v01@openssh.com
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
```

Set the signed certificate as `HostCertificate` in `sshd_config` on the host
machine. In order to make things work more automagically (e.g. if you don't
want to specifically pick the host key/cert) it's a good idea to call this
`<key>-cert.pub`, e.g. `/etc/ssh/ssh_host_rsa_key-cert.pub`. It's also not a
bad idea to specifically pick the `HostKey` that will be used:

```text
cat /etc/ssh/sshd_config
...
...
TrustedUserCAKeys /etc/ssh/trusted-user-ca-keys.pem
HostKey /etc/ssh/ssh_host_rsa_key
HostCertificate /etc/ssh/ssh_host_rsa_key-cert.pub
```

It's also a good idea to make sure that the permissions on the
`HostCertificate` file are `0640` (on most systems).


### Sign the client key

This is any key you want to use, e.g. `~/.ssh/id_rsa.pub`:

```text
cd ~/.ssh

cat id_rsa.pub | vault write ssh-client-signer/sign/clientrole public_key=-
Key             Value
---             -----
serial_number   c73f26d2340276aa
signed_key      ssh-rsa-cert-v01@openssh.com
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
```

Save the signed key in a file (make sure you set the permissions to only be
writeable by the owner); if you're saving it along with your public/private key
(e.g. `~/.ssh/id_rsa(.pub)`), saving it with `-cert.pub` (e.g.
`~/.ssh/id_rsa-cert.pub`) allows SSH to automatically discover it. Otherwise,
you can tell the SSH client where to find it (and the associated private key)
by using `-i` such as `ssh -i ~/.ssh/id_rsa -i my-signed-cert.pub localhost`.

If you want to see the configured extensions, principals, etc., you can use the
`-L` flag to `ssh-keygen`, e.g. `ssh-keygen -Lf ~/.ssh/id_rsa-cert.pub`.

### SSH into the host machine

```text
ssh -i signed-client-cert.pub username@<IP of remote host>
username@<IP of remote host>:~$
```

### Troubleshooting

If you are not able to successfully make a client connection, looking at SSH server logs is likely to be your best bet. Some known potential issues:

* If on an SELinux-enforcing system, you may need to adjust related types so
  that the SSH daemon is able to read it, for instance by adjusting the signed
  host certificate to be an `sshd_key_t` type.
* If you can an error on the client indicating `"no separate private key for
  certificate"` you may be hitting a bug introduced into OpenSSH version 7.2
  and fixed in 7.5. See [OpenSSH bug
  2617](https://bugzilla.mindrot.org/show_bug.cgi?id=2617) for details. If you
  are able to find a workaround without patching OpenSSH, please submit a PR to
  update this documentation!

----------------------------------------------------
## II. One-Time-Password (OTP) Type

This backend type allows a Vault server to issue an OTP every time a client
wants to SSH into a remote host, using a helper command on the remote host to
perform verification.

An authenticated client requests credentials from the Vault server and, if
authorized, is issued an OTP. When the client establishes an SSH connection to
the desired remote host, the OTP used during SSH authentication is received by
the Vault helper, which then validates the OTP with the Vault server. The Vault
server then deletes this OTP, ensuring that it is only used once.

Since the Vault server is contacted during SSH connection establishment, every
login attempt and the correlating Vault lease information is logged to the
audit backend.

See [Vault-SSH-Helper](https://github.com/hashicorp/vault-ssh-helper) for
details on the helper.

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

----------------------------------------------------
## III. Dynamic Key Type (Deprecated)

**Note**: There are several serious drawbacks (detailed below), including some
with security implications, inherent in this method. Because of these
drawbacks, the Vault team recommends use of the CA or OTP types whenever
possible. Care should be taken with respect to the above issues with any
deployments using the dynamic key type.

When using this type, the administrator registers a secret key with appropriate
`sudo` privileges on the remote machines; for every authorized credential
request, Vault creates a new SSH key pair and appends the newly-generated
public key to the `authorized_keys` file for the configured username on the
remote host. Vault uses a configurable install script to achieve this.

The backend does not prompt for `sudo` passwords; the `NOPASSWD` option for
sudoers should be enabled at all remote hosts for the Vault administrative
user.

The private key returned to the user will be leased and can be renewed if
desired. Once the key is given to the user, Vault will not know when it gets
used or how many time it gets used. Therefore, Vault **WILL NOT** and cannot
audit the SSH session establishments.

When the credential lease expires, Vault removes the secret key from the remote
machine.

### Drawbacks

The dynamic key type has several serious drawbacks:

1. _Audit logs are unreliable_: Vault can only log when users request
   credentials, not when they use the given keys. If user A and user B both
   request access to a machine, and are given a lease valid for five minutes,
   it is impossible to know whether two accesses to that user account on the
   remote machine were A, A; A, B; B, A; or B, B.
2. _Generating dynamic keys consumes entropy_: Unless equipped with a hardware
   entropy generating device, a machine can quickly run out of entropy when
   generating SSH keys. This will cause further requests for various Vault
   operations to stall until more entropy is available, which could take a
   significant amount of time, after which the next request for a new SSH key
   will use the generated entropy and cause stalling again.
3. This type makes connections to client hosts; when this happens the host key
   is *not* verified.

### sudo

In order to adjust the `authorized_keys` file for the desired user, Vault
connects via SSH to the remote machine as a separate user, and uses `sudo` to
gain the privileges required. An example `sudoers` file is shown below.

File: `/etc/sudoers`

```hcl
# This is a sample sudoers statement; you should modify it
# as appropriate to satisfy your security needs.
vaultadmin   ALL=(ALL)NOPASSWD: ALL
```

### Configuration

Next, infrastructure configuration must be registered with Vault via roles.
First, however, the shared secret key must be specified.

### Mount the backend

```text
$ vault mount ssh
Successfully mounted 'ssh' at 'ssh'!
```

#### Registering the shared secret key

Register a key with a name; this key must have administrative capabilities on
the remote hosts.

```text
$ vault write ssh/keys/dev_key \
    key=@dev_shared_key.pem
```

#### Create a Role

Next, create a role. All of the machines contained within this CIDR block list
should be accessible using the registered shared secret key.

```text
$ vault write ssh/roles/dynamic_key_role \
    key_type=dynamic \
    key=dev_key \
    admin_user=username \
    default_user=username \
    cidr_list=x.x.x.x/y
Success! Data written to: ssh/roles/dynamic_key_role
```

`cidr_list` is a comma separated list of CIDR blocks for which a role can
generate credentials. If this is empty, the role can only generate credentials
if it belongs to the set of zero-address roles.

Zero-address roles, configured via `/ssh/config/zeroaddress` endpoint, takes
comma separated list of role names that can generate credentials for any IP
address.

Use the `install_script` option to provide an install script if the remote
hosts do not resemble a typical Linux machine. The default script is compiled
into the Vault binary, but it is straight forward to specify an alternate.  The
script takes three arguments which are explained in the comments.

To see the default, see
[linux_install_script.go](https://github.com/hashicorp/vault/blob/master/builtin/logical/ssh/linux_install_script.go)

### Create a credential

Create a dynamic key for an IP of the remote host that is covered by
`dynamic_key_role`'s CIDR list.

```text
$ vault write ssh/creds/dynamic_key_role ip=x.x.x.x
Key            	Value
lease_id       	ssh/creds/dynamic_key_role/8c4d2042-23bc-d6a8-42c2-6ff01cb83cf8
lease_duration 	600
lease_renewable	true
ip             	x.x.x.x
key            	-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEA5V/Y95qfGaUXRPkKNK9jgDHXPD2n5Ein+QTNnLSGrHtJUH7+
pgs/5Hc4//124P9qHNmjIYQVyvcLreFgSrQCq4K8193hmypBYtsvCgvpc+jEwaGA
zK0QV7uc1z8KL7FuRAxpHJwB6+nubOzzqM03xsViHRhaWhYVHw2Vl4oputSHE7R9
ugaTRg67wge4Nyi5RRL0RQcmW15/Vop8B6HpBSmZQy3enjg+32KbOWCMMTAPuF9/
DgxSgZQaFMjGN4RjDreZI8Vv5zIiFJzZ3KVOWy8piI0PblLnDpU4Q0QSQ9A+Vr7b
JS22Lbet1Zbapl/n947/r1wGObLCc5Lilu//1QIDAQABAoIBAHWLfdO9sETjHp6h
BULkkpgScpuTeSN6vGHXvUrOFKn1cCfJPNR4tWBuXI6LJM2+9nEccwXs+4IMwjZ0
ZfVCdI/SKtZxBXmP2PxBGMUMP7G/mn0kN64sDlD3ezOvQZgZVEmZFpCrvixYsG+v
qlpZ+HhrlJEWds7tvBsyyfNjwWjVIpm08zBmteFj4zu7OEcmGXEHDoxDXxyVP2BG
eLU/fM5JA2UEjfCQ1MIZ3rBtPePdz4LRpb+ajklqrUj1OHoiDrXa8EAf0/wDP9re
c1iH4bn7ZjYK0+IhZ+Pmw6gUftzZNWSC2kOLnZLdN/K7hgh0l0r0K/1eeXt43upB
WALNuiECgYEA8PM2Ob3XXKALF86PUewne4fCz9iixr/cIpvrEGrh9lyQRO8X5Jxb
ug38jEql4a574C6TSXfzxURza4P6lnfa0LvymmW0bhxZ5nev9kcAVnLKvpOUArTR
32k9bKXd6zp8Q9ZyVNwHRxcVs4YgwfJlcx8geC4o6YRiIjvcBQ9RVHkCgYEA87OK
lZDFBeEY/HVOxAQNXS5fgTd4U4DbwEJLv7SPk02v9oDkGHkpgMs4PcsIpCzsTpJ0
oXMfLSxZ1lmZiuUvAupKj/7RjJ0XyjSMfm1Zs81epWj+boVfM4amZNHVLIWgddmM
XzXEZKByvi1gs7qFcjQz2DEbZltWO6dX14O4Fz0CgYEAlWSWyHJWZ02r0xT1c7vS
NxtTxH7zXftzR9oYgtNiStfVc4gy7kGr9c3aOjnGZAlFMRhvpevDrxnj3lO0OTsS
5rzBjM1mc6cMboLjDPW01eTSpBroeE0Ym0arGQQ2djSK+5yowsixknhTsj2FbfsW
v6wa+6jTIQY9ujAXGOQIbzECgYAYuXlw7SwgCZNYYappFqQodQD5giAyEJu66L74
px/96N7WWoNJvFkqmPOOyV+KEIi0/ATbMGvUUHCY36RFRDU9zXldHJQz+Ogl+qja
VsvIAyj8DSfrHJrpBlsxVVyUVMZPzo+ARVs0flbF1qK9+Ul6qbMs1uaZvuCD0tmF
ovZ1XQKBgQDB0s7SDmAMgVjG8UBZgUru9vsDrxERT2BloptnnAjSiarLF5M+qeZO
7L4NLyVP39Z83eerEonzDAHHbvhPyi6n2YmnYhGjeP+lPZIVqGF9cpZD3q48YHZc
3ePn2/oLZrXKWOMyMwp2Uj+0SArCW+xMnoNp50sYNVR/JK3BPIdkag==
-----END RSA PRIVATE KEY-----
key_type       	dynamic
port           	22
username       	username
```

### Establish an SSH session

Save the key to a file (e.g. `dyn_key.pem`) and then use it to establish an SSH
session.

```text
$ ssh -i dyn_key.pem username@<IP of remote host>
username@<IP of remote host>:~$
```

### Automate it!

Creation of new key, saving to a file, and using it to establish an SSH session
can all be done with a single Vault CLI command.

```text
$ vault ssh -role dynamic_key_role username@<IP of remote host>
username@<IP of remote host>:~$
```

----------------------------------------------------
## API

The SSH secret backend has a full HTTP API. Please see the
[SSH secret backend API](/api/secret/ssh/index.html) for more
details.
