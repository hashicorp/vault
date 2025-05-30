#!/usr/bin/env bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

set -euo pipefail

sudo sed -i '$a client_addr = "0.0.0.0"\n' /etc/consul.d/consul.hcl
sudo sed -i '$a telemetry {\n  prometheus_retention_time = "24h"\n  disable_hostname = true\n}' /etc/consul.d/consul.hcl
