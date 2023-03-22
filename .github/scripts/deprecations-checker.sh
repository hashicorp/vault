# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

# This script is sourced into the shell running in a Github workflow.

# Usage:
# To check deprecations locally using the script, follow these steps:
# Run the command ./deprecations-checker.sh <optional branch name to compare with> 
# on vault/vault-enterprise repository path or any package path
# Optionally specify a branch name to compare with. This will only show deprecations present in files 
#    that have changed when compared to the specified branch.
#    For example: ./.github/scripts/deprecations-checker.sh main (or)
#                 ./.github/scripts/deprecations-checker.sh (Running this from vault repo path)
# If no branch name is specified, the command will show all the deprecations in the code.

# GHA runs this against the PR's base ref branch. 

# Staticcheck uses static analysis to finds bugs and performance issues, offers simplifications, 
# and enforces style rules.
# Here, it is used to check if a deprecated function, variable, constant or field is used.
echo "Installing staticcheck"
go install honnef.co/go/tools/cmd/staticcheck@2023.1.2 #v0.4.2

# revgrep is a CLI tool used to filter static analysis tools to only lines changed based on a commit reference
echo "Installing revgrep"
go install github.com/golangci/revgrep/cmd/revgrep

# Run staticcheck 
echo "Performing Deprecations Check: Running Staticcheck"

# If no compare branch name is specified, output all deprecations
# Else only output the deprecations from the changes added
if [ -z $1 ]
    then
        staticcheck ./... | grep deprecated
    else
        # GHA will run this to find only changes wrt PR's base ref branch
        # revgrep CLI tool will return an exit status of 1 if any issues match, else it will return 0
        staticcheck ./... | grep deprecated 2>&1 | revgrep "$(git merge-base HEAD "origin/$1")"
fi
