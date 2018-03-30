package dbutil

import (
	"reflect"
	"testing"

	"github.com/hashicorp/vault/builtin/logical/database/dbplugin"
)

func TestStatementCompatibilityHelper(t *testing.T) {
	const (
		creationStatement = "creation"
		renewStatement    = "renew"
		revokeStatement   = "revoke"
		rollbackStatement = "rollback"
	)

	expectedStatements := dbplugin.Statements{
		Creation:             []string{creationStatement},
		Rollback:             []string{rollbackStatement},
		Revocation:           []string{revokeStatement},
		Renewal:              []string{renewStatement},
		CreationStatements:   creationStatement,
		RenewStatements:      renewStatement,
		RollbackStatements:   rollbackStatement,
		RevocationStatements: revokeStatement,
	}

	statements1 := dbplugin.Statements{
		CreationStatements:   creationStatement,
		RenewStatements:      renewStatement,
		RollbackStatements:   rollbackStatement,
		RevocationStatements: revokeStatement,
	}

	if !reflect.DeepEqual(expectedStatements, StatementCompatibilityHelper(statements1)) {
		t.Fatalf("mismatch: %#v, %#v", expectedStatements, statements1)
	}

	statements2 := dbplugin.Statements{
		Creation:   []string{creationStatement},
		Rollback:   []string{rollbackStatement},
		Revocation: []string{revokeStatement},
		Renewal:    []string{renewStatement},
	}

	if !reflect.DeepEqual(expectedStatements, StatementCompatibilityHelper(statements2)) {
		t.Fatalf("mismatch: %#v, %#v", expectedStatements, statements2)
	}

	statements3 := dbplugin.Statements{
		CreationStatements: creationStatement,
	}
	expectedStatements3 := dbplugin.Statements{
		Creation:           []string{creationStatement},
		CreationStatements: creationStatement,
	}
	if !reflect.DeepEqual(expectedStatements3, StatementCompatibilityHelper(statements3)) {
		t.Fatalf("mismatch: %#v, %#v", expectedStatements3, statements3)
	}

}
