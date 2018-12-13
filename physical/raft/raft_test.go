package raft

import (
	"context"
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"testing"

	log "github.com/hashicorp/go-hclog"
	uuid "github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/helper/logging"
	"github.com/hashicorp/vault/physical"
)

func TestRaftBoltBackend(t *testing.T) {
	dir, err := ioutil.TempDir("", "vault")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	defer os.RemoveAll(dir)

	logger := logging.NewVaultLogger(log.Debug)

	b, err := NewRaftBackend(map[string]string{
		"path": dir,
	}, logger)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	physical.ExerciseBackend(t, b)
}

func TestRaftBadgerBackend(t *testing.T) {
	dir, err := ioutil.TempDir("", "vault")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	defer os.RemoveAll(dir)

	logger := logging.NewVaultLogger(log.Debug)

	b, err := NewRaftBackend(map[string]string{
		"path":    dir,
		"db_type": "badger",
	}, logger)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	physical.ExerciseBackend(t, b)
}

func BenchmarkDB_Puts(b *testing.B) {
	dir, err := ioutil.TempDir("", "vault")
	if err != nil {
		b.Fatalf("err: %s", err)
	}
	defer os.RemoveAll(dir)

	logger := logging.NewVaultLogger(log.Debug)

	badger, err := NewRaftBackend(map[string]string{
		"path":    dir,
		"db_type": "badger",
	}, logger)
	if err != nil {
		b.Fatalf("err: %s", err)
	}

	bolt, err := NewRaftBackend(map[string]string{
		"path": dir,
	}, logger)
	if err != nil {
		b.Fatalf("err: %s", err)
	}

	bench := func(b *testing.B, s physical.Backend, dataSize int, parallel bool) {
		data, err := uuid.GenerateRandomBytes(dataSize)
		if err != nil {
			b.Fatal(err)
		}

		ctx := context.Background()
		pe := &physical.Entry{
			Value: data,
		}
		testName := b.Name()

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			pe.Key = fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("%s-%d", testName, i))))
			err := s.Put(ctx, pe)
			if err != nil {
				b.Fatal(err)
			}
		}
	}

	b.Run("256b-badger", func(b *testing.B) { bench(b, badger, 256, false) })
	b.Run("256b-bolt", func(b *testing.B) { bench(b, bolt, 256, false) })
	b.Run("256kb-badger", func(b *testing.B) { bench(b, badger, 256*1024, false) })
	b.Run("256kb-bolt", func(b *testing.B) { bench(b, bolt, 256*1024, false) })

}

func BenchmarkDB_Snapshot(b *testing.B) {
	dir, err := ioutil.TempDir("", "vault")
	if err != nil {
		b.Fatalf("err: %s", err)
	}
	defer os.RemoveAll(dir)

	logger := logging.NewVaultLogger(log.Debug)

	badger, err := NewRaftBackend(map[string]string{
		"path":    dir,
		"db_type": "badger",
	}, logger)
	if err != nil {
		b.Fatalf("err: %s", err)
	}

	bolt, err := NewRaftBackend(map[string]string{
		"path": dir,
	}, logger)
	if err != nil {
		b.Fatalf("err: %s", err)
	}

	data, err := uuid.GenerateRandomBytes(256 * 1024)
	if err != nil {
		b.Fatal(err)
	}

	ctx := context.Background()
	pe := &physical.Entry{
		Value: data,
	}
	testName := b.Name()

	for i := 0; i < 1000; i++ {
		pe.Key = fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("%s-%d", testName, i))))
		err := badger.Put(ctx, pe)
		if err != nil {
			b.Fatal(err)
		}
		err = bolt.Put(ctx, pe)
		if err != nil {
			b.Fatal(err)
		}
	}

	type snap interface {
		Snapshot(context.Context, io.Writer) error
	}
	bench := func(b *testing.B, s snap) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			pe.Key = fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("%s-%d", testName, i))))
			err := s.Snapshot(ctx, ioutil.Discard)
			if err != nil {
				b.Fatal(err)
			}
		}
	}

	b.Run("256b-badger", func(b *testing.B) { bench(b, badger) })
	b.Run("256b-bolt", func(b *testing.B) { bench(b, bolt) })

}
