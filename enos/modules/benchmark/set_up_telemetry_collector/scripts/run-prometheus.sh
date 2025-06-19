#!/usr/bin/env bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

set -euo pipefail

pkill prometheus || true
sleep 2
prom_dir="$HOME/prom"
"$prom_dir"/prometheus --web.enable-remote-write-receiver --config.file="$prom_dir/prometheus.yml" >> "$prom_dir"/prom.log 2>&1 &
