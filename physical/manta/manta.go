package manta

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"
	"time"

	metrics "github.com/armon/go-metrics"
	"github.com/hashicorp/errwrap"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/physical"
	triton "github.com/joyent/triton-go"
	"github.com/joyent/triton-go/authentication"
	"github.com/joyent/triton-go/errors"
	"github.com/joyent/triton-go/storage"
)

const mantaDefaultRootStore = "/stor"

type MantaBackend struct {
	logger     log.Logger
	permitPool *physical.PermitPool
	client     *storage.StorageClient
	directory  string
}

func NewMantaBackend(conf map[string]string, logger log.Logger) (physical.Backend, error) {
	user := os.Getenv("MANTA_USER")
	if user == "" {
		user = conf["user"]
	}

	keyId := os.Getenv("MANTA_KEY_ID")
	if keyId == "" {
		keyId = conf["key_id"]
	}

	url := os.Getenv("MANTA_URL")
	if url == "" {
		url = conf["url"]
	} else {
		url = "https://us-east.manta.joyent.com"
	}

	subuser := os.Getenv("MANTA_SUBUSER")
	if subuser == "" {
		if confUser, ok := conf["subuser"]; ok {
			subuser = confUser
		}
	}

	input := authentication.SSHAgentSignerInput{
		KeyID:       keyId,
		AccountName: user,
		Username:    subuser,
	}
	signer, err := authentication.NewSSHAgentSigner(input)
	if err != nil {
		return nil, errwrap.Wrapf("Error Creating SSH Agent Signer: {{err}}", err)
	}

	maxParStr, ok := conf["max_parallel"]
	var maxParInt int
	if ok {
		maxParInt, err = strconv.Atoi(maxParStr)
		if err != nil {
			return nil, errwrap.Wrapf("failed parsing max_parallel parameter: {{err}}", err)
		}
		if logger.IsDebug() {
			logger.Debug("max_parallel set", "max_parallel", maxParInt)
		}
	}

	config := &triton.ClientConfig{
		MantaURL:    url,
		AccountName: user,
		Signers:     []authentication.Signer{signer},
	}

	client, err := storage.NewClient(config)
	if err != nil {
		return nil, errwrap.Wrapf("failed initialising Storage client: {{err}}", err)
	}

	return &MantaBackend{
		client:     client,
		directory:  conf["directory"],
		logger:     logger,
		permitPool: physical.NewPermitPool(maxParInt),
	}, nil
}

// Put is used to insert or update an entry
func (m *MantaBackend) Put(ctx context.Context, entry *physical.Entry) error {
	defer metrics.MeasureSince([]string{"manta", "put"}, time.Now())

	m.permitPool.Acquire()
	defer m.permitPool.Release()

	r := bytes.NewReader(entry.Value)
	r.Seek(0, 0)

	return m.client.Objects().Put(ctx, &storage.PutObjectInput{
		ObjectPath:    path.Join(mantaDefaultRootStore, m.directory, entry.Key, ".vault_value"),
		ObjectReader:  r,
		ContentLength: uint64(len(entry.Value)),
		ForceInsert:   true,
	})
}

// Get is used to fetch an entry
func (m *MantaBackend) Get(ctx context.Context, key string) (*physical.Entry, error) {
	defer metrics.MeasureSince([]string{"manta", "get"}, time.Now())

	m.permitPool.Acquire()
	defer m.permitPool.Release()

	output, err := m.client.Objects().Get(ctx, &storage.GetObjectInput{
		ObjectPath: path.Join(mantaDefaultRootStore, m.directory, key, ".vault_value"),
	})
	if err != nil {
		if strings.Contains(err.Error(), "ResourceNotFound") {
			return nil, nil
		}
		return nil, err
	}

	defer output.ObjectReader.Close()

	data := make([]byte, output.ContentLength)
	_, err = io.ReadFull(output.ObjectReader, data)
	if err != nil {
		return nil, err
	}

	ent := &physical.Entry{
		Key:   key,
		Value: data,
	}

	return ent, nil
}

// Delete is used to permanently delete an entry
func (m *MantaBackend) Delete(ctx context.Context, key string) error {
	defer metrics.MeasureSince([]string{"manta", "delete"}, time.Now())

	m.permitPool.Acquire()
	defer m.permitPool.Release()

	if strings.HasSuffix(key, "/") {
		err := m.client.Dir().Delete(ctx, &storage.DeleteDirectoryInput{
			DirectoryName: path.Join(mantaDefaultRootStore, m.directory, key),
			ForceDelete:   true,
		})
		if err != nil {
			return err
		}
	} else {
		err := m.client.Objects().Delete(ctx, &storage.DeleteObjectInput{
			ObjectPath: path.Join(mantaDefaultRootStore, m.directory, key, ".vault_value"),
		})
		if err != nil {
			if errors.IsResourceNotFound(err) {
				return nil
			}
			return err
		}

		return tryDeleteDirectory(ctx, m, path.Join(mantaDefaultRootStore, m.directory, key))
	}

	return nil
}

func tryDeleteDirectory(ctx context.Context, m *MantaBackend, directoryPath string) error {
	objs, err := m.client.Dir().List(ctx, &storage.ListDirectoryInput{
		DirectoryName: directoryPath,
	})
	if err != nil {
		if errors.IsResourceNotFound(err) {
			return nil
		}
		return err
	}
	if objs != nil && len(objs.Entries) == 0 {
		err := m.client.Dir().Delete(ctx, &storage.DeleteDirectoryInput{
			DirectoryName: directoryPath,
		})
		if err != nil {
			return err
		}

		return tryDeleteDirectory(ctx, m, path.Dir(directoryPath))
	}
	return nil
}

// List is used to list all the keys under a given
// prefix, up to the next prefix.
func (m *MantaBackend) List(ctx context.Context, prefix string) ([]string, error) {
	defer metrics.MeasureSince([]string{"manta", "list"}, time.Now())

	m.permitPool.Acquire()
	defer m.permitPool.Release()

	objs, err := m.client.Dir().List(ctx, &storage.ListDirectoryInput{
		DirectoryName: path.Join(mantaDefaultRootStore, m.directory, prefix),
	})
	if err != nil {
		if errors.IsResourceNotFound(err) {
			return []string{}, nil
		}
		return nil, err
	}

	keys := []string{}
	for _, obj := range objs.Entries {
		if obj.Type == "directory" {
			objs, err := m.client.Dir().List(ctx, &storage.ListDirectoryInput{
				DirectoryName: path.Join(mantaDefaultRootStore, m.directory, prefix, obj.Name),
			})
			if err != nil {
				if !errors.IsResourceNotFound(err) {
					return nil, err
				}
			}

			//We need to check to see if there is something more than just the `value` file
			//if the length of the children is:
			// > 1 and includes the value `index` then we need to add foo and foo/
			// = 1 and the value is `index` then we need to add foo
			// = 1 and the value is not `index` then we need to add foo/
			if len(objs.Entries) == 1 {
				if objs.Entries[0].Name != ".vault_value" {
					keys = append(keys, fmt.Sprintf("%s/", obj.Name))
				} else {
					keys = append(keys, obj.Name)
				}
			} else if len(objs.Entries) > 1 {
				for _, childObj := range objs.Entries {
					if childObj.Name == ".vault_value" {
						keys = append(keys, obj.Name)
					} else {
						keys = append(keys, fmt.Sprintf("%s/", obj.Name))
					}
				}
			} else {
				keys = append(keys, obj.Name)
			}
		}
	}

	sort.Strings(keys)

	return keys, nil
}
