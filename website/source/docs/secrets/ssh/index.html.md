---
layout: "docs"
page_title: "Secret Backend: SSH"
sidebar_current: "docs-secrets-ssh"
description: |-
  The SSH secret backend for Vault generates dynamic SSH keys or One-Time-Passwords. 
---

# SSH Secret Backend

Name: `ssh`

The SSH secret backend for Vault generates SSH credentials dynamically. This solves
the problem of managing private keys of the infrastructure. There are 2 options
available with this backend. 
1) Dynamic Type: Registering private keys (having root privileges) with Vault. Vault then issues
leased dynamic credentials to Vault authenticated users. Vault uses the registered
private key to install a new key for the user in the target host. This key will
be a long lived key and gets deleted only after the lease is expired. After the
user receiving the dynamic keys, Vault will have no control on the sessions created
with that key and hence the sessions will not be audited. Which brings us to option 2.

2) One-Time-Password (OTP) Type: Installing Vault-SSH agent in the target machines
and enabling challenge response mechanism for client authentication. Vault server
issues a OTP upon user request. During authentication, agent acts as a PAM module
and validates the password with Vault server. Since Vault server is contacted for
every SSH session establishment, they all get audited.

## Quick Start

`ssh` backend is not mounted by default. So, the first step in using the SSH backend
is to mount it.

```text
$ vault mount ssh
Successfully mounted 'ssh' at 'ssh'!
```

Next, we must register infrastructures with Vault. This is done by writing the role
information. The type of credentials created are determined by the key_type option.

### Dynamic key type

Create a named key, say "dev_key", in Vault which represents a registered shared key.

```text
$ vault write ssh/keys/dev_key key=@dev_shared_key.pem
```

Assuming that the target machine is hosted on Linux, create a script "key-linux-install.sh"
that can install the given public key in the authorized keys file.

```text
# This script file installs or uninstalls an RSA public key to/from authoried_keys
# file in a typical linux machine. This script should be registered with vault
# server while creating a role for key type 'dynamic'.

if [ $1 != "install" && $1 != "uninstall" ]; then
	exit 1
fi

# If the key being installed is already present in the authorized_keys file, it is
# removed and the result is stored in a temporary file.
grep -vFf $2 $3 > temp_$2

# Contents of temporary file will be the contents of authorized_keys file.
cat temp_$2 > $3

if [ $1 == "install" ]; then
# New public key is appended to authorized_keys file
cat $2 >> $3
fi

# Auxiliary files are deleted
rm -f $2 temp_$2
```

Create a role "dynamic_key_role". All the machines represented by CIDR block
should be accessible through "dev_key" with root privileges.

```text
$ vault write ssh/roles/dynamic_key_role key_type=dynamic key=dev_key admin_user=foo default_user=bar cidr=x.x.x.x/y install_script=@key-linux-install.sh
```

Create a dynamic key for an IP that belongs to the "dynamic_key_role".

```text
$ vault write ssh/creds/dynamic_key_role ip=x.x.x.x
Key            	Value
lease_id       	ssh/creds/dynamic_key_role/862d55c9-e54e-917c-4c8f-f6e1a54b2e51
lease_duration 	600
lease_renewable	true
key_type       	dynamic
key            	-----BEGIN RSA PRIVATE KEY-----
MIIEogIBAAKCAQEArp4Y31kwSaIVcZ/geLCfhrbG4fpBXTTcPgefo/9YNUGCmbiC
pqHcW7TJ7wpLdWYTxEoD8fZJ5GKIYKvesGkiG2as6iBXrxYp+byZkZ8TmAbYyxhk
j5RN2Arxb7tWL/9FuLNrH0sa/xPX117mhKNdV1RquSNehqGvfC4Vd2Rl43tXyNpM
WSr787ERfAq4EqQfZC17QauUCy0DJwy5vP7t0QzIuCh9GZT+pFvXNJcEm4NkhJbh
jp+cU9JTEQW+Tw6BwDtGFhgSQd7KFd+7Asx4T8UfDb3461cRqLcAcfM1+Y18DNcP
chf3OP0qDJw2ovvGZ6X3f/+6GIttSttciqmF2QIDAQABAoIBACqVJ1+gMmRigHQ7
FtSXze9eN1X4X2RJdcQyu72UkYA7P4wZMNNN+Zzrk6sViZ1RjVR68EdbVl25oaRh
hWbj3ItuGJDn3jo2X3olghW/A1o5oTi19CAHfIxI7uPefYAq8me+aUsyV50Iy8Qb
wn9qD2MylOwdMfoHB/Jyko2RED/O9zBtlCz6qObFOimLNRKKNoK1gz0KctRQaV6j
2PHrnyF2OCuxFcEU9gOEW5rGlBxkhQbiBWYC8HXALgcpQ4FyF4MxnQuyGAQVQZh2
FhhuOBW0iiElK8U+WOwMTyZZBHhbszFF05CM8IsWvqJLkKuCmbHG2Mq0Irigo9gR
HfNDhnkCgYEAyKemyGv27bXjSEhudtBcP+EoTrfqhLjGNCO7J/LlBvUpDo4yjXq3
z3J+jPuzakfOfN27xBNODP1tWyIwZ45Aozoa3enuk/kFhNAJVkgFa7JEuoqBFgCH
pj51JXLtF6K+2JbpJfYNSGbVPHNfSvs4uJoWKXZt4QATbdt9gBvww1sCgYEA3sfv
v9to8vSyD1Du9kxyl6PjiXc+CNagYenUmoHJaRFupDBIoUH65XCuXwcJvloakIX4
XAwuHtkPJcFKGX0mh5btbwoOtz3nb4hp2LEY/T05Jam5bYfRJZop112Xc/MLUm0J
/oazn5p06kg/Z2SwrY+IAs7VMm6PT/6NZvt25dsCgYBOfJ2Vef29n88GgCaNXRUo
e4cLu48FWU1WKb/UcYM6hH0Jz39grebmQy/TL8VPRkUzvHvsx2xZUmwLIMV0TEVm
U50cvptu0BJjkAiG8mcEaFfP68twcsactYOXIWwyOZuTFvydt7AcaPTxz2Mv7jKS
qtsOXt++CgyPhTKDAOrdTwKBgE3JW+IOl0d1vxJv/PAM41olRFaERynI3vkxLyW/
uXaxOoOjxEhiBFvGi2vsxi8rwOjDjmN9cUEeIxbYtanOs/xV65OA3ICI4d1ksSiT
NZl+ngyThYZEDPfnK0Lij/ZRX5upLPstR1ysDrSbA2BznOkNG713QKO6TNnulKrn
lK1PAoGAfK2HtnHwiBQC2OobA84tlx6571zuTcoFl0FN74fDUIChh4YVzVXsYBcp
1PFYe3YpCpgNwjmX8uNBHWVL/m21c55C88QAExtfoUsQ3dAJvunpJ6MTGBoTDMWi
HRTKLBJUNd9V410xllz+uupFayoMJyfesULETjuT/UYXBoon46I=
-----END RSA PRIVATE KEY-----
```

Save the key to a file, say "dyn_key.pem", and then use it to establish an SSH session.

```text
$ ssh -i dyn_key.pem username@ip
username@ip:~$
```

Creating new keys, saving it in a file and establishing an SSH session will all be done
via a single Vault CLI.

```text
$ vault ssh -role dynamic_key_role username@ip
username@ip:~$
```

### OTP key type

Create a role "otp_key_role" of key type "otp". All the machines represented by CIDR
block should have Vault SSH agent installed and challenge response mechanism enabled
as detailed in https://github.com/hashicorp/vault-ssh-agent.

```text
$ vault write ssh/roles/otp_key_role key_type=otp default_user=bar cidr=x.x.x.x/y
```

Create an OTP 
