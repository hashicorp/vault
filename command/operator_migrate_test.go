package command

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/go-test/deep"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/command/server"
	"github.com/hashicorp/vault/helper/base62"
	"github.com/hashicorp/vault/helper/testhelpers"
	"github.com/hashicorp/vault/physical"
	"github.com/hashicorp/vault/vault"
)

const trailing_slash_key = "trailing_slash/"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func TestMigration(t *testing.T) {
	t.Run("Default", func(t *testing.T) {
		data := generateData()

		fromFactory := physicalBackends["file"]

		folder := filepath.Join(os.TempDir(), testhelpers.RandomWithPrefix("migrator"))
		defer os.RemoveAll(folder)
		confFrom := map[string]string{
			"path": folder,
		}

		from, err := fromFactory(confFrom, nil)
		if err != nil {
			t.Fatal(err)
		}
		if err := storeData(from, data); err != nil {
			t.Fatal(err)
		}

		toFactory := physicalBackends["inmem"]
		confTo := map[string]string{}
		to, err := toFactory(confTo, nil)
		if err != nil {
			t.Fatal(err)
		}

		cmd := OperatorMigrateCommand{
			logger: log.NewNullLogger(),
		}
		if err := cmd.migrateAll(context.Background(), from, to); err != nil {
			t.Fatal(err)
		}

		if err := compareStoredData(to, data, ""); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Start option", func(t *testing.T) {
		data := generateData()

		fromFactory := physicalBackends["inmem"]
		confFrom := map[string]string{}
		from, err := fromFactory(confFrom, nil)
		if err != nil {
			t.Fatal(err)
		}
		if err := storeData(from, data); err != nil {
			t.Fatal(err)
		}

		toFactory := physicalBackends["file"]
		folder := filepath.Join(os.TempDir(), testhelpers.RandomWithPrefix("migrator"))
		defer os.RemoveAll(folder)
		confTo := map[string]string{
			"path": folder,
		}

		to, err := toFactory(confTo, nil)
		if err != nil {
			t.Fatal(err)
		}

		const start = "m"

		cmd := OperatorMigrateCommand{
			logger:    log.NewNullLogger(),
			flagStart: start,
		}
		if err := cmd.migrateAll(context.Background(), from, to); err != nil {
			t.Fatal(err)
		}

		if err := compareStoredData(to, data, start); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Config parsing", func(t *testing.T) {
		cmd := new(OperatorMigrateCommand)

		cfgName := filepath.Join(os.TempDir(), testhelpers.RandomWithPrefix("migrator"))
		ioutil.WriteFile(cfgName, []byte(`
storage_source "src_type" {
  path = "src_path"
}

storage_destination "dest_type" {
  path = "dest_path"
}`), 0644)
		defer os.Remove(cfgName)

		expCfg := &migratorConfig{
			StorageSource: &server.Storage{
				Type: "src_type",
				Config: map[string]string{
					"path": "src_path",
				},
			},
			StorageDestination: &server.Storage{
				Type: "dest_type",
				Config: map[string]string{
					"path": "dest_path",
				},
			},
		}
		cfg, err := cmd.loadMigratorConfig(cfgName)
		if err != nil {
			t.Fatal(cfg)
		}
		if diff := deep.Equal(cfg, expCfg); diff != nil {
			t.Fatal(diff)
		}

		verifyBad := func(cfg string) {
			ioutil.WriteFile(cfgName, []byte(cfg), 0644)
			_, err := cmd.loadMigratorConfig(cfgName)
			if err == nil {
				t.Fatalf("expected error but none received from: %v", cfg)
			}
		}

		// missing source
		verifyBad(`
storage_destination "dest_type" {
  path = "dest_path"
}`)

		// missing destination
		verifyBad(`
storage_source "src_type" {
  path = "src_path"
}`)

		// duplicate source
		verifyBad(`
storage_source "src_type" {
  path = "src_path"
}

storage_source "src_type2" {
  path = "src_path"
}

storage_destination "dest_type" {
  path = "dest_path"
}`)

		// duplicate destination
		verifyBad(`
storage_source "src_type" {
  path = "src_path"
}

storage_destination "dest_type" {
  path = "dest_path"
}

storage_destination "dest_type2" {
  path = "dest_path"
}`)

	})
	t.Run("DFS Scan", func(t *testing.T) {
		s, _ := physicalBackends["inmem"](map[string]string{}, nil)

		data := generateData()
		data["cc"] = []byte{}
		data["c/d/e/f"] = []byte{}
		data["c/d/e/g"] = []byte{}
		data["c"] = []byte{}
		storeData(s, data)

		l := randomLister{s}

		var out []string
		dfsScan(context.Background(), l, func(ctx context.Context, path string) error {
			out = append(out, path)
			return nil
		})

		delete(data, trailing_slash_key)
		delete(data, "")

		var keys []string
		for key := range data {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		if !reflect.DeepEqual(keys, out) {
			t.Fatalf("expected equal: %v, %v", keys, out)
		}
	})
}

// randomLister wraps a physical backend, providing a List method
// that returns results in a random order.
type randomLister struct {
	b physical.Backend
}

func (l randomLister) List(ctx context.Context, path string) ([]string, error) {
	result, err := l.b.List(ctx, path)
	if err != nil {
		return nil, err
	}
	rand.Shuffle(len(result), func(i, j int) {
		result[i], result[j] = result[j], result[i]
	})
	return result, err
}

func (l randomLister) Get(ctx context.Context, path string) (*physical.Entry, error) {
	return l.b.Get(ctx, path)
}

func (l randomLister) Put(ctx context.Context, entry *physical.Entry) error {
	return l.b.Put(ctx, entry)
}

func (l randomLister) Delete(ctx context.Context, path string) error {
	return l.b.Delete(ctx, path)
}

// generateData creates a map of 500 random keys and values
func generateData() map[string][]byte {
	result := make(map[string][]byte)
	for i := 0; i < 500; i++ {
		segments := make([]string, rand.Intn(8)+1)
		for j := 0; j < len(segments); j++ {
			s, _ := base62.Random(6)
			segments[j] = s
		}
		data := make([]byte, 100)
		rand.Read(data)
		result[strings.Join(segments, "/")] = data
	}

	// Add special keys that should be excluded from migration
	result[storageMigrationLock] = []byte{}
	result[vault.CoreLockPath] = []byte{}

	// Empty keys are now prevented in Vault, but older data sets
	// might contain them.
	result[""] = []byte{}
	result[trailing_slash_key] = []byte{}

	return result
}

func storeData(s physical.Backend, ref map[string][]byte) error {
	for k, v := range ref {
		entry := physical.Entry{
			Key:   k,
			Value: v,
		}

		err := s.Put(context.Background(), &entry)
		if err != nil {
			return err
		}
	}
	return nil
}

func compareStoredData(s physical.Backend, ref map[string][]byte, start string) error {
	for k, v := range ref {
		entry, err := s.Get(context.Background(), k)
		if err != nil {
			return err
		}

		if k == storageMigrationLock || k == vault.CoreLockPath || k == "" || strings.HasSuffix(k, "/") {
			if entry == nil {
				continue
			}
			return fmt.Errorf("key found that should have been excluded: %s", k)
		}

		if k >= start {
			if entry == nil {
				return fmt.Errorf("key not found: %s", k)
			}
			if !bytes.Equal(v, entry.Value) {
				return fmt.Errorf("values differ for key: %s", k)
			}
		} else {
			if entry != nil {
				return fmt.Errorf("found key the should have been skipped by start option: %s", k)
			}
		}
	}

	return nil
}
