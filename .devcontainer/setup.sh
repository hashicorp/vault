#!/bin/bash
set -e

echo "📦 Setting up Vault UI development environment..."

# Install Vault by downloading the binary directly
echo "🔐 Downloading Vault binary..."
VAULT_VERSION="1.21.3"
wget -q "https://releases.hashicorp.com/vault/${VAULT_VERSION}/vault_${VAULT_VERSION}_linux_amd64.zip"

echo "🔐 Installing Vault..."
sudo apt-get update && sudo apt-get install -y unzip
unzip -o -q "vault_${VAULT_VERSION}_linux_amd64.zip"
sudo mv vault /usr/local/bin/vault
sudo chmod +x /usr/local/bin/vault
rm "vault_${VAULT_VERSION}_linux_amd64.zip"

echo "📦 Installing pnpm..."
npm install -g pnpm@10.22.0

echo "📦 Installing UI dependencies (this may take a few minutes)..."
cd ui && pnpm install

echo "✅ Setup complete!"
