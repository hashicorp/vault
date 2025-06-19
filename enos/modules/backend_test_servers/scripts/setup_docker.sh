#!/usr/bin/env bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

set -e

[[ -z "$DISTRO" ]] && fail "DISTRO env variable has not been set"

check_docker_running() {
  if docker info > /dev/null 2>&1; then
    return 0
  fi
  echo "Docker daemon not running."
  if [[ "$DISTRO" == "leap" || "$DISTRO" == "sles" ]]; then
    echo "Detected distro: $DISTRO. Attempting to start and enable Docker..."
    sudo systemctl start docker || true
    sudo systemctl enable docker || true
  fi
  echo "Waiting for Docker to start..."
  docker info
  echo "Docker is now running."
}

# Check if Docker is already installed
if command -v sudo docker &>/dev/null; then
  echo "Docker is already installed: $(sudo docker --version)"
  check_docker_running
  exit 0
fi

echo "Installing Docker for distro: $DISTRO"
case "$DISTRO" in
  ubuntu)
    sudo apt update -y
    sudo apt install -y apt-transport-https ca-certificates curl software-properties-common
    curl -fsSL https://get.docker.com | sudo sh
    ;;
  rhel)
    sudo yum update -y
    sudo yum install -y docker
    ;;
  amzn)
    sudo yum update -y
    sudo amazon-linux-extras enable docker
    sudo yum install -y docker
    ;;
  sles | leap)
    sudo zypper refresh
    sudo zypper install -y docker
    ;;
  *)
    echo "Unsupported OS: $DISTRO"
    exit 1
    ;;
esac

echo "Enabling and starting Docker service..."
sudo systemctl start docker || true
sudo systemctl enable docker || true

echo "Docker installation complete."
sudo docker --version
