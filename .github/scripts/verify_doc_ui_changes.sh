#!/bin/bash
#  This script validates if the git diff contains on docs changes

event_type=$1 # GH event type (pull_request)
ref_name=$2 # branch reference that triggered the workflow
base_ref=$3 # PR branch base ref

contains() {
    target=$1; shift
    for i; do
        if [[ "$i" == "$target" ]]; then
            return 0
        fi
    done
    return 1
}

if [[ "$event_type" == "pull_request" ]]; then
  git fetch --no-tags --prune origin $base_ref
  head_commit="HEAD"
  base_commit="origin/$base_ref"
else
  git fetch --no-tags --prune origin $ref_name
  head_commit=$(git log origin/$ref_name --oneline | head -1 | awk '{print $1}')
  base_commit=$(git log origin/$ref_name --oneline | head -2 | awk 'NR==2 {print $1}')
fi

# git diff with ... shows the differences between base_commit and head_commit starting at the last common commit
changed_dir=$(git diff $base_commit...$head_commit --name-only | awk -F"/" '{ print $1}' | uniq)
change_count=$(echo "$changed_dir" | wc -l)

if [[ $change_count -gt 2 && $change_count -ne 0 ]]; then
  echo "is_docs_ui_change=false" >> "$GITHUB_OUTPUT"
  exit 0
elif [[ $change_count -eq 1 && "$changed_dir" == "website" ]]; then
  echo "is_docs_ui_change=true" >> "$GITHUB_OUTPUT"
  exit 0
elif [[ $change_count -eq 1 && "$changed_dir" == "ui" ]]; then
  echo "is_docs_ui_change=true" >> "$GITHUB_OUTPUT"
  exit 0
else
  if ! contains "website" ${changed_dir[@]} || ! contains "ui" ${changed_dir[@]}; then
    echo "is_docs_ui_change=false" >> "$GITHUB_OUTPUT"
    exit 0
  else
    echo "is_docs_ui_change=true" >> "$GITHUB_OUTPUT"
  fi
fi
