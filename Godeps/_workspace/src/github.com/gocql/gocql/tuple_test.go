// +build all integration

package gocql

import "testing"

func TestTupleSimple(t *testing.T) {
	if *flagProto < protoVersion3 {
		t.Skip("tuple types are only available of proto>=3")
	}

	session := createSession(t)
	defer session.Close()

	err := createTable(session, `CREATE TABLE tuple_test(
		id int,
		coord frozen<tuple<int, int>>,

		primary key(id))`)
	if err != nil {
		t.Fatal(err)
	}

	err = session.Query("INSERT INTO tuple_test(id, coord) VALUES(?, (?, ?))", 1, 100, -100).Exec()
	if err != nil {
		t.Fatal(err)
	}

	var (
		id    int
		coord struct {
			x int
			y int
		}
	)

	iter := session.Query("SELECT id, coord FROM tuple_test WHERE id=?", 1)
	if err := iter.Scan(&id, &coord.x, &coord.y); err != nil {
		t.Fatal(err)
	}

	if id != 1 {
		t.Errorf("expected to get id=1 got: %v", id)
	}
	if coord.x != 100 {
		t.Errorf("expected to get coord.x=100 got: %v", coord.x)
	}
	if coord.y != -100 {
		t.Errorf("expected to get coord.y=-100 got: %v", coord.y)
	}
}
