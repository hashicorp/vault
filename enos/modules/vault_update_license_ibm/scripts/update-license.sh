#!/usr/bin/env bash
# Copyright IBM Corp. 2016, 2025
# SPDX-License-Identifier: BUSL-1.1

set -e

fail() {
  echo "$1" 1>&2
  exit 1
}

[[ -z "$VAULT_IBM_LICENSE" ]] && fail "VAULT_IBM_LICENSE env variable has not been set"
[[ -z "$VAULT_IBM_LICENSE_EDITION" ]] && fail "VAULT_IBM_LICENSE_EDITION env variable has not been set"

# Validate that $VAULT_IBM_LICENSE_EDITION is a valid edition option among "standard", "plus", and "premium"
case "${VAULT_IBM_LICENSE_EDITION}" in
  standard | plus | premium) ;;
  *)
    fail "IBM license edition '${VAULT_IBM_LICENSE_EDITION}' is not a valid edition. Valid editions are 'standard', 'plus', or 'premium'. Check your IBM license entitlements for the correct edition."
    ;;
esac

# Update the license file with the new license and add the license entitlement to the vault configuration
# Do this as a superuser because the files are owned by vault
echo "$VAULT_IBM_LICENSE" | sudo tee /etc/vault.d/vault.lic > /dev/null
echo "license_entitlement {
  edition = \"$VAULT_IBM_LICENSE_EDITION\"
}" | sudo tee -a /etc/vault.d/vault.hcl > /dev/null

echo "License updated successfully"
