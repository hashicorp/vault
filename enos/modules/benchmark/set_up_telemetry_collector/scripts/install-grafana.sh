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

[[ -z "${GRAFANA_VERSION}" ]] && fail "GRAFANA_VERSION env variable has not been set"
[[ -z "${RETRY_INTERVAL}" ]] && fail "RETRY_INTERVAL env variable has not been set"
[[ -z "${TIMEOUT_SECONDS}" ]] && fail "TIMEOUT_SECONDS env variable has not been set"

install_grafana() {
  if command -v grafana-server &> /dev/null; then
    logger "grafana already installed"
    return 0
  fi

  logger "installing grafana"
  cd "$HOME"

  file_name="grafana-enterprise_${GRAFANA_VERSION}_amd64.deb"
  sudo apt-get update
  sudo apt-get install -y adduser libfontconfig1 musl
  wget "https://dl.grafana.com/enterprise/release/$file_name"
  sudo dpkg -i "$file_name"

  prom_ds="/etc/grafana/provisioning/datasources/prometheus.yaml"
  dash_config="/etc/grafana/provisioning/dashboards/benchmark.yaml"
  dash_dir="/etc/grafana/dashboards"

  logger "writing out grafana datasource for prometheus"
  sudo tee "$prom_ds" << EOF
  apiVersion: 1

  datasources:
  - name: default-prom
    type: prometheus
    access: proxy
    orgId: 1
    url: http://localhost:9090
    isDefault: true
    version: 1
    editable: true
EOF
  sudo chown root:grafana "$prom_ds"

  logger "removing sample.yaml dashboard config"
  sudo rm -f /etc/grafana/provisioning/dashboards/sample.yaml

  logger "writing out grafana dashboard config"
  sudo mkdir -p "$dash_dir"
  sudo tee "$dash_config" << EOF
  apiVersion: 1

  providers:
  - name: 'vault'
    orgId: 1
    folder: 'Vault dashboards'
    type: file
    options:
      path: /etc/grafana/dashboards
EOF
  sudo chown root:grafana "$dash_config"
  sudo chown root:grafana "$dash_dir"
}

begin_time=$(date +%s)
end_time=$((begin_time + TIMEOUT_SECONDS))
while [[ "$(date +%s)" -lt "${end_time}" ]]; do
  if install_grafana; then
    exit 0
  fi

  sleep "${RETRY_INTERVAL}"
done

fail "Timed out waiting for grafana to install"
