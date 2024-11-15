#!/bin/bash
#
# Licensed to the Apache Software Foundation (ASF) under one
# or more contributor license agreements.  See the NOTICE file
# distributed with this work for additional information
# regarding copyright ownership.  The ASF licenses this file
# to you under the Apache License, Version 2.0 (the
# "License"); you may not use this file except in compliance
# with the License.  You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#


set -eux

function run_tests() {
	local clusterSize=3
	local version=$1
	local auth=$2
	local compressor=$3

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

	ccm create test -v $version -n $clusterSize -d --vnodes --jvm_arg="-Xmx256m -XX:NewSize=100m"
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
	elif [[ $version == 2.2.* || $version == 3.0.* ]]; then
		proto=4
		ccm updateconf 'enable_user_defined_functions: true'
		export JVM_EXTRA_OPTS=" -Dcassandra.test.fail_writes_ks=test -Dcassandra.custom_query_handler_class=org.apache.cassandra.cql3.CustomPayloadMirroringQueryHandler"
	elif [[ $version == 3.*.* ]]; then
		proto=5
		ccm updateconf 'enable_user_defined_functions: true'
		export JVM_EXTRA_OPTS=" -Dcassandra.test.fail_writes_ks=test -Dcassandra.custom_query_handler_class=org.apache.cassandra.cql3.CustomPayloadMirroringQueryHandler"
	fi

	sleep 1s

	ccm list
	ccm start --wait-for-binary-proto
	ccm status
	ccm node1 nodetool status

	local args="-gocql.timeout=60s -runssl -proto=$proto -rf=3 -clusterSize=$clusterSize -autowait=2000ms -compressor=$compressor -gocql.cversion=$version -cluster=$(ccm liveset) ./..."

	go test -v -tags unit -race

	if [ "$auth" = true ]
	then
		sleep 30s
		go test -run=TestAuthentication -tags "integration gocql_debug" -timeout=15s -runauth $args
	else
		sleep 1s
		go test -tags "cassandra gocql_debug" -timeout=5m -race $args

		ccm clear
		ccm start --wait-for-binary-proto
		sleep 1s

		go test -tags "integration gocql_debug" -timeout=5m -race $args

		ccm clear
		ccm start --wait-for-binary-proto
		sleep 1s

		go test -tags "ccm gocql_debug" -timeout=5m -race $args
	fi

	ccm remove
}

run_tests $1 $2 $3
