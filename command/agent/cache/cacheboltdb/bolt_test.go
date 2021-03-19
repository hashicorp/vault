package cacheboltdb

import (
	"context"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/command/agent/cache/keymanager"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getTestKeyManager(t *testing.T) keymanager.KeyManager {
	t.Helper()

	km, err := keymanager.NewPassthroughKeyManager(nil)
	require.NoError(t, err)

	return km
}

func TestBolt_SetGet(t *testing.T) {
	ctx := context.Background()

	path, err := ioutil.TempDir("", "bolt-test")
	require.NoError(t, err)
	defer os.RemoveAll(path)

	b, err := NewBoltStorage(&BoltStorageConfig{
		Path:    path,
		Logger:  hclog.Default(),
		Wrapper: getTestKeyManager(t).Wrapper(),
	})
	require.NoError(t, err)

	secrets, err := b.GetByType(ctx, SecretLeaseType)
	assert.NoError(t, err)
	require.Len(t, secrets, 0)

	err = b.Set(ctx, "test1", []byte("hello"), SecretLeaseType)
	assert.NoError(t, err)
	secrets, err = b.GetByType(ctx, SecretLeaseType)
	assert.NoError(t, err)
	require.Len(t, secrets, 1)
	assert.Equal(t, []byte("hello"), secrets[0])
}

func TestBoltDelete(t *testing.T) {
	ctx := context.Background()

	path, err := ioutil.TempDir("", "bolt-test")
	require.NoError(t, err)
	defer os.RemoveAll(path)

	b, err := NewBoltStorage(&BoltStorageConfig{
		Path:    path,
		Logger:  hclog.Default(),
		Wrapper: getTestKeyManager(t).Wrapper(),
	})
	require.NoError(t, err)

	err = b.Set(ctx, "secret-test1", []byte("hello1"), SecretLeaseType)
	require.NoError(t, err)
	err = b.Set(ctx, "secret-test2", []byte("hello2"), SecretLeaseType)
	require.NoError(t, err)

	secrets, err := b.GetByType(ctx, SecretLeaseType)
	require.NoError(t, err)
	assert.Len(t, secrets, 2)
	assert.ElementsMatch(t, [][]byte{[]byte("hello1"), []byte("hello2")}, secrets)

	err = b.Delete("secret-test1")
	require.NoError(t, err)
	secrets, err = b.GetByType(ctx, SecretLeaseType)
	require.NoError(t, err)
	require.Len(t, secrets, 1)
	assert.Equal(t, []byte("hello2"), secrets[0])
}

func TestBoltClear(t *testing.T) {
	ctx := context.Background()

	path, err := ioutil.TempDir("", "bolt-test")
	require.NoError(t, err)
	defer os.RemoveAll(path)

	b, err := NewBoltStorage(&BoltStorageConfig{
		Path:    path,
		Logger:  hclog.Default(),
		Wrapper: getTestKeyManager(t).Wrapper(),
	})
	require.NoError(t, err)

	// Populate the bolt db
	err = b.Set(ctx, "secret-test1", []byte("hello"), SecretLeaseType)
	require.NoError(t, err)
	secrets, err := b.GetByType(ctx, SecretLeaseType)
	require.NoError(t, err)
	require.Len(t, secrets, 1)
	assert.Equal(t, []byte("hello"), secrets[0])

	err = b.Set(ctx, "auth-test1", []byte("hello"), AuthLeaseType)
	require.NoError(t, err)
	auths, err := b.GetByType(ctx, AuthLeaseType)
	require.NoError(t, err)
	require.Len(t, auths, 1)
	assert.Equal(t, []byte("hello"), auths[0])

	err = b.Set(ctx, "token-test1", []byte("hello"), TokenType)
	require.NoError(t, err)
	tokens, err := b.GetByType(ctx, TokenType)
	require.NoError(t, err)
	require.Len(t, tokens, 1)
	assert.Equal(t, []byte("hello"), tokens[0])

	// Clear the bolt db, and check that it's indeed clear
	err = b.Clear()
	require.NoError(t, err)
	secrets, err = b.GetByType(ctx, SecretLeaseType)
	require.NoError(t, err)
	assert.Len(t, secrets, 0)
	auths, err = b.GetByType(ctx, AuthLeaseType)
	require.NoError(t, err)
	assert.Len(t, auths, 0)
	tokens, err = b.GetByType(ctx, TokenType)
	require.NoError(t, err)
	assert.Len(t, tokens, 0)
}

func TestBoltSetAutoAuthToken(t *testing.T) {
	ctx := context.Background()

	path, err := ioutil.TempDir("", "bolt-test")
	require.NoError(t, err)
	defer os.RemoveAll(path)

	b, err := NewBoltStorage(&BoltStorageConfig{
		Path:    path,
		Logger:  hclog.Default(),
		Wrapper: getTestKeyManager(t).Wrapper(),
	})
	require.NoError(t, err)

	token, err := b.GetAutoAuthToken(ctx)
	assert.NoError(t, err)
	assert.Nil(t, token)

	// set first token
	err = b.Set(ctx, "token-test1", []byte("hello 1"), TokenType)
	require.NoError(t, err)
	secrets, err := b.GetByType(ctx, TokenType)
	require.NoError(t, err)
	require.Len(t, secrets, 1)
	assert.Equal(t, []byte("hello 1"), secrets[0])
	token, err = b.GetAutoAuthToken(ctx)
	assert.NoError(t, err)
	assert.Equal(t, []byte("hello 1"), token)

	// set second token
	err = b.Set(ctx, "token-test2", []byte("hello 2"), TokenType)
	require.NoError(t, err)
	secrets, err = b.GetByType(ctx, TokenType)
	require.NoError(t, err)
	require.Len(t, secrets, 2)
	assert.ElementsMatch(t, [][]byte{[]byte("hello 1"), []byte("hello 2")}, secrets)
	token, err = b.GetAutoAuthToken(ctx)
	assert.NoError(t, err)
	assert.Equal(t, []byte("hello 2"), token)
}

func TestDBFileExists(t *testing.T) {
	testCases := []struct {
		name        string
		mkDir       bool
		createFile  bool
		expectExist bool
	}{
		{
			name:        "all exists",
			mkDir:       true,
			createFile:  true,
			expectExist: true,
		},
		{
			name:        "dir exist, file missing",
			mkDir:       true,
			createFile:  false,
			expectExist: false,
		},
		{
			name:        "all missing",
			mkDir:       false,
			createFile:  false,
			expectExist: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var tmpPath string
			var err error
			if tc.mkDir {
				tmpPath, err = ioutil.TempDir("", "test-db-path")
				require.NoError(t, err)
			}
			if tc.createFile {
				err = ioutil.WriteFile(path.Join(tmpPath, DatabaseFileName), []byte("test-db-path"), 0600)
				require.NoError(t, err)
			}
			exists, err := DBFileExists(tmpPath)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectExist, exists)
		})
	}

}

func Test_SetGetRetrievalToken(t *testing.T) {
	testCases := []struct {
		name          string
		tokenToSet    []byte
		expectedToken []byte
	}{
		{
			name:          "normal set and get",
			tokenToSet:    []byte("test token"),
			expectedToken: []byte("test token"),
		},
		{
			name:          "no token set",
			tokenToSet:    nil,
			expectedToken: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			path, err := ioutil.TempDir("", "bolt-test")
			require.NoError(t, err)
			defer os.RemoveAll(path)

			b, err := NewBoltStorage(&BoltStorageConfig{
				Path:    path,
				Logger:  hclog.Default(),
				Wrapper: getTestKeyManager(t).Wrapper(),
			})
			require.NoError(t, err)
			defer b.Close()

			if tc.tokenToSet != nil {
				err := b.StoreRetrievalToken(tc.tokenToSet)
				require.NoError(t, err)
			}
			gotKey, err := b.GetRetrievalToken()
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedToken, gotKey)
		})
	}
}
