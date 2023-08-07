// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package dbutil

import (
	"reflect"
	"testing"

	"github.com/hashicorp/vault/sdk/database/dbplugin"
	"github.com/stretchr/testify/assert"
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

func TestQueryHelper(t *testing.T) {
	data := map[string]string{
		// These are typical keys you find in a data map used with QueryHelper
		"username":   "hello",
		"name":       "hello",
		"password":   "world",
		"expiration": "24h",
	}
	for _, tc := range []struct {
		tpl      string
		expected string
	}{
		{"", ""},
		{"somedb://{{username}}:{{password}}@something", "somedb://hello:world@something"},
		// Unknown placeholders pass through as is
		{"user={{name}} other={{unknown}}", "user=hello other={{unknown}}"},
		// Various incorrect delimiters pass through as is
		{"{{{{{{{{", "{{{{{{{{"},
		{"{{username}} {{incomplete", "hello {{incomplete"},
		{"VALID UNTIL '{{expiration}}'; {{", "VALID UNTIL '24h'; {{"},
		// This case tests whether `{{!{{password}}` successfully looks past the earlier unmatched {{
		{"}}backwards{{!{{password}}!", "}}backwards{{!world!"},
	} {
		assert.Equal(t, tc.expected, QueryHelper(tc.tpl, data),
			"template processing produced unexpected result")
	}
}

// TestQueryHelper_Recursion confirms QueryHelper does not replace placeholders that were themselves added as part of
// a replacement value.
func TestQueryHelper_Recursion(t *testing.T) {
	data := map[string]string{
		"a": "A{{a}}{{b}}{{c}}{{d}}",
		"b": "B{{a}}{{b}}{{c}}{{d}}",
		"c": "C{{a}}{{b}}{{c}}{{d}}",
		"d": "D{{a}}{{b}}{{c}}{{d}}",
	}
	assert.Equal(t, "A{{a}}{{b}}{{c}}{{d}}B{{a}}{{b}}{{c}}{{d}}C{{a}}{{b}}{{c}}{{d}}D{{a}}{{b}}{{c}}{{d}}",
		QueryHelper("{{a}}{{b}}{{c}}{{d}}", data),
		"template processing produced unexpected result")
}
