package dbutil

import (
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/vault/sdk/database/dbplugin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrEmptyCreationStatement = errors.New("empty creation statements")
	ErrEmptyRotationStatement = errors.New("empty rotation statements")
)

// Query templates a query for us.
func QueryHelper(tpl string, data map[string]string) string {
	for k, v := range data {
		tpl = strings.Replace(tpl, fmt.Sprintf("{{%s}}", k), v, -1)
	}

	return tpl
}

// StatementCompatibilityHelper will populate the statements fields to support
// compatibility
func StatementCompatibilityHelper(statements dbplugin.Statements) dbplugin.Statements {
	switch {
	case len(statements.Creation) > 0 && len(statements.CreationStatements) == 0:
		statements.CreationStatements = strings.Join(statements.Creation, ";")
	case len(statements.CreationStatements) > 0:
		statements.Creation = []string{statements.CreationStatements}
	}
	switch {
	case len(statements.Revocation) > 0 && len(statements.RevocationStatements) == 0:
		statements.RevocationStatements = strings.Join(statements.Revocation, ";")
	case len(statements.RevocationStatements) > 0:
		statements.Revocation = []string{statements.RevocationStatements}
	}
	switch {
	case len(statements.Renewal) > 0 && len(statements.RenewStatements) == 0:
		statements.RenewStatements = strings.Join(statements.Renewal, ";")
	case len(statements.RenewStatements) > 0:
		statements.Renewal = []string{statements.RenewStatements}
	}
	switch {
	case len(statements.Rollback) > 0 && len(statements.RollbackStatements) == 0:
		statements.RollbackStatements = strings.Join(statements.Rollback, ";")
	case len(statements.RollbackStatements) > 0:
		statements.Rollback = []string{statements.RollbackStatements}
	}
	return statements
}

// Unimplemented returns a gRPC error with the Unimplemented code
func Unimplemented() error {
	return status.Error(codes.Unimplemented, "Not yet implemented")
}
