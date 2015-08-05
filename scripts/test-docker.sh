#!/bin/bash
set -e

# Build Docker image unless we opt out of it
if [[ -z "$SKIP_BUILD" ]]; then
   docker build --rm=true --force-rm=true -t vault-builder -f Dockerfile.build .
fi

# Get rid of existing binaries
docker run --rm -e "TF_ACC=" -v `pwd`:/go/src/github.com/hashicorp/vault:Z vault-builder make test
