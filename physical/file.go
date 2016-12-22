package physical

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	log "github.com/mgutz/logxi/v1"

	"github.com/hashicorp/vault/helper/jsonutil"
)

// FileBackend is a physical backend that stores data on disk
// at a given file path. It can be used for durable single server
// situations, or to develop locally where durability is not critical.
//
// WARNING: the file backend implementation is currently extremely unsafe
// and non-performant. It is meant mostly for local testing and development.
// It can be improved in the future.
type FileBackend struct {
	Path   string
	l      sync.Mutex
	logger log.Logger
}

// newFileBackend constructs a Filebackend using the given directory
func newFileBackend(conf map[string]string, logger log.Logger) (Backend, error) {
	path, ok := conf["path"]
	if !ok {
		return nil, fmt.Errorf("'path' must be set")
	}

	return &FileBackend{
		Path:   path,
		logger: logger,
	}, nil
}

func (b *FileBackend) Delete(path string) error {
	if path == "" {
		return nil
	}

	b.l.Lock()
	defer b.l.Unlock()

	basePath, key := b.path(path)
	fullPath := filepath.Join(basePath, "_"+key)

	err := os.Remove(fullPath)

	// if the file does not exist and if the key is base64 URL encoded, try to delete the file with decoded key
	if err != nil && os.IsNotExist(err) {
		keyDecodedBytes, err := base64.URLEncoding.DecodeString(key)
		if err == nil {
			fullPath = filepath.Join(basePath, "_"+string(keyDecodedBytes))
			err = os.Remove(fullPath)
		}
	}

	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("Failed to remove %q: %v", fullPath, err)
	}

	err = b.cleanupLogicalPath(path)

	return err
}

// cleanupLogicalPath is used to remove all empty nodes, begining with deepest
// one, aborting on first non-empty one, up to top-level node.
func (b *FileBackend) cleanupLogicalPath(path string) error {
	nodes := strings.Split(path, fmt.Sprintf("%c", os.PathSeparator))
	for i := len(nodes) - 1; i > 0; i-- {
		fullPath := filepath.Join(b.Path, filepath.Join(nodes[:i]...))

		dir, err := os.Open(fullPath)
		if err != nil {
			if os.IsNotExist(err) {
				return nil
			} else {
				return err
			}
		}

		list, err := dir.Readdir(1)
		dir.Close()
		if err != nil && err != io.EOF {
			return err
		}

		// If we have no entries, it's an empty directory; remove it
		if err == io.EOF || list == nil || len(list) == 0 {
			err = os.Remove(fullPath)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (b *FileBackend) Get(k string) (*Entry, error) {
	b.l.Lock()
	defer b.l.Unlock()

	path, key := b.path(k)
	fullPath := filepath.Join(path, "_"+key)

	f, err := os.Open(fullPath)

	if err != nil && os.IsNotExist(err) {
		fullPath = filepath.Join(path, "_"+base64.URLEncoding.EncodeToString([]byte(key)))
		f, err = os.Open(fullPath)
	}

	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	defer f.Close()

	var entry Entry
	if err := jsonutil.DecodeJSONFromReader(f, &entry); err != nil {
		return nil, err
	}

	return &entry, nil
}

func (b *FileBackend) Put(entry *Entry) error {
	path, key := b.path(entry.Key)

	key = base64.URLEncoding.EncodeToString([]byte(key))

	b.l.Lock()
	defer b.l.Unlock()

	// Make the parent tree
	if err := os.MkdirAll(path, 0755); err != nil {
		return err
	}

	fullPath := filepath.Join(path, "_"+key)

	// JSON encode the entry and write it
	f, err := os.OpenFile(
		fullPath,
		os.O_CREATE|os.O_TRUNC|os.O_WRONLY,
		0600)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	return enc.Encode(entry)
}

func (b *FileBackend) List(prefix string) ([]string, error) {
	b.l.Lock()
	defer b.l.Unlock()

	path := b.Path
	if prefix != "" {
		path = filepath.Join(path, prefix)
	}

	// Read the directory contents
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}

		return nil, err
	}
	defer f.Close()

	names, err := f.Readdirnames(-1)
	if err != nil {
		return nil, err
	}

	for i, name := range names {
		if name[0] == '_' {
			names[i] = name[1:]
			// TODO: Decode the name
		} else {
			names[i] = name + "/"
		}
	}

	return names, nil
}

func (b *FileBackend) path(k string) (string, string) {
	path := filepath.Join(b.Path, k)
	key := filepath.Base(path)
	path = filepath.Dir(path)
	return path, key
}
