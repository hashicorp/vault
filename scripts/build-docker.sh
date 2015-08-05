#!/bin/bash
set -e

if [ -z "$1" ]; then
    OS_PLATFORM_ARG=(-os="darwin linux windows")
else
    OS_PLATFORM_ARG=($1)
fi

if [ -z "$2" ]; then
    OS_ARCH_ARG=(-arch="386 amd64 arm")
else
    OS_ARCH_ARG=($2)
fi

# Build Docker image unless we opt out of it
if [[ -z "$SKIP_BUILD" ]]; then
   docker build --rm=true --force-rm=true -t vault-builder -f Dockerfile.build .
fi

# Get rid of existing binaries
rm -f *-386
rm -f *-amd64
rm -f dist/*
docker run --rm -v `pwd`:/go/src/github.com/hashicorp/vault:Z vault-builder gox "${OS_PLATFORM_ARG[@]}" "${OS_ARCH_ARG[@]}" -output="dist/{{.Dir}}_{{.OS}}-{{.Arch}}" -ldflags="-w"
