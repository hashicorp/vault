package dbtxn

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

type Config struct {
	Name       string
	Username   string
	Password   string
	Expiration string
}

// ExecuteDBQuery handles executing one single statement, while properly releasing its resources.
// - ctx: 	Optional, may be nil
// - db: 	Required
// - config: 	Optional, may be nil
// - query: 	Required
func ExecuteDBQuery(ctx context.Context, db *sql.DB, config *Config, query string) error {

	parsedQuery := parseQuery(config, query)

	var stmt *sql.Stmt
	var err error
	if ctx != nil {
		stmt, err = db.PrepareContext(ctx, parsedQuery)
	} else {
		stmt, err = db.Prepare(parsedQuery)
	}
	if err != nil {
		return err
	}
	defer stmt.Close()

	return execute(ctx, stmt)
}

// ExecuteTxQuery handles executing one single statement, while properly releasing its resources.
// - ctx: 	Optional, may be nil
// - tx: 	Required
// - config: 	Optional, may be nil
// - query: 	Required
func ExecuteTxQuery(ctx context.Context, tx *sql.Tx, config *Config, query string) error {

	parsedQuery := parseQuery(config, query)

	var stmt *sql.Stmt
	var err error
	if ctx != nil {
		stmt, err = tx.PrepareContext(ctx, parsedQuery)
	} else {
		stmt, err = tx.Prepare(parsedQuery)
	}
	if err != nil {
		return err
	}
	defer stmt.Close()

	return execute(ctx, stmt)
}

func execute(ctx context.Context, stmt *sql.Stmt) error {

	if ctx != nil {
		if _, err := stmt.ExecContext(ctx); err != nil {
			return err
		}
		return nil
	}

	if _, err := stmt.Exec(); err != nil {
		return err
	}
	return nil
}

func parseQuery(c *Config, tpl string) string {

	if c == nil {
		return tpl
	}

	if c.Name == "" && c.Username == "" && c.Password == "" && c.Expiration == "" {
		return tpl
	}

	data := make(map[string]string)
	if c.Name != "" {
		data["name"] = c.Name
	}
	if c.Username != "" {
		data["username"] = c.Username
	}
	if c.Password != "" {
		data["password"] = c.Password
	}
	if c.Expiration != "" {
		data["expiration"] = c.Expiration
	}

	for k, v := range data {
		tpl = strings.Replace(tpl, fmt.Sprintf("{{%s}}", k), v, -1)
	}

	return tpl
}
