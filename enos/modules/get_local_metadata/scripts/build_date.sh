#!/usr/bin/env bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

set -eu -o pipefail

pushd "$(git rev-parse --show-toplevel)" > /dev/null
make ci-get-date
popd > /dev/null
