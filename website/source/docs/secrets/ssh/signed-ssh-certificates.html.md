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
complexity and in terms of being platform agnostic.

When using this type, an SSH CA signing key is generated or configured at the
backend's mount. This key will be used to sign other SSH keys. The private half
of the signing key always stays within Vault and the public half is exposed via
the API. Each mount of this backend represents a unique signing key pair. It is
recommended that the host keys and client keys are signed using different mounts
of this backend.

This page will show a quick start for this backend. For detailed documentation
on every path, use `vault path-help` after mounting the backend.

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
* If you encounter an error on the client indicating `"no separate private key
  for certificate"` you may be hitting a bug introduced into OpenSSH version
  7.2 and fixed in 7.5. See [OpenSSH bug
  2617](https://bugzilla.mindrot.org/show_bug.cgi?id=2617) for details. If you
  are able to find a workaround without patching OpenSSH, please submit a PR to
  update this documentation!

## API

The SSH secret backend has a full HTTP API. Please see the
[SSH secret backend API](/api/secret/ssh/index.html) for more
details.
