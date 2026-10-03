#!/usr/bin/env bash
# Copyright IBM Corp. 2016, 2025
# SPDX-License-Identifier: BUSL-1.1

set -euo pipefail

script="$${1:-}"

if [[ -z "$${script}" ]]; then
  echo "Usage: $0 approle-login | kvv1 | kvv2 | lease-revocation"
  exit 1
fi

K6_PROMETHEUS_RW_SERVER_URL=http://${metrics_addr}:9090/api/v1/write k6 run -o experimental-prometheus-rw scripts/k6-$${script}.js
