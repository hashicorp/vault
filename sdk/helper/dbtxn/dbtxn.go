package dbtxn

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

// ExecuteDBQuery handles executing one single statement, while properly releasing its resources.
// - ctx: 	Required
// - db: 	Required
// - config: 	Optional, may be nil
// - query: 	Required
func ExecuteDBQuery(ctx context.Context, db *sql.DB, params map[string]string, query string) error {

	parsedQuery := parseQuery(params, query)

	stmt, err := db.PrepareContext(ctx, parsedQuery)
	if err != nil {
		return err
	}
	defer stmt.Close()

	return execute(ctx, stmt)
}

// ExecuteTxQuery handles executing one single statement, while properly releasing its resources.
// - ctx: 	Required
// - tx: 	Required
// - config: 	Optional, may be nil
// - query: 	Required
func ExecuteTxQuery(ctx context.Context, tx *sql.Tx, params map[string]string, query string) error {

	parsedQuery := parseQuery(params, query)

	stmt, err := tx.PrepareContext(ctx, parsedQuery)
	if err != nil {
		return err
	}
	defer stmt.Close()

	return execute(ctx, stmt)
}

func execute(ctx context.Context, stmt *sql.Stmt) error {
	if _, err := stmt.ExecContext(ctx); err != nil {
		return err
	}
	return nil
}

func parseQuery(m map[string]string, tpl string) string {

	if m == nil || len(m) <= 0 {
		return tpl
	}

	for k, v := range m {
		tpl = strings.Replace(tpl, fmt.Sprintf("{{%s}}", k), v, -1)
	}
	return tpl
}
