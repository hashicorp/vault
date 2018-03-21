# Vagrant Vault Demo

This demo provides a `Vagrantfile` that creates a simple Vault server instance with [filesystem based storage backend](https://www.vaultproject.io/docs/configuration/storage/filesystem.html) and available on the network at the IPv4 address *172.20.20.20*.

Getting started should be straightforward — make sure both [VirtualBox and](https://www.virtualbox.org/) and Vagrant are installed for your host system, you have cloned the [Vault GitHub repository](https://github.com/hashicorp/vault), and then open a terminal:

```
$ cd $GIT_CHECKOUT_DIR/demo/vagrant
$ vagrant up
```

where `$GIT_CHECKOUT_DIR` is the directory you cloned the Vault git repository into.

> NOTE: If you prefer a different Vagrant box, you can set the `DEMO_BOX_NAME`
> environment variable before starting `vagrant` like this:
> `DEMO_BOX_NAME="ubuntu/xenial64" vagrant up`

Once Vagrant is finished, you can check the status:

```
$ vagrant status
Current machine states:

v1                        running (virtualbox)
```

## Use Vault from the VM via ssh

At this point the Vault server's virtual machine is running and you can use it. One way to do so is to `ssh` into the VM and then export `VAULT_ADDR` specifying the loopback address for communication with Vault:

```
$ vagrant ssh v1
$ export VAULT_ADDR=http://127.0.0.1:8200
```

Then check Vault version and status:

```
$ vault version
Vault v0.9.5 ('36edb4d42380d89a897e7f633046423240b710d9')

$ vault status
Error checking seal status: Error making API request.

URL: GET http://127.0.0.1:8200/v1/sys/seal-status
Code: 400. Errors:

* server is not yet initialized
```

This output is expected, and means that the Vault server is running, but has not yet been initialized.

You can now initialize, unseal, and use Vault.

Initialize Vault:

```
$ vault operator init
Unseal Key 1: NrNq9ESXpifDAR+L3WSiTFHIVQPczpTDwGHBhLa08MJ4
Unseal Key 2: o5apSWUeDbwPK++w1bqlGvqUM/YM6DeB/V5BHzMyagoG
Unseal Key 3: 02luET3dhhgsPmC0uer1zDT74QUsSC3qUFUrLFy03cVd
Unseal Key 4: LK3ejDsGnfWStufKgTGcpJ/aHgCtXnJ+BZB/is97mkur
Unseal Key 5: h+NZ4Qc7aBUkXpMfmqtvWudmqK/XnI97hjr0yxormZWt

Initial Root Token: 715a8534-89c1-4270-f5ea-874994cc5702

Vault initialized with 5 key shares and a key threshold of 3. Please securely
distribute the key shares printed above. When the Vault is re-sealed,
restarted, or stopped, you must supply at least 3 of these keys to unseal it
before it can start servicing requests.

Vault does not store the generated master key. Without at least 3 key to
reconstruct the master key, Vault will remain permanently sealed!

It is possible to generate new unseal keys, provided you have a quorum of
existing unseal keys shares. See "vault rekey" for more information.
```

Unseal Vault:

```
$ vault operator unseal NrNq9ESXpifDAR+L3WSiTFHIVQPczpTDwGHBhLa08MJ4
Key                Value
---                -----
Seal Type          shamir
Sealed             true
Total Shares       5
Threshold          3
Unseal Progress    1/3
Unseal Nonce       56b2d361-aa29-4945-4089-b83de2e1e910
Version            0.9.5
HA Enabled         true
$ vault operator unseal o5apSWUeDbwPK++w1bqlGvqUM/YM6DeB/V5BHzMyagoG
Key                Value
---                -----
Seal Type          shamir
Sealed             true
Total Shares       5
Threshold          3
Unseal Progress    2/3
Unseal Nonce       56b2d361-aa29-4945-4089-b83de2e1e910
Version            0.9.5
HA Enabled         true
$ vault operator unseal 02luET3dhhgsPmC0uer1zDT74QUsSC3qUFUrLFy03cVd
Key             Value
---             -----
Seal Type       shamir
Sealed          false
Total Shares    5
Threshold       3
Version         0.9.5
Cluster Name    vault-cluster-2393d85a
Cluster ID      7f5dc0bf-fae8-160a-88de-367c692fe2a9
HA Enabled      false
```

Login with the initial root token that was presented when Vault was initialized:

```
$ vault login 715a8534-89c1-4270-f5ea-874994cc5702
is already stored in the token helper. You do NOT need to run "vault login"
again. Future Vault requests will automatically use this token.

Key                Value
---                -----
token              715a8534-89c1-4270-f5ea-874994cc5702
token_accessor     64d530a5-f8e8-79db-45d8-e3af52f9d523
token_duration     ∞
token_renewable    false
token_policies     [root]
```

Write a secret:

```
$ vault write secret/foo zip=zap
Success! Data written to: secret/foo
```

## User Vault from the Host

If you'd rather not have to login to the VM via `ssh`, you can also install a Vault binary on the host system and export the `VAULT_ADDR` environment variable to directly access Vault from the host system:

```
$ export VAULT_ADDR=http://172.20.20.20:8200
$ vault status
Key             Value
---             -----
Seal Type       shamir
Sealed          false
Total Shares    5
Threshold       3
Version         0.9.5
Cluster Name    vault-cluster-2393d85a
Cluster ID      7f5dc0bf-fae8-160a-88de-367c692fe2a9
HA Enabled      false
```

If you instead receive an error about the server not being initialized, follow the steps above to initialize, unseal, and login into your Vault VM.

## Where to Next?

If you're new, you can learn more about Vault from the [Getting Started guide](https://www.vaultproject.io/intro/getting-started/install.html), or check out other available [Vault Guides](https://www.vaultproject.io/guides/index.html) for deep dives into particular Vault topics.

Of course the [Vault documentation](https://www.vaultproject.io/docs/index.html) and the [Vault HTTP API documentation](https://www.vaultproject.io/api/index.html) are helpful resources at all times and great to explore as well. Be sure to review the links int he [Vault GitHub README](https://github.com/hashicorp/vault) for additional Vault resources, too!

> NOTE: This demo will use the latest Vault open source release by default,
> but if you need a different Vault version, set the `Vault_DEMO_VERSION`
> environment variable before `vagrant up` like this:
> `VAULT_DEMO_VERSION=0.9.2 vagrant up`

## Resources

1. [Vault website](https://www.vaultproject.io/)
2. [Vault documentation](https://www.vaultproject.io/docs/index.html)
3. [Filesystem Storage Backend](https://www.vaultproject.io/docs/configuration/storage/filesystem.html)
4. [Vault GitHub repository](https://github.com/hashicorp/vault)
5. [VirtualBox](https://www.virtualbox.org/)
6. [Vagrant](https://www.vagrantup.com/)
7. [Getting Started guide](https://www.vaultproject.io/intro/getting-started/install.html)
8. [Vault Guides](https://www.vaultproject.io/guides/index.html)
9. [Vault HTTP API documentation](https://www.vaultproject.io/api/index.html)
