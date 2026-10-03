#!/usr/bin/env bash
# Copyright IBM Corp. 2016, 2025
# SPDX-License-Identifier: BUSL-1.1

set -e

fail() {
  echo "$1" 1>&2
  return 1
}

[[ -z "$AGENT_LISTEN_ADDR" ]] && fail "AGENT_LISTEN_ADDR env variable has not been set"
[[ -z "$VAULT_ADDR" ]] && fail "VAULT_ADDR env variable has not been set"
[[ -z "$VAULT_INSTALL_DIR" ]] && fail "VAULT_INSTALL_DIR env variable has not been set"
[[ -z "$VAULT_AGENT_TEMPLATE_CONTENTS" ]] && fail "VAULT_AGENT_TEMPLATE_CONTENTS env variable has not been set"
[[ -z "$VAULT_AGENT_TEMPLATE_DESTINATION" ]] && fail "VAULT_AGENT_TEMPLATE_DESTINATION env variable has not been set"
[[ -z "$VAULT_TOKEN" ]] && fail "VAULT_TOKEN env variable has not been set"

binpath=${VAULT_INSTALL_DIR}/vault
test -x "$binpath" || fail "unable to locate vault binary at $binpath"

# If approle was already enabled, disable it as we're about to re-enable it (the || true is so we don't fail if it doesn't already exist)
$binpath auth disable approle || true

$binpath auth enable approle

$binpath write auth/approle/role/agent-role secret_id_ttl=700h token_num_uses=1000 token_ttl=600h token_max_ttl=700h secret_id_num_uses=1000

ROLEID=$($binpath read --format=json auth/approle/role/agent-role/role-id   | jq -r '.data.role_id')

if [[ "$ROLEID" == '' ]]; then
  fail "expected ROLEID to be nonempty, but it is empty"
fi

SECRETID=$($binpath write -f --format=json  auth/approle/role/agent-role/secret-id  | jq -r '.data.secret_id')

if [[ "$SECRETID" == '' ]]; then
  fail "expected SECRETID to be nonempty, but it is empty"
fi

echo "$ROLEID" > /tmp/role-id
echo "$SECRETID" > /tmp/secret-id

cat > /tmp/vault-agent.hcl <<- EOM
pid_file = "/tmp/pidfile"

vault {
  address = "${VAULT_ADDR}"
  tls_skip_verify = true
  retry {
    num_retries = 10
  }
}

cache {
  enforce_consistency = "always"
  use_auto_auth_token = true
}

listener "tcp" {
  address = "${AGENT_LISTEN_ADDR}"
  tls_disable = true
}

template {
  destination  = "${VAULT_AGENT_TEMPLATE_DESTINATION}"
  contents     = "${VAULT_AGENT_TEMPLATE_CONTENTS}"
  exec {
    command = "pkill -F /tmp/pidfile"
  }
}

auto_auth {
  method {
    type      = "approle"
    config = {
      role_id_file_path   = "/tmp/role-id"
      secret_id_file_path = "/tmp/secret-id"
    }
  }
  sink {
    type = "file"
    config = {
      path = "/tmp/token"
    }
  }
}
EOM

# If Agent is still running from a previous run, kill it
pkill -F /tmp/pidfile || true

# If the template file already exists, remove it
rm "${VAULT_AGENT_TEMPLATE_DESTINATION}" || true

# Run agent (it will kill itself when it finishes rendering the template)
if ! $binpath agent -config=/tmp/vault-agent.hcl > /tmp/agent-logs.txt 2>&1; then
  fail "failed to run vault agent: $(cat /tmp/agent-logs.txt)"
fi
