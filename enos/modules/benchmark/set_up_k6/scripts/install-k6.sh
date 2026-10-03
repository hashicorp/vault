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

[[ -z "${RETRY_INTERVAL}" ]] && fail "RETRY_INTERVAL env variable has not been set"
[[ -z "${TIMEOUT_SECONDS}" ]] && fail "TIMEOUT_SECONDS env variable has not been set"

install_k6() {
  if command -v k6 &> /dev/null; then
    logger "k6 already installed"
    return 0
  fi

  sudo gpg -k
  sudo gpg --no-default-keyring --keyring /usr/share/keyrings/k6-archive-keyring.gpg --keyserver hkp://keyserver.ubuntu.com:80 --recv-keys C5AD17C747E3415A3642D57D77C6C491D6AC1D69
  echo "deb [signed-by=/usr/share/keyrings/k6-archive-keyring.gpg] https://dl.k6.io/deb stable main" | sudo tee /etc/apt/sources.list.d/k6.list
  sudo apt-get update
  sudo apt-get install -y k6

  logger "tweak some kernel parameters so k6 doesn't barf"
  sudo sysctl -w net.ipv4.ip_local_port_range="1024 65535"
  sudo sysctl -w net.ipv4.tcp_fin_timeout=10
  sudo sysctl -w net.ipv4.tcp_tw_reuse=1

  logger "setup k6 scripts dir"
  sudo mkdir -p /home/ubuntu/scripts
  sudo chown ubuntu:ubuntu /home/ubuntu/scripts
}

begin_time=$(date +%s)
end_time=$((begin_time + TIMEOUT_SECONDS))
while [[ "$(date +%s)" -lt "${end_time}" ]]; do
  if install_k6; then
    exit 0
  fi

  sleep "${RETRY_INTERVAL}"
done

fail "Timed out waiting for k6 to install"
