package kubernetes

import (
	"context"
	"crypto/sha256"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/leaderelection"
	"k8s.io/client-go/tools/leaderelection/resourcelock"

	"github.com/armon/go-metrics"
	"github.com/hashicorp/errwrap"
	log "github.com/hashicorp/go-hclog"
	uuid "github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/sdk/helper/strutil"
	"github.com/hashicorp/vault/sdk/physical"
)

// Verify KubernetesBackend satisfies the correct interfaces
var _ physical.Backend = (*KubernetesBackend)(nil)

// KubernetesBackend is a physical backend that stores data
// within Kubernetes secrets.
type KubernetesBackend struct {
	namespace  string
	client     *kubernetes.Clientset
	logger     log.Logger
	permitPool *physical.PermitPool
	haEnabled  bool
}

// NewKubernetesBackend constructs a Kubernetes backend using the given API client and
// server address and credential for accessing Kubernetes API.
func NewKubernetesBackend(conf map[string]string, logger log.Logger) (physical.Backend, error) {
	var err error

	kubeconfig := os.Getenv(clientcmd.RecommendedConfigPathEnvVar)
	var config *rest.Config

	if kubeconfig != "" {
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	} else if _, err := os.Stat(clientcmd.RecommendedConfigDir + "/config"); !os.IsNotExist(err) {
		config, err = clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedConfigDir+"/config")
	} else {
		config, err = rest.InClusterConfig()
	}
	if err != nil {
		return nil, err
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	// Get the namespace to perform read/write operations.
	namespace, ok := conf["namespace"]
	if !ok || namespace == "" {
		return nil, fmt.Errorf("missing namespace")
	}

	_, err = client.CoreV1().Namespaces().Get(context.Background(), namespace, metav1.GetOptions{})
	if err != nil && errors.IsNotFound(err) {
		var ns v1.Namespace
		ns.Name = namespace
		_, err = client.CoreV1().Namespaces().Create(context.Background(), &ns, metav1.CreateOptions{})
	}
	if err != nil {
		return nil, errwrap.Wrapf("failed to create namespace: {{err}}", err)
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
	} else {
		maxParInt = physical.DefaultParallelOperations
	}

	// Default value for ha_enabled
	haEnabledStr, ok := conf["ha_enabled"]
	if !ok {
		haEnabledStr = "false"
	}
	haEnabled, err := strconv.ParseBool(haEnabledStr)
	if err != nil {
		return nil, fmt.Errorf("value [%v] of 'ha_enabled' could not be understood", haEnabledStr)
	}

	// TODO
	// config.Burst
	// config.RateLimiter
	// client.Settings().RESTClient().GetRateLimiter().

	// Setup the backend.
	k := &KubernetesBackend{
		namespace:  namespace,
		client:     client,
		logger:     logger,
		permitPool: physical.NewPermitPool(maxParInt),
		haEnabled:  haEnabled,
	}

	return k, nil
}

func prefixToLabels(prefix string) map[string]string {
	levels := strings.Split(prefix, "/")

	labels := map[string]string{}
	labels["prefix-0"] = ""
	for i, level := range levels {
		if level == "" {
			continue
		}
		labels[fmt.Sprintf("prefix-%d", i+1)] = level
	}

	return labels
}

// Put is used to insert or update an entry.
func (k *KubernetesBackend) Put(ctx context.Context, entry *physical.Entry) error {
	defer metrics.MeasureSince([]string{"kubernetes", "put"}, time.Now())

	k.permitPool.Acquire()
	defer k.permitPool.Release()

	secretName := toSHASum(entry.Key)

	lastSlash := strings.LastIndex(entry.Key, "/")

	prefix := entry.Key[:lastSlash+1]
	// key := entry.Key[lastSlash+1:]

	labels := prefixToLabels(prefix)

	secret := v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: k.namespace,
			Labels:    labels,
		},
		Data: map[string][]byte{
			"key":   []byte(entry.Key),
			"value": entry.Value,
		},
	}

	_, err := k.client.CoreV1().Secrets(k.namespace).Create(ctx, &secret, metav1.CreateOptions{})

	if err != nil && apierrors.IsAlreadyExists(err) {
		_, err = k.client.CoreV1().Secrets(k.namespace).Update(ctx, &secret, metav1.UpdateOptions{})
	}

	return err
}

func toSHASum(s string) string {
	shaSum := sha256.Sum224([]byte(s))
	return fmt.Sprintf("vault-%x", shaSum)
}

// Get is used to fetch and entry.
func (k *KubernetesBackend) Get(ctx context.Context, key string) (*physical.Entry, error) {
	defer metrics.MeasureSince([]string{"kubernetes", "get"}, time.Now())

	k.permitPool.Acquire()
	defer k.permitPool.Release()

	secretName := toSHASum(key)

	secret, err := k.client.CoreV1().Secrets(k.namespace).Get(ctx, secretName, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}

	ent := &physical.Entry{
		Key:   key,
		Value: secret.Data["value"],
	}
	return ent, nil
}

// Delete is used to permanently delete an entry
func (k *KubernetesBackend) Delete(ctx context.Context, key string) error {
	defer metrics.MeasureSince([]string{"kubernetes", "delete"}, time.Now())

	k.permitPool.Acquire()
	defer k.permitPool.Release()

	key = toSHASum(key)

	err := k.client.CoreV1().Secrets(k.namespace).Delete(ctx, key, metav1.DeleteOptions{})
	if err != nil && !apierrors.IsNotFound(err) {
		return err
	}
	return nil
}

// List is used to list all the keys under a given
// prefix, up to the next prefix.
func (k *KubernetesBackend) List(ctx context.Context, prefix string) ([]string, error) {
	defer metrics.MeasureSince([]string{"kubernetes", "list"}, time.Now())

	k.permitPool.Acquire()
	defer k.permitPool.Release()

	levelLabels := prefixToLabels(prefix)
	selector := labels.SelectorFromSet(levelLabels)

	listOptions := metav1.ListOptions{LabelSelector: selector.String(), Limit: 100}

	var keys []string
	for {
		secrets, err := k.client.CoreV1().Secrets(k.namespace).List(ctx, listOptions)
		if err != nil {
			return nil, errwrap.Wrapf("failed to list secrets: {{err}}", err)
		}

		for _, secret := range secrets.Items {
			key := string(secret.Data["key"])
			key = strings.TrimPrefix(key, prefix)

			if i := strings.Index(key, "/"); i == -1 {
				// Add objects only from the current 'folder'
				keys = append(keys, key)
			} else if i != -1 {
				// Add truncated 'folder' paths
				keys = strutil.AppendIfMissing(keys, string(key[:i+1]))
			}
		}

		if secrets.Continue == "" {
			break
		}

		listOptions.Continue = secrets.Continue
	}

	sort.Strings(keys)
	return keys, nil
}

// LockWith is used for mutual exclusion based on the given key.
func (k *KubernetesBackend) LockWith(key, value string) (physical.Lock, error) {
	identity, err := uuid.GenerateUUID()
	if err != nil {
		return nil, err
	}
	l := &KubernetesHALock{
		in:        k,
		key:       key,
		value:     value,
		logger:    k.logger,
		client:    k.client,
		namespace: k.namespace,
		identity:  identity,
	}
	return l, nil
}

func (k *KubernetesBackend) HAEnabled() bool {
	return k.haEnabled
}

type KubernetesHALock struct {
	in            *KubernetesBackend
	client        *kubernetes.Clientset
	key           string
	namespace     string
	value         string
	logger        log.Logger
	leaderElector *leaderelection.LeaderElector
	cancel        context.CancelFunc
	identity      string
	lock          sync.Mutex
	held          bool
}

func (l *KubernetesHALock) Lock(stopCh <-chan struct{}) (<-chan struct{}, error) {
	l.lock.Lock()
	defer l.lock.Unlock()
	if l.held {
		return nil, fmt.Errorf("lock already held")
	}

	rl, err := resourcelock.New(
		resourcelock.LeasesResourceLock,
		l.namespace,
		toSHASum(l.key),
		l.client.CoreV1(),
		l.client.CoordinationV1(),
		resourcelock.ResourceLockConfig{
			Identity: l.identity,
		})
	if err != nil {
		return nil, errwrap.Wrapf("error creating lock: {{err}}", err)
	}

	acquiredChan := make(chan struct{})
	closeChan := make(chan struct{})

	lec := leaderelection.LeaderElectionConfig{
		Lock:          rl,
		LeaseDuration: 15 * time.Second,
		RenewDeadline: 10 * time.Second,
		RetryPeriod:   2 * time.Second,
		Callbacks: leaderelection.LeaderCallbacks{
			OnStartedLeading: func(ctx context.Context) {
				close(acquiredChan)
			},
			OnStoppedLeading: func() {
				close(closeChan)
			},
		},
		Name: l.value,
	}

	le, err := leaderelection.NewLeaderElector(lec)
	if err != nil {
		return nil, err
	}

	l.leaderElector = le

	ctx, cancel := context.WithCancel(context.Background())
	l.cancel = cancel

	go func() {
		le.Run(ctx)
	}()

	select {
	case <-acquiredChan:
		l.held = true
		if stopCh != nil {
			go func() {
				<-stopCh
				cancel()
			}()
		}
		return closeChan, nil
	case <-stopCh:
		cancel()
		return nil, nil
	}
}

func (l *KubernetesHALock) Unlock() error {
	l.lock.Lock()
	defer l.lock.Unlock()
	if !l.held {
		return nil
	}

	l.cancel()
	l.held = false

	return nil
}

func (l *KubernetesHALock) Value() (bool, string, error) {
	return l.leaderElector.IsLeader(), l.value, nil
}
