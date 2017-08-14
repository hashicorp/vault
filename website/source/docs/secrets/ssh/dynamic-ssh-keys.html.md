---
layout: "docs"
page_title: "Dynamic SSH Keys - SSH Secret Backend"
sidebar_current: "docs-secrets-ssh-dynamic-ssh-keys"
description: |-
  When using this type, the administrator registers a secret key with
  appropriate sudo privileges on the remote machines. For every authorized
  credential request, Vault creates a new SSH key pair and appends the
  newly-generated public key to the authorized_keys file for the configured
  username on the remote host. Vault uses a configurable install script to
  achieve this.
---

# Dynamic SSH Keys

~> **Deprecated**: There are several serious drawbacks and security implications
inherent in this type. Because of these drawbacks, please use the SSH CA or OTP
types whenever possible.

When using this type, the administrator registers a secret key with appropriate
`sudo` privileges on the remote machines; for every authorized credential
request, Vault creates a new SSH key pair and appends the newly-generated public
key to the `authorized_keys` file for the configured username on the remote
host. Vault uses a configurable install script to achieve this.

The backend does not prompt for `sudo` passwords; the `NOPASSWD` option for
sudoers should be enabled at all remote hosts for the Vault administrative
user.

The private key returned to the user will be leased and can be renewed if
desired. Once the key is given to the user, Vault will not know when it gets
used or how many time it gets used. Therefore, Vault **WILL NOT** and cannot
audit the SSH session establishments.

When the credential lease expires, Vault removes the secret key from the remote
machine.

This page will show a quick start for this backend. For detailed documentation
on every path, use `vault path-help` after mounting the backend.

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

## API

The SSH secret backend has a full HTTP API. Please see the
[SSH secret backend API](/api/secret/ssh/index.html) for more
details.
