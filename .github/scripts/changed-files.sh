#!/usr/bin/env bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

# Determine what files have changed between two git references.
#
# * For pull_request event_type's we'll the merge target (base_ref) with the pull requests reference,
#   (ref_name) which is usually a branch name.
# * For other event types (push, workflow_call) we don't have a base_ref target to merge into, so
#   instead we'll compare the last commit.
#
# Write the resulting metadata to STDOUT and $GITHUB_OUTPUT if it's defined.

event_type=$1 # GH event type (pull_request, push, workflow_call)
ref_name=$2 # branch reference that triggered the workflow
base_ref=$3 # PR branch base ref

if [[ "$event_type" == "pull_request" ]]; then
  git fetch --no-tags --prune origin "$base_ref"
  head_commit="HEAD"
  base_commit="origin/$base_ref"
else
  git fetch --no-tags --prune origin "$ref_name"
  head_commit=$(git log "origin/$ref_name" --oneline | head -1 | awk '{print $1}')
  base_commit=$(git log "origin/$ref_name" --oneline | head -2 | awk 'NR==2 {print $1}')
fi

docs_changed=false
ui_changed=false
app_changed=false

if ! files="$(git diff "${base_commit}...${head_commit}" --name-only)"; then
  echo "failed to get changed files from git"
  exit 1
fi

for file in $(awk -F "/" '{ print $1}' <<< "$files" | uniq); do
  if [[ "$file" == "changelog" ]]; then
    continue
  fi

  if [[ "$file" == "website" ]]; then
    docs_changed=true
    continue
  fi

  if [[ "$file" == "ui" ]]; then
    ui_changed=true
    continue
  fi

  # Anything that isn't either a changelog, ui, or docs change we'll consider an app change.
  app_changed=true
done

echo "app-changed=${app_changed}"
echo "docs-changed=${docs_changed}"
echo "ui-changed=${ui_changed}"
echo "files='${files}'"
[ -n "$GITHUB_OUTPUT" ] && {
  echo "app-changed=${app_changed}"
  echo "docs-changed=${docs_changed}"
  echo "ui-changed=${ui_changed}"
  # Use a random delimiter for multiline strings.
  # https://docs.github.com/en/actions/using-workflows/workflow-commands-for-github-actions#multiline-strings
  delimiter="$(openssl rand -hex 8)"
  echo "files<<${delimiter}"
  echo "${files}"
  echo "${delimiter}"
} >> "$GITHUB_OUTPUT"
