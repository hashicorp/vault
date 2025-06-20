#!/usr/bin/env bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

set -e

[[ -z "$DISTRO" ]] && fail "DISTRO env variable has not been set"

# Checking if docker is installed
if command -v docker &> /dev/null; then
  echo "Docker is already installed: $(docker --version)"
  if [[ "$DISTRO" == "leap" || "$DISTRO" == "sles" ]]; then
    echo "Detected distro: $DISTRO. Attempting to start and enable Docker..."
    sudo systemctl start docker || true
    sudo systemctl enable docker || true
  fi
  exit 0
fi

# Install Docker based on distro
echo "Installing Docker for distro: $DISTRO"
case "$DISTRO" in
  "ubuntu")
    sudo apt update -y
    sudo apt install apt-transport-https ca-certificates curl software-properties-common -y
    curl -fsSL https://get.docker.com | sudo sh
    ;;
  "rhel")
    sudo yum update -y
    sudo yum install -y docker
    ;;
  "amzn")
    sudo yum update -y
    sudo amazon-linux-extras enable docker
    sudo yum install -y docker
    ;;
  "sles" | "leap")
    sudo zypper refresh
    sudo zypper install -y docker
    ;;
  *)
    echo "Unsupported OS: $DISTRO"
    exit 1
    ;;
esac

echo "Docker installation complete."
sudo docker info
