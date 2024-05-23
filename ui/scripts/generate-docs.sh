#!/usr/bin/env bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

echo "Generating markdown files in core addon..." 

for FILE in ./lib/core/addon/components/*.js;  do
  component=`eval "echo $FILE | cut -d/ -f6"`; 
  if [[ "$component" == replication* ]]; then
    echo "🔃 skipping $component" 
    continue
  fi

  yarn docfy-md $component core

  echo "✅ $component" 
done