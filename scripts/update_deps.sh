#!/bin/sh

set -e

TOOL=vault

## Make a temp dir
tempdir=$(mktemp -d update-${TOOL}-deps.XXXXXX)

## Set paths
export GOPATH="$(pwd)/${tempdir}"
export PATH="${GOPATH}/bin:${PATH}"
cd $tempdir

## Get Vault
mkdir -p src/github.com/hashicorp
cd src/github.com/hashicorp
echo "Fetching ${TOOL}..."
git clone https://github.com/hashicorp/${TOOL}
cd ${TOOL}

## Clean out earlier vendoring
rm -rf Godeps vendor

## Get govendor
go get github.com/kardianos/govendor

## Init
govendor init

## Fetch deps
echo "Fetching deps, will take some time..."
govendor fetch +missing

govendor remove github.com/Sirupsen/logrus
cd vendor
find -type f | grep '.go' | xargs sed -i -e 's/Sirupsen/sirupsen/'

echo "Done; to commit run \n\ncd ${GOPATH}/src/github.com/hashicorp/${TOOL}\n"
