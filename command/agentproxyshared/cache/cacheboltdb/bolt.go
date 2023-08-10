// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cacheboltdb

import (
	"context"
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/hashicorp/go-hclog"
	wrapping "github.com/hashicorp/go-kms-wrapping/v2"
	"github.com/hashicorp/go-multierror"
	bolt "go.etcd.io/bbolt"
)

const (
	// Keep track of schema version for future migrations
	storageVersionKey = "version"
	storageVersion    = "2" // v2 merges auth-lease and secret-lease buckets into one ordered bucket

	// DatabaseFileName - filename for the persistent cache file
	DatabaseFileName = "vault-agent-cache.db"

	// metaBucketName - naming the meta bucket that holds the version and
	// bootstrapping keys
	metaBucketName = "meta"

	// DEPRECATED: secretLeaseType - v1 Bucket/type for leases with secret info
	secretLeaseType = "secret-lease"

	// DEPRECATED: authLeaseType - v1 Bucket/type for leases with auth info
	authLeaseType = "auth-lease"

	// TokenType - Bucket/type for auto-auth tokens
	TokenType = "token"

	// LeaseType - v2 Bucket/type for auth AND secret leases.
	//
	// This bucket stores keys in the same order they were created using
	// auto-incrementing keys and the fact that BoltDB stores keys in byte
	// slice order. This means when we iterate through this bucket during
	// restore, we will always restore parent tokens before their children,
	// allowing us to correctly attach child contexts to their parent's context.
	LeaseType = "lease"

	// lookupType - v2 Bucket/type to map from a memcachedb index ID to an
	// auto-incrementing BoltDB key. Facilitates deletes from the lease
	// bucket using an ID instead of the auto-incrementing BoltDB key.
	lookupType = "lookup"

	// AutoAuthToken - key for the latest auto-auth token
	AutoAuthToken = "auto-auth-token"

	// RetrievalTokenMaterial is the actual key or token in the key bucket
	RetrievalTokenMaterial = "retrieval-token-material"
)

// BoltStorage is a persistent cache using a bolt db. Items are organized with
// the version and bootstrapping items in the "meta" bucket, and tokens, auth
// leases, and secret leases in their own buckets.
type BoltStorage struct {
	db      *bolt.DB
	logger  hclog.Logger
	wrapper wrapping.Wrapper
	aad     string
}

// BoltStorageConfig is the collection of input parameters for setting up bolt
// storage
type BoltStorageConfig struct {
	Path    string
	Logger  hclog.Logger
	Wrapper wrapping.Wrapper
	AAD     string
}

// NewBoltStorage opens a new bolt db at the specified file path and returns it.
// If the db already exists the buckets will just be created if they don't
// exist.
func NewBoltStorage(config *BoltStorageConfig) (*BoltStorage, error) {
	dbPath := filepath.Join(config.Path, DatabaseFileName)
	db, err := bolt.Open(dbPath, 0o600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, err
	}
	err = db.Update(func(tx *bolt.Tx) error {
		return createBoltSchema(tx, storageVersion)
	})
	if err != nil {
		return nil, err
	}
	bs := &BoltStorage{
		db:      db,
		logger:  config.Logger,
		wrapper: config.Wrapper,
		aad:     config.AAD,
	}
	return bs, nil
}

func createBoltSchema(tx *bolt.Tx, createVersion string) error {
	switch {
	case createVersion == "1":
		if err := createV1BoltSchema(tx); err != nil {
			return err
		}
	case createVersion == "2":
		if err := createV2BoltSchema(tx); err != nil {
			return err
		}
	default:
		return fmt.Errorf("schema version %s not supported", createVersion)
	}

	meta, err := tx.CreateBucketIfNotExists([]byte(metaBucketName))
	if err != nil {
		return fmt.Errorf("failed to create bucket %s: %w", metaBucketName, err)
	}

	// Check and set file version in the meta bucket.
	version := meta.Get([]byte(storageVersionKey))
	switch {
	case version == nil:
		err = meta.Put([]byte(storageVersionKey), []byte(createVersion))
		if err != nil {
			return fmt.Errorf("failed to set storage version: %w", err)
		}

		return nil

	case string(version) == createVersion:
		return nil

	case string(version) == "1" && createVersion == "2":
		return migrateFromV1ToV2Schema(tx)

	default:
		return fmt.Errorf("storage migration from %s to %s not implemented", string(version), createVersion)
	}
}

func createV1BoltSchema(tx *bolt.Tx) error {
	// Create the buckets for tokens and leases.
	for _, bucket := range []string{TokenType, authLeaseType, secretLeaseType} {
		if _, err := tx.CreateBucketIfNotExists([]byte(bucket)); err != nil {
			return fmt.Errorf("failed to create %s bucket: %w", bucket, err)
		}
	}

	return nil
}

func createV2BoltSchema(tx *bolt.Tx) error {
	// Create the buckets for tokens and leases.
	for _, bucket := range []string{TokenType, LeaseType, lookupType} {
		if _, err := tx.CreateBucketIfNotExists([]byte(bucket)); err != nil {
			return fmt.Errorf("failed to create %s bucket: %w", bucket, err)
		}
	}

	return nil
}

func migrateFromV1ToV2Schema(tx *bolt.Tx) error {
	if err := createV2BoltSchema(tx); err != nil {
		return err
	}

	for _, v1BucketType := range []string{authLeaseType, secretLeaseType} {
		if bucket := tx.Bucket([]byte(v1BucketType)); bucket != nil {
			bucket.ForEach(func(key, value []byte) error {
				autoIncKey, err := autoIncrementedLeaseKey(tx, string(key))
				if err != nil {
					return fmt.Errorf("error migrating %s %q key to auto incremented key: %w", v1BucketType, string(key), err)
				}
				if err := tx.Bucket([]byte(LeaseType)).Put(autoIncKey, value); err != nil {
					return fmt.Errorf("error migrating %s %q from v1 to v2 schema: %w", v1BucketType, string(key), err)
				}
				return nil
			})

			if err := tx.DeleteBucket([]byte(v1BucketType)); err != nil {
				return fmt.Errorf("failed to clean up %s bucket during v1 to v2 schema migration: %w", v1BucketType, err)
			}
		}
	}

	meta, err := tx.CreateBucketIfNotExists([]byte(metaBucketName))
	if err != nil {
		return fmt.Errorf("failed to create meta bucket: %w", err)
	}
	if err := meta.Put([]byte(storageVersionKey), []byte(storageVersion)); err != nil {
		return fmt.Errorf("failed to update schema from v1 to v2: %w", err)
	}

	return nil
}

func autoIncrementedLeaseKey(tx *bolt.Tx, id string) ([]byte, error) {
	leaseBucket := tx.Bucket([]byte(LeaseType))
	keyValue, err := leaseBucket.NextSequence()
	if err != nil {
		return nil, fmt.Errorf("failed to generate lookup key for id %q: %w", id, err)
	}

	key := make([]byte, 8)
	// MUST be big endian, because keys are ordered by byte slice comparison
	// which progressively compares each byte in the slice starting at index 0.
	// BigEndian in the range [255-257] looks like this:
	// [0 0 0 0 0 0 0 255]
	// [0 0 0 0 0 0 1 0]
	// [0 0 0 0 0 0 1 1]
	// LittleEndian in the same range looks like this:
	// [255 0 0 0 0 0 0 0]
	// [0 1 0 0 0 0 0 0]
	// [1 1 0 0 0 0 0 0]
	binary.BigEndian.PutUint64(key, keyValue)

	err = tx.Bucket([]byte(lookupType)).Put([]byte(id), key)
	if err != nil {
		return nil, err
	}

	return key, nil
}

// Set an index (token or lease) in bolt storage
func (b *BoltStorage) Set(ctx context.Context, id string, plaintext []byte, indexType string) error {
	blob, err := b.wrapper.Encrypt(ctx, plaintext, wrapping.WithAad([]byte(b.aad)))
	if err != nil {
		return fmt.Errorf("error encrypting %s index: %w", indexType, err)
	}

	protoBlob, err := proto.Marshal(blob)
	if err != nil {
		return err
	}

	return b.db.Update(func(tx *bolt.Tx) error {
		var key []byte
		switch indexType {
		case LeaseType:
			// If this is a lease type, generate an auto-incrementing key and
			// store an ID -> key lookup entry
			key, err = autoIncrementedLeaseKey(tx, id)
			if err != nil {
				return err
			}
		case TokenType:
			// If this is an auto-auth token, also stash it in the meta bucket for
			// easy retrieval upon restore
			key = []byte(id)
			meta := tx.Bucket([]byte(metaBucketName))
			if err := meta.Put([]byte(AutoAuthToken), protoBlob); err != nil {
				return fmt.Errorf("failed to set latest auto-auth token: %w", err)
			}
		default:
			return fmt.Errorf("called Set for unsupported type %q", indexType)
		}
		s := tx.Bucket([]byte(indexType))
		if s == nil {
			return fmt.Errorf("bucket %q not found", indexType)
		}
		return s.Put(key, protoBlob)
	})
}

// Delete an index (token or lease) by key from bolt storage
func (b *BoltStorage) Delete(id string, indexType string) error {
	return b.db.Update(func(tx *bolt.Tx) error {
		key := []byte(id)
		if indexType == LeaseType {
			key = tx.Bucket([]byte(lookupType)).Get(key)
			if key == nil {
				return fmt.Errorf("failed to lookup bolt DB key for id %q", id)
			}

			err := tx.Bucket([]byte(lookupType)).Delete([]byte(id))
			if err != nil {
				return fmt.Errorf("failed to delete %q from lookup bucket: %w", id, err)
			}
		}

		bucket := tx.Bucket([]byte(indexType))
		if bucket == nil {
			return fmt.Errorf("bucket %q not found during delete", indexType)
		}
		if err := bucket.Delete(key); err != nil {
			return fmt.Errorf("failed to delete %q from %q bucket: %w", id, indexType, err)
		}
		b.logger.Trace("deleted index from bolt db", "id", id)
		return nil
	})
}

func (b *BoltStorage) decrypt(ctx context.Context, ciphertext []byte) ([]byte, error) {
	var blob wrapping.BlobInfo
	if err := proto.Unmarshal(ciphertext, &blob); err != nil {
		return nil, err
	}

	return b.wrapper.Decrypt(ctx, &blob, wrapping.WithAad([]byte(b.aad)))
}

// GetByType returns a list of stored items of the specified type
func (b *BoltStorage) GetByType(ctx context.Context, indexType string) ([][]byte, error) {
	var returnBytes [][]byte

	err := b.db.View(func(tx *bolt.Tx) error {
		var errors *multierror.Error

		bucket := tx.Bucket([]byte(indexType))
		if bucket == nil {
			return fmt.Errorf("bucket %q not found", indexType)
		}
		bucket.ForEach(func(key, ciphertext []byte) error {
			plaintext, err := b.decrypt(ctx, ciphertext)
			if err != nil {
				errors = multierror.Append(errors, fmt.Errorf("error decrypting entry %s: %w", key, err))
				return nil
			}

			returnBytes = append(returnBytes, plaintext)
			return nil
		})
		return errors.ErrorOrNil()
	})

	return returnBytes, err
}

// GetAutoAuthToken retrieves the latest auto-auth token, and returns nil if non
// exists yet
func (b *BoltStorage) GetAutoAuthToken(ctx context.Context) ([]byte, error) {
	var encryptedToken []byte

	err := b.db.View(func(tx *bolt.Tx) error {
		meta := tx.Bucket([]byte(metaBucketName))
		if meta == nil {
			return fmt.Errorf("bucket %q not found", metaBucketName)
		}
		value := meta.Get([]byte(AutoAuthToken))
		if value != nil {
			encryptedToken = make([]byte, len(value))
			copy(encryptedToken, value)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	if encryptedToken == nil {
		return nil, nil
	}

	plaintext, err := b.decrypt(ctx, encryptedToken)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt auto-auth token: %w", err)
	}
	return plaintext, nil
}

// GetRetrievalToken retrieves a plaintext token from the KeyBucket, which will
// be used by the key manager to retrieve the encryption key, nil if none set
func (b *BoltStorage) GetRetrievalToken() ([]byte, error) {
	var token []byte

	err := b.db.View(func(tx *bolt.Tx) error {
		metaBucket := tx.Bucket([]byte(metaBucketName))
		if metaBucket == nil {
			return fmt.Errorf("bucket %q not found", metaBucketName)
		}
		value := metaBucket.Get([]byte(RetrievalTokenMaterial))
		if value != nil {
			token = make([]byte, len(value))
			copy(token, value)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return token, err
}

// StoreRetrievalToken sets plaintext token material in the RetrievalTokenBucket
func (b *BoltStorage) StoreRetrievalToken(token []byte) error {
	return b.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(metaBucketName))
		if bucket == nil {
			return fmt.Errorf("bucket %q not found", metaBucketName)
		}
		return bucket.Put([]byte(RetrievalTokenMaterial), token)
	})
}

// Close the boltdb
func (b *BoltStorage) Close() error {
	b.logger.Trace("closing bolt db", "path", b.db.Path())
	return b.db.Close()
}

// Clear the boltdb by deleting all the token and lease buckets and recreating
// the schema/layout
func (b *BoltStorage) Clear() error {
	return b.db.Update(func(tx *bolt.Tx) error {
		for _, name := range []string{TokenType, LeaseType, lookupType} {
			b.logger.Trace("deleting bolt bucket", "name", name)
			if err := tx.DeleteBucket([]byte(name)); err != nil {
				return err
			}
		}
		return createBoltSchema(tx, storageVersion)
	})
}

// DBFileExists checks whether the vault agent cache file at `filePath` exists
func DBFileExists(path string) (bool, error) {
	checkFile, err := os.OpenFile(filepath.Join(path, DatabaseFileName), os.O_RDWR, 0o600)
	defer checkFile.Close()
	switch {
	case err == nil:
		return true, nil
	case os.IsNotExist(err):
		return false, nil
	default:
		return false, fmt.Errorf("failed to check if bolt file exists at path %s: %w", path, err)
	}
}
