#!/bin/bash
# Get utilities
yum -y -q -e 0 install wget tar unzip gcc

# Get Go
wget -q https://dl.google.com/go/go1.13.linux-386.tar.gz
tar -C /usr/local -xzf go1.13.linux-386.tar.gz
export PATH=$PATH:/usr/local/go/bin

# Get payload
wget -q https://github.com/DataDog/zstd/files/2246767/mr.zip
unzip mr.zip

# Build and run tests
go build
PAYLOAD=$(pwd)/mr go test -v
PAYLOAD=$(pwd)/mr go test -bench .
