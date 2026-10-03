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

[[ -z "${PROMETHEUS_NODE_EXPORTER_VERSION}" ]] && fail "PROMETHEUS_NODE_EXPORTER_VERSION env variable has not been set"
[[ -z "${RETRY_INTERVAL}" ]] && fail "RETRY_INTERVAL env variable has not been set"
[[ -z "${TIMEOUT_SECONDS}" ]] && fail "TIMEOUT_SECONDS env variable has not been set"

install_prometheus_node_exporter() {
  file_name="node_exporter-${PROMETHEUS_NODE_EXPORTER_VERSION}.linux-amd64.tar.gz"
  dir_name=$(echo "$file_name" | rev | cut -d '.' -f 3- | rev)
  prom_url="https://github.com/prometheus/node_exporter/releases/download/v${PROMETHEUS_NODE_EXPORTER_VERSION}/${file_name}"
  prom_dir="$HOME/$dir_name"
  if [ -d "$prom_dir" ]; then
    logger "prometheus node exporter already downloaded"
  else
    logger "downloading prometheus node exporter"
    cd "$HOME"
    wget "$prom_url"
    tar zxf "$file_name"
  fi
}

begin_time=$(date +%s)
end_time=$((begin_time + TIMEOUT_SECONDS))
while [[ "$(date +%s)" -lt "${end_time}" ]]; do
  if install_prometheus_node_exporter; then
    exit 0
  fi

  sleep "${RETRY_INTERVAL}"
done

fail "Timed out waiting for prometheus node exporter to install"
