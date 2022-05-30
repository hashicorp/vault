package dbtxn

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

// ExecuteDBQuery handles executing one single statement while properly releasing its resources.
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

// ExecuteDBQueryDirect handles executing one single statement without preparing the query
// before executing it, which can be more efficient.
// - ctx: 	Required
// - db: 	Required
// - config: 	Optional, may be nil
// - query: 	Required
func ExecuteDBQueryDirect(ctx context.Context, db *sql.DB, params map[string]string, query string) error {
	parsedQuery := parseQuery(params, query)
	_, err := db.ExecContext(ctx, parsedQuery)
	return err
}

// ExecuteTxQuery handles executing one single statement while properly releasing its resources.
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

// ExecuteTxQueryDirect handles executing one single statement.
// - ctx: 	Required
// - tx: 	Required
// - config: 	Optional, may be nil
// - query: 	Required
func ExecuteTxQueryDirect(ctx context.Context, tx *sql.Tx, params map[string]string, query string) error {
	parsedQuery := parseQuery(params, query)
	_, err := tx.ExecContext(ctx, parsedQuery)
	return err
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
