#!/usr/bin/env bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1


set -e

binpath=${vault_install_dir}/vault

fail() {
  echo "$1" 1>&2
  return 1
}

test -x "$binpath" || fail "unable to locate vault binary at $binpath"

export VAULT_ADDR='http://127.0.0.1:8200'
export VAULT_TOKEN='${vault_token}'

# If approle was already enabled, disable it as we're about to re-enable it (the || true is so we don't fail if it doesn't already exist)
$binpath auth disable approle || true

approle_create_status=$($binpath auth enable approle)

approle_status=$($binpath write auth/approle/role/agent-role secret_id_ttl=700h token_num_uses=1000 token_ttl=600h token_max_ttl=700h secret_id_num_uses=1000)

ROLEID=$($binpath read --format=json auth/approle/role/agent-role/role-id   | jq -r '.data.role_id')

if [[ "$ROLEID" == '' ]]; then
  fail "expected ROLEID to be nonempty, but it is empty"
fi

SECRETID=$($binpath write -f --format=json  auth/approle/role/agent-role/secret-id  | jq -r '.data.secret_id')

if [[ "$SECRETID" == '' ]]; then
  fail "expected SECRETID to be nonempty, but it is empty"
fi

echo $ROLEID > /tmp/role-id
echo $SECRETID > /tmp/secret-id

cat > /tmp/vault-agent.hcl <<- EOM
pid_file = "/tmp/pidfile"

vault {
  address = "http://127.0.0.1:8200"
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
    address = "127.0.0.1:8100"
    tls_disable = true
}

template {
  destination  = "${vault_agent_template_destination}"
  contents     = "${vault_agent_template_contents}"
  exec {
    command = "pkill -F /tmp/pidfile"
  }
}

auto_auth {
  method {
    type      = "approle"
    config = {
      role_id_file_path = "/tmp/role-id"
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
rm ${vault_agent_template_destination} || true

# Run agent (it will kill itself when it finishes rendering the template)
$binpath agent -config=/tmp/vault-agent.hcl > /tmp/agent-logs.txt 2>&1
