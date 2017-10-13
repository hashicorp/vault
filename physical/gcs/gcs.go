package gcs

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/errwrap"
	uuid "github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/physical"
	log "github.com/mgutz/logxi/v1"

	"cloud.google.com/go/storage"
	"github.com/armon/go-metrics"
	"golang.org/x/net/context"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

const (
	// GCSLockPrefix is the prefix used to mark GCS records
	// as locks. This prefix causes them not to be returned by
	// List operations.
	GCSLockPrefix = "_"

	// GCSLockTTL The lock TTL matches the default that Consul API uses, 15 seconds.
	GCSLockTTL = 15 * time.Second

	// GCSLockRenewInterval The amount of time to wait between the lock renewals
	GCSLockRenewInterval = 5 * time.Second

	// GCSLockRetryInterval is the amount of time to wait
	// if a lock fails before trying again.
	GCSLockRetryInterval = time.Second

	// GCSWatchRetryMax is the number of times to re-try a
	// failed watch before signaling that leadership is lost.
	GCSWatchRetryMax = 5

	// GCSWatchRetryInterval is the amount of time to wait
	// if a watch fails before trying again.
	GCSWatchRetryInterval = 5 * time.Second
)

// GCSBackend is a physical backend that stores data
// within an Google Cloud Storage bucket.
type GCSBackend struct {
	bucketName string
	client     *storage.Client
	permitPool *physical.PermitPool
	logger     log.Logger

	haEnabled bool
}

// NewGCSBackend constructs a Google Cloud Storage backend using a pre-existing
// bucket. Credentials can be provided to the backend, sourced
// from environment variables or a service account file
func NewGCSBackend(conf map[string]string, logger log.Logger) (physical.Backend, error) {
	bucketName := os.Getenv("GOOGLE_STORAGE_BUCKET")

	if bucketName == "" {
		bucketName = conf["bucket"]
		if bucketName == "" {
			return nil, fmt.Errorf("env var GOOGLE_STORAGE_BUCKET or configuration parameter 'bucket' must be set")
		}
	}

	ctx := context.Background()
	client, err := newGCSClient(ctx, conf, logger)
	if err != nil {
		return nil, errwrap.Wrapf("error establishing strorage client: {{err}}", err)
	}

	// check client connectivity by getting bucket attributes
	_, err = client.Bucket(bucketName).Attrs(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to access bucket '%s': '%v'", bucketName, err)
	}

	maxParStr, ok := conf["max_parallel"]
	var maxParInt int
	if ok {
		maxParInt, err = strconv.Atoi(maxParStr)
		if err != nil {
			return nil, errwrap.Wrapf("failed parsing max_parallel parameter: {{err}}", err)
		}
		if logger.IsWarn() {
			logger.Warn("physical/gcs: max_parallel set", "max_parallel", maxParInt)
		}
	}

	g := GCSBackend{
		bucketName: bucketName,
		client:     client,
		permitPool: physical.NewPermitPool(maxParInt),
		logger:     logger,
	}

	return &g, nil
}

func newGCSClient(ctx context.Context, conf map[string]string, logger log.Logger) (*storage.Client, error) {
	// if credentials_file is configured, try to use it
	// else use application default credentials
	credentialsFile, ok := conf["credentials_file"]
	if ok {
		client, err := storage.NewClient(
			ctx,
			option.WithServiceAccountFile(credentialsFile),
		)

		if err != nil {
			return nil, fmt.Errorf("error with provided credentials: '%v'", err)
		}
		return client, nil
	}

	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, errwrap.Wrapf("error with application default credentials: {{err}}", err)
	}
	return client, nil
}

// Put is used to insert or update an entry
func (g *GCSBackend) Put(entry *physical.Entry) error {
	defer metrics.MeasureSince([]string{"gcs", "put"}, time.Now())

	bucket := g.client.Bucket(g.bucketName)
	writer := bucket.Object(entry.Key).NewWriter(context.Background())

	g.permitPool.Acquire()
	defer g.permitPool.Release()

	defer writer.Close()
	_, err := writer.Write(entry.Value)

	return err
}

// Get is used to fetch an entry
func (g *GCSBackend) Get(key string) (*physical.Entry, error) {
	defer metrics.MeasureSince([]string{"gcs", "get"}, time.Now())

	bucket := g.client.Bucket(g.bucketName)
	reader, err := bucket.Object(key).NewReader(context.Background())

	// return (nil, nil) if object doesn't exist
	if err == storage.ErrObjectNotExist {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("error creating bucket reader: '%v'", err)
	}

	g.permitPool.Acquire()
	defer g.permitPool.Release()

	defer reader.Close()
	value, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("error reading object '%v': '%v'", key, err)
	}

	ent := physical.Entry{
		Key:   key,
		Value: value,
	}

	return &ent, nil
}

// Delete is used to permanently delete an entry
func (g *GCSBackend) Delete(key string) error {
	defer metrics.MeasureSince([]string{"gcs", "delete"}, time.Now())

	bucket := g.client.Bucket(g.bucketName)

	g.permitPool.Acquire()
	defer g.permitPool.Release()

	err := bucket.Object(key).Delete(context.Background())

	// deletion of non existent object is OK
	if err == storage.ErrObjectNotExist {
		return nil
	} else if err != nil {
		return fmt.Errorf("error deleting object '%v': '%v'", key, err)
	}

	return nil
}

// List is used to list all the keys under a given
// prefix, up to the next prefix.
func (g *GCSBackend) List(prefix string) ([]string, error) {
	defer metrics.MeasureSince([]string{"gcs", "list"}, time.Now())

	bucket := g.client.Bucket(g.bucketName)

	objects_it := bucket.Objects(
		context.Background(),
		&storage.Query{
			Prefix:    prefix,
			Delimiter: "/",
			Versions:  false,
		})

	keys := []string{}

	g.permitPool.Acquire()
	defer g.permitPool.Release()

	for {
		objAttrs, err := objects_it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error listing bucket '%v': '%v'", g.bucketName, err)
		}

		path := ""
		if objAttrs.Prefix != "" {
			// "subdirectory"
			path = objAttrs.Prefix
		} else {
			// file
			path = objAttrs.Name
		}

		// get relative file/dir just like "basename"
		key := strings.TrimPrefix(path, prefix)
		keys = append(keys, key)
	}

	sort.Strings(keys)

	return keys, nil
}

// GCSLock implements a lock using an GCS client.
type GCSLock struct {
	backend    *GCSBackend
	value, key string
	identity   string
	held       bool
	lock       sync.Mutex
	recovery   bool
	// Allow modifying the Lock durations for ease of unit testing.
	renewInterval      time.Duration
	ttl                time.Duration
	watchRetryInterval time.Duration
}

type GCSLockRecord struct {
	Path     string
	Key      string
	Value    []byte
	Identity []byte
	Expires  int64
}

// LockWith is used for mutual exclusion based on the given key.
func (g *GCSBackend) LockWith(key, value string) (physical.Lock, error) {
	identity, err := uuid.GenerateUUID()
	log.Warn("Creating new identity: ", identity)
	if err != nil {
		return nil, err
	}
	return &GCSLock{
		backend:  g,
		key:      key,
		value:    value,
		identity: identity,
		// recovery:           g.recovery,
		renewInterval:      GCSLockRenewInterval,
		ttl:                GCSLockTTL,
		watchRetryInterval: GCSWatchRetryInterval,
	}, nil
}

func (g *GCSBackend) HAEnabled() bool {
	return true
}

// Lock tries to acquire the lock by repeatedly trying to create
// a record in the GCS. It will block until either the
// stop channel is closed or the lock could be acquired successfully.
// The returned channel will be closed once the lock is deleted or
// changed in the GCS.
func (l *GCSLock) Lock(stopCh <-chan struct{}) (doneCh <-chan struct{}, retErr error) {
	l.lock.Lock()
	defer l.lock.Unlock()
	// log.Warn("Attempting to lock")
	if l.held {
		return nil, fmt.Errorf("lock already held")
	}

	done := make(chan struct{})
	// close done channel even in case of error
	defer func() {
		if retErr != nil {
			close(done)
		}
	}()

	var (
		stop    = make(chan struct{})
		success = make(chan struct{})
		errors  = make(chan error)
		leader  = make(chan struct{})
	)
	// log.Warn("Attempting to lock async")
	// try to acquire the lock asynchronously
	go l.tryToLock(stop, success, errors)

	select {
	case <-success:
		l.held = true
		// after acquiring it successfully, we must renew the lock periodically,
		// and watch the lock in order to close the leader channel
		// once it is lost.
		go l.periodicallyRenewLock(leader)
		go l.watch(leader)
	case retErr = <-errors:
		close(stop)
		return nil, retErr
	case <-stopCh:
		close(stop)
		return nil, nil
	}

	return leader, retErr
}

// Unlock releases the lock by deleting the lock record from the
// GCS.
func (l *GCSLock) Unlock() error {
	l.lock.Lock()
	defer l.lock.Unlock()
	if !l.held {
		return nil
	}

	l.held = false
	if err := l.backend.Delete(l.key); err != nil {
		return err
	}
	return nil
}

// Value checks whether or not the lock is held by any instance of GCSLock,
// including this one, and returns the current value.
func (l *GCSLock) Value() (bool, string, error) {
	entry, err := l.backend.Get(l.key)
	if err != nil {
		return false, "", err
	}
	if entry == nil {
		return false, "", nil
	}

	return true, string(entry.Value), nil
}

// tryToLock tries to create a new item in GCS
// every `GCSLockRetryInterval`. As long as the item
// cannot be created (because it already exists), it will
// be retried. If the operation fails due to an error, it
// is sent to the errors channel.
// When the lock could be acquired successfully, the success
// channel is closed.
func (l *GCSLock) tryToLock(stop, success chan struct{}, errors chan error) {
	ticker := time.NewTicker(GCSLockRetryInterval)

	for {
		select {
		case <-stop:
			ticker.Stop()
		case <-ticker.C:
			// log.Warn("Before Write Item")
			err := l.writeItem()
			// Don't report a condition check failure, this means that the lock
			// is already being held.
			log.Warn("After Write Item: Erorr: ", err)
			if err != nil {
				if err.Error() != "ConditionalCheckFailedException" {
					errors <- err
				}
			} else {
				ticker.Stop()
				close(success)
			}
			return
		}
	}
}

func (l *GCSLock) periodicallyRenewLock(done chan struct{}) {
	ticker := time.NewTicker(l.renewInterval)
	for {
		select {
		case <-ticker.C:
			l.writeItem()
			log.Warn("TRYING TO RENEW LOCK")
		case <-done:
			ticker.Stop()
			return
		}
	}
}

// Attempts to put/update the gcs item using condition expressions to
// evaluate the TTL.
func (l *GCSLock) writeItem() error {
	now := time.Now()
	// // If both key and path already exist, we can only write if
	// // A. identity is equal to our identity (or the identity doesn't exist)
	// // or
	// // B. The ttl on the item is <= to the current time
	// _, err := l.backend.Put()
	defer metrics.MeasureSince([]string{"gcs", "get"}, time.Now())
	var err error

	canwriteExpired := false
	newObj := true
	canwriteIdentity := false
	bucket := l.backend.client.Bucket(l.backend.bucketName)

	exists, err := l.backend.Get(l.key)
	if err != nil {
		return errors.New("ERROR")
	}

	if exists != nil {
		obj := bucket.Object(l.key)
		attrs, _ := obj.Attrs(context.Background())

		newObj = false
		if identity, ok := attrs.Metadata["identity"]; ok && identity == l.identity {
			canwriteIdentity = true
		}
		// else {
		// 	log.Warn("NO IDENTITY MATCH")
		// 	return errors.New("ConditionalCheckFailedException")
		// }

		if ts, ok := attrs.Metadata["expires"]; ok {
			i, err := strconv.ParseInt(ts, 10, 64)
			if err != nil {
				return errors.New("couldn't parse unix timestamp")
			}
			expires := time.Unix(i, 0)

			// log.Warn("Expires: ", expires.Unix())
			// log.Warn("Now: ", now.Unix())

			log.Warn("Now: ", strconv.FormatInt(now.Unix(), 10))
			log.Warn("Expires: ", strconv.FormatInt(expires.Unix(), 10))

			canwriteExpired = true
			if expires.Unix() > now.Unix() {
				log.Warn("NOT EXPIRED")
				return errors.New("ConditionalCheckFailedException")
			}
		}
	}

	if newObj || (canwriteExpired || canwriteIdentity) {
		obj := bucket.Object(l.key)
		rw := obj.NewWriter(context.Background())
		defer rw.Close()

		log.Warn("\n!!!Writing Object!!!: ", l.key, string(l.identity), "Value: ", string(l.value), "Expires: ", strconv.FormatInt(now.Add(l.ttl).Unix(), 10))

		// update the expire time
		rw.ObjectAttrs.Metadata = map[string]string{}
		rw.ObjectAttrs.Metadata["expires"] = strconv.FormatInt(now.Add(l.ttl).Unix(), 10)
		// rw.ObjectAttrs.Metadata["now"] = strconv.FormatInt(now.UnixNano(), 10)
		rw.ObjectAttrs.Metadata["identity"] = l.identity

		_, err = rw.Write([]byte(l.value))
	}

	log.Warn("WRITING OBJECT    NEWOBJ: ", fmt.Sprintf("%t", newObj))
	log.Warn("WRITING OBJECT    EXPIRED: ", fmt.Sprintf("%t", canwriteExpired))
	log.Warn("WRITING OBJECT    IDENTITY MATCH", fmt.Sprintf("%t", canwriteIdentity))

	return err
}

// watch checks whether the lock has changed in the
// GCS and closes the leader channel if so.
// The interval is set by `GCSWatchRetryInterval`.
// If an error occurs during the check, watch will retry
// the operation for `GCSWatchRetryMax` times and
// close the leader channel if it can't succeed.
func (l *GCSLock) watch(lost chan struct{}) {
	retries := GCSWatchRetryMax

	ticker := time.NewTicker(l.watchRetryInterval)
WatchLoop:
	for {
		select {
		case <-ticker.C:
			bucket := l.backend.client.Bucket(l.backend.bucketName)
			obj := bucket.Object(l.key)
			attrs, err := obj.Attrs(context.Background())

			if err != nil {
				log.Warn("watch. Retrying. Retries left: ", retries)
				retries--
				if retries == 0 {
					log.Warn("Breaking watch loop")
					break WatchLoop
				}
				continue
			}

			if attrs == nil {
				break WatchLoop
			}

			record := &GCSLockRecord{}
			if identity, ok := attrs.Metadata["identity"]; ok {
				record.Identity, err = uuid.ParseUUID(identity)
				if err != nil {
					break WatchLoop
				}
			}

			reader, err := obj.NewReader(context.Background())
			if err != nil {
				break WatchLoop
			}
			defer reader.Close()
			value, err := ioutil.ReadAll(reader)
			if err != nil {
				break WatchLoop
			}

			expires, err := strconv.ParseInt(attrs.Metadata["expires"], 10, 64)
			record.Expires = expires
			record.Value = value

			compareId, err := uuid.FormatUUID(record.Identity)
			if err != nil {
				log.Warn("Coudln't format identity")
				break WatchLoop
			}

			if err != nil || compareId != string(l.identity) {
				log.Warn("Identity doesn't match.", "Expected:", string(l.identity), "Got: ", compareId)
				break WatchLoop
			}
		}
	}

	close(lost)
}
