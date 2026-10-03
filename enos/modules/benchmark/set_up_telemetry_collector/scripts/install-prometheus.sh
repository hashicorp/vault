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

env

[[ -z "${PROMETHEUS_VERSION}" ]] && fail "PROMETHEUS_VERSION env variable has not been set"
[[ -z "${RETRY_INTERVAL}" ]] && fail "RETRY_INTERVAL env variable has not been set"
[[ -z "${TIMEOUT_SECONDS}" ]] && fail "TIMEOUT_SECONDS env variable has not been set"
[[ -z "${VAULT_1_ADDR}" ]] && fail "VAULT_1_ADDR env variable has not been set"
[[ -z "${VAULT_2_ADDR}" ]] && fail "VAULT_2_ADDR env variable has not been set"
[[ -z "${VAULT_3_ADDR}" ]] && fail "VAULT_3_ADDR env variable has not been set"
[[ -z "${K6_ADDR}" ]] && fail "K6_ADDR env variable has not been set"

install_prometheus() {
  file_name="prometheus-${PROMETHEUS_VERSION}.linux-amd64.tar.gz"
  dir_name=$(echo "$file_name" | rev | cut -d '.' -f 3- | rev)
  prom_url="https://github.com/prometheus/prometheus/releases/download/v${PROMETHEUS_VERSION}/$file_name"
  prom_dir="$HOME/prom"

  if [ -d "$HOME/$prom_dir" ]; then
    logger "prometheus already downloaded"
    return 0
  fi

  logger "downloading prometheus"
  cd "$HOME"
  mkdir -p "$prom_dir"
  wget "$prom_url"
  tar zxf "$file_name"
  cp "$dir_name/prometheus" "$prom_dir"

  logger "writing out prometheus config file"
  tee "$prom_dir/prometheus.yml" << EOF
global:
  scrape_interval:     5s
  evaluation_interval: 5s

scrape_configs:
  - job_name: 'prometheus'
    static_configs:
      - targets: ['localhost:9090']
  - job_name: 'vault-nodes'
    static_configs:
      - targets: ['${VAULT_1_ADDR}:9100']
        labels:
          instance: 'vault-1'
      - targets: ['${VAULT_2_ADDR}:9100']
        labels:
          instance: 'vault-2'
      - targets: ['${VAULT_3_ADDR}:9100']
        labels:
          instance: 'vault-3'
      - targets: ['${K6_ADDR}:9100']
        labels:
          instance: 'k6'
  - job_name: 'vault-metrics'
    metrics_path: "/v1/sys/metrics"
    params:
      format: ['prometheus']
    static_configs:
      - targets: ['${VAULT_1_ADDR}:8200']
        labels:
          instance: 'vault-1'
      - targets: ['${VAULT_2_ADDR}:8200']
        labels:
          instance: 'vault-2'
      - targets: ['${VAULT_3_ADDR}:8200']
        labels:
          instance: 'vault-3'
EOF

  if [[ -n "${CONSUL_1_ADDR:-}" && -n "${CONSUL_2_ADDR:-}" && -n "${CONSUL_3_ADDR:-}" ]]; then
    tee -a "$prom_dir/prometheus.yml" << EOF
  - job_name: 'consul-nodes'
    static_configs:
      - targets: ['${CONSUL_1_ADDR}:9100']
        labels:
          instance: 'consul-1'
      - targets: ['${CONSUL_2_ADDR}:9100']
        labels:
          instance: 'consul-2'
      - targets: ['${CONSUL_3_ADDR}:9100']
        labels:
          instance: 'consul-3'
  - job_name: 'consul-metrics'
    metrics_path: "/v1/agent/metrics"
    params:
      format: ['prometheus']
    static_configs:
      - targets: ['${CONSUL_1_ADDR}:8500']
        labels:
          instance: 'consul-1'
      - targets: ['${CONSUL_2_ADDR}:8500']
        labels:
          instance: 'consul-2'
      - targets: ['${CONSUL_3_ADDR}:8500']
        labels:
          instance: 'consul-3'
EOF
  fi
}

begin_time=$(date +%s)
end_time=$((begin_time + TIMEOUT_SECONDS))
while [[ "$(date +%s)" -lt "${end_time}" ]]; do
  if install_prometheus; then
    exit 0
  fi

  sleep "${RETRY_INTERVAL}"
done

fail "Timed out waiting for prometheus to install"
