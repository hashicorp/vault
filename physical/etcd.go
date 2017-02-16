package physical

import (
	"context"
	"errors"
	"os"

	"github.com/coreos/etcd/client"
	"github.com/coreos/go-semver/semver"
	log "github.com/mgutz/logxi/v1"
)

var (
	EtcdSyncConfigError          = errors.New("client setup failed: unable to parse etcd sync field in config")
	EtcdSyncClusterError         = errors.New("client setup failed: unable to sync etcd cluster")
	EtcdAddressError             = errors.New("client setup failed: address must be valid URL (ex. 'scheme://host:port')")
	EtcdSemaphoreKeysEmptyError  = errors.New("lock queue is empty")
	EtcdLockHeldError            = errors.New("lock already held")
	EtcdLockNotHeldError         = errors.New("lock not held")
	EtcdSemaphoreKeyRemovedError = errors.New("semaphore key removed before lock aquisition")
	EtcdVersionUnknow            = errors.New("etcd: unknown API version")
)

// newEtcdBackend constructs a etcd backend using a given machine address.
func newEtcdBackend(conf map[string]string, logger log.Logger) (Backend, error) {
	var (
		apiVersion string
		ok         bool
	)

	// v2 client can talk to both etcd2 and etcd3 thought API v2
	c, err := newEtcdV2Client(conf)
	if err != nil {
		return nil, errors.New("failed to create etcd client: " + err.Error())
	}

	remoteAPIVersion, err := getEtcdAPIVersion(c)
	if err != nil {
		return nil, errors.New("failed to get etcd API version: " + err.Error())
	}

	if apiVersion, ok = conf["etcd_api"]; !ok {
		apiVersion = os.Getenv("ETCD_API")
	}

	if apiVersion == "" {
		path, ok := conf["path"]
		if !ok {
			path = "/vault"
		}
		kAPI := client.NewKeysAPI(c)

		// keep using v2 if vault data exists in v2 and user does not explicitly
		// ask for v3.
		_, err := kAPI.Get(context.Background(), path, &client.GetOptions{})
		if errorIsMissingKey(err) {
			apiVersion = remoteAPIVersion
		} else if err == nil {
			apiVersion = "2"
		} else {
			return nil, errors.New("failed to check etcd status: " + err.Error())
		}
	}

	switch apiVersion {
	case "2", "etcd2", "v2":
		return newEtcd2Backend(conf, logger)
	case "3", "etcd3", "v3":
		if remoteAPIVersion == "2" {
			return nil, errors.New("etcd3 is required: etcd2 is running")
		}
		return newEtcd3Backend(conf, logger)
	default:
		return nil, EtcdVersionUnknow
	}
}

// getEtcdAPIVersion gets the latest supported API version.
// If etcd cluster version >= 3.1, "3" will be returned.
// Otherwise, "2" will be returned.
func getEtcdAPIVersion(c client.Client) (string, error) {
	v, err := c.GetVersion(context.Background())
	if err != nil {
		return "", err
	}

	sv, err := semver.NewVersion(v.Cluster)
	if err != nil {
		return "", nil
	}

	if sv.LessThan(*semver.Must(semver.NewVersion("3.1.0"))) {
		return "2", nil
	}

	return "3", nil
}
