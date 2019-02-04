package etcd

import (
	"context"
	"errors"
	"net/url"
	"os"
	"strings"

	"github.com/coreos/go-semver/semver"
	"github.com/hashicorp/errwrap"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/physical"
	"go.etcd.io/etcd/client"
)

var (
	EtcdSyncConfigError          = errors.New("client setup failed: unable to parse etcd sync field in config")
	EtcdSyncClusterError         = errors.New("client setup failed: unable to sync etcd cluster")
	EtcdMultipleBootstrapError   = errors.New("client setup failed: multiple discovery or bootstrap flags specified, use either \"address\" or \"discovery_srv\"")
	EtcdAddressError             = errors.New("client setup failed: address must be valid URL (ex. 'scheme://host:port')")
	EtcdSemaphoreKeysEmptyError  = errors.New("lock queue is empty")
	EtcdLockHeldError            = errors.New("lock already held")
	EtcdLockNotHeldError         = errors.New("lock not held")
	EtcdSemaphoreKeyRemovedError = errors.New("semaphore key removed before lock acquisition")
	EtcdVersionUnknown           = errors.New("etcd: unknown API version")
)

// NewEtcdBackend constructs a etcd backend using a given machine address.
func NewEtcdBackend(conf map[string]string, logger log.Logger) (physical.Backend, error) {
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
		return nil, EtcdVersionUnknown
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

// Retrieves the config option in order of priority:
//  1. The named environment variable if it exist
//  2. The key in the config map
func getEtcdOption(conf map[string]string, confKey, envVar string) (string, bool) {
	confVal, inConf := conf[confKey]
	envVal, inEnv := os.LookupEnv(envVar)
	if inEnv {
		return envVal, true
	}
	return confVal, inConf
}

func getEtcdEndpoints(conf map[string]string) ([]string, error) {
	address, staticBootstrap := getEtcdOption(conf, "address", "ETCD_ADDR")
	domain, useSrv := getEtcdOption(conf, "discovery_srv", "ETCD_DISCOVERY_SRV")
	if useSrv && staticBootstrap {
		return nil, EtcdMultipleBootstrapError
	}

	if staticBootstrap {
		endpoints := strings.Split(address, Etcd2MachineDelimiter)
		// Verify that the machines are valid URLs
		for _, e := range endpoints {
			u, urlErr := url.Parse(e)
			if urlErr != nil || u.Scheme == "" {
				return nil, EtcdAddressError
			}
		}
		return endpoints, nil
	}

	if useSrv {
		srvName, _ := getEtcdOption(conf, "discovery_srv_name", "ETCD_DISCOVERY_SRV_NAME")
		discoverer := client.NewSRVDiscover()
		endpoints, err := discoverer.Discover(domain, srvName)
		if err != nil {
			return nil, errwrap.Wrapf("failed to discover etcd endpoints through SRV discovery: {{err}}", err)
		}
		return endpoints, nil
	}

	// Set a default endpoints list if no option was set
	return []string{"http://127.0.0.1:2379"}, nil
}
