#!/usr/bin/env bash
# Copyright IBM Corp. 2016, 2025
# SPDX-License-Identifier: BUSL-1.1

set -euo pipefail

fail() {
  echo "$1" 1>&2
  exit 1
}

logger() {
  DT=$(date '+%Y/%m/%d %H:%M:%S')
  echo "$DT $0: $1"
}

file_name="node_exporter-${PROMETHEUS_NODE_EXPORTER_VERSION}.linux-amd64.tar.gz"
dir_name=$(echo "$file_name" | rev | cut -d '.' -f 3- | rev)
prom_dir="$HOME/$dir_name"

if [ -d "$prom_dir" ]; then
  if pgrep node_exporter > /dev/null; then
    logger "killing prometheus node exporter"
    pkill node_exporter
    sleep 3
  fi

  logger "starting prometheus node exporter"
  "$prom_dir"/node_exporter >> "$prom_dir"/node_exporter.log 2>&1 &
else
  logger "prometheus node exporter couldn't be found"
fi
