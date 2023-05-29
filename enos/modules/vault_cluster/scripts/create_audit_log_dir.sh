#!/bin/env sh

set -eux

LOG_DIR=$(dirname "$LOG_FILE_PATH")

sudo mkdir -p "$LOG_DIR"
sudo chown "$SERVICE_USER":"$SERVICE_USER" "$LOG_DIR"
