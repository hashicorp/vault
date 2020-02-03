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
govendor fetch -v +missing

# Clean up after the logrus mess
govendor remove -v github.com/Sirupsen/logrus
cd vendor
find -type f | grep '.go' | xargs sed -i -e 's/Sirupsen/sirupsen/'

# Need the v2 branch for Azure
govendor fetch -v github.com/coreos/go-oidc@v2

# Need the v3 branch for dockertest
govendor fetch -v github.com/ory/dockertest@v3

# Current influx master is alpha, pin to v1.7.3
govendor fetch github.com/influxdata/influxdb/client/v2@v1.7.4
govendor fetch github.com/influxdata/influxdb/models@v1.7.4
govendor fetch github.com/influxdata/influxdb/pkg/escape@v1.7.4

# Current circonus needs v3
grep circonus-gometrics vendor.json | cut -d '"' -f 4 | while read -r i; do govendor fetch $i@v2; done

# API breakage
govendor fetch github.com/satori/go.uuid@f58768cc1a7a7e77a3bd49e98cdd21419399b6a3

echo "Done; to commit run \n\ncd ${GOPATH}/src/github.com/hashicorp/${TOOL}\n"
