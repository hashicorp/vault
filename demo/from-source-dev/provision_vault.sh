#!/bin/bash
# make sure the os is up to date on patches
echo "###"
echo "### make sure the os is up to date on patches"
echo "###"
sudo apt-get update -y
echo "###"
echo "### some of the default intro examples/docs use this tool, so include it in the build"
echo "###"
sudo apt-get install jq
# this package is not good enough. :( go is not a high enough version from this package (as of 2015.10.01)
#sudo apt-get install golang-go -y
# so go get a higher version from the go project and build it
echo "###"
echo "### get go from https://storage.googleapis.com/golang/go1.5.1.linux-amd64.tar.gz"
echo "###"
cd ~
# maybe this location should be moved to /vagrant to prevent repeated downloads?
# look to see if we have the file before downloading it again?
echo "### cd ~"
curl -O https://storage.googleapis.com/golang/go1.5.1.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.5.1.linux-amd64.tar.gz
echo "export PATH=$PATH:/usr/local/go/bin" >> ~/.bashrc
export PATH=$PATH:/usr/local/go/bin
# setting GOBIN is not requried after go 1.0, it defaults to $GOROOT/bin
#echo "export export GOBIN='/usr/bin/go'" >> ~/.bashrc
echo "export GOPATH=$HOME/work" >> ~/.bashrc
# this added a env to allow ""vault status" to resolve the localhost vault server
echo "export VAULT_ADDR=http://127.0.0.1:8200" >> ~/.bashrc
export VAULT_ADDR=http://127.0.0.1:8200
echo "set VAULT_ADDR in .bashrc"
export GOPATH=$HOME/work
echo "export PATH=$PATH:$GOPATH/bin" >> ~/.bashrc
export PATH=$PATH:$GOPATH/bin
mkdir -p  $GOPATH/src/github.com/hashicorp/vault/
cd $GOPATH
echo "###"
echo "### install git client"
echo "###"
sudo apt-get install git -y
echo "###"
echo "### fetch the current vault git repo"
echo "###"
git clone https://github.com/hashicorp/vault $GOPATH/src/github.com/hashicorp/vault
cd $GOPATH/src/github.com/hashicorp/vault
echo "###"
echo "### bootstrap the environment for the project"
echo "###"
make bootstrap
echo "###"
echo "### LAST build a dev build of the repo"
echo "###"
make dev
echo "###"
echo "### Copy a possible Vault server config file into /etc/vault.conf"
echo "###"
sudo cp /vagrant/etc_vault.conf /etc/vault.conf
echo "###"
echo "### Record when the env was provisioned in /etc/vagrant_provisioned_at"
echo "###"
sudo touch /etc/vagrant_provisioned_at
sudo chown vagrant:vagrant /etc/vagrant_provisioned_at
sudo chmod u=rw,go=r /etc/vagrant_provisioned_at
date >> /etc/vagrant_provisioned_at
