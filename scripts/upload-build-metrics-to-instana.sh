#!/usr/bin/env bash
# Copyright IBM Corp. 2016, 2026
# SPDX-License-Identifier: BUSL-1.1

set -euo pipefail

# Upload a build stage metric to Instana by creating a span for it.
# Thin wrapper around tools/build-metric-uploader.
#
# Usage: upload-build-metrics-to-instana.sh <stage-name> <duration-seconds> <job-name>
#
# Environment Variables:
#   INSTANA_AGENT_KEY:    Required. Instana agent key for authentication.
#   INSTANA_ENDPOINT_URL: Optional. Instana backend URL for serverless/direct mode.

if [ $# -ne 3 ]; then
    echo "Usage: $0 <stage-name> <duration-seconds> <job-name>" >&2
    exit 1
fi

STAGE_NAME="$1"
DURATION="$2"
JOB_NAME="$3"

if [ -z "$STAGE_NAME" ] || [ -z "$DURATION" ] || [ -z "$JOB_NAME" ]; then
    echo "stage-name, duration-seconds, and job-name must all be non-empty" >&2
    exit 1
fi

if ! [[ "$DURATION" =~ ^[0-9]+$ ]]; then
    echo "duration-seconds must be a non-negative integer, got: $DURATION" >&2
    exit 1
fi

if [ -z "${INSTANA_AGENT_KEY:-}" ]; then
    echo "INSTANA_AGENT_KEY not set, skipping Instana upload" >&2
    exit 0
fi

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
VAULT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

go run "$VAULT_ROOT/tools/build-metric-uploader/" "$STAGE_NAME" "$DURATION" "$JOB_NAME"
