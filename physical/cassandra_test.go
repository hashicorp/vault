package physical

import (
	"os"
	"reflect"
	"testing"

	"github.com/hashicorp/vault/helper/logformat"
	log "github.com/mgutz/logxi/v1"

	_ "github.com/lib/pq"
)

func TestCassandraBackend(t *testing.T) {
	hosts := os.Getenv("CASSANDRA_HOSTS")
	if hosts == "" {
		t.SkipNow()
	}

	// Run vault tests
	logger := logformat.NewVaultLogger(log.LevelTrace)

	b, err := NewBackend("cassandra", logger, map[string]string{
		"hosts": hosts,
	})

	if err != nil {
		t.Fatalf("Failed to create new backend: %v", err)
	}

	testBackend(t, b)
	testBackend_ListPrefix(t, b)
}

func TestCassandraBackendBuckets(t *testing.T) {
	t.Parallel()
	expectations := map[string][]string{
		"":          {"."},
		"a":         {"."},
		"a/b":       {".", "a"},
		"a/b/c/d/e": {".", "a", "a/b", "a/b/c", "a/b/c/d"}}

	b := &CassandraBackend{}
	for input, expected := range expectations {
		actual := b.buckets(input)
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("bad: %v expected: %v", actual, expected)
		}
	}
}
