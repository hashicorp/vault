#!/usr/bin/env bash

if [[ ! -e http/bindata_assetfs.go ]]
then
  printf "Compiled UI assets not found. They can be built with: make static-dist\n\n"
else
  if [[ `find http/bindata_assetfs.go -mmin +10080` ]]
  then
    printf "Compiled UI assets are more than one week old. They can be rebuilt with: make static-dist\n\n"
  fi
fi
