# Copyright IBM Corp. 2016, 2025
# SPDX-License-Identifier: BUSL-1.1

observations {
    ledger_path = "/var/ledger.log"
    type_prefix_denylist = ["deny1", "deny2"]
    type_prefix_allowlist = ["allow1", "allow2", "allow3"]
    file_mode = "0777"
}
