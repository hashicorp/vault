package cacheboltdb

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/hashicorp/go-hclog"
	wrapping "github.com/hashicorp/go-kms-wrapping"
	"github.com/hashicorp/go-multierror"
	bolt "go.etcd.io/bbolt"
)

const (
	// Keep track of schema version for future migrations
	storageVersionKey = "version"
	storageVersion    = "1"

	// DatabaseFileName - filename for the persistent cache file
	DatabaseFileName = "vault-agent-cache.db"

	// metaBucketName - naming the meta bucket that holds the version and
	// bootstrapping keys
	metaBucketName = "meta"

	// SecretLeaseType - Bucket/type for leases with secret info
	SecretLeaseType = "secret-lease"

	// AuthLeaseType - Bucket/type for leases with auth info
	AuthLeaseType = "auth-lease"

	// TokenType - Bucket/type for auto-auth tokens
	TokenType = "token"

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
	db, err := bolt.Open(dbPath, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, err
	}
	err = db.Update(func(tx *bolt.Tx) error {
		return createBoltSchema(tx)
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

func createBoltSchema(tx *bolt.Tx) error {
	// create the meta bucket at the top level
	meta, err := tx.CreateBucketIfNotExists([]byte(metaBucketName))
	if err != nil {
		return fmt.Errorf("failed to create bucket %s: %w", metaBucketName, err)
	}
	// check and set file version in the meta bucket
	version := meta.Get([]byte(storageVersionKey))
	switch {
	case version == nil:
		err = meta.Put([]byte(storageVersionKey), []byte(storageVersion))
		if err != nil {
			return fmt.Errorf("failed to set storage version: %w", err)
		}
	case string(version) != storageVersion:
		return fmt.Errorf("storage migration from %s to %s not implemented", string(version), storageVersion)
	}

	// create the buckets for tokens and leases
	_, err = tx.CreateBucketIfNotExists([]byte(TokenType))
	if err != nil {
		return fmt.Errorf("failed to create token bucket: %w", err)
	}
	_, err = tx.CreateBucketIfNotExists([]byte(AuthLeaseType))
	if err != nil {
		return fmt.Errorf("failed to create auth lease bucket: %w", err)
	}
	_, err = tx.CreateBucketIfNotExists([]byte(SecretLeaseType))
	if err != nil {
		return fmt.Errorf("failed to create secret lease bucket: %w", err)
	}

	return nil
}

// Set an index (token or lease) in bolt storage
func (b *BoltStorage) Set(ctx context.Context, id string, plaintext []byte, indexType string) error {
	blob, err := b.wrapper.Encrypt(ctx, plaintext, []byte(b.aad))
	if err != nil {
		return fmt.Errorf("error encrypting %s index: %w", indexType, err)
	}

	protoBlob, err := proto.Marshal(blob)
	if err != nil {
		return err
	}

	return b.db.Update(func(tx *bolt.Tx) error {
		s := tx.Bucket([]byte(indexType))
		if s == nil {
			return fmt.Errorf("bucket %q not found", indexType)
		}
		// If this is an auto-auth token, also stash it in the meta bucket for
		// easy retrieval upon restore
		if indexType == TokenType {
			meta := tx.Bucket([]byte(metaBucketName))
			if err := meta.Put([]byte(AutoAuthToken), protoBlob); err != nil {
				return fmt.Errorf("failed to set latest auto-auth token: %w", err)
			}
		}
		return s.Put([]byte(id), protoBlob)
	})
}

func getBucketIDs(b *bolt.Bucket) ([][]byte, error) {
	ids := [][]byte{}
	err := b.ForEach(func(k, v []byte) error {
		ids = append(ids, k)
		return nil
	})
	return ids, err
}

// Delete an index (token or lease) by id from bolt storage
func (b *BoltStorage) Delete(id string) error {
	return b.db.Update(func(tx *bolt.Tx) error {
		// Since Delete returns a nil error if the key doesn't exist, just call
		// delete in all three index buckets without checking existence first
		if err := tx.Bucket([]byte(TokenType)).Delete([]byte(id)); err != nil {
			return fmt.Errorf("failed to delete %q from token bucket: %w", id, err)
		}
		if err := tx.Bucket([]byte(AuthLeaseType)).Delete([]byte(id)); err != nil {
			return fmt.Errorf("failed to delete %q from auth lease bucket: %w", id, err)
		}
		if err := tx.Bucket([]byte(SecretLeaseType)).Delete([]byte(id)); err != nil {
			return fmt.Errorf("failed to delete %q from secret lease bucket: %w", id, err)
		}
		b.logger.Trace("deleted index from bolt db", "id", id)
		return nil
	})
}

func (b *BoltStorage) decrypt(ctx context.Context, ciphertext []byte) ([]byte, error) {
	var blob wrapping.EncryptedBlobInfo
	if err := proto.Unmarshal(ciphertext, &blob); err != nil {
		return nil, err
	}

	return b.wrapper.Decrypt(ctx, &blob, []byte(b.aad))
}

// GetByType returns a list of stored items of the specified type
func (b *BoltStorage) GetByType(ctx context.Context, indexType string) ([][]byte, error) {
	var returnBytes [][]byte

	err := b.db.View(func(tx *bolt.Tx) error {
		var errors *multierror.Error

		tx.Bucket([]byte(indexType)).ForEach(func(id, ciphertext []byte) error {
			plaintext, err := b.decrypt(ctx, ciphertext)
			if err != nil {
				errors = multierror.Append(errors, fmt.Errorf("error decrypting index id %s: %w", id, err))
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
		encryptedToken = meta.Get([]byte(AutoAuthToken))
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
		keyBucket := tx.Bucket([]byte(metaBucketName))
		if keyBucket == nil {
			return fmt.Errorf("bucket %q not found", metaBucketName)
		}
		token = keyBucket.Get([]byte(RetrievalTokenMaterial))
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
		for _, name := range []string{AuthLeaseType, SecretLeaseType, TokenType} {
			b.logger.Trace("deleting bolt bucket", "name", name)
			if err := tx.DeleteBucket([]byte(name)); err != nil {
				return err
			}
		}
		return createBoltSchema(tx)
	})
}

// DBFileExists checks whether the vault agent cache file at `filePath` exists
func DBFileExists(path string) (bool, error) {
	checkFile, err := os.OpenFile(filepath.Join(path, DatabaseFileName), os.O_RDWR, 0600)
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

// GetServiceAccountJWT reads the service account jwt from `tokenFile`. Default is
// the default service account file path in kubernetes.
func GetServiceAccountJWT(tokenFile string) (string, error) {
	if len(tokenFile) == 0 {
		tokenFile = "/var/run/secrets/kubernetes.io/serviceaccount/token"
	}
	token, err := ioutil.ReadFile(tokenFile)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(token)), nil
}
