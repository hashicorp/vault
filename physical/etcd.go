package physical

import (
	"errors"
	"os"

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

	if apiVersion, ok = conf["etcd_api"]; !ok {
		apiVersion = os.Getenv("ETCD_API")
	}
	if apiVersion == "" {
		// TODO: auto discover latest version
		apiVersion = "2"
	}

	// TODO: check etcd server version. Fail if there is a version mismatch

	switch apiVersion {
	case "2", "etcd2", "v2":
		return newEtcd2Backend(conf, logger)
	case "3", "etcd3", "v3":
		return newEtcd3Backend(conf, logger)
	default:
		return nil, EtcdVersionUnknow
	}
}
