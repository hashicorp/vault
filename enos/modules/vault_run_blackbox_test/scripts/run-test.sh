#!/usr/bin/env bash
# Copyright IBM Corp. 2016, 2025
# SPDX-License-Identifier: BUSL-1.1

set -euo pipefail

fail() {
    echo "$1" 1>&2
    exit 1
}

# Check required environment variables
[[ -z "${VAULT_TOKEN}" ]] && fail "VAULT_TOKEN env variable has not been set"
[[ -z "${VAULT_ADDR}" ]] && fail "VAULT_ADDR env variable has not been set"
[[ -z "${VAULT_TEST_PACKAGE}" ]] && fail "VAULT_TEST_PACKAGE env variable has not been set"
[[ -z "${VAULT_EDITION}" ]] && fail "VAULT_EDITION env variable has not been set"

# Check required dependencies
echo "Checking required dependencies..."

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "ERROR: Go is not installed or not found in PATH."
    echo ""
    echo "To resolve this issue:"
    echo "  • On a developer machine: Install Go from https://golang.org/dl/"
    echo "  • In CI: Ensure the setup-go action is configured properly"
    echo "  • If Go is installed elsewhere, add it to your PATH environment variable"
    echo ""
    fail "Go is required to run blackbox tests."
fi

echo "Go version: $(go version)"

# Check if gotestsum is installed (required)
if ! command -v gotestsum &> /dev/null; then
    echo "ERROR: gotestsum is not installed or not found in PATH."
    echo ""
    echo "To resolve this issue:"
    echo "  • Run 'make tools' to install required development tools"
    echo "  • Ensure GOPATH/bin is in your PATH environment variable"
    echo "  • Or manually install: go install gotest.tools/gotestsum@v1.13.0"
    echo ""
    fail "gotestsum is required to run blackbox tests."
fi

# Check if jq is available (needed for parsing test matrix)
if ! command -v jq &> /dev/null; then
    fail "jq is not installed or not in PATH. jq is required to parse test matrix files."
fi

# Check if git is available (needed for git rev-parse)
if ! command -v git &> /dev/null; then
    fail "Git is not installed or not in PATH. Git is required to determine the repository root."
fi

# Verify we're in a git repository and get the root directory
if ! root_dir="$(git rev-parse --show-toplevel 2> /dev/null)"; then
    fail "Not in a git repository. Tests must be run from within the Vault repository."
fi

echo "All required dependencies are available."
pushd "$root_dir" > /dev/null

# Create unique output files for test results
timestamp="$(date +%s)_$$"
json_output="/tmp/vault_test_results_${timestamp}.json"
junit_output="/tmp/vault_test_results_${timestamp}.xml"

echo "Test results will be written to: $json_output"

# Run tests using gotestsum with JSON output and JUnit reporting
echo "Using gotestsum for enhanced test output and JUnit reporting"
echo "JUnit results will be written to: $junit_output"

echo "Running tests..."
echo "Vault environment variables:"
env | grep VAULT | sed 's/VAULT_TOKEN=.*/VAULT_TOKEN=***REDACTED***/'

case $VAULT_EDITION in
  ent | ent.hsm | ent.hsm.fips1402 | ent.hsm.fips1403 | ent.fips1403 | ent.fips1402)
    tags="-tags=ent,enterprise"
    ;;
  ce)
    tags=""
    ;;
  *)
    fail "unknown VAULT_EDITION: $VAULT_EDITION"
    ;;
esac

# Build gotestsum command based on whether we have specific tests
# Convert VAULT_TEST_PACKAGE to array to handle multiple package paths properly
IFS=$' ' read -r -d '' -a packages <<< "$VAULT_TEST_PACKAGE" || true

set -x # Show commands being executed
set +e # Temporarily disable exit on error
if [ -n "$VAULT_TEST_MATRIX" ] && [ -f "$VAULT_TEST_MATRIX" ]; then
    echo "Using test matrix from: $VAULT_TEST_MATRIX"
    # Extract test names from matrix and create regex pattern
    test_pattern=$(jq -r '.include[].test' "$VAULT_TEST_MATRIX" | paste -sd '|' -)
    echo "Running specific tests: $test_pattern"
    gotestsum --junitfile="$junit_output" --format=standard-verbose --jsonfile="$json_output" -- -count=1 "${tags}" -run="$test_pattern" "${packages[@]}"
else
    echo "Running all tests in package"
    gotestsum --junitfile="$junit_output" --format=standard-verbose --jsonfile="$json_output" -- -count=1 "${tags}" "${packages[@]}"
fi
test_exit_code=$?
set -e # Re-enable exit on error
set +x # Turn off command tracing

echo "Test execution completed with exit code: $test_exit_code"

# Check if JSON file was created successfully
if [ -f "$json_output" ] && [ -s "$json_output" ]; then
    echo "JSON file created successfully: $(wc -l < "$json_output") lines"
    echo "JSON_RESULTS_FILE=$json_output"

    # Check if JUnit file was created (only when using gotestsum)
    if [ -f "$junit_output" ] && [ -s "$junit_output" ]; then
        echo "JUnit file created successfully: $(wc -l < "$junit_output") lines"
        echo "JUNIT_RESULTS_FILE=$junit_output"
  else
        echo "JUNIT_RESULTS_FILE="
  fi
else
    echo "WARNING: Test results file not created or empty" >&2
    echo "TEST_STATUS=ERROR"
    echo "TEST_EXIT_CODE=$test_exit_code"
    echo "JSON_RESULTS_FILE="
    echo "JUNIT_RESULTS_FILE="
    # Don't exit here - continue to show what we can
fi

# Also output human-readable results to stdout
echo "=== TEST EXECUTION SUMMARY ==="
if [ $test_exit_code -eq 0 ]; then
    echo "✅ Tests PASSED"
else
    echo "❌ Tests FAILED (exit code: $test_exit_code)"
fi

# Parse JSON results and create a summary
echo "=== DETAILED RESULTS ==="
if [ -f "$json_output" ] && [ -s "$json_output" ]; then
    if command -v jq &> /dev/null; then
        # Use jq if available for better parsing
        echo "Test Results Summary (JSON):"
        set +e # Temporarily disable exit on error for jq parsing
        if jq -e . "$json_output" > /dev/null 2>&1; then
            # JSON is valid, proceed with parsing
            jq -r 'select(.Action == "pass" or .Action == "fail") | "\(.Time) \(.Action | ascii_upcase) \(.Test // "PACKAGE")"' "$json_output" 2> /dev/null || echo "Failed to parse test results"
    else
            echo "Invalid JSON in test results file, showing raw output:"
            cat "$json_output" 2> /dev/null || echo "Failed to read JSON file"
    fi
        set -e # Re-enable exit on error
  else
        # Fallback: show raw JSON output without jq
        echo "Test Results (JSON output - install jq for better formatting):"
        set +e # Temporarily disable exit on error
        if grep -q '"Action"' "$json_output" 2> /dev/null; then
            cat "$json_output" 2> /dev/null || echo "Failed to read JSON file"
    else
            echo "No structured test results found, showing raw output:"
            cat "$json_output" 2> /dev/null || echo "Failed to read JSON file"
    fi
        set -e # Re-enable exit on error
  fi
else
    echo "No JSON file to parse"
fi

# Output the JSON file path so Terraform can capture it (if not already output above)
if [ -f "$json_output" ] && [ -s "$json_output" ]; then
    echo "JSON_RESULTS_FILE=$json_output"
fi

popd > /dev/null

# Always output exit code for Terraform to capture, but exit 0 so script doesn't fail
echo "Final test exit code: $test_exit_code"

# Exit with the actual test exit code so Terraform fails on test failures
exit $test_exit_code
