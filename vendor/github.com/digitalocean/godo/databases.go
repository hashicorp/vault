package godo

import (
	"context"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"time"
)

const (
	databaseBasePath                    = "/v2/databases"
	databaseSinglePath                  = databaseBasePath + "/%s"
	databaseCAPath                      = databaseBasePath + "/%s/ca"
	databaseConfigPath                  = databaseBasePath + "/%s/config"
	databaseResizePath                  = databaseBasePath + "/%s/resize"
	databaseMigratePath                 = databaseBasePath + "/%s/migrate"
	databaseMaintenancePath             = databaseBasePath + "/%s/maintenance"
	databaseUpdateInstallationPath      = databaseBasePath + "/%s/install_update"
	databaseBackupsPath                 = databaseBasePath + "/%s/backups"
	databaseUsersPath                   = databaseBasePath + "/%s/users"
	databaseUserPath                    = databaseBasePath + "/%s/users/%s"
	databaseResetUserAuthPath           = databaseUserPath + "/reset_auth"
	databaseDBPath                      = databaseBasePath + "/%s/dbs/%s"
	databaseDBsPath                     = databaseBasePath + "/%s/dbs"
	databasePoolPath                    = databaseBasePath + "/%s/pools/%s"
	databasePoolsPath                   = databaseBasePath + "/%s/pools"
	databaseReplicaPath                 = databaseBasePath + "/%s/replicas/%s"
	databaseReplicasPath                = databaseBasePath + "/%s/replicas"
	databaseEvictionPolicyPath          = databaseBasePath + "/%s/eviction_policy"
	databaseSQLModePath                 = databaseBasePath + "/%s/sql_mode"
	databaseFirewallRulesPath           = databaseBasePath + "/%s/firewall"
	databaseOptionsPath                 = databaseBasePath + "/options"
	databaseUpgradeMajorVersionPath     = databaseBasePath + "/%s/upgrade"
	databasePromoteReplicaToPrimaryPath = databaseReplicaPath + "/promote"
	databaseTopicPath                   = databaseBasePath + "/%s/topics/%s"
	databaseTopicsPath                  = databaseBasePath + "/%s/topics"
	databaseMetricsCredentialsPath      = databaseBasePath + "/metrics/credentials"
	databaseEvents                      = databaseBasePath + "/%s/events"
	databaseIndexesPath                 = databaseBasePath + "/%s/indexes"
	databaseIndexPath                   = databaseBasePath + "/%s/indexes/%s"
	databaseLogsinkPath                 = databaseBasePath + "/%s/logsink/%s"
	databaseLogsinksPath                = databaseBasePath + "/%s/logsink"
)

// SQL Mode constants allow for MySQL-specific SQL flavor configuration.
const (
	SQLModeAllowInvalidDates     = "ALLOW_INVALID_DATES"
	SQLModeANSIQuotes            = "ANSI_QUOTES"
	SQLModeHighNotPrecedence     = "HIGH_NOT_PRECEDENCE"
	SQLModeIgnoreSpace           = "IGNORE_SPACE"
	SQLModeNoAuthCreateUser      = "NO_AUTO_CREATE_USER"
	SQLModeNoAutoValueOnZero     = "NO_AUTO_VALUE_ON_ZERO"
	SQLModeNoBackslashEscapes    = "NO_BACKSLASH_ESCAPES"
	SQLModeNoDirInCreate         = "NO_DIR_IN_CREATE"
	SQLModeNoEngineSubstitution  = "NO_ENGINE_SUBSTITUTION"
	SQLModeNoFieldOptions        = "NO_FIELD_OPTIONS"
	SQLModeNoKeyOptions          = "NO_KEY_OPTIONS"
	SQLModeNoTableOptions        = "NO_TABLE_OPTIONS"
	SQLModeNoUnsignedSubtraction = "NO_UNSIGNED_SUBTRACTION"
	SQLModeNoZeroDate            = "NO_ZERO_DATE"
	SQLModeNoZeroInDate          = "NO_ZERO_IN_DATE"
	SQLModeOnlyFullGroupBy       = "ONLY_FULL_GROUP_BY"
	SQLModePadCharToFullLength   = "PAD_CHAR_TO_FULL_LENGTH"
	SQLModePipesAsConcat         = "PIPES_AS_CONCAT"
	SQLModeRealAsFloat           = "REAL_AS_FLOAT"
	SQLModeStrictAllTables       = "STRICT_ALL_TABLES"
	SQLModeStrictTransTables     = "STRICT_TRANS_TABLES"
	SQLModeANSI                  = "ANSI"
	SQLModeDB2                   = "DB2"
	SQLModeMaxDB                 = "MAXDB"
	SQLModeMSSQL                 = "MSSQL"
	SQLModeMYSQL323              = "MYSQL323"
	SQLModeMYSQL40               = "MYSQL40"
	SQLModeOracle                = "ORACLE"
	SQLModePostgreSQL            = "POSTGRESQL"
	SQLModeTraditional           = "TRADITIONAL"
)

// SQL Auth constants allow for MySQL-specific user auth plugins
const (
	SQLAuthPluginNative      = "mysql_native_password"
	SQLAuthPluginCachingSHA2 = "caching_sha2_password"
)

// Redis eviction policies supported by the managed Redis product.
const (
	EvictionPolicyNoEviction     = "noeviction"
	EvictionPolicyAllKeysLRU     = "allkeys_lru"
	EvictionPolicyAllKeysRandom  = "allkeys_random"
	EvictionPolicyVolatileLRU    = "volatile_lru"
	EvictionPolicyVolatileRandom = "volatile_random"
	EvictionPolicyVolatileTTL    = "volatile_ttl"
)

// evictionPolicyMap is used to normalize the eviction policy string in requests
// to the advanced Redis configuration endpoint from the consts used with SetEvictionPolicy.
var evictionPolicyMap = map[string]string{
	EvictionPolicyAllKeysLRU:     "allkeys-lru",
	EvictionPolicyAllKeysRandom:  "allkeys-random",
	EvictionPolicyVolatileLRU:    "volatile-lru",
	EvictionPolicyVolatileRandom: "volatile-random",
	EvictionPolicyVolatileTTL:    "volatile-ttl",
}

// The DatabasesService provides access to the DigitalOcean managed database
// suite of products through the public API. Customers can create new database
// clusters, migrate them  between regions, create replicas and interact with
// their configurations. Each database service is referred to as a Database. A
// SQL database service can have multiple databases residing in the system. To
// help make these entities distinct from Databases in godo, we refer to them
// here as DatabaseDBs.
//
// See: https://docs.digitalocean.com/reference/api/api-reference/#tag/Databases
type DatabasesService interface {
	List(context.Context, *ListOptions) ([]Database, *Response, error)
	Get(context.Context, string) (*Database, *Response, error)
	GetCA(context.Context, string) (*DatabaseCA, *Response, error)
	Create(context.Context, *DatabaseCreateRequest) (*Database, *Response, error)
	Delete(context.Context, string) (*Response, error)
	Resize(context.Context, string, *DatabaseResizeRequest) (*Response, error)
	Migrate(context.Context, string, *DatabaseMigrateRequest) (*Response, error)
	UpdateMaintenance(context.Context, string, *DatabaseUpdateMaintenanceRequest) (*Response, error)
	InstallUpdate(context.Context, string) (*Response, error)
	ListBackups(context.Context, string, *ListOptions) ([]DatabaseBackup, *Response, error)
	GetUser(context.Context, string, string) (*DatabaseUser, *Response, error)
	ListUsers(context.Context, string, *ListOptions) ([]DatabaseUser, *Response, error)
	CreateUser(context.Context, string, *DatabaseCreateUserRequest) (*DatabaseUser, *Response, error)
	UpdateUser(context.Context, string, string, *DatabaseUpdateUserRequest) (*DatabaseUser, *Response, error)
	DeleteUser(context.Context, string, string) (*Response, error)
	ResetUserAuth(context.Context, string, string, *DatabaseResetUserAuthRequest) (*DatabaseUser, *Response, error)
	ListDBs(context.Context, string, *ListOptions) ([]DatabaseDB, *Response, error)
	CreateDB(context.Context, string, *DatabaseCreateDBRequest) (*DatabaseDB, *Response, error)
	GetDB(context.Context, string, string) (*DatabaseDB, *Response, error)
	DeleteDB(context.Context, string, string) (*Response, error)
	ListPools(context.Context, string, *ListOptions) ([]DatabasePool, *Response, error)
	CreatePool(context.Context, string, *DatabaseCreatePoolRequest) (*DatabasePool, *Response, error)
	GetPool(context.Context, string, string) (*DatabasePool, *Response, error)
	DeletePool(context.Context, string, string) (*Response, error)
	UpdatePool(context.Context, string, string, *DatabaseUpdatePoolRequest) (*Response, error)
	GetReplica(context.Context, string, string) (*DatabaseReplica, *Response, error)
	ListReplicas(context.Context, string, *ListOptions) ([]DatabaseReplica, *Response, error)
	CreateReplica(context.Context, string, *DatabaseCreateReplicaRequest) (*DatabaseReplica, *Response, error)
	DeleteReplica(context.Context, string, string) (*Response, error)
	PromoteReplicaToPrimary(context.Context, string, string) (*Response, error)
	GetEvictionPolicy(context.Context, string) (string, *Response, error)
	SetEvictionPolicy(context.Context, string, string) (*Response, error)
	GetSQLMode(context.Context, string) (string, *Response, error)
	SetSQLMode(context.Context, string, ...string) (*Response, error)
	GetFirewallRules(context.Context, string) ([]DatabaseFirewallRule, *Response, error)
	UpdateFirewallRules(context.Context, string, *DatabaseUpdateFirewallRulesRequest) (*Response, error)
	GetPostgreSQLConfig(context.Context, string) (*PostgreSQLConfig, *Response, error)
	GetRedisConfig(context.Context, string) (*RedisConfig, *Response, error)
	GetMySQLConfig(context.Context, string) (*MySQLConfig, *Response, error)
	GetMongoDBConfig(context.Context, string) (*MongoDBConfig, *Response, error)
	GetOpensearchConfig(context.Context, string) (*OpensearchConfig, *Response, error)
	GetKafkaConfig(context.Context, string) (*KafkaConfig, *Response, error)
	UpdatePostgreSQLConfig(context.Context, string, *PostgreSQLConfig) (*Response, error)
	UpdateRedisConfig(context.Context, string, *RedisConfig) (*Response, error)
	UpdateMySQLConfig(context.Context, string, *MySQLConfig) (*Response, error)
	UpdateMongoDBConfig(context.Context, string, *MongoDBConfig) (*Response, error)
	UpdateOpensearchConfig(context.Context, string, *OpensearchConfig) (*Response, error)
	UpdateKafkaConfig(context.Context, string, *KafkaConfig) (*Response, error)
	ListOptions(todo context.Context) (*DatabaseOptions, *Response, error)
	UpgradeMajorVersion(context.Context, string, *UpgradeVersionRequest) (*Response, error)
	ListTopics(context.Context, string, *ListOptions) ([]DatabaseTopic, *Response, error)
	CreateTopic(context.Context, string, *DatabaseCreateTopicRequest) (*DatabaseTopic, *Response, error)
	GetTopic(context.Context, string, string) (*DatabaseTopic, *Response, error)
	DeleteTopic(context.Context, string, string) (*Response, error)
	UpdateTopic(context.Context, string, string, *DatabaseUpdateTopicRequest) (*Response, error)
	GetMetricsCredentials(context.Context) (*DatabaseMetricsCredentials, *Response, error)
	UpdateMetricsCredentials(context.Context, *DatabaseUpdateMetricsCredentialsRequest) (*Response, error)
	ListDatabaseEvents(context.Context, string, *ListOptions) ([]DatabaseEvent, *Response, error)
	ListIndexes(context.Context, string, *ListOptions) ([]DatabaseIndex, *Response, error)
	DeleteIndex(context.Context, string, string) (*Response, error)
	CreateLogsink(ctx context.Context, databaseID string, createLogsink *DatabaseCreateLogsinkRequest) (*DatabaseLogsink, *Response, error)
	GetLogsink(ctx context.Context, databaseID string, logsinkID string) (*DatabaseLogsink, *Response, error)
	ListLogsinks(ctx context.Context, databaseID string, opts *ListOptions) ([]DatabaseLogsink, *Response, error)
	UpdateLogsink(ctx context.Context, databaseID string, logsinkID string, updateLogsink *DatabaseUpdateLogsinkRequest) (*Response, error)
	DeleteLogsink(ctx context.Context, databaseID, logsinkID string) (*Response, error)
}

// DatabasesServiceOp handles communication with the Databases related methods
// of the DigitalOcean API.
type DatabasesServiceOp struct {
	client *Client
}

var _ DatabasesService = &DatabasesServiceOp{}

// Database represents a DigitalOcean managed database product. These managed databases
// are usually comprised of a cluster of database nodes, a primary and 0 or more replicas.
// The EngineSlug is a string which indicates the type of database service. Some examples are
// "pg", "mysql" or "redis". A Database also includes connection information and other
// properties of the service like region, size and current status.
type Database struct {
	ID                       string                     `json:"id,omitempty"`
	Name                     string                     `json:"name,omitempty"`
	EngineSlug               string                     `json:"engine,omitempty"`
	VersionSlug              string                     `json:"version,omitempty"`
	Connection               *DatabaseConnection        `json:"connection,omitempty"`
	UIConnection             *DatabaseConnection        `json:"ui_connection,omitempty"`
	PrivateConnection        *DatabaseConnection        `json:"private_connection,omitempty"`
	StandbyConnection        *DatabaseConnection        `json:"standby_connection,omitempty"`
	StandbyPrivateConnection *DatabaseConnection        `json:"standby_private_connection,omitempty"`
	Users                    []DatabaseUser             `json:"users,omitempty"`
	NumNodes                 int                        `json:"num_nodes,omitempty"`
	SizeSlug                 string                     `json:"size,omitempty"`
	DBNames                  []string                   `json:"db_names,omitempty"`
	RegionSlug               string                     `json:"region,omitempty"`
	Status                   string                     `json:"status,omitempty"`
	MaintenanceWindow        *DatabaseMaintenanceWindow `json:"maintenance_window,omitempty"`
	CreatedAt                time.Time                  `json:"created_at,omitempty"`
	PrivateNetworkUUID       string                     `json:"private_network_uuid,omitempty"`
	Tags                     []string                   `json:"tags,omitempty"`
	ProjectID                string                     `json:"project_id,omitempty"`
	StorageSizeMib           uint64                     `json:"storage_size_mib,omitempty"`
	MetricsEndpoints         []*ServiceAddress          `json:"metrics_endpoints,omitempty"`
}

// DatabaseCA represents a database ca.
type DatabaseCA struct {
	Certificate []byte `json:"certificate"`
}

// DatabaseConnection represents a database connection
type DatabaseConnection struct {
	Protocol         string            `json:"protocol"`
	URI              string            `json:"uri,omitempty"`
	Database         string            `json:"database,omitempty"`
	Host             string            `json:"host,omitempty"`
	Port             int               `json:"port,omitempty"`
	User             string            `json:"user,omitempty"`
	Password         string            `json:"password,omitempty"`
	SSL              bool              `json:"ssl,omitempty"`
	ApplicationPorts map[string]uint32 `json:"application_ports,omitempty"`
}

// ServiceAddress represents a host:port for a generic service (e.g. metrics endpoint)
type ServiceAddress struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

// DatabaseUser represents a user in the database
type DatabaseUser struct {
	Name          string                     `json:"name,omitempty"`
	Role          string                     `json:"role,omitempty"`
	Password      string                     `json:"password,omitempty"`
	AccessCert    string                     `json:"access_cert,omitempty"`
	AccessKey     string                     `json:"access_key,omitempty"`
	MySQLSettings *DatabaseMySQLUserSettings `json:"mysql_settings,omitempty"`
	Settings      *DatabaseUserSettings      `json:"settings,omitempty"`
}

// KafkaACL contains Kafka specific user access control information
type KafkaACL struct {
	ID         string `json:"id,omitempty"`
	Permission string `json:"permission,omitempty"`
	Topic      string `json:"topic,omitempty"`
}

// OpenSearchACL contains OpenSearch specific user access control information
type OpenSearchACL struct {
	Permission string `json:"permission,omitempty"`
	Index      string `json:"index,omitempty"`
}

// DatabaseUserSettings contains user settings
type DatabaseUserSettings struct {
	ACL           []*KafkaACL      `json:"acl,omitempty"`
	OpenSearchACL []*OpenSearchACL `json:"opensearch_acl,omitempty"`
}

// DatabaseMySQLUserSettings contains MySQL-specific user settings
type DatabaseMySQLUserSettings struct {
	AuthPlugin string `json:"auth_plugin"`
}

// DatabaseMaintenanceWindow represents the maintenance_window of a database
// cluster
type DatabaseMaintenanceWindow struct {
	Day         string   `json:"day,omitempty"`
	Hour        string   `json:"hour,omitempty"`
	Pending     bool     `json:"pending,omitempty"`
	Description []string `json:"description,omitempty"`
}

// DatabaseBackup represents a database backup.
type DatabaseBackup struct {
	CreatedAt     time.Time `json:"created_at,omitempty"`
	SizeGigabytes float64   `json:"size_gigabytes,omitempty"`
}

// DatabaseBackupRestore contains information needed to restore a backup.
type DatabaseBackupRestore struct {
	DatabaseName    string `json:"database_name,omitempty"`
	BackupCreatedAt string `json:"backup_created_at,omitempty"`
}

// DatabaseCreateRequest represents a request to create a database cluster
type DatabaseCreateRequest struct {
	Name               string                 `json:"name,omitempty"`
	EngineSlug         string                 `json:"engine,omitempty"`
	Version            string                 `json:"version,omitempty"`
	SizeSlug           string                 `json:"size,omitempty"`
	Region             string                 `json:"region,omitempty"`
	NumNodes           int                    `json:"num_nodes,omitempty"`
	PrivateNetworkUUID string                 `json:"private_network_uuid"`
	Tags               []string               `json:"tags,omitempty"`
	BackupRestore      *DatabaseBackupRestore `json:"backup_restore,omitempty"`
	ProjectID          string                 `json:"project_id"`
	StorageSizeMib     uint64                 `json:"storage_size_mib,omitempty"`
}

// DatabaseResizeRequest can be used to initiate a database resize operation.
type DatabaseResizeRequest struct {
	SizeSlug       string `json:"size,omitempty"`
	NumNodes       int    `json:"num_nodes,omitempty"`
	StorageSizeMib uint64 `json:"storage_size_mib,omitempty"`
}

// DatabaseMigrateRequest can be used to initiate a database migrate operation.
type DatabaseMigrateRequest struct {
	Region             string `json:"region,omitempty"`
	PrivateNetworkUUID string `json:"private_network_uuid"`
}

// DatabaseUpdateMaintenanceRequest can be used to update the database's maintenance window.
type DatabaseUpdateMaintenanceRequest struct {
	Day  string `json:"day,omitempty"`
	Hour string `json:"hour,omitempty"`
}

// DatabaseDB represents an engine-specific database created within a database cluster. For SQL
// databases like PostgreSQL or MySQL, a "DB" refers to a database created on the RDBMS. For instance,
// a PostgreSQL database server can contain many database schemas, each with its own settings, access
// permissions and data. ListDBs will return all databases present on the server.
type DatabaseDB struct {
	Name string `json:"name"`
}

// DatabaseTopic represents a Kafka topic
type DatabaseTopic struct {
	Name              string            `json:"name"`
	Partitions        []*TopicPartition `json:"partitions,omitempty"`
	ReplicationFactor *uint32           `json:"replication_factor,omitempty"`
	State             string            `json:"state,omitempty"`
	Config            *TopicConfig      `json:"config,omitempty"`
}

// DatabaseLogsink represents a logsink
type DatabaseLogsink struct {
	ID     string                 `json:"sink_id"`
	Name   string                 `json:"sink_name,omitempty"`
	Type   string                 `json:"sink_type,omitempty"`
	Config *DatabaseLogsinkConfig `json:"config,omitempty"`
}

// TopicPartition represents the state of a Kafka topic partition
type TopicPartition struct {
	EarliestOffset uint64                `json:"earliest_offset,omitempty"`
	InSyncReplicas uint32                `json:"in_sync_replicas,omitempty"`
	Id             uint32                `json:"id,omitempty"`
	Size           uint64                `json:"size,omitempty"`
	ConsumerGroups []*TopicConsumerGroup `json:"consumer_groups,omitempty"`
}

// TopicConsumerGroup represents a consumer group for a particular Kafka topic
type TopicConsumerGroup struct {
	Name   string `json:"name,omitempty"`
	Offset uint64 `json:"offset,omitempty"`
}

// TopicConfig represents all configurable options for a Kafka topic
type TopicConfig struct {
	CleanupPolicy                   string   `json:"cleanup_policy,omitempty"`
	CompressionType                 string   `json:"compression_type,omitempty"`
	DeleteRetentionMS               *uint64  `json:"delete_retention_ms,omitempty"`
	FileDeleteDelayMS               *uint64  `json:"file_delete_delay_ms,omitempty"`
	FlushMessages                   *uint64  `json:"flush_messages,omitempty"`
	FlushMS                         *uint64  `json:"flush_ms,omitempty"`
	IndexIntervalBytes              *uint64  `json:"index_interval_bytes,omitempty"`
	MaxCompactionLagMS              *uint64  `json:"max_compaction_lag_ms,omitempty"`
	MaxMessageBytes                 *uint64  `json:"max_message_bytes,omitempty"`
	MessageDownConversionEnable     *bool    `json:"message_down_conversion_enable,omitempty"`
	MessageFormatVersion            string   `json:"message_format_version,omitempty"`
	MessageTimestampDifferenceMaxMS *uint64  `json:"message_timestamp_difference_max_ms,omitempty"`
	MessageTimestampType            string   `json:"message_timestamp_type,omitempty"`
	MinCleanableDirtyRatio          *float32 `json:"min_cleanable_dirty_ratio,omitempty"`
	MinCompactionLagMS              *uint64  `json:"min_compaction_lag_ms,omitempty"`
	MinInsyncReplicas               *uint32  `json:"min_insync_replicas,omitempty"`
	Preallocate                     *bool    `json:"preallocate,omitempty"`
	RetentionBytes                  *int64   `json:"retention_bytes,omitempty"`
	RetentionMS                     *int64   `json:"retention_ms,omitempty"`
	SegmentBytes                    *uint64  `json:"segment_bytes,omitempty"`
	SegmentIndexBytes               *uint64  `json:"segment_index_bytes,omitempty"`
	SegmentJitterMS                 *uint64  `json:"segment_jitter_ms,omitempty"`
	SegmentMS                       *uint64  `json:"segment_ms,omitempty"`
}

// DatabaseCreateTopicRequest is used to create a new topic within a kafka cluster
type DatabaseCreateTopicRequest struct {
	Name              string       `json:"name"`
	PartitionCount    *uint32      `json:"partition_count,omitempty"`
	ReplicationFactor *uint32      `json:"replication_factor,omitempty"`
	Config            *TopicConfig `json:"config,omitempty"`
}

// DatabaseUpdateTopicRequest ...
type DatabaseUpdateTopicRequest struct {
	PartitionCount    *uint32      `json:"partition_count,omitempty"`
	ReplicationFactor *uint32      `json:"replication_factor,omitempty"`
	Config            *TopicConfig `json:"config,omitempty"`
}

// DatabaseReplica represents a read-only replica of a particular database
type DatabaseReplica struct {
	ID                 string              `json:"id"`
	Name               string              `json:"name"`
	Connection         *DatabaseConnection `json:"connection"`
	PrivateConnection  *DatabaseConnection `json:"private_connection,omitempty"`
	Region             string              `json:"region"`
	Status             string              `json:"status"`
	CreatedAt          time.Time           `json:"created_at"`
	PrivateNetworkUUID string              `json:"private_network_uuid,omitempty"`
	Tags               []string            `json:"tags,omitempty"`
	StorageSizeMib     uint64              `json:"storage_size_mib,omitempty"`
	Size               string              `json:"size"`
}

// DatabasePool represents a database connection pool
type DatabasePool struct {
	User                     string              `json:"user"`
	Name                     string              `json:"name"`
	Size                     int                 `json:"size"`
	Database                 string              `json:"db"`
	Mode                     string              `json:"mode"`
	Connection               *DatabaseConnection `json:"connection"`
	PrivateConnection        *DatabaseConnection `json:"private_connection,omitempty"`
	StandbyConnection        *DatabaseConnection `json:"standby_connection,omitempty"`
	StandbyPrivateConnection *DatabaseConnection `json:"standby_private_connection,omitempty"`
}

// DatabaseCreatePoolRequest is used to create a new database connection pool
type DatabaseCreatePoolRequest struct {
	User     string `json:"user"`
	Name     string `json:"name"`
	Size     int    `json:"size"`
	Database string `json:"db"`
	Mode     string `json:"mode"`
}

// DatabaseUpdatePoolRequest is used to update a database connection pool
type DatabaseUpdatePoolRequest struct {
	User     string `json:"user,omitempty"`
	Size     int    `json:"size"`
	Database string `json:"db"`
	Mode     string `json:"mode"`
}

// DatabaseCreateUserRequest is used to create a new database user
type DatabaseCreateUserRequest struct {
	Name          string                     `json:"name"`
	MySQLSettings *DatabaseMySQLUserSettings `json:"mysql_settings,omitempty"`
	Settings      *DatabaseUserSettings      `json:"settings,omitempty"`
}

// DatabaseUpdateUserRequest is used to update an existing database user
type DatabaseUpdateUserRequest struct {
	Settings *DatabaseUserSettings `json:"settings,omitempty"`
}

// DatabaseResetUserAuthRequest is used to reset a users DB auth
type DatabaseResetUserAuthRequest struct {
	MySQLSettings *DatabaseMySQLUserSettings `json:"mysql_settings,omitempty"`
	Settings      *DatabaseUserSettings      `json:"settings,omitempty"`
}

// DatabaseCreateDBRequest is used to create a new engine-specific database within the cluster
type DatabaseCreateDBRequest struct {
	Name string `json:"name"`
}

// DatabaseCreateReplicaRequest is used to create a new read-only replica
type DatabaseCreateReplicaRequest struct {
	Name               string   `json:"name"`
	Region             string   `json:"region"`
	Size               string   `json:"size"`
	PrivateNetworkUUID string   `json:"private_network_uuid"`
	Tags               []string `json:"tags,omitempty"`
	StorageSizeMib     uint64   `json:"storage_size_mib,omitempty"`
}

// DatabaseUpdateFirewallRulesRequest is used to set the firewall rules for a database
type DatabaseUpdateFirewallRulesRequest struct {
	Rules []*DatabaseFirewallRule `json:"rules"`
}

// DatabaseFirewallRule is a rule describing an inbound source to a database
type DatabaseFirewallRule struct {
	UUID        string    `json:"uuid"`
	ClusterUUID string    `json:"cluster_uuid"`
	Type        string    `json:"type"`
	Value       string    `json:"value"`
	CreatedAt   time.Time `json:"created_at"`
}

// DatabaseCreateLogsinkRequest is used to create logsink for a database cluster
type DatabaseCreateLogsinkRequest struct {
	Name   string                 `json:"sink_name"`
	Type   string                 `json:"sink_type"`
	Config *DatabaseLogsinkConfig `json:"config"`
}

// DatabaseUpdateLogsinkRequest is used to update logsink for a database cluster
type DatabaseUpdateLogsinkRequest struct {
	Config *DatabaseLogsinkConfig `json:"config"`
}

// DatabaseLogsinkConfig represents one of the configurable options (rsyslog_logsink, elasticsearch_logsink, or opensearch_logsink) for a logsink.
type DatabaseLogsinkConfig struct {
	URL          string  `json:"url,omitempty"`
	IndexPrefix  string  `json:"index_prefix,omitempty"`
	IndexDaysMax int     `json:"index_days_max,omitempty"`
	Timeout      float32 `json:"timeout,omitempty"`
	Server       string  `json:"server,omitempty"`
	Port         int     `json:"port,omitempty"`
	TLS          bool    `json:"tls,omitempty"`
	Format       string  `json:"format,omitempty"`
	Logline      string  `json:"logline,omitempty"`
	SD           string  `json:"sd,omitempty"`
	CA           string  `json:"ca,omitempty"`
	Key          string  `json:"key,omitempty"`
	Cert         string  `json:"cert,omitempty"`
}

// PostgreSQLConfig holds advanced configurations for PostgreSQL database clusters.
type PostgreSQLConfig struct {
	AutovacuumFreezeMaxAge          *int                         `json:"autovacuum_freeze_max_age,omitempty"`
	AutovacuumMaxWorkers            *int                         `json:"autovacuum_max_workers,omitempty"`
	AutovacuumNaptime               *int                         `json:"autovacuum_naptime,omitempty"`
	AutovacuumVacuumThreshold       *int                         `json:"autovacuum_vacuum_threshold,omitempty"`
	AutovacuumAnalyzeThreshold      *int                         `json:"autovacuum_analyze_threshold,omitempty"`
	AutovacuumVacuumScaleFactor     *float32                     `json:"autovacuum_vacuum_scale_factor,omitempty"`
	AutovacuumAnalyzeScaleFactor    *float32                     `json:"autovacuum_analyze_scale_factor,omitempty"`
	AutovacuumVacuumCostDelay       *int                         `json:"autovacuum_vacuum_cost_delay,omitempty"`
	AutovacuumVacuumCostLimit       *int                         `json:"autovacuum_vacuum_cost_limit,omitempty"`
	BGWriterDelay                   *int                         `json:"bgwriter_delay,omitempty"`
	BGWriterFlushAfter              *int                         `json:"bgwriter_flush_after,omitempty"`
	BGWriterLRUMaxpages             *int                         `json:"bgwriter_lru_maxpages,omitempty"`
	BGWriterLRUMultiplier           *float32                     `json:"bgwriter_lru_multiplier,omitempty"`
	DeadlockTimeoutMillis           *int                         `json:"deadlock_timeout,omitempty"`
	DefaultToastCompression         *string                      `json:"default_toast_compression,omitempty"`
	IdleInTransactionSessionTimeout *int                         `json:"idle_in_transaction_session_timeout,omitempty"`
	JIT                             *bool                        `json:"jit,omitempty"`
	LogAutovacuumMinDuration        *int                         `json:"log_autovacuum_min_duration,omitempty"`
	LogErrorVerbosity               *string                      `json:"log_error_verbosity,omitempty"`
	LogLinePrefix                   *string                      `json:"log_line_prefix,omitempty"`
	LogMinDurationStatement         *int                         `json:"log_min_duration_statement,omitempty"`
	MaxFilesPerProcess              *int                         `json:"max_files_per_process,omitempty"`
	MaxPreparedTransactions         *int                         `json:"max_prepared_transactions,omitempty"`
	MaxPredLocksPerTransaction      *int                         `json:"max_pred_locks_per_transaction,omitempty"`
	MaxLocksPerTransaction          *int                         `json:"max_locks_per_transaction,omitempty"`
	MaxStackDepth                   *int                         `json:"max_stack_depth,omitempty"`
	MaxStandbyArchiveDelay          *int                         `json:"max_standby_archive_delay,omitempty"`
	MaxStandbyStreamingDelay        *int                         `json:"max_standby_streaming_delay,omitempty"`
	MaxReplicationSlots             *int                         `json:"max_replication_slots,omitempty"`
	MaxLogicalReplicationWorkers    *int                         `json:"max_logical_replication_workers,omitempty"`
	MaxParallelWorkers              *int                         `json:"max_parallel_workers,omitempty"`
	MaxParallelWorkersPerGather     *int                         `json:"max_parallel_workers_per_gather,omitempty"`
	MaxWorkerProcesses              *int                         `json:"max_worker_processes,omitempty"`
	PGPartmanBGWRole                *string                      `json:"pg_partman_bgw.role,omitempty"`
	PGPartmanBGWInterval            *int                         `json:"pg_partman_bgw.interval,omitempty"`
	PGStatStatementsTrack           *string                      `json:"pg_stat_statements.track,omitempty"`
	TempFileLimit                   *int                         `json:"temp_file_limit,omitempty"`
	Timezone                        *string                      `json:"timezone,omitempty"`
	TrackActivityQuerySize          *int                         `json:"track_activity_query_size,omitempty"`
	TrackCommitTimestamp            *string                      `json:"track_commit_timestamp,omitempty"`
	TrackFunctions                  *string                      `json:"track_functions,omitempty"`
	TrackIOTiming                   *string                      `json:"track_io_timing,omitempty"`
	MaxWalSenders                   *int                         `json:"max_wal_senders,omitempty"`
	WalSenderTimeout                *int                         `json:"wal_sender_timeout,omitempty"`
	WalWriterDelay                  *int                         `json:"wal_writer_delay,omitempty"`
	SharedBuffersPercentage         *float32                     `json:"shared_buffers_percentage,omitempty"`
	PgBouncer                       *PostgreSQLBouncerConfig     `json:"pgbouncer,omitempty"`
	BackupHour                      *int                         `json:"backup_hour,omitempty"`
	BackupMinute                    *int                         `json:"backup_minute,omitempty"`
	WorkMem                         *int                         `json:"work_mem,omitempty"`
	TimeScaleDB                     *PostgreSQLTimeScaleDBConfig `json:"timescaledb,omitempty"`
	SynchronousReplication          *string                      `json:"synchronous_replication,omitempty"`
	StatMonitorEnable               *bool                        `json:"stat_monitor_enable,omitempty"`
	MaxFailoverReplicationTimeLag   *int64                       `json:"max_failover_replication_time_lag,omitempty"`
}

// PostgreSQLBouncerConfig configuration
type PostgreSQLBouncerConfig struct {
	ServerResetQueryAlways  *bool     `json:"server_reset_query_always,omitempty"`
	IgnoreStartupParameters *[]string `json:"ignore_startup_parameters,omitempty"`
	MinPoolSize             *int      `json:"min_pool_size,omitempty"`
	ServerLifetime          *int      `json:"server_lifetime,omitempty"`
	ServerIdleTimeout       *int      `json:"server_idle_timeout,omitempty"`
	AutodbPoolSize          *int      `json:"autodb_pool_size,omitempty"`
	AutodbPoolMode          *string   `json:"autodb_pool_mode,omitempty"`
	AutodbMaxDbConnections  *int      `json:"autodb_max_db_connections,omitempty"`
	AutodbIdleTimeout       *int      `json:"autodb_idle_timeout,omitempty"`
}

// PostgreSQLTimeScaleDBConfig configuration
type PostgreSQLTimeScaleDBConfig struct {
	MaxBackgroundWorkers *int `json:"max_background_workers,omitempty"`
}

// RedisConfig holds advanced configurations for Redis database clusters.
type RedisConfig struct {
	RedisMaxmemoryPolicy               *string `json:"redis_maxmemory_policy,omitempty"`
	RedisPubsubClientOutputBufferLimit *int    `json:"redis_pubsub_client_output_buffer_limit,omitempty"`
	RedisNumberOfDatabases             *int    `json:"redis_number_of_databases,omitempty"`
	RedisIOThreads                     *int    `json:"redis_io_threads,omitempty"`
	RedisLFULogFactor                  *int    `json:"redis_lfu_log_factor,omitempty"`
	RedisLFUDecayTime                  *int    `json:"redis_lfu_decay_time,omitempty"`
	RedisSSL                           *bool   `json:"redis_ssl,omitempty"`
	RedisTimeout                       *int    `json:"redis_timeout,omitempty"`
	RedisNotifyKeyspaceEvents          *string `json:"redis_notify_keyspace_events,omitempty"`
	RedisPersistence                   *string `json:"redis_persistence,omitempty"`
	RedisACLChannelsDefault            *string `json:"redis_acl_channels_default,omitempty"`
}

// MySQLConfig holds advanced configurations for MySQL database clusters.
type MySQLConfig struct {
	ConnectTimeout               *int     `json:"connect_timeout,omitempty"`
	DefaultTimeZone              *string  `json:"default_time_zone,omitempty"`
	InnodbLogBufferSize          *int     `json:"innodb_log_buffer_size,omitempty"`
	InnodbOnlineAlterLogMaxSize  *int     `json:"innodb_online_alter_log_max_size,omitempty"`
	InnodbLockWaitTimeout        *int     `json:"innodb_lock_wait_timeout,omitempty"`
	InteractiveTimeout           *int     `json:"interactive_timeout,omitempty"`
	MaxAllowedPacket             *int     `json:"max_allowed_packet,omitempty"`
	NetReadTimeout               *int     `json:"net_read_timeout,omitempty"`
	SortBufferSize               *int     `json:"sort_buffer_size,omitempty"`
	SQLMode                      *string  `json:"sql_mode,omitempty"`
	SQLRequirePrimaryKey         *bool    `json:"sql_require_primary_key,omitempty"`
	WaitTimeout                  *int     `json:"wait_timeout,omitempty"`
	NetWriteTimeout              *int     `json:"net_write_timeout,omitempty"`
	GroupConcatMaxLen            *int     `json:"group_concat_max_len,omitempty"`
	InformationSchemaStatsExpiry *int     `json:"information_schema_stats_expiry,omitempty"`
	InnodbFtMinTokenSize         *int     `json:"innodb_ft_min_token_size,omitempty"`
	InnodbFtServerStopwordTable  *string  `json:"innodb_ft_server_stopword_table,omitempty"`
	InnodbPrintAllDeadlocks      *bool    `json:"innodb_print_all_deadlocks,omitempty"`
	InnodbRollbackOnTimeout      *bool    `json:"innodb_rollback_on_timeout,omitempty"`
	InternalTmpMemStorageEngine  *string  `json:"internal_tmp_mem_storage_engine,omitempty"`
	MaxHeapTableSize             *int     `json:"max_heap_table_size,omitempty"`
	TmpTableSize                 *int     `json:"tmp_table_size,omitempty"`
	SlowQueryLog                 *bool    `json:"slow_query_log,omitempty"`
	LongQueryTime                *float32 `json:"long_query_time,omitempty"`
	BackupHour                   *int     `json:"backup_hour,omitempty"`
	BackupMinute                 *int     `json:"backup_minute,omitempty"`
	BinlogRetentionPeriod        *int     `json:"binlog_retention_period,omitempty"`
	InnodbChangeBufferMaxSize    *int     `json:"innodb_change_buffer_max_size,omitempty"`
	InnodbFlushNeighbors         *int     `json:"innodb_flush_neighbors,omitempty"`
	InnodbReadIoThreads          *int     `json:"innodb_read_io_threads,omitempty"`
	InnodbThreadConcurrency      *int     `json:"innodb_thread_concurrency,omitempty"`
	InnodbWriteIoThreads         *int     `json:"innodb_write_io_threads,omitempty"`
	NetBufferLength              *int     `json:"net_buffer_length,omitempty"`
	LogOutput                    *string  `json:"log_output,omitempty"`
}

// MongoDBConfig holds advanced configurations for MongoDB database clusters.
type MongoDBConfig struct {
	DefaultReadConcern              *string `json:"default_read_concern,omitempty"`
	DefaultWriteConcern             *string `json:"default_write_concern,omitempty"`
	TransactionLifetimeLimitSeconds *int    `json:"transaction_lifetime_limit_seconds,omitempty"`
	SlowOpThresholdMs               *int    `json:"slow_op_threshold_ms,omitempty"`
	Verbosity                       *int    `json:"verbosity,omitempty"`
}

// KafkaConfig holds advanced configurations for Kafka database clusters.
type KafkaConfig struct {
	GroupInitialRebalanceDelayMs       *int     `json:"group_initial_rebalance_delay_ms,omitempty"`
	GroupMinSessionTimeoutMs           *int     `json:"group_min_session_timeout_ms,omitempty"`
	GroupMaxSessionTimeoutMs           *int     `json:"group_max_session_timeout_ms,omitempty"`
	MessageMaxBytes                    *int     `json:"message_max_bytes,omitempty"`
	LogCleanerDeleteRetentionMs        *int64   `json:"log_cleaner_delete_retention_ms,omitempty"`
	LogCleanerMinCompactionLagMs       *uint64  `json:"log_cleaner_min_compaction_lag_ms,omitempty"`
	LogFlushIntervalMs                 *uint64  `json:"log_flush_interval_ms,omitempty"`
	LogIndexIntervalBytes              *int     `json:"log_index_interval_bytes,omitempty"`
	LogMessageDownconversionEnable     *bool    `json:"log_message_downconversion_enable,omitempty"`
	LogMessageTimestampDifferenceMaxMs *uint64  `json:"log_message_timestamp_difference_max_ms,omitempty"`
	LogPreallocate                     *bool    `json:"log_preallocate,omitempty"`
	LogRetentionBytes                  *big.Int `json:"log_retention_bytes,omitempty"`
	LogRetentionHours                  *int     `json:"log_retention_hours,omitempty"`
	LogRetentionMs                     *big.Int `json:"log_retention_ms,omitempty"`
	LogRollJitterMs                    *uint64  `json:"log_roll_jitter_ms,omitempty"`
	LogSegmentDeleteDelayMs            *int     `json:"log_segment_delete_delay_ms,omitempty"`
	AutoCreateTopicsEnable             *bool    `json:"auto_create_topics_enable,omitempty"`
}

// OpensearchConfig holds advanced configurations for Opensearch database clusters.
type OpensearchConfig struct {
	HttpMaxContentLengthBytes                        *int     `json:"http_max_content_length_bytes,omitempty"`
	HttpMaxHeaderSizeBytes                           *int     `json:"http_max_header_size_bytes,omitempty"`
	HttpMaxInitialLineLengthBytes                    *int     `json:"http_max_initial_line_length_bytes,omitempty"`
	IndicesQueryBoolMaxClauseCount                   *int     `json:"indices_query_bool_max_clause_count,omitempty"`
	IndicesFielddataCacheSizePercentage              *int     `json:"indices_fielddata_cache_size_percentage,omitempty"`
	IndicesMemoryIndexBufferSizePercentage           *int     `json:"indices_memory_index_buffer_size_percentage,omitempty"`
	IndicesMemoryMinIndexBufferSizeMb                *int     `json:"indices_memory_min_index_buffer_size_mb,omitempty"`
	IndicesMemoryMaxIndexBufferSizeMb                *int     `json:"indices_memory_max_index_buffer_size_mb,omitempty"`
	IndicesQueriesCacheSizePercentage                *int     `json:"indices_queries_cache_size_percentage,omitempty"`
	IndicesRecoveryMaxMbPerSec                       *int     `json:"indices_recovery_max_mb_per_sec,omitempty"`
	IndicesRecoveryMaxConcurrentFileChunks           *int     `json:"indices_recovery_max_concurrent_file_chunks,omitempty"`
	ThreadPoolSearchSize                             *int     `json:"thread_pool_search_size,omitempty"`
	ThreadPoolSearchThrottledSize                    *int     `json:"thread_pool_search_throttled_size,omitempty"`
	ThreadPoolGetSize                                *int     `json:"thread_pool_get_size,omitempty"`
	ThreadPoolAnalyzeSize                            *int     `json:"thread_pool_analyze_size,omitempty"`
	ThreadPoolWriteSize                              *int     `json:"thread_pool_write_size,omitempty"`
	ThreadPoolForceMergeSize                         *int     `json:"thread_pool_force_merge_size,omitempty"`
	ThreadPoolSearchQueueSize                        *int     `json:"thread_pool_search_queue_size,omitempty"`
	ThreadPoolSearchThrottledQueueSize               *int     `json:"thread_pool_search_throttled_queue_size,omitempty"`
	ThreadPoolGetQueueSize                           *int     `json:"thread_pool_get_queue_size,omitempty"`
	ThreadPoolAnalyzeQueueSize                       *int     `json:"thread_pool_analyze_queue_size,omitempty"`
	ThreadPoolWriteQueueSize                         *int     `json:"thread_pool_write_queue_size,omitempty"`
	IsmEnabled                                       *bool    `json:"ism_enabled,omitempty"`
	IsmHistoryEnabled                                *bool    `json:"ism_history_enabled,omitempty"`
	IsmHistoryMaxAgeHours                            *int     `json:"ism_history_max_age_hours,omitempty"`
	IsmHistoryMaxDocs                                *int64   `json:"ism_history_max_docs,omitempty"`
	IsmHistoryRolloverCheckPeriodHours               *int     `json:"ism_history_rollover_check_period_hours,omitempty"`
	IsmHistoryRolloverRetentionPeriodDays            *int     `json:"ism_history_rollover_retention_period_days,omitempty"`
	SearchMaxBuckets                                 *int     `json:"search_max_buckets,omitempty"`
	ActionAutoCreateIndexEnabled                     *bool    `json:"action_auto_create_index_enabled,omitempty"`
	EnableSecurityAudit                              *bool    `json:"enable_security_audit,omitempty"`
	ActionDestructiveRequiresName                    *bool    `json:"action_destructive_requires_name,omitempty"`
	ClusterMaxShardsPerNode                          *int     `json:"cluster_max_shards_per_node,omitempty"`
	OverrideMainResponseVersion                      *bool    `json:"override_main_response_version,omitempty"`
	ScriptMaxCompilationsRate                        *string  `json:"script_max_compilations_rate,omitempty"`
	ClusterRoutingAllocationNodeConcurrentRecoveries *int     `json:"cluster_routing_allocation_node_concurrent_recoveries,omitempty"`
	ReindexRemoteWhitelist                           []string `json:"reindex_remote_whitelist,omitempty"`
	PluginsAlertingFilterByBackendRolesEnabled       *bool    `json:"plugins_alerting_filter_by_backend_roles_enabled,omitempty"`
}

type databaseUserRoot struct {
	User *DatabaseUser `json:"user"`
}

type databaseUsersRoot struct {
	Users []DatabaseUser `json:"users"`
}

type databaseDBRoot struct {
	DB *DatabaseDB `json:"db"`
}

type databaseDBsRoot struct {
	DBs []DatabaseDB `json:"dbs"`
}

type databasesRoot struct {
	Databases []Database `json:"databases"`
}

type databaseRoot struct {
	Database *Database `json:"database"`
}

type databaseCARoot struct {
	CA *DatabaseCA `json:"ca"`
}

type databasePostgreSQLConfigRoot struct {
	Config *PostgreSQLConfig `json:"config"`
}

type databaseRedisConfigRoot struct {
	Config *RedisConfig `json:"config"`
}

type databaseMySQLConfigRoot struct {
	Config *MySQLConfig `json:"config"`
}

type databaseMongoDBConfigRoot struct {
	Config *MongoDBConfig `json:"config"`
}

type databaseOpensearchConfigRoot struct {
	Config *OpensearchConfig `json:"config"`
}

type databaseKafkaConfigRoot struct {
	Config *KafkaConfig `json:"config"`
}

type databaseBackupsRoot struct {
	Backups []DatabaseBackup `json:"backups"`
}

type databasePoolRoot struct {
	Pool *DatabasePool `json:"pool"`
}

type databasePoolsRoot struct {
	Pools []DatabasePool `json:"pools"`
}

type databaseReplicaRoot struct {
	Replica *DatabaseReplica `json:"replica"`
}

type databaseReplicasRoot struct {
	Replicas []DatabaseReplica `json:"replicas"`
}

type evictionPolicyRoot struct {
	EvictionPolicy string `json:"eviction_policy"`
}

type UpgradeVersionRequest struct {
	Version string `json:"version"`
}

type sqlModeRoot struct {
	SQLMode string `json:"sql_mode"`
}

type databaseFirewallRuleRoot struct {
	Rules []DatabaseFirewallRule `json:"rules"`
}

// databaseOptionsRoot represents the root of all available database options (i.e. engines, regions, version, etc.)
type databaseOptionsRoot struct {
	Options *DatabaseOptions `json:"options"`
}

type databaseTopicRoot struct {
	Topic *DatabaseTopic `json:"topic"`
}

type databaseTopicsRoot struct {
	Topics []DatabaseTopic `json:"topics"`
}

type databaseLogsinksRoot struct {
	Sinks []DatabaseLogsink `json:"sinks"`
}

type databaseMetricsCredentialsRoot struct {
	Credentials *DatabaseMetricsCredentials `json:"credentials"`
}

type DatabaseMetricsCredentials struct {
	BasicAuthUsername string `json:"basic_auth_username"`
	BasicAuthPassword string `json:"basic_auth_password"`
}

type DatabaseUpdateMetricsCredentialsRequest struct {
	Credentials *DatabaseMetricsCredentials `json:"credentials"`
}

// DatabaseOptions represents the available database engines
type DatabaseOptions struct {
	MongoDBOptions     DatabaseEngineOptions `json:"mongodb"`
	MySQLOptions       DatabaseEngineOptions `json:"mysql"`
	PostgresSQLOptions DatabaseEngineOptions `json:"pg"`
	RedisOptions       DatabaseEngineOptions `json:"redis"`
	KafkaOptions       DatabaseEngineOptions `json:"kafka"`
	OpensearchOptions  DatabaseEngineOptions `json:"opensearch"`
}

// DatabaseEngineOptions represents the configuration options that are available for a given database engine
type DatabaseEngineOptions struct {
	Regions  []string         `json:"regions"`
	Versions []string         `json:"versions"`
	Layouts  []DatabaseLayout `json:"layouts"`
}

// DatabaseLayout represents the slugs available for a given database engine at various node counts
type DatabaseLayout struct {
	NodeNum int      `json:"num_nodes"`
	Sizes   []string `json:"sizes"`
}

// ListDatabaseEvents contains a list of project events.
type ListDatabaseEvents struct {
	Events []DatabaseEvent `json:"events"`
}

// DatbaseEvent contains the information about a Datbase event.
type DatabaseEvent struct {
	ID          string `json:"id"`
	ServiceName string `json:"cluster_name"`
	EventType   string `json:"event_type"`
	CreateTime  string `json:"create_time"`
}

type ListDatabaseEventsRoot struct {
	Events []DatabaseEvent `json:"events"`
}

type DatabaseIndex struct {
	IndexName        string            `json:"index_name"`
	NumberofShards   uint64            `json:"number_of_shards"`
	NumberofReplicas uint64            `json:"number_of_replicas"`
	Size             int64             `json:"size,omitempty"`
	Health           string            `json:"health,omitempty"`
	Status           string            `json:"status,omitempty"`
	Docs             int64             `json:"docs,omitempty"`
	CreateTime       string            `json:"create_time"`
	Replication      *IndexReplication `json:"replication,omitempty"`
}

type IndexReplication struct {
	LeaderIndex   string `json:"leader_index,omitempty"`
	LeaderProject string `json:"leader_project,omitempty"`
	LeaderService string `json:"leader_service,omitempty"`
}

type databaseIndexesRoot struct {
	Indexes []DatabaseIndex `json:"indexes"`
}

// URN returns a URN identifier for the database
func (d Database) URN() string {
	return ToURN("dbaas", d.ID)
}

// List returns a list of the Databases visible with the caller's API token
func (svc *DatabasesServiceOp) List(ctx context.Context, opts *ListOptions) ([]Database, *Response, error) {
	path := databaseBasePath
	path, err := addOptions(path, opts)
	if err != nil {
		return nil, nil, err
	}
	req, err := svc.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}
	root := new(databasesRoot)
	resp, err := svc.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	return root.Databases, resp, nil
}

// Get retrieves the details of a database cluster
func (svc *DatabasesServiceOp) Get(ctx context.Context, databaseID string) (*Database, *Response, error) {
	path := fmt.Sprintf(databaseSinglePath, databaseID)
	req, err := svc.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}
	root := new(databaseRoot)
	resp, err := svc.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	return root.Database, resp, nil
}

// GetCA retrieves the CA of a database cluster.
func (svc *DatabasesServiceOp) GetCA(ctx context.Context, databaseID string) (*DatabaseCA, *Response, error) {
	path := fmt.Sprintf(databaseCAPath, databaseID)
	req, err := svc.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}
	root := new(databaseCARoot)
	resp, err := svc.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	return root.CA, resp, nil
}

// Create creates a database cluster
func (svc *DatabasesServiceOp) Create(ctx context.Context, create *DatabaseCreateRequest) (*Database, *Response, error) {
	path := databaseBasePath
	req, err := svc.client.NewRequest(ctx, http.MethodPost, path, create)
	if err != nil {
		return nil, nil, err
	}
	root := new(databaseRoot)
	resp, err := svc.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	return root.Database, resp, nil
}

// Delete deletes a database cluster. There is no way to recover a cluster once
// it has been destroyed.
func (svc *DatabasesServiceOp) Delete(ctx context.Context, databaseID string) (*Response, error) {
	path := fmt.Sprintf("%s/%s", databaseBasePath, databaseID)
	req, err := svc.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}
	resp, err := svc.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

// Resize resizes a database cluster by number of nodes or size
func (svc *DatabasesServiceOp) Resize(ctx context.Context, databaseID string, resize *DatabaseResizeRequest) (*Response, error) {
	path := fmt.Sprintf(databaseResizePath, databaseID)
	req, err := svc.client.NewRequest(ctx, http.MethodPut, path, resize)
	if err != nil {
		return nil, err
	}
	resp, err := svc.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

// Migrate migrates a database cluster to a new region
func (svc *DatabasesServiceOp) Migrate(ctx context.Context, databaseID string, migrate *DatabaseMigrateRequest) (*Response, error) {
	path := fmt.Sprintf(databaseMigratePath, databaseID)
	req, err := svc.client.NewRequest(ctx, http.MethodPut, path, migrate)
	if err != nil {
		return nil, err
	}
	resp, err := svc.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

// UpdateMaintenance updates the maintenance window on a cluster
func (svc *DatabasesServiceOp) UpdateMaintenance(ctx context.Context, databaseID string, maintenance *DatabaseUpdateMaintenanceRequest) (*Response, error) {
	path := fmt.Sprintf(databaseMaintenancePath, databaseID)
	req, err := svc.client.NewRequest(ctx, http.MethodPut, path, maintenance)
	if err != nil {
		return nil, err
	}
	resp, err := svc.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

// InstallUpdate starts installation of updates
func (svc *DatabasesServiceOp) InstallUpdate(ctx context.Context, databaseID string) (*Response, error) {
	path := fmt.Sprintf(databaseUpdateInstallationPath, databaseID)
	req, err := svc.client.NewRequest(ctx, http.MethodPut, path, nil)
	if err != nil {
		return nil, err
	}
	resp, err := svc.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

// ListBackups returns a list of the current backups of a database
func (svc *DatabasesServiceOp) ListBackups(ctx context.Context, databaseID string, opts *ListOptions) ([]DatabaseBackup, *Response, error) {
	path := fmt.Sprintf(databaseBackupsPath, databaseID)
	path, err := addOptions(path, opts)
	if err != nil {
		return nil, nil, err
	}
	req, err := svc.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}
	root := new(databaseBackupsRoot)
	resp, err := svc.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	return root.Backups, resp, nil
}

// GetUser returns the database user identified by userID
func (svc *DatabasesServiceOp) GetUser(ctx context.Context, databaseID, userID string) (*DatabaseUser, *Response, error) {
	path := fmt.Sprintf(databaseUserPath, databaseID, userID)
	req, err := svc.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}
	root := new(databaseUserRoot)
	resp, err := svc.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	return root.User, resp, nil
}

// ListUsers returns all database users for the database
func (svc *DatabasesServiceOp) ListUsers(ctx context.Context, databaseID string, opts *ListOptions) ([]DatabaseUser, *Response, error) {
	path := fmt.Sprintf(databaseUsersPath, databaseID)
	path, err := addOptions(path, opts)
	if err != nil {
		return nil, nil, err
	}
	req, err := svc.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}
	root := new(databaseUsersRoot)
	resp, err := svc.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	return root.Users, resp, nil
}

// CreateUser will create a new database user
func (svc *DatabasesServiceOp) CreateUser(ctx context.Context, databaseID string, createUser *DatabaseCreateUserRequest) (*DatabaseUser, *Response, error) {
	path := fmt.Sprintf(databaseUsersPath, databaseID)
	req, err := svc.client.NewRequest(ctx, http.MethodPost, path, createUser)
	if err != nil {
		return nil, nil, err
	}
	root := new(databaseUserRoot)
	resp, err := svc.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	return root.User, resp, nil
}

// UpdateUser will update an existing database user
func (svc *DatabasesServiceOp) UpdateUser(ctx context.Context, databaseID, userID string, updateUser *DatabaseUpdateUserRequest) (*DatabaseUser, *Response, error) {
	path := fmt.Sprintf(databaseUserPath, databaseID, userID)
	req, err := svc.client.NewRequest(ctx, http.MethodPut, path, updateUser)
	if err != nil {
		return nil, nil, err
	}
	root := new(databaseUserRoot)
	resp, err := svc.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	return root.User, resp, nil
}

// ResetUserAuth will reset user authentication
func (svc *DatabasesServiceOp) ResetUserAuth(ctx context.Context, databaseID, userID string, resetAuth *DatabaseResetUserAuthRequest) (*DatabaseUser, *Response, error) {
	path := fmt.Sprintf(databaseResetUserAuthPath, databaseID, userID)
	req, err := svc.client.NewRequest(ctx, http.MethodPost, path, resetAuth)
	if err != nil {
		return nil, nil, err
	}
	root := new(databaseUserRoot)
	resp, err := svc.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	return root.User, resp, nil
}

// DeleteUser will delete an existing database user
func (svc *DatabasesServiceOp) DeleteUser(ctx context.Context, databaseID, userID string) (*Response, error) {
	path := fmt.Sprintf(databaseUserPath, databaseID, userID)
	req, err := svc.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}
	resp, err := svc.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

// ListDBs returns all databases for a given database cluster
func (svc *DatabasesServiceOp) ListDBs(ctx context.Context, databaseID string, opts *ListOptions) ([]DatabaseDB, *Response, error) {
	path := fmt.Sprintf(databaseDBsPath, databaseID)
	path, err := addOptions(path, opts)
	if err != nil {
		return nil, nil, err
	}
	req, err := svc.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}
	root := new(databaseDBsRoot)
	resp, err := svc.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	return root.DBs, resp, nil
}

// GetDB returns a single database by name
func (svc *DatabasesServiceOp) GetDB(ctx context.Context, databaseID, name string) (*DatabaseDB, *Response, error) {
	path := fmt.Sprintf(databaseDBPath, databaseID, name)
	req, err := svc.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}
	root := new(databaseDBRoot)
	resp, err := svc.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	return root.DB, resp, nil
}

// CreateDB will create a new database
func (svc *DatabasesServiceOp) CreateDB(ctx context.Context, databaseID string, createDB *DatabaseCreateDBRequest) (*DatabaseDB, *Response, error) {
	path := fmt.Sprintf(databaseDBsPath, databaseID)
	req, err := svc.client.NewRequest(ctx, http.MethodPost, path, createDB)
	if err != nil {
		return nil, nil, err
	}
	root := new(databaseDBRoot)
	resp, err := svc.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	return root.DB, resp, nil
}

// DeleteDB will delete an existing database
func (svc *DatabasesServiceOp) DeleteDB(ctx context.Context, databaseID, name string) (*Response, error) {
	path := fmt.Sprintf(databaseDBPath, databaseID, name)
	req, err := svc.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}
	resp, err := svc.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

// ListPools returns all connection pools for a given database cluster
func (svc *DatabasesServiceOp) ListPools(ctx context.Context, databaseID string, opts *ListOptions) ([]DatabasePool, *Response, error) {
	path := fmt.Sprintf(databasePoolsPath, databaseID)
	path, err := addOptions(path, opts)
	if err != nil {
		return nil, nil, err
	}
	req, err := svc.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}
	root := new(databasePoolsRoot)
	resp, err := svc.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	return root.Pools, resp, nil
}

// GetPool returns a single database connection pool by name
func (svc *DatabasesServiceOp) GetPool(ctx context.Context, databaseID, name string) (*DatabasePool, *Response, error) {
	path := fmt.Sprintf(databasePoolPath, databaseID, name)
	req, err := svc.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}
	root := new(databasePoolRoot)
	resp, err := svc.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	return root.Pool, resp, nil
}

// CreatePool will create a new database connection pool
func (svc *DatabasesServiceOp) CreatePool(ctx context.Context, databaseID string, createPool *DatabaseCreatePoolRequest) (*DatabasePool, *Response, error) {
	path := fmt.Sprintf(databasePoolsPath, databaseID)
	req, err := svc.client.NewRequest(ctx, http.MethodPost, path, createPool)
	if err != nil {
		return nil, nil, err
	}
	root := new(databasePoolRoot)
	resp, err := svc.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	return root.Pool, resp, nil
}

// DeletePool will delete an existing database connection pool
func (svc *DatabasesServiceOp) DeletePool(ctx context.Context, databaseID, name string) (*Response, error) {
	path := fmt.Sprintf(databasePoolPath, databaseID, name)
	req, err := svc.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}
	resp, err := svc.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

// UpdatePool will update an existing database connection pool
func (svc *DatabasesServiceOp) UpdatePool(ctx context.Context, databaseID, name string, updatePool *DatabaseUpdatePoolRequest) (*Response, error) {
	path := fmt.Sprintf(databasePoolPath, databaseID, name)

	if updatePool == nil {
		return nil, NewArgError("updatePool", "cannot be nil")
	}

	if updatePool.Mode == "" {
		return nil, NewArgError("mode", "cannot be empty")
	}

	if updatePool.Database == "" {
		return nil, NewArgError("database", "cannot be empty")
	}

	if updatePool.Size < 1 {
		return nil, NewArgError("size", "cannot be less than 1")
	}

	req, err := svc.client.NewRequest(ctx, http.MethodPut, path, updatePool)
	if err != nil {
		return nil, err
	}
	resp, err := svc.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

// GetReplica returns a single database replica
func (svc *DatabasesServiceOp) GetReplica(ctx context.Context, databaseID, name string) (*DatabaseReplica, *Response, error) {
	path := fmt.Sprintf(databaseReplicaPath, databaseID, name)
	req, err := svc.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}
	root := new(databaseReplicaRoot)
	resp, err := svc.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	return root.Replica, resp, nil
}

// ListReplicas returns all read-only replicas for a given database cluster
func (svc *DatabasesServiceOp) ListReplicas(ctx context.Context, databaseID string, opts *ListOptions) ([]DatabaseReplica, *Response, error) {
	path := fmt.Sprintf(databaseReplicasPath, databaseID)
	path, err := addOptions(path, opts)
	if err != nil {
		return nil, nil, err
	}
	req, err := svc.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}
	root := new(databaseReplicasRoot)
	resp, err := svc.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	return root.Replicas, resp, nil
}

// CreateReplica will create a new database connection pool
func (svc *DatabasesServiceOp) CreateReplica(ctx context.Context, databaseID string, createReplica *DatabaseCreateReplicaRequest) (*DatabaseReplica, *Response, error) {
	path := fmt.Sprintf(databaseReplicasPath, databaseID)
	req, err := svc.client.NewRequest(ctx, http.MethodPost, path, createReplica)
	if err != nil {
		return nil, nil, err
	}
	root := new(databaseReplicaRoot)
	resp, err := svc.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	return root.Replica, resp, nil
}

// DeleteReplica will delete an existing database replica
func (svc *DatabasesServiceOp) DeleteReplica(ctx context.Context, databaseID, name string) (*Response, error) {
	path := fmt.Sprintf(databaseReplicaPath, databaseID, name)
	req, err := svc.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}
	resp, err := svc.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

// PromoteReplicaToPrimary will sever the read replica integration and then promote the replica cluster to be a R/W cluster
func (svc *DatabasesServiceOp) PromoteReplicaToPrimary(ctx context.Context, databaseID, name string) (*Response, error) {
	path := fmt.Sprintf(databasePromoteReplicaToPrimaryPath, databaseID, name)
	req, err := svc.client.NewRequest(ctx, http.MethodPut, path, nil)
	if err != nil {
		return nil, err
	}
	resp, err := svc.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

// GetEvictionPolicy loads the eviction policy for a given Redis cluster.
func (svc *DatabasesServiceOp) GetEvictionPolicy(ctx context.Context, databaseID string) (string, *Response, error) {
	path := fmt.Sprintf(databaseEvictionPolicyPath, databaseID)
	req, err := svc.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return "", nil, err
	}
	root := new(evictionPolicyRoot)
	resp, err := svc.client.Do(ctx, req, root)
	if err != nil {
		return "", resp, err
	}
	return root.EvictionPolicy, resp, nil
}

// SetEvictionPolicy updates the eviction policy for a given Redis cluster.
//
// The valid eviction policies are documented by the exported string constants
// with the prefix `EvictionPolicy`.
func (svc *DatabasesServiceOp) SetEvictionPolicy(ctx context.Context, databaseID, policy string) (*Response, error) {
	path := fmt.Sprintf(databaseEvictionPolicyPath, databaseID)
	root := &evictionPolicyRoot{EvictionPolicy: policy}
	req, err := svc.client.NewRequest(ctx, http.MethodPut, path, root)
	if err != nil {
		return nil, err
	}
	resp, err := svc.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

// GetSQLMode loads the SQL Mode settings for a given MySQL cluster.
func (svc *DatabasesServiceOp) GetSQLMode(ctx context.Context, databaseID string) (string, *Response, error) {
	path := fmt.Sprintf(databaseSQLModePath, databaseID)
	req, err := svc.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return "", nil, err
	}
	root := &sqlModeRoot{}
	resp, err := svc.client.Do(ctx, req, root)
	if err != nil {
		return "", resp, err
	}
	return root.SQLMode, resp, nil
}

// SetSQLMode updates the SQL Mode settings for a given MySQL cluster.
func (svc *DatabasesServiceOp) SetSQLMode(ctx context.Context, databaseID string, sqlModes ...string) (*Response, error) {
	path := fmt.Sprintf(databaseSQLModePath, databaseID)
	root := &sqlModeRoot{SQLMode: strings.Join(sqlModes, ",")}
	req, err := svc.client.NewRequest(ctx, http.MethodPut, path, root)
	if err != nil {
		return nil, err
	}
	resp, err := svc.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

// GetFirewallRules loads the inbound sources for a given cluster.
func (svc *DatabasesServiceOp) GetFirewallRules(ctx context.Context, databaseID string) ([]DatabaseFirewallRule, *Response, error) {
	path := fmt.Sprintf(databaseFirewallRulesPath, databaseID)
	root := new(databaseFirewallRuleRoot)
	req, err := svc.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	resp, err := svc.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Rules, resp, nil
}

// UpdateFirewallRules sets the inbound sources for a given cluster.
func (svc *DatabasesServiceOp) UpdateFirewallRules(ctx context.Context, databaseID string, firewallRulesReq *DatabaseUpdateFirewallRulesRequest) (*Response, error) {
	path := fmt.Sprintf(databaseFirewallRulesPath, databaseID)
	req, err := svc.client.NewRequest(ctx, http.MethodPut, path, firewallRulesReq)
	if err != nil {
		return nil, err
	}
	return svc.client.Do(ctx, req, nil)
}

// GetPostgreSQLConfig retrieves the config for a PostgreSQL database cluster.
func (svc *DatabasesServiceOp) GetPostgreSQLConfig(ctx context.Context, databaseID string) (*PostgreSQLConfig, *Response, error) {
	path := fmt.Sprintf(databaseConfigPath, databaseID)
	req, err := svc.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}
	root := new(databasePostgreSQLConfigRoot)
	resp, err := svc.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	return root.Config, resp, nil
}

// UpdatePostgreSQLConfig updates the config for a PostgreSQL database cluster.
func (svc *DatabasesServiceOp) UpdatePostgreSQLConfig(ctx context.Context, databaseID string, config *PostgreSQLConfig) (*Response, error) {
	path := fmt.Sprintf(databaseConfigPath, databaseID)
	root := &databasePostgreSQLConfigRoot{
		Config: config,
	}
	req, err := svc.client.NewRequest(ctx, http.MethodPatch, path, root)
	if err != nil {
		return nil, err
	}
	resp, err := svc.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

// GetRedisConfig retrieves the config for a Redis database cluster.
func (svc *DatabasesServiceOp) GetRedisConfig(ctx context.Context, databaseID string) (*RedisConfig, *Response, error) {
	path := fmt.Sprintf(databaseConfigPath, databaseID)
	req, err := svc.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}
	root := new(databaseRedisConfigRoot)
	resp, err := svc.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	return root.Config, resp, nil
}

// UpdateRedisConfig updates the config for a Redis database cluster.
func (svc *DatabasesServiceOp) UpdateRedisConfig(ctx context.Context, databaseID string, config *RedisConfig) (*Response, error) {
	path := fmt.Sprintf(databaseConfigPath, databaseID)

	// We provide consts for use with SetEvictionPolicy method. Unfortunately, those are
	// in a different format than what can be used for RedisConfig.RedisMaxmemoryPolicy.
	// So we attempt to normalize them here to use dashes as separators if provided in
	// the old format (underscores). Other values are passed through untouched.
	if config.RedisMaxmemoryPolicy != nil {
		if policy, ok := evictionPolicyMap[*config.RedisMaxmemoryPolicy]; ok {
			config.RedisMaxmemoryPolicy = &policy
		}
	}

	root := &databaseRedisConfigRoot{
		Config: config,
	}
	req, err := svc.client.NewRequest(ctx, http.MethodPatch, path, root)
	if err != nil {
		return nil, err
	}
	resp, err := svc.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

// GetMySQLConfig retrieves the config for a MySQL database cluster.
func (svc *DatabasesServiceOp) GetMySQLConfig(ctx context.Context, databaseID string) (*MySQLConfig, *Response, error) {
	path := fmt.Sprintf(databaseConfigPath, databaseID)
	req, err := svc.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}
	root := new(databaseMySQLConfigRoot)
	resp, err := svc.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	return root.Config, resp, nil
}

// UpdateMySQLConfig updates the config for a MySQL database cluster.
func (svc *DatabasesServiceOp) UpdateMySQLConfig(ctx context.Context, databaseID string, config *MySQLConfig) (*Response, error) {
	path := fmt.Sprintf(databaseConfigPath, databaseID)
	root := &databaseMySQLConfigRoot{
		Config: config,
	}
	req, err := svc.client.NewRequest(ctx, http.MethodPatch, path, root)
	if err != nil {
		return nil, err
	}
	resp, err := svc.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

// GetMongoDBConfig retrieves the config for a MongoDB database cluster.
func (svc *DatabasesServiceOp) GetMongoDBConfig(ctx context.Context, databaseID string) (*MongoDBConfig, *Response, error) {
	path := fmt.Sprintf(databaseConfigPath, databaseID)
	req, err := svc.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}
	root := new(databaseMongoDBConfigRoot)
	resp, err := svc.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	return root.Config, resp, nil
}

// UpdateMongoDBConfig updates the config for a MongoDB database cluster.
func (svc *DatabasesServiceOp) UpdateMongoDBConfig(ctx context.Context, databaseID string, config *MongoDBConfig) (*Response, error) {
	path := fmt.Sprintf(databaseConfigPath, databaseID)
	root := &databaseMongoDBConfigRoot{
		Config: config,
	}
	req, err := svc.client.NewRequest(ctx, http.MethodPatch, path, root)
	if err != nil {
		return nil, err
	}
	resp, err := svc.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

// GetKafkaConfig retrieves the config for a Kafka database cluster.
func (svc *DatabasesServiceOp) GetKafkaConfig(ctx context.Context, databaseID string) (*KafkaConfig, *Response, error) {
	path := fmt.Sprintf(databaseConfigPath, databaseID)
	req, err := svc.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}
	root := new(databaseKafkaConfigRoot)
	resp, err := svc.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	return root.Config, resp, nil
}

// UpdateKafkaConfig updates the config for a Kafka database cluster.
func (svc *DatabasesServiceOp) UpdateKafkaConfig(ctx context.Context, databaseID string, config *KafkaConfig) (*Response, error) {
	path := fmt.Sprintf(databaseConfigPath, databaseID)
	root := &databaseKafkaConfigRoot{
		Config: config,
	}
	req, err := svc.client.NewRequest(ctx, http.MethodPatch, path, root)
	if err != nil {
		return nil, err
	}
	resp, err := svc.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

// GetOpensearchConfig retrieves the config for a Opensearch database cluster.
func (svc *DatabasesServiceOp) GetOpensearchConfig(ctx context.Context, databaseID string) (*OpensearchConfig, *Response, error) {
	path := fmt.Sprintf(databaseConfigPath, databaseID)
	req, err := svc.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}
	root := new(databaseOpensearchConfigRoot)
	resp, err := svc.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	return root.Config, resp, nil
}

// UpdateOpensearchConfig updates the config for a Opensearch database cluster.
func (svc *DatabasesServiceOp) UpdateOpensearchConfig(ctx context.Context, databaseID string, config *OpensearchConfig) (*Response, error) {
	path := fmt.Sprintf(databaseConfigPath, databaseID)
	root := &databaseOpensearchConfigRoot{
		Config: config,
	}
	req, err := svc.client.NewRequest(ctx, http.MethodPatch, path, root)
	if err != nil {
		return nil, err
	}
	resp, err := svc.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

// ListOptions gets the database options available.
func (svc *DatabasesServiceOp) ListOptions(ctx context.Context) (*DatabaseOptions, *Response, error) {
	root := new(databaseOptionsRoot)
	req, err := svc.client.NewRequest(ctx, http.MethodGet, databaseOptionsPath, nil)
	if err != nil {
		return nil, nil, err
	}

	resp, err := svc.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Options, resp, nil
}

// UpgradeMajorVersion upgrades the major version of a cluster.
func (svc *DatabasesServiceOp) UpgradeMajorVersion(ctx context.Context, databaseID string, upgradeReq *UpgradeVersionRequest) (*Response, error) {
	path := fmt.Sprintf(databaseUpgradeMajorVersionPath, databaseID)
	req, err := svc.client.NewRequest(ctx, http.MethodPut, path, upgradeReq)
	if err != nil {
		return nil, err
	}

	resp, err := svc.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// ListTopics returns all topics for a given kafka cluster
func (svc *DatabasesServiceOp) ListTopics(ctx context.Context, databaseID string, opts *ListOptions) ([]DatabaseTopic, *Response, error) {
	path := fmt.Sprintf(databaseTopicsPath, databaseID)
	path, err := addOptions(path, opts)
	if err != nil {
		return nil, nil, err
	}
	req, err := svc.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}
	root := new(databaseTopicsRoot)
	resp, err := svc.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	return root.Topics, resp, nil
}

// GetTopic returns a single kafka topic by name
func (svc *DatabasesServiceOp) GetTopic(ctx context.Context, databaseID, name string) (*DatabaseTopic, *Response, error) {
	path := fmt.Sprintf(databaseTopicPath, databaseID, name)
	req, err := svc.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}
	root := new(databaseTopicRoot)
	resp, err := svc.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	return root.Topic, resp, nil
}

// CreateTopic will create a new kafka topic
func (svc *DatabasesServiceOp) CreateTopic(ctx context.Context, databaseID string, createTopic *DatabaseCreateTopicRequest) (*DatabaseTopic, *Response, error) {
	path := fmt.Sprintf(databaseTopicsPath, databaseID)
	req, err := svc.client.NewRequest(ctx, http.MethodPost, path, createTopic)
	if err != nil {
		return nil, nil, err
	}
	root := new(databaseTopicRoot)
	resp, err := svc.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	return root.Topic, resp, nil
}

// UpdateTopic updates a single kafka topic
func (svc *DatabasesServiceOp) UpdateTopic(ctx context.Context, databaseID string, name string, updateTopic *DatabaseUpdateTopicRequest) (*Response, error) {
	path := fmt.Sprintf(databaseTopicPath, databaseID, name)
	req, err := svc.client.NewRequest(ctx, http.MethodPut, path, updateTopic)
	if err != nil {
		return nil, err
	}
	root := new(databaseTopicRoot)
	resp, err := svc.client.Do(ctx, req, root)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

// DeleteTopic will delete an existing kafka topic
func (svc *DatabasesServiceOp) DeleteTopic(ctx context.Context, databaseID, name string) (*Response, error) {
	path := fmt.Sprintf(databaseTopicPath, databaseID, name)
	req, err := svc.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}
	resp, err := svc.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

// GetMetricsCredentials gets the credentials required to access a user's metrics endpoints
func (svc *DatabasesServiceOp) GetMetricsCredentials(ctx context.Context) (*DatabaseMetricsCredentials, *Response, error) {
	req, err := svc.client.NewRequest(ctx, http.MethodGet, databaseMetricsCredentialsPath, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(databaseMetricsCredentialsRoot)
	resp, err := svc.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	return root.Credentials, resp, nil
}

// UpdateMetricsAuth updates the credentials required to access a user's metrics endpoints
func (svc *DatabasesServiceOp) UpdateMetricsCredentials(ctx context.Context, updateCreds *DatabaseUpdateMetricsCredentialsRequest) (*Response, error) {
	req, err := svc.client.NewRequest(ctx, http.MethodPut, databaseMetricsCredentialsPath, updateCreds)
	if err != nil {
		return nil, err
	}
	resp, err := svc.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

// ListDatabaseEvents returns all the events for a given cluster
func (svc *DatabasesServiceOp) ListDatabaseEvents(ctx context.Context, databaseID string, opts *ListOptions) ([]DatabaseEvent, *Response, error) {
	path := fmt.Sprintf(databaseEvents, databaseID)
	path, err := addOptions(path, opts)
	if err != nil {
		return nil, nil, err
	}
	root := new(ListDatabaseEventsRoot)
	req, err := svc.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	resp, err := svc.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Events, resp, nil
}

// ListIndexes returns all indexes for a given opensearch cluster
func (svc *DatabasesServiceOp) ListIndexes(ctx context.Context, databaseID string, opts *ListOptions) ([]DatabaseIndex, *Response, error) {
	path := fmt.Sprintf(databaseIndexesPath, databaseID)
	path, err := addOptions(path, opts)
	if err != nil {
		return nil, nil, err
	}
	req, err := svc.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}
	root := new(databaseIndexesRoot)
	resp, err := svc.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	return root.Indexes, resp, nil
}

// DeleteIndex will delete an existing opensearch index
func (svc *DatabasesServiceOp) DeleteIndex(ctx context.Context, databaseID, name string) (*Response, error) {
	path := fmt.Sprintf(databaseIndexPath, databaseID, name)
	req, err := svc.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}
	resp, err := svc.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

// CreateLogsink creates a new logsink for a database
func (svc *DatabasesServiceOp) CreateLogsink(ctx context.Context, databaseID string, createLogsink *DatabaseCreateLogsinkRequest) (*DatabaseLogsink, *Response, error) {
	path := fmt.Sprintf(databaseLogsinksPath, databaseID)
	req, err := svc.client.NewRequest(ctx, http.MethodPost, path, createLogsink)
	if err != nil {
		return nil, nil, err
	}

	root := new(DatabaseLogsink)
	resp, err := svc.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	return root, resp, nil
}

// GetLogsink gets a logsink for a database
func (svc *DatabasesServiceOp) GetLogsink(ctx context.Context, databaseID string, logsinkID string) (*DatabaseLogsink, *Response, error) {
	path := fmt.Sprintf(databaseLogsinkPath, databaseID, logsinkID)
	req, err := svc.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(DatabaseLogsink)
	resp, err := svc.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	return root, resp, nil
}

// ListTopics returns all topics for a given kafka cluster
func (svc *DatabasesServiceOp) ListLogsinks(ctx context.Context, databaseID string, opts *ListOptions) ([]DatabaseLogsink, *Response, error) {
	path := fmt.Sprintf(databaseLogsinksPath, databaseID)
	path, err := addOptions(path, opts)
	if err != nil {
		return nil, nil, err
	}
	req, err := svc.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}
	root := new(databaseLogsinksRoot)
	resp, err := svc.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	return root.Sinks, resp, nil
}

// UpdateLogsink updates a logsink for a database cluster
func (svc *DatabasesServiceOp) UpdateLogsink(ctx context.Context, databaseID string, logsinkID string, updateLogsink *DatabaseUpdateLogsinkRequest) (*Response, error) {
	path := fmt.Sprintf(databaseLogsinkPath, databaseID, logsinkID)
	req, err := svc.client.NewRequest(ctx, http.MethodPut, path, updateLogsink)
	if err != nil {
		return nil, err
	}

	resp, err := svc.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

// DeleteLogsink deletes a logsink for a database cluster
func (svc *DatabasesServiceOp) DeleteLogsink(ctx context.Context, databaseID, logsinkID string) (*Response, error) {
	path := fmt.Sprintf(databaseLogsinkPath, databaseID, logsinkID)
	req, err := svc.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}
	resp, err := svc.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}
	return resp, nil
}
