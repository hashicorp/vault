#What is this configuration about?

This is a single vm(ubuntu/trusty64) with a standalone Vault server in "demo mode" built from the most current source code.
This is intended to be the fastest, and simplest way to PLAY with vault.
This should not ever be used for any real secrets. This is a demo/experment environment and is NOT SECURE or production ready.
It is however built from the latest source and should be an easy path to put your hands on vault for the first time. :)

Howto get the party started?:
  Install the required host OS software ( https://vagrantup.com and https://virtualbox.org/)
  cd <into this directory>
  vagrant up
  vagrant ssh
  # start the server in demo mode and route stdout to a log file in /tmp
  #   note the lets you keep using the same shell. You can also opt to not do the redirect and start a second vagrant ssh shell
  vault server -dev > /tmp/vault.dev.demo.log &
  # then enter any other vault commands you want to....s
  vault status
  vault mounts
  vault write secret/hello value=world excited=yes
  vault read secret/hello
  ...
