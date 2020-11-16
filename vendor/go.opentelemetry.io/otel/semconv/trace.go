// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package semconv

import "go.opentelemetry.io/otel/label"

// Semantic conventions for attribute keys used for network related
// operations.
const (
	// Transport protocol used.
	NetTransportKey = label.Key("net.transport")

	// Remote address of the peer.
	NetPeerIPKey = label.Key("net.peer.ip")

	// Remote port number.
	NetPeerPortKey = label.Key("net.peer.port")

	// Remote hostname or similar.
	NetPeerNameKey = label.Key("net.peer.name")

	// Local host IP. Useful in case of a multi-IP host.
	NetHostIPKey = label.Key("net.host.ip")

	// Local host port.
	NetHostPortKey = label.Key("net.host.port")

	// Local hostname or similar.
	NetHostNameKey = label.Key("net.host.name")
)

var (
	NetTransportTCP    = NetTransportKey.String("IP.TCP")
	NetTransportUDP    = NetTransportKey.String("IP.UDP")
	NetTransportIP     = NetTransportKey.String("IP")
	NetTransportUnix   = NetTransportKey.String("Unix")
	NetTransportPipe   = NetTransportKey.String("pipe")
	NetTransportInProc = NetTransportKey.String("inproc")
	NetTransportOther  = NetTransportKey.String("other")
)

// General attribute keys for spans.
const (
	// Service name of the remote service. Should equal the actual
	// `service.name` resource attribute of the remote service, if any.
	PeerServiceKey = label.Key("peer.service")
)

// Semantic conventions for attribute keys used to identify an authorized
// user.
const (
	// Username or the client identifier extracted from the access token or
	// authorization header in the inbound request from outside the system.
	EnduserIDKey = label.Key("enduser.id")

	// Actual or assumed role the client is making the request with.
	EnduserRoleKey = label.Key("enduser.role")

	// Scopes or granted authorities the client currently possesses.
	EnduserScopeKey = label.Key("enduser.scope")
)

// Semantic conventions for attribute keys for HTTP.
const (
	// HTTP request method.
	HTTPMethodKey = label.Key("http.method")

	// Full HTTP request URL in the form:
	// scheme://host[:port]/path?query[#fragment].
	HTTPUrlKey = label.Key("http.url")

	// The full request target as passed in a HTTP request line or
	// equivalent, e.g. "/path/12314/?q=ddds#123".
	HTTPTargetKey = label.Key("http.target")

	// The value of the HTTP host header.
	HTTPHostKey = label.Key("http.host")

	// The URI scheme identifying the used protocol.
	HTTPSchemeKey = label.Key("http.scheme")

	// HTTP response status code.
	HTTPStatusCodeKey = label.Key("http.status_code")

	// Kind of HTTP protocol used.
	HTTPFlavorKey = label.Key("http.flavor")

	// Value of the HTTP User-Agent header sent by the client.
	HTTPUserAgentKey = label.Key("http.user_agent")

	// The primary server name of the matched virtual host.
	HTTPServerNameKey = label.Key("http.server_name")

	// The matched route served (path template). For example,
	// "/users/:userID?".
	HTTPRouteKey = label.Key("http.route")

	// The IP address of the original client behind all proxies, if known
	// (e.g. from X-Forwarded-For).
	HTTPClientIPKey = label.Key("http.client_ip")

	// The size of the request payload body in bytes.
	HTTPRequestContentLengthKey = label.Key("http.request_content_length")

	// The size of the uncompressed request payload body after transport decoding.
	// Not set if transport encoding not used.
	HTTPRequestContentLengthUncompressedKey = label.Key("http.request_content_length_uncompressed")

	// The size of the response payload body in bytes.
	HTTPResponseContentLengthKey = label.Key("http.response_content_length")

	// The size of the uncompressed response payload body after transport decoding.
	// Not set if transport encoding not used.
	HTTPResponseContentLengthUncompressedKey = label.Key("http.response_content_length_uncompressed")
)

var (
	HTTPSchemeHTTP  = HTTPSchemeKey.String("http")
	HTTPSchemeHTTPS = HTTPSchemeKey.String("https")

	HTTPFlavor1_0  = HTTPFlavorKey.String("1.0")
	HTTPFlavor1_1  = HTTPFlavorKey.String("1.1")
	HTTPFlavor2    = HTTPFlavorKey.String("2")
	HTTPFlavorSPDY = HTTPFlavorKey.String("SPDY")
	HTTPFlavorQUIC = HTTPFlavorKey.String("QUIC")
)

// Semantic conventions for attribute keys for database connections.
const (
	// Identifier for the database system (DBMS) being used.
	DBSystemKey = label.Key("db.system")

	// Database Connection String with embedded credentials removed.
	DBConnectionStringKey = label.Key("db.connection_string")

	// Username for accessing database.
	DBUserKey = label.Key("db.user")
)

var (
	DBSystemDB2       = DBSystemKey.String("db2")        // IBM DB2
	DBSystemDerby     = DBSystemKey.String("derby")      // Apache Derby
	DBSystemHive      = DBSystemKey.String("hive")       // Apache Hive
	DBSystemMariaDB   = DBSystemKey.String("mariadb")    // MariaDB
	DBSystemMSSql     = DBSystemKey.String("mssql")      // Microsoft SQL Server
	DBSystemMySQL     = DBSystemKey.String("mysql")      // MySQL
	DBSystemOracle    = DBSystemKey.String("oracle")     // Oracle Database
	DBSystemPostgres  = DBSystemKey.String("postgresql") // PostgreSQL
	DBSystemSqlite    = DBSystemKey.String("sqlite")     // SQLite
	DBSystemTeradata  = DBSystemKey.String("teradata")   // Teradata
	DBSystemOtherSQL  = DBSystemKey.String("other_sql")  // Some other Sql database. Fallback only
	DBSystemCassandra = DBSystemKey.String("cassandra")  // Cassandra
	DBSystemCosmosDB  = DBSystemKey.String("cosmosdb")   // Microsoft Azure CosmosDB
	DBSystemCouchbase = DBSystemKey.String("couchbase")  // Couchbase
	DBSystemCouchDB   = DBSystemKey.String("couchdb")    // CouchDB
	DBSystemDynamoDB  = DBSystemKey.String("dynamodb")   // Amazon DynamoDB
	DBSystemHBase     = DBSystemKey.String("hbase")      // HBase
	DBSystemMongodb   = DBSystemKey.String("mongodb")    // MongoDB
	DBSystemNeo4j     = DBSystemKey.String("neo4j")      // Neo4j
	DBSystemRedis     = DBSystemKey.String("redis")      // Redis
)

// Semantic conventions for attribute keys for database calls.
const (
	// Database instance name.
	DBNameKey = label.Key("db.name")

	// A database statement for the given database type.
	DBStatementKey = label.Key("db.statement")

	// A database operation for the given database type.
	DBOperationKey = label.Key("db.operation")
)

// Database technology-specific attributes
const (
	// Name of the Cassandra keyspace accessed. Use instead of `db.name`.
	DBCassandraKeyspaceKey = label.Key("db.cassandra.keyspace")

	// HBase namespace accessed. Use instead of `db.name`.
	DBHBaseNamespaceKey = label.Key("db.hbase.namespace")

	// Index of Redis database accessed. Use instead of `db.name`.
	DBRedisDBIndexKey = label.Key("db.redis.database_index")

	// Collection being accessed within the database in `db.name`.
	DBMongoDBCollectionKey = label.Key("db.mongodb.collection")
)

// Semantic conventions for attribute keys for RPC.
const (
	// A string identifying the remoting system.
	RPCSystemKey = label.Key("rpc.system")

	// The full name of the service being called.
	RPCServiceKey = label.Key("rpc.service")

	// The name of the method being called.
	RPCMethodKey = label.Key("rpc.method")

	// Name of message transmitted or received.
	RPCNameKey = label.Key("name")

	// Type of message transmitted or received.
	RPCMessageTypeKey = label.Key("message.type")

	// Identifier of message transmitted or received.
	RPCMessageIDKey = label.Key("message.id")

	// The compressed size of the message transmitted or received in bytes.
	RPCMessageCompressedSizeKey = label.Key("message.compressed_size")

	// The uncompressed size of the message transmitted or received in
	// bytes.
	RPCMessageUncompressedSizeKey = label.Key("message.uncompressed_size")
)

var (
	RPCSystemGRPC = RPCSystemKey.String("grpc")

	RPCNameMessage = RPCNameKey.String("message")

	RPCMessageTypeSent     = RPCMessageTypeKey.String("SENT")
	RPCMessageTypeReceived = RPCMessageTypeKey.String("RECEIVED")
)

// Semantic conventions for attribute keys for messaging systems.
const (
	// A unique identifier describing the messaging system. For example,
	// kafka, rabbitmq or activemq.
	MessagingSystemKey = label.Key("messaging.system")

	// The message destination name, e.g. MyQueue or MyTopic.
	MessagingDestinationKey = label.Key("messaging.destination")

	// The kind of message destination.
	MessagingDestinationKindKey = label.Key("messaging.destination_kind")

	// Describes if the destination is temporary or not.
	MessagingTempDestinationKey = label.Key("messaging.temp_destination")

	// The name of the transport protocol.
	MessagingProtocolKey = label.Key("messaging.protocol")

	// The version of the transport protocol.
	MessagingProtocolVersionKey = label.Key("messaging.protocol_version")

	// Messaging service URL.
	MessagingURLKey = label.Key("messaging.url")

	// Identifier used by the messaging system for a message.
	MessagingMessageIDKey = label.Key("messaging.message_id")

	// Identifier used by the messaging system for a conversation.
	MessagingConversationIDKey = label.Key("messaging.conversation_id")

	// The (uncompressed) size of the message payload in bytes.
	MessagingMessagePayloadSizeBytesKey = label.Key("messaging.message_payload_size_bytes")

	// The compressed size of the message payload in bytes.
	MessagingMessagePayloadCompressedSizeBytesKey = label.Key("messaging.message_payload_compressed_size_bytes")

	// Identifies which part and kind of message consumption is being
	// preformed.
	MessagingOperationKey = label.Key("messaging.operation")

	// RabbitMQ specific attribute describing the destination routing key.
	MessagingRabbitMQRoutingKeyKey = label.Key("messaging.rabbitmq.routing_key")
)

var (
	MessagingDestinationKindKeyQueue = MessagingDestinationKindKey.String("queue")
	MessagingDestinationKindKeyTopic = MessagingDestinationKindKey.String("topic")

	MessagingTempDestination = MessagingTempDestinationKey.Bool(true)

	MessagingOperationReceive = MessagingOperationKey.String("receive")
	MessagingOperationProcess = MessagingOperationKey.String("process")
)

// Semantic conventions for attribute keys for FaaS systems.
const (

	// Type of the trigger on which the function is executed.
	FaaSTriggerKey = label.Key("faas.trigger")

	// String containing the execution identifier of the function.
	FaaSExecutionKey = label.Key("faas.execution")

	// A boolean indicating that the serverless function is executed
	// for the first time (aka cold start).
	FaaSColdstartKey = label.Key("faas.coldstart")

	// The name of the source on which the operation was performed.
	// For example, in Cloud Storage or S3 corresponds to the bucket name,
	// and in Cosmos DB to the database name.
	FaaSDocumentCollectionKey = label.Key("faas.document.collection")

	// The type of the operation that was performed on the data.
	FaaSDocumentOperationKey = label.Key("faas.document.operation")

	// A string containing the time when the data was accessed.
	FaaSDocumentTimeKey = label.Key("faas.document.time")

	// The document name/table subjected to the operation.
	FaaSDocumentNameKey = label.Key("faas.document.name")

	// The function invocation time.
	FaaSTimeKey = label.Key("faas.time")

	// The schedule period as Cron Expression.
	FaaSCronKey = label.Key("faas.cron")
)

var (
	FaasTriggerDatasource = FaaSTriggerKey.String("datasource")
	FaasTriggerHTTP       = FaaSTriggerKey.String("http")
	FaasTriggerPubSub     = FaaSTriggerKey.String("pubsub")
	FaasTriggerTimer      = FaaSTriggerKey.String("timer")
	FaasTriggerOther      = FaaSTriggerKey.String("other")

	FaaSDocumentOperationInsert = FaaSDocumentOperationKey.String("insert")
	FaaSDocumentOperationEdit   = FaaSDocumentOperationKey.String("edit")
	FaaSDocumentOperationDelete = FaaSDocumentOperationKey.String("delete")
)
