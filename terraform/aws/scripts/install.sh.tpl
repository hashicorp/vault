#!/usr/bin/env bash
set -e

# Install packages
sudo apt-get update -y
sudo apt-get install -y curl unzip

# Download Vault into some temporary directory
curl -L "${download_url}" > /tmp/vault.zip

# Unzip it
cd /tmp
sudo unzip vault.zip
sudo mv vault /usr/local/bin
sudo chmod 0755 /usr/local/bin/vault
sudo chown root:root /usr/local/bin/vault

# Setup the configuration
cat <<EOF >/tmp/vault-config
${config}
EOF
sudo mv /tmp/vault-config /usr/local/etc/vault-config.json

# Setup the init script
cat <<EOF >/tmp/vault.service
#https://devopscube.com/setup-hashicorp-vault-beginners-guide/
#https://www.hashicorp.com/resources/hashicorp-vault-administrative-guide

[Unit]
Description=vault service
Requires=network-online.target
After=network-online.target
ConditionFileNotEmpty=/etc/vault/config.json
 
[Service]
EnvironmentFile=-/etc/sysconfig/vault
Environment=GOMAXPROCS=2
Restart=on-failure
ExecStart=/usr/bin/vault server -config=/usr/local/etc/vault-config.json
StandardOutput=/logs/vault/output.log
StandardError=/logs/vault/error.log
LimitMEMLOCK=infinity
ExecReload=/bin/kill -HUP $MAINPID
KillSignal=SIGTERM
 
[Install]
WantedBy=multi-user.target
EOF
sudo mv /tmp/vault.service /etc/systemd/system/vault.service
sudo chmod 0644 /etc/systemd/system/vault.service

# Extra install steps (if any)
${extra-install}

# Start Vault
sudo systemctl enable vault.service
sudo systemctl start vault