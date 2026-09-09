#!/usr/bin/env bash
# Copyright IBM Corp. 2016, 2026
# SPDX-License-Identifier: BUSL-1.1

# Builds a manifest of relevant Go test targets based on the set of changed files
# in a pull request. For each changed Go file, it locates the corresponding
# *_test.go file, extracts the top-level Test functions and their line numbers,
# and emits a JSON array describing the directory, test file, and tests.
#
# Inputs (environment variables):
#   CHANGED_FILES_JSON - JSON payload describing changed files. Supported shapes:
#                        {"files": [...]}, {"changed_files": [...]},
#                        {"changedFiles": [...]}; entries may be strings or
#                        objects with a "path", "filename", or "file" key.
#   SUITE_NAME         - Identifier used to name the output manifest file.
#
# Outputs:
#   Writes the manifest to "relevant-tests-${SUITE_NAME}.json".
#   If GITHUB_OUTPUT is set, appends "manifest=<file>" to it.

set -euo pipefail

suite_name="${SUITE_NAME:-}"
out_file="relevant-tests-${suite_name}.json"
changed_files_json="${CHANGED_FILES_JSON:-}"
if [ -z "$changed_files_json" ]; then
  changed_files_json='{}'
fi

emit_manifest_output() {
  if [ -n "${GITHUB_OUTPUT:-}" ]; then
    echo "manifest=$out_file" >> "$GITHUB_OUTPUT"
  fi
}

# Extract changed paths from known shapes of changed-files payload.
changed_paths=$(jq -r '
  (.files // .changed_files // .changedFiles // [])
  | map(if type == "string" then . else (.path // .filename // .file // empty) end)
  | .[]?
' <<< "${changed_files_json}" | awk 'NF' | sort -u)

if [ -z "${changed_paths}" ]; then
  echo '[]' > "$out_file"
  emit_manifest_output
  exit 0
fi

map_file="$(mktemp)"
while IFS= read -r path; do
  case "$path" in
    *.go)
      dir="$(dirname "$path")"
      base="$(basename "$path")"
      if [[ "$base" == *_test.go ]]; then
        test_file="$path"
      else
        test_file="${path%.go}_test.go"
      fi
      printf '%s\t%s\n' "$dir" "$test_file" >> "$map_file"
      ;;
  esac
done <<< "$changed_paths"

if [ ! -s "$map_file" ]; then
  echo '[]' > "$out_file"
  emit_manifest_output
  exit 0
fi

jq -n '[]' > "$out_file"
while IFS=$'\t' read -r dir test_file; do
  if [ ! -f "$test_file" ]; then
    continue
  fi

  # Capture top-level Go test functions and their line numbers in the target test file.
  tests=$(awk '
    /^func[[:space:]]+Test[[:alnum:]_]+\(t[[:space:]]+\*testing\.T\)/ {
      name=$2;
      sub(/\(.*/, "", name);
      print NR ":" name;
    }
  ' "$test_file" | jq -Rsc '
    split("\n")
    | map(select(length > 0))
    | map(split(":") | {"line": (.[0] | tonumber), "name": .[1]})
  ')

  if [ "$tests" = "[]" ]; then
    continue
  fi

  tmp_file="$(mktemp)"
  jq \
    --arg dir "$dir" \
    --arg test_file "$test_file" \
    --argjson tests "$tests" \
    '. + [{"dir": $dir, "test_file": $test_file, "tests": $tests}]' \
    "$out_file" > "$tmp_file"
  mv "$tmp_file" "$out_file"
done < <(sort -u "$map_file")

rm -f "$map_file"

emit_manifest_output
