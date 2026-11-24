#!/usr/bin/env bash
# Copyright IBM Corp. 2026
# SPDX-License-Identifier: BUSL-1.1

set -euo pipefail

image=${1:?Usage: check-emudhra.sh IMAGE}
chains_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")/reject_chains" && pwd)"

docker run --rm --pull=never --user 0 \
  --mount "type=bind,source=$chains_dir,target=/reject_chains,readonly" \
  --entrypoint /bin/sh "$image" -ec '
  # Package installation must not change the trust store being tested.
  bundle=/etc/ssl/certs/ca-certificates.crt
  if [ -s /etc/pki/tls/certs/ca-bundle.crt ]; then
    bundle=/etc/pki/tls/certs/ca-bundle.crt
  fi
  cp "$bundle" /tmp/image-ca-bundle.pem
  mkdir /tmp/empty-ca-path
  if ! command -v openssl >/dev/null 2>&1; then
    if command -v apk >/dev/null 2>&1; then
      apk add --no-cache openssl
    else
      microdnf install -y openssl
    fi
  fi
  failed=0
  for chain in /reject_chains/*.pem; do
    test -s "$chain"
    # Expired fixtures must not hide continued trust in a compromised chain.
    if output=$(openssl verify -show_chain -no_check_time -purpose sslserver \
      -CAfile /tmp/image-ca-bundle.pem -CApath /tmp/empty-ca-path -no-CAstore \
      -untrusted "$chain" "$chain" 2>&1); then
      echo "ERROR: image trusts compromised chain $chain" >&2
      failed=1
    elif printf "%s\n" "$output" | grep -Eq "^error (2|18|19|20|21|27|28) at [0-9]+ depth lookup:"; then
      echo "PASS: rejected $chain"
    else
      echo "ERROR: unexpected verification failure for $chain" >&2
      failed=1
    fi
    printf "%s\n" "$output"
  done
  exit "$failed"
'
