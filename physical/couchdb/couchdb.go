package couchdb

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	metrics "github.com/armon/go-metrics"
	"github.com/hashicorp/errwrap"
	cleanhttp "github.com/hashicorp/go-cleanhttp"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/physical"
)

// CouchDBBackend allows the management of couchdb users
type CouchDBBackend struct {
	logger     log.Logger
	client     *couchDBClient
	permitPool *physical.PermitPool
}

// Verify CouchDBBackend satisfies the correct interfaces
var _ physical.Backend = (*CouchDBBackend)(nil)
var _ physical.PseudoTransactional = (*CouchDBBackend)(nil)
var _ physical.PseudoTransactional = (*TransactionalCouchDBBackend)(nil)

type couchDBClient struct {
	endpoint string
	username string
	password string
	*http.Client
}

type couchDBListItem struct {
	ID    string `json:"id"`
	Key   string `json:"key"`
	Value struct {
		Revision string
	} `json:"value"`
}

type couchDBList struct {
	TotalRows int               `json:"total_rows"`
	Offset    int               `json:"offset"`
	Rows      []couchDBListItem `json:"rows"`
}

func (m *couchDBClient) rev(key string) (string, error) {
	req, err := http.NewRequest("HEAD", fmt.Sprintf("%s/%s", m.endpoint, key), nil)
	if err != nil {
		return "", err
	}
	req.SetBasicAuth(m.username, m.password)

	resp, err := m.Client.Do(req)
	if err != nil {
		return "", err
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", nil
	}
	etag := resp.Header.Get("Etag")
	if len(etag) < 2 {
		return "", nil
	}
	return etag[1 : len(etag)-1], nil
}

func (m *couchDBClient) put(e couchDBEntry) error {
	bs, err := json.Marshal(e)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/%s", m.endpoint, e.ID), bytes.NewReader(bs))
	if err != nil {
		return err
	}
	req.SetBasicAuth(m.username, m.password)
	_, err = m.Client.Do(req)

	return err
}

func (m *couchDBClient) get(key string) (*physical.Entry, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", m.endpoint, url.PathEscape(key)), nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(m.username, m.password)
	resp, err := m.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	} else if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET returned %q", resp.Status)
	}
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	entry := couchDBEntry{}
	if err := json.Unmarshal(bs, &entry); err != nil {
		return nil, err
	}
	return entry.Entry, nil
}

func (m *couchDBClient) list(prefix string) ([]couchDBListItem, error) {
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/_all_docs", m.endpoint), nil)
	req.SetBasicAuth(m.username, m.password)
	values := req.URL.Query()
	values.Set("skip", "0")
	values.Set("include_docs", "false")
	if prefix != "" {
		values.Set("startkey", fmt.Sprintf("%q", prefix))
		values.Set("endkey", fmt.Sprintf("%q", prefix+"{}"))
	}
	req.URL.RawQuery = values.Encode()

	resp, err := m.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	results := couchDBList{}
	if err := json.Unmarshal(data, &results); err != nil {
		return nil, err
	}

	return results.Rows, nil
}

func buildCouchDBBackend(conf map[string]string, logger log.Logger) (*CouchDBBackend, error) {
	endpoint := os.Getenv("COUCHDB_ENDPOINT")
	if endpoint == "" {
		endpoint = conf["endpoint"]
	}
	if endpoint == "" {
		return nil, fmt.Errorf("missing endpoint")
	}

	username := os.Getenv("COUCHDB_USERNAME")
	if username == "" {
		username = conf["username"]
	}

	password := os.Getenv("COUCHDB_PASSWORD")
	if password == "" {
		password = conf["password"]
	}

	maxParStr, ok := conf["max_parallel"]
	var maxParInt int
	var err error
	if ok {
		maxParInt, err = strconv.Atoi(maxParStr)
		if err != nil {
			return nil, errwrap.Wrapf("failed parsing max_parallel parameter: {{err}}", err)
		}
		if logger.IsDebug() {
			logger.Debug("max_parallel set", "max_parallel", maxParInt)
		}
	}

	return &CouchDBBackend{
		client: &couchDBClient{
			endpoint: endpoint,
			username: username,
			password: password,
			Client:   cleanhttp.DefaultPooledClient(),
		},
		logger:     logger,
		permitPool: physical.NewPermitPool(maxParInt),
	}, nil
}

func NewCouchDBBackend(conf map[string]string, logger log.Logger) (physical.Backend, error) {
	return buildCouchDBBackend(conf, logger)
}

type couchDBEntry struct {
	Entry   *physical.Entry `json:"entry"`
	Rev     string          `json:"_rev,omitempty"`
	ID      string          `json:"_id"`
	Deleted *bool           `json:"_deleted,omitempty"`
}

// Put is used to insert or update an entry
func (m *CouchDBBackend) Put(ctx context.Context, entry *physical.Entry) error {
	m.permitPool.Acquire()
	defer m.permitPool.Release()

	return m.PutInternal(ctx, entry)
}

// Get is used to fetch an entry
func (m *CouchDBBackend) Get(ctx context.Context, key string) (*physical.Entry, error) {
	m.permitPool.Acquire()
	defer m.permitPool.Release()

	return m.GetInternal(ctx, key)
}

// Delete is used to permanently delete an entry
func (m *CouchDBBackend) Delete(ctx context.Context, key string) error {
	m.permitPool.Acquire()
	defer m.permitPool.Release()

	return m.DeleteInternal(ctx, key)
}

// List is used to list all the keys under a given prefix
func (m *CouchDBBackend) List(ctx context.Context, prefix string) ([]string, error) {
	defer metrics.MeasureSince([]string{"couchdb", "list"}, time.Now())

	m.permitPool.Acquire()
	defer m.permitPool.Release()

	items, err := m.client.list(prefix)
	if err != nil {
		return nil, err
	}

	var out []string
	seen := make(map[string]interface{})
	for _, result := range items {
		trimmed := strings.TrimPrefix(result.ID, prefix)
		sep := strings.Index(trimmed, "/")
		if sep == -1 {
			out = append(out, trimmed)
		} else {
			trimmed = trimmed[:sep+1]
			if _, ok := seen[trimmed]; !ok {
				out = append(out, trimmed)
				seen[trimmed] = struct{}{}
			}
		}
	}
	return out, nil
}

// TransactionalCouchDBBackend creates a couchdb backend that forces all operations to happen
// in serial
type TransactionalCouchDBBackend struct {
	CouchDBBackend
}

func NewTransactionalCouchDBBackend(conf map[string]string, logger log.Logger) (physical.Backend, error) {
	backend, err := buildCouchDBBackend(conf, logger)
	if err != nil {
		return nil, err
	}
	backend.permitPool = physical.NewPermitPool(1)

	return &TransactionalCouchDBBackend{
		CouchDBBackend: *backend,
	}, nil
}

// GetInternal is used to fetch an entry
func (m *CouchDBBackend) GetInternal(ctx context.Context, key string) (*physical.Entry, error) {
	defer metrics.MeasureSince([]string{"couchdb", "get"}, time.Now())

	return m.client.get(key)
}

// PutInternal is used to insert or update an entry
func (m *CouchDBBackend) PutInternal(ctx context.Context, entry *physical.Entry) error {
	defer metrics.MeasureSince([]string{"couchdb", "put"}, time.Now())

	revision, _ := m.client.rev(url.PathEscape(entry.Key))

	return m.client.put(couchDBEntry{
		Entry: entry,
		Rev:   revision,
		ID:    url.PathEscape(entry.Key),
	})
}

// DeleteInternal is used to permanently delete an entry
func (m *CouchDBBackend) DeleteInternal(ctx context.Context, key string) error {
	defer metrics.MeasureSince([]string{"couchdb", "delete"}, time.Now())

	revision, _ := m.client.rev(url.PathEscape(key))
	deleted := true
	return m.client.put(couchDBEntry{
		ID:      url.PathEscape(key),
		Rev:     revision,
		Deleted: &deleted,
	})
}
