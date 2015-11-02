// +build all integration

package gocql

import (
	"testing"
)

func TestErrorsParse(t *testing.T) {
	session := createSession(t)
	defer session.Close()

	if err := createTable(session, `CREATE TABLE gocql_test.errors_parse (id int primary key)`); err != nil {
		t.Fatal("create:", err)
	}

	if err := createTable(session, `CREATE TABLE gocql_test.errors_parse (id int primary key)`); err == nil {
		t.Fatal("Should have gotten already exists error from cassandra server.")
	} else {
		switch e := err.(type) {
		case *RequestErrAlreadyExists:
			if e.Table != "errors_parse" {
				t.Fatalf("expected error table to be 'errors_parse' but was %q", e.Table)
			}
		default:
			t.Fatalf("expected to get RequestErrAlreadyExists instead got %T", e)
		}
	}
}
