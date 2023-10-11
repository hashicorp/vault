#!/usr/bin/env bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

echo "==> Checking that code complies with gofmt requirements..."

files=$(echo $1 | xargs)
if [ -n "$files" ]; then
    echo "Checking changed files..."
    gofmt_files="$(echo $1 | grep -v pb.go | grep -v vendor | xargs go run mvdan.cc/gofumpt -l)"
else
    echo "Checking all files..."
    gofmt_files="$(find . -name '*.go' | grep -v pb.go | grep -v vendor | xargs go run mvdan.cc/gofumpt -l)"
fi

if [[ -n "${gofmt_files}" ]]; then
    echo 'gofumpt needs running on the following files:'
    echo "${gofmt_files}"
    echo "You can use the command: \`make fmt\` to reformat code."
    exit 1
fi
