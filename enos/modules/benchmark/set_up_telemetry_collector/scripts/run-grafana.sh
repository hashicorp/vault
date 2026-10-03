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

sudo systemctl stop grafana-server || true
sudo systemctl start grafana-server

# let grafana start up
begin_time=$(date +%s)
end_time=$((begin_time + TIMEOUT_SECONDS))
started="0"
while [[ "$(date +%s)" -lt "${end_time}" ]]; do
  if systemctl is-active grafana-server; then
    logger "grafana is running"
    started="1"
    break
  fi

  sleep "${RETRY_INTERVAL}"
done

if [[ "${started}" == "0" ]]; then
  fail "Error: Grafana did not start after waiting for ${TIMEOUT_SECONDS} seconds"
fi

# Get the actual UID of the default-prom datasource
GRAFANA_USER="admin"
GRAFANA_PASSWORD="admin"  # Or whatever you've set
AUTH=$(echo -n "$GRAFANA_USER:$GRAFANA_PASSWORD" | base64)

# Use a function to safely handle the API call and potential failure
get_datasource_uid() {
  local response
  response=$(curl -s -H "Authorization: Basic $AUTH" \
    "http://localhost:3000/api/datasources/name/default-prom")

  # Check if the response contains the UID
  if echo "$response" | grep -q '"uid"'; then
    echo "$response" | grep -o '"uid":"[^"]*"' | sed 's/"uid":"//;s/"//'
    return 0
  else
    fail "Error: Failed to get datasource UID. Response: $response"
  fi
}

begin_time=$(date +%s)
end_time=$((begin_time + TIMEOUT_SECONDS))
while [[ "$(date +%s)" -lt "${end_time}" ]]; do
  if DATASOURCE_UID=$(get_datasource_uid); then
    logger "Datasource UID: $DATASOURCE_UID"
    break
  fi

  sleep "${RETRY_INTERVAL}"
done

if [[ -z "${DATASOURCE_UID:-}" ]]; then
  fail "Error: Failed to get datasource UID after waiting for ${TIMEOUT_SECONDS} seconds"
fi

# Directory containing your dashboard templates
DASHBOARD_DIR="/etc/grafana/dashboards"

# Install jq if needed
if ! command -v jq &> /dev/null; then
    logger "Installing jq..."
    sudo apt-get update
    sudo apt-get install -y jq
fi

# Create a function to process each dashboard file
process_dashboard() {
  local dashboard_file="$1"
  local temp_file

  logger "Processing $dashboard_file"

  # Create a temporary file
  temp_file=$(mktemp)

  # Use jq to update all datasource references
  # The || clause will only execute if jq fails, due to set -e
  jq --arg uid "$DATASOURCE_UID" '
    walk(
      if type == "object" and .datasource != null and .datasource.type == "prometheus" then
        .datasource.uid = $uid
      else
        .
      end
    )
  ' "$dashboard_file" > "$temp_file" || {
    rm -f "$temp_file"
    fail "Error: jq processing failed for $dashboard_file"
  }

  # Replace the original file with the updated content
  sudo mv "$temp_file" "$dashboard_file"
  logger "Updated $dashboard_file with datasource UID: $DATASOURCE_UID"
  return 0
}

# Process all JSON files, with error count tracking
error_count=0

for dashboard_file in "$DASHBOARD_DIR"/*.json; do
  if ! process_dashboard "$dashboard_file"; then
    error_count=$((error_count + 1))
  fi
done

if [ $error_count -gt 0 ]; then
  logger "Warning: $error_count files could not be processed"
else
  logger "All dashboards processed successfully"
fi

# Restart Grafana to apply changes
sudo systemctl restart grafana-server
