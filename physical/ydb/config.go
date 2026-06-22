package ydb

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	ydbconsts "github.com/hashicorp/vault/physical/ydb/consts"
	"github.com/ydb-platform/ydb-go-sdk/v3"
	"github.com/ydb-platform/ydb-go-sdk/v3/balancers"
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

	if authOpt := getYDBAuthOptionFromConfMap(conf); authOpt != nil {
		opts = append(opts, authOpt)
	}

	return opts, nil
}

func getYDBAuthOptionFromConfMap(conf map[string]string) ydb.Option {
	switch auth := resolveYDBAuth(conf); auth.kind {
	case "token":
		return ydb.WithAccessTokenCredentials(auth.value)
	case "static":
		return ydb.WithStaticCredentials(auth.value, auth.value2)
	case "anonymous":
		return ydb.WithAnonymousCredentials()
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
	if user, password := lookupStaticCredentials(conf); user != "" && password != "" {
		return ydbAuthConfig{kind: "static", value: user, value2: password}
	}
	if lookupFirstBool(ydbconsts.EnvAnonymousCredentials, conf["anonymous_credentials"]) {
		return ydbAuthConfig{kind: "anonymous"}
	}
	return ydbAuthConfig{}
}

func lookupStaticCredentials(conf map[string]string) (string, string) {
	user := lookupFirstNonEmpty(ydbconsts.EnvStaticCredentialsUser, conf["static_credentials_user"])
	password := lookupFirstNonEmpty(ydbconsts.EnvStaticCredentialsPassword, conf["static_credentials_password"])
	return user, password
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
