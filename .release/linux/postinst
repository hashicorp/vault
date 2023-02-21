#!/bin/bash

if [[ -f /opt/vault/tls/tls.crt ]] && [[ -f /opt/vault/tls/tls.key ]]; then
  echo "Vault TLS key and certificate already exist. Exiting."
  exit 0
fi

echo "Generating Vault TLS key and self-signed certificate..."

# Create TLS and Data directory
mkdir --parents /opt/vault/tls
mkdir --parents /opt/vault/data

# Generate TLS key and certificate
cd /opt/vault/tls
openssl req \
  -out tls.crt \
  -new \
  -keyout tls.key \
  -newkey rsa:4096 \
  -nodes \
  -sha256 \
  -x509 \
  -subj "/O=HashiCorp/CN=Vault" \
  -days 1095 # 3 years

# Update file permissions
chown --recursive vault:vault /etc/vault.d
chown --recursive vault:vault /opt/vault
chmod 600 /opt/vault/tls/tls.crt /opt/vault/tls/tls.key
chmod 700 /opt/vault/tls

echo "Vault TLS key and self-signed certificate have been generated in '/opt/vault/tls'."

# Set IPC_LOCK capabilities on vault
setcap cap_ipc_lock=+ep /usr/bin/vault

if [ -d /run/systemd/system ]; then
    systemctl --system daemon-reload >/dev/null || true
fi

if [[ $(vault version) == *+ent* ]]; then
echo "
The following shall apply unless your organization has a separately signed Enterprise License Agreement or Evaluation Agreement governing your use of the software:
Software in this repository is subject to the license terms located in the software, copies of which are also available at https://eula.hashicorp.com/ClickThruELA-Global.pdf or https://www.hashicorp.com/terms-of-evaluation as applicable. Please read the license terms prior to using the software. Your installation and use of the software constitutes your acceptance of these terms. If you do not accept the terms, do not use the software.
"
fi
