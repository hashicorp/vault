# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

# This script is sourced into the shell running in a Github workflow.

# Usage:
# To test for deprecations locally using the script, follow these steps:
# Run the command ./deprecation-checker.sh <optional branch name to compare with>
# Optionally specify a branch name to compare with. This will only show deprecations present in files 
#    that have changed when compared to the specified branch.
#    For example: ./deprecation-checker.sh main
# If no branch name is specified, the command will show all the deprecations in the code.

# GHA runs this against the PR's base ref branch. 

# Staticcheck uses static analysis to finds bugs and performance issues, offers simplifications, 
# and enforces style rules.
# Here, it is used to check if a deprecated function, variable, constant or field is used.

echo "Installing staticcheck"
go install honnef.co/go/tools/cmd/staticcheck@2023.1.2 #v0.4.2

# Check if the compare branch parameter is present in commandline
compareBranch="$1"

# Run staticcheck 
echo "Performing Deprecations Check: Running Staticcheck"
staticcheck ./... | grep deprecated > staticcheckOutput.txt

# If no compare branch name specified, output all deprecations. 
if [ -z $1 ]; then
    echo "Results:"
    if [ -s staticcheckOutput.txt ]
    then
     echo "Use of deprecated function, variable, constant or field found"
     cat staticcheckOutput.txt 
     # output file clean up 
     rm staticcheckOutput.txt  
     exit 1 
    else
     echo "No deprecated function, variable, constant or field found!!"
     # output file clean up 
     rm staticcheckOutput.txt  
     exit 0
    fi
fi

# Get changed files names from the PR
changedFiles=$(git --no-pager diff --name-only HEAD "$(git merge-base HEAD "origin/$compareBranch")")

# Include deprecations details of only changed files in the PR
echo "Results:"

# deprecationsCount checks if any deprecations were found to fail later 
deprecationsCount=0

for fileName in ${changedFiles[@]}; do
if grep -q $fileName staticcheckOutput.txt; then

    # output deprecations in the file 
    grep $fileName staticcheckOutput.txt

    # deprecation found, increment count
    deprecationsCount=$((deprecationsCount+1))
fi
done

# Cleanup deprecations file
rm staticcheckOutput.txt  

if [ "$deprecationsCount" -ne "0" ]
then
    echo "Deprecations check failed. This check examines the entire file included in the PR for any deprecated functions, variables, constants, or fields. Please review your changes to ensure that you have not included any deprecated elements."
    exit 1 
else
    echo "No deprecated functions, variables, constants, or fields were found in the PR!"
fi
