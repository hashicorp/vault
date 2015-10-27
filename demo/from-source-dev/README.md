#What is this configuration about?

This is a single vm(ubuntu/trusty64) with a standalone Vault server in "demo mode" built from the most current source code.
This is intended to be the fastest, and simplest way to PLAY with vault.
This should not ever be used for any real secrets. This is a demo/experment environment and is NOT SECURE or production ready.
It is however built from the latest source and should be an easy path to put your hands on vault for the first time. :)

* Howto get the party started?:
  * Install the required host OS software ( [vagrantup](https://vagrantup.com) and [Virtualbox](https://virtualbox.org/) )
  * cd <into this directory>
  * vagrant up
 <p>Note: The first "vagrant up" takes a few (6+) minutes. It downloads are needed dependencies, current vault source code and compiles it. The duration depends on the speed of your host, guest OS and network. </p>
  * vagrant ssh
<p><div>
# Start the server in demo mode and route stdout to a log file in /tmp<br/>
# This lets you keep using the same shell. <br/>
# You can also opt to not do the redirect and start a second vagrant ssh shell
</div></p>
 * vault server -dev > /tmp/vault.dev.demo.log &
<p><div># Then enter any other vault commands you want to....</div></p>
  * vault status
  * vault mounts
  * vault write secret/hello value=world excited=yes
  * vault read secret/hello
  * vault read -format=json secret/hello
  * vault read -format=json secret/hello | jq -r .data.excited
  * vault delete secret/hello
  * ...
<p><div>
  #<br/>
  # when your done for the day....<br/>
  #<br/>
  # exit the vagrant ssh shell to the vm<br/>
  #<br/>
</div></p>
    exit
<p><div>
  #<br/>
  # then stop the vm<br/>
  #  This preserves the VM as it is. It just stops the VM.<br/>
  #  You can start again with a "vagrant up" tomorrow. :)<br/>
  #<br/>
  # or destroy it AND DELETE any "data/changes" you made to the VM host<br/>
  #   You can always return the the "starting state" by issuing a "vagrant up" tomorrow. :)<br/>
  # vagrant destroy<br/>
  #<br/>
</div></p>
    vagrant halt
