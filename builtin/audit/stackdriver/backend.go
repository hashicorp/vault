package stackdriver

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"sync/atomic"

	"cloud.google.com/go/logging"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/audit"
	"github.com/hashicorp/vault/sdk/helper/salt"
	"github.com/hashicorp/vault/sdk/helper/useragent"
	"github.com/hashicorp/vault/sdk/logical"
	"google.golang.org/api/option"
	mrpb "google.golang.org/genproto/googleapis/api/monitoredres"
)

func Factory(ctx context.Context, conf *audit.BackendConfig) (audit.Backend, error) {
	// Generic audit logging config
	if conf.SaltConfig == nil {
		return nil, fmt.Errorf("nil salt config")
	}
	if conf.SaltView == nil {
		return nil, fmt.Errorf("nil salt view")
	}

	// Check if hashing of accessor is disabled
	hmacAccessor := true
	if hmacAccessorRaw, ok := conf.Config["hmac_accessor"]; ok {
		value, err := strconv.ParseBool(hmacAccessorRaw)
		if err != nil {
			return nil, err
		}
		hmacAccessor = value
	}

	logRaw := false
	if raw, ok := conf.Config["log_raw"]; ok {
		b, err := strconv.ParseBool(raw)
		if err != nil {
			return nil, err
		}
		logRaw = b
	}

	// Stackdriver specific audit logging config
	//
	// parent is where the logs will be written to in GCP and should be of the form:
	//
	// projects/PROJECT_ID
	// folders/FOLDER_ID
	// billingAccounts/ACCOUNT_ID
	// organizations/ORG_ID
	//
	// The caller must have the IAM Permission logging.logEntries.create (roles/logging.logWriter)
	// at the appropriate level in the resource heirarchy (project, folder, etc.)
	parent, ok := conf.Config["parent"]
	if !ok {
		return nil, fmt.Errorf("parent is required")
	}

	logID, ok := conf.Config["log_id"]
	if !ok {
		return nil, fmt.Errorf("log_id is required")
	}

	// Support writing logs asynchronously. This is less secure since there is no guarantee that
	// the log is written before responding to a request.
	async := false
	if asyncRaw, ok := conf.Config["async"]; ok {
		value, err := strconv.ParseBool(asyncRaw)
		if err != nil {
			return nil, err
		}
		async = value
	}

	client, err := logging.NewClient(ctx, parent, option.WithUserAgent(useragent.String()))
	if err != nil {
		return nil, err
	}
	lg := client.Logger(logID)

	// perform a test log
	err = lg.LogSync(ctx, logging.Entry{
		Resource: &mrpb.MonitoredResource{
			Type: "generic_task",
		},
	})
	if err != nil {
		return nil, errwrap.Wrapf("Unable to log to stackdriver: {{err}}", err)
	}

	b := &Backend{
		client:     client,
		lg:         lg,
		async:      async,
		saltConfig: conf.SaltConfig,
		saltView:   conf.SaltView,
		salt:       new(atomic.Value),
		formatConfig: audit.FormatterConfig{
			Raw:          logRaw,
			HMACAccessor: hmacAccessor,
		},
	}

	// Ensure we are working with the right type by explicitly storing a nil of
	// the right type
	b.salt.Store((*salt.Salt)(nil))

	b.formatter.AuditFormatWriter = &audit.JSONFormatWriter{
		SaltFunc: b.Salt,
	}

	return b, nil
}

// Backend is the audit backend for Stackdriver audit logging.
type Backend struct {
	client *logging.Client
	lg     *logging.Logger
	async  bool

	formatter    audit.AuditFormatter
	formatConfig audit.FormatterConfig

	saltMutex  sync.RWMutex
	salt       *atomic.Value
	saltConfig *salt.Config
	saltView   logical.Storage
}

var _ audit.Backend = (*Backend)(nil)

func (b *Backend) Salt(ctx context.Context) (*salt.Salt, error) {
	s := b.salt.Load().(*salt.Salt)
	if s != nil {
		return s, nil
	}

	b.saltMutex.Lock()
	defer b.saltMutex.Unlock()

	s = b.salt.Load().(*salt.Salt)
	if s != nil {
		return s, nil
	}

	newSalt, err := salt.NewSalt(ctx, b.saltView, b.saltConfig)
	if err != nil {
		b.salt.Store((*salt.Salt)(nil))
		return nil, err
	}

	b.salt.Store(newSalt)
	return newSalt, nil
}

func (b *Backend) GetHash(ctx context.Context, data string) (string, error) {
	salt, err := b.Salt(ctx)
	if err != nil {
		return "", err
	}
	return audit.HashString(salt, data), nil
}

func (b *Backend) LogRequest(ctx context.Context, in *logical.LogInput) error {
	var buf bytes.Buffer
	b.formatter.FormatRequest(ctx, &buf, b.formatConfig, in)
	log := logging.Entry{
		Payload: json.RawMessage(buf.String()),
		Resource: &mrpb.MonitoredResource{
			Type: "generic_task",
		},
	}
	if b.async {
		b.lg.Log(log)
		return nil
	}
	return b.lg.LogSync(ctx, log)
}

func (b *Backend) LogResponse(ctx context.Context, in *logical.LogInput) error {
	var buf bytes.Buffer
	b.formatter.FormatResponse(ctx, &buf, b.formatConfig, in)
	log := logging.Entry{
		Payload: json.RawMessage(buf.String()),
		Resource: &mrpb.MonitoredResource{
			Type: "generic_task",
		},
	}
	if b.async {
		b.lg.Log(log)
		return nil
	}
	return b.lg.LogSync(ctx, log)
}

func (b *Backend) Reload(_ context.Context) error {
	return nil
}

func (b *Backend) Invalidate(_ context.Context) {
	b.saltMutex.Lock()
	defer b.saltMutex.Unlock()
	b.salt.Store((*salt.Salt)(nil))
}
