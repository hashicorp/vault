package physical

import (
	"os"
	"testing"

	_ "github.com/lib/pq"
)

func TestPostgresqlBackend(t *testing.T) {
	url := os.Getenv("POSTGRESQL_URL")
	if url == "" {
		t.SkipNow()
	}

	// Run vault tests
	b, err := NewBackend("postgresql", map[string]string{
		"url": url,
	})

	if err != nil {
		t.Fatalf("Failed to create new backend: %v", err)
	}

	defer func() {
		pg := b.(*PostgresqlBackend)
		_, err := pg.client.Exec("DROP TABLE " + pg.dbTable)
		if err != nil {
			t.Fatalf("Failed to drop table: %v", err)
		}
	}()

	testBackend(t, b)
	testBackend_ListPrefix(t, b)

}
