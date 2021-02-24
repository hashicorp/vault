package cacheboltdb

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/vault/command/agent/cache/crypto"
	bolt "go.etcd.io/bbolt"
)

const (
	// Keep track of schema version for future migrations
	storageVersionKey = "version"
	storageVersion    = "1"

	// DatabaseFileName - filename for the persistent cache file
	DatabaseFileName = "vault-agent-cache.db"

	// rootBucketName - naming the root bucket
	rootBucketName = "root"

	// SecretLeaseType - Bucket/type for leases with secret info
	SecretLeaseType = "secret-lease"

	// AuthLeaseType - Bucket/type for leases with auth info
	AuthLeaseType = "auth-lease"

	// TokenType - Bucket/type for auto-auth tokens
	TokenType = "token"

	// AutoAuthToken - key for the latest auto-auth token
	AutoAuthToken = "auto-auth-token"

	// RetrievalTokenBucket names the bucket that stores the retrieval material
	// needed for unlocking (decrypting) the persisted tokens and leases
	RetrievalTokenBucket = "retrieval-token"

	// RetrievalTokenMaterial is the actual key or token in the key bucket
	RetrievalTokenMaterial = "retrieval-token-material"
)

// BoltStorage is a persistent cache using a bolt db. Items are organized with
// the encryption key as the top-level bucket, and then leases and tokens are
// stored in sub buckets.
type BoltStorage struct {
	db         *bolt.DB
	rootBucket string
	logger     hclog.Logger
	keyMgr     crypto.KeyManager
	aad        string
}

// BoltStorageConfig is the collection of input parameters for setting up bolt
// storage
type BoltStorageConfig struct {
	Path       string
	Logger     hclog.Logger
	AAD        string
	KeyManager crypto.KeyManager
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
		db:         db,
		rootBucket: rootBucketName,
		logger:     config.Logger,
		keyMgr:     config.KeyManager,
		aad:        config.AAD,
	}
	if config.KeyManager != nil {
		bs.SetKeyMgr(config.KeyManager)
	}
	return bs, nil
}

func createBoltSchema(tx *bolt.Tx) error {
	_, err := tx.CreateBucketIfNotExists([]byte(RetrievalTokenBucket))
	if err != nil {
		return fmt.Errorf("failed to create key bucket: %w", err)
	}
	root, err := tx.CreateBucketIfNotExists([]byte(rootBucketName))
	if err != nil {
		return fmt.Errorf("failed to create bucket %s: %w", rootBucketName, err)
	}
	_, err = root.CreateBucketIfNotExists([]byte(TokenType))
	if err != nil {
		return fmt.Errorf("failed to create token sub-bucket: %w", err)
	}
	_, err = root.CreateBucketIfNotExists([]byte(AuthLeaseType))
	if err != nil {
		return fmt.Errorf("failed to create auth lease sub-bucket: %w", err)
	}
	_, err = root.CreateBucketIfNotExists([]byte(SecretLeaseType))
	if err != nil {
		return fmt.Errorf("failed to create secret lease sub-bucket: %w", err)
	}

	// check and set file version in the root bucket
	version := root.Get([]byte(storageVersionKey))
	switch {
	case version == nil:
		err = root.Put([]byte(storageVersionKey), []byte(storageVersion))
		if err != nil {
			return fmt.Errorf("failed to set storage version: %w", err)
		}
	case string(version) != storageVersion:
		return fmt.Errorf("storage migration from %s to %s not implemented", string(version), storageVersion)
	}
	return nil
}

// SetKeyMgr sets the encryption for a bolt storage
func (b *BoltStorage) SetKeyMgr(km crypto.KeyManager) {
	b.keyMgr = km
}

// Set an index in bolt storage
func (b *BoltStorage) Set(ctx context.Context, id string, plainText []byte, indexType string) error {
	cipherText, err := b.keyMgr.Encrypt(ctx, plainText, []byte(b.aad))
	if err != nil {
		return fmt.Errorf("error encrypting %s index: %w", indexType, err)
	}

	return b.db.Update(func(tx *bolt.Tx) error {
		top := tx.Bucket([]byte(b.rootBucket))
		if top == nil {
			return fmt.Errorf("bucket %q not found", b.rootBucket)
		}
		s := top.Bucket([]byte(indexType))
		if s == nil {
			return fmt.Errorf("bucket %q not found", indexType)
		}
		// If this is an auto-auth token, also stash it in the root bucket for
		// easy retrieval upon restore
		if indexType == TokenType {
			if err := top.Put([]byte(AutoAuthToken), cipherText); err != nil {
				return fmt.Errorf("failed to set latest auto-auth token: %w", err)
			}
		}
		return s.Put([]byte(id), cipherText)
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

// Delete an index by id from bolt storage
func (b *BoltStorage) Delete(id string) error {
	return b.db.Update(func(tx *bolt.Tx) error {
		top := tx.Bucket([]byte(b.rootBucket))
		if top == nil {
			return fmt.Errorf("bucket %q not found", b.rootBucket)
		}
		// Since Delete returns a nil error if the key doesn't exist, just call
		// delete in all three sub-buckets without checking existence first
		if err := top.Bucket([]byte(TokenType)).Delete([]byte(id)); err != nil {
			return fmt.Errorf("failed to delete %q from token bucket: %w", id, err)
		}
		if err := top.Bucket([]byte(AuthLeaseType)).Delete([]byte(id)); err != nil {
			return fmt.Errorf("failed to delete %q from auth lease bucket: %w", id, err)
		}
		if err := top.Bucket([]byte(SecretLeaseType)).Delete([]byte(id)); err != nil {
			return fmt.Errorf("failed to delete %q from secret lease bucket: %w", id, err)
		}
		b.logger.Trace("deleted index from bolt db", "id", id)
		return nil
	})
}

// GetByType returns a list of stored items of the specified type
func (b *BoltStorage) GetByType(ctx context.Context, indexType string) ([][]byte, error) {
	returnBytes := [][]byte{}

	err := b.db.View(func(tx *bolt.Tx) error {
		var errors *multierror.Error

		top := tx.Bucket([]byte(b.rootBucket))
		if top == nil {
			return fmt.Errorf("bucket %q not found", b.rootBucket)
		}
		top.Bucket([]byte(indexType)).ForEach(func(id, cipherText []byte) error {
			plainText, err := b.keyMgr.Decrypt(ctx, cipherText, []byte(b.aad))
			if err != nil {
				errors = multierror.Append(errors, fmt.Errorf("error decrypting index id %s: %w", id, err))
				return nil
			}
			returnBytes = append(returnBytes, plainText)
			return nil
		})
		return errors.ErrorOrNil()
	})

	return returnBytes, err
}

// GetAutoAuthToken retrieves the latest auto-auth token, and returns nil if non
// exists yet
func (b *BoltStorage) GetAutoAuthToken(ctx context.Context) ([]byte, error) {
	token := []byte{}

	err := b.db.View(func(tx *bolt.Tx) error {
		top := tx.Bucket([]byte(b.rootBucket))
		if top == nil {
			return fmt.Errorf("bucket %q not found", b.rootBucket)
		}
		token = top.Get([]byte(AutoAuthToken))
		return nil
	})
	if err != nil {
		return nil, err
	}

	if token == nil {
		return nil, nil
	}

	plainText, err := b.keyMgr.Decrypt(ctx, token, []byte(b.aad))
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt auto-auth token: %w", err)
	}
	return plainText, nil
}

// GetRetrievalToken retrieves a plaintext token from the KeyBucket, which will
// be used by the key manager to retrieve the encryption key, nil if none set
func (b *BoltStorage) GetRetrievalToken() ([]byte, error) {
	token := []byte{}

	err := b.db.View(func(tx *bolt.Tx) error {
		keyBucket := tx.Bucket([]byte(RetrievalTokenBucket))
		if keyBucket == nil {
			return fmt.Errorf("bucket %q not found", RetrievalTokenBucket)
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
		bucket := tx.Bucket([]byte(RetrievalTokenBucket))
		if bucket == nil {
			return fmt.Errorf("bucket %q not found", RetrievalTokenBucket)
		}
		return bucket.Put([]byte(RetrievalTokenMaterial), token)
	})
}

// Close the boltdb
func (b *BoltStorage) Close() error {
	b.logger.Trace("closing bolt db", "path", b.db.Path())
	return b.db.Close()
}

// Clear the boltdb by deleting all the root buckets and recreating the
// schema/layout
func (b *BoltStorage) Clear() error {
	return b.db.Update(func(tx *bolt.Tx) error {
		err := tx.ForEach(func(name []byte, bucket *bolt.Bucket) error {
			b.logger.Trace("deleting bolt bucket", "name", name)
			if err := tx.DeleteBucket(name); err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			return err
		}
		return createBoltSchema(tx)
	})
}

// DBFileExists checks whether the vault agent cache file at `filePath` exists
func DBFileExists(filePath string) (bool, error) {
	checkFile, err := os.OpenFile(path.Join(filePath, DatabaseFileName), os.O_RDWR, 0600)
	defer checkFile.Close()
	switch {
	case err == nil:
		return true, nil
	case os.IsNotExist(err):
		return false, nil
	default:
		return false, fmt.Errorf("failed to check if bolt file exists at path %s: %w", filePath, err)
	}
}

func GetServiceAccountJWT(path string) (string, error) {
	if len(path) == 0 {
		path = "/var/run/secrets/kubernetes.io/serviceaccount/token"
	}
	token, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(token)), nil
}
