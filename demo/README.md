# Why a Demo?

This is a set of demo "vagrant based hosts/environments" to provide quick starting points for using Vault. You will need [Vagrant](https://vagrantup.com) installed and working with [VirtualBox](https://virtualbox.org/). ( Other VM providers should work too, but for now these examples are designed only for Virtualbox.)

* Before you do anything else:

	* Install Vagrant (https://vagrantup.com)
	* Install VirtualBox (https://virtualbox.org/)
	* Reboot if (or when) you are asked to.

-----
Now to get started, you select which demo you want. (Please see the README's in the respective directories for details.)
-----

* ../from-source-dev
	* This is a single vm(ubuntu/trusty64) with a standalone Vault server in "demo mode" built from the most current source code.
	* This is intended to be the fastest, and simplest way to PLAY with vault.
	* This should not ever be used for any real secrets. This is a demo/experment environment and is NOT SECURE or production ready.
	* It is however built from the latest source and should be an easy path to put your hands on vault for the first time. :)

* ../HA-example
  * (This is TBD.)
	* This is a set of 4 seperate VM's. There are two Vault hosts that uses two Consul hosts to provide HA support.
	* The two Vault hosts are exposed to only the guest OS by default.
			(The Consul hosts are only exposed to the Vault hosts.)

# Basic howto get the party started? :
* Start a command line shell. (If you are not already in one.)
* cd <into this the selected directory>
* vagrant up
