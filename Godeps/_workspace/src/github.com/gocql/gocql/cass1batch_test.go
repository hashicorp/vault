// +build all integration

package gocql

import (
	"strings"
	"testing"
)

func TestProto1BatchInsert(t *testing.T) {
	session := createSession(t)
	if err := session.Query("CREATE TABLE large (id int primary key)").Exec(); err != nil {
		t.Fatal("create table:", err)
	}
	defer session.Close()

	begin := "BEGIN BATCH"
	end := "APPLY BATCH"
	query := "INSERT INTO large (id) VALUES (?)"
	fullQuery := strings.Join([]string{begin, query, end}, "\n")
	args := []interface{}{5}
	if err := session.Query(fullQuery, args...).Consistency(Quorum).Exec(); err != nil {
		t.Fatal(err)
	}

}

func TestShouldPrepareFunction(t *testing.T) {
	var shouldPrepareTests = []struct {
		Stmt   string
		Result bool
	}{
		{`
      BEGIN BATCH
        INSERT INTO users (userID, password)
        VALUES ('smith', 'secret')
      APPLY BATCH
    ;
      `, true},
		{`INSERT INTO users (userID, password, name) VALUES ('user2', 'ch@ngem3b', 'second user')`, true},
		{`BEGIN COUNTER BATCH UPDATE stats SET views = views + 1 WHERE pageid = 1 APPLY BATCH`, true},
		{`delete name from users where userID = 'smith';`, true},
		{`  UPDATE users SET password = 'secret' WHERE userID = 'smith'   `, true},
		{`CREATE TABLE users (
        user_name varchar PRIMARY KEY,
        password varchar,
        gender varchar,
        session_token varchar,
        state varchar,
        birth_year bigint
      );`, false},
	}

	for _, test := range shouldPrepareTests {
		q := &Query{stmt: test.Stmt}
		if got := q.shouldPrepare(); got != test.Result {
			t.Fatalf("%q: got %v, expected %v\n", test.Stmt, got, test.Result)
		}
	}
}
