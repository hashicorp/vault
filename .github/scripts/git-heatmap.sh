#!/usr/bin/env bash
# Copyright IBM Corp. 2016, 2025
# SPDX-License-Identifier: BUSL-1.1

# Build a commit frequency histogram of oft-edited files, scoped to the files
# changed since a base branch. Used in CI to flag "hot" files touched by a PR.
#
# Adapted from https://github.com/jez/git-heatmap/blob/master/git-heatmap
# Requires either a 'bars' (https://github.com/jez/bars) or 'barchart'
# (https://github.com/jez/barchart) command on PATH.

set -euo pipefail

ccyan="$(echo -ne '\033[0;36m')"
cnone="$(echo -ne '\033[0m')"

ARGV=()
while [[ $# -gt 0 ]]; do
    key="$1"
    case $key in
    -n)
        LIMIT="$2"
        shift
        shift
        ;;
    -b | --base)
        REVIEW_BASE="$2"
        shift
        shift
        ;;
    -c | --char)
        CHAR="$2"
        shift
        shift
        ;;
    --width)
        WIDTH="$2"
        shift
        shift
        ;;
    -f | --filter)
        FILTER="$2"
        shift
        shift
        ;;
    --exclude)
        EXCLUDE="$2"
        shift
        shift
        ;;
    -h)
        cat <<EOF
Heatmap of oft-edited files.
Usage:
  git-heatmap.sh [options] [<path>...]
Options:
  -n <top>                      Limit to top <n> files. [default: 30]
  --width <n>                   Limit histogram to <n> chars.
  -b <branch>, --base <branch>  Compare relative to <branch>. If on <branch>,
                                show heatmap for entire repo. [default: master]
  -c <char>, --char <char>      Use <char> to draw the bars. [default: █]
  -f <cmd>, --filter <cmd>      Filter output through <cmd> before creating the
                                the histogram.
  --exclude <regex>             Drop files matching this extended regex (e.g.
                                lockfiles, generated code) before tallying.
  -h                            Show this message.
EOF
        exit
        ;;
    *)
        ARGV+=("$1")
        shift
        ;;
    esac
done

LIMIT=${LIMIT:-30}
REVIEW_BASE=${REVIEW_BASE:-master}
CHAR=${CHAR:-█}
WIDTH=${WIDTH:-60}
FILTER=${FILTER:-cat -}
EXCLUDE=${EXCLUDE:-}

files() {
    # Tallies edit frequency across the entire repo history (not scoped to
    # any range), so counts reflect how "hot" a file has always been. The
    # caller narrows the resulting list down to this PR's own files below.
    # https://stackoverflow.com/questions/7577052/
    #
    # A plain `cut -f 2-` breaks on renames: --name-status emits those as
    # "R100\t<old>\t<new>" (3 fields), so cutting from field 2 onward glues
    # the old and new paths into one bogus string that matches neither real
    # path downstream. Printing the last field handles both the 2-field
    # (status, path) and 3-field (status, old, new) cases correctly.
    git log --name-status --pretty=format: -- "${ARGV[@]+"${ARGV[@]}"}" |
        awk -F'\t' 'NF>=2 {print $NF}'
}

color_name() {
    if [ -t 1 ]; then
        sed -e "s/\(..*\/\)*\(.[^|]*\) |/\1$ccyan\2$cnone |/"
    else
        cat -
    fi
}

filter() {
    grep '.' |
        (if [[ -n "$EXCLUDE" ]]; then grep -vE "$EXCLUDE"; else cat -; fi) |
        eval "$FILTER" |
        sort |
        uniq -c |
        sort -nr |
        head -n "$LIMIT"
}

histogram() {
    "$bars" --bar "$CHAR" --width "$WIDTH"
}

if command -v barchart &>/dev/null; then
    bars="barchart"
elif command -v bars &>/dev/null; then
    bars="bars"
else
    echo >&2 "$0: This command requires a command called 'bars' or 'barchart'. See the README."
    exit 1
fi

if [[ "$(git rev-parse --abbrev-ref HEAD)" == "$REVIEW_BASE" ]]; then
    # If on master, show heatmap for whole repo
    files | filter | histogram | color_name
else
    MERGE_BASE="$(git merge-base HEAD "$REVIEW_BASE")"
    files |
        # If on separate branch, show heatmap for files changed since master.
        # The ..HEAD form (not a bare "$MERGE_BASE") matters: a bare ref diffs
        # against the working tree, not just committed history.
        grep -xF -f <(git diff --name-only "$MERGE_BASE"..HEAD) |
        filter |
        histogram |
        color_name
fi
