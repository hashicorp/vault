// Copyright 2013-2020 Aerospike, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package types

import "fmt"

// ResultCode signifies the database operation error codes.
// The positive numbers align with the server side file proto.h.
type ResultCode int

const (
	// Requested Rack for node/namespace was not defined in the cluster.
	RACK_NOT_DEFINED ResultCode = -13

	// Cluster has an invalid partition map, usually due to bad configuration.
	INVALID_CLUSTER_PARTITION_MAP ResultCode = -12

	// Server is not accepting requests.
	SERVER_NOT_AVAILABLE ResultCode = -11

	// Cluster Name does not match the ClientPolicy.ClusterName value.
	CLUSTER_NAME_MISMATCH_ERROR ResultCode = -10

	// Recordset has already been closed or cancelled
	RECORDSET_CLOSED ResultCode = -9

	// There were no connections available to the node in the pool, and the pool was limited
	NO_AVAILABLE_CONNECTIONS_TO_NODE ResultCode = -8

	// Data type is not supported by aerospike server.
	TYPE_NOT_SUPPORTED ResultCode = -7

	// Info Command was rejected by the server.
	COMMAND_REJECTED ResultCode = -6

	// Query was terminated by user.
	QUERY_TERMINATED ResultCode = -5

	// Scan was terminated by user.
	SCAN_TERMINATED ResultCode = -4

	// Chosen node is not currently active.
	INVALID_NODE_ERROR ResultCode = -3

	// Client parse error.
	PARSE_ERROR ResultCode = -2

	// Client serialization error.
	SERIALIZE_ERROR ResultCode = -1

	// Operation was successful.
	OK ResultCode = 0

	// Unknown server failure.
	SERVER_ERROR ResultCode = 1

	// On retrieving, touching or replacing a record that doesn't exist.
	KEY_NOT_FOUND_ERROR ResultCode = 2

	// On modifying a record with unexpected generation.
	GENERATION_ERROR ResultCode = 3

	// Bad parameter(s) were passed in database operation call.
	PARAMETER_ERROR ResultCode = 4

	// On create-only (write unique) operations on a record that already
	// exists.
	KEY_EXISTS_ERROR ResultCode = 5

	// Bin already exists on a create-only operation.
	BIN_EXISTS_ERROR ResultCode = 6

	// Expected cluster ID was not received.
	CLUSTER_KEY_MISMATCH ResultCode = 7

	// Server has run out of memory.
	SERVER_MEM_ERROR ResultCode = 8

	// Client or server has timed out.
	TIMEOUT ResultCode = 9

	// Operation not allowed in current configuration.
	ALWAYS_FORBIDDEN ResultCode = 10

	// Partition is unavailable.
	PARTITION_UNAVAILABLE ResultCode = 11

	// Operation is not supported with configured bin type (single-bin or
	// multi-bin).
	BIN_TYPE_ERROR ResultCode = 12

	// Record size exceeds limit.
	RECORD_TOO_BIG ResultCode = 13

	// Too many concurrent operations on the same record.
	KEY_BUSY ResultCode = 14

	// Scan aborted by server.
	SCAN_ABORT ResultCode = 15

	// Unsupported Server Feature (e.g. Scan + UDF)
	UNSUPPORTED_FEATURE ResultCode = 16

	// Bin not found on update-only operation.
	BIN_NOT_FOUND ResultCode = 17

	// Device not keeping up with writes.
	DEVICE_OVERLOAD ResultCode = 18

	// Key type mismatch.
	KEY_MISMATCH ResultCode = 19

	// Invalid namespace.
	INVALID_NAMESPACE ResultCode = 20

	// Bin name length greater than 14 characters,
	// or maximum number of unique bin names are exceeded.
	BIN_NAME_TOO_LONG ResultCode = 21

	// Operation not allowed at this time.
	FAIL_FORBIDDEN ResultCode = 22

	// Element Not Found in CDT
	FAIL_ELEMENT_NOT_FOUND ResultCode = 23

	// Element Already Exists in CDT
	FAIL_ELEMENT_EXISTS ResultCode = 24

	// Attempt to use an Enterprise feature on a Community server or a server
	// without the applicable feature key.
	ENTERPRISE_ONLY ResultCode = 25

	// The operation cannot be applied to the current bin value on the server.
	OP_NOT_APPLICABLE ResultCode = 26

	// The transaction was not performed because the predexp was false.
	FILTERED_OUT ResultCode = 27

	// There are no more records left for query.
	QUERY_END ResultCode = 50

	// Security type not supported by connected server.
	SECURITY_NOT_SUPPORTED ResultCode = 51

	// Administration command is invalid.
	SECURITY_NOT_ENABLED ResultCode = 52

	// Administration field is invalid.
	SECURITY_SCHEME_NOT_SUPPORTED ResultCode = 53

	// Administration command is invalid.
	INVALID_COMMAND ResultCode = 54

	// Administration field is invalid.
	INVALID_FIELD ResultCode = 55

	// Security protocol not followed.
	ILLEGAL_STATE ResultCode = 56

	// User name is invalid.
	INVALID_USER ResultCode = 60

	// User was previously created.
	USER_ALREADY_EXISTS ResultCode = 61

	// Password is invalid.
	INVALID_PASSWORD ResultCode = 62

	// Security credential is invalid.
	EXPIRED_PASSWORD ResultCode = 63

	// Forbidden password (e.g. recently used)
	FORBIDDEN_PASSWORD ResultCode = 64

	// Security credential is invalid.
	INVALID_CREDENTIAL ResultCode = 65

	// Login session expired.
	EXPIRED_SESSION ResultCode = 66

	// Role name is invalid.
	INVALID_ROLE ResultCode = 70

	// Role already exists.
	ROLE_ALREADY_EXISTS ResultCode = 71

	// Privilege is invalid.
	INVALID_PRIVILEGE ResultCode = 72

	// Invalid IP address whiltelist
	INVALID_WHITELIST = 73

	// User must be authentication before performing database operations.
	NOT_AUTHENTICATED ResultCode = 80

	// User does not posses the required role to perform the database operation.
	ROLE_VIOLATION ResultCode = 81

	// Command not allowed because sender IP address not whitelisted.
	NOT_WHITELISTED = 82

	// A user defined function returned an error code.
	UDF_BAD_RESPONSE ResultCode = 100

	// Batch functionality has been disabled.
	BATCH_DISABLED ResultCode = 150

	// Batch max requests have been exceeded.
	BATCH_MAX_REQUESTS_EXCEEDED ResultCode = 151

	// All batch queues are full.
	BATCH_QUEUES_FULL ResultCode = 152

	// Invalid GeoJSON on insert/update
	GEO_INVALID_GEOJSON ResultCode = 160

	// Secondary index already exists.
	INDEX_FOUND ResultCode = 200

	// Requested secondary index does not exist.
	INDEX_NOTFOUND ResultCode = 201

	// Secondary index memory space exceeded.
	INDEX_OOM ResultCode = 202

	// Secondary index not available.
	INDEX_NOTREADABLE ResultCode = 203

	// Generic secondary index error.
	INDEX_GENERIC ResultCode = 204

	// Index name maximum length exceeded.
	INDEX_NAME_MAXLEN ResultCode = 205

	// Maximum number of indexes exceeded.
	INDEX_MAXCOUNT ResultCode = 206

	// Secondary index query aborted.
	QUERY_ABORTED ResultCode = 210

	// Secondary index queue full.
	QUERY_QUEUEFULL ResultCode = 211

	// Secondary index query timed out on server.
	QUERY_TIMEOUT ResultCode = 212

	// Generic query error.
	QUERY_GENERIC ResultCode = 213

	// Query NetIO error on server
	QUERY_NETIO_ERR ResultCode = 214

	// Duplicate TaskId sent for the statement
	QUERY_DUPLICATE ResultCode = 215

	// UDF does not exist.
	AEROSPIKE_ERR_UDF_NOT_FOUND ResultCode = 1301

	// LUA file does not exist.
	AEROSPIKE_ERR_LUA_FILE_NOT_FOUND ResultCode = 1302
)

// Should connection be put back into pool.
func KeepConnection(err error) bool {
	// if error is not an AerospikeError, Throw the connection away conservatively
	ae, ok := err.(AerospikeError)
	if !ok {
		return false
	}

	switch ae.resultCode {
	case 0, // Zero Value
		QUERY_TERMINATED,
		SCAN_TERMINATED,
		PARSE_ERROR,
		SERIALIZE_ERROR,
		SERVER_NOT_AVAILABLE,
		SCAN_ABORT,
		QUERY_ABORTED,

		INVALID_NODE_ERROR,
		SERVER_MEM_ERROR,
		TIMEOUT,
		INDEX_OOM,
		QUERY_TIMEOUT:
		return false
	default:
		return true
	}
}

// Return result code as a string.
func ResultCodeToString(resultCode ResultCode) string {
	switch ResultCode(resultCode) {
	case RACK_NOT_DEFINED:
		return "Requested Rack for node/namespace was not defined in the cluster."

	case INVALID_CLUSTER_PARTITION_MAP:
		return "Cluster has an invalid partition map, usually due to bad configuration."

	case SERVER_NOT_AVAILABLE:
		return "Server is not accepting requests."

	case CLUSTER_NAME_MISMATCH_ERROR:
		return "Cluster Name does not match the ClientPolicy.ClusterName value"

	case RECORDSET_CLOSED:
		return "Recordset has already been closed or cancelled."

	case NO_AVAILABLE_CONNECTIONS_TO_NODE:
		return "No available connections to the node. Connection Pool was empty, and limited to certain number of connections."

	case TYPE_NOT_SUPPORTED:
		return "Type cannot be converted to Value Type."

	case COMMAND_REJECTED:
		return "command rejected"

	case QUERY_TERMINATED:
		return "Query terminated"

	case SCAN_TERMINATED:
		return "Scan terminated"

	case INVALID_NODE_ERROR:
		return "Invalid node"

	case PARSE_ERROR:
		return "Parse error"

	case SERIALIZE_ERROR:
		return "Serialize error"

	case OK:
		return "ok"

	case SERVER_ERROR:
		return "Server error"

	case KEY_NOT_FOUND_ERROR:
		return "Key not found"

	case GENERATION_ERROR:
		return "Generation error"

	case PARAMETER_ERROR:
		return "Parameter error"

	case KEY_EXISTS_ERROR:
		return "Key already exists"

	case BIN_EXISTS_ERROR:
		return "Bin already exists"

	case CLUSTER_KEY_MISMATCH:
		return "Cluster key mismatch"

	case SERVER_MEM_ERROR:
		return "Server memory error"

	case TIMEOUT:
		return "Timeout"

	case ALWAYS_FORBIDDEN:
		return "Operation not allowed in current configuration."

	case PARTITION_UNAVAILABLE:
		return "Partition not available"

	case BIN_TYPE_ERROR:
		return "Bin type error"

	case RECORD_TOO_BIG:
		return "Record too big"

	case KEY_BUSY:
		return "Hot key"

	case SCAN_ABORT:
		return "Scan aborted"

	case UNSUPPORTED_FEATURE:
		return "Unsupported Server Feature"

	case BIN_NOT_FOUND:
		return "Bin not found"

	case DEVICE_OVERLOAD:
		return "Device overload"

	case KEY_MISMATCH:
		return "Key mismatch"

	case INVALID_NAMESPACE:
		return "Namespace not found"

	case BIN_NAME_TOO_LONG:
		return "Bin name length greater than 15 characters, or maximum number of unique bin names are exceeded"

	case FAIL_FORBIDDEN:
		return "Operation not allowed at this time"

	case FAIL_ELEMENT_NOT_FOUND:
		return "Element not found."

	case FAIL_ELEMENT_EXISTS:
		return "Element exists"

	case ENTERPRISE_ONLY:
		return "Enterprise only"

	case OP_NOT_APPLICABLE:
		return "Operation not applicable"

	case FILTERED_OUT:
		return "Transaction filtered out by predexp"

	case QUERY_END:
		return "Query end"

	case SECURITY_NOT_SUPPORTED:
		return "Security not supported"

	case SECURITY_NOT_ENABLED:
		return "Security not enabled"

	case SECURITY_SCHEME_NOT_SUPPORTED:
		return "Security scheme not supported"

	case INVALID_COMMAND:
		return "Invalid command"

	case INVALID_FIELD:
		return "Invalid field"

	case ILLEGAL_STATE:
		return "Illegal state"

	case INVALID_USER:
		return "Invalid user"

	case USER_ALREADY_EXISTS:
		return "User already exists"

	case INVALID_PASSWORD:
		return "Invalid password"

	case EXPIRED_PASSWORD:
		return "Expired password"

	case FORBIDDEN_PASSWORD:
		return "Forbidden password"

	case INVALID_CREDENTIAL:
		return "Invalid credential"

	case EXPIRED_SESSION:
		return "Login session expired"

	case INVALID_ROLE:
		return "Invalid role"

	case ROLE_ALREADY_EXISTS:
		return "Role already exists"

	case INVALID_PRIVILEGE:
		return "Invalid privilege"

	case INVALID_WHITELIST:
		return "Invalid whitelist"

	case NOT_AUTHENTICATED:
		return "Not authenticated"

	case ROLE_VIOLATION:
		return "Role violation"

	case NOT_WHITELISTED:
		return "Command not whitelisted"

	case UDF_BAD_RESPONSE:
		return "UDF returned error"

	case BATCH_DISABLED:
		return "Batch functionality has been disabled"

	case BATCH_MAX_REQUESTS_EXCEEDED:
		return "Batch max requests have been exceeded"

	case BATCH_QUEUES_FULL:
		return "All batch queues are full"

	case GEO_INVALID_GEOJSON:
		return "Invalid GeoJSON on insert/update"

	case INDEX_FOUND:
		return "Index already exists"

	case INDEX_NOTFOUND:
		return "Index not found"

	case INDEX_OOM:
		return "Index out of memory"

	case INDEX_NOTREADABLE:
		return "Index not readable"

	case INDEX_GENERIC:
		return "Index error"

	case INDEX_NAME_MAXLEN:
		return "Index name max length exceeded"

	case INDEX_MAXCOUNT:
		return "Index count exceeds max"

	case QUERY_ABORTED:
		return "Query aborted"

	case QUERY_QUEUEFULL:
		return "Query queue full"

	case QUERY_TIMEOUT:
		return "Query timeout"

	case QUERY_GENERIC:
		return "Query error"

	case QUERY_NETIO_ERR:
		return "Query NetIO error on server"

	case QUERY_DUPLICATE:
		return "Duplicate TaskId sent for the statement"

	case AEROSPIKE_ERR_UDF_NOT_FOUND:
		return "UDF does not exist."

	case AEROSPIKE_ERR_LUA_FILE_NOT_FOUND:
		return "LUA package/file does not exist."

	default:
		return fmt.Sprintf("Error code (%v) not available yet - please file an issue on github.", resultCode)
	}
}
