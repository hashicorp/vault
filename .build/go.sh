#!/bin/bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1
set -e

host_arch="$(dpkg --print-architecture)"
host_arch="${host_arch##*-}"
curl -L "https://go.dev/dl/go${GO_VERSION}.linux-${host_arch}.tar.gz" | tar -C /opt -zxv
