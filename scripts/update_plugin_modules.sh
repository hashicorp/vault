#!/bin/sh

set -e

## Make a temp dir
tempdir=$(mktemp -d update-plugin-modules.XXXXXX)

## Set paths
cd $tempdir

## Get Vault
echo "Fetching vault..."
git clone https://github.com/hashicorp/vault

for plugin in $(grep github.com/hashicorp/vault-plugin- vault/go.mod | cut -f 2 | cut -d ' ' -f 1 | cut -d '/' -f 3)
do
	echo "Fetching $plugin..."
	git clone https://github.com/hashicorp/$plugin
	cd $plugin
	rm -rf vendor
	go get github.com/hashicorp/vault/api@master
	go mod tidy
	go mod vendor
	git add .
	git commit --allow-empty -m "Updating vault dep"
	if [ ! -z $PUSH_COMMITS ]
	then
		git push
	fi
	cd ..
	cd vault
	go get github.com/hashicorp/$plugin@master
	cd ..
done

cd vault
go mod tidy
rm -rf vendor
go mod vendor
git add .
git commit --allow-empty -m "Updating plugin deps"
if [ ! -z $PUSH_COMMITS ]
then
	git push
fi
