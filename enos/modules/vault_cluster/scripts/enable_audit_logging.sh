#!/bin/env sh

set -eux

LOG_DIR=$(dirname "$LOG_FILE_PATH")
FILE_NAME=$(basename "$LOG_FILE_PATH")

# setup dir
if [ ! -d "$LOG_DIR" ]
then
  sudo mkdir -p "$LOG_DIR"
  sudo chmod 600 "$LOG_DIR"
  sudo chown "$SERVICE_USER":"$SERVICE_USER" "$LOG_DIR"
fi

# create log file
sudo touch /tmp/"$FILE_NAME"
sudo mv /tmp/"$FILE_NAME" "$LOG_FILE_PATH"
sudo chmod 600 "$LOG_FILE_PATH"
sudo chown "$SERVICE_USER":"$SERVICE_USER" "$LOG_FILE_PATH"

sudo su vault -c "VAULT_TOKEN=$VAULT_TOKEN VAULT_ADDR=$VAULT_ADDR $VAULT_BIN_PATH audit enable file file_path=$LOG_FILE_PATH"
