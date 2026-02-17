#!/bin/bash
set -e

echo "📦 Setting up Vault UI development environment..."

# Install Vault by downloading the binary directly
echo "🔐 Downloading Vault binary..."
VAULT_VERSION="1.21.3"
VAULT_ZIP="vault_${VAULT_VERSION}_linux_amd64.zip"

# Work in /tmp to avoid polluting the workspace
cd /tmp
wget -q "https://releases.hashicorp.com/vault/${VAULT_VERSION}/${VAULT_ZIP}"

echo "🔐 Installing Vault..."
sudo apt-get update && sudo apt-get install -y unzip

# Clean up any existing vault installation
sudo rm -rf /usr/local/bin/vault

# Extract and install
unzip -o -q "${VAULT_ZIP}"
sudo mv vault /usr/local/bin/vault
sudo chmod +x /usr/local/bin/vault
rm -f "${VAULT_ZIP}" vault

# Return to workspace
cd /workspaces/vault

echo "✅ Vault $(vault version) installed"

echo "📦 Installing pnpm..."
npm install -g pnpm@10.22.0

echo "📦 Installing UI dependencies (this may take a few minutes)..."
cd ui && pnpm install

echo "✅ Setup complete!"
