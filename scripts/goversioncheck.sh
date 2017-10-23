#!/usr/bin/env bash

GO_VERSION_MIN=$1
echo "==> Checking that build is using go version >= $1..."

GO_VERSION=$(go version | grep -o 'go[0-9]\+\.[0-9]\+\(\.[0-9]\+\)\?' | tr -d 'go')


IFS="." read -r -a GO_VERSION_ARR <<< "$GO_VERSION"
IFS="." read -r -a GO_VERSION_REQ <<< "$GO_VERSION_MIN"

if [[ ${GO_VERSION_ARR[0]} -lt ${GO_VERSION_REQ[0]} ||
    ( ${GO_VERSION_ARR[0]} -eq ${GO_VERSION_REQ[0]} &&
      ( ${GO_VERSION_ARR[1]} -lt ${GO_VERSION_REQ[1]} ||
        ( ${GO_VERSION_ARR[1]} -eq ${GO_VERSION_REQ[1]} && ${GO_VERSION_ARR[2]} -lt ${GO_VERSION_REQ[2]} )))
    ]]; then
    echo "Vault requires go $GO_VERSION_MIN to build; found $GO_VERSION."
    exit 1
fi
