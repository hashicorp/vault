#!/bin/bash
set -e

echo "🚀 Starting Vault UI development environment..."
echo ""

# Start Vault dev server in the background
echo "🔐 Starting Vault dev server on port 8200..."
vault server \
  -dev \
  -dev-root-token-id=root \
  -dev-listen-address=0.0.0.0:8200 \
  -dev-ha \
  -dev-transactional \
  -log-level=error &

# Wait for Vault to be ready
echo "⏳ Waiting for Vault dev server..."
until curl -sf http://127.0.0.1:8200/v1/sys/health > /dev/null 2>&1; do
  sleep 0.5
done
echo "✅ Vault is running"
echo "   Root token: root"
echo "   API: http://localhost:8200"
echo ""

# Start the Ember dev server
cd /workspaces/vault/ui
echo "🌐 Starting Ember dev server on port 4200..."
echo "   This may take a minute on first run..."
echo ""
pnpm ember server --proxy http://127.0.0.1:8200
