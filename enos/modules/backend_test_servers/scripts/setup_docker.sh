#!/usr/bin/env bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

set -e

# Checking if docker is installed
if command -v docker &> /dev/null; then
  echo "Docker is already installed: $(docker --version)"
  exit 0
fi

# Detect OS
if [ -f /etc/os-release ]; then
  # shellcheck disable=SC1091
  . /etc/os-release
  OS_ID=$ID
  VERSION=$VERSION_ID
else
  echo "Cannot detect OS. /etc/os-release not found."
  exit 1
fi
echo "Detected OS: name=$NAME, id=$OS_ID, version=$VERSION"

# Install Docker based on distro
case "$OS_ID" in
  "ubuntu" | "debian")
    sudo apt update -y
    sudo apt install apt-transport-https ca-certificates curl software-properties-common -y
    curl -fsSL https://get.docker.com | sudo sh
    ;;
  "centos" | "rhel" | "fedora")
    sudo yum update -y
    sudo yum install -y docker
    ;;
  "amzn")
    sudo yum update -y
    sudo amazon-linux-extras enable docker
    sudo yum install -y docker
    ;;
  "sles" | "opensuse-leap" | "leap")
    sudo zypper refresh
    sudo zypper install -y docker
    ;;
  "alpine")
    sudo apk update
    sudo apk add docker
    ;;
  *)
    echo "Unsupported OS: $OS_ID"
    exit 1
    ;;
esac

# Enable and start Docker
echo "Enabling and starting Docker service..."
sudo systemctl enable docker || true
sudo systemctl start docker || true
echo "Docker installation complete."
docker --version
