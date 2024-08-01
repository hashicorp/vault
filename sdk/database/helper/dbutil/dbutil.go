// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package dbutil

import (
	"errors"
	"regexp"
	"strings"

	"github.com/hashicorp/vault/sdk/database/dbplugin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrEmptyCreationStatement = errors.New("empty creation statements")
	ErrEmptyRotationStatement = errors.New("empty rotation statements")
)

var queryHelperRegex = regexp.MustCompile(`{{[^{}]+}}`)

// QueryHelper evaluates a simple string template syntax by replacing {{value}}
// placeholders with the values from the supplied data map. Despite the name,
// it is NOT only used to template queries - it is also used for connection
// URIs or DSNs. Since it has no idea of the specific syntax into which it is
// templating, it does not perform any escaping. Unbalanced opening and closing
// brace sequences are passed through as is, as are any placeholders for which
// there is no key found in the map.
func QueryHelper(tpl string, data map[string]string) string {
	return queryHelperRegex.ReplaceAllStringFunc(tpl, func(s string) string {
		replacement, ok := data[s[2:len(s)-2]]
		if ok {
			return replacement
		} else {
			return s
		}
	})
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
