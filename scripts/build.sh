#!/usr/bin/env bash
# Copyright IBM Corp. 2016, 2026
# SPDX-License-Identifier: BUSL-1.1

#
# This script builds the application from source for multiple platforms.
set -e

GO_CMD=${GO_CMD:-go}

# Get the parent directory of where this script is.
SOURCE="${BASH_SOURCE[0]}"
SOURCE_DIR=$( dirname "$SOURCE" )
while [ -h "$SOURCE" ] ; do SOURCE="$(readlink "$SOURCE")"; done
DIR="$( cd -P "$SOURCE_DIR/.." && pwd )"

# Change into that directory
cd "$DIR"

# Set build tags
BUILD_TAGS="${BUILD_TAGS:-"vault"}"

# BUILD_SOURCE selects the entrypoint package. Defaults to the full server binary.
BUILD_SOURCE="${BUILD_SOURCE:-"./command/server/cmd"}"

# Entrypoint packages are named cmd, so the product name must be explicit.
BIN_NAME="${BIN_NAME:-"vault"}"

# Get the git commit
GIT_COMMIT="$("$SOURCE_DIR"/ci-helper.sh revision)"
GIT_DIRTY="$(test -n "`git status --porcelain`" && echo "+CHANGES" || true)"

BUILD_DATE="$("$SOURCE_DIR"/ci-helper.sh date)"

GOPATH=${GOPATH:-$(${GO_CMD} env GOPATH)}
case $(uname) in
    CYGWIN*)
        GOPATH="$(cygpath $GOPATH)"
        ;;
esac

# Delete the old binary
echo "==> Removing old binary..."
rm -f "bin/${BIN_NAME}"
rm -rf pkg/*
mkdir -p bin/

# Build!
echo "==> Building..."
${GO_CMD} build \
    -gcflags "${GCFLAGS}" \
    -ldflags "${LD_FLAGS} -X github.com/hashicorp/vault/version.GitCommit='${GIT_COMMIT}${GIT_DIRTY}' -X github.com/hashicorp/vault/version.BuildDate=${BUILD_DATE}" \
    -o "bin/${BIN_NAME}" \
    -tags "${BUILD_TAGS}" \
    "${BUILD_SOURCE}"

# Move all the compiled things to the $GOPATH/bin
OLDIFS=$IFS
IFS=: FIRST=($GOPATH) BIN_PATH=${GOBIN:-${FIRST}/bin}
IFS=$OLDIFS

# Ensure the go bin folder exists
mkdir -p ${BIN_PATH}
rm -f "${BIN_PATH}/${BIN_NAME}"
cp "bin/${BIN_NAME}" "${BIN_PATH}"

# Done!
echo
echo "==> Results:"
ls -hl bin/
