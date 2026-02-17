#!/bin/bash
set -e

echo "📦 Setting up Vault UI development environment..."

# Install Vault from HashiCorp's official APT repository
echo "🔐 Adding HashiCorp APT repository..."
wget -O- https://apt.releases.hashicorp.com/gpg | sudo gpg --dearmor -o /usr/share/keyrings/hashicorp-archive-keyring.gpg
echo "deb [signed-by=/usr/share/keyrings/hashicorp-archive-keyring.gpg] https://apt.releases.hashicorp.com $(lsb_release -cs) main" \
  | sudo tee /etc/apt/sources.list.d/hashicorp.list

echo "🔐 Installing Vault..."
sudo apt-get update && sudo SKIP_SETCAP=1 apt-get install -y vault

echo "📦 Installing pnpm..."
npm install -g pnpm@10.22.0

echo "📦 Installing UI dependencies (this may take a few minutes)..."
cd ui && pnpm install

echo "✅ Setup complete!"
