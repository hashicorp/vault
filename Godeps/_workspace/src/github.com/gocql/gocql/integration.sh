#!/bin/bash

set -e

function run_tests() {
	local clusterSize=3
	local version=$1
	local auth=$2

	if [ "$auth" = true ]; then
		clusterSize=1
	fi

	local keypath="$(pwd)/testdata/pki"

	local conf=(
	    "client_encryption_options.enabled: true"
	    "client_encryption_options.keystore: $keypath/.keystore"
	    "client_encryption_options.keystore_password: cassandra"
	    "client_encryption_options.require_client_auth: true"
	    "client_encryption_options.truststore: $keypath/.truststore"
	    "client_encryption_options.truststore_password: cassandra"
	    "concurrent_reads: 2"
	    "concurrent_writes: 2"
	    "rpc_server_type: sync"
	    "rpc_min_threads: 2"
	    "rpc_max_threads: 2"
	    "write_request_timeout_in_ms: 5000"
	    "read_request_timeout_in_ms: 5000"
	)

	ccm remove test || true

	ccm create test -v binary:$version -n $clusterSize -d --vnodes --jvm_arg="-Xmx256m -XX:NewSize=100m"
    ccm updateconf "${conf[@]}"

    if [ "$auth" = true ]
	then
		ccm updateconf 'authenticator: PasswordAuthenticator' 'authorizer: CassandraAuthorizer'
		rm -rf $HOME/.ccm/test/node1/data/system_auth
	fi

	local proto=2
	if [[ $version == 1.2.* ]]; then
		proto=1
	elif [[ $version == 2.0.* ]]; then
		proto=2
	elif [[ $version == 2.1.* ]]; then
		proto=3
	elif [[ $version == 2.2.* ]]; then
		proto=4
		ccm updateconf 'enable_user_defined_functions: true'
	fi

	ccm start -v
	ccm status
	ccm node1 nodetool status

	if [ "$auth" = true ]
	then
    	sleep 30s
    	go test -v . -timeout 15s -run=TestAuthentication -tags integration -runssl -runauth -proto=$proto -cluster=$(ccm liveset) -clusterSize=$clusterSize -autowait=1000ms
	else

		go test -timeout 5m -tags integration -v -gocql.timeout=10s -runssl -proto=$proto -rf=3 -cluster=$(ccm liveset) -clusterSize=$clusterSize -autowait=2000ms -compressor=snappy ./...

		if [ ${PIPESTATUS[0]} -ne 0 ]; then
			echo "--- FAIL: ccm status follows:"
			ccm status
			ccm node1 nodetool status
			ccm node1 showlog > status.log
			cat status.log
			echo "--- FAIL: Received a non-zero exit code from the go test execution, please investigate this"
			exit 1
		fi
	fi

	ccm remove
}

run_tests $1 $2
