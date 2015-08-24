---
layout: "docs"
page_title: "Secret Backend: SSH"
sidebar_current: "docs-secrets-ssh"
description: |-
  The SSH secret backend for Vault generates dynamic SSH keys or One-Time-Passwords. 
---

# SSH Secret Backend

Name: `ssh`

Vault SSH backend generates SSH credentials for remote hosts dynamically. This
backend increases the security by removing the need to share the private key to
everyone who needs access to infrastructures. It also solves the problem of
management and distribution of keys belonging to remote hosts.

This backend supports two types of credential creation: Dynamic and OTP. Both of
them addresses the problems in different ways.

Read and carefully understand both of them and choose the one which best suits
your needs.

This page will show a quick start for this backend. For detailed documentation
on every path, use `vault path-help` after mounting the backend.

----------------------------------------------------
## I. Dynamic Type

Register the shared secret key (having super user privileges) with Vault and let
Vault take care of issuing a dynamic secret key every time a client wants to SSH
into the remote host.

When a Vault authenticated client requests for a dynamic credential, Vault server
creates a key-pair, uses the previously shared secret key to login to the remote
host and appends the newly generated public key to `~/.ssh/authorized_keys` file for 
the desired username. Vault uses an install script (configurable) to achieve this.
To run this script in super user mode without password prompts, `NOPASSWD` option
for sudoers should be enabled at all remote hosts.

File: `/etc/sudoers`

```hcl
%sudo   ALL=(ALL)NOPASSWD: ALL 
```

The private key returned to the user will be leased and can be renewed if desired.
Once the key is given to the user, Vault will not know when it gets used or how many
time it gets used. Therefore, Vault **WILL NOT** and cannot audit the SSH session
establishments. An alternative is to use OTP type, which audits every SSH request
(see below).

### Mounting SSH

`ssh` backend is not mounted by default. So, the first step in using the SSH backend
is to mount it.

```shell
$ vault mount ssh
Successfully mounted 'ssh' at 'ssh'!
```

Next, we must register infrastructures with Vault. This is done by writing the role
information. The type of credentials created are determined by the `key_type` option.
To do this, first create a named key and then create a role.

### Registering shared secret key

Create a named key, say `dev_key`, which represents a registered shared private key.
Remember that this key should be of admin user with super user privileges.

```shell
$ vault write ssh/keys/dev_key key=@dev_shared_key.pem
```

### Create a Role

Create a role, say `dynamic_key_role`. All the machines represented by CIDR block
should be accessible through `dev_key` with root privileges.

```shell
$ vault write ssh/roles/dynamic_key_role key_type=dynamic key=dev_key admin_user=username default_user=username cidr_list=x.x.x.x/y
Success! Data written to: ssh/roles/dynamic_key_role
```

Use the `install_script` option to provide an install script if hosts does not
resemble typical Linux machine. The default script is compiled into the binary.
It is straight forward and is shown below. The script takes three arguments which
are explained in the comments.

```shell
# This script file installs or uninstalls an RSA public key to/from authoried_keys
# file in a typical linux machine. This script should be registered with vault
# server while creating a role for key type 'dynamic'.

# $1: "install" or "uninstall"
#
# $2: File name containing public key to be installed. Vault server uses UUID
# as file name to avoid collisions with public keys generated for requests.
#
# $3: Absolute path of the authorized_keys file.

if [ $1 != "install" && $1 != "uninstall" ]; then
	exit 1
fi

# If the key being installed is already present in the authorized_keys file, it is
# removed and the result is stored in a temporary file.
grep -vFf $2 $3 > temp_$2

# Contents of temporary file will be the contents of authorized_keys file.
cat temp_$2 | sudo tee $3

if [ $1 == "install" ]; then
# New public key is appended to authorized_keys file
cat $2 | sudo tee --append $3
fi

# Auxiliary files are deleted
rm -f $2 temp_$2
```

### Create a credential

Create a dynamic key for an IP that belongs to `dynamic_key_role`.

```shell
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

Save the key to a file, say `dyn_key.pem`, and then use it to establish an SSH session.

```shell
$ ssh -i dyn_key.pem username@ip
username@ip:~$
```

### Automate it!

Creation of new key, saving it in a file and establishing an SSH session will all be done
via a single Vault CLI.

```shell
$ vault ssh -role dynamic_key_role username@ip
username@ip:~$
```
----------------------------------------------------
## II. One-Time-Password (OTP) Type

Install Vault SSH Agent in remote hosts and let Vault server issue an OTP every time
a client wants to SSH into remote hosts.

Vault authenticated clients request for a credential from Vault server and get an OTP
issued. When clients try to establish SSH connection with the remote host, OTP typed
in at the password prompt will be received by the Vault agent and gets validated
by the Vault server. Vault server deletes the OTP after validating it once (hence one-time).

Since Vault server is contacted for every successful connection establishment, unlike
Dynamic type, every login attempt **WILL** be audited.

See [Vault-SSH-Agent](https://github.com/hashicorp/vault-ssh-agent) for details
on how to configure the agent.

### Mounting SSH

`ssh` backend is not mounted by default and needs to be explicitly mounted. This is
a common step for both OTP and Dynamic types. 

```shell
$ vault mount ssh
Successfully mounted 'ssh' at 'ssh'!
```

### Creating a Role

Create a role, say `otp_key_role` for key type `otp`. All the machines represented
by CIDR block should have agent installed in them and have their SSH configuration
modified to support Vault SSH Agent client authentication.

```shell
$ vault write ssh/roles/otp_key_role key_type=otp default_user=username cidr_list=x.x.x.x/y,m.m.m.m/n
Success! Data written to: ssh/roles/otp_key_role
```

### Create a Credential

Create an OTP credential for an IP that belongs to `otp_key_role`.

```shell
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

```shell
$ ssh username@localhost
Password: <Enter OTP>
username@ip:~$
```

### Automate it!

Creation of new OTP and running SSH command can be done via a single CLI.

```shell
$ vault ssh -role otp_key_role username@x.x.x.x
OTP for the session is `b4d47e1b-4879-5f4e-ce5c-7988d7986f37`
[Note: Install `sshpass` to automate typing in OTP]
Password: <Enter OTP>
```

OTP will be typed in using `sshpass` if it is installed.

```shell
$ vault ssh -role otp_key_role username@x.x.x.x
username@ip:~$ 
``` 
----------------------------------------------------

## API

### /ssh/config/lease
#### POST

<dl class="api">
  <dt>Description</dt>
  <dd>
    Configures the lease settings for generated credentials.
    This is a root protected endpoint.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/ssh/config/lease`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">lease</span>
        <span class="param-flags">required</span>
        (String)
	The lease value provided as a duration
        with time suffix. Hour is the largest suffix.
      </li>
      <li>
        <span class="param">lease_max</span>
        <span class="param-flags">required</span>
        (String)
	The maximum lease value provided as a duration
        with time suffix. Hour is the largest suffix.
      </li>
    </ul>
  </dd>

  <dt>Returns</dt>
  <dd>
    A `204` response code.
  </dd>
</dl>

### /ssh/keys/
#### POST

<dl class="api">
  <dt>Description</dt>
  <dd>
    Creates or updates a named key. This is a root protected endpoint.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/ssh/keys/<key name>`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">key</span>
        <span class="param-flags">required</span>
        (String)
	SSH private key with super user privileges in host 
      </li>
    </ul>
  </dd>

  <dt>Returns</dt>
  <dd>
    A `204` response code.
  </dd>

#### GET 

<dl class="api">
  <dt>Description</dt>
  <dd>
    Queries a named key. This is a root protected endpoint.
  </dd>

  <dt>Method</dt>
  <dd>GET</dd>

  <dt>URL</dt>
  <dd>`/ssh/keys/<key name>`</dd>

  <dt>Parameters</dt>
  <dd>None</dd>
  
  <dt>Returns</dt>
  <dd>

```javascript
{
		"key": "-----BEGIN RSA PRIVATE KEY-----\nMIIEogIBAAKCAQEAvYvoRcWRxqOim5VZnuM6wHCbLUeiND0yaM1tvOl+Fsrz55DG\nA0OZp4RGAu1Fgr46E1mzxFz1+zY4UbcEExg+u21fpa8YH8sytSWW1FyuD8ICib0A\n/l8slmDMw4BkkGOtSlEqgscpkpv/TWZD1NxJWkPcULk8z6c7TOETn2/H9mL+v2RE\nmbE6NDEwJKfD3MvlpIqCP7idR+86rNBAODjGOGgyUbtFLT+K01XmDRALkV3V/nh+\nGltyjL4c6RU4zG2iRyV5RHlJtkml+UzUMkzr4IQnkCC32CC/wmtoo/IsAprpcHVe\nnkBn3eFQ7uND70p5n6GhN/KOh2j519JFHJyokwIDAQABAoIBAHX7VOvBC3kCN9/x\n+aPdup84OE7Z7MvpX6w+WlUhXVugnmsAAVDczhKoUc/WktLLx2huCGhsmKvyVuH+\nMioUiE+vx75gm3qGx5xbtmOfALVMRLopjCnJYf6EaFA0ZeQ+NwowNW7Lu0PHmAU8\nZ3JiX8IwxTz14DU82buDyewO7v+cEr97AnERe3PUcSTDoUXNaoNxjNpEJkKREY6h\n4hAY676RT/GsRcQ8tqe/rnCqPHNd7JGqL+207FK4tJw7daoBjQyijWuB7K5chSal\noPInylM6b13ASXuOAOT/2uSUBWmFVCZPDCmnZxy2SdnJGbsJAMl7Ma3MUlaGvVI+\nTfh1aQkCgYEA4JlNOabTb3z42wz6mz+Nz3JRwbawD+PJXOk5JsSnV7DtPtfgkK9y\n6FTQdhnozGWShAvJvc+C4QAihs9AlHXoaBY5bEU7R/8UK/pSqwzam+MmxmhVDV7G\nIMQPV0FteoXTaJSikhZ88mETTegI2mik+zleBpVxvfdhE5TR+lq8Br0CgYEA2AwJ\nCUD5CYUSj09PluR0HHqamWOrJkKPFPwa+5eiTTCzfBBxImYZh7nXnWuoviXC0sg2\nAuvCW+uZ48ygv/D8gcz3j1JfbErKZJuV+TotK9rRtNIF5Ub7qysP7UjyI7zCssVM\nkuDd9LfRXaB/qGAHNkcDA8NxmHW3gpln4CFdSY8CgYANs4xwfercHEWaJ1qKagAe\nrZyrMpffAEhicJ/Z65lB0jtG4CiE6w8ZeUMWUVJQVcnwYD+4YpZbX4S7sJ0B8Ydy\nAhkSr86D/92dKTIt2STk6aCN7gNyQ1vW198PtaAWH1/cO2UHgHOy3ZUt5X/Uwxl9\ncex4flln+1Viumts2GgsCQKBgCJH7psgSyPekK5auFdKEr5+Gc/jB8I/Z3K9+g4X\n5nH3G1PBTCJYLw7hRzw8W/8oALzvddqKzEFHphiGXK94Lqjt/A4q1OdbCrhiE68D\nMy21P/dAKB1UYRSs9Y8CNyHCjuZM9jSMJ8vv6vG/SOJPsnVDWVAckAbQDvlTHC9t\nO98zAoGAcbW6uFDkrv0XMCpB9Su3KaNXOR0wzag+WIFQRXCcoTvxVi9iYfUReQPi\noOyBJU/HMVvBfv4g+OVFLVgSwwm6owwsouZ0+D/LasbuHqYyqYqdyPJQYzWA2Y+F\n+B6f4RoPdSXj24JHPg/ioRxjaj094UXJxua2yfkcecGNEuBQHSs=\n-----END RSA PRIVATE KEY-----\n"
}
```
  </dd>


#### DELETE 

<dl class="api">
  <dt>Description</dt>
  <dd>
    Deletes a named key. This is a root protected endpoint.
  </dd>

  <dt>Method</dt>
  <dd>DELETE</dd>

  <dt>URL</dt>
  <dd>`/ssh/keys/<key name>`</dd>

  <dt>Parameters</dt>
  <dd>None</dd>

  <dt>Returns</dt>
  <dd>
    A `204` response code.
  </dd>

### /ssh/roles/
#### POST

<dl class="api">
  <dt>Description</dt>
  <dd>
    Creates or updates a named role.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/ssh/roles/<role name>`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">key</span>
        <span class="param-flags">required for Dynamic type, NA for OTP type</span>
	(String)
        Name of the registered key in Vault. Before creating the role, use
        the `keys/` endpoint to create a named key.
      </li>
      <li>
        <span class="param">admin_user</span>
        <span class="param-flags">required for Dynamic type, NA for OTP type</span>
	(String)
	Admin user at remote host. The shared key being registered should be
	for this user and should have root privileges. Everytime a dynamic 
	credential is being generated for other users, Vault uses this admin
	username to login to remote host and install the generated credential
	for the other user.
      </li>
      <li>
        <span class="param">default_user</span>
        <span class="param-flags">required for both types</span>
	(String)
	Default username for which a credential will be generated.
	When the endpoint 'creds/' is used without a username, this
	value will be used as default username.
      </li>
      <li>
        <span class="param">cidr_list</span>
        <span class="param-flags">required for both types</span>
	(String)
	Comma separated list of CIDR blocks for which the role is applicable for.
	CIDR blocks can belong to more than one role.
      </li>
      <li>
        <span class="param">port</span>
        <span class="param-flags">optional for both types</span>
	(Integer)
	Port number for SSH connection. Default is '22'. Port number does not
	play any role in creation of OTP. For 'otp' type, this is just a way
	to inform client about the port number to use. Port number will be
	returned to client by Vault server along with OTP.
      </li>
      <li>
        <span class="param">key_type</span>
        <span class="param-flags">required for both types</span>
	(String)
	Type of key used to login to hosts. It can be either `otp` or `dynamic`.
	`otp` type requires agent to be installed in remote hosts.
      </li>
      <li>
        <span class="param">key_bits</span>
        <span class="param-flags">optional for Dynamic type, NA for OTP type</span>
	(Integer)
	Length of the RSA dynamic key in bits. It is 1024 by default or it can be 2048.
      </li>
      <li>
        <span class="param">install_script</span>
        <span class="param-flags">optional for Dynamic type, NA for OTP type</span>
	(String)
	Script used to install and uninstall public keys in the target machine.
	The inbuilt default install script will be for Linux hosts.
      </li>
      <li>
        <span class="param">allowed_users</span>
        <span class="param-flags">optional for both types</span>
	(String)
	If this option is not specified, client can request for a credential for
	any valid user at the remote host, including the admin user. If only certain
	usernames are to be allowed, then this list enforces it. If this field is
	set, then credentials can only be created for default_user and usernames
	present in this list.
      </li>
    </ul>
  </dd>

  <dt>Returns</dt>
  <dd>
    A `204` response code.
  </dd>

#### GET 

<dl class="api">
  <dt>Description</dt>
  <dd>
    Queries a named role.
  </dd>

  <dt>Method</dt>
  <dd>GET</dd>

  <dt>URL</dt>
  <dd>`/ssh/roles/<role name>`</dd>

  <dt>Parameters</dt>
  <dd>None</dd>
  
  <dt>Returns</dt>
  <dd>For dynamic role:

```json
{
		"admin_user": "username",
		"cidr_list": "x.x.x.x/y",
		"default_user": "username",
		"key": "<key name>",
		"key_type": "dynamic",
		"port": 22
}
```
  </dd>

  <dd>For OTP role:

```json
{
		"cidr_list": "x.x.x.x/y",
		"default_user": "username",
		"key_type": "otp",
		"port": 22
}
```
  </dd>


#### DELETE 

<dl class="api">
  <dt>Description</dt>
  <dd>
    Deletes a named role.
  </dd>

  <dt>Method</dt>
  <dd>DELETE</dd>

  <dt>URL</dt>
  <dd>`/ssh/roles/<role name>`</dd>

  <dt>Parameters</dt>
  <dd>None</dd>

  <dt>Returns</dt>
  <dd>
    A `204` response code.
  </dd>
### /ssh/creds/
#### POST

<dl class="api">
  <dt>Description</dt>
  <dd>
    Creates a credential for a specific username and IP under the given role.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/ssh/creds/<role name>`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">username</span>
        <span class="param-flags">optional</span>
	(String)
        Username in remote host.
      </li>
      <li>
        <span class="param">ip</span>
        <span class="param-flags">required</span>
	(String)
        IP of the remote host.
      </li>
    </ul>
  </dd>
  
  <dt>Returns</dt>
  <dd>
    A `204` response code.
  </dd>

### /ssh/lookup
#### POST

<dl class="api">
  <dt>Description</dt>
  <dd>
    Lists all the roles given IP is associated with.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/ssh/lookup`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">ip</span>
        <span class="param-flags">required</span>
	(String)
        IP of the remote host.
      </li>
    </ul>
  </dd>

  <dt>Returns</dt>
  <dd>
    A `204` response code.
  </dd>

### /ssh/verify
#### POST

<dl class="api">
  <dt>Description</dt>
  <dd>
    Verifies if the given OTP is valid. This is an unauthenticated endpoint.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/ssh/verify`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">otp</span>
        <span class="param-flags">required</span>
	(String)
        One-Time-Key that needs to be validated.
      </li>
    </ul>
  </dd>

  <dt>Returns</dt>
  <dd>
    A `204` response code.
  </dd>

