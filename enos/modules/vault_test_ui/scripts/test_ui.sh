#!/usr/bin/env bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1


set -eux -o pipefail

project_root=$(git rev-parse --show-toplevel)
pushd "$project_root" > /dev/null

echo "running test-ember-enos"
make test-ember-enos
popd > /dev/null
