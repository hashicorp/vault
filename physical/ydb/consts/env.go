package ydbconst

const (
	EnvDSN                       = "VAULT_YDB_DSN"
	EnvTable                     = "VAULT_YDB_TABLE"
	EnvToken                     = "VAULT_YDB_TOKEN"
	EnvInternalCA                = "VAULT_YDB_INTERNAL_CA"
	EnvSAKeyFile                 = "VAULT_YDB_SA_KEYFILE"
	EnvSAKey                     = "VAULT_YDB_SA_KEY"
	EnvStaticCredentialsUser     = "VAULT_YDB_STATIC_CREDENTIALS_USER"
	EnvStaticCredentialsPassword = "VAULT_YDB_STATIC_CREDENTIALS_PASSWORD"
	EnvMetadataAuth              = "VAULT_YDB_METADATA_AUTH"
	EnvAnonymousCredentials      = "VAULT_YDB_ANONYMOUS_CREDENTIALS"
	EnvHACoordinationNode        = "VAULT_YDB_HA_COORDINATION_NODE"
	EnvHAEnabled                 = "VAULT_YDB_HA_ENABLED"
	EnvTransactionMaxEntries     = "VAULT_YDB_TRANSACTION_MAX_ENTRIES"
	EnvTransactionMaxSize        = "VAULT_YDB_TRANSACTION_MAX_SIZE"
	EnvBalancer                  = "VAULT_YDB_BALANCER"

	VAULT_TABLE = "vault_kv"
)
