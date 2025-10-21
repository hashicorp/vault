#!/usr/bin/env bash
# Copyright IBM Corp. 2016, 2025
# SPDX-License-Identifier: BUSL-1.1

echo "Create components/ directory"
mkdir docs/components/ 

echo "Generating markdown files for components in core addon..." 

# iterate over every .ts and .js file in core/addon/components (including nested files)
# skip .hbs files and shamir/ directory 
find "./lib/core/addon/components" -type f ! -name "*.hbs" -not -path "*/shamir*"  -print0 | while IFS= read -r -d '' file; do
  component=`eval "echo $file | cut -d/ -f6"`; 

 # skip replication components
  if [[ "$component" == replication* ]]; then
    echo "🔃 skipping $component" 
    continue
  fi

  yarn docfy-md $component core $file
done