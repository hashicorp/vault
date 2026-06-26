#!/usr/bin/env bash
# Copyright IBM Corp. 2016, 2026
# SPDX-License-Identifier: BUSL-1.1

set -euo pipefail

# Upload test results to Instana by creating spans for each test.
# Processes gotestsum JSON output via tools/instana-uploader.
#
# Usage: upload-test-results-to-instana.sh <json-file> <workflow-name> <matrix-id>
#
# Environment Variables:
#   INSTANA_AGENT_KEY:    Required. Instana agent key for authentication.
#   INSTANA_ENDPOINT_URL: Optional. Instana backend URL for serverless/direct mode.

if [ $# -ne 3 ]; then
    echo "Usage: $0 <json-file> <workflow-name> <matrix-id>" >&2
    exit 1
fi

JSON_FILE="$1"
WORKFLOW_NAME="$2"
MATRIX_ID="$3"

if [ -z "${INSTANA_AGENT_KEY:-}" ]; then
    echo "INSTANA_AGENT_KEY not set, skipping Instana upload" >&2
    exit 0
fi

# Resolve to absolute path before passing to the Go tool
JSON_FILE_ABS="$(cd "$(dirname "$JSON_FILE")" && pwd)/$(basename "$JSON_FILE")"

if [ ! -f "$JSON_FILE_ABS" ]; then
    echo "Error: test results file not found: $JSON_FILE_ABS" >&2
    exit 1
fi

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
VAULT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

go run "$VAULT_ROOT/tools/instana-uploader/" "$JSON_FILE_ABS" "$WORKFLOW_NAME" "$MATRIX_ID"
