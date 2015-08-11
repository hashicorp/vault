package physical

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"
)

// PostGresBackend is a physical backend that stores data on postgres database
type PostGresBackend struct {
	Url string

	l sync.Mutex
}

var (
	Trace   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

// newFileBackend constructs a Filebackend using the given directory
func newPostGresBackend(conf map[string]string) (Backend, error) {

	Trace = log.New(ioutil.Discard, "Trace: ", log.Ldate|log.Ltime|log.Lshortfile)
	Info = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	Warning = log.New(os.Stdout, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(os.Stdout, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

	Trace.Println("> newPostGresBackend")
	url, ok := conf["url"]
	if !ok {
		return nil, fmt.Errorf("'url' must be set")
	}

	db, err := sql.Open("postgres", url)
	if err != nil {
		Error.Fatal(err)
	}
	defer db.Close()
	sqlStmt := "CREATE EXTENSION if not exists \"uuid-ossp\";"
	_, err = db.Exec(sqlStmt)
	if err != nil {
		Error.Printf("%q: %s\n", err, sqlStmt)
		Error.Println("< newPostGresBackend(err)")
		Error.Panic(err)
	}

	sqlStmt = `
CREATE TABLE if not exists vault (
  id UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
  key TEXT not null,
  value bytea,
  created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL,
  updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL
);`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		Error.Printf("%q: %s\n", err, sqlStmt)
		Error.Println("< newPostGresBackend(err)")
		Error.Panic(err)
	}

	Trace.Println("< newPostGresBackend")
	return &PostGresBackend{Url: url}, nil
}

func (b *PostGresBackend) Delete(k string) error {
	Trace.Println("> Delete")
	b.l.Lock()
	defer b.l.Unlock()

	db, err := sql.Open("postgres", b.Url)
	if err != nil {
		Trace.Println("< Delete error db")
		Error.Fatal(err)
	}
	defer db.Close()

	txn, err := db.Begin()
	if err != nil {
		Trace.Println("< Delete error begin")
		Trace.Fatal(err)
	}

	stmt, err := txn.Prepare("delete from vault where key = $1")
	if err != nil {
		Trace.Println("< Delete error prepare")
		Error.Fatal(err)
	}

	_, err = stmt.Exec(k)
	if err != nil {
		Trace.Println("< Delete error exec")
		Error.Fatal(err)
	}

	err = stmt.Close()
	if err != nil {
		Trace.Println("< Delete error close")
		Error.Fatal(err)
	}

	err = txn.Commit()
	if err != nil {
		Trace.Println("< Delete error commit")
		Error.Fatal(err)
	}

	Trace.Println("< Delete")
	return err
}

func (b *PostGresBackend) Get(k string) (*Entry, error) {
	Trace.Println("> Get")
	b.l.Lock()
	defer b.l.Unlock()

	db, err := sql.Open("postgres", b.Url)
	if err != nil {
		Trace.Println("< Get error db")
		Error.Fatal(err)
	}
	defer db.Close()

	stmt, err := db.Prepare("select value from vault where key = $1")
	defer stmt.Close()
	if err != nil {
		Trace.Println("< Get error select")
		Error.Fatal(err)
	}

	var value []byte
	row, err := stmt.Query(k)
	if err != nil {
		Error.Println(err)
		Trace.Println("< Get error query")
		return nil, err
	}

	if row.Next() {
		row.Scan(&value)
		entry := Entry{k, value}

		return &entry, nil
	}
	Trace.Println("< Get")
	return nil, err

}

func (b *PostGresBackend) Put(entry *Entry) error {
	Trace.Println("> Put")
	b.l.Lock()
	defer b.l.Unlock()

	db, err := sql.Open("postgres", b.Url)
	if err != nil {
		Trace.Println("< Put error db")
		Error.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("select key from vault where key = $1", entry.Key)
	if err != nil {
		Trace.Println("< Put error select")
		Error.Fatal(err)
	}
	defer rows.Close()

	txn, err := db.Begin()
	if err != nil {
		Trace.Println("< Put error tran begin")
		Error.Fatal(err)
	}
	next := rows.Next()
	Trace.Printf("D Put - is there next row? %t", next)

	// need to update if already there
	if next {
		stmt, err := txn.Prepare("update vault set value =  $1, updated_at = $2 where key = $3")
		if err != nil {
			Trace.Println("< Put error update")
			Error.Fatal(err)
		}

		time := time.Now()
		_, err = stmt.Exec(entry.Value, time, entry.Key)
		if err != nil {
			Trace.Println("< Put error exec")
			Error.Fatal(err)
		}

		err = stmt.Close()
		if err != nil {
			Trace.Println("< Put error close")
			Error.Fatal(err)
		}

	} else {
		stmt, err := txn.Prepare("insert into vault (key, value, created_at, updated_at) values ($1, $2, $3, $4)")
		if err != nil {
			Trace.Println("< Put error insert prepare")
			Error.Fatal(err)
		}

		time := time.Now()
		_, err = stmt.Exec(entry.Key, entry.Value, time, time)
		if err != nil {
			Trace.Println("< Put error exec")
			Error.Fatal(err)
		}

		err = stmt.Close()
		if err != nil {
			Trace.Println("< Put error close")
			Error.Fatal(err)
		}

	}
	err = txn.Commit()
	if err != nil {
		Trace.Println("< Put error commit")
		Error.Fatal(err)
	}
	Trace.Println("< Put")

	return err
}

func (b *PostGresBackend) List(prefix string) ([]string, error) {
	Trace.Println("> List")
	b.l.Lock()
	defer b.l.Unlock()

	db, err := sql.Open("postgres", b.Url)
	if err != nil {
		Trace.Println("< List error db")
		Error.Fatal(err)
	}
	defer db.Close()

	query := "%" + prefix + "%"

	rows, err := db.Query("select key from vault where key like $1", query)
	if err != nil {
		Trace.Println("< List error query")
		Error.Fatal(err)
	}
	defer rows.Close()

	result := make([]string, 0)

	for rows.Next() {
		var message string
		rows.Scan(&message)
		result = append(result, message)
		Trace.Println(message)
	}

	Trace.Println("< List")
	return result, nil
}
