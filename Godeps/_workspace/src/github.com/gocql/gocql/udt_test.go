// +build all integration

package gocql

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

type position struct {
	Lat     int    `cql:"lat"`
	Lon     int    `cql:"lon"`
	Padding string `json:"padding"`
}

// NOTE: due to current implementation details it is not currently possible to use
// a pointer receiver type for the UDTMarshaler interface to handle UDT's
func (p position) MarshalUDT(name string, info TypeInfo) ([]byte, error) {
	switch name {
	case "lat":
		return Marshal(info, p.Lat)
	case "lon":
		return Marshal(info, p.Lon)
	case "padding":
		return Marshal(info, p.Padding)
	default:
		return nil, fmt.Errorf("unknown column for position: %q", name)
	}
}

func (p *position) UnmarshalUDT(name string, info TypeInfo, data []byte) error {
	switch name {
	case "lat":
		return Unmarshal(info, data, &p.Lat)
	case "lon":
		return Unmarshal(info, data, &p.Lon)
	case "padding":
		return Unmarshal(info, data, &p.Padding)
	default:
		return fmt.Errorf("unknown column for position: %q", name)
	}
}

func TestUDT_Marshaler(t *testing.T) {
	if *flagProto < protoVersion3 {
		t.Skip("UDT are only available on protocol >= 3")
	}

	session := createSession(t)
	defer session.Close()

	err := createTable(session, `CREATE TYPE gocql_test.position(
		lat int,
		lon int,
		padding text);`)
	if err != nil {
		t.Fatal(err)
	}

	err = createTable(session, `CREATE TABLE gocql_test.houses(
		id int,
		name text,
		loc frozen<position>,

		primary key(id)
	);`)
	if err != nil {
		t.Fatal(err)
	}

	const (
		expLat = -1
		expLon = 2
	)
	pad := strings.Repeat("X", 1000)

	err = session.Query("INSERT INTO houses(id, name, loc) VALUES(?, ?, ?)", 1, "test", &position{expLat, expLon, pad}).Exec()
	if err != nil {
		t.Fatal(err)
	}

	pos := &position{}

	err = session.Query("SELECT loc FROM houses WHERE id = ?", 1).Scan(pos)
	if err != nil {
		t.Fatal(err)
	}

	if pos.Lat != expLat {
		t.Errorf("expeceted lat to be be %d got %d", expLat, pos.Lat)
	}
	if pos.Lon != expLon {
		t.Errorf("expeceted lon to be be %d got %d", expLon, pos.Lon)
	}
	if pos.Padding != pad {
		t.Errorf("expected to get padding %q got %q\n", pad, pos.Padding)
	}
}

func TestUDT_Reflect(t *testing.T) {
	if *flagProto < protoVersion3 {
		t.Skip("UDT are only available on protocol >= 3")
	}

	// Uses reflection instead of implementing the marshaling type
	session := createSession(t)
	defer session.Close()

	err := createTable(session, `CREATE TYPE gocql_test.horse(
		name text,
		owner text);`)
	if err != nil {
		t.Fatal(err)
	}

	err = createTable(session, `CREATE TABLE gocql_test.horse_race(
		position int,
		horse frozen<horse>,

		primary key(position)
	);`)
	if err != nil {
		t.Fatal(err)
	}

	type horse struct {
		Name  string `cql:"name"`
		Owner string `cql:"owner"`
	}

	insertedHorse := &horse{
		Name:  "pony",
		Owner: "jim",
	}

	err = session.Query("INSERT INTO horse_race(position, horse) VALUES(?, ?)", 1, insertedHorse).Exec()
	if err != nil {
		t.Fatal(err)
	}

	retrievedHorse := &horse{}
	err = session.Query("SELECT horse FROM horse_race WHERE position = ?", 1).Scan(retrievedHorse)
	if err != nil {
		t.Fatal(err)
	}

	if *retrievedHorse != *insertedHorse {
		t.Fatal("exepcted to get %+v got %+v", insertedHorse, retrievedHorse)
	}
}

func TestUDT_Proto2error(t *testing.T) {
	if *flagProto < protoVersion3 {
		t.Skip("UDT are only available on protocol >= 3")
	}

	cluster := createCluster()
	cluster.ProtoVersion = 2
	cluster.Keyspace = "gocql_test"

	// Uses reflection instead of implementing the marshaling type
	session, err := cluster.CreateSession()
	if err != nil {
		t.Fatal(err)
	}
	defer session.Close()

	err = createTable(session, `CREATE TYPE gocql_test.fish(
		name text,
		owner text);`)
	if err != nil {
		t.Fatal(err)
	}

	err = createTable(session, `CREATE TABLE gocql_test.fish_race(
		position int,
		fish frozen<fish>,

		primary key(position)
	);`)
	if err != nil {
		t.Fatal(err)
	}

	type fish struct {
		Name  string `cql:"name"`
		Owner string `cql:"owner"`
	}

	insertedFish := &fish{
		Name:  "pony",
		Owner: "jim",
	}

	err = session.Query("INSERT INTO fish_race(position, fish) VALUES(?, ?)", 1, insertedFish).Exec()
	if err != ErrorUDTUnavailable {
		t.Fatalf("expected to get %v got %v", ErrorUDTUnavailable, err)
	}
}

func TestUDT_NullObject(t *testing.T) {
	if *flagProto < protoVersion3 {
		t.Skip("UDT are only available on protocol >= 3")
	}

	session := createSession(t)
	defer session.Close()

	err := createTable(session, `CREATE TYPE gocql_test.udt_null_type(
		name text,
		owner text);`)
	if err != nil {
		t.Fatal(err)
	}

	err = createTable(session, `CREATE TABLE gocql_test.udt_null_table(
		id uuid,
		udt_col frozen<udt_null_type>,

		primary key(id)
	);`)
	if err != nil {
		t.Fatal(err)
	}

	type col struct {
		Name  string `cql:"name"`
		Owner string `cql:"owner"`
	}

	id := TimeUUID()
	err = session.Query("INSERT INTO udt_null_table(id) VALUES(?)", id).Exec()
	if err != nil {
		t.Fatal(err)
	}

	readCol := &col{
		Name:  "temp",
		Owner: "temp",
	}

	err = session.Query("SELECT udt_col FROM udt_null_table WHERE id = ?", id).Scan(readCol)
	if err != nil {
		t.Fatal(err)
	}

	if readCol.Name != "" {
		t.Errorf("expected empty string to be returned for null udt: got %q", readCol.Name)
	}
	if readCol.Owner != "" {
		t.Errorf("expected empty string to be returned for null udt: got %q", readCol.Owner)
	}
}

func TestMapScanUDT(t *testing.T) {
	if *flagProto < protoVersion3 {
		t.Skip("UDT are only available on protocol >= 3")
	}

	session := createSession(t)
	defer session.Close()

	err := createTable(session, `CREATE TYPE gocql_test.log_entry (
		created_timestamp timestamp,
		message text
	);`)
	if err != nil {
		t.Fatal(err)
	}

	err = createTable(session, `CREATE TABLE gocql_test.requests_by_id (
		id uuid PRIMARY KEY,
		type int,
		log_entries list<frozen <log_entry>>
	);`)
	if err != nil {
		t.Fatal(err)
	}

	entry := []struct {
		CreatedTimestamp time.Time `cql:"created_timestamp"`
		Message          string    `cql:"message"`
	}{
		{
			CreatedTimestamp: time.Now().Truncate(time.Millisecond),
			Message:          "test time now",
		},
	}

	id, _ := RandomUUID()
	const typ = 1

	err = session.Query("INSERT INTO requests_by_id(id, type, log_entries) VALUES (?, ?, ?)", id, typ, entry).Exec()
	if err != nil {
		t.Fatal(err)
	}

	rawResult := map[string]interface{}{}
	err = session.Query(`SELECT * FROM requests_by_id WHERE id = ?`, id).MapScan(rawResult)
	if err != nil {
		t.Fatal(err)
	}

	logEntries, ok := rawResult["log_entries"].([]map[string]interface{})
	if !ok {
		t.Fatal("log_entries not in scanned map")
	}

	if len(logEntries) != 1 {
		t.Fatalf("expected to get 1 log_entry got %d", len(logEntries))
	}

	logEntry := logEntries[0]

	timestamp, ok := logEntry["created_timestamp"]
	if !ok {
		t.Error("created_timestamp not unmarshalled into map")
	} else {
		if ts, ok := timestamp.(time.Time); ok {
			if !ts.In(time.UTC).Equal(entry[0].CreatedTimestamp.In(time.UTC)) {
				t.Errorf("created_timestamp not equal to stored: got %v expected %v", ts.In(time.UTC), entry[0].CreatedTimestamp.In(time.UTC))
			}
		} else {
			t.Errorf("created_timestamp was not time.Time got: %T", timestamp)
		}
	}

	message, ok := logEntry["message"]
	if !ok {
		t.Error("message not unmarshalled into map")
	} else {
		if ts, ok := message.(string); ok {
			if ts != message {
				t.Errorf("message not equal to stored: got %v expected %v", ts, entry[0].Message)
			}
		} else {
			t.Errorf("message was not string got: %T", message)
		}
	}

}
