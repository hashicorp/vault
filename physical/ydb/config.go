package ydb

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	ydbconsts "github.com/hashicorp/vault/physical/ydb/consts"
	env "github.com/ydb-platform/ydb-go-sdk-auth-environ"
	"github.com/ydb-platform/ydb-go-sdk/v3"
	"github.com/ydb-platform/ydb-go-sdk/v3/balancers"
	yc "github.com/ydb-platform/ydb-go-yc"
)

const (
	defaultYDBTransactionMaxEntries = 63
	defaultYDBTransactionMaxSize    = 128 * 1024
)

func getYDBOptionsFromConfMap(conf map[string]string) ([]ydb.Option, error) {
	var opts []ydb.Option

	if balancerConfig := lookupFirstNonEmpty(ydbconsts.EnvBalancer, conf["balancer"]); balancerConfig != "" {
		balancer, err := balancers.CreateFromConfig(balancerConfig)
		if err != nil {
			return nil, fmt.Errorf("balancer: %w", err)
		}
		opts = append(opts, ydb.WithBalancer(balancer))
	}

	internalCAVal := ""
	if envv := os.Getenv(ydbconsts.EnvInternalCA); envv != "" {
		internalCAVal = envv
	} else if v, ok := conf["internal_ca"]; ok {
		internalCAVal = v
	}

	if parseYDBBool(internalCAVal) {
		opts = append(opts, yc.WithInternalCA())
	}
	if authOpt := getYDBAuthOptionFromConfMap(conf); authOpt != nil {
		opts = append(opts, authOpt)
	}

	return opts, nil
}

func getYDBAuthOptionFromConfMap(conf map[string]string) ydb.Option {
	switch auth := resolveYDBAuth(conf); auth.kind {
	case "token":
		return ydb.WithAccessTokenCredentials(auth.value)
	case "service_account_key_file":
		return yc.WithServiceAccountKeyFileCredentials(auth.value)
	case "service_account_key":
		return yc.WithServiceAccountKeyCredentials(auth.value)
	case "static":
		return ydb.WithStaticCredentials(auth.value, auth.value2)
	case "metadata":
		return yc.WithMetadataCredentials()
	case "anonymous":
		return ydb.WithAnonymousCredentials()
	case "environ":
		return env.WithEnvironCredentials()
	default:
		return nil
	}
}

type ydbAuthConfig struct {
	kind   string
	value  string
	value2 string
}

func resolveYDBAuth(conf map[string]string) ydbAuthConfig {
	if value := lookupFirstNonEmpty(ydbconsts.EnvToken, conf["token"]); value != "" {
		return ydbAuthConfig{kind: "token", value: value}
	}
	if value := lookupFirstNonEmpty(ydbconsts.EnvSAKeyFile, conf["service_account_key_file"]); value != "" {
		return ydbAuthConfig{kind: "service_account_key_file", value: value}
	}
	if value := lookupFirstNonEmpty(ydbconsts.EnvSAKey, conf["service_account_key"]); value != "" {
		return ydbAuthConfig{kind: "service_account_key", value: value}
	}
	if user, password := lookupStaticCredentials(conf); user != "" && password != "" {
		return ydbAuthConfig{kind: "static", value: user, value2: password}
	}
	if lookupFirstBool(ydbconsts.EnvMetadataAuth, conf["metadata_auth"]) {
		return ydbAuthConfig{kind: "metadata"}
	}
	if lookupFirstBool(ydbconsts.EnvAnonymousCredentials, conf["anonymous_credentials"]) {
		return ydbAuthConfig{kind: "anonymous"}
	}
	if hasYDBEnvironCredentials() {
		return ydbAuthConfig{kind: "environ"}
	}
	return ydbAuthConfig{}
}

func lookupStaticCredentials(conf map[string]string) (string, string) {
	user := lookupFirstNonEmpty(ydbconsts.EnvStaticCredentialsUser, conf["static_credentials_user"])
	password := lookupFirstNonEmpty(ydbconsts.EnvStaticCredentialsPassword, conf["static_credentials_password"])
	return user, password
}

func hasYDBEnvironCredentials() bool {
	if lookupEnvNonEmpty("YDB_SERVICE_ACCOUNT_KEY_CREDENTIALS") != "" {
		return true
	}
	if lookupEnvNonEmpty("YDB_SERVICE_ACCOUNT_KEY_FILE_CREDENTIALS") != "" {
		return true
	}
	if lookupEnvBool("YDB_METADATA_CREDENTIALS") {
		return true
	}
	if lookupEnvNonEmpty("YDB_ACCESS_TOKEN_CREDENTIALS") != "" {
		return true
	}
	if lookupEnvNonEmpty("YDB_STATIC_CREDENTIALS_USER") != "" &&
		lookupEnvNonEmpty("YDB_STATIC_CREDENTIALS_PASSWORD") != "" &&
		lookupEnvNonEmpty("YDB_STATIC_CREDENTIALS_ENDPOINT") != "" {
		return true
	}
	if lookupEnvNonEmpty("YDB_OAUTH2_KEY_FILE") != "" {
		return true
	}
	if value, ok := os.LookupEnv("YDB_ANONYMOUS_CREDENTIALS"); ok && strings.TrimSpace(value) == "0" {
		return true
	}
	return false
}

func lookupFirstNonEmpty(envKey, confValue string) string {
	if envv := strings.TrimSpace(os.Getenv(envKey)); envv != "" {
		return envv
	}
	return strings.TrimSpace(confValue)
}

func lookupFirstBool(envKey, confValue string) bool {
	if envv, ok := os.LookupEnv(envKey); ok {
		if strings.TrimSpace(envv) == "" {
			return parseYDBBool(confValue)
		}
		return parseYDBBool(envv)
	}
	return parseYDBBool(confValue)
}

func lookupEnvNonEmpty(envKey string) string {
	return strings.TrimSpace(os.Getenv(envKey))
}

func lookupEnvBool(envKey string) bool {
	envv, ok := os.LookupEnv(envKey)
	if !ok {
		return false
	}
	return parseYDBBool(envv)
}

func parseYDBBool(value string) bool {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "1", "true", "yes":
		return true
	default:
		return false
	}
}

func getYDBHACoordinationNodePath(conf map[string]string, dbName, table string) string {
	if path := lookupFirstNonEmpty(ydbconsts.EnvHACoordinationNode, conf["ha_coordination_node"]); path != "" {
		if strings.HasPrefix(path, "/") {
			return path
		}
		return dbName + "/" + path
	}

	if strings.HasPrefix(table, "/") {
		return table + "_ha"
	}
	return dbName + "/" + table + "_ha"
}

func getYDBHAEnabled(conf map[string]string) bool {
	return lookupFirstBool(ydbconsts.EnvHAEnabled, conf["ha_enabled"])
}

func getYDBTransactionLimits(conf map[string]string) (int, int, error) {
	maxEntries, err := lookupPositiveInt(
		ydbconsts.EnvTransactionMaxEntries,
		conf["transaction_max_entries"],
		defaultYDBTransactionMaxEntries,
	)
	if err != nil {
		return 0, 0, fmt.Errorf("transaction_max_entries: %w", err)
	}

	maxSize, err := lookupPositiveInt(
		ydbconsts.EnvTransactionMaxSize,
		conf["transaction_max_size"],
		defaultYDBTransactionMaxSize,
	)
	if err != nil {
		return 0, 0, fmt.Errorf("transaction_max_size: %w", err)
	}

	return maxEntries, maxSize, nil
}

func lookupPositiveInt(envKey, confValue string, defaultValue int) (int, error) {
	value := lookupFirstNonEmpty(envKey, confValue)
	if value == "" {
		return defaultValue, nil
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("must be an integer")
	}
	if parsed <= 0 {
		return 0, fmt.Errorf("must be greater than zero")
	}
	return parsed, nil
}
