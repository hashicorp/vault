package swift

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	log "github.com/hashicorp/go-hclog"

	metrics "github.com/armon/go-metrics"
	"github.com/hashicorp/errwrap"
	cleanhttp "github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/vault/helper/strutil"
	"github.com/hashicorp/vault/physical"
	"github.com/ncw/swift"
)

// Verify SwiftBackend satisfies the correct interfaces
var _ physical.Backend = (*SwiftBackend)(nil)

// SwiftBackend is a physical backend that stores data
// within an OpenStack Swift container.
type SwiftBackend struct {
	container  string
	client     *swift.Connection
	logger     log.Logger
	permitPool *physical.PermitPool
}

// NewSwiftBackend constructs a Swift backend using a pre-existing
// container. Credentials can be provided to the backend, sourced
// from the environment.
func NewSwiftBackend(conf map[string]string, logger log.Logger) (physical.Backend, error) {
	var ok bool

	username := os.Getenv("OS_USERNAME")
	if username == "" {
		username = conf["username"]
		if username == "" {
			return nil, fmt.Errorf("missing username")
		}
	}
	password := os.Getenv("OS_PASSWORD")
	if password == "" {
		password = conf["password"]
		if password == "" {
			return nil, fmt.Errorf("missing password")
		}
	}
	authUrl := os.Getenv("OS_AUTH_URL")
	if authUrl == "" {
		authUrl = conf["auth_url"]
		if authUrl == "" {
			return nil, fmt.Errorf("missing auth_url")
		}
	}
	container := os.Getenv("OS_CONTAINER")
	if container == "" {
		container = conf["container"]
		if container == "" {
			return nil, fmt.Errorf("missing container")
		}
	}
	project := os.Getenv("OS_PROJECT_NAME")
	if project == "" {
		if project, ok = conf["project"]; !ok {
			// Check for KeyStone naming prior to V3
			project = os.Getenv("OS_TENANT_NAME")
			if project == "" {
				project = conf["tenant"]
			}
		}
	}

	domain := os.Getenv("OS_USER_DOMAIN_NAME")
	if domain == "" {
		domain = conf["domain"]
	}
	projectDomain := os.Getenv("OS_PROJECT_DOMAIN_NAME")
	if projectDomain == "" {
		projectDomain = conf["project-domain"]
	}

	region := os.Getenv("OS_REGION_NAME")
	if region == "" {
		region = conf["region"]
	}
	tenantID := os.Getenv("OS_TENANT_ID")
	if tenantID == "" {
		tenantID = conf["tenant_id"]
	}
	trustID := os.Getenv("OS_TRUST_ID")
	if trustID == "" {
		trustID = conf["trust_id"]
	}
	storageUrl := os.Getenv("OS_STORAGE_URL")
	if storageUrl == "" {
		storageUrl = conf["storage_url"]
	}
	authToken := os.Getenv("OS_AUTH_TOKEN")
	if authToken == "" {
		authToken = conf["auth_token"]
	}

	c := swift.Connection{
		Domain:       domain,
		UserName:     username,
		ApiKey:       password,
		AuthUrl:      authUrl,
		Tenant:       project,
		TenantDomain: projectDomain,
		Region:       region,
		TenantId:     tenantID,
		TrustId:      trustID,
		StorageUrl:   storageUrl,
		AuthToken:    authToken,
		Transport:    cleanhttp.DefaultPooledTransport(),
	}

	err := c.Authenticate()
	if err != nil {
		return nil, err
	}

	_, _, err = c.Container(container)
	if err != nil {
		return nil, errwrap.Wrapf(fmt.Sprintf("Unable to access container %q: {{err}}", container), err)
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

	s := &SwiftBackend{
		client:     &c,
		container:  container,
		logger:     logger,
		permitPool: physical.NewPermitPool(maxParInt),
	}
	return s, nil
}

// Put is used to insert or update an entry
func (s *SwiftBackend) Put(ctx context.Context, entry *physical.Entry) error {
	defer metrics.MeasureSince([]string{"swift", "put"}, time.Now())

	s.permitPool.Acquire()
	defer s.permitPool.Release()

	err := s.client.ObjectPutBytes(s.container, entry.Key, entry.Value, "")

	if err != nil {
		return err
	}

	return nil
}

// Get is used to fetch an entry
func (s *SwiftBackend) Get(ctx context.Context, key string) (*physical.Entry, error) {
	defer metrics.MeasureSince([]string{"swift", "get"}, time.Now())

	s.permitPool.Acquire()
	defer s.permitPool.Release()

	//Do a list of names with the key first since eventual consistency means
	//it might be deleted, but a node might return a read of bytes which fails
	//the physical test
	list, err := s.client.ObjectNames(s.container, &swift.ObjectsOpts{Prefix: key})
	if err != nil {
		return nil, err
	}
	if 0 == len(list) {
		return nil, nil
	}
	data, err := s.client.ObjectGetBytes(s.container, key)
	if err == swift.ObjectNotFound {
		return nil, nil
	}
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
func (s *SwiftBackend) Delete(ctx context.Context, key string) error {
	defer metrics.MeasureSince([]string{"swift", "delete"}, time.Now())

	s.permitPool.Acquire()
	defer s.permitPool.Release()

	err := s.client.ObjectDelete(s.container, key)

	if err != nil && err != swift.ObjectNotFound {
		return err
	}

	return nil
}

// List is used to list all the keys under a given
// prefix, up to the next prefix.
func (s *SwiftBackend) List(ctx context.Context, prefix string) ([]string, error) {
	defer metrics.MeasureSince([]string{"swift", "list"}, time.Now())

	s.permitPool.Acquire()
	defer s.permitPool.Release()

	list, err := s.client.ObjectNamesAll(s.container, &swift.ObjectsOpts{Prefix: prefix})
	if nil != err {
		return nil, err
	}

	keys := []string{}
	for _, key := range list {
		key := strings.TrimPrefix(key, prefix)

		if i := strings.Index(key, "/"); i == -1 {
			// Add objects only from the current 'folder'
			keys = append(keys, key)
		} else if i != -1 {
			// Add truncated 'folder' paths
			keys = strutil.AppendIfMissing(keys, key[:i+1])
		}
	}

	sort.Strings(keys)

	return keys, nil
}
