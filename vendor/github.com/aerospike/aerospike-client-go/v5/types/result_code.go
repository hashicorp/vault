// Copyright 2014-2021 Aerospike, Inc.
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
	// NETWORK_ERROR defines a network error. Checked the wrapped error for detail.
	NETWORK_ERROR ResultCode = -18

	// COMMON_ERROR defines a common, none-aerospike error. Checked the wrapped error for detail.
	COMMON_ERROR ResultCode = -17

	// MAX_RETRIES_EXCEEDED defines max retries limit reached.
	MAX_RETRIES_EXCEEDED ResultCode = -16

	// MAX_ERROR_RATE defines max errors limit reached.
	MAX_ERROR_RATE ResultCode = -15

	// RACK_NOT_DEFINED defines requested Rack for node/namespace was not defined in the cluster.
	RACK_NOT_DEFINED ResultCode = -13

	// INVALID_CLUSTER_PARTITION_MAP defines cluster has an invalid partition map, usually due to bad configuration.
	INVALID_CLUSTER_PARTITION_MAP ResultCode = -12

	// SERVER_NOT_AVAILABLE defines server is not accepting requests.
	SERVER_NOT_AVAILABLE ResultCode = -11

	// CLUSTER_NAME_MISMATCH_ERROR defines cluster Name does not match the ClientPolicy.ClusterName value.
	CLUSTER_NAME_MISMATCH_ERROR ResultCode = -10

	// RECORDSET_CLOSED defines recordset has already been closed or cancelled
	RECORDSET_CLOSED ResultCode = -9

	// NO_AVAILABLE_CONNECTIONS_TO_NODE defines there were no connections available to the node in the pool, and the pool was limited
	NO_AVAILABLE_CONNECTIONS_TO_NODE ResultCode = -8

	// TYPE_NOT_SUPPORTED defines data type is not supported by aerospike server.
	TYPE_NOT_SUPPORTED ResultCode = -7

	// COMMAND_REJECTED defines info Command was rejected by the server.
	COMMAND_REJECTED ResultCode = -6

	// QUERY_TERMINATED defines query was terminated by user.
	QUERY_TERMINATED ResultCode = -5

	// SCAN_TERMINATED defines scan was terminated by user.
	SCAN_TERMINATED ResultCode = -4

	// INVALID_NODE_ERROR defines chosen node is not currently active.
	INVALID_NODE_ERROR ResultCode = -3

	// PARSE_ERROR defines client parse error.
	PARSE_ERROR ResultCode = -2

	// SERIALIZE_ERROR defines client serialization error.
	SERIALIZE_ERROR ResultCode = -1

	// OK defines operation was successful.
	OK ResultCode = 0

	// SERVER_ERROR defines unknown server failure.
	SERVER_ERROR ResultCode = 1

	// KEY_NOT_FOUND_ERROR defines on retrieving, touching or replacing a record that doesn't exist.
	KEY_NOT_FOUND_ERROR ResultCode = 2

	// GENERATION_ERROR defines on modifying a record with unexpected generation.
	GENERATION_ERROR ResultCode = 3

	// PARAMETER_ERROR defines bad parameter(s) were passed in database operation call.
	PARAMETER_ERROR ResultCode = 4

	// KEY_EXISTS_ERROR defines on create-only (write unique) operations on a record that already
	// exists.
	KEY_EXISTS_ERROR ResultCode = 5

	// BIN_EXISTS_ERROR defines bin already exists on a create-only operation.
	BIN_EXISTS_ERROR ResultCode = 6

	// CLUSTER_KEY_MISMATCH defines expected cluster ID was not received.
	CLUSTER_KEY_MISMATCH ResultCode = 7

	// SERVER_MEM_ERROR defines server has run out of memory.
	SERVER_MEM_ERROR ResultCode = 8

	// TIMEOUT defines client or server has timed out.
	TIMEOUT ResultCode = 9

	// ALWAYS_FORBIDDEN defines operation not allowed in current configuration.
	ALWAYS_FORBIDDEN ResultCode = 10

	// PARTITION_UNAVAILABLE defines partition is unavailable.
	PARTITION_UNAVAILABLE ResultCode = 11

	// BIN_TYPE_ERROR defines operation is not supported with configured bin type (single-bin or
	// multi-bin).
	BIN_TYPE_ERROR ResultCode = 12

	// RECORD_TOO_BIG defines record size exceeds limit.
	RECORD_TOO_BIG ResultCode = 13

	// KEY_BUSY defines too many concurrent operations on the same record.
	KEY_BUSY ResultCode = 14

	// SCAN_ABORT defines scan aborted by server.
	SCAN_ABORT ResultCode = 15

	// UNSUPPORTED_FEATURE defines unsupported Server Feature (e.g. Scan + UDF)
	UNSUPPORTED_FEATURE ResultCode = 16

	// BIN_NOT_FOUND defines bin not found on update-only operation.
	BIN_NOT_FOUND ResultCode = 17

	// DEVICE_OVERLOAD defines device not keeping up with writes.
	DEVICE_OVERLOAD ResultCode = 18

	// KEY_MISMATCH defines key type mismatch.
	KEY_MISMATCH ResultCode = 19

	// INVALID_NAMESPACE defines invalid namespace.
	INVALID_NAMESPACE ResultCode = 20

	// BIN_NAME_TOO_LONG defines bin name length greater than 14 characters,
	// or maximum number of unique bin names are exceeded.
	BIN_NAME_TOO_LONG ResultCode = 21

	// FAIL_FORBIDDEN defines operation not allowed at this time.
	FAIL_FORBIDDEN ResultCode = 22

	// FAIL_ELEMENT_NOT_FOUND defines element Not Found in CDT
	FAIL_ELEMENT_NOT_FOUND ResultCode = 23

	// FAIL_ELEMENT_EXISTS defines element Already Exists in CDT
	FAIL_ELEMENT_EXISTS ResultCode = 24

	// ENTERPRISE_ONLY defines attempt to use an Enterprise feature on a Community server or a server
	// without the applicable feature key.
	ENTERPRISE_ONLY ResultCode = 25

	// OP_NOT_APPLICABLE defines the operation cannot be applied to the current bin value on the server.
	OP_NOT_APPLICABLE ResultCode = 26

	// FILTERED_OUT defines the transaction was not performed because the filter was false.
	FILTERED_OUT ResultCode = 27

	// LOST_CONFLICT defines write command loses conflict to XDR.
	LOST_CONFLICT = 28

	// QUERY_END defines there are no more records left for query.
	QUERY_END ResultCode = 50

	// SECURITY_NOT_SUPPORTED defines security type not supported by connected server.
	SECURITY_NOT_SUPPORTED ResultCode = 51

	// SECURITY_NOT_ENABLED defines administration command is invalid.
	SECURITY_NOT_ENABLED ResultCode = 52

	// SECURITY_SCHEME_NOT_SUPPORTED defines administration field is invalid.
	SECURITY_SCHEME_NOT_SUPPORTED ResultCode = 53

	// INVALID_COMMAND defines administration command is invalid.
	INVALID_COMMAND ResultCode = 54

	// INVALID_FIELD defines administration field is invalid.
	INVALID_FIELD ResultCode = 55

	// ILLEGAL_STATE defines security protocol not followed.
	ILLEGAL_STATE ResultCode = 56

	// INVALID_USER defines user name is invalid.
	INVALID_USER ResultCode = 60

	// USER_ALREADY_EXISTS defines user was previously created.
	USER_ALREADY_EXISTS ResultCode = 61

	// INVALID_PASSWORD defines password is invalid.
	INVALID_PASSWORD ResultCode = 62

	// EXPIRED_PASSWORD defines security credential is invalid.
	EXPIRED_PASSWORD ResultCode = 63

	// FORBIDDEN_PASSWORD defines forbidden password (e.g. recently used)
	FORBIDDEN_PASSWORD ResultCode = 64

	// INVALID_CREDENTIAL defines security credential is invalid.
	INVALID_CREDENTIAL ResultCode = 65

	// EXPIRED_SESSION defines login session expired.
	EXPIRED_SESSION ResultCode = 66

	// INVALID_ROLE defines role name is invalid.
	INVALID_ROLE ResultCode = 70

	// ROLE_ALREADY_EXISTS defines role already exists.
	ROLE_ALREADY_EXISTS ResultCode = 71

	// INVALID_PRIVILEGE defines privilege is invalid.
	INVALID_PRIVILEGE ResultCode = 72

	// INVALID_WHITELIST defines invalid IP address whiltelist
	INVALID_WHITELIST = 73

	// QUOTAS_NOT_ENABLED defines Quotas not enabled on server.
	QUOTAS_NOT_ENABLED = 74

	// INVALID_QUOTA defines invalid quota value.
	INVALID_QUOTA = 75

	// NOT_AUTHENTICATED defines user must be authentication before performing database operations.
	NOT_AUTHENTICATED ResultCode = 80

	// ROLE_VIOLATION defines user does not posses the required role to perform the database operation.
	ROLE_VIOLATION ResultCode = 81

	// NOT_WHITELISTED defines command not allowed because sender IP address not whitelisted.
	NOT_WHITELISTED = 82

	// QUOTA_EXCEEDED defines Quota exceeded.
	QUOTA_EXCEEDED = 83

	// UDF_BAD_RESPONSE defines a user defined function returned an error code.
	UDF_BAD_RESPONSE ResultCode = 100

	// BATCH_DISABLED defines batch functionality has been disabled.
	BATCH_DISABLED ResultCode = 150

	// BATCH_MAX_REQUESTS_EXCEEDED defines batch max requests have been exceeded.
	BATCH_MAX_REQUESTS_EXCEEDED ResultCode = 151

	// BATCH_QUEUES_FULL defines all batch queues are full.
	BATCH_QUEUES_FULL ResultCode = 152

	// GEO_INVALID_GEOJSON defines invalid GeoJSON on insert/update
	GEO_INVALID_GEOJSON ResultCode = 160

	// INDEX_FOUND defines secondary index already exists.
	INDEX_FOUND ResultCode = 200

	// INDEX_NOTFOUND defines requested secondary index does not exist.
	INDEX_NOTFOUND ResultCode = 201

	// INDEX_OOM defines secondary index memory space exceeded.
	INDEX_OOM ResultCode = 202

	// INDEX_NOTREADABLE defines secondary index not available.
	INDEX_NOTREADABLE ResultCode = 203

	// INDEX_GENERIC defines generic secondary index error.
	INDEX_GENERIC ResultCode = 204

	// INDEX_NAME_MAXLEN defines index name maximum length exceeded.
	INDEX_NAME_MAXLEN ResultCode = 205

	// INDEX_MAXCOUNT defines maximum number of indexes exceeded.
	INDEX_MAXCOUNT ResultCode = 206

	// QUERY_ABORTED defines secondary index query aborted.
	QUERY_ABORTED ResultCode = 210

	// QUERY_QUEUEFULL defines secondary index queue full.
	QUERY_QUEUEFULL ResultCode = 211

	// QUERY_TIMEOUT defines secondary index query timed out on server.
	QUERY_TIMEOUT ResultCode = 212

	// QUERY_GENERIC defines generic query error.
	QUERY_GENERIC ResultCode = 213

	// QUERY_NETIO_ERR defines query NetIO error on server
	QUERY_NETIO_ERR ResultCode = 214

	// QUERY_DUPLICATE defines duplicate TaskId sent for the statement
	QUERY_DUPLICATE ResultCode = 215

	// AEROSPIKE_ERR_UDF_NOT_FOUND defines uDF does not exist.
	AEROSPIKE_ERR_UDF_NOT_FOUND ResultCode = 1301

	// AEROSPIKE_ERR_LUA_FILE_NOT_FOUND defines lUA file does not exist.
	AEROSPIKE_ERR_LUA_FILE_NOT_FOUND ResultCode = 1302
)

// ResultCodeToString returns a human readable errors message based on the result code.
func ResultCodeToString(resultCode ResultCode) string {
	switch ResultCode(resultCode) {

	case NETWORK_ERROR:
		return "network error. Checked the wrapped error for detail"

	case COMMON_ERROR:
		return "common, none-aerospike error. Checked the wrapped error for detail"

	case MAX_RETRIES_EXCEEDED:
		return "Max retries exceeded"

	case MAX_ERROR_RATE:
		return "Max errors limit reached for node"

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
		return "Transaction filtered out"

	case LOST_CONFLICT:
		return "Write command loses conflict to XDR."

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

	case QUOTAS_NOT_ENABLED:
		return "Quotas not enabled"

	case INVALID_QUOTA:
		return "Invalid quota"

	case NOT_AUTHENTICATED:
		return "Not authenticated"

	case ROLE_VIOLATION:
		return "Role violation"

	case NOT_WHITELISTED:
		return "Command not whitelisted"

	case QUOTA_EXCEEDED:
		return "Quota exceeded"

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

func (rc ResultCode) String() string {
	switch rc {
	case NETWORK_ERROR:
		return "NETWORK_ERROR"
	case COMMON_ERROR:
		return "COMMON_ERROR"
	case MAX_RETRIES_EXCEEDED:
		return "MAX_RETRIES_EXCEEDED"
	case MAX_ERROR_RATE:
		return "MAX_ERROR_RATE"
	case RACK_NOT_DEFINED:
		return "RACK_NOT_DEFINED"
	case INVALID_CLUSTER_PARTITION_MAP:
		return "INVALID_CLUSTER_PARTITION_MAP"
	case SERVER_NOT_AVAILABLE:
		return "SERVER_NOT_AVAILABLE"
	case CLUSTER_NAME_MISMATCH_ERROR:
		return "CLUSTER_NAME_MISMATCH_ERROR"
	case RECORDSET_CLOSED:
		return "RECORDSET_CLOSED"
	case NO_AVAILABLE_CONNECTIONS_TO_NODE:
		return "NO_AVAILABLE_CONNECTIONS_TO_NODE"
	case TYPE_NOT_SUPPORTED:
		return "TYPE_NOT_SUPPORTED"
	case COMMAND_REJECTED:
		return "COMMAND_REJECTED"
	case QUERY_TERMINATED:
		return "QUERY_TERMINATED"
	case SCAN_TERMINATED:
		return "SCAN_TERMINATED"
	case INVALID_NODE_ERROR:
		return "INVALID_NODE_ERROR"
	case PARSE_ERROR:
		return "PARSE_ERROR"
	case SERIALIZE_ERROR:
		return "SERIALIZE_ERROR"
	case OK:
		return "OK"
	case SERVER_ERROR:
		return "SERVER_ERROR"
	case KEY_NOT_FOUND_ERROR:
		return "KEY_NOT_FOUND_ERROR"
	case GENERATION_ERROR:
		return "GENERATION_ERROR"
	case PARAMETER_ERROR:
		return "PARAMETER_ERROR"
	case KEY_EXISTS_ERROR:
		return "KEY_EXISTS_ERROR"
	case BIN_EXISTS_ERROR:
		return "BIN_EXISTS_ERROR"
	case CLUSTER_KEY_MISMATCH:
		return "CLUSTER_KEY_MISMATCH"
	case SERVER_MEM_ERROR:
		return "SERVER_MEM_ERROR"
	case TIMEOUT:
		return "TIMEOUT"
	case ALWAYS_FORBIDDEN:
		return "ALWAYS_FORBIDDEN"
	case PARTITION_UNAVAILABLE:
		return "PARTITION_UNAVAILABLE"
	case BIN_TYPE_ERROR:
		return "BIN_TYPE_ERROR"
	case RECORD_TOO_BIG:
		return "RECORD_TOO_BIG"
	case KEY_BUSY:
		return "KEY_BUSY"
	case SCAN_ABORT:
		return "SCAN_ABORT"
	case UNSUPPORTED_FEATURE:
		return "UNSUPPORTED_FEATURE"
	case BIN_NOT_FOUND:
		return "BIN_NOT_FOUND"
	case DEVICE_OVERLOAD:
		return "DEVICE_OVERLOAD"
	case KEY_MISMATCH:
		return "KEY_MISMATCH"
	case INVALID_NAMESPACE:
		return "INVALID_NAMESPACE"
	case BIN_NAME_TOO_LONG:
		return "BIN_NAME_TOO_LONG"
	case FAIL_FORBIDDEN:
		return "FAIL_FORBIDDEN"
	case FAIL_ELEMENT_NOT_FOUND:
		return "FAIL_ELEMENT_NOT_FOUND"
	case FAIL_ELEMENT_EXISTS:
		return "FAIL_ELEMENT_EXISTS"
	case ENTERPRISE_ONLY:
		return "ENTERPRISE_ONLY"
	case OP_NOT_APPLICABLE:
		return "OP_NOT_APPLICABLE"
	case FILTERED_OUT:
		return "FILTERED_OUT"
	case LOST_CONFLICT:
		return "LOST_CONFLICT"
	case QUERY_END:
		return "QUERY_END"
	case SECURITY_NOT_SUPPORTED:
		return "SECURITY_NOT_SUPPORTED"
	case SECURITY_NOT_ENABLED:
		return "SECURITY_NOT_ENABLED"
	case SECURITY_SCHEME_NOT_SUPPORTED:
		return "SECURITY_SCHEME_NOT_SUPPORTED"
	case INVALID_COMMAND:
		return "INVALID_COMMAND"
	case INVALID_FIELD:
		return "INVALID_FIELD"
	case ILLEGAL_STATE:
		return "ILLEGAL_STATE"
	case INVALID_USER:
		return "INVALID_USER"
	case USER_ALREADY_EXISTS:
		return "USER_ALREADY_EXISTS"
	case INVALID_PASSWORD:
		return "INVALID_PASSWORD"
	case EXPIRED_PASSWORD:
		return "EXPIRED_PASSWORD"
	case FORBIDDEN_PASSWORD:
		return "FORBIDDEN_PASSWORD"
	case INVALID_CREDENTIAL:
		return "INVALID_CREDENTIAL"
	case EXPIRED_SESSION:
		return "EXPIRED_SESSION"
	case INVALID_ROLE:
		return "INVALID_ROLE"
	case ROLE_ALREADY_EXISTS:
		return "ROLE_ALREADY_EXISTS"
	case INVALID_PRIVILEGE:
		return "INVALID_PRIVILEGE"
	case INVALID_WHITELIST:
		return "INVALID_WHITELIST"
	case QUOTAS_NOT_ENABLED:
		return "QUOTAS_NOT_ENABLED"
	case INVALID_QUOTA:
		return "INVALID_QUOTA"
	case NOT_AUTHENTICATED:
		return "NOT_AUTHENTICATED"
	case ROLE_VIOLATION:
		return "ROLE_VIOLATION"
	case NOT_WHITELISTED:
		return "NOT_WHITELISTED"
	case QUOTA_EXCEEDED:
		return "QUOTA_EXCEEDED"
	case UDF_BAD_RESPONSE:
		return "UDF_BAD_RESPONSE"
	case BATCH_DISABLED:
		return "BATCH_DISABLED"
	case BATCH_MAX_REQUESTS_EXCEEDED:
		return "BATCH_MAX_REQUESTS_EXCEEDED"
	case BATCH_QUEUES_FULL:
		return "BATCH_QUEUES_FULL"
	case GEO_INVALID_GEOJSON:
		return "GEO_INVALID_GEOJSON"
	case INDEX_FOUND:
		return "INDEX_FOUND"
	case INDEX_NOTFOUND:
		return "INDEX_NOTFOUND"
	case INDEX_OOM:
		return "INDEX_OOM"
	case INDEX_NOTREADABLE:
		return "INDEX_NOTREADABLE"
	case INDEX_GENERIC:
		return "INDEX_GENERIC"
	case INDEX_NAME_MAXLEN:
		return "INDEX_NAME_MAXLEN"
	case INDEX_MAXCOUNT:
		return "INDEX_MAXCOUNT"
	case QUERY_ABORTED:
		return "QUERY_ABORTED"
	case QUERY_QUEUEFULL:
		return "QUERY_QUEUEFULL"
	case QUERY_TIMEOUT:
		return "QUERY_TIMEOUT"
	case QUERY_GENERIC:
		return "QUERY_GENERIC"
	case QUERY_NETIO_ERR:
		return "QUERY_NETIO_ERR"
	case QUERY_DUPLICATE:
		return "QUERY_DUPLICATE"
	case AEROSPIKE_ERR_UDF_NOT_FOUND:
		return "AEROSPIKE_ERR_UDF_NOT_FOUND"
	case AEROSPIKE_ERR_LUA_FILE_NOT_FOUND:
		return "AEROSPIKE_ERR_LUA_FILE_NOT_FOUND"
	default:
		return "invalid ResultCode. Please report on https://github.com/aerospike/aerospike-client.go"
	}
}
