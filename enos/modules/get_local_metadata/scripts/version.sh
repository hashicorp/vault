#!/bin/env bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

set -eu -o pipefail

pushd "$(git rev-parse --show-toplevel)" > /dev/null
make ci-get-version
popd > /dev/null
