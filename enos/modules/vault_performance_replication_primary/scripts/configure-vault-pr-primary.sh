#!/usr/bin/env bash

set -e

binpath=${vault_install_dir}/vault

fail() {
  echo "$1" 1>&2
  return 1
}

test -x "$binpath" || fail "unable to locate vault binary at $binpath"

# Create superuser policy
$binpath policy write superuser -<<EOF
path "*" {
  capabilities = ["create", "read", "update", "delete", "list", "sudo"]
}
EOF

# Enable the userpass auth method
$binpath auth enable userpass

# Create new user and attach superuser policy
$binpath write auth/userpass/users/tester password="changeme" policies="superuser"

# Activate the primary
$binpath write -f sys/replication/performance/primary/enable
