package postgresql

import (
	"os"
)

// connectionURL first check the environment variables for a connection URL. If
// no connection URL exists in the environment variable, the Vault config file is
// checked. If neither the environment variables or the config file set the connection
// URL for the Postgres backend, because it is a required field, an error is returned.
func connectionURL(conf map[string]string) string {
	connURL := os.Getenv("PG_CONNECTION_URL")
	if connURL == "" {
		connURL = conf["connection_url"]
		if connURL == "" {
			return ""
		}
	}

	return connURL
}
