package gcs

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/physical"
	log "github.com/mgutz/logxi/v1"

	"cloud.google.com/go/storage"
	"github.com/armon/go-metrics"
	"golang.org/x/net/context"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

// GCSBackend is a physical backend that stores data
// within an Google Cloud Storage bucket.
type GCSBackend struct {
	bucketName string
	client     *storage.Client
	permitPool *physical.PermitPool
	logger     log.Logger
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
		if logger.IsDebug() {
			logger.Debug("physical/gcs: max_parallel set", "max_parallel", maxParInt)
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
