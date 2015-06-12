package physical

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestMySQLBackend(t *testing.T) {
	address := os.Getenv("MYSQL_ADDR")
	if address == "" {
		t.SkipNow()
	}

	database := os.Getenv("MYSQL_DB")
	if database == "" {
		database = "test"
	}

	table := os.Getenv("MYSQL_TABLE")
	if table == "" {
		table = "test"
	}

	username := os.Getenv("MYSQL_USERNAME")
	password := os.Getenv("MYSQL_PASSWORD")

	// Create MySQL handle for the database.
	db, err := sql.Open("mysql", username+":"+password+"@tcp("+address+")/"+database)

	if err != nil {
		t.Fatalf("Failed to open an handler with database: %v", err)
	}
	defer db.Close()

	// Prepare statement for creating table.
	create_stmt := "CREATE TABLE IF NOT EXISTS " + database + "." + table + "(num int, sqr int, PRIMARY KEY (num))"
	stmtCrt, err := db.Prepare(create_stmt)
	if err != nil {
		t.Fatalf("Failed to prepare statement: %v", err)
	}
	defer stmtCrt.Close()

	// Create table
	_, err = stmtCrt.Exec()
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}

	// Prepare statement for inserting data.
	insert_stmt := "INSERT INTO " + database + "." + table + " VALUES( ?, ? ) ON DUPLICATE KEY UPDATE sqr=VALUES(sqr)"
	stmtIns, err := db.Prepare(insert_stmt)
	if err != nil {
		t.Fatalf("Failed to prepare statement: %v", err)
	}
	defer stmtIns.Close()

	// Prepare statement for reading data.
	select_stmt := "SELECT sqr FROM " + database + "." + table + " WHERE num = ?"
	stmtOut, err := db.Prepare(select_stmt)
	if err != nil {
		t.Fatalf("Failed to prepare statement: %v", err)
	}
	defer stmtOut.Close()

	// Insert square numbers for 0-24 in the database
	for i := 0; i < 25; i++ {
		_, err = stmtIns.Exec(i, (i * i)) // Insert tuples (i, i^2)
		if err != nil {
			t.Fatalf("Failed to insert data: %v", err)
		}
	}

	var square int

	// Query the square-number of 13
	err = stmtOut.QueryRow(13).Scan(&square)
	if err != nil {
		t.Fatalf("Failed to query data: %v", err)
	}
	fmt.Printf("The square number of 13 is: %d", square)

	b, err := NewBackend("mysql", map[string]string{
		"address":  address,
		"database": database,
		"table":    table,
		"username": username,
		"password": password,
	})

	if err != nil {
		t.Fatalf("Failed to create new backend: %v", err)
	}

	testBackend(t, b)
	testBackend_ListPrefix(t, b)
}
